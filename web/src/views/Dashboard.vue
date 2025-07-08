<script setup lang="ts">
import {
  Activity,
  AlertTriangle,
  ArrowRight,
  BarChart3,
  CheckCircle,
  Clock,
  Database,
  Eye,
  FileText,
  Key,
  Shield,
  Users,
  XCircle,
} from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { type AccessRequest, accessRequestAPI, auditAPI, type AuditLog, type SecretItem, secretItemAPI, userAPI } from '@/api'
import { PermissionButton } from '@/components'
import { Card } from '@/components/ui/card'
import { AUDIT_LOG_ACTION_MAP, AUDIT_LOG_RESOURCE_MAP } from '@/constants/auditLog'
import { formatRelativeTime, getFileIcon } from '@/lib/utils'
import { usePermissionStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { getActionColor } from './audit/helper'

const router = useRouter()
const authStore = useAuthStore()
const { isAdmin, isSecurityManager, isDeveloper, isAuditor, currentRole } = storeToRefs(usePermissionStore())
const { hasPermission, roleChecker } = usePermissionStore()
const { hasRole } = roleChecker

// 基础统计数据
const stats = ref({
  totalItems: 0,
  expiringItems: 0,
  todayAccess: 0,
  securityScore: 95,
})

// 管理员专用统计
const adminStats = ref({
  totalUsers: 0,
  activeUsers: 0,
  pendingRequests: 0,
  systemHealth: 98,
})

// 安全管理员专用统计
const securityStats = ref({
  riskItems: 0,
  pendingApprovals: 0,
  securityAlerts: 0,
  complianceScore: 92,
})

// 开发人员专用统计
const devStats = ref({
  myItems: 0,
  myRequests: 0,
  accessibleItems: 0,
  requestsApproved: 0,
})

// 审计员专用统计
const auditStats = ref({
  todayLogs: 0,
  criticalLogs: 0,
  complianceIssues: 0,
  auditScore: 96,
})

const recentItems = ref<SecretItem[]>([])
const recentLogs = ref<AuditLog[]>([])
const recentRequests = ref<AccessRequest[]>([])
const recentUsers = ref<any[]>([])

// 权限检查（使用同步检查，提供回退值为true以避免空白页面）
const canViewAudit = computed(() =>
  hasPermission('audit', 'read', true)
  || isAdmin.value
  || isSecurityManager.value
  || isDeveloper.value
  || isAuditor.value,
)

const canViewUsers = computed(() =>
  hasPermission('user', 'read', true)
  || isAdmin.value,
)

const canViewRequests = computed(() =>
  hasPermission('access_request', 'read', true)
  || isAdmin.value
  || isSecurityManager.value
  || isDeveloper.value,
)

const canViewSecrets = computed(() =>
  hasPermission('secret', 'read', true)
  || isAdmin.value
  || isSecurityManager.value
  || isDeveloper.value,
)

function getActionDisplayName(action: string) {
  return AUDIT_LOG_ACTION_MAP[action as keyof typeof AUDIT_LOG_ACTION_MAP] || action
}

function getResourceDisplayName(resource: string) {
  return AUDIT_LOG_RESOURCE_MAP[resource as keyof typeof AUDIT_LOG_RESOURCE_MAP] || resource
}

// 根据角色获取仪表板标题
const getDashboardTitle = computed(() => {
  switch (currentRole.value) {
    case 'super_admin':
      return '系统管理仪表板'
    case 'sec_mgr':
      return '安全管理仪表板'
    case 'dev':
      return '开发者仪表板'
    case 'auditor':
      return '审计仪表板'
    default:
      return '仪表板'
  }
})

// 根据角色获取描述
const getRoleDescription = computed(() => {
  switch (currentRole.value) {
    case 'super_admin':
      return '系统管理看板'
    case 'sec_mgr':
      return '系统安全和风险管理'
    case 'dev':
      return '开发工作相关的资源'
    case 'auditor':
      return '系统审计和合规监控'
    default:
      return ''
  }
})

// 加载管理员专用数据
async function loadAdminData() {
  try {
    if (canViewUsers.value) {
      const usersResponse = await userAPI.getUsers({ page: 1, page_size: 10 })
      adminStats.value.totalUsers = usersResponse.pagination?.total || 0
      adminStats.value.activeUsers = usersResponse.data?.filter(user => user.status === 'active').length || 0
      recentUsers.value = usersResponse.data?.slice(0, 5) || []
    }

    if (canViewRequests.value) {
      const requestsResponse = await accessRequestAPI.getRequests({ page: 1, page_size: 5, status: 'pending' })
      adminStats.value.pendingRequests = requestsResponse.pagination?.total || 0
      recentRequests.value = requestsResponse.data || []
    }
  }
  catch (error) {
    console.error('Failed to load admin data:', error)
  }
}

// 加载安全管理员专用数据
async function loadSecurityData() {
  try {
    if (canViewSecrets.value) {
      // 计算风险项目数量（即将过期的项目）
      const itemsResponse = await secretItemAPI.getItems({ page: 1, page_size: 100 })
      const now = new Date()
      const thirtyDaysLater = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)

      securityStats.value.riskItems = itemsResponse.data?.filter((item: SecretItem) => {
        if (!item.expires_at)
          return false
        const expiresAt = new Date(item.expires_at)
        return expiresAt <= thirtyDaysLater
      }).length || 0
    }

    if (canViewRequests.value) {
      const pendingResponse = await accessRequestAPI.getRequests({ page: 1, page_size: 5, status: 'pending' })
      securityStats.value.pendingApprovals = pendingResponse.pagination?.total || 0
      recentRequests.value = pendingResponse.data || []
    }

    if (canViewAudit.value) {
      const today = new Date().toISOString().split('T')[0]
      const logsResponse = await auditAPI.getLogs({ page: 1, page_size: 10 })
      securityStats.value.securityAlerts = logsResponse.data?.filter((log: AuditLog) =>
        log.created_at.startsWith(today) && (log.action === 'delete' || log.action === 'create'),
      ).length || 0
    }
  }
  catch (error) {
    console.error('Failed to load security data:', error)
  }
}

