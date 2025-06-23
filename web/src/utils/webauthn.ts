/**
 * WebAuthn 工具函数
 */

/**
 * 检查浏览器是否支持 WebAuthn
 */
export function isWebAuthnSupported(): boolean {
  return !!(
    window.PublicKeyCredential
    && navigator.credentials
    && typeof navigator.credentials.create === 'function'
    && typeof navigator.credentials.get === 'function'
  )
}

/**
 * 将 ArrayBuffer 转换为 Base64 URL 编码字符串
 */
function bufferToBase64URLString(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let str = ''
  for (const byte of bytes) {
    str += String.fromCharCode(byte)
  }
  return btoa(str)
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')
}

/**
 * 将 Base64 URL 编码字符串转换为 ArrayBuffer
 */
function base64URLStringToBuffer(base64URLString: string): ArrayBuffer {
  const base64 = base64URLString
    .replace(/-/g, '+')
    .replace(/_/g, '/')
  const padLength = (4 - (base64.length % 4)) % 4
  const padded = base64 + '='.repeat(padLength)
  const binary = atob(padded)
  const buffer = new ArrayBuffer(binary.length)
  const bytes = new Uint8Array(buffer)
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i)
  }
  return buffer
}

/**
 * 准备注册选项（将服务器返回的选项转换为浏览器可用的格式）
 */
export function prepareRegistrationOptions(options: any): PublicKeyCredentialCreationOptions {
  console.log('options', options)
  return {
    ...options,
    challenge: base64URLStringToBuffer(options.challenge),
    user: {
      ...options.user,
      id: base64URLStringToBuffer(options.user.id),
    },
    excludeCredentials: options.excludeCredentials?.map((cred: any) => ({
      ...cred,
      id: base64URLStringToBuffer(cred.id),
    })),
  }
}

/**
 * 准备登录选项（将服务器返回的选项转换为浏览器可用的格式）
 */
export function prepareLoginOptions(options: any): PublicKeyCredentialRequestOptions {
  return {
    ...options,
    challenge: base64URLStringToBuffer(options.challenge),
    allowCredentials: options.allowCredentials?.map((cred: any) => ({
      ...cred,
      id: base64URLStringToBuffer(cred.id),
    })),
  }
}

/**
 * 准备注册响应（将浏览器返回的凭证转换为服务器可接受的格式）
 */
export function prepareRegistrationResponse(credential: PublicKeyCredential): any {
  const response = credential.response as AuthenticatorAttestationResponse
  const clientExtensionResults = credential.getClientExtensionResults()

  return {
    id: credential.id,
    rawId: bufferToBase64URLString(credential.rawId),
    type: credential.type,
    response: {
      clientDataJSON: bufferToBase64URLString(response.clientDataJSON),
      attestationObject: bufferToBase64URLString(response.attestationObject),
    },
    clientExtensionResults,
  }
}

/**
 * 准备登录响应（将浏览器返回的凭证转换为服务器可接受的格式）
 */
export function prepareLoginResponse(credential: PublicKeyCredential): any {
  const response = credential.response as AuthenticatorAssertionResponse
  const clientExtensionResults = credential.getClientExtensionResults()

  return {
    id: credential.id,
    rawId: bufferToBase64URLString(credential.rawId),
    type: credential.type,
    response: {
      clientDataJSON: bufferToBase64URLString(response.clientDataJSON),
      authenticatorData: bufferToBase64URLString(response.authenticatorData),
      signature: bufferToBase64URLString(response.signature),
      userHandle: response.userHandle ? bufferToBase64URLString(response.userHandle) : null,
    },
    clientExtensionResults,
  }
}

/**
 * 创建 WebAuthn 凭证（注册）
 */
export async function createCredential(options: PublicKeyCredentialCreationOptions): Promise<PublicKeyCredential> {
  try {
    const credential = await navigator.credentials.create({ publicKey: options })
    if (!credential) {
      throw new Error('创建凭证失败')
    }
    return credential as PublicKeyCredential
  }
  catch (error: any) {
    console.log(error)
    if (error.name === 'NotAllowedError') {
      throw new Error('用户取消了操作或操作超时')
    }
    else if (error.name === 'InvalidStateError') {
      throw new Error('该认证器已经注册过了')
    }
    else if (error.name === 'NotSupportedError') {
      throw new Error('浏览器不支持该操作')
    }
    throw error
  }
}

/**
 * 获取 WebAuthn 凭证（登录）
 */
export async function getCredential(options: PublicKeyCredentialRequestOptions): Promise<PublicKeyCredential> {
  try {
    const credential = await navigator.credentials.get({ publicKey: options })
    if (!credential) {
      throw new Error('获取凭证失败')
    }
    return credential as PublicKeyCredential
  }
  catch (error: any) {
    if (error.name === 'NotAllowedError') {
      throw new Error('用户取消了操作或操作超时')
    }
    else if (error.name === 'NotSupportedError') {
      throw new Error('浏览器不支持该操作')
    }
    throw error
  }
}
