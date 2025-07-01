<script setup lang="ts">
import type { ButtonPermission } from '@/stores/permission'
import { computed, onMounted, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { usePermission } from '@/composables/usePermission'

interface Props {
  permission: ButtonPermission
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
  loading?: boolean
  fallbackVisible?: boolean // 权限检查失败时是否显示，默认为 false
  useAsync?: boolean // 是否使用异步检查，默认为 true
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false,
  loading: false,
  fallbackVisible: false,
  useAsync: true,
})

const { hasButtonPermission, checkPermission } = usePermission()

// 权限状态
const hasPermission = ref(props.fallbackVisible)
const isCheckingPermission = ref(false)

/**
 * 检查权限
 */
async function checkButtonPermission() {
  if (!props.useAsync) {
    // 同步检查
    const result = hasButtonPermission(props.permission, false)
    if (result instanceof Promise) {
      hasPermission.value = props.fallbackVisible
    }
    else {
      hasPermission.value = result
    }
    return
  }

  // 异步检查
  isCheckingPermission.value = true
  try {
    const result = await checkPermission(
      props.permission.resource,
      props.permission.action,
      props.fallbackVisible,
    )
    hasPermission.value = result
  }
  catch (error) {
    console.error('权限检查失败:', error)
    hasPermission.value = props.fallbackVisible
  }
  finally {
    isCheckingPermission.value = false
  }
}

// 计算最终的禁用状态
const isDisabled = computed(() => {
  return props.disabled || props.loading || isCheckingPermission.value || !hasPermission.value
})

// 是否显示加载状态
const showLoading = computed(() => {
  return props.loading || isCheckingPermission.value
})

// 监听权限配置变化
watch(
  () => props.permission,
  () => {
    checkButtonPermission()
  },
  { deep: true },
)

// 监听异步模式变化
watch(
  () => props.useAsync,
  () => {
    checkButtonPermission()
  },
)

// 组件挂载时检查权限
onMounted(() => {
  checkButtonPermission()
})
</script>

<template>
  <Button
    v-if="hasPermission || isCheckingPermission"
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
