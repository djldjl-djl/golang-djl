// stores/useAppData.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
interface namespacejson{
    name:string
}
export const useAppData = defineStore('appData', () => {
  const data = ref<namespacejson[]>([])

  const getNamespace = async () => {
    try {
      const res = await axios.get('http://192.168.85.80:8080/namespace')
      data.value = res.data.data
    } catch (error) {
      console.error('获取名称空间失败:', error)
    }
  }

  return { data, getNamespace }
})

export const useUserStore = defineStore('user', {
  state: () => ({
    username: '',
    token: ''
  }),
  actions: {
    setUser(username: string, token: string) {
      this.username = username
      this.token = token
      localStorage.setItem('username', username)
      localStorage.setItem('token', token)
    },
    loadFromStorage() {
      this.username = localStorage.getItem('username') || ''
      this.token = localStorage.getItem('token') || ''
    },
    logout() {
      this.username = ''
      this.token = ''
      localStorage.removeItem('username')
      localStorage.removeItem('token')
    }
  }, persist: true
})
