<script setup lang="ts">
import Header from './components/header.vue'
import Left from './components/left.vue'
import Main from './components/main.vue'
</script>

<template>
  <div class="app">
    <div id="header" class="head">
      <Header msg="这是参数" />
    </div>
    <div class="ti">
      <div id="left" class="left">
        <Left />
      </div>
      <div id="main" class="main">
        <Main />
      </div>
    </div>
  </div>
</template>

<style scoped>
.app{
  height: 100%;
}
.ti{
  display: flex;
  height: 85%;
  width: 100%;
  overflow: hidden;/*防止整个页面滚动*/
}
.head{
  width: 100%;
  height: 15%;
  background-color: rgb(134, 135, 134);
}
.left{
  width: 20%;
  background-color: rgb(203, 206, 210);
}
.main{
  overflow-y: auto; /*添加垂直滚动条*/
  width: 80%;
  background-color: rgb(196, 199, 197);
}
</style>
