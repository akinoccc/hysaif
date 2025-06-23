<script setup lang="ts">
import { SecretInput } from '@/components'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
</script>

<template>
  <div class="space-y-4">
    <!-- 令牌值 -->
    <FormField v-slot="{ componentField }" name="data.token">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>令牌值
        </FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入令牌值"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          主要的令牌或访问令牌值
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

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
