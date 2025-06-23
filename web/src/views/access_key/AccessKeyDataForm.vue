<script setup lang="ts">
import { SecretInput } from '@/components'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
</script>

<template>
  <div class="space-y-4">
    <!-- Access Key -->
    <FormField v-slot="{ componentField }" name="data.access_key">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>Access Key
        </FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入 Access Key"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          主要的访问密钥标识符
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Secret Key -->
    <FormField v-slot="{ componentField }" name="data.secret_key">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>Secret Key
        </FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入 Secret Key"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          与 Access Key 配对的密钥
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 区域 -->
    <FormField v-slot="{ componentField }" name="data.region">
      <FormItem>
        <FormLabel>区域（可选）</FormLabel>
        <FormControl>
          <Input v-bind="componentField" placeholder="输入区域信息，如 us-east-1" />
        </FormControl>
        <FormDescription>
          密钥所适用的区域或地区
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
