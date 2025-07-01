import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { permissionAPI } from '@/api'
import { useAuthStore } from './auth'

// 权限缓存接口
interface PermissionCache {
  [key: string]: boolean
}

// 菜单权限配置
export interface MenuPermission {
  resource: string
  action: string
}

// 按钮权限配置
export interface ButtonPermission {
  resource: string
  action: string
  label?: string
}

// 权限检查状态
interface PermissionState {
  [key: string]: 'checking' | 'loaded'
}

export const usePermissionStore = defineStore('permission', () => {
  const authStore = useAuthStore()

  // 权限缓存
  const permissionCache = ref<PermissionCache>({})
  // 权限检查状态，避免重复请求
  const permissionStates = ref<PermissionState>({})
  const loading = ref(false)

  /**
   * 生成权限缓存键
   */
  function generateCacheKey(role: string, resource: string, action: string): string {
    return `${role}:${resource}:${action}`
  }

  /**
   * 统一的权限检查入口
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 缓存为空时的回退值，默认为 false
   * @returns 是否有权限
   */
  async function hasPermission(resource: string, action: string, fallbackValue = false): Promise<boolean> {
    const userRole = authStore.user?.role
    if (!userRole) {
      return false
    }

    // 超级管理员拥有所有权限
    if (userRole === 'super_admin') {
      return true
    }

    const cacheKey = generateCacheKey(userRole, resource, action)

    // 如果缓存中有值，直接返回
    if (permissionCache.value[cacheKey] !== undefined) {
      return permissionCache.value[cacheKey]
    }

    // 如果正在检查中，等待结果
    if (permissionStates.value[cacheKey] === 'checking') {
      // 使用轮询等待检查完成
      while (permissionStates.value[cacheKey] === 'checking') {
        await new Promise(resolve => setTimeout(resolve, 50))
      }
      return permissionCache.value[cacheKey] ?? fallbackValue
    }

    // 开始权限检查
    permissionStates.value[cacheKey] = 'checking'

    try {
      loading.value = true
      const response = await permissionAPI.checkPermission(userRole, resource, action)
      const hasAccess = !!response.data?.has_permission

      // 更新缓存
      permissionCache.value[cacheKey] = hasAccess
      permissionStates.value[cacheKey] = 'loaded'

      return hasAccess
    }
    catch (error) {
      console.error('权限检查失败:', error)
      // 检查失败时，使用回退值
      permissionCache.value[cacheKey] = fallbackValue
      permissionStates.value[cacheKey] = 'loaded'
      return fallbackValue
    }
    finally {
      loading.value = false
    }
  }

  /**
   * 同步权限检查（仅使用缓存）
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 缓存为空时的回退值，默认为 false
   * @returns 是否有权限
   */
  function hasPermissionSync(resource: string, action: string, fallbackValue = false): boolean {
    const userRole = authStore.user?.role
    if (!userRole) {
      return false
    }

    // 超级管理员拥有所有权限
    if (userRole === 'super_admin') {
      return true
    }

    const cacheKey = generateCacheKey(userRole, resource, action)
    const cached = permissionCache.value[cacheKey]

    return cached !== undefined ? cached : fallbackValue
  }

  /**
   * 批量权限检查
   * @param permissions 权限列表
   * @returns 权限结果映射
   */
  async function checkMultiplePermissions(permissions: Array<{ resource: string, action: string }>): Promise<Record<string, boolean>> {
    const results: Record<string, boolean> = {}

    // 并行检查所有权限
    const promises = permissions.map(async (permission) => {
      const key = `${permission.resource}:${permission.action}`
      const hasAccess = await hasPermission(permission.resource, permission.action)
      results[key] = hasAccess
    })

    await Promise.all(promises)
    return results
  }

  /**
   * 预加载权限
   * @param permissions 权限列表
   */
  async function preloadPermissions(permissions: Array<{ resource: string, action: string }>) {
    const userRole = authStore.user?.role
    if (!userRole || userRole === 'super_admin') {
      return // 超级管理员无需预加载
    }

    // 过滤掉已经缓存的权限
    const uncachedPermissions = permissions.filter(({ resource, action }) => {
      const cacheKey = generateCacheKey(userRole, resource, action)
      return permissionCache.value[cacheKey] === undefined
        && permissionStates.value[cacheKey] !== 'checking'
    })

    if (uncachedPermissions.length === 0) {
      return
    }

    // 并行预加载权限
    await checkMultiplePermissions(uncachedPermissions)
  }

  /**
   * 检查菜单权限
   * @param menuPermission 菜单权限配置
   * @param useAsync 是否使用异步检查
   * @returns 是否有权限
   */
  function checkMenuPermission(menuPermission: MenuPermission, useAsync = false): boolean | Promise<boolean> {
    if (useAsync) {
      return hasPermission(menuPermission.resource, menuPermission.action)
    }
    return hasPermissionSync(menuPermission.resource, menuPermission.action)
  }

  /**
   * 检查按钮权限
   * @param buttonPermission 按钮权限配置
   * @param useAsync 是否使用异步检查
   * @returns 是否有权限
   */
  function checkButtonPermission(buttonPermission: ButtonPermission, useAsync = false): boolean | Promise<boolean> {
    if (useAsync) {
      return hasPermission(buttonPermission.resource, buttonPermission.action)
    }
    return hasPermissionSync(buttonPermission.resource, buttonPermission.action)
  }

  /**
   * 清除权限缓存
   */
  function clearCache() {
    permissionCache.value = {}
    permissionStates.value = {}
  }

  /**
   * 清除特定用户角色的缓存
   * @param role 用户角色
   */
  function clearRoleCache(role: string) {
    Object.keys(permissionCache.value).forEach((key) => {
      if (key.startsWith(`${role}:`)) {
        delete permissionCache.value[key]
        delete permissionStates.value[key]
      }
    })
  }

  // 计算属性：当前用户角色
  const currentRole = computed(() => authStore.user?.role)
  const isAdmin = computed(() => currentRole.value === 'super_admin')
  const isSecurityManager = computed(() => currentRole.value === 'sec_mgr')
  const isDeveloper = computed(() => currentRole.value === 'dev')
  const isAuditor = computed(() => currentRole.value === 'auditor')

  return {
    permissionCache,
    loading,
    hasPermission,
    hasPermissionSync,
    checkMultiplePermissions,
    checkMenuPermission,
    checkButtonPermission,
    preloadPermissions,
    clearCache,
    clearRoleCache,
    currentRole,
    isAdmin,
    isSecurityManager,
    isDeveloper,
    isAuditor,
  }
})
