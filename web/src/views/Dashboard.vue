<script setup lang="ts">
import {
  ArrowRight,
  Clock,
  Database,
  Eye,
  Shield,
} from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { type SecretItem, secretItemAPI } from '@/api'
import { PermissionButton } from '@/components'
import { Card } from '@/components/ui/card'
import { formatRelativeTime, getFileIcon } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const stats = ref({
  totalItems: 0,
  expiringItems: 0,
  todayAccess: 0,
  securityScore: 95,
})

const recentItems = ref<SecretItem[]>([])

async function loadDashboardData() {
  try {
    // 加载信息项统计
    const itemsResponse = await secretItemAPI.getItems({ page: 1, page_size: 5 })
    stats.value.totalItems = itemsResponse.pagination.total
    recentItems.value = itemsResponse.data || []

    // 计算即将过期的项目数量
    const now = new Date()
    const thirtyDaysLater = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)
    stats.value.expiringItems = recentItems.value.filter((item: SecretItem) => {
      if (!item.expires_at)
        return false
      const expiresAt = new Date(item.expires_at)
      return expiresAt <= thirtyDaysLater
    }).length
  }
  catch (error) {
    console.error('Failed to load dashboard data:', error)
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
          仪表板
        </h1>
        <p class="text-muted-foreground mt-2">
          欢迎回来，{{ authStore.user?.name }}！
        </p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
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

    <div class="grid gap-6 md:grid-cols-2">
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
              @click="router.push('/items')"
            >
              查看全部
              <ArrowRight class="ml-2 h-4 w-4" />
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
    </div>
  </div>
</template>
