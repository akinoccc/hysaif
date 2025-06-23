<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import {
  ArrowUpDown,
  Calendar,
  Filter,
  Globe,
  Hash,
  Plus,
  Tag,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { secretItemAPI } from '@/api/secret'
import AccessRequestDialog from '@/components/business/secret-list/AccessRequestDialog.vue'
import DataTable from '@/components/common/data-table/DataTable.vue'
import { DataFilter, EmptyState, type FilterField, PageHeader } from '@/components/common/layout'
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
import { Card, CardContent } from '@/components/ui/card'
import { SECRET_ITEM_TYPE, SECRET_ITEM_TYPE_MAP, type SecretItemTypeT } from '@/constants'
import { getCategoriesByGroup, getCategoryGroupsByType, type SecretItemCategoryGroupKey } from '@/constants'
import { generateColumns } from './columns'

const props = defineProps<{
  secretType: SecretItemTypeT
}>()

const router = useRouter()

// 根据类型获取相关信息
const typeInfo = computed(() => {
  const typeInfoMap = {
    [SECRET_ITEM_TYPE.ApiKey]: {
      title: 'API 密钥管理',
      description: '管理您的 API 密钥和访问凭据',
      searchPlaceholder: 'API 密钥名称、描述或标签',
      createButtonText: '新建 API 密钥',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.ApiKey].icon,
      emptyTitle: '暂无 API 密钥',
      emptyDescription: '开始创建您的第一个 API 密钥',
    },
    [SECRET_ITEM_TYPE.Password]: {
      title: '密码管理',
      description: '管理您的网站、应用和系统密码',
      searchPlaceholder: '密码名称、描述或标签',
      createButtonText: '新建密码',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Password].icon,
      emptyTitle: '暂无密码',
      emptyDescription: '开始创建您的第一个密码',
    },
    [SECRET_ITEM_TYPE.SshKey]: {
      title: 'SSH 密钥管理',
      description: '管理您的服务器和开发环境 SSH 密钥',
      searchPlaceholder: 'SSH 密钥名称、描述或标签',
      createButtonText: '新建 SSH 密钥',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.SshKey].icon,
      emptyTitle: '暂无 SSH 密钥',
      emptyDescription: '开始创建您的第一个 SSH 密钥',
    },
    [SECRET_ITEM_TYPE.Token]: {
      title: '令牌管理',
      description: '管理您的访问令牌和认证凭据',
      searchPlaceholder: '令牌名称、描述或标签',
      createButtonText: '新建令牌',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Token].icon,
      emptyTitle: '暂无令牌',
      emptyDescription: '开始创建您的第一个令牌',
    },
    [SECRET_ITEM_TYPE.AccessKey]: {
      title: '访问密钥管理',
      description: '管理您的云服务和平台访问密钥',
      searchPlaceholder: '访问密钥名称、描述或标签',
      createButtonText: '新建访问密钥',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.AccessKey].icon,
      emptyTitle: '暂无访问密钥',
      emptyDescription: '开始创建您的第一个访问密钥',
    },
    [SECRET_ITEM_TYPE.Custom]: {
      title: '自定义密钥管理',
      description: '管理您的自定义密钥和配置',
      searchPlaceholder: '自定义密钥名称、描述或标签',
      createButtonText: '新建自定义密钥',
      emptyIcon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Custom].icon,
      emptyTitle: '暂无自定义密钥',
      emptyDescription: '开始创建您的第一个自定义密钥',
    },
  }

  return typeInfoMap[props.secretType] || typeInfoMap[SECRET_ITEM_TYPE.Custom]
})

// 根据类型获取列定义
const columns = computed(() => {
  return generateColumns(props.secretType)
})

// 筛选状态接口
interface FilterState {
  searchQuery?: string
  selectedCategory?: string
  selectedEnvironment?: string
  selectedStatus?: string
  searchCreator?: string
  dateRange?: [number, number]
  searchTags?: string
  sortBy: string
}

// Header相关方法
function handleCreate() {
  router.push(`/${props.secretType}/create`)
}

// Filter相关配置
const availableGroups = computed(() => {
  return getCategoryGroupsByType(props.secretType)
})

function getFilteredCategoriesByGroup(groupKey: SecretItemCategoryGroupKey) {
  return getCategoriesByGroup(groupKey, props.secretType)
}

