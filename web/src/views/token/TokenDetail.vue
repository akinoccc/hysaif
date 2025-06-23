<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { Coins } from 'lucide-vue-next'
import { ref } from 'vue'
import { SecretDetail, SecretInput } from '@/components'
import { Label } from '@/components/ui/label'

const item = ref<SecretItem | null>(null)
function handleLoadSuccess(loadedItem: SecretItem) {
  item.value = loadedItem
}

function handleLoadError(error: any) {
  console.error('加载令牌失败:', error)
  item.value = null
}
</script>

<template>
  <SecretDetail
    title="令牌详情" description="查看和管理令牌信息" secret-type="token" secret-data-title="令牌信息" :secret-data-icon="Coins"
    error-text="令牌信息" @load-success="handleLoadSuccess" @load-error="handleLoadError"
  >
    <template #secret-data>
      <!-- 令牌值 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">令牌值</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data?.token || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
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
