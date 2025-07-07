<script setup lang="ts">
import type { ButtonPermission } from '@/stores/permission'
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { usePermissionStore } from '@/stores/permission'

interface Props {
  permission: ButtonPermission
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
  loading?: boolean
  fallbackVisible?: boolean // 权限检查失败时是否显示，默认为 false
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false,
  loading: false,
  fallbackVisible: false,
})

const permissionStore = usePermissionStore()

// 使用同步权限检查
const hasPermission = computed(() => {
  return permissionStore.hasPermission(
    props.permission.resource,
    props.permission.action,
    props.fallbackVisible,
  )
})

// 计算最终的禁用状态
const isDisabled = computed(() => {
  return props.disabled || props.loading || !hasPermission.value
})

// 是否显示加载状态
const showLoading = computed(() => {
  return props.loading
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
    <div v-if="showLoading" class="flex items-center space-x-2">
      <div class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
      <span v-if="$slots.default"><slot /></span>
    </div>
    <slot v-else />
  </Button>
</template>
