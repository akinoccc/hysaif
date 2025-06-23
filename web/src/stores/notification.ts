import type { Pagination } from '@/api'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { type GetNotificationsRequest, type Notification, notificationApi, type NotificationStats } from '@/api/notification'

export const useNotificationStore = defineStore('notification', () => {
  // 状态
  const notifications = ref<Notification[]>([])
  const pagination = ref<Pagination>({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0,
  })
  const stats = ref<NotificationStats>({
    total_count: 0,
    unread_count: 0,
    read_count: 0,
    by_priority: {},
    by_type: {},
  })

  const filters = ref<GetNotificationsRequest>({
    status: 'all',
    type: '',
    priority: '',
  })

  const loading = ref(false)
  const lastUpdateTime = ref<number>(0)

  // 计算属性
  const unreadCount = computed(() => {
    console.log('unreadCount', stats.value)
    return stats.value.unread_count
  })
  const hasUnread = computed(() => stats.value.unread_count > 0)

  // 获取未读通知数量
  const loadUnreadCount = async () => {
    try {
      const response = await notificationApi.getUnreadCount()
      stats.value.unread_count = response.unread_count || 0
      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('加载未读通知数量失败:', error)
    }
  }

  // 获取通知列表
  const loadNotifications = async (params: GetNotificationsRequest = {}) => {
    loading.value = true
    try {
      const response = await notificationApi.getNotifications(params)
      notifications.value = response.data || []
      pagination.value = response.pagination
      return response
    }
    catch (error) {
      console.error('加载通知列表失败:', error)
      throw error
    }
    finally {
      loading.value = false
    }
  }

  // 获取统计信息
  const loadStats = async () => {
    try {
      const response = await notificationApi.getStats()
      stats.value = response
      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('加载统计信息失败:', error)
    }
  }

  // 标记单个通知为已读
  const markAsRead = async (id: string) => {
    try {
      await notificationApi.markAsRead(id)

      // 更新本地状态
      const notification = notifications.value.find(n => n.id === id)
      if (notification && notification.status === 'unread') {
        notification.status = 'read'
        notification.read_at = Date.now()
        stats.value.unread_count = Math.max(0, stats.value.unread_count - 1)
        stats.value.read_count += 1
      }

      console.log('markAsRead', stats.value, unreadCount.value)

      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('标记通知已读失败:', error)
      throw error
    }
  }

  // 标记所有通知为已读
  const markAllAsRead = async () => {
    try {
      await notificationApi.markAllAsRead()

      // 更新本地状态
      notifications.value.forEach((notification) => {
        if (notification.status === 'unread') {
          notification.status = 'read'
          notification.read_at = Date.now()
        }
      })

      stats.value.read_count += stats.value.unread_count
      stats.value.unread_count = 0
      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('标记所有通知已读失败:', error)
      throw error
    }
  }

  // 删除通知
  const deleteNotification = async (id: string) => {
    try {
      await notificationApi.deleteNotification(id)

      // 更新本地状态
      const index = notifications.value.findIndex(n => n.id === id)
      if (index > -1) {
        const notification = notifications.value[index]
        if (notification.status === 'unread') {
          stats.value.unread_count = Math.max(0, stats.value.unread_count - 1)
        }
        else {
          stats.value.read_count = Math.max(0, stats.value.read_count - 1)
        }
        stats.value.total_count = Math.max(0, stats.value.total_count - 1)
        notifications.value.splice(index, 1)
      }

      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('删除通知失败:', error)
      throw error
    }
  }

  // 批量删除通知
  const bulkDeleteNotifications = async (ids: string[]) => {
    try {
      let unreadCount = 0
      let readCount = 0

      for (const id of ids) {
        await notificationApi.deleteNotification(id)
        const notification = notifications.value.find(n => n.id === id)
        if (notification) {
          if (notification.status === 'unread') {
            unreadCount++
          }
          else {
            readCount++
          }
        }
      }

      // 更新本地状态
      notifications.value = notifications.value.filter(n => !ids.includes(n.id))
      stats.value.unread_count = Math.max(0, stats.value.unread_count - unreadCount)
      stats.value.read_count = Math.max(0, stats.value.read_count - readCount)
      stats.value.total_count = Math.max(0, stats.value.total_count - ids.length)

      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('批量删除通知失败:', error)
      throw error
    }
  }

  // 批量标记为已读
  const bulkMarkAsRead = async (ids: string[]) => {
    try {
      const unreadIds = ids.filter((id) => {
        const notification = notifications.value.find(n => n.id === id)
        return notification && notification.status === 'unread'
      })

      for (const id of unreadIds) {
        await notificationApi.markAsRead(id)
        const notification = notifications.value.find(n => n.id === id)
        if (notification) {
          notification.status = 'read'
          notification.read_at = Date.now()
        }
      }

      stats.value.unread_count = Math.max(0, stats.value.unread_count - unreadIds.length)
      stats.value.read_count += unreadIds.length
      lastUpdateTime.value = Date.now()
    }
    catch (error) {
      console.error('批量标记已读失败:', error)
      throw error
    }
  }

  // 重置状态
  const reset = () => {
    notifications.value = []
    stats.value = {
      total_count: 0,
      unread_count: 0,
      read_count: 0,
      by_priority: {},
      by_type: {},
    }
    loading.value = false
    lastUpdateTime.value = 0
  }

  // 定时更新未读数量
  let intervalId: number | null = null

  const startPolling = (interval = 30000) => {
    if (intervalId) {
      clearInterval(intervalId)
    }
    intervalId = setInterval(loadUnreadCount, interval)
  }

  const stopPolling = () => {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  return {
    // 状态
    notifications,
    pagination,
    stats,
    loading,
    lastUpdateTime,
    filters,

    // 计算属性
    unreadCount,
    hasUnread,

    // 方法
    loadUnreadCount,
    loadNotifications,
    loadStats,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    bulkDeleteNotifications,
    bulkMarkAsRead,
    reset,
    startPolling,
    stopPolling,
  }
})
