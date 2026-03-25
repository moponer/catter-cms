import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get } from '@/utils/request'
import { getToken } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const currentUser = ref(null)
  const token = ref(getToken())

  const getCurrentUser = async () => {
    if (!token.value) return null
    
    try {
      const res = await get('/api/auth/me')
      currentUser.value = res.data || res
      return currentUser.value
    } catch (error) {
      // 错误已在 request 拦截器中处理
      return null
    }
  }

  const logout = () => {
    currentUser.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    currentUser,
    token,
    getCurrentUser,
    logout
  }
})
