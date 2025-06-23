<script setup lang="ts">
import type { WebAuthnCredentialResponse } from '@/api/types'
import { AlertCircle, Fingerprint, Loader2, Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuthStore } from '@/stores/auth'
import {
  createCredential,
  isWebAuthnSupported,
  prepareRegistrationOptions,
  prepareRegistrationResponse,
} from '@/utils/webauthn'

const authStore = useAuthStore()

const credentials = ref<WebAuthnCredentialResponse[]>([])
const loading = ref(false)
const registering = ref(false)
const error = ref('')
const showRegisterDialog = ref(false)
const credentialName = ref('')
const webAuthnSupported = ref(false)
const showDeleteDialog = ref(false)
const deleteTarget = ref<string>('')

// 获取凭证列表
async function fetchCredentials() {
  loading.value = true
  error.value = ''

  try {
    const result = await authStore.getWebAuthnCredentials()
    if (result.success) {
      credentials.value = result.credentials || []
    }
    else {
      error.value = result.message || '获取 Passkey 列表失败'
    }
  }
  catch (err) {
    error.value = '获取 Passkey 列表失败'
  }
  finally {
    loading.value = false
  }
}

// 注册新的 Passkey
async function registerPasskey() {
  if (!credentialName.value.trim()) {
    error.value = '请输入 Passkey 名称'
    return
  }

  registering.value = true
  error.value = ''

  try {
    // 开始注册
    const beginResult = await authStore.beginWebAuthnRegistration(credentialName.value)
    if (!beginResult.success || !beginResult.options) {
      throw new Error(beginResult.message || '开始注册失败')
    }

    // 准备注册选项
    const options = prepareRegistrationOptions(beginResult.options.publicKey)

    // 创建凭证
    const credential = await createCredential(options)

    // 准备响应数据
    const response = prepareRegistrationResponse(credential)

    // 完成注册
    const finishResult = await authStore.finishWebAuthnRegistration(response, credentialName.value)
    if (finishResult.success) {
      showRegisterDialog.value = false
      credentialName.value = ''
      await fetchCredentials()
    }
    else {
      throw new Error(finishResult.message || '注册失败')
    }
  }
  catch (err: any) {
    error.value = err.message || 'Passkey 注册失败'
    console.error('Passkey registration error:', err)
  }
  finally {
    registering.value = false
  }
}

// 删除 Passkey
async function deletePasskey(id: string) {
  deleteTarget.value = id
  showDeleteDialog.value = true
}

// 确认删除
async function confirmDelete() {
  if (!deleteTarget.value)
    return

  loading.value = true
  error.value = ''
  showDeleteDialog.value = false

  try {
    const result = await authStore.deleteWebAuthnCredential(deleteTarget.value)
    if (result.success) {
      await fetchCredentials()
    }
    else {
      error.value = result.message || '删除 Passkey 失败'
    }
  }
  catch (err) {
    error.value = '删除 Passkey 失败'
  }
  finally {
    loading.value = false
    deleteTarget.value = ''
  }
}

// 格式化日期
function formatDate(dateString: string) {
  if (!dateString)
    return '从未使用'

  const date = new Date(dateString)
  if (Number.isNaN(date.getTime()))
    return '从未使用'

  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

onMounted(() => {
  webAuthnSupported.value = isWebAuthnSupported()
  if (webAuthnSupported.value) {
    fetchCredentials()
  }
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <Fingerprint class="h-5 w-5" />
        Passkey 管理
      </CardTitle>
      <CardDescription>
        使用 Passkey 可以更安全、更便捷地登录系统，无需记住密码
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- 浏览器不支持提示 -->
      <Alert v-if="!webAuthnSupported" variant="destructive">
        <AlertCircle class="h-4 w-4" />
        <AlertDescription>
          您的浏览器不支持 Passkey 功能。请使用最新版本的 Chrome、Safari、Edge 或 Firefox。
        </AlertDescription>
      </Alert>

      <!-- 错误提示 -->
      <Alert v-if="error && webAuthnSupported" variant="destructive">
        <AlertCircle class="h-4 w-4" />
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <!-- 添加按钮 -->
      <div v-if="webAuthnSupported" class="flex justify-end">
        <Button :disabled="loading || registering" @click="showRegisterDialog = true">
          <Plus class="mr-2 h-4 w-4" />
          添加 Passkey
        </Button>
      </div>

      <!-- Passkey 列表 -->
      <div v-if="webAuthnSupported" class="space-y-2">
        <div v-if="loading && credentials.length === 0" class="text-center py-8">
          <Loader2 class="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
          <p class="mt-2 text-sm text-muted-foreground">
            加载中...
          </p>
        </div>

        <div v-else-if="credentials.length === 0" class="text-center py-8">
          <Fingerprint class="h-12 w-12 mx-auto text-muted-foreground/50" />
          <p class="mt-2 text-sm text-muted-foreground">
            您还没有添加任何 Passkey
          </p>
        </div>

        <div
          v-for="credential in credentials"
          v-else
          :key="credential.id"
          class="flex items-center justify-between p-4 border rounded-lg"
        >
          <div class="flex-1">
            <div class="font-medium">
              {{ credential.name }}
            </div>
            <div class="text-sm text-muted-foreground">
              创建时间：{{ formatDate(credential.created_at) }}
            </div>
            <div class="text-sm text-muted-foreground">
              最后使用：{{ formatDate(credential.last_used_at) }}
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            :disabled="loading"
            @click="deletePasskey(credential.id)"
          >
            <Trash2 class="h-4 w-4" />
          </Button>
        </div>
      </div>

      <!-- 注册对话框 -->
      <Dialog v-model:open="showRegisterDialog">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>添加新的 Passkey</DialogTitle>
            <DialogDescription>
              给您的 Passkey 起一个容易识别的名称，例如 "MacBook Pro" 或 "iPhone"
            </DialogDescription>
          </DialogHeader>
          <div class="space-y-4 py-4">
            <div class="space-y-2">
              <Label for="name">Passkey 名称</Label>
              <Input
                id="name"
                v-model="credentialName"
                placeholder="例如：我的 MacBook"
                @keyup.enter="registerPasskey"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" :disabled="registering" @click="showRegisterDialog = false">
              取消
            </Button>
            <Button :disabled="registering" @click="registerPasskey">
              <Loader2 v-if="registering" class="mr-2 h-4 w-4 animate-spin" />
              {{ registering ? '注册中...' : '确定' }}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <!-- 删除确认对话框 -->
      <Dialog v-model:open="showDeleteDialog">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>确认删除</DialogTitle>
            <DialogDescription>
              确定要删除这个 Passkey 吗？删除后您将无法使用它登录。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" @click="showDeleteDialog = false">
              取消
            </Button>
            <Button variant="destructive" @click="confirmDelete">
              删除
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </CardContent>
  </Card>
</template>
