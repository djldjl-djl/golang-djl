<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const namespace = ref<string>('')

interface PodsInfo {
  name: string
  status: boolean
  node_name: string
  cname: string
}

const pods = ref<PodsInfo[]>([])
const showModal = ref(false)
const modalMessage = ref('')
const isDeleting = ref(false)

// 新增日志弹窗相关变量
const showLogModal = ref(false)
const logContent = ref('')
const currentLogPodName = ref('')
let logWs: WebSocket | null = null

const getPods = async (namespace1: string) => {
  try {
    const res = await axios.get('http://192.168.85.80:8080/k8s/select/' + namespace1)
    pods.value = res.data.pods
  } catch (error) {
    console.error('获取节点失败:', error)
  }
}

onMounted(async () => {
  namespace.value = (route.query.namespace as string) || ''
  console.log('初始化 namespace:', namespace.value)
  await getPods(namespace.value)
})

watch(() => route.query.namespace, async (newNamespace) => {
  namespace.value = (newNamespace as string) || ''
  await getPods(namespace.value)
})

const waitForPodDeletion = async (name: string, interval = 1000, maxTries = 10) => {
  for (let i = 0; i < maxTries; i++) {
    await new Promise(resolve => setTimeout(resolve, interval))
    await getPods(namespace.value)
    const exists = pods.value.some(p => p.name === name)
    if (!exists) {
      console.log(`Pod ${name} 已删除`)
      return true
    }
    console.log(`第 ${i + 1} 次轮询，Pod 仍存在...`)
  }
  return false
}

const handleDelete = async (name: string) => {
  console.log('删除 Pod：', name)
  showModal.value = true
  modalMessage.value = `正在删除 Pod：${name} ...`
  isDeleting.value = true

  try {
    await axios.get(`http://192.168.85.80:8080/k8s/delete/${namespace.value}/${name}`)
    console.log('删除请求已发送')
    const deleted = await waitForPodDeletion(name)
    if (deleted) {
      modalMessage.value = `✅ Pod ${name} 删除成功`
    } else {
      modalMessage.value = `⚠️ Pod ${name} 删除超时，请稍后重试`
    }
  } catch (error) {
    console.error('删除失败:', error)
    modalMessage.value = `❌ 删除失败：${error}`
  } finally {
    isDeleting.value = false
    setTimeout(() => {
      showModal.value = false
    }, 2000)
  }
}

const handleLog = (name: string) => {
  openLogModal(name)
}

const handleShell = (name: string,cname:string) => {
  console.log('进入 Shell：', name)
  openShellModal(name,cname)
}

// 日志弹窗相关逻辑
const openLogModal = (podName: string) => {
  currentLogPodName.value = podName
  logContent.value = ''
  showLogModal.value = true

  logWs = new WebSocket(`ws://192.168.85.80:8080/k8s/wslog/${namespace.value}/${podName}`)
  logWs.onmessage = (event) => {
    logContent.value += event.data
    setTimeout(() => {
      const el = document.getElementById('logArea')
      if (el) el.scrollTop = el.scrollHeight
    }, 10)
  }
  logWs.onerror = (err) => {
    console.error('日志 WebSocket 错误', err)
    logContent.value += '\n[错误] WebSocket 连接错误\n'
  }
  logWs.onclose = () => {
    logContent.value += '\n[关闭] WebSocket 连接关闭\n'
  }
}

const closeLogModal = () => {
  showLogModal.value = false
  if (logWs) {
    logWs.close()
    logWs = null
  }
  currentLogPodName.value = ''
  logContent.value = ''
}


import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'

// WebShell 弹窗控制
const showShellModal = ref(false)
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let shellWs: WebSocket | null = null
const currentShellPodName = ref('')

const openShellModal = (podName: string,cName:string) => {
  currentShellPodName.value = podName
  showShellModal.value = true

  // 等 DOM 渲染完成后再初始化终端
  nextTick(() => {
    const terminalContainer = document.getElementById('terminal')
    if (terminalContainer) {
      term = new Terminal({ cursorBlink: true, fontSize: 14 })
      fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      term.open(terminalContainer)
      fitAddon.fit()

      shellWs = new WebSocket(`ws://192.168.85.80:8080/k8s/webshell/${namespace.value}/${podName}/${cName}`)
      term.onData(data => {
        if (shellWs?.readyState === WebSocket.OPEN) {
          shellWs.send(JSON.stringify({ operation: 'stdin', data }))
        }
      })
      shellWs.onmessage = event => {
        const message = JSON.parse(event.data)
        if (message.operation === 'stdout') {
          term?.write(message.data)
        }
        if (message.data.includes('关闭终端')) {
          term?.write('\r\n*** 终端连接已关闭 ***\r\n')
          shellWs?.close()
        }
      }
      shellWs.onerror = () => {
        term?.write(`\r\n*** WebSocket 错误 ***\r\n`)
      }
      shellWs.onclose = () => {
        term?.write('\r\n*** WebSocket 连接已关闭 ***\r\n')
      }
    }
  })
}

