<script setup lang="ts">
import router from '../router'
import { useUserStore } from '../stores/useAppData'
const userStore = useUserStore()

const outlogin=()=>{
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<template>
  <div class="header">
    <div class="tou1"><img src="../../img/login.jpg" height="100%"></div>
    <div class="tou2"><h1>kubernetes管理平台</h1></div>
    <div class="tou3">
      <div>欢迎登录:{{userStore.username}}</div>
      <div>&nbsp;&nbsp;|&nbsp;&nbsp;</div>
      <div class="outlogin" @click="outlogin">退出登录</div>
    </div>
  </div>
</template>

<style scoped>
.header{
  display: flex;
  justify-content:space-between;
  height: 100%;
  width: 100%;
}
.tou1{
  width: 33%;
  height: 100%;
}
.tou2{
  width: 33%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.tou3{
  width: 33%;
  display: flex;
  justify-content: flex-end;
  align-items:center;
  padding: 10px;
}
.outlogin{
  cursor: pointer;
}
</style>
