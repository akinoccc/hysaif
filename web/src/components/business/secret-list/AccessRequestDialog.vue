<script setup lang="ts">
import type { SecretItem } from '@/api/types'
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
  item?: SecretItem
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

// 监听 props.open 变化
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

// 监听 isOpen 变化，同步到父组件
watch(isOpen, (newValue) => {
  emit('update:open', newValue)
})

async function handleSubmit() {
  if (!props.item || !form.reason.trim()) {
    return
  }

  loading.value = true
  try {
    await accessRequestAPI.createRequest({
      secret_item_id: props.item.id,
      reason: form.reason.trim(),
    })

    toast.success('申请已提交，请等待管理员审批')
    isOpen.value = false
    emit('success')
  }
  catch (error: any) {
    console.error('提交申请失败:', error)
    toast.error(error.response?.data?.error || '提交申请失败')
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
        <DialogTitle>申请访问密钥</DialogTitle>
        <DialogDescription>
          您正在申请访问密钥项：{{ item?.name }}
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-2">
          <Label for="reason">申请理由</Label>
          <Textarea
            id="reason"
            v-model="form.reason"
            placeholder="请说明申请访问此密钥的理由..."
            rows="4"
            required
          />
        </div>

        <div class="flex justify-end space-x-2">
          <Button type="button" variant="outline" @click="handleCancel">
            取消
          </Button>
          <Button type="submit" :disabled="loading || !form.reason.trim()">
            <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
            提交申请
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
