import type {
  ApiListResponse,
  ApiMethod,
  ChangePasswordRequest,
  CreateUserRequest,
  UpdateProfileRequest,
  UpdateUserRequest,
  User,
  UserListParams,
} from './types'
import api from './config'

/**
 * 用户相关API
 */
export const userAPI = {
  /**
   * 获取用户资料
   * @returns 用户资料信息
   */
  getProfile: (): ApiMethod<User> => {
    return api.get('/users/profile')
  },

  /**
   * 更新用户资料
   * @param data 更新的用户资料数据
   * @returns 更新后的用户资料
   */
  updateProfile: (data: UpdateProfileRequest): ApiMethod<User> => {
    return api.put('/users/profile', data)
  },

  /**
   * 修改密码
   * @param data 修改密码的数据
   * @returns 修改密码响应
   */
  changePassword: (data: ChangePasswordRequest): ApiMethod<{ message: string }> => {
    return api.put('/users/change-password', data)
  },

  /**
   * 获取用户列表
   * @param params 查询参数
   * @returns 用户列表
   */
  getUsers: (params?: UserListParams): ApiMethod<ApiListResponse<User>> => {
    return api.get('/users', { params })
  },

  /**
   * 获取指定用户信息
   * @param id 用户ID
   * @returns 用户信息
   */
  getUser: (id: string | number): ApiMethod<User> => {
    return api.get(`/users/${id}`)
  },

  /**
   * 创建用户
   * @param data 用户数据
   * @returns 创建的用户信息
   */
  createUser: (data: CreateUserRequest): ApiMethod<User> => {
    return api.post('/users', data)
  },

  /**
   * 更新用户信息
   * @param id 用户ID
   * @param data 更新的用户数据
   * @returns 更新后的用户信息
   */
  updateUser: (id: string | number, data: UpdateUserRequest): ApiMethod<User> => {
    return api.put(`/users/${id}`, data)
  },

  /**
   * 删除用户
   * @param id 用户ID
   * @returns 删除响应
   */
  deleteUser: (id: string | number): ApiMethod<{ message: string }> => {
    return api.delete(`/users/${id}`)
  },
}

export default userAPI
