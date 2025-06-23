import type { RouteRecordNormalized } from 'vue-router'
import { getAccessibleMenus } from './permission'

// 菜单项接口
export interface MenuItem {
  path: string
  title: string
  icon?: any
  order?: number
  children?: MenuItem[]
}

/**
 * 从路由生成菜单
 * @param routes 路由列表
 * @param userRole 用户角色
 * @param hasPermission 权限检查函数
 * @returns 菜单项列表
 */
export function generateMenuFromRoutes(
  routes: RouteRecordNormalized[],
  userRole?: string,
  hasPermission?: (resource: string, action: string) => boolean,
): MenuItem[] {
  const menuItems: MenuItem[] = []

  routes.forEach((route) => {
    // 检查路由是否有菜单配置
    if (route.meta?.menu && route.meta.menu.showInMenu) {
      // 检查角色权限
      if (route.meta.roles && userRole) {
        const allowedRoles = route.meta.roles as string[]
        if (!allowedRoles.includes(userRole)) {
          return // 跳过没有权限的菜单项
        }
      }

      const menuItem: MenuItem = {
        path: route.path,
        title: route.meta.menu.title,
        icon: route.meta.menu.icon,
        order: route.meta.menu.order || 999,
      }

      menuItems.push(menuItem)
    }
  })

  // 使用权限工具函数进一步过滤菜单
  const accessibleMenus = hasPermission
    ? getAccessibleMenus(menuItems, userRole, hasPermission)
    : menuItems

  // 按order排序
  return accessibleMenus.sort((a, b) => (a.order || 999) - (b.order || 999))
}

/**
 * 检查菜单是否激活
 * @param menuPath 菜单路径
 * @param currentPath 当前路径
 * @returns 是否激活
 */
export function isMenuActive(menuPath: string, currentPath: string): boolean {
  if (menuPath === currentPath) {
    return true
  }

  // 检查是否为子路径
  if (currentPath.startsWith(menuPath) && menuPath !== '/') {
    return true
  }

  return false
}
