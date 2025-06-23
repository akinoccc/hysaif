import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// 格式化日期
export function formatDate(date: number | string | Date) {
  const d = new Date(date)
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 格式化相对时间
export function formatRelativeTime(date: number | string | Date) {
  const now = new Date()
  const target = new Date(date)
  const diff = now.getTime() - target.getTime()

  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) {
    return `${days}天前`
  }
  else if (hours > 0) {
    return `${hours}小时前`
  }
  else if (minutes > 0) {
    return `${minutes}分钟前`
  }
  else {
    return '刚刚'
  }
}

// 复制到剪贴板
export async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    return true
  }
  catch (err) {
    console.error('Failed to copy: ', err)
    return false
  }
}

// 生成随机密码
export function generatePassword(length: number = 16) {
  const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?'
  let password = ''
  for (let i = 0; i < length; i++) {
    password += charset.charAt(Math.floor(Math.random() * charset.length))
  }
  return password
}

// 验证密码强度
export function validatePasswordStrength(password: string) {
  const minLength = 8
  const hasUpperCase = /[A-Z]/.test(password)
  const hasLowerCase = /[a-z]/.test(password)
  const hasNumbers = /\d/.test(password)
  const hasSpecialChar = /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password)

  const score = [
    password.length >= minLength,
    hasUpperCase,
    hasLowerCase,
    hasNumbers,
    hasSpecialChar,
  ].filter(Boolean).length

  if (score < 3)
    return { strength: 'weak', score }
  if (score < 4)
    return { strength: 'medium', score }
  return { strength: 'strong', score }
}

// 获取文件图标
export function getFileIcon(type: string) {
  const icons: Record<string, string> = {
    password: '🔐',
    api_key: '🔑',
    access_key: '🗝️',
    ssh_key: '🔒',
    certificate: '📜',
    token: '🎫',
    custom: '📄',
  }
  return icons[type] || '📄'
}

// 获取角色显示名称
export function getRoleDisplayName(role: string) {
  const roleNames: Record<string, string> = {
    super_admin: '超级管理员',
    sec_mgr: '安全管理员',
    dev: '开发人员',
    auditor: '审计员',
  }
  return roleNames[role] || role
}

// 获取操作显示名称
export function getActionDisplayName(action: string) {
  const actionNames: Record<string, string> = {
    login: '登录',
    logout: '登出',
    create: '创建',
    read: '查看',
    update: '更新',
    delete: '删除',
  }
  return actionNames[action] || action
}
