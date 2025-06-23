<script setup lang="ts">
import type { ButtonPermission } from '@/stores/permission'
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { usePermission } from '@/composables/usePermission'

interface Props {
  permission: ButtonPermission
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false,
  loading: false,
})

const { hasButtonPermission } = usePermission()

// 检查是否有权限
const hasPermission = computed(() => {
  return hasButtonPermission(props.permission)
})

// 计算最终的禁用状态
const isDisabled = computed(() => {
  return props.disabled || props.loading || !hasPermission.value
})
</script>

<template>
  <Button
    v-if="hasPermission"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    v-bind="$attrs"
  >
    <slot />
  </Button>
</template>
