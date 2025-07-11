<script setup lang="ts">
import type { VersionComparisonResponse } from '@/api/types'
import { GitCompare, Loader2, X } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import { secretItemAPI } from '@/api/secret'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'

const props = defineProps<{
  open: boolean
  itemId: string | number
  version1: number
  version2: number
  comparing: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'close'): void
}>()

const loading = ref(false)
const comparisonResult = ref<VersionComparisonResponse | null>(null)

async function loadComparison() {
  if (!props.itemId || !props.version1 || !props.version2) {
    return
  }

  try {
    loading.value = true
    const response = await secretItemAPI.compareVersions(props.itemId, {
      version1: props.version1,
      version2: props.version2,
    })
    comparisonResult.value = response
  }
  catch (error) {
    console.error('比较版本失败:', error)
    toast.error('比较版本失败')
  }
  finally {
    loading.value = false
  }
}

function handleClose() {
  emit('update:open', false)
  emit('close')
}

function getChangeType(oldValue: any, newValue: any): 'added' | 'removed' | 'modified' {
  if (oldValue === undefined || oldValue === null || oldValue === '') {
    return 'added'
  }
  if (newValue === undefined || newValue === null || newValue === '') {
    return 'removed'
  }
  return 'modified'
}

function formatValue(value: any): string {
  if (value === null || value === undefined) {
    return '(空)'
  }
  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2)
  }
  return String(value)
}

onMounted(() => {
  if (props.open) {
    loadComparison()
  }
})
</script>

<template>
  <Dialog
    :open="open"
    @update:open="handleClose"
  >
    <DialogContent class="max-w-4xl max-h-[80vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2">
          <GitCompare class="h-5 w-5" />
          版本对比
        </DialogTitle>
      </DialogHeader>

      <div class="space-y-6">
        <!-- 版本信息 -->
        <div class="flex items-center justify-between p-4 bg-muted rounded-lg">
          <div class="flex items-center justify-around gap-4 w-full">
            <div class="text-center">
              <Badge variant="secondary">
                旧版本
              </Badge>
              <div class="text-lg font-semibold">
                v{{ version1 }}
              </div>
            </div>
            <div class="text-muted-foreground">
              vs
            </div>
            <div class="text-center">
              <Badge variant="secondary">
                新版本
              </Badge>
              <div class="text-lg font-semibold">
                v{{ version2 }}
              </div>
            </div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center py-8">
          <Loader2 class="h-6 w-6 animate-spin" />
        </div>

        <!-- 比较结果 -->
        <div v-else-if="comparisonResult && comparisonResult.changes" class="space-y-6">
          <div
            v-for="(change, field) in comparisonResult.changes"
            :key="field"
            class="space-y-4"
          >
            <div class="flex items-center gap-2">
              <Label class="text-sm font-medium">
                {{ field }}
              </Label>
              <Badge
                :variant="getChangeType(change.old, change.new) === 'removed' ? 'destructive' : 'secondary'"
                class="text-xs"
              >
                {{
                  getChangeType(change.old, change.new) === 'added' ? '新增'
                  : getChangeType(change.old, change.new) === 'removed' ? '删除' : '修改'
                }}
              </Badge>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- 旧值 -->
              <div class="space-y-2">
                <Label class="text-xs text-muted-foreground">
                  旧值 (v{{ version1 }})
                </Label>
                <div class="p-3 bg-red-50 border border-red-200 rounded-md">
                  <pre class="text-sm text-red-800 whitespace-pre-wrap">{{ formatValue(change.old) }}</pre>
                </div>
              </div>

              <!-- 新值 -->
              <div class="space-y-2">
                <Label class="text-xs text-muted-foreground">
                  新值 (v{{ version2 }})
                </Label>
                <div class="p-3 bg-green-50 border border-green-200 rounded-md">
                  <pre class="text-sm text-green-800 whitespace-pre-wrap">{{ formatValue(change.new) }}</pre>
                </div>
              </div>
            </div>

            <Separator />
          </div>
        </div>

        <!-- 无变化 -->
        <div v-else-if="comparisonResult && !comparisonResult.changes" class="text-center py-8">
          <GitCompare class="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
          <h3 class="text-lg font-medium mb-2">
            无变化
          </h3>
          <p class="text-muted-foreground">
            两个版本之间没有发现任何差异
          </p>
        </div>

        <!-- 错误状态 -->
        <div v-else class="text-center py-8">
          <GitCompare class="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
          <h3 class="text-lg font-medium mb-2">
            加载失败
          </h3>
          <p class="text-muted-foreground">
            无法加载版本比较结果
          </p>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
