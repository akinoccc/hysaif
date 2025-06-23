<script setup lang="ts">
import { Copy, Eye, EyeOff } from 'lucide-vue-next'
import { computed, type HTMLAttributes, ref } from 'vue'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
  readonly?: boolean
  toggleable?: boolean
  copyable?: boolean
  placeholder?: string
}>()

const modelValue = defineModel<string>()
const showSecret = ref(false)

// 默认值处理
const readonly = computed(() => props.readonly !== undefined ? props.readonly : true)
const toggleable = computed(() => props.toggleable !== undefined ? props.toggleable : true)
const copyable = computed(() => props.copyable !== undefined ? props.copyable : true)

// 复制到剪贴板
async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(modelValue.value || '')
    toast.success('已复制到剪贴板')
  }
  catch (error) {
    console.error('复制失败:', error)
    toast.error('复制失败')
  }
}
</script>

<template>
  <div class="relative">
    <Input
      v-bind="$attrs"
      v-model="modelValue"
      :type="showSecret ? 'text' : 'password'"
      :readonly="readonly"
      :placeholder="placeholder"
      :class="cn('pr-20', props.class, readonly ? 'cursor-not-allowed' : '')"
    />
    <div class="absolute right-0 top-0 h-full flex">
      <Button
        v-if="toggleable" type="button" variant="ghost" size="sm" class="h-full px-2"
        @click="showSecret = !showSecret"
      >
        <Eye v-if="!showSecret" class="h-4 w-4" />
        <EyeOff v-else class="h-4 w-4" />
      </Button>
      <Button v-if="copyable" type="button" variant="ghost" size="sm" class="h-full px-2" @click="copyToClipboard">
        <Copy class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>
