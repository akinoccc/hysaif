<script setup lang="ts">
import {
  Activity,
  Calendar,
  ChevronRight,
  Clock,
  Database,
  Download,
  FileText,
  Globe,
  Loader2,
  Monitor,
  RefreshCw,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { auditAPI, type AuditLog } from '@/api'
import { type ActionButton, DataFilter, EmptyState, type FilterField, PageHeader } from '@/components'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationFirst,
  PaginationItem,
  PaginationLast,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import {
  Sheet,
  SheetClose,
  SheetContent,
} from '@/components/ui/sheet'
import { AUDIT_LOG_ACTION_MAP, AUDIT_LOG_RESOURCE_MAP } from '@/constants'
import { formatDate, formatRelativeTime } from '@/lib/utils'
import { getActionColor, getActionDisplayName, getActionIcon, getResourceDisplayName, getUserAgentInfo } from './helper'
import LogDetail from './LogDetail.vue'

export interface FilterValues {
  searchQuery: string
  action: string
  resource: string
  user: string
  dateRange?: [number, number]
}

const logs = ref<AuditLog[]>([])
const expandedLogs = ref<Set<string>>(new Set())
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalCount = ref(0)
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value))
const selectedLog = ref<AuditLog | undefined>()
const isSheetOpen = ref(false)

const filters = ref<FilterValues>({
  searchQuery: '',
  action: '',
  resource: '',
  user: '',
  dateRange: undefined as [number, number] | undefined,
})

// 快速筛选配置
const quickFilters = computed((): FilterField[] => [
  {
    key: 'action',
    label: '操作类型',
    type: 'select',
    icon: Activity,
    placeholder: '全部操作',
    options: [
      ...Object.entries(AUDIT_LOG_ACTION_MAP).map(([key, label]) => ({
        value: key,
        label,
      })),
    ],
  },
  {
    key: 'resource',
    label: '资源类型',
    type: 'select',
    icon: Database,
    placeholder: '全部资源',
    options: [
      ...Object.entries(AUDIT_LOG_RESOURCE_MAP).map(([key, label]) => ({
        value: key,
        label,
      })),
    ],
  },
])

// 高级筛选配置
const advancedFilters = computed((): FilterField[] => [
  {
    key: 'user',
    type: 'user-search',
    label: '用户',
    icon: User,
    placeholder: '输入用户名',
  },
  {
    key: 'dateRange',
    type: 'date-range',
    label: '日期范围',
    icon: Calendar,
    placeholder: '选择日期范围',
  },
])

// 操作按钮配置
const actionButtons = computed((): ActionButton[] => [
  {
    text: '导出',
    icon: Download,
    variant: 'outline',
    permission: { resource: 'audit', action: 'export' },
    onClick: () => exportLogs(),
  },
  {
    text: '刷新',
    icon: RefreshCw,
    variant: 'default',
    permission: { resource: 'audit', action: 'read' },
    onClick: () => refreshLogs(),
  },
])

// 结果文本
const resultText = computed(() => {
  return `共找到 ${totalCount.value} 条审计日志`
})

async function loadLogs() {
  loading.value = true
  try {
    const auditParams = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: filters.value.searchQuery,
      action: filters.value.action,
      resource_type: filters.value.resource,
      user_id: filters.value.user,
      created_at_from: filters.value.dateRange?.[0],
      created_at_to: filters.value.dateRange?.[1],
    }

    const response = await auditAPI.getLogs(auditParams)
    logs.value = response.data || []
    totalCount.value = response.pagination.total
  }
  catch (error) {
    console.error('Failed to load audit logs:', error)
  }
  finally {
    loading.value = false
  }
}

function handleFilterChange(newFilters: FilterValues) {
  filters.value = newFilters
  currentPage.value = 1
  loadLogs()
}

function handleFilterReset() {
  filters.value = {
    searchQuery: '',
    action: '',
    resource: '',
    user: '',
    dateRange: undefined,
  }
  currentPage.value = 1
  loadLogs()
}

function toggleLogExpansion(logId: string) {
  if (expandedLogs.value.has(logId)) {
    expandedLogs.value.delete(logId)
  }
  else {
    expandedLogs.value.add(logId)
  }
}

function openLogDetail(log: AuditLog) {
  selectedLog.value = log
  isSheetOpen.value = true
}

function refreshLogs() {
  loadLogs()
}

function goToPage(page: number) {
  currentPage.value = page
  loadLogs()
}

async function exportLogs() {
  try {
    // 这里应该调用导出API，暂时用alert提示
    // alert('导出功能开发中...')
  }
  catch (error) {
    console.error('Failed to export logs:', error)
    // alert('导出失败，请重试')
  }
}

