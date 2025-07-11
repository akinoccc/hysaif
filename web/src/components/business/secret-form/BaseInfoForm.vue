<script setup lang="ts">
import type { SecretBaseInfo } from '@/api/types'
import type { SecretItemTypeT } from '@/constants/secretItem'
import { X } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { DatePicker } from '@/components/common/date-picker'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { getCategoriesByGroup, getCategoryGroupsByType, SECRET_ITEM_CATEGORY_GROUPS, SECRET_ITEM_TYPE_MAP, type SecretItemCategoryGroupKey } from '@/constants/secretItem'

const props = defineProps<{
  modelValue: Partial<SecretBaseInfo>
  type: keyof typeof SECRET_ITEM_TYPE_MAP
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<SecretBaseInfo>): void
}>()

const typeInfo = computed(() => {
  return SECRET_ITEM_TYPE_MAP[props.type]
})

// 标签输入
const tagInput = ref<HTMLInputElement | null>(null)
const newTag = ref('')

function addTag() {
  if (newTag.value.trim()) {
    const updatedTags = [...(props.modelValue.tags || []), newTag.value.trim()]
    emit('update:modelValue', { ...props.modelValue, tags: updatedTags })
    newTag.value = ''
  }
}

function removeTag(index: number) {
  const updatedTags = [...(props.modelValue.tags || [])]
  updatedTags.splice(index, 1)
  emit('update:modelValue', { ...props.modelValue, tags: updatedTags })
}

// 分类相关函数 - 根据类型筛选分类
function getFilteredCategoriesByGroup(groupKey: SecretItemCategoryGroupKey) {
  if (props.type) {
    return getCategoriesByGroup(groupKey, props.type as SecretItemTypeT)
  }
  return getCategoriesByGroup(groupKey)
}

// 根据类型获取可用的分类组
const availableGroups = computed(() => {
  if (props.type) {
    return getCategoryGroupsByType(props.type as SecretItemTypeT)
  }
  return Object.values(SECRET_ITEM_CATEGORY_GROUPS)
})
</script>

<template>
  <div class="space-y-4">
    <!-- 名称 -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>名称
        </FormLabel>
        <FormControl>
          <Input v-bind="componentField" placeholder="输入密钥名称" />
        </FormControl>
        <FormDescription>
          为您的{{ typeInfo.label }}提供一个描述性名称
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 描述 -->
    <FormField v-slot="{ componentField }" name="description">
      <FormItem>
        <FormLabel>描述</FormLabel>
        <FormControl>
          <Textarea v-bind="componentField" placeholder="输入密钥描述信息" rows="3" />
        </FormControl>
        <FormDescription>
          添加详细描述，帮助您识别此{{ typeInfo.label }}的用途
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 分类 -->
    <FormField v-slot="{ componentField }" name="category">
      <FormItem>
        <FormLabel>
          <span class="text-red-500">*</span>分类
        </FormLabel>
        <FormControl>
          <Select v-bind="componentField">
            <SelectTrigger>
              <SelectValue placeholder="选择分类" />
            </SelectTrigger>
            <SelectContent>
              <template v-for="group in availableGroups" :key="group.key">
                <SelectGroup>
                  <SelectLabel class="flex items-center gap-2">
                    <component :is="group.icon" class="h-4 w-4" />
                    {{ group.label }}
                  </SelectLabel>
                  <SelectItem
                    v-for="category in getFilteredCategoriesByGroup(group.key)" :key="category.key"
                    :value="category.key"
                  >
                    {{ category.label }}
                  </SelectItem>
                </SelectGroup>
              </template>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          选择密钥所属的服务或平台分类
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 环境 -->
    <FormField v-slot="{ componentField }" name="environment">
      <FormItem>
        <FormLabel>
          <span class="text-red-500">*</span>环境
        </FormLabel>
      </FormItem>
      <FormControl>
        <Select v-bind="componentField">
          <SelectTrigger>
            <SelectValue placeholder="选择环境" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="production">
              生产环境
            </SelectItem>
            <SelectItem value="staging">
              预发环境
            </SelectItem>
            <SelectItem value="test">
              测试环境
            </SelectItem>
            <SelectItem value="development">
              开发环境
            </SelectItem>
            <SelectItem value="local">
              本地环境
            </SelectItem>
          </SelectContent>
        </Select>
      </FormControl>
      <FormDescription>
        选择{{ typeInfo.label }}所属的环境
      </FormDescription>
      <FormMessage />
    </FormField>

    <!-- 标签 -->
    <FormField name="tags">
      <FormItem>
        <FormLabel>标签</FormLabel>
        <FormControl>
          <div class="flex gap-2 flex-wrap">
            <div
              v-for="(tag, index) in modelValue.tags" :key="index"
              class="flex items-center w-fit bg-primary/10 text-primary px-2 py-1 rounded-md"
            >
              <span class="text-sm">{{ tag }}</span>
              <button type="button" class="ml-1 text-primary hover:text-primary/80" @click="removeTag(index)">
                <X class="h-3 w-3" />
              </button>
            </div>
          </div>
          <Input
            ref="tagInput" v-model="newTag" class="flex-1 min-w-[100px] outline-none bg-transparent theme-transition"
            placeholder="输入标签并按回车添加" @keydown.enter.prevent="addTag"
          />
        </FormControl>
        <FormDescription>
          添加标签以便更好地组织和搜索
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 过期时间 -->
    <FormField v-slot="{ componentField }" name="expires_at">
      <FormItem>
        <FormLabel>过期时间</FormLabel>
        <FormControl>
          <DatePicker v-bind="componentField" class="w-full" />
        </FormControl>
        <FormDescription>
          设置{{ typeInfo.label }}的过期时间，到期前会收到提醒
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>
  </div>
</template>
