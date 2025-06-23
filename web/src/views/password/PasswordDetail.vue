<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { User } from 'lucide-vue-next'
import { ref } from 'vue'
import { SecretDetail, SecretInput } from '@/components'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const item = ref<SecretItem | null>(null)

function handleLoadSuccess(loadedItem: SecretItem) {
  item.value = loadedItem
}

function handleLoadError(error: any) {
  console.error('加载账号密码失败:', error)
  item.value = null
}
</script>

<template>
  <SecretDetail
    title="账号密码详情" description="查看和管理账号密码信息" secret-type="password" secret-data-title="账号信息"
    :secret-data-icon="User" error-text="账号密码信息" @load-success="handleLoadSuccess" @load-error="handleLoadError"
  >
    <template #secret-data>
      <!-- 用户名 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">用户名</Label>
        <div class="mt-1">
          <Input readonly copyable :model-value="item?.data?.username || '未填写'" class=" cursor-not-allowed" />
        </div>
      </div>

      <!-- 密码 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">密码</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data?.password || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- 网址 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">网址</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.address || '未填写' }}
        </p>
      </div>

      <!-- 备注 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">备注</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.notes || '未填写' }}
        </p>
      </div>
    </template>
  </SecretDetail>
</template>
