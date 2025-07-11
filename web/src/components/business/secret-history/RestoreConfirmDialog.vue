<script setup lang="ts">
import type { SecretItemHistory } from '@/api/types'
import { AlertTriangle } from 'lucide-vue-next'
import { ref } from 'vue'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

defineProps<{
  open: boolean
  history: SecretItemHistory | null
  restoring: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm', reason?: string): void
  (e: 'close'): void
}>()

const reason = ref('')

function handleConfirm() {
  emit('confirm', reason.value)
  reason.value = ''
}

function handleClose() {
  emit('update:open', false)
  emit('close')
  reason.value = ''
}
</script>

<template>
  <AlertDialog
    :open="open"
    @update:open="handleClose"
  >
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle class="flex items-center gap-2">
          <AlertTriangle class="h-5 w-5 text-amber-500" />
          确认恢复版本
        </AlertDialogTitle>
        <AlertDialogDescription>
          此操作将会恢复到历史版本，当前数据将会被覆盖。
        </AlertDialogDescription>
      </AlertDialogHeader>

      <div class="space-y-4">
        <div class="p-4 bg-muted rounded-lg">
          <h4 class="font-medium mb-2">
            恢复信息
          </h4>
          <div class="space-y-2 text-sm">
            <div>
              <span class="text-muted-foreground">版本：</span>
              <span class="font-medium">v{{ history?.version }}</span>
            </div>
            <div>
              <span class="text-muted-foreground">修改说明：</span>
              <span>{{ history?.change_reason || '无' }}</span>
            </div>
            <div>
              <span class="text-muted-foreground">创建时间：</span>
              <span>{{ history?.created_at ? new Date(history.created_at).toLocaleString() : '未知' }}</span>
            </div>
            <div v-if="history?.created_by">
              <span class="text-muted-foreground">创建者：</span>
              <span>{{ history.created_by.name }}</span>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <Label for="reason">恢复原因（可选）</Label>
          <Textarea
            id="reason"
            v-model="reason"
            placeholder="请输入恢复原因..."
            rows="3"
          />
        </div>
      </div>

      <AlertDialogFooter>
        <AlertDialogCancel @click="handleClose">
          取消
        </AlertDialogCancel>
        <AlertDialogAction as-child>
          <Button
            variant="destructive"
            :disabled="restoring"
            @click="handleConfirm"
          >
            {{ restoring ? '恢复中...' : '确认恢复' }}
          </Button>
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
