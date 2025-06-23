import type { RouteRecordNormalized } from 'vue-router'
import type { MenuPermission } from '@/stores/permission'

// 权限配置映射
export const PERMISSION_CONFIG = {
  // 用户管理权限
  USER: {
    CREATE: { resource: 'user', action: 'create' },
    READ: { resource: 'user', action: 'read' },
    UPDATE: { resource: 'user', action: 'update' },
    DELETE: { resource: 'user', action: 'delete' },
  },
  // 敏感信息权限
  SECRET: {
    CREATE: { resource: 'secret', action: 'create' },
    READ: { resource: 'secret', action: 'read' },
    UPDATE: { resource: 'secret', action: 'update' },
    DELETE: { resource: 'secret', action: 'delete' },
    REQUEST: { resource: 'secret', action: 'request' },
    TEMP: { resource: 'secret', action: 'temp' },
  },
  // 策略管理权限
  POLICY: {
    CREATE: { resource: 'policy', action: 'create' },
    READ: { resource: 'policy', action: 'read' },
    UPDATE: { resource: 'policy', action: 'update' },
    DELETE: { resource: 'policy', action: 'delete' },
  },
  // 审计日志权限
  AUDIT: {
    READ: { resource: 'audit', action: 'read' },
  },
  // 仪表盘权限
  DASHBOARD: {
    READ: { resource: 'dashboard', action: 'read' },
  },
} as const

// 菜单权限映射
export const MENU_PERMISSIONS: Record<string, MenuPermission> = {
  '/dashboard': PERMISSION_CONFIG.DASHBOARD.READ,
  '/users': PERMISSION_CONFIG.USER.READ,
  '/audit': PERMISSION_CONFIG.AUDIT.READ,
  '/permission': PERMISSION_CONFIG.POLICY.READ,
  '/api_key': PERMISSION_CONFIG.SECRET.READ,
  '/access_key': PERMISSION_CONFIG.SECRET.READ,
  '/ssh_key': PERMISSION_CONFIG.SECRET.READ,
  '/password': PERMISSION_CONFIG.SECRET.READ,
  '/certificate': PERMISSION_CONFIG.SECRET.READ,
  '/token': PERMISSION_CONFIG.SECRET.READ,
  '/custom': PERMISSION_CONFIG.SECRET.READ,
}

// 按钮权限映射
export const BUTTON_PERMISSIONS = {
  // 用户管理按钮
  USER_CREATE: PERMISSION_CONFIG.USER.CREATE,
  USER_EDIT: PERMISSION_CONFIG.USER.UPDATE,
  USER_DELETE: PERMISSION_CONFIG.USER.DELETE,

  // 敏感信息按钮
  SECRET_CREATE: PERMISSION_CONFIG.SECRET.CREATE,
  SECRET_EDIT: PERMISSION_CONFIG.SECRET.UPDATE,
  SECRET_DELETE: PERMISSION_CONFIG.SECRET.DELETE,
  SECRET_VIEW: PERMISSION_CONFIG.SECRET.READ,

  // 策略管理按钮
  POLICY_CREATE: PERMISSION_CONFIG.POLICY.CREATE,
  POLICY_EDIT: PERMISSION_CONFIG.POLICY.UPDATE,
  POLICY_DELETE: PERMISSION_CONFIG.POLICY.DELETE,
} as const

/**
 * 根据用户角色过滤路由
 * @param routes 路由列表
 * @param userRole 用户角色
 * @returns 过滤后的路由列表
 */
export function filterRoutesByRole(routes: RouteRecordNormalized[], userRole?: string): RouteRecordNormalized[] {
  if (!userRole)
    return []

  return routes.filter((route) => {
    // 如果路由没有角色限制，则显示
    if (!route.meta?.roles)
      return true

    // 检查用户角色是否在允许的角色列表中
    const allowedRoles = route.meta.roles as string[]
    return allowedRoles.includes(userRole)
  })
}

