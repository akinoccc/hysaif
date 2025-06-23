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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
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
  validDuration: '24', // 默认1天
  note: '',
})

watch(
  () => props.open,
  (newValue) => {
    isOpen.value = newValue
    if (newValue) {
      // 重置表单
      form.validDuration = '24'
      form.note = ''
    }
  },
)

// 监听 isOpen 变化，同步到父组件
watch(isOpen, (newValue) => {
  emit('update:open', newValue)
})

async function handleSubmit() {
  if (!props.request || !form.validDuration) {
    return
  }

  loading.value = true
  try {
    await accessRequestAPI.approveRequest(props.request.id, {
      valid_duration: Number.parseInt(form.validDuration),
      note: form.note.trim() || undefined,
    })

    toast.success('申请已批准')
    isOpen.value = false
    emit('success')
  }
  catch (error: any) {
    console.error('批准申请失败:', error)
    toast.error(error.response?.data?.error || '批准申请失败')
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
        <DialogTitle>批准访问申请</DialogTitle>
        <DialogDescription>
          批准用户【<span class="font-bold">{{ request?.applicant?.name }}</span>】访问密钥项({{ request?.secret_item?.type }})【<span class="font-bold">{{ request?.secret_item?.name }}</span>】
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-2">
          <Label for="duration">有效时长（小时）</Label>
          <Select v-model="form.validDuration">
            <SelectTrigger>
              <SelectValue placeholder="选择有效时长" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">
                1小时
              </SelectItem>
              <SelectItem value="4">
                4小时
              </SelectItem>
              <SelectItem value="8">
                8小时
              </SelectItem>
              <SelectItem value="24">
                1天
              </SelectItem>
              <SelectItem value="72">
                3天
              </SelectItem>
              <SelectItem value="168">
                1周
              </SelectItem>
              <SelectItem value="720">
                1个月
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <Label for="note">审批备注（可选）</Label>
          <Textarea
            id="note"
            v-model="form.note"
            placeholder="添加审批备注..."
            rows="3"
          />
        </div>

        <div class="flex justify-end space-x-2">
          <Button type="button" variant="outline" @click="handleCancel">
            取消
          </Button>
          <Button type="submit" :disabled="loading || !form.validDuration">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            批准申请
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
