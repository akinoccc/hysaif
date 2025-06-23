import type {
  ApiListResponse,
  ApiMethod,
  ItemsListParams,
  PostItemRequest,
  SecretItem,
} from './types'
import api from './http'

/**
 * 信息项相关API
 */
export const secretItemAPI = {
  /**
   * 获取信息项列表
   * @param params 查询参数，包括分页、搜索、分类等
   * @returns 信息项列表响应
   */
  getItems: (params?: ItemsListParams): ApiMethod<ApiListResponse<SecretItem>> => {
    return api.get('/items', { params })
  },

  /**
   * 创建新的信息项
   * @param data 创建信息项的数据
   * @returns 创建的信息项
   */
  createItem: <D>(data: PostItemRequest<D>): ApiMethod<SecretItem> => {
    return api.post('/items', data)
  },

  /**
   * 获取单个信息项详情
   * @param id 信息项ID
   * @returns 信息项详情
   */
  getItem: (id: string | number): ApiMethod<SecretItem> => {
    return api.get(`/items/${id}`)
  },

  /**
   * 更新信息项
   * @param id 信息项ID
   * @param data 更新的数据
   * @returns 更新后的信息项
   */
  updateItem: <D>(id: string | number, data: PostItemRequest<D>): ApiMethod<SecretItem> => {
    return api.put(`/items/${id}`, data)
  },

  /**
   * 删除信息项
   * @param id 信息项ID
   * @returns 删除响应
   */
  deleteItem: (id: string | number): ApiMethod<{ message: string }> => {
    return api.delete(`/items/${id}`)
  },

  /**
   * 通过申请访问密钥项详情
   * @param id 信息项ID
   * @returns 信息项详情
   */
  getItemWithAccess: (id: string | number): ApiMethod<SecretItem> => {
    return api.get(`/items/${id}/access`)
  },
}

export default secretItemAPI
