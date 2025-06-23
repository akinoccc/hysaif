<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { Key } from 'lucide-vue-next'
import { ref } from 'vue'
import { SecretDetail, SecretInput } from '@/components'
import { Label } from '@/components/ui/label'

const item = ref<SecretItem | null>(null)

function handleLoadSuccess(loadedItem: SecretItem) {
  item.value = loadedItem
}

function handleLoadError(error: any) {
  console.error('加载API密钥失败:', error)
  item.value = null
}
</script>

<template>
  <SecretDetail
    title="API 密钥详情" description="查看和管理 API 密钥信息" secret-type="api_key" secret-data-title="密钥信息"
    :secret-data-icon="Key" error-text="API 密钥信息" @load-success="handleLoadSuccess" @load-error="handleLoadError"
  >
    <template #secret-data>
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">API 密钥</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data.api_key || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">API Secret</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data.api_secret || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">API 端点</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.endpoint || '未填写' }}
        </p>
      </div>

      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">备注</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.notes || '未填写' }}
        </p>
      </div>
    </template>
  </SecretDetail>
</template>
