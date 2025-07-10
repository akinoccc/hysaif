import type {
  ApiMethod,
  LoginRequest,
  LoginResponse,
  LogoutResponse,
} from './types'
import { api } from './http'

/**
 * 认证相关API
 */
export const authAPI = {
  /**
   * 用户登录
   * @param email 邮箱
   * @param password 密码
   * @returns 登录响应，包含token和用户信息
   */
  login: (email: string, password: string): ApiMethod<LoginResponse> => {
    const loginData: LoginRequest = { email, password }
    return api.post('/auth/login', loginData)
  },

  /**
   * 用户登出
   * @returns 登出响应
   */
  logout: (): ApiMethod<LogoutResponse> => {
    return api.post('/auth/logout')
  },
}
