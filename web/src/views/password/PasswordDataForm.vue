<script setup lang="ts">
import { SecretInput } from '@/components'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
</script>

<template>
  <div class="space-y-4">
    <!-- 用户名 -->
    <FormField v-slot="{ componentField }" name="data.username">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>用户名
        </FormLabel>
        <FormControl>
          <Input
            :model-value="componentField.modelValue" placeholder="输入用户名"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          登录账号的用户名或邮箱
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 密码 -->
    <FormField v-slot="{ componentField }" name="data.password">
      <FormItem>
        <FormLabel class="gap-0.5">
          <span class="text-red-500">*</span>密码
        </FormLabel>
        <FormControl>
          <SecretInput
            :model-value="componentField.modelValue" toggleable placeholder="输入密码"
            @update:model-value="componentField.onChange"
          />
        </FormControl>
        <FormDescription>
          账号的登录密码
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- 网址 -->
    <FormField v-slot="{ componentField }" name="data.address">
      <FormItem>
        <FormLabel>地址</FormLabel>
        <FormControl>
          <Input v-bind="componentField" placeholder="输入网址/地址" />
        </FormControl>
        <FormDescription>
          账号所属的域名或地址
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
          记录账号用途、权限范围等重要信息
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>
  </div>
</template>