// 加载开发人员专用数据
async function loadDevData() {
  try {
    if (canViewSecrets.value) {
      // 获取用户可访问的项目
      const itemsResponse = await secretItemAPI.getAccessedItems({ page: 1, page_size: 5 })
      devStats.value.accessibleItems = itemsResponse.pagination?.total || 0
      recentItems.value = itemsResponse.data || []
    }

    if (canViewRequests.value) {
      // 获取用户的访问请求
      const requestsResponse = await accessRequestAPI.getRequests({ page: 1, page_size: 5 })
      const userRequests = requestsResponse.data?.filter(req => req.applicant.id === authStore.user?.id) || []
      devStats.value.myRequests = userRequests.length
      devStats.value.requestsApproved = userRequests.filter(req => req.status === 'approved').length
      recentRequests.value = userRequests
    }
  }
  catch (error) {
    console.error('Failed to load dev data:', error)
  }
}

// 加载审计员专用数据
async function loadAuditData() {
  try {
    if (canViewAudit.value) {
      const today = new Date().toISOString().split('T')[0]
      const logsResponse = await auditAPI.getLogs({ page: 1, page_size: 5 })
      recentLogs.value = logsResponse.data || []

      auditStats.value.todayLogs = recentLogs.value.filter((log: AuditLog) =>
        log.created_at.startsWith(today),
      ).length

      auditStats.value.criticalLogs = recentLogs.value.filter((log: AuditLog) =>
        log.action === 'delete' || log.action === 'update',
      ).length
    }
  }
  catch (error) {
    console.error('Failed to load audit data:', error)
  }
}

