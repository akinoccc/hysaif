import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
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
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 缓存为空时的回退值，默认为 false
   * @returns 是否有权限
   */
  function hasPermission(resource: string, action: string, fallbackValue = false): boolean {
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
   * 初始化用户权限缓存
   * 在用户登录、页面刷新、应用初始化时调用
   */
  async function initializePermissions(): Promise<void> {
    const userRole = authStore.user?.role
    if (!userRole) {
      clearCache()
      return
    }

    // 超级管理员无需初始化权限缓存
    if (userRole === 'super_admin') {
      return
    }

    try {
      loading.value = true
      const response = await permissionAPI.getUserAllPermissions()
      const permissions = response.data?.permissions || {}

      // 清空现有缓存
      clearRoleCache(userRole)

      // 更新权限缓存
      Object.entries(permissions).forEach(([key, hasAccess]) => {
        const cacheKey = generateCacheKey(userRole, ...key.split(':') as [string, string])
        permissionCache.value[cacheKey] = hasAccess
        permissionStates.value[cacheKey] = 'loaded'
      })
    }
    catch (error) {
      console.error('权限缓存初始化失败:', error)
      // 初始化失败时清空缓存，避免使用过期数据
      clearRoleCache(userRole)
    }
    finally {
      loading.value = false
    }
  }

  /**
   * 响应式权限检查
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 回退值
   * @returns 响应式权限状态
   */
  function usePermissionState(resource: string, action: string, fallbackValue = false) {
    const hasAccess = ref(fallbackValue)

    const checkAccess = () => {
      hasAccess.value = hasPermission(resource, action, fallbackValue)
    }

    // 监听用户变化和权限缓存变化，重新检查权限
    watch(
      () => [authStore.user?.role, permissionCache.value],
      () => {
        if (authStore.user) {
          checkAccess()
        }
        else {
          hasAccess.value = fallbackValue
        }
      },
      { immediate: true, deep: true },
    )

    return {
      hasAccess,
      loading: computed(() => false), // 同步检查，无需loading状态
      checkAccess,
    }
  }

  /**
   * 检查菜单权限
   * @param menuPermission 菜单权限配置
   * @returns 是否有权限
   */
  function checkMenuPermission(menuPermission: MenuPermission): boolean {
    return hasPermission(menuPermission.resource, menuPermission.action)
  }

  /**
   * 检查按钮权限
   * @param buttonPermission 按钮权限配置
   * @returns 是否有权限
   */
  function checkButtonPermission(buttonPermission: ButtonPermission): boolean {
    return hasPermission(buttonPermission.resource, buttonPermission.action)
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

  /**
   * 权限指令相关函数
   */
  const permissionDirective = {
    /**
     * 检查元素是否应该显示
     * @param resource 资源
     * @param action 操作
     * @returns 是否显示
     */
    shouldShow: (resource: string, action: string): boolean => {
      return hasPermission(resource, action, false)
    },

    /**
     * 检查元素是否应该禁用
     * @param resource 资源
     * @param action 操作
     * @returns 是否禁用
     */
    shouldDisable: (resource: string, action: string): boolean => {
      return !hasPermission(resource, action, false)
    },
  }

  /**
   * 角色检查函数
   */
  const roleChecker = {
    hasRole: (role: string | string[]): boolean => {
      const userRole = authStore.user?.role
      if (!userRole)
        return false

      if (Array.isArray(role)) {
        return role.includes(userRole)
      }
      return userRole === role
    },

    hasAnyRole: (roles: string[]): boolean => {
      const userRole = authStore.user?.role
      if (!userRole)
        return false
      return roles.includes(userRole)
    },

    hasAllRoles: (roles: string[]): boolean => {
      const userRole = authStore.user?.role
      if (!userRole)
        return false
      // 对于单个用户角色，检查是否在所需角色列表中
      return roles.includes(userRole)
    },
  }

  /**
   * 模块权限检查（计算属性）
   */
  const modulePermissions = {
    // 用户管理权限
    user: computed(() => ({
      canCreate: hasPermission('user', 'create'),
      canRead: hasPermission('user', 'read'),
      canUpdate: hasPermission('user', 'update'),
      canDelete: hasPermission('user', 'delete'),
      canManage: hasPermission('user', 'create') || hasPermission('user', 'update') || hasPermission('user', 'delete'),
    })),

    // 敏感信息权限
    secret: computed(() => ({
      canCreate: hasPermission('secret', 'create'),
      canRead: hasPermission('secret', 'read'),
      canUpdate: hasPermission('secret', 'update'),
      canDelete: hasPermission('secret', 'delete'),
      canRequest: hasPermission('secret', 'request'),
      canTemp: hasPermission('secret', 'temp'),
      canManage: hasPermission('secret', 'create') || hasPermission('secret', 'update') || hasPermission('secret', 'delete'),
    })),

    // 策略管理权限
    policy: computed(() => ({
      canCreate: hasPermission('policy', 'create'),
      canRead: hasPermission('policy', 'read'),
      canUpdate: hasPermission('policy', 'update'),
      canDelete: hasPermission('policy', 'delete'),
      canManage: hasPermission('policy', 'create') || hasPermission('policy', 'update') || hasPermission('policy', 'delete'),
    })),

    // 审计日志权限
    audit: computed(() => ({
      canRead: hasPermission('audit', 'read'),
    })),

    // 仪表盘权限
    dashboard: computed(() => ({
      canRead: hasPermission('dashboard', 'read'),
    })),

    // 访问请求权限
    accessRequest: computed(() => ({
      canRead: hasPermission('access_request', 'read'),
      canCreate: hasPermission('access_request', 'create'),
      canUpdate: hasPermission('access_request', 'update'),
      canApprove: hasPermission('access_request', 'approve'),
      canReject: hasPermission('access_request', 'reject'),
    })),
  }

  // 计算属性：当前用户角色
  const currentRole = computed(() => authStore.user?.role)
  const isAdmin = computed(() => currentRole.value === 'super_admin')
  const isSecurityManager = computed(() => currentRole.value === 'sec_mgr')
  const isDeveloper = computed(() => currentRole.value === 'dev')
  const isAuditor = computed(() => currentRole.value === 'auditor')

  return {
    // 基础权限检查
    permissionCache,
    loading,
    hasPermission,
    usePermissionState,
    checkMenuPermission,
    checkButtonPermission,
    initializePermissions,
    clearCache,
    clearRoleCache,

    // 角色相关
    currentRole,
    isAdmin,
    isSecurityManager,
    isDeveloper,
    isAuditor,
    roleChecker,

    // 模块权限
    modulePermissions,

    // 权限指令
    permissionDirective,
  }
})
