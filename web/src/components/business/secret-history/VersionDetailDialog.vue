<script setup lang="ts">
import type { SecretItemHistory } from '@/api/types'
import { formatDate } from 'date-fns'
import { Clock, Eye } from 'lucide-vue-next'
import { computed } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { getCategoryByKey } from '@/constants'

const props = defineProps<{
  open: boolean
  history: SecretItemHistory | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'close'): void
}>()

const categoryInfo = computed(() => {
  if (!props.history) {
    return null
  }
  return getCategoryByKey(props.history.category)
})

function handleClose() {
  emit('update:open', false)
  emit('close')
}
</script>

<template>
  <Dialog
    :open="open"
    @update:open="handleClose"
  >
    <DialogContent class="max-w-2xl max-h-[80vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2">
          <Eye class="h-5 w-5" />
          版本详情 - v{{ history?.version }}
        </DialogTitle>
      </DialogHeader>

      <div v-if="history" class="space-y-6">
        <!-- 基本信息 -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <Clock class="h-4 w-4 text-muted-foreground" />
            <h3 class="font-semibold">
              基本信息
            </h3>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">名称</Label>
              <p class="text-sm font-medium">
                {{ history.name }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">版本</Label>
              <p class="text-sm font-medium">
                v{{ history.version }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">类型</Label>
              <p class="text-sm">
                {{ history.type }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">分类</Label>
              <p class="text-sm">
                {{ categoryInfo?.label || history.category }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">环境</Label>
              <p class="text-sm">
                {{ history.environment || '无' }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">修改类型</Label>
              <p class="text-sm">
                {{ history.change_type || '更新' }}
              </p>
            </div>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">描述</Label>
            <p class="text-sm">
              {{ history.description || '无' }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">修改说明</Label>
            <p class="text-sm">
              {{ history.change_reason || '无修改说明' }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">标签</Label>
            <div class="flex flex-wrap gap-2">
              <template v-if="history.tags && history.tags.length > 0">
                <span
                  v-for="tag in history.tags"
                  :key="tag"
                  class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-muted text-muted-foreground"
                >
                  {{ tag }}
                </span>
              </template>
              <template v-else>
                <span class="text-sm text-muted-foreground">无</span>
              </template>
            </div>
          </div>
        </div>

        <Separator />

        <!-- 系统信息 -->
        <div class="space-y-4">
          <h3 class="font-semibold">
            系统信息
          </h3>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">创建时间</Label>
              <p class="text-sm">
                {{ formatDate(new Date(history.created_at), 'yyyy-MM-dd HH:mm:ss') }}
              </p>
            </div>

            <div class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">创建者</Label>
              <p class="text-sm">
                {{ history.created_by?.name || '未知' }}
              </p>
            </div>

            <div v-if="history.expires_at" class="space-y-2">
              <Label class="text-sm font-medium text-muted-foreground">过期时间</Label>
              <p class="text-sm">
                {{ formatDate(new Date(history.expires_at), 'yyyy-MM-dd HH:mm:ss') }}
              </p>
            </div>
          </div>
        </div>

        <Separator />

        <!-- 数据信息 -->
        <div class="space-y-4">
          <h3 class="font-semibold">
            数据信息
          </h3>

          <div class="p-4 bg-muted/50 rounded-lg">
            <p class="text-sm text-muted-foreground mb-2">
              此版本的数据信息已被安全存储，出于安全考虑，在此处不显示敏感数据内容。
            </p>
            <p class="text-xs text-muted-foreground">
              如需查看完整数据，请恢复到此版本。
            </p>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