onMounted(() => {
  loadLogs()
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <PageHeader
      title="审计日志"
      description="查看系统操作记录和安全审计信息"
      :actions="actionButtons"
    />

    <Card>
      <CardContent class="px-6 space-y-6">
        <!-- 筛选器 -->
        <DataFilter
          v-model="filters"
          :total-items="totalCount"
          search-placeholder="搜索审计日志..."
          :quick-filters="quickFilters"
          :advanced-filters="advancedFilters"
          :result-text="resultText"
          @update:model-value="handleFilterChange"
          @reset="handleFilterReset"
        />

        <!-- 日志列表 -->
        <div v-if="loading" class="flex items-center justify-center py-8">
          <Loader2 class="h-6 w-6 animate-spin" />
          <span class="ml-2">加载中...</span>
        </div>

        <EmptyState
          v-else-if="logs.length === 0"
          :icon="FileText"
          title="暂无日志记录"
          description="当前筛选条件下没有找到审计日志"
        />

        <div v-else class="space-y-4">
          <div
            v-for="log in logs"
            :key="log.id"
            class="border border-border rounded-lg p-4 hover:bg-muted/50 transition-all duration-200 hover:shadow-md cursor-pointer"
            @click="openLogDetail(log)"
          >
            <div class="flex items-start justify-between">
              <div class="flex items-start space-x-4 flex-1">
                <div class="flex-shrink-0">
                  <div
                    class="w-10 h-10 rounded-full flex items-center justify-center"
                    :class="getActionColor(log.action)"
                  >
                    <component :is="getActionIcon(log.action)" class="h-5 w-5" />
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center space-x-2 mb-2">
                    <h4 class="text-sm font-semibold">
                      {{ getActionDisplayName(log.action) }}
                    </h4>
                  </div>

                  <!-- 简要信息 - 始终显示 -->
                  <div class="mb-2">
                    <p class="text-sm text-foreground">
                      <strong>{{ log.user?.name || '未知用户' }}</strong>
                      对 <strong>{{ getResourceDisplayName(log.resource) }}</strong>
                      <span v-if="log.resource_id"> "{{ log.resource_id }}"</span>
                      <span v-else> 列表</span>
                      执行了 <strong>{{ getActionDisplayName(log.action) }}</strong> 操作
                    </p>

                    <div class="flex items-center space-x-4 text-xs text-muted-foreground mt-1">
                      <span class="flex items-center space-x-1">
                        <Clock class="h-3 w-3" />
                        <span>{{ formatDate(log.created_at) }}</span>
                      </span>
                      <span class="flex items-center space-x-1">
                        <Globe class="h-3 w-3" />
                        <span>{{ log.ip_address || '未知IP' }}</span>
                      </span>
                    </div>
                  </div>

                  <!-- 详细信息 - 可展开 -->
                  <div v-if="expandedLogs.has(log.id)" class="mb-3 border-t pt-3">
                    <div class="space-y-3">
                      <!-- 操作详情 -->
                      <div v-if="log.details" class="text-sm">
                        <strong class="text-foreground">操作详情：</strong>
                        <span class="text-muted-foreground">{{ log.details }}</span>
                      </div>

                      <!-- 资源信息 -->
                      <div v-if="log.resource_id" class="text-sm">
                        <strong class="text-foreground">资源信息：</strong>
                        <div class="text-muted-foreground ml-4">
                          <div v-if="log.resource_id">
                            ID: {{ log.resource_id }}
                          </div>
                        </div>
                      </div>

                      <!-- 数据变更 -->
                      <div v-if="log.details" class="text-sm">
                        <strong class="text-foreground">操作详情：</strong>
                        <span class="text-muted-foreground">{{ log.details }}</span>
                      </div>

                      <!-- 技术信息 -->
                      <div class="text-sm">
                        <strong class="text-foreground">技术信息：</strong>
                        <div class="grid grid-cols-2 gap-2 mt-1 text-xs text-muted-foreground">
                          <div v-if="log.user?.role" class="flex items-center space-x-1">
                            <User class="h-3 w-3" />
                            <span>角色: {{ log.user.role }}</span>
                          </div>
                          <div class="flex items-center space-x-1">
                            <Monitor class="h-3 w-3" />
                            <span>设备: {{ log.user_agent ? getUserAgentInfo(log.user_agent) : '未知设备' }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="flex-shrink-0 flex flex-col items-end space-y-2">
                <div class="text-right">
                  <p class="text-xs text-muted-foreground">
                    {{ formatRelativeTime(log.created_at) }}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  class="h-8 w-8 p-0"
                  @click="toggleLogExpansion(log.id)"
                >
                  <ChevronRight
                    class="h-4 w-4 transition-transform duration-200" :class="[
                      expandedLogs.has(log.id) ? 'rotate-90' : '',
                    ]"
                  />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="flex items-center justify-between mt-6">
          <Pagination
            :page="currentPage"
            :total="totalCount || 0"
            :items-per-page="pageSize || 20"
            :sibling-count="1"
            :show-edges="true"
            @update:page="goToPage"
          >
            <PaginationContent v-slot="{ items }">
              <PaginationFirst />
              <PaginationPrevious />
              <template v-for="(item, index) in items" :key="index">
                <PaginationItem
                  v-if="item.type === 'page'"
                  :value="item.value"
                  :is-active="item.value === currentPage"
                >
                  {{ item.value }}
                </PaginationItem>
                <PaginationEllipsis v-else :index="index">
                  &#8230;
                </PaginationEllipsis>
              </template>
              <PaginationNext />
              <PaginationLast />
            </PaginationContent>
          </Pagination>
        </div>
      </CardContent>
    </Card>

    <!-- 详细信息Sheet -->
    <Sheet v-model:open="isSheetOpen">
      <SheetContent class="w-[600px] sm:w-[800px] overflow-y-auto p-0">
        <LogDetail v-if="selectedLog" v-model="selectedLog" />

        <!-- 底部固定操作栏 -->
        <div class="sticky bottom-0 z-10 px-6 py-4 flex justify-end">
          <SheetClose class="w-full" as-child>
            <Button variant="outline">
              关闭
            </Button>
          </SheetClose>
        </div>
      </SheetContent>
    </Sheet>
  </div>
</template>
