<script setup lang="ts">
import type { AuditLog } from '@/api'
import { Monitor, Search, User } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { AUDIT_LOG_ACTION_LIST } from '@/constants'
import { formatDate, formatRelativeTime } from '@/lib/utils'
import { getActionColor, getActionDisplayName, getActionIcon, getResourceDisplayName, getUserAgentInfo } from './helper'

const selectedLog = defineModel<AuditLog>({ required: true })
const router = useRouter()

function navigateToResourceDetail() {
  if (selectedLog.value.action === AUDIT_LOG_ACTION_LIST.Login || selectedLog.value.action === AUDIT_LOG_ACTION_LIST.Logout) {
    return
  }
  router.push(`/${selectedLog.value.resource}/${selectedLog.value.resource_id}`)
}
</script>

<template>
  <!-- 顶部标题栏 -->
  <div class="sticky top-0 z-10 bg-background border-b px-6 py-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div
          v-if="selectedLog"
          class="w-10 h-10 rounded-full flex items-center justify-center"
          :class="getActionColor(selectedLog.action)"
        >
          <component :is="getActionIcon(selectedLog.action)" class="h-5 w-5" />
        </div>
        <div>
          <h2 class="text-lg font-semibold">
            {{ selectedLog ? getActionDisplayName(selectedLog.action) : '审计日志详情' }}
          </h2>
          <p class="text-sm text-muted-foreground">
            完整的审计日志详细信息
          </p>
        </div>
      </div>
    </div>
  </div>

  <div v-if="selectedLog" class="px-6 py-4 space-y-6">
    <!-- 基本信息卡片 -->
    <div class="bg-card rounded-xl border shadow-sm p-5">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-semibold flex items-center gap-x-2">
          <User /> 基本信息
        </h3>
        <div class="flex items-center space-x-2">
          <span
            class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium"
            :class="getActionColor(selectedLog.action)"
          >
            {{ getActionDisplayName(selectedLog.action) }}
          </span>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-x-8 gap-y-4">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">用户</label>
          <p class="text-sm font-medium">
            {{ selectedLog.user?.name || '未知用户' }}
          </p>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">用户角色</label>
          <p class="text-sm font-medium">
            {{ selectedLog.user?.role || '未知角色' }}
          </p>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">时间</label>
          <p class="text-sm font-medium">
            {{ formatDate(selectedLog.created_at) }}
          </p>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">相对时间</label>
          <p class="text-sm font-medium">
            {{ formatRelativeTime(selectedLog.created_at) }}
          </p>
        </div>
      </div>
    </div>

    <!-- 资源信息卡片 -->
    <div class="bg-card rounded-xl border shadow-sm p-5">
      <h3 class="font-semibold flex items-center gap-x-2 mb-4">
        <Search /> 资源信息
      </h3>
      <div class="grid grid-cols-2 gap-x-8 gap-y-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-muted-foreground">资源类型</label>
          <p class="text-sm font-medium">
            {{ selectedLog.resource }}
          </p>
        </div>
        <div v-if="selectedLog.resource" class="flex flex-col gap-1">
          <label class="text-sm font-medium text-muted-foreground">资源名称</label>
          <p class="text-sm font-medium">
            {{ getResourceDisplayName(selectedLog.resource) }}
          </p>
        </div>
        <div v-if="selectedLog.resource_id" class="col-span-2 flex flex-col gap-1">
          <label class="text-sm font-medium text-muted-foreground">资源ID</label>
          <p class="text-sm text-blue-600 cursor-pointer bg-muted/40 px-2 py-1 rounded-md" @click="navigateToResourceDetail">
            {{ selectedLog.resource_id }}
          </p>
        </div>
      </div>
    </div>

    <!-- 操作详情卡片 -->
    <div v-if="selectedLog.details" class="bg-card rounded-xl border shadow-sm p-5">
      <h3 class="font-semibold flex items-center gap-x-2 mb-4">
        <FileText /> 操作详情
      </h3>
      <div class="bg-muted/30 rounded-lg p-4 border">
        <p class="text-sm whitespace-pre-wrap break-all">
          {{ JSON.parse(selectedLog.details) }}
        </p>
      </div>
    </div>

    <!-- 系统信息卡片 -->
    <div class="bg-card rounded-xl border shadow-sm p-5">
      <h3 class="font-semibold flex items-center gap-x-2 mb-4">
        <Monitor /> 系统信息
      </h3>
      <div class="grid grid-cols-2 gap-x-8 gap-y-4 mb-4">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">IP地址</label>
          <p class="text-sm  bg-muted/40 px-2 py-1 rounded-md">
            {{ selectedLog.ip_address || '未知IP' }}
          </p>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-muted-foreground">用户代理</label>
          <p class="text-sm">
            {{ selectedLog.user_agent ? getUserAgentInfo(selectedLog.user_agent) : '未知设备' }}
          </p>
        </div>
      </div>
      <div v-if="selectedLog.user_agent" class="flex flex-col gap-2">
        <label class="text-sm font-medium text-muted-foreground">完整用户代理</label>
        <div class="bg-muted/30 rounded-lg p-3 border">
          <p class="text-xs  break-all">
            {{ selectedLog.user_agent }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
