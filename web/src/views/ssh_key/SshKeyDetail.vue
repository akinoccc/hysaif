<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { Download, Terminal } from 'lucide-vue-next'
import { SecretDetail, SecretInput } from '@/components'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'

// 下载密钥文件
function downloadKey(keyType: 'private' | 'public', item: SecretItem) {
  if (!item?.data) {
    return
  }

  const keyData = keyType === 'private' ? item.data.private_key : item.data.public_key
  if (!keyData || keyData === '未填写') {
    return
  }

  const fileName = keyType === 'private'
    ? `${item.name}_private_key.pem`
    : `${item.name}_public_key.pub`

  const blob = new Blob([keyData], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
</script>

<template>
  <SecretDetail
    title="SSH 密钥详情" description="查看和管理 SSH 密钥信息" secret-type="ssh_key" secret-data-title="密钥信息"
    :secret-data-icon="Terminal" error-text="SSH 密钥信息"
  >
    <template #secret-actions="{ item }">
      <div class="flex items-center space-x-2">
        <Button
          v-if="item?.data?.private_key"
          variant="outline"
          size="sm"
          @click="downloadKey('private', item)"
        >
          <Download class="h-4 w-4" />
          私钥
        </Button>
        <Button
          v-if="item?.data?.public_key"
          variant="outline"
          size="sm"
          @click="downloadKey('public', item)"
        >
          <Download class="h-4 w-4" />
          公钥
        </Button>
      </div>
    </template>
    <template #secret-data="{ item }">
      <!-- 私钥 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">私钥</Label>
        <div class="mt-1">
          <SecretInput
            readonly
            placeholder="未填写"
            toggleable
            copyable
            :model-value="item?.data?.private_key"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- 公钥 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">公钥</Label>
        <div class="mt-1">
          <SecretInput
            readonly
            placeholder="未填写"
            toggleable
            copyable
            :model-value="item?.data?.public_key"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- 密码短语 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">密码短语</Label>
        <div class="mt-1">
          <SecretInput
            readonly toggleable copyable :model-value="item?.data?.passphrase || '未填写'"
            class=" cursor-not-allowed"
          />
        </div>
      </div>

      <!-- 备注 -->
      <div class="space-y-2">
        <Label class="text-sm font-medium text-muted-foreground">备注</Label>
        <p class="text-sm mt-1 p-3 bg-muted rounded-md">
          {{ item?.data?.notes || '未填写' }}
        </p>
      </div>
    </template>
  </SecretDetail>
</template>