/**
 * 根据权限过滤菜单项
 * @param menuItems 菜单项列表
 * @param hasPermission 权限检查函数
 * @returns 过滤后的菜单项列表
 */
export function filterMenuByPermission(
  menuItems: any[],
  hasPermission: (resource: string, action: string) => boolean,
): any[] {
  return menuItems.filter((item) => {
    const permission = MENU_PERMISSIONS[item.path]
    if (!permission)
      return true // 没有权限配置的菜单项默认显示

    return hasPermission(permission.resource, permission.action)
  })
}

/**
 * 检查路由权限
 * @param routePath 路由路径
 * @param hasPermission 权限检查函数
 * @returns 是否有权限访问
 */
export function checkRoutePermission(
  routePath: string,
  hasPermission: (resource: string, action: string) => boolean,
): boolean {
  const permission = MENU_PERMISSIONS[routePath]
  if (!permission)
    return true // 没有权限配置的路由默认允许访问

  return hasPermission(permission.resource, permission.action)
}

/**
 * 获取用户可访问的菜单列表
 * @param allMenus 所有菜单
 * @param userRole 用户角色
 * @param hasPermission 权限检查函数
 * @returns 可访问的菜单列表
 */
export function getAccessibleMenus(
  allMenus: any[],
  userRole?: string,
  hasPermission?: (resource: string, action: string) => boolean,
): any[] {
  if (!userRole)
    return []

  // 首先根据角色过滤
  let filteredMenus = allMenus.filter((menu) => {
    if (!menu.meta?.roles)
      return true
    return (menu.meta.roles as string[]).includes(userRole)
  })

  // 然后根据权限过滤
  if (hasPermission) {
    filteredMenus = filterMenuByPermission(filteredMenus, hasPermission)
  }

  return filteredMenus
}

/**
 * 预加载页面权限
 * @param routePath 路由路径
 * @returns 需要预加载的权限列表
 */
export function getPagePermissions(routePath: string): Array<{ resource: string, action: string }> {
  const permissions: Array<{ resource: string, action: string }> = []

  // 根据路由路径确定需要预加载的权限
  switch (true) {
    case routePath.startsWith('/users'):
      permissions.push(
        PERMISSION_CONFIG.USER.READ,
        PERMISSION_CONFIG.USER.CREATE,
        PERMISSION_CONFIG.USER.UPDATE,
        PERMISSION_CONFIG.USER.DELETE,
      )
      break

    case routePath.startsWith('/permission'):
      permissions.push(
        PERMISSION_CONFIG.POLICY.READ,
        PERMISSION_CONFIG.POLICY.CREATE,
        PERMISSION_CONFIG.POLICY.UPDATE,
        PERMISSION_CONFIG.POLICY.DELETE,
      )
      break

    case routePath.startsWith('/audit'):
      permissions.push(PERMISSION_CONFIG.AUDIT.READ)
      break

    case routePath.includes('api_key')
      || routePath.includes('access_key')
      || routePath.includes('ssh_key')
      || routePath.includes('password')
      || routePath.includes('certificate')
      || routePath.includes('token')
      || routePath.includes('custom'):
      permissions.push(
        PERMISSION_CONFIG.SECRET.READ,
        PERMISSION_CONFIG.SECRET.CREATE,
        PERMISSION_CONFIG.SECRET.UPDATE,
        PERMISSION_CONFIG.SECRET.DELETE,
      )
      break

    case routePath.startsWith('/dashboard'):
      permissions.push(PERMISSION_CONFIG.DASHBOARD.READ)
      break

    default:
      // 默认权限
      permissions.push({ resource: 'dashboard', action: 'read' })
  }

  return permissions
}

/**
 * 权限常量
 */
export const PERMISSIONS = PERMISSION_CONFIG
export const MENU_PERMS = MENU_PERMISSIONS
export const BUTTON_PERMS = BUTTON_PERMISSIONS
