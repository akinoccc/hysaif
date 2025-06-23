import type {
  ApiMethod,
  LoginRequest,
  LoginResponse,
  LogoutResponse,
  WebAuthnCredentialResponse,
  WeWorkAuthRequest,
  WeWorkAuthURLResponse,
} from './types'
import api from './config'

/**
 * 认证相关API
 */
export const authAPI = {
  /**
   * 用户登录
   * @param Email 邮箱
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

  /**
   * 获取企业微信授权URL
   * @param state 状态参数
   * @returns 授权URL响应
   */
  getWeWorkAuthURL: (state?: string): ApiMethod<WeWorkAuthURLResponse> => {
    const params = state ? { state } : {}
    return api.get('/auth/wework/url', { params })
  },

  /**
   * 获取企业微信配置状态
   * @returns 企业微信配置状态
   */
  getWeWorkConfig: (): ApiMethod<{ enabled: boolean }> => {
    return api.get('/auth/wework/config')
  },

  /**
   * 企业微信登录
   * @param code 授权码
   * @param state 状态参数
   * @returns 登录响应
   */
  weWorkLogin: (code: string, state?: string): ApiMethod<LoginResponse> => {
    const loginData: WeWorkAuthRequest = { code, state }
    return api.post('/auth/wework/login', loginData)
  },

  /**
   * WebAuthn 相关 API
   */

  /**
   * 获取用户的 WebAuthn 凭证列表
   * @returns 凭证列表
   */
  getWebAuthnCredentials: (): ApiMethod<{ credentials: WebAuthnCredentialResponse[] }> => {
    return api.get('/users/webauthn/credentials')
  },

  /**
   * 删除 WebAuthn 凭证
   * @param credentialId 凭证ID
   * @returns 删除结果
   */
  deleteWebAuthnCredential: (credentialId: string): ApiMethod<{ message: string }> => {
    return api.delete(`/users/webauthn/credentials/${credentialId}`)
  },

  /**
   * 开始 WebAuthn 注册
   * @param credentialName 凭证名称
   * @returns 注册选项
   */
  beginWebAuthnRegistration: (credentialName: string): ApiMethod<{ options: any }> => {
    return api.post('/users/webauthn/register/begin', { credential_name: credentialName })
  },

  /**
   * 完成 WebAuthn 注册
   * @param response 注册响应
   * @param credentialName 凭证名称
   * @returns 注册结果
   */
  finishWebAuthnRegistration: (response: any, credentialName: string): ApiMethod<{ message: string }> => {
    return api.post('/users/webauthn/register/finish', {
      response,
      credential_name: credentialName,
    })
  },

  /**
   * 开始 Passkey 登录
   * @returns 登录选项
   */
  beginWebAuthnLogin: (): ApiMethod<{ options: any }> => {
    return api.get('/auth/webauthn/login/begin')
  },

  /**
   * 完成 Passkey 登录
   * @param response 登录响应
   * @returns 登录结果
   */
  finishWebAuthnLogin: (response: any): ApiMethod<LoginResponse> => {
    return api.post('/auth/webauthn/login/finish', { response })
  },
}

export default authAPI
