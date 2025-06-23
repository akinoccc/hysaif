import type { Component } from 'vue'
import type { ApiListResponse } from './types'
import { AlertCircle, AlertTriangle, Bell, Clock, Info, UserCheck, UserPlus, UserX } from 'lucide-vue-next'
import { api } from './http'

// 通知相关类型定义
export interface Notification {
  id: string
  recipient_id: string
  type: string
  title: string
  content: string
  status: 'unread' | 'read'
  priority: 'low' | 'normal' | 'high' | 'urgent'
  related_id?: string
  related_type?: string
  read_at?: number
  expires_at?: number
  metadata?: string
  created_at: number
  updated_at: number
  is_expired?: boolean
}

export interface NotificationListResponse extends ApiListResponse<Notification> {
  unread_count: number
}

export interface NotificationStats {
  total_count: number
  unread_count: number
  read_count: number
  by_type: Record<string, number>
  by_priority: Record<string, number>
}

export interface CreateNotificationRequest {
  recipient_ids: string[]
  type: string
  title: string
  content: string
  priority?: 'low' | 'normal' | 'high' | 'urgent'
  expires_at?: number
}

export interface BulkNotificationRequest {
  user_roles?: string[]
  user_ids?: string[]
  type: string
  title: string
  content: string
  priority?: 'low' | 'normal' | 'high' | 'urgent'
  expires_at?: number
}

export interface NotificationTemplate {
  type: string
  title: string
  content: string
  priority: string
  description: string
  variables: string[]
}

export interface GetNotificationsRequest {
  page?: number
  page_size?: number
  status?: 'unread' | 'read' | 'all'
  type?: string
  priority?: string
}

// 通知API接口
export const notificationApi = {
  // 获取通知列表
  getNotifications: (params?: GetNotificationsRequest) => {
    return api.get<any, NotificationListResponse>('/notifications', { params })
  },

  // 标记通知为已读
  markAsRead: (id: string) => {
    return api.put(`/notifications/${id}/read`)
  },

  // 标记所有通知为已读
  markAllAsRead: () => {
    return api.put('/notifications/read-all')
  },

  // 获取未读通知数量
  getUnreadCount: () => {
    return api.get<any, { unread_count: number }>('/notifications/unread-count')
  },

  // 删除通知
  deleteNotification: (id: string) => {
    return api.delete(`/notifications/${id}`)
  },

  // 获取通知统计
  getStats: () => {
    return api.get<any, NotificationStats>('/notifications/stats')
  },

  // 创建通知（管理员）
  createNotification: (data: CreateNotificationRequest) => {
    return api.post('/notifications', data)
  },

  // 批量发送通知（管理员）
  sendBulkNotification: (data: BulkNotificationRequest) => {
    return api.post('/notifications/bulk', data)
  },

  // 获取通知模板
  getTemplates: () => {
    return api.get<NotificationTemplate[]>('/notifications/templates')
  },
}

// 通知类型映射
export const notificationTypeMap: Record<string, { icon: Component, label: string }> = {
  access_request_created: {
    icon: UserPlus,
    label: '新的访问申请',
  },
  access_request_approved: {
    icon: UserCheck,
    label: '访问申请已批准',
  },
  access_request_rejected: {
    icon: UserX,
    label: '访问申请已拒绝',
  },
  access_request_expired: {
    icon: Clock,
    label: '访问权限已过期',
  },
  secret_item_expiring: {
    icon: Clock,
    label: '密钥项即将过期',
  },
  secret_item_expired: {
    icon: Clock,
    label: '密钥项已过期',
  },
  security_alert: {
    icon: AlertCircle,
    label: '安全警报',
  },
}

// 通知优先级映射
export const notificationPriorityMap: Record<string, { label: string, color: string }> = {
  low: { label: '低', color: 'bg-gray-500' },
  normal: { label: '普通', color: 'bg-blue-500' },
  high: { label: '高', color: 'bg-orange-500' },
  urgent: { label: '紧急', color: 'bg-red-500' },
}

// 通知优先级图标
export const notificationPriorityIcons: Record<string, Component> = {
  low: Info,
  normal: Bell,
  high: AlertTriangle,
  urgent: AlertCircle,
}
