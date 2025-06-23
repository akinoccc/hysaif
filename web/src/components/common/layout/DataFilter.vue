<script setup lang="ts" generic="T extends Record<string, any>">
import type { FilterField } from '.'
import {
  ChevronDown,
  FileSearch,
  RotateCcw,
  Search,
  SlidersHorizontal,
  X,
} from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { RangeDatePicker } from '@/components'
import {
  Button,
} from '@/components/ui/button'
import {
  Input,
} from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { formatDate } from '@/utils/date'

interface SortOption {
  value: string
  label: string
}

interface ActiveFilter {
  key: string
  label: string
  value: string
  displayValue: string
}

const props = defineProps<{
  modelValue: T
  totalItems: number
  searchPlaceholder?: string
  quickFilters?: FilterField[]
  advancedFilters?: FilterField[]
  sortOptions?: SortOption[]
  resultText?: string
  clearAllText?: string
  advancedSearchText?: string
  resetText?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: T]
  'reset': []
}>()

// 本地状态
const showAdvancedSearch = ref(false)

// 计算属性
const hasActiveFilters = computed(() => {
  const filters = props.modelValue
  return Object.entries(filters).some(([key, value]) => {
    if (key === 'searchQuery')
      return value && value.trim() !== ''
    if (Array.isArray(value))
      return value && value.length > 0
    if (typeof value === 'string')
      return value && value !== 'all' && value !== ''
    return value != null
  })
})

// 获取活跃筛选条件
const activeFilters = computed((): ActiveFilter[] => {
  const filters = props.modelValue
  const active: ActiveFilter[] = []

  // 处理搜索查询
  if (filters.searchQuery && filters.searchQuery.trim() !== '') {
    active.push({
      key: 'searchQuery',
      label: '搜索',
      value: filters.searchQuery,
      displayValue: filters.searchQuery,
    })
  }

  // 处理快速筛选
  props.quickFilters?.forEach((field) => {
    const value = filters[field.key]
    if (value && value !== 'all' && value !== '') {
      const option = field.options?.find(opt => opt.value === value)
      active.push({
        key: field.key,
        label: field.label,
        value,
        displayValue: option?.label || value,
      })
    }
  })

  // 处理高级筛选
  props.advancedFilters?.forEach((field) => {
    const value = filters[field.key]
    if (field.type === 'user-search' && value && value.trim() !== '') {
      active.push({
        key: field.key,
        label: field.label,
        value,
        displayValue: value,
      })
    }
    else if (field.type === 'date-range' && value && Array.isArray(value) && value.length > 0) {
      active.push({
        key: field.key,
        label: field.label,
        value: value.join('~'),
        displayValue: `${formatDate(value[0])} ~ ${formatDate(value[1])}`,
      })
    }
  })

  return active
})

// 方法
function updateFilter(key: keyof T, value: any) {
  emit('update:modelValue', {
    ...props.modelValue,
    [key]: value,
  })
}

function resetFilters() {
  emit('reset')
}

function removeFilter(key: string) {
  if (key === 'searchQuery') {
    updateFilter(key as keyof T, '')
  }
  else {
    const field = [...(props.quickFilters || []), ...(props.advancedFilters || [])]
      .find(f => f.key === key)

    if (field?.type === 'date-range') {
      updateFilter(key as keyof T, undefined)
    }
    else if (field?.type === 'select') {
      updateFilter(key as keyof T, 'all')
    }
    else {
      updateFilter(key as keyof T, '')
    }
  }
}

// 渲染选择器选项
function renderSelectOptions(field: FilterField) {
  if (field.groups && field.groups.length > 0) {
    // 分组选项
    return field.groups.map(group => ({
      group: group.key,
      groupLabel: group.label,
      groupIcon: group.icon,
      options: field.options?.filter(opt => opt.group === group.key) || [],
    }))
  }
  return [{ group: undefined, options: field.options || [] }]
}
</script>

