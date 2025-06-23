<script setup lang="ts">
import type { User } from '@/api/types'
import {
  Activity,
  Plus,
  Shield,
  Users,
} from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { userAPI } from '@/api/user'
import { DataFilter, DataTable, type FilterField, PermissionButton, PermissionWrapper } from '@/components'
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
import {
  Card,
  CardContent,
} from '@/components/ui/card'
import { generateColumns } from './columns'
import UserSheet from './UserSheet.vue'

// 页面信息
const pageInfo = {
  title: '用户管理',
  description: '管理系统用户账户和权限',
  searchPlaceholder: '用户名、姓名或邮箱',
  createButtonText: '新建用户',
  emptyIcon: Users,
  emptyTitle: '暂无用户',
  emptyDescription: '开始创建您的第一个用户账户',
}

// 筛选字段配置
const quickFilters: FilterField[] = [
  {
    key: 'role',
    label: '角色',
    type: 'select',
    icon: Shield,
    placeholder: '选择角色',
    options: [
      { value: 'super_admin', label: '超级管理员' },
      { value: 'sec_mgr', label: '安全管理员' },
      { value: 'dev', label: '开发人员' },
      { value: 'auditor', label: '审计员' },
      { value: 'bot', label: '服务账号' },
    ],
  },
  {
    key: 'status',
    label: '状态',
    type: 'select',
    icon: Activity,
    placeholder: '选择状态',
    options: [
      { value: 'active', label: '活跃' },
      { value: 'disabled', label: '禁用' },
      { value: 'locked', label: '锁定' },
      { value: 'expired', label: '过期' },
    ],
  },
]

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(0)

// 筛选状态
const filters = ref({
  searchQuery: undefined,
  role: undefined,
  status: undefined,
})

// Sheet 相关
const showUserSheet = ref(false)
const sheetMode = ref<'create' | 'edit' | 'view'>('create')
const selectedUser = ref<User | null>(null)
const showDeleteDialog = ref(false)
const userToDelete = ref<User | null>(null)

const columns = computed(() => generateColumns({
  onView: handleViewUser,
  onEdit: handleEditUser,
  onDelete: openDeleteDialog,
}))

// 获取用户列表
async function fetchUsers() {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      name: filters.value.searchQuery,
      role: filters.value.role,
      status: filters.value.status,
    }

    const response = await userAPI.getUsers(params)
    users.value = response.data || []
    total.value = response.pagination.total
    totalPages.value = response.pagination.total_pages
  }
  catch (error) {
    console.error('获取用户列表失败:', error)
    users.value = []
    total.value = 0
  }
  finally {
    loading.value = false
  }
}

// 事件处理
function handleFiltersChange() {
  currentPage.value = 1
  fetchUsers()
}

function handleFiltersReset() {
  filters.value = {
    searchQuery: undefined,
    role: undefined,
    status: undefined,
  }
  currentPage.value = 1
  fetchUsers()
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchUsers()
}

function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchUsers()
}

function handleCreateUser() {
  selectedUser.value = null
  sheetMode.value = 'create'
  showUserSheet.value = true
}

function handleViewUser(user: User) {
  selectedUser.value = user
  sheetMode.value = 'view'
  showUserSheet.value = true
}

function handleEditUser(user: User) {
  selectedUser.value = user
  sheetMode.value = 'edit'
  showUserSheet.value = true
}

function openDeleteDialog(user: User) {
  userToDelete.value = user
  showDeleteDialog.value = true
}

async function confirmDeleteUser() {
  if (!userToDelete.value)
    return

  try {
    await userAPI.deleteUser(userToDelete.value.id)
    await fetchUsers()
    showDeleteDialog.value = false
    userToDelete.value = null
  }
  catch (error) {
    console.error('删除用户失败:', error)
  }
}

function handleUserSaved() {
  showUserSheet.value = false
  fetchUsers()
}

// 监听筛选变化
watch(filters, () => {
  handleFiltersChange()
}, { deep: true })

// 组件挂载时获取数据
onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">
          {{ pageInfo.title }}
        </h1>
        <p class="text-muted-foreground mt-2">
          {{ pageInfo.description }}
        </p>
      </div>
      <PermissionButton
        :permission="{ resource: 'user', action: 'create' }"
        @click="handleCreateUser"
      >
        <Plus class=" h-4 w-4" />
        {{ pageInfo.createButtonText }}
      </PermissionButton>
    </div>

    <!-- 搜索和筛选 -->
    <PermissionWrapper resource="user" action="read">
      <Card>
        <CardContent class="px-6 space-y-6">
          <DataFilter
            v-model="filters"
            :total-items="total"
            :search-placeholder="pageInfo.searchPlaceholder"
            :quick-filters="quickFilters"
            @reset="handleFiltersReset"
          />
          <DataTable
            v-model:current-page="currentPage"
            :columns="columns"
            :data="users"
            :loading="loading"
            :total="total"
            :page-size="pageSize"
            :total-pages="totalPages"
            :empty-icon="pageInfo.emptyIcon"
            :empty-title="pageInfo.emptyTitle"
            :empty-description="pageInfo.emptyDescription"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          />
        </CardContent>
      </Card>
      <template #fallback>
        <Card>
          <CardContent class="flex items-center justify-center py-12">
            <div class="text-center space-y-2">
              <Users class="h-12 w-12 mx-auto text-muted-foreground" />
              <h3 class="text-lg font-medium">
                权限不足
              </h3>
              <p class="text-muted-foreground">
                您没有权限查看用户管理功能
              </p>
            </div>
          </CardContent>
        </Card>
      </template>
    </PermissionWrapper>

    <!-- 用户详情/编辑 Sheet -->
    <UserSheet
      v-model:open="showUserSheet"
      :mode="sheetMode"
      :user="selectedUser"
      @saved="handleUserSaved"
    />

    <!-- 删除确认对话框 -->
    <AlertDialog v-model:open="showDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除用户 "{{ userToDelete?.name }}" 吗？此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">
            取消
          </AlertDialogCancel>
          <AlertDialogAction class="bg-destructive hover:bg-destructive/90" @click="confirmDeleteUser">
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
