import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authAPI, type User } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'super_admin')
  const isSecurityManager = computed(() => user.value?.role === 'sec_mgr')

  // 登录
  const login = async (email: string, password: string) => {
    loading.value = true
    try {
      const response = await authAPI.login(email, password)
      token.value = response.token
      user.value = response.user

      return { success: true }
    }
    catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.error || '登录失败',
      }
    }
    finally {
      loading.value = false
    }
  }

  // 登出
  const logout = async () => {
    try {
      await authAPI.logout()
    }
    catch (error) {
      console.error('Logout error:', error)
    }
    finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }

  const setAuth = (t: string, u: User) => {
    token.value = t
    user.value = u
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    isAdmin,
    isSecurityManager,
    login,
    logout,
    setAuth,
  }
}, {
  persist: true,
})
