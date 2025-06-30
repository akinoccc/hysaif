import type {
  ApiListResponse,
  ApiMethod,
  AuditLog,
  AuditLogsParams,
} from './types'
import { api } from './http'

/**
 * 审计日志相关API
 */
export const auditAPI = {
  /**
   * 获取审计日志列表
   * @param params 查询参数，包括分页、用户ID、操作类型、时间范围等
   * @returns 审计日志列表响应
   */
  getLogs: (params?: AuditLogsParams): ApiMethod<ApiListResponse<AuditLog>> => {
    return api.get('/audit/logs', { params })
  },
}

export default auditAPI
