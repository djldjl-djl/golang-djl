<template>
  <div class="container">
    <h2>自定义资源列表</h2>
    <button class="add-btn" @click="showDialog = true">添加资源</button>

    <table class="custom-table">
      <thead>
        <tr>
          <th>名称</th>
          <th>命名空间</th>
          <th>Pod数</th>
          <th>镜像</th>
          <th>镜像拉取策略</th>
          <th>ServerName</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in DjlD1s" :key="item.name + item.namespace">
          <td>{{ item.name }}</td>
          <td>{{ item.namespace }}</td>
          <td>{{ item.spec.size }}</td>
          <td>{{ item.spec.image }}</td>
          <td>{{ item.spec.imagePullPolicy }}</td>
          <td>{{ item.spec.serverName }}</td>
          <td><button class="delete-btn" @click="deleteDjlD1(item)">删除</button></td>
        </tr>
      </tbody>
    </table>

    <!-- 添加弹窗 -->
    <div v-if="showDialog" class="dialog" @click.self="showDialog = false">
      <div class="dialog-content" @click.stop>
        <h3>添加自定义资源</h3>

        <div class="form-row">
          <label for="name">名称:</label>
          <input id="name" v-model="newDjlD1.name" />
        </div>

        <div class="form-row">
          <label for="namespace">命名空间:</label>
          <input id="namespace" v-model="newDjlD1.namespace" />
        </div>

        <fieldset class="fieldset">
          <legend>Labels</legend>
          <div v-for="(_, key) in newDjlD1.labels" :key="key" class="label-row">
            <input
              v-model="labelKeys[key]"
              placeholder="Key"
              @input="updateLabelKey(key, labelKeys[key])"
              class="label-input"
            />
            <input v-model="newDjlD1.labels[key]" placeholder="Value" class="label-input" />
            <button @click="removeLabel(key)" class="remove-btn">×</button>
          </div>
          <button @click="addLabel" class="add-label-btn">+ 添加标签</button>
        </fieldset>

        <div class="form-row">
          <label for="size">Pod数:</label>
          <input id="size" type="number" v-model.number="newDjlD1.spec.size" />
        </div>

        <div class="form-row">
          <label for="image">镜像:</label>
          <input id="image" v-model="newDjlD1.spec.image" />
        </div>

        <div class="form-row">
          <label for="serverName">ServerName:</label>
          <input id="serverName" v-model="newDjlD1.spec.serverName" />
        </div>

        <div class="form-row">
          <label for="imagePullPolicy">镜像拉取策略:</label>
          <select id="imagePullPolicy" v-model="newDjlD1.spec.imagePullPolicy">
            <option>Always</option>
            <option>IfNotPresent</option>
            <option>Never</option>
          </select>
        </div>

        <fieldset class="fieldset">
          <legend>Ports</legend>
          <div v-for="(port, index) in newDjlD1.spec.ports" :key="index" class="port-row">
            <input v-model="port.name" placeholder="name" class="port-input" />
            <input v-model="port.protocol" placeholder="protocol" class="port-input" />
            <input type="number" v-model.number="port.port" placeholder="port" class="port-input" />
            <input type="number" v-model.number="port.targetPort" placeholder="targetPort" class="port-input" />
            <button @click="removePort(index)" class="remove-btn">×</button>
          </div>
          <button @click="addPort" class="add-port-btn">+ 添加端口</button>
        </fieldset>

        <div class="dialog-actions">
          <button @click="createDjlD1">创建</button>
          <button @click="showDialog = false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import axios from 'axios'

interface Port {
  name: string
  protocol: string
  port: number
  targetPort: number
}

interface DjlD1Spec {
  size: number
  image: string
  imagePullPolicy: string
  ports: Port[]
  serverName: string
}

interface DjlD1 {
  name: string
  namespace: string
  labels: Record<string, string>
  spec: DjlD1Spec
}

const DjlD1s = ref<DjlD1[]>([])
const showDialog = ref(false)

const newDjlD1 = reactive<DjlD1>({
  name: '',
  namespace: '',
  labels: { app: '', env: '' },
  spec: {
    size: 1,
    image: '',
    imagePullPolicy: 'IfNotPresent',
    ports: [],
    serverName: '',
  }
})

// labels的key临时储存
const labelKeys = reactive<Record<string, string>>({
  app: 'app',
  env: 'env',
})

const addLabel = () => {
  const newKey = 'key' + Date.now()
  newDjlD1.labels[newKey] = ''
  labelKeys[newKey] = newKey
}

