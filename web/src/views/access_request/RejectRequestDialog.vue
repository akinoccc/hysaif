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
  request?: AccessRequest | null
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isOpen = ref(props.open)
const loading = ref(false)

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

watch(isOpen, (newValue) => {
  emit('update:open', newValue)
})

async function handleSubmit() {
  if (!props.request || !form.reason.trim()) {
    return
  }

  loading.value = true
  try {
    await accessRequestAPI.rejectRequest(props.request.id, {
      reason: form.reason.trim(),
    })

    toast.success('申请已拒绝')
    isOpen.value = false
    emit('success')
  }
  catch (error: any) {
    console.error('拒绝申请失败:', error)
    toast.error(error.response?.data?.error || '拒绝申请失败')
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
        <DialogTitle>拒绝访问申请</DialogTitle>
        <DialogDescription>
          拒绝用户【<span class="font-bold">{{ request?.applicant?.name }}</span>】访问密钥项【<span class="font-bold">{{ request?.secret_item?.name }}</span>】
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-2">
          <Label for="reason">拒绝理由</Label>
          <Textarea
            id="reason"
            v-model="form.reason"
            placeholder="请说明拒绝此申请的理由..."
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
            <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
            拒绝申请
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
