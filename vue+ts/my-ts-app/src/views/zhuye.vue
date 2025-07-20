<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as echarts from 'echarts'
import axios from 'axios'
import { useAppData } from '../stores/useAppData'
import { storeToRefs } from 'pinia'

const store = useAppData()
const { data } = storeToRefs(store) 

interface NodeInfo {
  name: string
  status: any
  fuzai: number
}
interface PodsInfo {
  name: string
  stytus: boolean
  namespace: string
}
const chartRef = ref<HTMLDivElement | null>(null)
const chartRef1 = ref<HTMLDivElement | null>(null)
const chartRef2 = ref<HTMLDivElement | null>(null)
const nodes = ref<NodeInfo[]>([])
const pods =ref<PodsInfo[]>([])
const getNodes = async () => {
  try {
    const res = await axios.get('http://192.168.85.80:8080/node/')
    nodes.value = res.data.data
  } catch (error) {
    console.error('获取节点失败:', error)
  }
}
const getPods = async () => {
  try {
    const res = await axios.get('http://192.168.85.80:8080/podes')
    pods.value = res.data.data
  } catch (error) {
    console.error('获取pod失败:', error)
  }
}

onMounted(async() => {
  await store.getNamespace()
  await getNodes()
  await getPods()
  let podnamespace:string[]=[]
  let podready=ref<number[]>(new Array(data.value.length).fill(0))
  let podnotready=ref<number[]>(new Array(data.value.length).fill(0))
  for (let i=0;i<data.value.length;i++){
    podnamespace.push(data.value[i].name)
  }
  for (let i=0;i<pods.value.length;i++){
    let index = podnamespace.indexOf(pods.value[i].namespace)
    console.log(pods.value[i].stytus)
    if (pods.value[i].stytus){ 
      podready.value[index]++
    }else{
      podnotready.value[index]++
    }
  }
  
  console.log(podnamespace)
  const chart2 =echarts.init(chartRef2.value!)
  chart2.setOption(
    {
  xAxis: {
    data: podnamespace,
    axisLabel:{
      interval: 0 
    }
  },
  yAxis: {},
  series: [
    {
      data: podready.value,
      type: 'bar',
      stack: 'x'
    },
    {
      data: podnotready.value,
      type: 'bar',
      stack: 'x'
    }
  ]
}
  )
  let nodenames:string[] =[]
   for(let i=0;i<nodes.value.length;i++){
    nodenames.push(nodes.value[i].name)
  }
  let nodeneicun:number[] =[]
   for(let i=0;i<nodes.value.length;i++){
    nodeneicun.push(nodes.value[i].fuzai/1024/1024/1024)
  }
  const chart1 =echarts.init(chartRef1.value!)
  chart1.setOption(
    {
    xAxis: {
    data: nodenames
  },
  yAxis: {},
  series: [
    {
      type: 'bar',
      data: nodeneicun
    }
  ] 
    }
  )
  let Ready=0
  let NotReady=0
  for(let i=0;i<nodes.value.length;i++){
    if (nodes.value[i].status){
        Ready++
    }else{
        NotReady++
    }
  }
  const chart = echarts.init(chartRef.value!)
  chart.setOption({
    title: {
      text: Ready+NotReady,
      left: 'center',
      top: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: Ready, name: 'Ready' },
          { value: NotReady, name: 'NotReady' }
        ]
      }
    ]
  })
})
</script>

<template>
  <div class="node">
    <div class="node1">
      <p class="miaoshu">节点状态:</p>
      <div ref="chartRef" style="width: 100%; height: 300px;"></div>
    </div>
    <div class="node1">
      <p class="miaoshu">节点内存负载:</p>
      <div ref="chartRef1" style="width: 100%; height: 300px;"></div>
    </div>
  </div>
  <div class="pod">
    <p class="miaoshu">pod数量:</p>
    <div ref="chartRef2" style="width: 100%; height: 300px;"></div>
  </div>
</template>

<style scoped>
.node {
  display: flex;
  justify-content: space-between;
}
.node1 {
    width: 40%;
    background-color: white;
    margin-top: 20px;
    margin-bottom: 20px;
    border-radius: 12px;
    margin-left: 5%;
    margin-right: 5%;
    
}
.pod{
    width: 90%;
    background-color: white;
    margin: auto;
    margin-top: 20px;
    margin-bottom: 20px;
    border-radius: 12px;
   
}
.miaoshu{
    font-size: 20px;
    padding-top: 10px;
    padding-left: 20px;
}
</style>