// 根据角色加载对应数据
async function loadDashboardData() {
  // 根据角色加载专用数据
  if (isAdmin.value) {
    await loadAdminData()
  }
  else if (isSecurityManager.value) {
    await loadSecurityData()
  }
  else if (isDeveloper.value) {
    await loadDevData()
  }
  else if (isAuditor.value) {
    await loadAuditData()
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">
          {{ getDashboardTitle }}
        </h1>
        <p class="text-muted-foreground mt-2">
          欢迎回来，{{ authStore.user?.name }}！{{ getRoleDescription }}
        </p>
      </div>
    </div>

    <!-- 超级管理员统计卡片 -->
    <div v-if="hasRole('super_admin')" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              总用户数
            </p>
            <p class="text-2xl font-bold">
              {{ adminStats.totalUsers }}
            </p>
          </div>
          <Users class="h-8 w-8 text-blue-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              活跃用户
            </p>
            <p class="text-2xl font-bold text-green-600">
              {{ adminStats.activeUsers }}
            </p>
          </div>
          <Activity class="h-8 w-8 text-green-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              待处理申请
            </p>
            <p class="text-2xl font-bold text-orange-600">
              {{ adminStats.pendingRequests }}
            </p>
          </div>
          <AlertTriangle class="h-8 w-8 text-orange-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              系统健康度
            </p>
            <p class="text-2xl font-bold text-green-600">
              {{ adminStats.systemHealth }}%
            </p>
          </div>
          <Shield class="h-8 w-8 text-green-500" />
        </div>
      </Card>
    </div>

    <!-- 安全管理员统计卡片 -->
    <div v-else-if="hasRole('sec_mgr')" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              风险项目
            </p>
            <p class="text-2xl font-bold text-red-600">
              {{ securityStats.riskItems }}
            </p>
          </div>
          <AlertTriangle class="h-8 w-8 text-red-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              待审批申请
            </p>
            <p class="text-2xl font-bold text-orange-600">
              {{ securityStats.pendingApprovals }}
            </p>
          </div>
          <Clock class="h-8 w-8 text-orange-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              安全警报
            </p>
            <p class="text-2xl font-bold text-yellow-600">
              {{ securityStats.securityAlerts }}
            </p>
          </div>
          <Shield class="h-8 w-8 text-yellow-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              合规评分
            </p>
            <p class="text-2xl font-bold text-green-600">
              {{ securityStats.complianceScore }}%
            </p>
          </div>
          <CheckCircle class="h-8 w-8 text-green-500" />
        </div>
      </Card>
    </div>

    <!-- 开发人员统计卡片 -->
    <div v-else-if="hasRole('dev')" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              可访问项目
            </p>
            <p class="text-2xl font-bold">
              {{ devStats.accessibleItems }}
            </p>
          </div>
          <Key class="h-8 w-8 text-blue-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              我的申请
            </p>
            <p class="text-2xl font-bold text-orange-600">
              {{ devStats.myRequests }}
            </p>
          </div>
          <FileText class="h-8 w-8 text-orange-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              已批准申请
            </p>
            <p class="text-2xl font-bold text-green-600">
              {{ devStats.requestsApproved }}
            </p>
          </div>
          <CheckCircle class="h-8 w-8 text-green-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              今日访问
            </p>
            <p class="text-2xl font-bold">
              {{ stats.todayAccess }}
            </p>
          </div>
          <Eye class="h-8 w-8 text-purple-500" />
        </div>
      </Card>
    </div>

    <!-- 审计员统计卡片 -->
    <div v-else-if="hasRole('auditor')" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              今日日志
            </p>
            <p class="text-2xl font-bold">
              {{ auditStats.todayLogs }}
            </p>
          </div>
          <BarChart3 class="h-8 w-8 text-blue-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              关键操作
            </p>
            <p class="text-2xl font-bold text-orange-600">
              {{ auditStats.criticalLogs }}
            </p>
          </div>
          <AlertTriangle class="h-8 w-8 text-orange-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              合规问题
            </p>
            <p class="text-2xl font-bold text-red-600">
              {{ auditStats.complianceIssues }}
            </p>
          </div>
          <XCircle class="h-8 w-8 text-red-500" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              审计评分
            </p>
            <p class="text-2xl font-bold text-green-600">
              {{ auditStats.auditScore }}%
            </p>
          </div>
          <CheckCircle class="h-8 w-8 text-green-500" />
        </div>
      </Card>
    </div>

    <!-- 默认统计卡片 -->
    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              总信息项
            </p>
            <p class="text-2xl font-bold">
              {{ stats.totalItems }}
            </p>
          </div>
          <Database class="h-8 w-8 text-muted-foreground" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              即将过期
            </p>
            <p class="text-2xl font-bold text-warning">
              {{ stats.expiringItems }}
            </p>
          </div>
          <Clock class="h-8 w-8 text-warning" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              今日访问
            </p>
            <p class="text-2xl font-bold">
              {{ stats.todayAccess }}
            </p>
          </div>
          <Eye class="h-8 w-8 text-muted-foreground" />
        </div>
      </Card>

      <Card class="p-6 theme-transition">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-muted-foreground">
              安全评分
            </p>
            <p class="text-2xl font-bold text-success">
              {{ stats.securityScore }}
            </p>
          </div>
          <Shield class="h-8 w-8 text-success" />
        </div>
      </Card>
    </div>

    <!-- 超级管理员面板 -->
    <div v-if="hasRole('super_admin')" class="grid gap-6 md:grid-cols-2">
      <!-- 最近用户 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              最近注册用户
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'user', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/users')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="user in recentUsers"
              :key="user.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm">
                  {{ user.name?.[0]?.toUpperCase() }}
                </div>
                <div>
                  <p class="font-medium">
                    {{ user.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    {{ user.role }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(user.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentUsers.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>

      <!-- 待处理申请 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              待处理申请
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'access_request', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/access_requests')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="request in recentRequests"
              :key="request.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex items-center space-x-3">
                <AlertTriangle class="h-5 w-5 text-orange-500" />
                <div>
                  <p class="font-medium">
                    {{ request.secret_item?.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    {{ request.applicant?.name }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(request.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentRequests.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>
    </div>

    <!-- 安全管理员面板 -->
    <div v-else-if="hasRole('sec_mgr')" class="grid gap-6 md:grid-cols-2">
      <!-- 待审批申请 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              待审批申请
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'access_request', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/access_requests')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="request in recentRequests"
              :key="request.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border hover:bg-accent cursor-pointer theme-transition"
              @click="router.push('/access_requests')"
            >
              <div class="flex items-center space-x-3">
                <Clock class="h-5 w-5 text-orange-500" />
                <div>
                  <p class="font-medium">
                    {{ request.secret_item?.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    申请人：{{ request.applicant?.name }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(request.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentRequests.length === 0" class="text-center py-8 text-muted-foreground">
              暂无待审批申请
            </div>
          </div>
        </div>
      </Card>

      <!-- 安全事件 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              安全事件
            </h3>
            <PermissionButton
              v-if="canViewAudit"
              variant="ghost"
              size="sm"
              :permission="{ resource: 'audit', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/audit')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="log in recentLogs"
              :key="log.id"
              class="flex items-center space-x-3 p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex-shrink-0">
                <div class="w-2 h-2 rounded-full" :class="getActionColor(log.action)" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate">
                  {{ log.user?.name }} {{ getActionDisplayName(log.action) }}了 {{ getResourceDisplayName(log.resource) }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ formatRelativeTime(log.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentLogs.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>
    </div>

    <!-- 开发人员面板 -->
    <div v-else-if="hasRole('dev')" class="grid gap-6 md:grid-cols-2">
      <!-- 我的访问权限 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              我的访问权限
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'secret', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/items')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="item in recentItems"
              :key="item.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border hover:bg-accent cursor-pointer theme-transition"
              @click="router.push(`/${item.type}/${item.id}`)"
            >
              <div class="flex items-center space-x-3">
                <div class="text-2xl">
                  {{ getFileIcon(item.type) }}
                </div>
                <div>
                  <p class="font-medium">
                    {{ item.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    {{ item.category }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(item.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentItems.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>

      <!-- 我的申请状态 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              我的申请
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'access_request', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/access_requests')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="request in recentRequests"
              :key="request.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex items-center space-x-3">
                <CheckCircle v-if="request.status === 'approved'" class="h-6 w-6 text-green-500" />
                <Clock v-else-if="request.status === 'pending'" class="h-6 w-6 text-orange-500" />
                <XCircle v-else class="h-6 w-6 text-red-500" />
                <div>
                  <p class="font-medium">
                    {{ request.secret_item?.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    状态：
                    <span v-if="request.status === 'approved'" class="text-green-600">已批准</span>
                    <span v-else-if="request.status === 'pending'" class="text-orange-600">待审批</span>
                    <span v-else class="text-red-600">已拒绝</span>
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(request.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentRequests.length === 0" class="text-center py-8 text-muted-foreground">
              暂无申请记录
            </div>
          </div>
        </div>
      </Card>
    </div>

    <!-- 审计员面板 -->
    <div v-else-if="hasRole('auditor')" class="grid gap-6 md:grid-cols-1">
      <!-- 审计日志 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              最近审计日志
            </h3>
            <PermissionButton
              v-if="canViewAudit"
              variant="ghost"
              size="sm"
              :permission="{ resource: 'audit', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/audit')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="log in recentLogs"
              :key="log.id"
              class="flex items-center space-x-3 p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex-shrink-0">
                <div class="w-2 h-2 rounded-full" :class="getActionColor(log.action)" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate">
                  {{ log.user?.name }} {{ getActionDisplayName(log.action) }}了 {{ getResourceDisplayName(log.resource) }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ formatRelativeTime(log.created_at) }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-xs text-muted-foreground">
                  {{ log.ip_address }}
                </p>
              </div>
            </div>
            <div v-if="recentLogs.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>
    </div>

    <!-- 默认面板 -->
    <div v-else class="grid gap-6 md:grid-cols-2">
      <!-- 最近创建的信息项 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              最近创建
            </h3>
            <PermissionButton
              variant="ghost"
              size="sm"
              :permission="{ resource: 'secret', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/items')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="item in recentItems"
              :key="item.id"
              class="flex items-center justify-between p-3 rounded-lg border border-border hover:bg-accent cursor-pointer theme-transition"
              @click="router.push(`/${item.type}/${item.id}`)"
            >
              <div class="flex items-center space-x-3">
                <div class="text-2xl">
                  {{ getFileIcon(item.type) }}
                </div>
                <div>
                  <p class="font-medium">
                    {{ item.name }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    {{ item.category }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-muted-foreground">
                  {{ formatRelativeTime(item.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentItems.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>

      <!-- 最近活动 -->
      <Card class="theme-transition">
        <div class="px-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">
              最近活动
            </h3>
            <PermissionButton
              v-if="canViewAudit"
              variant="ghost"
              size="sm"
              :permission="{ resource: 'audit', action: 'read' }"
              :fallback-visible="true"
              @click="router.push('/audit')"
            >
              查看全部
              <ArrowRight class="h-4 w-4" />
            </PermissionButton>
          </div>
          <div class="space-y-3">
            <div
              v-for="log in recentLogs"
              :key="log.id"
              class="flex items-center space-x-3 p-3 rounded-lg border border-border theme-transition"
            >
              <div class="flex-shrink-0">
                <div class="w-2 h-2 rounded-full" :class="getActionColor(log.action)" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate">
                  {{ log.user?.name }} {{ getActionDisplayName(log.action) }}了 {{ getResourceDisplayName(log.resource) }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ formatRelativeTime(log.created_at) }}
                </p>
              </div>
            </div>
            <div v-if="recentLogs.length === 0" class="text-center py-8 text-muted-foreground">
              暂无数据
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
