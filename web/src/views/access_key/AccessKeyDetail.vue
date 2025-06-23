<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { Cloud } from 'lucide-vue-next'
import { ref } from 'vue'
import { SecretDetail, SecretInput } from '@/components'
import { Label } from '@/components/ui/label'

const item = ref<SecretItem | null>(null)

function handleLoadSuccess(loadedItem: SecretItem) {
  item.value = loadedItem
}

function handleLoadError(error: any) {
  console.error('加载访问密钥失败:', error)
  item.value = null
}
</script>

<template>
  <SecretDetail
    title="访问密钥详情" description="查看和管理访问密钥信息" secret-type="access_key" secret-data-title="密钥信息"
    :secret-data-icon="Cloud" error-text="访问密钥信息" @load-success="handleLoadSuccess" @load-error="handleLoadError"
  >
    <template #secret-data>
      <!-- Access Key -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">Access Key</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data?.access_key || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- Secret Key -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">Secret Key</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data?.secret_key || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- 区域 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">区域</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.region || '未填写' }}
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
