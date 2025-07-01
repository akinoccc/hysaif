import { computed, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { type ButtonPermission, type MenuPermission, usePermissionStore } from '@/stores/permission'

/**
 * 权限检查组合式函数
 */
export function usePermission() {
  const permissionStore = usePermissionStore()
  const authStore = useAuthStore()

  /**
   * 异步权限检查
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 回退值
   * @returns 权限检查结果
   */
  const checkPermission = async (resource: string, action: string, fallbackValue = false) => {
    return await permissionStore.hasPermission(resource, action, fallbackValue)
  }

  /**
   * 同步权限检查（仅使用缓存）
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 回退值
   * @returns 是否有权限
   */
  const hasPermission = (resource: string, action: string, fallbackValue = false) => {
    return permissionStore.hasPermissionSync(resource, action, fallbackValue)
  }

  /**
   * 响应式权限检查
   * @param resource 资源
   * @param action 操作
   * @param fallbackValue 回退值
   * @returns 响应式权限状态
   */
  const usePermissionState = (resource: string, action: string, fallbackValue = false) => {
    const hasAccess = ref(fallbackValue)
    const loading = ref(false)

    const checkAccess = async () => {
      loading.value = true
      try {
        hasAccess.value = await checkPermission(resource, action, fallbackValue)
      }
      finally {
        loading.value = false
      }
    }

    // 监听用户变化，重新检查权限
    watch(
      () => authStore.user?.role,
      () => {
        if (authStore.user) {
          checkAccess()
        }
        else {
          hasAccess.value = fallbackValue
        }
      },
      { immediate: true },
    )

    return {
      hasAccess,
      loading,
      checkAccess,
    }
  }

  /**
   * 菜单权限检查
   * @param menuPermission 菜单权限配置
   * @param useAsync 是否使用异步检查
   * @returns 是否有菜单权限
   */
  const hasMenuPermission = (menuPermission: MenuPermission, useAsync = false) => {
    return permissionStore.checkMenuPermission(menuPermission, useAsync)
  }

  /**
   * 按钮权限检查
   * @param buttonPermission 按钮权限配置
   * @param useAsync 是否使用异步检查
   * @returns 是否有按钮权限
   */
  const hasButtonPermission = (buttonPermission: ButtonPermission, useAsync = false) => {
    return permissionStore.checkButtonPermission(buttonPermission, useAsync)
  }

  /**
   * 批量权限检查
   * @param permissions 权限列表
   * @returns 权限检查结果
   */
  const checkMultiplePermissions = async (permissions: Array<{ resource: string, action: string }>) => {
    return await permissionStore.checkMultiplePermissions(permissions)
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

  /**
   * 预加载权限
   * @param permissions 权限列表
   */
  const preloadPermissions = async (permissions: Array<{ resource: string, action: string }>) => {
    await permissionStore.preloadPermissions(permissions)
  }

  /**
   * 清除权限缓存
   */
  const clearPermissionCache = () => {
    permissionStore.clearCache()
  }

  return {
    checkPermission,
    hasPermission,
    usePermissionState,
    hasMenuPermission,
    hasButtonPermission,
    checkMultiplePermissions,
    modulePermissions,
    preloadPermissions,
    clearPermissionCache,
    // 导出store中的计算属性
    isAdmin: computed(() => permissionStore.isAdmin),
    isSecurityManager: computed(() => permissionStore.isSecurityManager),
    isDeveloper: computed(() => permissionStore.isDeveloper),
    isAuditor: computed(() => permissionStore.isAuditor),
    currentRole: computed(() => permissionStore.currentRole),
  }
}

/**
 * 权限指令相关函数
 */
export function usePermissionDirective() {
  const permissionStore = usePermissionStore()

  /**
   * 检查元素是否应该显示
   * @param resource 资源
   * @param action 操作
   * @returns 是否显示
   */
  const shouldShow = (resource: string, action: string): boolean => {
    return permissionStore.hasPermissionSync(resource, action, false)
  }

  /**
   * 检查元素是否应该禁用
   * @param resource 资源
   * @param action 操作
   * @returns 是否禁用
   */
  const shouldDisable = (resource: string, action: string): boolean => {
    return !permissionStore.hasPermissionSync(resource, action, false)
  }

  return {
    shouldShow,
    shouldDisable,
  }
}

/**
 * 角色检查组合式函数
 */
export function useRole() {
  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()

  const hasRole = (role: string | string[]): boolean => {
    const userRole = authStore.user?.role
    if (!userRole)
      return false

    if (Array.isArray(role)) {
      return role.includes(userRole)
    }
    return userRole === role
  }

  const hasAnyRole = (roles: string[]): boolean => {
    const userRole = authStore.user?.role
    if (!userRole)
      return false
    return roles.includes(userRole)
  }

  const hasAllRoles = (roles: string[]): boolean => {
    const userRole = authStore.user?.role
    if (!userRole)
      return false
    // 对于单个用户角色，检查是否在所需角色列表中
    return roles.includes(userRole)
  }

  return {
    hasRole,
    hasAnyRole,
    hasAllRoles,
    currentRole: computed(() => permissionStore.currentRole),
  }
}
