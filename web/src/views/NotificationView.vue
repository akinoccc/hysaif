<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import type { Notification } from '@/api/notification'
import type { FilterField } from '@/components/common/layout'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import { AlertTriangle, Bell, CheckCheck, CheckCircle, Clock, ExternalLink, Info, Trash2 } from 'lucide-vue-next'
import { computed, h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  notificationPriorityMap,
  notificationTypeMap,
} from '@/api/notification'
import { DataFilter, DataTable, PageHeader } from '@/components'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

const pageInfo = {
  title: '通知管理',
  description: '查看和管理系统通知消息',
}

const selectedNotifications = ref<string[]>([])
const showDeleteDialog = ref(false)
const showSingleDeleteDialog = ref(false)
const notificationToDelete = ref<string | null>(null)

// 筛选字段配置
const filterFields: FilterField[] = [
  {
    key: 'status',
    label: '状态',
    type: 'select',
    icon: Bell,
    options: [
      { value: 'all', label: '全部状态' },
      { value: 'unread', label: '未读' },
      { value: 'read', label: '已读' },
    ],
  },
  {
    key: 'type',
    label: '类型',
    type: 'select',
    icon: Info,
    options: [
      { value: 'all', label: '全部类型' },
      ...Object.entries(notificationTypeMap).map(([key, value]) => ({
        value: key,
        label: value.label,
      })),
    ],
  },
  {
    key: 'priority',
    label: '优先级',
    type: 'select',
    icon: AlertTriangle,
    options: [
      { value: 'all', label: '全部优先级' },
      ...Object.entries(notificationPriorityMap).map(([key, value]) => ({
        value: key,
        label: value.label,
      })),
    ],
  },
]

// 优先级颜色映射
const priorityColorMap: Record<string, string> = {
  low: 'bg-muted text-muted-foreground',
  normal: 'bg-info text-info-foreground',
  high: 'bg-warning text-warning-foreground',
  urgent: 'bg-destructive text-destructive-foreground',
}

// 计算属性
const isAllSelected = computed(() => {
  return notificationStore.notifications.length > 0 && selectedNotifications.value.length === notificationStore.notifications.length
})

const isIndeterminate = computed(() => {
  return selectedNotifications.value.length > 0 && selectedNotifications.value.length < notificationStore.notifications.length
})

// 表格列定义
const columns: ColumnDef<Notification>[] = [
  {
    id: 'select',
    header: ({ table }) => {
      return h(Checkbox, {
        'modelValue': isAllSelected.value,
        'indeterminate': isIndeterminate.value,
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => {
          toggleSelectAll()
          table.toggleAllPageRowsSelected(!!value)
        },
        'ariaLabel': '选择全部',
      })
    },
    cell: ({ row }) => {
      const notification = row.original
      return h(Checkbox, {
        'modelValue': selectedNotifications.value.includes(notification.id),
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => {
          toggleSelectNotification(notification.id)
          row.toggleSelected(!!value)
        },
        'ariaLabel': '选择行',
      })
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'title',
    header: '标题',
    cell: ({ row }) => {
      const notification = row.original
      return h('div', { class: 'flex items-center gap-3' }, [
        h('div', { class: 'flex items-center gap-2' }, [
          h('div', {
            class: {
              'w-1 h-1 rounded-full': true,
              'bg-primary': notification.status === 'unread',
            },
          }),
          h(notificationTypeMap[notification.type].icon, { class: 'h-4 w-4 text-muted-foreground' }),
          h('span', { class: notification.status === 'read' ? 'text-muted-foreground' : 'font-medium' }, notification.title),
        ]),
      ])
    },
  },
  {
    accessorKey: 'content',
    header: '內容',
    cell: ({ row }) => {
      return h('div', { class: 'text-sm text-muted-foreground' }, row.getValue('content'))
    },
  },
  {
    accessorKey: 'priority',
    header: '优先级',
    cell: ({ row }) => {
      const priority = row.getValue('priority') as string
      const colorClass = priorityColorMap[priority] || 'bg-muted text-muted-foreground'
      return h('div', { class: 'flex items-center gap-2' }, [
        h(Badge, { class: colorClass }, () => notificationPriorityMap[priority].label || priority),
      ])
    },
  },
  {
    accessorKey: 'created_at',
    header: '时间',
    cell: ({ row }) => {
      const date = row.getValue('created_at') as number
      return h('span', { class: 'text-sm text-muted-foreground' }, formatTime(date))
    },
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => {
      const notification = row.original
      return h('div', { class: 'flex items-center' }, [
        notification.status === 'unread' && h(Button, {
          variant: 'ghost',
          size: 'sm',
          onClick: () => notificationStore.markAsRead(notification.id),
        }, () => [h(CheckCheck, { class: 'h-4 w-4' }), '标记已读']),
        h(Button, {
          variant: 'ghost',
          size: 'sm',
          class: 'text-red-600 hover:bg-red-50 hover:text-red-600',
          onClick: () => openSingleDeleteDialog(notification.id),
        }, () => [h(Trash2, { class: 'h-4 w-4' }), '删除']),
        notification.related_type && notification.related_id && h(Button, {
          variant: 'ghost',
          size: 'sm',
          onClick: () => navigateToRelated(notification),
        }, () => [h(ExternalLink, { class: 'h-4 w-4' }), '查看详情']),
      ])
    },
  },
]

const selectedCount = computed(() => selectedNotifications.value.length)

const hasSelected = computed(() => selectedNotifications.value.length > 0)

// 筛选器处理
function handleFilterUpdate(newFilters: typeof notificationStore.filters) {
  Object.assign(notificationStore.filters, newFilters)
  notificationStore.pagination.page = 1
  selectedNotifications.value = [] // 清空选择状态
  notificationStore.loadNotifications()
}

function handleFilterReset() {
  Object.assign(notificationStore.filters, {
    searchQuery: '',
    status: 'all',
    type: 'all',
    priority: 'all',
  })
  notificationStore.pagination.page = 1
  selectedNotifications.value = [] // 清空选择状态
  notificationStore.loadNotifications()
}

// 切换选择通知
function toggleSelectNotification(id: string) {
  const index = selectedNotifications.value.indexOf(id)
  if (index > -1) {
    selectedNotifications.value.splice(index, 1)
  }
  else {
    selectedNotifications.value.push(id)
  }
}

// 切换全选
function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedNotifications.value = []
  }
  else {
    selectedNotifications.value = notificationStore.notifications.map(n => n.id)
  }
}

