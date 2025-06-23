<script setup lang="ts">
import {
  DateFormatter,
  type DateValue,
  fromDate,
  getLocalTimeZone,
} from '@internationalized/date'
import { CalendarIcon } from 'lucide-vue-next'

import { computed, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

interface Props {
  modelValue?: number // timestamp in milliseconds
  placeholder?: string
}

interface Emits {
  (e: 'update:modelValue', value: number | undefined): void
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '选择日期',
})

const emit = defineEmits<Emits>()

const df = new DateFormatter('zh-CN', {
  dateStyle: 'long',
})

const value = ref<DateValue>()

// Convert timestamp to DateValue when props.modelValue changes
watch(() => props.modelValue, (newTimestamp) => {
  if (newTimestamp) {
    const date = new Date(newTimestamp)
    value.value = fromDate(date, getLocalTimeZone())
  }
  else {
    value.value = undefined
  }
}, { immediate: true })

// Convert DateValue to timestamp when value changes
watch(value, (newValue) => {
  if (newValue) {
    const timestamp = newValue.toDate(getLocalTimeZone()).getTime()
    emit('update:modelValue', timestamp)
  }
  else {
    emit('update:modelValue', undefined)
  }
})

const displayValue = computed(() => {
  return value.value ? df.format(value.value.toDate(getLocalTimeZone())) : props.placeholder
})
</script>

<template>
  <Popover>
    <PopoverTrigger as-child>
      <Button
        variant="outline" :class="cn(
          'justify-start text-left font-normal',
          !value && 'text-muted-foreground',
        )"
      >
        <CalendarIcon class="h-4 w-4" />
        {{ displayValue }}
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-auto p-0">
      <Calendar v-model="value" initial-focus />
    </PopoverContent>
  </Popover>
</template>
