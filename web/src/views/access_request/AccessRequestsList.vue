<script setup lang="ts">
import type { AccessRequest } from '@/api/types'
import {
  ArrowUpDown,
  Calendar,
  CircleAlert,
  FileX,
  Plus,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { accessRequestAPI } from '@/api/access-request'
import { DataFilter, DataTable, EmptyState, type FilterField, PageHeader } from '@/components'
import { Card, CardContent } from '@/components/ui/card'
import ApproveRequestDialog from '@/views/access_request/ApproveRequestDialog.vue'
import { accessRequestColumns } from '@/views/access_request/columns'
import RejectRequestDialog from '@/views/access_request/RejectRequestDialog.vue'
import RevokeRequestDialog from '@/views/access_request/RevokeRequestDialog.vue'

const router = useRouter()

// 筛选状态接口
interface FilterState {
  searchQuery?: string
  selectedStatus?: string
  searchApplicant?: string
  dateRange?: [number, number]
  sortBy: string
}

function handleCreateRequest() {
  router.push('/access-request/create')
}

// Filter相关配置
const quickFilters = computed((): FilterField[] => [
  {
    key: 'selectedStatus',
    label: '状态',
    type: 'select',
    icon: CircleAlert,
    placeholder: '全部状态',
    options: [
      { value: 'pending', label: '待审批' },
      { value: 'approved', label: '已批准' },
      { value: 'rejected', label: '已拒绝' },
      { value: 'expired', label: '已过期' },
      { value: 'revoked', label: '已撤销' },
    ],
  },
])

const advancedFilters = computed((): FilterField[] => [
  {
    key: 'sortBy',
    label: '排序方式',
    type: 'select',
    icon: ArrowUpDown,
    placeholder: '排序方式',
    options: [
      { value: 'created_at_desc', label: '申请时间 ↓' },
      { value: 'created_at_asc', label: '申请时间 ↑' },
      { value: 'updated_at_desc', label: '更新时间 ↓' },
      { value: 'updated_at_asc', label: '更新时间 ↑' },
      { value: 'status_asc', label: '状态 A-Z' },
      { value: 'status_desc', label: '状态 Z-A' },
    ],
  },
  {
    key: 'searchApplicant',
    label: '申请人',
    type: 'user-search',
    icon: User,
    placeholder: '搜索申请人名称',
  },
  {
    key: 'dateRange',
    label: '申请日期范围',
    type: 'date-range',
    icon: Calendar,
    placeholder: '选择申请日期范围',
  },
])

// 响应式数据
const requests = ref<AccessRequest[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(0)
const totalItems = ref(0)

const resultText = computed(() => `共找到 ${totalItems.value} 个申请`)

// 对话框状态
const showApproveDialog = ref(false)
const showRejectDialog = ref(false)
const showRevokeDialog = ref(false)
const selectedRequest = ref<AccessRequest | null>(null)

// 筛选状态
const filterState = ref<FilterState>({
  searchQuery: undefined,
  selectedStatus: undefined,
  searchApplicant: undefined,
  dateRange: undefined,
  sortBy: 'created_at_desc',
})

const filter = computed(() => {
  const baseFilter = {
    search: filterState.value.searchQuery,
    status: filterState.value.selectedStatus,
    applicant_name: filterState.value.searchApplicant,
    created_at_from: filterState.value.dateRange?.[0],
    created_at_to: filterState.value.dateRange?.[1],
    sort_by: filterState.value.sortBy,
    page: currentPage.value,
    page_size: Number.parseInt(pageSize.value.toString()),
  }

  // 移除undefined值
  return Object.fromEntries(
    Object.entries(baseFilter).filter(([_, value]) => value !== undefined),
  )
})

async function loadRequests() {
  try {
    loading.value = true
    const response = await accessRequestAPI.getRequests(filter.value)
    requests.value = response.data || []
    totalPages.value = response.pagination.total_pages || 0
    totalItems.value = response.pagination.total || 0
  }
  catch (error: any) {
    console.error('加载申请列表失败:', error)
    toast.error(error.response?.data?.error || '加载申请列表失败')
  }
  finally {
    loading.value = false
  }
}

watch(filter, () => {
  loadRequests()
})

function resetFilters() {
  filterState.value = {
    searchQuery: undefined,
    selectedStatus: undefined,
    searchApplicant: undefined,
    dateRange: undefined,
    sortBy: 'created_at_desc',
  }
  currentPage.value = 1
}

function handlePageChange(page: number) {
  currentPage.value = page
}

function approveRequest(request: AccessRequest) {
  selectedRequest.value = request
  showApproveDialog.value = true
}

function rejectRequest(request: AccessRequest) {
  selectedRequest.value = request
  showRejectDialog.value = true
}

function revokeRequest(request: AccessRequest) {
  selectedRequest.value = request
  showRevokeDialog.value = true
}

function handleApproveSuccess() {
  showApproveDialog.value = false
  selectedRequest.value = null
  loadRequests()
}

function handleRejectSuccess() {
  showRejectDialog.value = false
  selectedRequest.value = null
  loadRequests()
}

function handleRevokeSuccess() {
  showRevokeDialog.value = false
  selectedRequest.value = null
  loadRequests()
}

// 监听筛选条件变化，重置分页
watch(filterState, () => {
  currentPage.value = 1
}, { deep: true })

// 监听页面大小变化，重置分页
watch(pageSize, () => {
  currentPage.value = 1
})

onMounted(() => {
  loadRequests()

  // 监听批准事件
  window.addEventListener('approve-request', (event: any) => {
    approveRequest(event.detail)
  })

  // 监听拒绝事件
  window.addEventListener('reject-request', (event: any) => {
    rejectRequest(event.detail)
  })

  // 监听作废事件
  window.addEventListener('revoke-request', (event: any) => {
    revokeRequest(event.detail)
  })
})
</script>

<template>
  <div class="space-y-6">
    <PageHeader
      title="访问申请管理"
      description="管理和审批用户的密钥访问申请"
      button-text="创建申请"
      :button-icon="Plus"
      @button-click="handleCreateRequest"
    />

    <Card class="border-none shadow-sm bg-gradient-to-br from-background to-muted/30">
      <CardContent class="px-6 space-y-6">
        <!-- 筛选和搜索 -->
        <DataFilter
          v-model="filterState"
          :total-items="totalItems"
          search-placeholder="搜索申请人、密钥名称或申请理由..."
          :quick-filters="quickFilters"
          :advanced-filters="advancedFilters"
          :result-text="resultText"
          @update:model-value="filterState = $event"
          @reset="resetFilters"
        />

        <!-- 列表 -->
        <div>
          <DataTable
            v-if="!loading && requests.length > 0"
            v-model:current-page="currentPage"
            :columns="accessRequestColumns"
            :data="requests"
            :page-count="totalPages"
            :page-size="pageSize"
            :loading="loading"
            @page-change="handlePageChange"
            @approve-request="approveRequest"
            @reject-request="rejectRequest"
          />

          <!-- 空状态 -->
          <EmptyState
            v-else-if="!loading"
            :icon="FileX"
            title="暂无访问申请"
            description="还没有任何访问申请记录。您可以创建新的申请来开始使用。"
            create-button-text="创建申请"
            :create-button-icon="Plus"
            @create="handleCreateRequest"
          />
        </div>
      </CardContent>
    </Card>

    <!-- 批准对话框 -->
    <ApproveRequestDialog
      v-model:open="showApproveDialog"
      :request="selectedRequest!"
      @success="handleApproveSuccess"
    />

    <!-- 拒绝对话框 -->
    <RejectRequestDialog
      v-model:open="showRejectDialog"
      :request="selectedRequest!"
      @success="handleRejectSuccess"
    />

    <!-- 作废对话框 -->
    <RevokeRequestDialog
      v-model:open="showRevokeDialog"
      :request="selectedRequest!"
      @success="handleRevokeSuccess"
    />
  </div>
</template>