const quickFilters = computed((): FilterField[] => [
  {
    key: 'selectedCategory',
    label: '分类',
    icon: Filter,
    type: 'select',
    placeholder: '全部分类',
    options: availableGroups.value.flatMap(group =>
      getFilteredCategoriesByGroup(group.key).map(category => ({
        value: category.key,
        label: category.label,
        group: group.label,
      })),
    ),
  },
  {
    key: 'selectedStatus',
    label: '状态',
    icon: Hash,
    type: 'select',
    placeholder: '全部状态',
    options: [
      { value: 'active', label: '正常' },
      { value: 'expiring', label: '即将过期' },
      { value: 'expired', label: '已过期' },
    ],
  },
  {
    key: 'selectedEnvironment',
    label: '环境',
    icon: Globe,
    type: 'select',
    placeholder: '全部环境',
    options: [
      { value: 'production', label: '生产环境' },
      { value: 'staging', label: '测试环境' },
      { value: 'development', label: '开发环境' },
      { value: 'sandbox', label: '沙盒环境' },
      { value: 'unspecified', label: '未指定' },
    ],
  },
])

const advancedFilters = computed((): FilterField[] => [
  {
    key: 'sortBy',
    type: 'select',
    label: '排序方式',
    icon: ArrowUpDown,
    placeholder: '排序方式',
    options: [
      { value: 'created_at_desc', label: '创建时间 ↓' },
      { value: 'created_at_asc', label: '创建时间 ↑' },
      { value: 'updated_at_desc', label: '更新时间 ↓' },
      { value: 'updated_at_asc', label: '更新时间 ↑' },
      { value: 'name_asc', label: '名称 A-Z' },
      { value: 'name_desc', label: '名称 Z-A' },
    ],
  },
  {
    key: 'searchCreator',
    type: 'user-search',
    label: '创建者',
    icon: User,
    placeholder: '搜索创建者名称',
  },
  {
    key: 'searchTags',
    type: 'text',
    label: '标签',
    icon: Tag,
    placeholder: '搜索标签',
  },
  {
    key: 'dateRange',
    type: 'date-range',
    label: '创建日期范围',
    icon: Calendar,
    placeholder: '选择创建日期范围',
  },
])

