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

export const usePermissionStore = defineStore('permission', () => {
  const authStore = useAuthStore()

  // 权限缓存
  const permissionCache = ref<PermissionCache>({})
  const loading = ref(false)

  /**
   * 生成权限缓存键
   */
  const generateCacheKey = (role: string, resource: string, action: string): string => {
    return `${role}:${resource}:${action}`
  }

  /**
   * 检查权限（带缓存）
   * @param resource 资源
   * @param action 操作
   * @param useCache 是否使用缓存
   * @returns 是否有权限
   */
  const checkPermission = async (resource: string, action: string, useCache = true): Promise<boolean> => {
    const userRole = authStore.user?.role
    if (!userRole) {
      return false
    }

    // 超级管理员拥有所有权限
    if (userRole === 'super_admin') {
      return true
    }

    const cacheKey = generateCacheKey(userRole, resource, action)

    // 检查缓存
    if (useCache && permissionCache.value[cacheKey] !== undefined) {
      return permissionCache.value[cacheKey]
    }

    try {
      loading.value = true
      const response = await permissionAPI.checkPermission(userRole, resource, action)
      const hasPermission = !!response.data?.has_permission

      // 更新缓存
      permissionCache.value[cacheKey] = hasPermission

      return hasPermission
    }
    catch (error) {
      console.error('权限检查失败:', error)
      return false
    }
    finally {
      loading.value = false
    }
  }

  /**
   * 批量检查权限
   * @param permissions 权限列表
   * @returns 权限结果映射
   */
  const checkPermissions = async (permissions: Array<{ resource: string, action: string }>): Promise<Record<string, boolean>> => {
    const results: Record<string, boolean> = {}

    for (const permission of permissions) {
      const key = `${permission.resource}:${permission.action}`
      results[key] = await checkPermission(permission.resource, permission.action)
    }

    return results
  }

  /**
   * 同步检查权限（仅使用缓存）
   * @param resource 资源
   * @param action 操作
   * @returns 是否有权限
   */
  const hasPermission = (resource: string, action: string): boolean => {
    const userRole = authStore.user?.role
    if (!userRole) {
      return false
    }

    // 超级管理员拥有所有权限
    if (userRole === 'super_admin') {
      return true
    }

    const cacheKey = generateCacheKey(userRole, resource, action)
    return permissionCache.value[cacheKey] || false
  }

  /**
   * 检查菜单权限
   * @param menuPermission 菜单权限配置
   * @returns 是否有权限
   */
  const hasMenuPermission = (menuPermission: MenuPermission): boolean => {
    return hasPermission(menuPermission.resource, menuPermission.action)
  }

  /**
   * 检查按钮权限
   * @param buttonPermission 按钮权限配置
   * @returns 是否有权限
   */
  const hasButtonPermission = (buttonPermission: ButtonPermission): boolean => {
    return hasPermission(buttonPermission.resource, buttonPermission.action)
  }

  /**
   * 预加载权限
   * @param permissions 权限列表
   */
  const preloadPermissions = async (permissions: Array<{ resource: string, action: string }>) => {
    await checkPermissions(permissions)
  }

  /**
   * 清除权限缓存
   */
  const clearCache = () => {
    permissionCache.value = {}
  }

  /**
   * 清除特定用户角色的缓存
   * @param role 用户角色
   */
  const clearRoleCache = (role: string) => {
    Object.keys(permissionCache.value).forEach((key) => {
      if (key.startsWith(`${role}:`)) {
        delete permissionCache.value[key]
      }
    })
  }

  // 计算属性：当前用户是否为管理员
  const isAdmin = computed(() => authStore.user?.role === 'super_admin')
  const isSecurityManager = computed(() => authStore.user?.role === 'sec_mgr')
  const isDeveloper = computed(() => authStore.user?.role === 'dev')
  const isAuditor = computed(() => authStore.user?.role === 'auditor')

  return {
    permissionCache,
    loading,
    checkPermission,
    checkPermissions,
    hasPermission,
    hasMenuPermission,
    hasButtonPermission,
    preloadPermissions,
    clearCache,
    clearRoleCache,
    isAdmin,
    isSecurityManager,
    isDeveloper,
    isAuditor,
  }
})
