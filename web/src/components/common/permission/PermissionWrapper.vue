<script setup lang="ts">
import { computed, watch } from 'vue'
import { usePermissionStore } from '@/stores'

interface Props {
  resource: string
  action: string
  mode?: 'hide' | 'disable' // 默认为 hide
  fallback?: boolean // 是否显示fallback内容
  fallbackValue?: boolean // 权限检查失败时的回退值，默认为 false
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'hide',
  fallback: false,
  fallbackValue: false,
})

const { hasPermission, usePermissionState } = usePermissionStore()

// 使用响应式权限状态
const { hasAccess } = usePermissionState(
  props.resource,
  props.action,
  props.fallbackValue,
)

// 是否显示内容
const shouldShow = computed(() => {
  if (props.mode === 'hide') {
    return hasAccess.value
  }
  return true // disable模式下总是显示
})

// 是否禁用内容
const shouldDisable = computed(() => {
  if (props.mode === 'disable') {
    return !hasAccess.value
  }
  return false // hide模式下不需要禁用
})

// 是否显示回退内容
const shouldShowFallback = computed(() => {
  if (!props.fallback)
    return false

  if (props.mode === 'hide') {
    return !hasAccess.value
  }
  else {
    return !hasAccess.value
  }
})

// 监听资源和操作变化
watch(
  () => [props.resource, props.action],
  () => {
    hasPermission(props.resource, props.action)
  },
)

// 权限现在通过响应式状态自动更新，无需手动检查
</script>

<template>
  <div v-if="shouldShow" :class="{ 'opacity-50 pointer-events-none': shouldDisable }">
    <!-- 有权限时显示内容 -->
    <slot v-if="hasAccess" />

    <!-- 回退内容 -->
    <slot v-else-if="shouldShowFallback" name="fallback" />
  </div>

  <!-- hide模式下的回退内容 -->
  <slot v-else-if="shouldShowFallback" name="fallback" />
</template>
