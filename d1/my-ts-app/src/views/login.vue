<template>
  <div class="login-overlay">
    <div class="login-container">
      <h2>登录系统</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-item">
          <label>用户名：</label>
          <input v-model="username" required />
        </div>
        <div class="form-item">
          <label>密码：</label>
          <input type="password" v-model="password" required />
        </div>
        <div class="form-item">
          <button type="submit">登录</button>
        </div>
        <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useUserStore } from '../stores/useAppData'

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const router = useRouter()

const handleLogin = async () => {
  try {
    const res = await axios.post('http://192.168.85.80:8080/login/login', {
      username: username.value,
      password: password.value
    })

    // 登录成功后存储 token
    const token = res.data.token
    localStorage.setItem('token', token)
    
const userStore = useUserStore()

userStore.setUser(username.value, token)

    // 跳转主页
    router.push('/zhuye')
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || '登录失败，请检查用户名密码'
  }
}
</script>

<style scoped>
.login-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(to right, #74ebd5, #acb6e5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.login-container {
  background-color: #fff;
  padding: 40px 30px;
  border-radius: 20px;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 400px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
}

.form-item {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #666;
}

input {
  width: calc(100% - 30px);
  padding: 10px 12px;
  font-size: 14px;
  border: 1px solid #ccc;
  border-radius: 8px;
  transition: border-color 0.3s;
}

input:focus {
  border-color: #42b983;
  outline: none;
}

button {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  background: #42b983;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s;
}

button:hover {
  background: #36976b;
}

.error {
  margin-top: 10px;
  color: #e74c3c;
  font-size: 14px;
  text-align: center;
}
</style>

