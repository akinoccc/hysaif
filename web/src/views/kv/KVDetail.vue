<script setup lang="ts">
import { Braces, Code2, Copy, Eye } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'
import { SecretDetail, SecretInput } from '@/components'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { copyToClipboard } from '@/lib/utils'

const isRawMode = ref(false)

const infoRef = ref<InstanceType<typeof SecretDetail>>()

// 生成 JSON 格式的数据
const jsonData = computed(() => {
  if (!infoRef.value?.item?.data?.custom_data) {
    return {}
  }

  const result: Record<string, string> = {}
  infoRef.value?.item?.data?.custom_data.forEach((entry: { key: string, value: string }) => {
    result[entry.key] = entry.value
  })

  return result
})

// 格式化后的 JSON 字符串
const formattedJson = computed(() => {
  return JSON.stringify(jsonData.value, null, 2)
})

// 切换视图模式
function toggleViewMode() {
  isRawMode.value = !isRawMode.value
}

// 复制 JSON 数据
async function copyJsonData() {
  const success = await copyToClipboard(formattedJson.value)
  if (success) {
    toast.success('JSON 数据已复制到剪贴板')
  }
  else {
    toast.error('复制失败')
  }
}
</script>

<template>
  <SecretDetail
    title="KV 键值对详情" description="查看和管理 KV 键值对信息" secret-type="custom" secret-data-title="KV 键值对信息"
    :secret-data-icon="Braces" error-text="KV 键值对信息"
  >
    <template #secret-actions>
      <!-- 切换视图模式按钮 -->
      <Button
        variant="ghost"
        size="sm"
        :title="isRawMode ? '切换到正常视图' : '切换到 Raw JSON 视图'"
        @click="toggleViewMode"
      >
        <Eye v-if="isRawMode" class="h-4 w-4" />
        <Code2 v-else class="h-4 w-4" />
        {{ isRawMode ? '正常视图' : 'Raw JSON' }}
      </Button>
    </template>

    <template #secret-data="{ item }">
      <!-- 正常视图 -->
      <div class="space-y-4">
        <!-- KV 键值对列表 -->
        <div v-if="item?.data?.custom_data?.length > 0" class="space-y-4">
          <template v-if="!isRawMode">
            <div v-for="(entry, index) in item?.data?.custom_data" :key="index" class="space-y-2 p-4 border rounded-md bg-card">
              <!-- 字段名称 -->
              <Label class="text-sm font-medium text-muted-foreground">{{ entry.key }}</Label>
              <div class="mt-1">
                <SecretInput readonly toggleable copyable :model-value="entry.value" class=" cursor-not-allowed" />
              </div>
            </div>
          </template>

          <template v-else>
            <div class="relative">
              <pre class="bg-muted/50 text-foreground border border-border rounded-lg p-4 font-mono text-sm overflow-x-auto whitespace-pre-wrap break-all"><code>{{ formattedJson }}</code></pre>
              <Button
                variant="ghost"
                size="sm"
                class="absolute top-2 right-2"
                title="复制 JSON 数据"
                @click="copyJsonData"
              >
                <Copy class="h-4 w-4" />
              </Button>
            </div>
          </template>
        </div>

        <!-- 无数据提示 -->
        <div v-else class="p-4 border rounded-md bg-muted">
          <p class="text-sm text-muted-foreground text-center">
            无 KV 键值对数据
          </p>
        </div>

        <!-- 备注 -->
        <div class="space-y-2">
          <Label class="text-sm font-medium text-muted-foreground">备注</Label>
          <p class="text-sm mt-1 p-3 bg-muted rounded-md">
            {{ item?.data?.notes || '未填写' }}
          </p>
        </div>
      </div>
    </template>
  </SecretDetail>
</template>
