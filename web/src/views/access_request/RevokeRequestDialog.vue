<script setup lang="ts">
import type { AccessRequest } from '@/api/types'
import { Loader2 } from 'lucide-vue-next'
import { reactive, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { accessRequestAPI } from '@/api/access-request'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

interface Props {
  open: boolean
  request: AccessRequest
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const loading = ref(false)
const isOpen = ref(props.open)

const form = reactive({
  reason: '',
})

watch(
  () => props.open,
  (newValue) => {
    isOpen.value = newValue
    if (newValue) {
      // 重置表单
      form.reason = ''
    }
  },
)

// 监听内部 isOpen 变化
watch(isOpen, (newValue) => {
  emit('update:open', newValue)
})

async function handleSubmit() {
  if (!props.request?.id || !form.reason.trim()) {
    return
  }

  try {
    loading.value = true
    await accessRequestAPI.revokeRequest(props.request.id, {
      reason: form.reason.trim(),
    })

    toast.success('申请已作废')
    emit('success')
  }
  catch (error: any) {
    console.error('作废申请失败:', error)
    toast.error(error.response?.data?.error || '作废申请失败')
  }
  finally {
    loading.value = false
  }
}

function handleCancel() {
  isOpen.value = false
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>作废访问申请</DialogTitle>
        <DialogDescription>
          作废用户【<span class="font-bold">{{ request?.applicant?.name }}</span>】访问密钥项({{ request?.secret_item?.type }})【<span class="font-bold">{{ request?.secret_item?.name }}</span>】
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-2">
          <Label for="reason">作废理由</Label>
          <Textarea
            id="reason"
            v-model="form.reason"
            placeholder="请说明作废此申请的理由..."
            rows="4"
            required
          />
        </div>

        <div class="flex justify-end space-x-2">
          <Button type="button" variant="outline" @click="handleCancel">
            取消
          </Button>
          <Button
            type="submit"
            variant="destructive"
            :disabled="loading || !form.reason.trim()"
          >
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            作废申请
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
