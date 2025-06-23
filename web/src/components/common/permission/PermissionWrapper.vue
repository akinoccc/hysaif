<script setup lang="ts">
import { computed } from 'vue'
import { usePermission } from '@/composables/usePermission'

interface Props {
  resource: string
  action: string
  mode?: 'hide' | 'disable' // 默认为 hide
  fallback?: boolean // 是否显示fallback内容
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'hide',
  fallback: false,
})

const { hasPermission } = usePermission()

// 检查是否有权限
const hasAccess = computed(() => {
  return hasPermission(props.resource, props.action)
})

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
</script>

<template>
  <div v-if="shouldShow" :class="{ 'opacity-50 pointer-events-none': shouldDisable }">
    <slot v-if="hasAccess" />
    <slot v-else-if="fallback" name="fallback" />
  </div>
  <slot v-else-if="fallback && !shouldShow" name="fallback" />
</template>
