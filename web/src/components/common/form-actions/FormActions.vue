<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

interface Props {
  loading: boolean
  isEdit: boolean
}

interface Emits {
  (e: 'cancel'): void
}

defineProps<Props>()
defineEmits<Emits>()
</script>

<template>
  <div class="flex justify-end space-x-4 py-4">
    <Button type="button" variant="outline" @click="$emit('cancel')">
      取消
    </Button>
    <Button type="submit" :disabled="loading">
      <Loader2 v-if="loading" class=" h-4 w-4 animate-spin" />
      {{ isEdit ? '更新' : '创建' }}
    </Button>
  </div>
</template>
