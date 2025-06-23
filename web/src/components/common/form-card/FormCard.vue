<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface Props {
  title: string
  description: string
  icon: any
  class?: HTMLAttributes['class']
}

defineOptions({
  inheritAttrs: false,
})
const props = defineProps<Props>()
</script>

<template>
  <Card :class="cn('', props.class)">
    <CardHeader>
      <CardTitle class="flex items-center gap-1">
        <component :is="icon" class=" h-5 w-5" />
        {{ title }}
      </CardTitle>
      <CardDescription>
        {{ description }}
      </CardDescription>
    </CardHeader>
    <CardContent>
      <slot />
    </CardContent>
  </Card>
</template>
