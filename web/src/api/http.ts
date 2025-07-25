import type { AxiosError, AxiosInstance, AxiosResponse } from 'axios'
import axios from 'axios'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'

const API_BASE_URL = '/api/v1'

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.setAuth('', null)
      await router.push('/login')
      return Promise.reject(error)
    }
    return Promise.reject(error)
  },
)
export {
  api,
  API_BASE_URL,
}