// 响应式数据
const items = ref<SecretItem[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(0)
const totalItems = ref(0)
const showDeleteDialog = ref(false)
const itemToDelete = ref<SecretItem | null>(null)

const resultText = computed(() => {
  return `共找到 ${totalItems.value} 个${typeInfo.value.searchPlaceholder}`
})

// 申请访问对话框状态
const showAccessRequestDialog = ref(false)
const selectedItemForRequest = ref<SecretItem | null>(null)

// 筛选状态
const filterState = ref<FilterState>({
  searchQuery: undefined,
  selectedCategory: undefined,
  selectedEnvironment: undefined,
  selectedStatus: undefined,
  searchCreator: undefined,
  dateRange: undefined,
  searchTags: undefined,
  sortBy: 'created_at_desc',
})

// 计算属性
const filter = computed(() => {
  const baseFilter = {
    type: props.secretType,
    search: filterState.value.searchQuery,
    category: filterState.value.selectedCategory,
    environment: filterState.value.selectedEnvironment,
    status: filterState.value.selectedStatus,
    creator_name: filterState.value.searchCreator,
    created_at_from: filterState.value.dateRange?.[0],
    created_at_to: filterState.value.dateRange?.[1],
    tags: filterState.value.searchTags ? filterState.value.searchTags.split(',').map(tag => tag.trim()).filter(Boolean) : undefined,
    sort_by: filterState.value.sortBy,
    page: currentPage.value,
    page_size: Number.parseInt(pageSize.value.toString()),
  }

  // 移除undefined值
  return Object.fromEntries(
    Object.entries(baseFilter).filter(([_, value]) => value !== undefined),
  )
})

// 方法
async function loadItems() {
  try {
    loading.value = true
    const response = await secretItemAPI.getItems(filter.value)
    items.value = response.data || []
    totalPages.value = response.pagination.total_pages || 0
    totalItems.value = response.pagination.total || 0
  }
  catch (error) {
    console.error(`加载${typeInfo.value.title}失败:`, error)
  }
  finally {
    loading.value = false
  }
}

watch(filter, () => {
  loadItems()
})

function resetFilters() {
  filterState.value = {
    searchQuery: undefined,
    selectedCategory: undefined,
    selectedEnvironment: undefined,
    selectedStatus: undefined,
    searchCreator: undefined,
    dateRange: undefined,
    searchTags: undefined,
    sortBy: 'created_at_desc',
  }
  currentPage.value = 1
}

// 处理分页变化
function handlePageChange(page: number) {
  currentPage.value = page
}

function viewItem(item: SecretItem) {
  router.push(`/${props.secretType}/${item.id}`)
}

function viewItemWithAccess(item: SecretItem) {
  // 通过访问申请查看密钥，使用特殊路由参数标识
  router.push(`/${props.secretType}/${item.id}?access=true`)
}

function editItem(item: SecretItem) {
  router.push(`/${props.secretType}/${item.id}/edit`)
}

function requestAccess(item: SecretItem) {
  selectedItemForRequest.value = item
  showAccessRequestDialog.value = true
}

function handleAccessRequestSuccess() {
  showAccessRequestDialog.value = false
  selectedItemForRequest.value = null
}

function openDeleteDialog(item: SecretItem) {
  itemToDelete.value = item
  showDeleteDialog.value = true
}

async function confirmDeleteItem() {
  if (!itemToDelete.value)
    return

  try {
    await secretItemAPI.deleteItem(itemToDelete.value.id)
    await loadItems()
    showDeleteDialog.value = false
    itemToDelete.value = null
  }
  catch (error) {
    console.error('删除失败:', error)
  }
}

// 监听筛选条件变化，重置分页
watch(filterState, () => {
  currentPage.value = 1
}, { deep: true })

// 监听页面大小变化，重置分页
watch(pageSize, () => {
  currentPage.value = 1
})

// 监听类型变化，重置所有筛选条件并重新加载数据
watch(() => props.secretType, () => {
  resetFilters()
  loadItems()
})

// 组件挂载时加载数据
onMounted(() => {
  loadItems()
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题和操作 -->
    <PageHeader
      :title="typeInfo.title"
      :description="typeInfo.description"
      :button-text="typeInfo.createButtonText"
      :button-icon="Plus"
      :create-button="{
        text: typeInfo.createButtonText,
        permission: {
          resource: props.secretType,
          action: 'create',
        },
        onClick: handleCreate,
      }"
    />

    <Card>
      <CardContent class="px-6 space-y-6">
        <!-- 筛选和搜索 -->
        <DataFilter
          v-model="filterState"
          :total-items="totalItems"
          :search-placeholder="`搜索${typeInfo.searchPlaceholder}...`"
          :quick-filters="quickFilters"
          :advanced-filters="advancedFilters"
          :result-text="resultText"
          @update:model-value="filterState = $event"
          @reset="resetFilters"
        />

        <!-- 列表 -->
        <div>
          <DataTable
            v-if="!loading && items.length > 0"
            v-model:current-page="currentPage"
            :columns="columns"
            :data="items"
            :page-count="totalPages"
            :page-size="pageSize"
            @page-change="handlePageChange"
            @view-item="viewItem"
            @view-item-with-access="viewItemWithAccess"
            @edit-item="editItem"
            @delete-item="openDeleteDialog"
            @request-access="requestAccess"
          />

          <!-- 空状态 -->
          <EmptyState
            v-else-if="!loading"
            :icon="typeInfo.emptyIcon"
            :title="typeInfo.emptyTitle"
            :description="typeInfo.emptyDescription"
            :create-button-text="typeInfo.createButtonText"
            :create-button-icon="Plus"
            @create="handleCreate"
          />
        </div>
      </CardContent>
    </Card>

    <!-- 申请访问对话框 -->
    <AccessRequestDialog
      v-model:open="showAccessRequestDialog"
      :item="selectedItemForRequest!"
      @success="handleAccessRequestSuccess"
    />

    <!-- 删除确认对话框 -->
    <AlertDialog v-model:open="showDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除 "{{ itemToDelete?.name }}" 吗？此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">
            取消
          </AlertDialogCancel>
          <AlertDialogAction class="bg-destructive hover:bg-destructive/90" @click="confirmDeleteItem">
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
