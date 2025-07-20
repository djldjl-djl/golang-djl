<script setup lang="ts">
import router from '../router'

const goToZhuye = () => {
  router.push('/zhuye')
}
const goToBushu = () => {
  router.push('/bushu')
}
const goToPodList = (i: any) => {
  router.push('/podlist?namespace='+i)
}
function openliebiao(){
  xuanting.value=true
}
function outliebiao(){
  xuanting.value=false
}
import { onMounted, ref } from 'vue'

const xuanting = ref(false)

import { useAppData } from '../stores/useAppData'
import { storeToRefs } from 'pinia'

const store = useAppData()
const { data } = storeToRefs(store) 
onMounted(async() => {
   await store.getNamespace()
  console.log(data.value)
})

</script>

<template>
  <div class="nav" @click="goToZhuye">集群概况</div>
  <div class="nav" @click="goToBushu">部署应用</div>
  <div class="nav"  @mouseenter="openliebiao" @mouseleave="outliebiao">
    <div class="xuan">pod列表</div>
    <transition name="expand-fade"><!-- 淡入淡出动画 -->
      <div class="beixuan" v-if="xuanting">
        <div class="xiaoxuan" @click="goToPodList(i.name)" v-for="i in data">{{i.name}}</div>
      </div>
    </transition>
  </div>
  <div class="nav" @click="goToPodList">未知</div>
</template>

<style scoped>
/* 淡入淡出动画 */
.fade-slide-enter-active, .fade-slide-leave-active {
  transition: all 0.3s ease;
}
.fade-slide-enter-from, .fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
.fade-slide-enter-to, .fade-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
}


/* 高度展开 + 淡入淡出动画 */
.expand-fade-enter-active,
.expand-fade-leave-active {
  transition: all 0.5s ease;
  overflow: hidden;
  max-height: 300px; /* 展开时的最大高度，根据内容估算设置 */
}

.expand-fade-enter-from,
.expand-fade-leave-to {
  opacity: 0;
  max-height: 0;
}

.expand-fade-enter-to,
.expand-fade-leave-from {
  opacity: 1;
  max-height: 300px; /* 同样保持一致 */
}


.nav{
  font-size: 20px;
  color: #000;
  text-decoration: none;
  font-weight: bold;
  padding: 10px;
  border: 1px solid #4d4848;
  cursor: pointer;
  transition: background-color 0.3s ease;
}
.nav:hover{
  background-color: rgb(131, 41, 131);
}
.beixuan{
background-color: rgb(60, 87, 110);
padding: 20px;
}
.xiaoxuan{
  padding: 5px;
  border-bottom: 1px solid #ccc;
  transition: background-color 0.3s ease;
}
.xiaoxuan:hover{
  background-color: rgb(131, 41, 131);
}
</style>