const closeShellModal = () => {
  showShellModal.value = false
  if (shellWs) {
    shellWs.close()
    shellWs = null
  }
  if (term) {
    term.dispose()
    term = null
  }
  currentShellPodName.value = ''
}

</script>

<template>
  <div class="pod-list">
    <h2 class="namespace-display">当前命名空间：<span>{{ namespace }}</span></h2>

    <div v-for="pod in pods" :key="pod.name" class="pod-item">
      <div class="pod-info">
        <div><strong>Pod名称：</strong>{{ pod.name }}</div>
        <div><strong>节点：</strong>{{ pod.node_name }}</div>
        <div>
          <strong>状态：</strong>
          <span :class="{'status-running': pod.status, 'status-error': !pod.status}">
            {{ pod.status ? '运行中' : '未就绪' }}
          </span>
        </div>
      </div>
      <div class="pod-actions">
        <button class="btn delete" @click="handleDelete(pod.name)" :disabled="isDeleting">删除</button>
        <button class="btn log" @click="handleLog(pod.name)">日志</button>
        <button class="btn shell" @click="handleShell(pod.name,pod.cname)">Shell</button>
      </div>
    </div>

    <!-- 删除提示弹窗 -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <div class="modal-content">{{ modalMessage }}</div>
      </div>
    </div>

    <!-- 日志弹窗 -->
    <div v-if="showLogModal" class="modal-overlay">
      <div class="modal log-modal">
        <div class="modal-header">
          <span>Pod {{ currentLogPodName }} 日志</span>
          <button @click="closeLogModal" class="close-btn">关闭</button>
        </div>
        <pre id="logArea" class="log-area">{{ logContent }}</pre>
      </div>
    </div>

    <!-- shell弹窗 -->
<div v-if="showShellModal" class="modal-overlay">
  <div class="modal shell-modal">
    <div class="modal-header">
      <span>Pod {{ currentShellPodName }} 终端</span>
      <button @click="closeShellModal" class="close-btn">关闭</button>
    </div>
    <div id="terminal" class="terminal-container"></div>
  </div>
</div>

  </div>
</template>

<style scoped>
.pod-list {
  padding: 16px;
  max-width: 900px;
  margin: auto;
}

.namespace-display {
  font-size: 18px;
  margin-bottom: 20px;
  color: #34495e;
}

.namespace-display span {
  color: #2980b9;
  font-weight: bold;
}

.pod-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  margin-bottom: 12px;
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.pod-info > div {
  margin-bottom: 4px;
  color: #333;
}

.status-running {
  color: green;
  font-weight: bold;
}

.status-error {
  color: red;
  font-weight: bold;
}

.pod-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.btn {
  padding: 6px 10px;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.3s ease;
  color: white;
}

.delete {
  background-color: #e74c3c;
}

.delete:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.delete:hover {
  background-color: #c0392b;
}

.log {
  background-color: #3498db;
}

.log:hover {
  background-color: #2980b9;
}

.shell {
  background-color: #2ecc71;
}

.shell:hover {
  background-color: #27ae60;
}

/* 删除弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal {
  background: white;
  padding: 20px 30px;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.3);
  font-size: 16px;
  color: #2c3e50;
  min-width: 280px;
  text-align: center;
}

/* 日志弹窗样式 */
.log-modal {
  width: 80%;
  max-width: 900px;
  height: 600px;
  display: flex;
  flex-direction: column;
  background: #000;
  color: #0f0;
  border-radius: 10px;
  padding: 10px;
  box-sizing: border-box;
  font-family: monospace;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 6px;
  border-bottom: 1px solid #0f0;
  font-weight: bold;
}

.close-btn {
  background: #c0392b;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  border-radius: 4px;
  padding: 2px 10px;
  font-size: 14px;
}

.log-area {
  text-align: left;
  white-space: pre-wrap;
  font-family: monospace;
  overflow-y: auto;
  margin-top: 8px;
}
.log-area::-webkit-scrollbar {
  width: 8px;               /* 滚动条宽度 */
}

.log-area::-webkit-scrollbar-track {
  background: #222;         /* 滚动条轨道背景 */
  border-radius: 4px;
}

.log-area::-webkit-scrollbar-thumb {
  background: #0f0;         /* 滚动条滑块颜色 */
  border-radius: 4px;
}

.log-area::-webkit-scrollbar-thumb:hover {
  background: #0c0;         /* 滑块悬停颜色 */
}

.shell-modal {
  width: 80%;
  max-width: 900px;
  height: 600px;
  display: flex;
  flex-direction: column;
  background: #000;
  border-radius: 10px;
  padding: 10px;
  box-sizing: border-box;
  font-family: monospace;
  overflow: hidden;
}

.terminal-container {
  text-align: left;
  flex: 1;
  overflow: hidden;
}
.terminal-container::-webkit-scrollbar {
  width: 8px;
}

.terminal-container::-webkit-scrollbar-track {
  background: #222;
  border-radius: 4px;
}

.terminal-container::-webkit-scrollbar-thumb {
  background: #0f0;
  border-radius: 4px;
}

.terminal-container::-webkit-scrollbar-thumb:hover {
  background: #0c0;
}

</style>
