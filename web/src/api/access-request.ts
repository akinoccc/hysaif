import type {
  AccessRequest,
  AccessRequestListParams,
  ApiListResponse,
  ApiMethod,
  ApproveAccessRequestRequest,
  CreateAccessRequestRequest,
  RejectAccessRequestRequest,
  RevokeAccessRequestRequest,
} from './types'
import api from './http'

/**
 * 访问申请相关API
 */
export const accessRequestAPI = {
  /**
   * 创建访问申请
   * @param data 申请数据
   * @returns 创建的申请
   */
  createRequest: (data: CreateAccessRequestRequest): ApiMethod<AccessRequest> => {
    return api.post('/access-requests', data)
  },

  /**
   * 获取访问申请列表
   * @param params 查询参数
   * @returns 申请列表响应
   */
  getRequests: (params?: AccessRequestListParams): ApiMethod<ApiListResponse<AccessRequest>> => {
    return api.get('/access-requests', { params })
  },

  /**
   * 批准访问申请
   * @param id 申请ID
   * @param data 批准数据
   * @returns 更新后的申请
   */
  approveRequest: (id: string | number, data: ApproveAccessRequestRequest): ApiMethod<AccessRequest> => {
    return api.put(`/access-requests/${id}/approve`, data)
  },

  /**
   * 拒绝访问申请
   * @param id 申请ID
   * @param data 拒绝数据
   * @returns 更新后的申请
   */
  rejectRequest: (id: string | number, data: RejectAccessRequestRequest): ApiMethod<AccessRequest> => {
    return api.put(`/access-requests/${id}/reject`, data)
  },

  /**
   * 作废访问申请
   * @param id 申请ID
   * @param data 作废数据
   * @returns 更新后的申请
   */
  revokeRequest: (id: string | number, data: RevokeAccessRequestRequest): ApiMethod<AccessRequest> => {
    return api.put(`/access-requests/${id}/revoke`, data)
  },
}

export default accessRequestAPI
