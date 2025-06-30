import type { User } from './types'
import { api } from './http'

export interface WebAuthnCredential {
  id: string
  credential_id: string
  credential_name: string
  created_at: number
  last_used_at?: number
}

// WebAuthn 登录相关
export const webauthnAPI = {
  // 开始登录
  beginLogin: (email: string) => {
    return api.post('/auth/webauthn/login/begin', { email })
  },

  // 完成登录
  finishLogin: (email: string, credential: any) => {
    return api.post<any, { token: string, user: User }>('/auth/webauthn/login/finish', {
      email,
      ...credential,
    })
  },

  // 开始注册
  beginRegistration: (credentialName: string) => {
    return api.post<any, { publicKey: any }>('/users/webauthn/register/begin', {
      credential_name: credentialName,
    })
  },

  // 完成注册
  finishRegistration: (credentialName: string, credential: any) => {
    return api.post('/users/webauthn/register/finish', {
      credential_name: credentialName,
      ...credential,
    })
  },

  // 获取凭证列表
  getCredentials: () => {
    return api.get<any, WebAuthnCredential[]>('/users/webauthn/credentials')
  },

  // 删除凭证
  deleteCredential: (id: string) => {
    return api.delete(`/users/webauthn/credentials/${id}`)
  },
}
