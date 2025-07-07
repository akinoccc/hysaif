import type {
  ApiMethod,
  ApiResponse,
} from './types'
import { api } from './http'

// 权限相关类型定义
export interface PermissionRequest {
  role: string
  resource: string
  action: string
}

export interface PermissionData {
  has_permission: boolean
  role: string
  resource: string
  action: string
}

export interface PermissionResponse extends ApiResponse<PermissionData> {}

// 批量权限检查请求接口
export interface BatchPermissionRequest {
  permissions: PermissionRequest[]
}

// 批量权限检查响应数据接口
export interface BatchPermissionData {
  results: Record<string, boolean>
}

// 批量权限检查响应接口
export interface BatchPermissionResponse extends ApiResponse<BatchPermissionData> {}

// 用户所有权限响应数据接口
export interface UserAllPermissionsData {
  role: string
  permissions: Record<string, boolean>
}

// 用户所有权限响应接口
export interface UserAllPermissionsResponse extends ApiResponse<UserAllPermissionsData> {}

export interface PolicyData {
  policies: string[][]
}

export interface PolicyResponse extends ApiResponse<PolicyData> {}

export interface RoleRequest {
  user: string
  role: string
}

export interface UserRolesData {
  user: string
  roles: string[]
}

export interface UserRolesResponse extends ApiResponse<UserRolesData> {}

export interface RoleUsersData {
  role: string
  users: string[]
}

export interface RoleUsersResponse extends ApiResponse<RoleUsersData> {}

export interface RolePermissionsData {
  role: string
  permissions: string[][]
}

export interface RolePermissionsResponse extends ApiResponse<RolePermissionsData> {}

export interface MenuItemData {
  path: string
  title: string
  icon: string
  order: number
}

export interface MenuData {
  menus: MenuItemData[]
}

export interface MenuResponse extends ApiResponse<MenuData> {}

/**
 * 权限管理相关API
 */
export const permissionAPI = {
  /**
   * 检查权限
   * @param role 角色
   * @param resource 资源
   * @param action 操作
   * @returns 权限检查结果
   */
  checkPermission: (role: string, resource: string, action: string): ApiMethod<PermissionResponse> => {
    const data: PermissionRequest = { role, resource, action }
    return api.post('/permissions/check', data)
  },

  /**
   * 批量检查权限
   * @param permissions 权限列表
   * @returns 批量权限检查结果
   */
  batchCheckPermissions: (permissions: PermissionRequest[]): ApiMethod<BatchPermissionResponse> => {
    const data: BatchPermissionRequest = { permissions }
    return api.post('/permissions/batch-check', data)
  },

  /**
   * 获取用户所有权限
   * @returns 用户所有权限
   */
  getUserAllPermissions: (): ApiMethod<UserAllPermissionsResponse> => {
    return api.get('/permissions/all')
  },

  /**
   * 添加权限策略
   * @param role 角色
   * @param resource 资源
   * @param action 操作
   * @returns 操作结果
   */
  addPolicy: (role: string, resource: string, action: string): ApiMethod<ApiResponse> => {
    const data: PermissionRequest = { role, resource, action }
    return api.post('/permissions/policies', data)
  },

  /**
   * 移除权限策略
   * @param role 角色
   * @param resource 资源
   * @param action 操作
   * @returns 操作结果
   */
  removePolicy: (role: string, resource: string, action: string): ApiMethod<ApiResponse> => {
    const data: PermissionRequest = { role, resource, action }
    return api.delete('/permissions/policies', { data })
  },

  /**
   * 获取所有权限策略
   * @returns 权限策略列表
   */
  getPolicies: (): ApiMethod<PolicyResponse> => {
    return api.get('/permissions/policies')
  },

  /**
   * 为用户添加角色
   * @param user 用户
   * @param role 角色
   * @returns 操作结果
   */
  addRoleForUser: (user: string, role: string): ApiMethod<ApiResponse> => {
    const data: RoleRequest = { user, role }
    return api.post('/permissions/users/roles', data)
  },

  /**
   * 删除用户角色
   * @param user 用户
   * @param role 角色
   * @returns 操作结果
   */
  deleteRoleForUser: (user: string, role: string): ApiMethod<ApiResponse> => {
    const data: RoleRequest = { user, role }
    return api.delete('/permissions/users/roles', { data })
  },

  /**
   * 获取用户的所有角色
   * @param user 用户
   * @returns 用户角色列表
   */
  getRolesForUser: (user: string): ApiMethod<UserRolesResponse> => {
    return api.get(`/permissions/users/${user}/roles`)
  },

  /**
   * 获取角色下的所有用户
   * @param role 角色
   * @returns 角色用户列表
   */
  getUsersForRole: (role: string): ApiMethod<RoleUsersResponse> => {
    return api.get(`/permissions/roles/${role}/users`)
  },

  /**
   * 获取角色的所有权限
   * @param role 角色
   * @returns 角色权限列表
   */
  getPermissionsForRole: (role: string): ApiMethod<RolePermissionsResponse> => {
    return api.get(`/permissions/roles/${role}/permissions`)
  },

  /**
   * 重新加载权限策略
   * @returns 操作结果
   */
  reloadPolicy: (): ApiMethod<ApiResponse> => {
    return api.post('/permissions/reload')
  },

  /**
   * 批量更新角色权限
   * @param role 角色
   * @param permissions 权限映射
   * @returns 操作结果
   */
  updateRolePermissions: (role: string, permissions: Record<string, string[]>): ApiMethod<ApiResponse> => {
    return api.put(`/permissions/roles/${role}/permissions`, { permissions })
  },

  /**
   * 获取用户可访问的菜单列表
   * @returns 菜单列表
   */
  getUserAccessibleMenus: (): ApiMethod<MenuResponse> => {
    return api.get('/permissions/menus')
  },
}

export default permissionAPI
