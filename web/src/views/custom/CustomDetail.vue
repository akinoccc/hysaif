<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { Braces } from 'lucide-vue-next'
import { ref } from 'vue'
import { SecretDetail, SecretInput } from '@/components'

import { Label } from '@/components/ui/label'

const item = ref<SecretItem | null>(null)

function handleLoadSuccess(loadedItem: SecretItem) {
  item.value = loadedItem
}

function handleLoadError(error: any) {
  console.error('加载自定义密钥失败:', error)
  item.value = null
}
</script>

<template>
  <SecretDetail
    title="自定义密钥详情" description="查看和管理自定义密钥信息" secret-type="custom" secret-data-title="密钥信息"
    :secret-data-icon="Braces" error-text="自定义密钥信息" @load-success="handleLoadSuccess" @load-error="handleLoadError"
  >
    <template #secret-data>
      <!-- 自定义字段列表 -->
      <div v-if="item?.data?.custom_data?.length > 0" class="space-y-4">
        <div v-for="(entry, index) in item?.data?.custom_data" :key="index" class="space-y-2 p-4 border rounded-md bg-card">
          <!-- 字段名称 -->
          <Label class="text-sm font-medium text-muted-foreground">{{ entry.key }}</Label>
          <div class="mt-1">
            <SecretInput readonly toggleable copyable :model-value="entry.value" class=" cursor-not-allowed" />
          </div>
        </div>
      </div>

      <!-- 无数据提示 -->
      <div v-else class="p-4 border rounded-md bg-muted">
        <p class="text-sm text-muted-foreground text-center">
          无自定义字段数据
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
