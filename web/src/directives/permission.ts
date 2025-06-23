import type { App, Directive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

// 权限指令参数类型
interface PermissionBinding {
  resource: string
  action: string
  mode?: 'hide' | 'disable' // 默认为 hide
}

/**
 * 权限指令
 * 用法：
 * v-permission="{ resource: 'user', action: 'create' }" // 隐藏元素
 * v-permission="{ resource: 'user', action: 'create', mode: 'disable' }" // 禁用元素
 */
const permissionDirective: Directive = {
  mounted(el: HTMLElement, binding) {
    checkPermission(el, binding)
  },
  updated(el: HTMLElement, binding) {
    checkPermission(el, binding)
  },
}

function checkPermission(el: HTMLElement, binding: any) {
  const permissionStore = usePermissionStore()
  const authStore = useAuthStore()

  // 如果用户未登录，隐藏元素
  if (!authStore.isAuthenticated) {
    hideElement(el)
    return
  }

  const value = binding.value as PermissionBinding

  if (!value || !value.resource || !value.action) {
    console.warn('v-permission directive requires resource and action')
    return
  }

  const { resource, action, mode = 'hide' } = value
  const hasPermission = permissionStore.hasPermission(resource, action)

  if (!hasPermission) {
    if (mode === 'disable') {
      disableElement(el)
    }
    else {
      hideElement(el)
    }
  }
  else {
    showElement(el)
    enableElement(el)
  }
}

function hideElement(el: HTMLElement) {
  el.style.display = 'none'
}

function showElement(el: HTMLElement) {
  el.style.display = ''
}

function disableElement(el: HTMLElement) {
  el.setAttribute('disabled', 'true')
  el.style.opacity = '0.5'
  el.style.cursor = 'not-allowed'

  // 如果是按钮或输入框，添加disabled属性
  if (el.tagName === 'BUTTON' || el.tagName === 'INPUT' || el.tagName === 'SELECT' || el.tagName === 'TEXTAREA') {
    (el as HTMLInputElement).disabled = true
  }
}

function enableElement(el: HTMLElement) {
  el.removeAttribute('disabled')
  el.style.opacity = ''
  el.style.cursor = ''

  // 如果是按钮或输入框，移除disabled属性
  if (el.tagName === 'BUTTON' || el.tagName === 'INPUT' || el.tagName === 'SELECT' || el.tagName === 'TEXTAREA') {
    (el as HTMLInputElement).disabled = false
  }
}

/**
 * 角色指令
 * 用法：
 * v-role="'super_admin'" // 只有超级管理员可见
 * v-role="['super_admin', 'sec_mgr']" // 超级管理员或安全管理员可见
 */
const roleDirective: Directive = {
  mounted(el: HTMLElement, binding) {
    checkRole(el, binding)
  },
  updated(el: HTMLElement, binding) {
    checkRole(el, binding)
  },
}

function checkRole(el: HTMLElement, binding: any) {
  const authStore = useAuthStore()

  // 如果用户未登录，隐藏元素
  if (!authStore.isAuthenticated || !authStore.user) {
    hideElement(el)
    return
  }

  const requiredRoles = binding.value
  const userRole = authStore.user.role

  if (!requiredRoles) {
    console.warn('v-role directive requires role value')
    return
  }

  let hasRole = false

  if (Array.isArray(requiredRoles)) {
    hasRole = requiredRoles.includes(userRole)
  }
  else {
    hasRole = userRole === requiredRoles
  }

  if (!hasRole) {
    hideElement(el)
  }
  else {
    showElement(el)
  }
}

/**
 * 安装权限指令
 */
export function setupPermissionDirectives(app: App) {
  app.directive('permission', permissionDirective)
  app.directive('role', roleDirective)
}

export { permissionDirective, roleDirective }
