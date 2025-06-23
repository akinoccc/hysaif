<script setup lang="ts">
import type { DateRange, DateValue } from 'reka-ui'
import {
  CalendarDate,
  DateFormatter,
  getLocalTimeZone,
} from '@internationalized/date'

import { CalendarIcon } from 'lucide-vue-next'
import { ref, type Ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { cn } from '@/lib/utils'

interface Props {
  modelValue?: [number, number] // timestamp in milliseconds
  placeholder?: string
}

interface Emits {
  (e: 'update:modelValue', value: [number, number] | undefined): void
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '选择日期',
})

const emit = defineEmits<Emits>()

const df = new DateFormatter('zh-CN', {
  dateStyle: 'long',
})

const value = ref({
  start: props.modelValue?.[0] ? new CalendarDate(new Date(props.modelValue[0]).getFullYear(), new Date(props.modelValue[0]).getMonth(), new Date(props.modelValue[0]).getDate()) : undefined,
  end: props.modelValue?.[1] ? new CalendarDate(new Date(props.modelValue[1]).getFullYear(), new Date(props.modelValue[1]).getMonth(), new Date(props.modelValue[1]).getDate()) : undefined,
}) as Ref<DateRange>

watch(value, (newValue) => {
  if (newValue) {
    emit('update:modelValue', [newValue.start!.toDate(getLocalTimeZone()).getTime(), newValue.end!.toDate(getLocalTimeZone()).getTime()])
  }
  else {
    emit('update:modelValue', undefined)
  }
})
</script>

<template>
  <Popover>
    <PopoverTrigger as-child>
      <Button
        variant="outline"
        :class="cn(
          'min-w-[280px] justify-start text-left font-normal',
          !value && 'text-muted-foreground',
        )"
      >
        <CalendarIcon class=" h-4 w-4" />
        <template v-if="value?.start">
          <template v-if="value?.end">
            {{ df.format(value?.start.toDate(getLocalTimeZone())) }} - {{ df.format(value?.end.toDate(getLocalTimeZone())) }}
          </template>

          <template v-else>
            {{ df.format(value?.start.toDate(getLocalTimeZone())) }}
          </template>
        </template>
        <template v-else>
          {{ placeholder }}
        </template>
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-auto p-0">
      <RangeCalendar
        v-model="value"
        initial-focus
        :number-of-months="2"
        @update:start-value="(startDate?: DateValue) => value!.start = startDate"
        @update:end-value="(endDate?: DateValue) => value!.end = endDate"
      />
    </PopoverContent>
  </Popover>
</template>