<template>
  <!-- 主搜索区域 -->
  <div class="space-y-6">
    <div class="flex flex-col space-y-4">
      <!-- 搜索框 -->
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-3">
          <Search class="h-5 w-5" />
        </div>
        <Input
          :model-value="modelValue.searchQuery"
          :placeholder="searchPlaceholder || '搜索...'"
          class="pl-10"
          @update:model-value="updateFilter('searchQuery', $event)"
        />
      </div>

      <!-- 快速筛选工具栏 -->
      <div class="flex flex-wrap items-center gap-3">
        <div class="flex-1 flex flex-wrap gap-2">
          <!-- 快速筛选器 -->
          <template v-for="field in quickFilters" :key="field.key">
            <Select
              v-if="field.type === 'select'"
              :model-value="modelValue[field.key]"
              class="min-w-[120px]"
              @update:model-value="updateFilter(field.key, $event)"
            >
              <SelectTrigger>
                <div class="flex items-center gap-2">
                  <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                  <SelectValue :placeholder="field.placeholder || field.label" />
                </div>
              </SelectTrigger>
              <SelectContent>
                <template v-for="group in renderSelectOptions(field)" :key="group.group || 'default'">
                  <template v-if="group.group">
                    <SelectGroup>
                      <SelectLabel class="flex items-center gap-2">
                        <component :is="group.groupIcon" v-if="group.groupIcon" class="h-4 w-4" />
                        {{ group.groupLabel }}
                      </SelectLabel>
                      <SelectItem
                        v-for="option in group.options"
                        :key="option.value"
                        :value="option.value"
                      >
                        {{ option.label }}
                      </SelectItem>
                    </SelectGroup>
                  </template>
                  <template v-else>
                    <SelectItem
                      v-for="option in group.options"
                      :key="option.value"
                      :value="option.value"
                    >
                      {{ option.label }}
                    </SelectItem>
                  </template>
                </template>
              </SelectContent>
            </Select>
          </template>
        </div>

        <div class="flex items-center gap-2">
          <!-- 高级筛选切换按钮 -->
          <Button
            v-if="advancedFilters && advancedFilters.length > 0"
            variant="outline"
            @click="showAdvancedSearch = !showAdvancedSearch"
          >
            <SlidersHorizontal class="h-4 w-4" />
            {{ advancedSearchText || '高级筛选' }}
            <ChevronDown
              class="h-4 w-4 transition-transform duration-200"
              :class="[
                showAdvancedSearch ? 'rotate-180' : '',
              ]"
            />
          </Button>

          <!-- 重置按钮 -->
          <Button
            variant="ghost"
            :disabled="!hasActiveFilters"
            :class="{ 'opacity-50 cursor-not-allowed': !hasActiveFilters }"
            @click="resetFilters"
          >
            <RotateCcw class="h-4 w-4" />
            {{ resetText || '重置' }}
          </Button>
        </div>
      </div>
    </div>

    <!-- 高级筛选区域（可折叠） -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 max-h-0"
      enter-to-class="opacity-100 max-h-[500px]"
      leave-active-class="transition-all duration-300 ease-in"
      leave-from-class="opacity-100 max-h-[500px]"
      leave-to-class="opacity-0 max-h-0"
    >
      <div v-if="showAdvancedSearch && advancedFilters && advancedFilters.length > 0">
        <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-3 gap-5">
          <!-- 高级筛选字段 -->
          <template v-for="field in advancedFilters" :key="field.key">
            <!-- 排序选择器 -->
            <div v-if="field.type === 'select'" class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                {{ field.label }}
              </label>
              <Select
                :model-value="modelValue[field.key]"
                @update:model-value="updateFilter(field.key, $event)"
              >
                <SelectTrigger class="w-full">
                  <SelectValue :placeholder="field.placeholder || field.label" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in field.options"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- 用户搜索 -->
            <div v-else-if="field.type === 'user-search'" class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                {{ field.label }}
              </label>
              <div class="relative w-full">
                <div class="absolute inset-y-0 left-0 flex items-center pl-3">
                  <Search class="h-4 w-4" />
                </div>
                <Input
                  :model-value="modelValue[field.key]"
                  :placeholder="field.placeholder || `搜索${field.label}`"
                  class="pl-10"
                  @update:model-value="updateFilter(field.key, $event)"
                />
              </div>
            </div>

            <!-- 日期范围 -->
            <div v-else-if="field.type === 'date-range'" class="space-y-2 w-full">
              <label class="text-sm font-medium flex items-center gap-2">
                <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                {{ field.label }}
              </label>
              <RangeDatePicker
                :model-value="modelValue[field.key]"
                :placeholder="field.placeholder || `选择${field.label}`"
                @update:model-value="updateFilter(field.key, $event)"
              />
            </div>

            <!-- 普通搜索 -->
            <div v-else-if="field.type === 'text'" class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                {{ field.label }}
              </label>
              <div class="relative w-full">
                <div class="absolute inset-y-0 left-0 flex items-center pl-3">
                  <component :is="field.icon" v-if="field.icon" class="h-4 w-4" />
                </div>
                <Input
                  :model-value="modelValue[field.key]"
                  :placeholder="field.placeholder || `搜索${field.label}`"
                  :class="field.icon ? 'pl-10' : ''"
                  @update:model-value="updateFilter(field.key, $event)"
                />
              </div>
            </div>
          </template>
        </div>
      </div>
    </Transition>

    <!-- 筛选结果统计 -->
    <div
      v-if="hasActiveFilters || totalItems > 0"
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 text-sm text-muted-foreground border-t border-border/40 pt-4"
    >
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <FileSearch class="h-4 w-4" />
          <span>{{ resultText || `共找到 ${totalItems} 个结果` }}</span>
        </div>
        <Button
          v-if="hasActiveFilters"
          variant="link"
          size="sm"
          class="h-auto p-0 text-xs text-primary hover:text-primary/80"
          @click="resetFilters"
        >
          {{ clearAllText || '清除所有筛选条件' }}
        </Button>
      </div>

      <!-- 活跃筛选条件标签 -->
      <div v-if="hasActiveFilters" class="flex flex-wrap gap-2">
        <span
          v-for="filter in activeFilters"
          :key="filter.key"
          class="inline-flex items-center px-3 py-1 rounded-full bg-primary/10 text-primary text-xs font-medium"
        >
          {{ filter.label }}: {{ filter.displayValue }}
          <button
            class="ml-1.5 hover:text-primary/70 focus:outline-none"
            @click="removeFilter(filter.key)"
          >
            <X class="h-3.5 w-3.5" />
          </button>
        </span>
      </div>
    </div>
  </div>
</template>