const removeLabel = (key: string) => {
  delete newDjlD1.labels[key]
  delete labelKeys[key]
}

const updateLabelKey = (oldKey: string, newKey: string) => {
  if (newKey === oldKey) return
  if (newKey in newDjlD1.labels) return // 防止重复
  newDjlD1.labels[newKey] = newDjlD1.labels[oldKey]
  delete newDjlD1.labels[oldKey]
  labelKeys[newKey] = newKey
  delete labelKeys[oldKey]
}

const addPort = () => {
  newDjlD1.spec.ports.push({
    name: '',
    protocol: 'TCP',
    port: 80,
    targetPort: 80,
  })
}

const removePort = (index: number) => {
  newDjlD1.spec.ports.splice(index, 1)
}

const getDjlD1s = async () => {
  try {
    const res = await axios.get('http://192.168.85.80:8080/ddapp')
    DjlD1s.value = res.data.data
  } catch (err) {
    console.error('获取资源失败:', err)
  }
}

const createDjlD1 = async () => {
  const payload = {
    name: newDjlD1.name,
    namespace: newDjlD1.namespace,
    labels: newDjlD1.labels,
    spce: newDjlD1.spec, // 注意拼写
  }
  try {
    await axios.post('http://192.168.85.80:8080/ddapp', payload)
    alert('创建成功')
    showDialog.value = false
    await getDjlD1s()
    // 重置表单
    newDjlD1.name = ''
    newDjlD1.namespace = ''
    newDjlD1.labels = { app: '', env: '' }
    Object.assign(labelKeys, { app: 'app', env: 'env' })
    newDjlD1.spec = {
      size: 1,
      image: '',
      imagePullPolicy: 'IfNotPresent',
      ports: [],
      serverName: '',
    }
  } catch (err) {
    console.error('创建失败:', err)
    alert('创建失败，请检查控制台')
  }
}

const deleteDjlD1 = async (item: DjlD1) => {
  try {
    await axios.delete(`http://192.168.85.80:8080/ddapp/${item.namespace}/${item.name}`)
    await getDjlD1s()
  } catch (err) {
    console.error('删除失败:', err)
    alert('删除失败，请检查控制台')
  }
}

onMounted(getDjlD1s)
</script>

<style scoped>
.container {
  padding: 20px;
  font-family: Arial, sans-serif;
}
.custom-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}
.custom-table th,
.custom-table td {
  border: 1px solid #ccc;
  padding: 8px 12px;
  text-align: center;
}
.custom-table th {
  background-color: #f5f5f5;
}
.add-btn {
  margin-bottom: 10px;
  padding: 6px 12px;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}
.add-btn:hover {
  background-color: #45a049;
}
.delete-btn {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 3px;
}
.delete-btn:hover {
  background-color: #c0392b;
}

.dialog {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  overflow: hidden;
}

.dialog-content {
  background: white;
  padding: 24px 36px;
  width: 900px;
  border-radius: 10px;
  max-height: 90vh;
  overflow-y: auto;
  box-sizing: border-box;
}

.dialog-content h3 {
  margin-bottom: 24px;
  text-align: center;
  font-size: 20px;
}

/* 统一form排版 */
.form-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.form-row label {
  width: 120px;
  text-align: right;
  font-weight: 600;
}

.form-row input,
.form-row select {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

/* 标签区域 */
.fieldset {
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 12px 16px;
  margin-bottom: 20px;
}

.fieldset legend {
  font-weight: 700;
  padding: 0 6px;
}

.label-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.label-input {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.port-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr) 40px;
  gap: 8px;
  margin-bottom: 10px;
}

.port-input {
  padding: 6px 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.add-label-btn,
.add-port-btn,
.remove-btn {
  background-color: #3498db;
  color: white;
  border: none;
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.add-label-btn:hover,
.add-port-btn:hover {
  background-color: #2980b9;
}

.remove-btn {
  background-color: #e74c3c;
  padding: 0;
  width: 30px;
  height: 30px;
  line-height: 30px;
  text-align: center;
}

.remove-btn:hover {
  background-color: #c0392b;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 14px;
  margin-top: 24px;
}

.dialog-actions button {
  padding: 8px 20px;
  font-weight: 600;
  border-radius: 5px;
  border: none;
  cursor: pointer;
}

.dialog-actions button:first-child {
  background-color: #4caf50;
  color: white;
}

.dialog-actions button:first-child:hover {
  background-color: #45a049;
}

.dialog-actions button:last-child {
  background-color: #aaa;
  color: white;
}

.dialog-actions button:last-child:hover {
  background-color: #888;
}
</style>
