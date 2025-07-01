<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { usePermission } from '@/composables/usePermission'

interface Props {
  resource: string
  action: string
  mode?: 'hide' | 'disable' // 默认为 hide
  fallback?: boolean // 是否显示fallback内容
  fallbackValue?: boolean // 权限检查失败时的回退值，默认为 false
  useAsync?: boolean // 是否使用异步检查，默认为 true
  showLoading?: boolean // 是否显示加载状态，默认为 false
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'hide',
  fallback: false,
  fallbackValue: false,
  useAsync: true,
  showLoading: false,
})

const { hasPermission, checkPermission } = usePermission()

// 权限状态
const hasAccess = ref(props.fallbackValue)
const isCheckingPermission = ref(false)

/**
 * 检查权限
 */
async function checkAccessPermission() {
  if (!props.useAsync) {
    // 同步检查
    const result = hasPermission(props.resource, props.action, props.fallbackValue)
    hasAccess.value = result
    return
  }

  // 异步检查
  isCheckingPermission.value = true
  try {
    const result = await checkPermission(
      props.resource,
      props.action,
      props.fallbackValue,
    )
    hasAccess.value = result
  }
  catch (error) {
    console.error('权限检查失败:', error)
    hasAccess.value = props.fallbackValue
  }
  finally {
    isCheckingPermission.value = false
  }
}

// 是否显示内容
const shouldShow = computed(() => {
  if (props.mode === 'hide') {
    return hasAccess.value || (isCheckingPermission.value && props.showLoading)
  }
  return true // disable模式下总是显示
})

// 是否禁用内容
const shouldDisable = computed(() => {
  if (props.mode === 'disable') {
    return !hasAccess.value || isCheckingPermission.value
  }
  return false // hide模式下不需要禁用
})

// 是否显示回退内容
const shouldShowFallback = computed(() => {
  if (!props.fallback)
    return false

  if (props.mode === 'hide') {
    return !hasAccess.value && !isCheckingPermission.value
  }
  else {
    return !hasAccess.value
  }
})

// 监听资源和操作变化
watch(
  () => [props.resource, props.action],
  () => {
    checkAccessPermission()
  },
)

// 监听异步模式变化
watch(
  () => props.useAsync,
  () => {
    checkAccessPermission()
  },
)

// 组件挂载时检查权限
onMounted(() => {
  checkAccessPermission()
})
</script>

<template>
  <div v-if="shouldShow" :class="{ 'opacity-50 pointer-events-none': shouldDisable }">
    <!-- 加载状态 -->
    <div v-if="isCheckingPermission && showLoading" class="flex items-center justify-center p-4">
      <div class="w-6 h-6 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin" />
      <span class="ml-2 text-sm text-gray-500">检查权限中...</span>
    </div>

    <!-- 有权限时显示内容 -->
    <slot v-else-if="hasAccess" />

    <!-- 回退内容 -->
    <slot v-else-if="shouldShowFallback" name="fallback" />
  </div>

  <!-- hide模式下的回退内容 -->
  <slot v-else-if="shouldShowFallback" name="fallback" />
</template>