// 批量标记为已读
async function bulkMarkAsRead() {
  if (selectedNotifications.value.length === 0) {
    return
  }

  try {
    await notificationStore.bulkMarkAsRead(selectedNotifications.value)
    selectedNotifications.value = []
  }
  catch (error) {
    console.error('批量标记已读失败:', error)
  }
}

// 打开删除确认对话框
function openDeleteDialog() {
  if (selectedNotifications.value.length === 0) {
    return
  }
  showDeleteDialog.value = true
}

// 打开单个删除确认对话框
function openSingleDeleteDialog(id: string) {
  notificationToDelete.value = id
  showSingleDeleteDialog.value = true
}

// 确认单个删除
async function confirmSingleDelete() {
  if (!notificationToDelete.value)
    return

  try {
    await notificationStore.deleteNotification(notificationToDelete.value)
    notificationStore.pagination.total = Math.max(0, notificationStore.pagination.total - 1)
    notificationToDelete.value = null
    showSingleDeleteDialog.value = false
  }
  catch (error) {
    console.error('删除通知失败:', error)
  }
}

// 确认批量删除
async function confirmBulkDelete() {
  try {
    const deleteCount = selectedNotifications.value.length
    await notificationStore.bulkDeleteNotifications(selectedNotifications.value)
    notificationStore.pagination.total = Math.max(0, notificationStore.pagination.total - deleteCount)
    selectedNotifications.value = []
    showDeleteDialog.value = false
  }
  catch (error) {
    console.error('批量删除失败:', error)
  }
}

// 页码变化处理
function handlePageChange(page: number) {
  notificationStore.pagination.page = page
  selectedNotifications.value = [] // 清空选择状态
  notificationStore.loadNotifications()
}

// 导航到相关页面
function navigateToRelated(notification: Notification) {
  if (notification.related_type === 'access_request' && notification.related_id) {
    router.push(`/access-requests/${notification.related_id}`)
  }
  else if (notification.related_type === 'secret_item' && notification.related_id) {
    router.push(`/items/${notification.related_id}`)
  }
}

// 格式化时间
function formatTime(timestamp: number) {
  return formatDistanceToNow(new Date(timestamp), {
    addSuffix: true,
    locale: zhCN,
  })
}

// 初始化
onMounted(() => {
  notificationStore.loadNotifications()
  notificationStore.loadStats()
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <PageHeader
      :title="pageInfo.title" :description="pageInfo.description" :actions="[
        {
          text: hasSelected ? `标记已读 (${selectedCount})` : '标记已读',
          variant: 'default',
          disabled: !hasSelected,
          onClick: bulkMarkAsRead,
        },
        {
          text: hasSelected ? `删除 (${selectedCount})` : '删除',
          variant: 'destructive',
          disabled: !hasSelected,
          onClick: openDeleteDialog,
        },
      ]"
    >
      <template #extra>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2 text-sm">
            <Bell class="h-4 w-4" />
            <span>总计: {{ notificationStore.stats.total_count }}</span>
          </div>
          <div class="flex items-center gap-2 text-sm text-orange-600">
            <Clock class="h-4 w-4" />
            <span>未读: {{ notificationStore.stats.unread_count }}</span>
          </div>
          <div class="flex items-center gap-2 text-sm text-green-600">
            <CheckCircle class="h-4 w-4" />
            <span>已读: {{ notificationStore.stats.read_count }}</span>
          </div>
        </div>
      </template>
    </PageHeader>

    <!-- 主要内容 -->
    <Card>
      <CardContent class="px-6 space-y-6">
        <!-- 筛选器 -->
        <DataFilter
          v-model="notificationStore.filters"
          :quick-filters="filterFields"
          :total-items="notificationStore.pagination.total"
          @update:model-value="handleFilterUpdate"
          @reset="handleFilterReset"
        />
        <!-- 数据表格 -->
        <DataTable
          v-model:current-page="notificationStore.pagination.page"
          :columns="columns"
          :data="notificationStore.notifications"
          :loading="notificationStore.loading"
          :total="notificationStore.pagination.total"
          :page-size="notificationStore.pagination.page_size"
          empty-title="暂无通知"
          empty-description="当前没有任何通知消息"
          :empty-icon="Bell"
          @page-change="handlePageChange"
        />
      </CardContent>
    </Card>

    <!-- 批量删除确认对话框 -->
    <AlertDialog v-model:open="showDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除选中的 {{ selectedNotifications.length }} 条通知吗？此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">
            取消
          </AlertDialogCancel>
          <AlertDialogAction class="bg-destructive hover:bg-destructive/90" @click="confirmBulkDelete">
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 单个删除确认对话框 -->
    <AlertDialog v-model:open="showSingleDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除通知</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除这条通知吗？此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showSingleDeleteDialog = false">
            取消
          </AlertDialogCancel>
          <AlertDialogAction class="bg-destructive hover:bg-destructive/90" @click="confirmSingleDelete">
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<style scoped>
/* 自定义样式 */
</style>
