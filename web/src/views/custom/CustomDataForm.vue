<script setup lang="ts">
import type { FormContext } from 'vee-validate'
import { Plus, Trash2 } from 'lucide-vue-next'
import { computed } from 'vue'
import { SecretInput } from '@/components'
import { Button } from '@/components/ui/button'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'

interface CustomDataItem {
  key: string
  value: string
}

interface Props {
  form: FormContext<any>
}

const props = defineProps<Props>()
const { form } = props

const customDataArray = computed({
  get: () => {
    const data = form.values.data?.custom_data
    return Array.isArray(data) ? data : [{ key: '', value: '' }]
  },
  set: (value: CustomDataItem[]) => {
    form.setFieldValue('data.custom_data', value)
  },
})

function addField() {
  const currentData = [...customDataArray.value]
  currentData.push({ key: '', value: '' })
  customDataArray.value = currentData
}

function removeField(index: number) {
  if (customDataArray.value.length > 1) {
    const currentData = [...customDataArray.value]
    currentData.splice(index, 1)
    customDataArray.value = currentData
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- 自定义字段列表 -->
    <div class="space-y-4">
      <div v-for="(_, index) in customDataArray" :key="index" class="p-4 border rounded-md bg-card">
        <div class="flex justify-between items-center mb-2">
          <h4 class="text-sm font-medium">
            字段 {{ index + 1 }}
          </h4>
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="removeField(index)">
            <Trash2 class="h-4 w-4" />
          </Button>
        </div>

        <div class="grid gap-4 grid-cols-1 xl:grid-cols-2">
          <!-- 字段名称 -->
          <FormField v-slot="{ componentField }" :name="`data.custom_data[${index}].key`">
            <FormItem>
              <FormLabel class="gap-0.5">
                <span class="text-red-500">*</span>字段名称
              </FormLabel>
              <FormControl>
                <Input v-bind="componentField" placeholder="输入字段名称" />
              </FormControl>
              <FormDescription>
                自定义字段的名称或标识符
              </FormDescription>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- 字段值 -->
          <FormField v-slot="{ componentField }" :name="`data.custom_data[${index}].value`">
            <FormItem>
              <FormLabel class="gap-0.5">
                <span class="text-red-500">*</span>字段值
              </FormLabel>
              <FormControl>
                <SecretInput
                  v-model="componentField.modelValue" toggleable
                  placeholder="输入字段值" @update:model-value="componentField.onChange"
                />
              </FormControl>
              <FormDescription>
                自定义字段的值或内容
              </FormDescription>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>
      </div>
    </div>

    <!-- 添加字段按钮 -->
    <Button type="button" variant="outline" class="w-full" @click="addField">
      <Plus class=" h-4 w-4" />
      添加字段
    </Button>

    <!-- 备注 -->
    <FormField v-slot="{ componentField }" name="data.notes">
      <FormItem>
        <FormLabel>备注</FormLabel>
        <FormControl>
          <Textarea v-bind="componentField" placeholder="添加使用说明或注意事项" rows="3" />
        </FormControl>
        <FormDescription>
          记录使用方法、权限范围等重要信息
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>
  </div>
</template>
