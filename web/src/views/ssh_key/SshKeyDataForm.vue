<script setup lang="ts">
import { SecretInput } from '@/components'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
</script>

<template>
  <div class="space-y-4">
    <!-- 私钥 -->
    <FormField v-slot="{ componentField }" name="data.private_key">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>私钥
        </FormLabel>
        <FormControl>
          <Textarea v-bind="componentField" placeholder="输入 SSH 私钥" />
        </FormControl>
        <FormDescription>
          SSH 私钥内容，通常以 -----BEGIN RSA PRIVATE KEY----- 开头
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 公钥（可选） -->
    <FormField v-slot="{ componentField }" name="data.public_key">
      <FormItem>
        <FormLabel>公钥（可选）</FormLabel>
        <FormControl>
          <Textarea v-bind="componentField" placeholder="输入 SSH 公钥" />
        </FormControl>
        <FormDescription>
          SSH 公钥内容，通常以 ssh-rsa 开头
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 密码短语（可选） -->
    <FormField v-slot="{ componentField }" name="data.passphrase">
      <FormItem>
        <FormLabel>密码短语（可选）</FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入密钥的密码短语"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          如果私钥设置了密码保护，请输入密码短语
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
