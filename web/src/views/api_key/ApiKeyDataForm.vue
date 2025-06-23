<script setup lang="ts">
import { SecretInput } from '@/components'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
</script>

<template>
  <div class="space-y-4">
    <!-- API 密钥 -->
    <FormField v-slot="{ componentField }" name="data.api_key">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>API 密钥
        </FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入 API 密钥"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          主要的 API 密钥或访问令牌
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- API 密钥 ID（可选） -->
    <FormField v-slot="{ componentField }" name="data.api_secret">
      <FormItem>
        <FormLabel>API 密钥 ID / Secret（可选）</FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入 API 密钥 ID 或 Secret"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          某些服务需要的额外密钥或 ID
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- API 端点 -->
    <FormField v-slot="{ componentField }" name="data.endpoint">
      <FormItem>
        <FormLabel>API 端点（可选）</FormLabel>
        <FormControl>
          <Input v-bind="componentField" placeholder="输入 API 端点 URL" />
        </FormControl>
        <FormDescription>
          API 服务的基础 URL 或端点地址
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
