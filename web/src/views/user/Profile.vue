<script setup lang="ts">
import {
  AlertTriangle,
  Calendar,
  CheckCircle,
  Eye,
  EyeOff,
  Fingerprint,
  History,
  Key,
  Loader2,
  LogIn,
  Mail,
  Plus,
  RefreshCw,
  Save,
  Shield,
  Trash2,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { userAPI, webauthnAPI, type WebAuthnCredential } from '@/api'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { formatDate, formatRelativeTime, generatePassword, getRoleDisplayName, validatePasswordStrength } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'
import { createCredential, isWebAuthnSupported } from '@/utils/webauthn'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const profileSubmitting = ref(false)
const passwordSubmitting = ref(false)
const loadingHistory = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// WebAuthn 相关
const webauthnSupported = ref(false)
const credentials = ref<WebAuthnCredential[]>([])
const loadingCredentials = ref(false)
const addingCredential = ref(false)
const deletingCredentialId = ref<string | null>(null)
const showAddCredentialDialog = ref(false)
const newCredentialName = ref('')

const profileForm = ref({
  name: '',
  email: '',
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const loginHistory = ref<any[]>([])

const passwordStrength = computed(() => {
  if (!passwordForm.value.newPassword)
    return null
  return validatePasswordStrength(passwordForm.value.newPassword)
})

const passwordStrengthColor = computed(() => {
  switch (passwordStrength.value?.strength) {
    case '强':
      return 'text-green-600'
    case '中等':
      return 'text-yellow-600'
    case '弱':
      return 'text-red-600'
    default:
      return 'text-gray-600'
  }
})

function getUserAgentInfo(userAgent: string) {
  if (userAgent.includes('Chrome'))
    return 'Chrome'
  if (userAgent.includes('Firefox'))
    return 'Firefox'
  if (userAgent.includes('Safari'))
    return 'Safari'
  if (userAgent.includes('Edge'))
    return 'Edge'
  if (userAgent.includes('Mobile'))
    return '移动设备'
  return '未知浏览器'
}

function generateRandomPassword() {
  passwordForm.value.newPassword = generatePassword(16)
}

async function updateProfile() {
  profileSubmitting.value = true
  try {
    await userAPI.updateProfile(profileForm.value)
    // alert('个人信息更新成功')
  }
  catch (error) {
    console.error('Failed to update profile:', error)
    // alert('更新失败，请重试')
  }
  finally {
    profileSubmitting.value = false
  }
}

async function changePassword() {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    // alert('新密码和确认密码不匹配')
    return
  }

  passwordSubmitting.value = true
  try {
    await userAPI.changePassword({
      current_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword,
    })

    // 清空表单
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    }

    // alert('密码修改成功')
  }
  catch (error) {
    console.error('Failed to change password:', error)
    // alert('密码修改失败，请检查当前密码是否正确')
  }
  finally {
    passwordSubmitting.value = false
  }
}

async function loadLoginHistory() {
  loadingHistory.value = true
  try {
    const response = await userAPI.getLoginHistory()
    loginHistory.value = response.data || []
  }
  catch (error) {
    console.error('Failed to load login history:', error)
  }
  finally {
    loadingHistory.value = false
  }
}

// WebAuthn 相关函数
async function loadCredentials() {
  loadingCredentials.value = true
  try {
    const data = await webauthnAPI.getCredentials()
    credentials.value = data || []
  }
  catch (error) {
    console.error('Failed to load credentials:', error)
  }
  finally {
    loadingCredentials.value = false
  }
}

async function addCredential() {
  if (!newCredentialName.value.trim()) {
    return
  }

  addingCredential.value = true
  try {
    // 开始注册
    const options = await webauthnAPI.beginRegistration(newCredentialName.value)

    // 创建凭证
    const credential = await createCredential(options)

    // 完成注册
    await webauthnAPI.finishRegistration(newCredentialName.value, credential)

    // 重新加载凭证列表
    await loadCredentials()

    // 关闭对话框并清空输入
    showAddCredentialDialog.value = false
    newCredentialName.value = ''
  }
  catch (error: any) {
    console.error('Failed to add credential:', error)
    // alert(error.message || '添加 Passkey 失败')
  }
  finally {
    addingCredential.value = false
  }
}

async function deleteCredential(id: string) {
  deletingCredentialId.value = id
  try {
    await webauthnAPI.deleteCredential(id)
    await loadCredentials()
  }
  catch (error) {
    console.error('Failed to delete credential:', error)
    // alert('删除失败')
  }
  finally {
    deletingCredentialId.value = null
  }
}

onMounted(() => {
  if (user.value) {
    profileForm.value = {
      name: user.value.name,
      email: user.value.email || '',
    }
  }
  loadLoginHistory()

  // 检查 WebAuthn 支持并加载凭证
  webauthnSupported.value = isWebAuthnSupported()
  if (webauthnSupported.value) {
    loadCredentials()
  }
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">
        个人资料
      </h1>
      <p class="text-muted-foreground">
        管理您的账户信息和安全设置
      </p>
    </div>

    <div class="grid gap-6 md:grid-cols-3">
      <!-- 主要内容 -->
      <div class="md:col-span-2 space-y-6">
        <!-- 基本信息 -->
        <Card>
          <div class="px-6 py-2">
            <h3 class="text-lg font-semibold mb-4">
              基本信息
            </h3>
            <form class="space-y-4" @submit.prevent="updateProfile">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="username" class="block text-sm font-medium mb-2">用户名称</label>
                  <Input
                    id="username"
                    v-model="profileForm.name"
                    placeholder="输入用户名称"
                    required
                  />
                </div>
                <div>
                  <label for="email" class="block text-sm font-medium mb-2">邮箱</label>
                  <Input
                    id="email"
                    v-model="profileForm.email"
                    type="email"
                    placeholder="输入邮箱地址"
                    required
                  />
                </div>
              </div>
              <div>
                <Button type="submit" :disabled="profileSubmitting">
                  <Loader2 v-if="profileSubmitting" class=" h-4 w-4 animate-spin" />
                  <Save v-else class=" h-4 w-4" />
                  {{ profileSubmitting ? '保存中...' : '保存修改' }}
                </Button>
              </div>
            </form>
          </div>
        </Card>

        <!-- 修改密码 -->
        <Card>
          <div class="px-6 py-2">
            <h3 class="text-lg font-semibold mb-4">
              修改密码
            </h3>
            <form class="space-y-4" @submit.prevent="changePassword">
              <div>
                <label for="current-password" class="block text-sm font-medium mb-2">当前密码</label>
                <div class="flex space-x-2">
                  <Input
                    id="current-password"
                    v-model="passwordForm.currentPassword"
                    :type="showCurrentPassword ? 'text' : 'password'"
                    placeholder="输入当前密码"
                    required
                  />
                  <Button type="button" variant="outline" @click="showCurrentPassword = !showCurrentPassword">
                    <Eye v-if="!showCurrentPassword" class="h-4 w-4" />
                    <EyeOff v-else class="h-4 w-4" />
                  </Button>
                </div>
              </div>
              <div>
                <label for="new-password" class="block text-sm font-medium mb-2">新密码</label>
                <div class="space-y-2">
                  <div class="flex space-x-2">
                    <Input
                      id="new-password"
                      v-model="passwordForm.newPassword"
                      :type="showNewPassword ? 'text' : 'password'"
                      placeholder="输入新密码"
                      required
                    />
                    <Button type="button" variant="outline" @click="showNewPassword = !showNewPassword">
                      <Eye v-if="!showNewPassword" class="h-4 w-4" />
                      <EyeOff v-else class="h-4 w-4" />
                    </Button>
                    <Button type="button" variant="outline" @click="generateRandomPassword">
                      <RefreshCw class="h-4 w-4" />
                    </Button>
                  </div>
                  <div v-if="passwordStrength" class="text-xs" :class="passwordStrengthColor">
                    密码强度: {{ passwordStrength }}
                  </div>
                </div>
              </div>
              <div>
                <label for="confirm-password" class="block text-sm font-medium mb-2">确认新密码</label>
                <div class="flex space-x-2">
                  <Input
                    id="confirm-password"
                    v-model="passwordForm.confirmPassword"
                    :type="showConfirmPassword ? 'text' : 'password'"
                    placeholder="再次输入新密码"
                    required
                  />
                  <Button type="button" variant="outline" @click="showConfirmPassword = !showConfirmPassword">
                    <Eye v-if="!showConfirmPassword" class="h-4 w-4" />
                    <EyeOff v-else class="h-4 w-4" />
                  </Button>
                </div>
                <div v-if="passwordForm.newPassword && passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword" class="text-xs text-red-600 mt-1">
                  密码不匹配
                </div>
              </div>
              <div>
                <Button
                  type="submit"
                  :disabled="passwordSubmitting || passwordForm.newPassword !== passwordForm.confirmPassword"
                >
                  <Loader2 v-if="passwordSubmitting" class=" h-4 w-4 animate-spin" />
                  <Key v-else class=" h-4 w-4" />
                  {{ passwordSubmitting ? '修改中...' : '修改密码' }}
                </Button>
              </div>
            </form>
          </div>
        </Card>

        <!-- Passkey 管理 -->
        <Card v-if="webauthnSupported">
          <div class="px-6 py-2">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold flex items-center">
                <Fingerprint class="h-5 w-5" />
                Passkey 管理
              </h3>
              <Button size="sm" @click="showAddCredentialDialog = true">
                <Plus class="h-4 w-4" />
                添加 Passkey
              </Button>
            </div>

            <div v-if="loadingCredentials" class="flex items-center justify-center py-4">
              <Loader2 class="h-5 w-5 animate-spin" />
              <span class="text-sm">加载中...</span>
            </div>

            <div v-else-if="credentials.length === 0" class="text-center py-8 text-muted-foreground">
              <Fingerprint class="mx-auto h-12 w-12 mb-2 opacity-50" />
              <p class="text-sm">
                还没有添加任何 Passkey
              </p>
              <p class="text-xs mt-1">
                添加 Passkey 后，您可以使用指纹、面容或设备密码快速登录
              </p>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="credential in credentials"
                :key="credential.id"
                class="flex items-center justify-between p-4 border border-border rounded-lg"
              >
                <div>
                  <p class="font-medium">
                    {{ credential.credential_name || '未命名的 Passkey' }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    创建于 {{ formatDate(credential.created_at) }}
                    <span v-if="credential.last_used_at">
                      · 最后使用 {{ formatRelativeTime(credential.last_used_at) }}
                    </span>
                  </p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  :disabled="deletingCredentialId === credential.id"
                  @click="deleteCredential(credential.id)"
                >
                  <Loader2 v-if="deletingCredentialId === credential.id" class="h-4 w-4 animate-spin" />
                  <Trash2 v-else class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </Card>

        <!-- 登录历史 -->
        <Card>
          <div class="px-6 py-2">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold">
                最近登录记录
              </h3>
              <Button variant="outline" size="sm" @click="loadLoginHistory">
                <RefreshCw class=" h-4 w-4" :class="{ 'animate-spin': loadingHistory }" />
                刷新
              </Button>
            </div>

            <div v-if="loadingHistory" class="flex items-center justify-center py-4">
              <Loader2 class="h-5 w-5 animate-spin" />
              <span class="text-sm">加载中...</span>
            </div>

            <div v-else-if="loginHistory.length === 0" class="text-center py-4 text-muted-foreground">
              <History class="mx-auto h-8 w-8 mb-2" />
              <p class="text-sm">
                暂无登录记录
              </p>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="record in loginHistory"
                :key="record.id"
                class="flex items-center justify-between p-3 border border-border rounded-lg"
              >
                <div class="flex items-center space-x-3">
                  <div class="w-8 h-8 rounded-full bg-green-100 text-green-600 flex items-center justify-center">
                    <LogIn class="h-4 w-4" />
                  </div>
                  <div>
                    <p class="text-sm font-medium">
                      {{ record.ip_address || '未知IP' }}
                    </p>
                    <p class="text-xs text-muted-foreground">
                      {{ record.user_agent ? getUserAgentInfo(record.user_agent) : '未知设备' }}
                    </p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="text-sm">
                    {{ formatDate(record.created_at) }}
                  </p>
                  <p class="text-xs text-muted-foreground">
                    {{ formatRelativeTime(record.created_at) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </div>

      <!-- 侧边栏 -->
      <div class="space-y-6">
        <!-- 账户信息 -->
        <Card>
          <div class="px-6">
            <h3 class="text-lg font-semibold mb-4">
              账户信息
            </h3>
            <div class="space-y-3">
              <div class="flex items-center space-x-2">
                <User class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm">{{ user?.name }}</span>
              </div>
              <div class="flex items-center space-x-2">
                <Mail class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm">{{ user?.email }}</span>
              </div>
              <div class="flex items-center space-x-2">
                <Shield class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm">{{ getRoleDisplayName(user?.role || '') }}</span>
              </div>
              <div class="flex items-center space-x-2">
                <Calendar class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm">{{ user?.created_at ? formatDate(user.created_at) : '未知' }}</span>
              </div>
            </div>
          </div>
        </Card>

        <!-- 安全设置 -->
        <Card>
          <div class="px-6">
            <h3 class="text-lg font-semibold mb-4 flex items-center">
              <Shield class=" h-5 w-5 text-blue-500" />
              安全设置
            </h3>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium">
                    密码强度
                  </p>
                  <p class="text-xs text-muted-foreground">
                    定期更新密码
                  </p>
                </div>
                <div class="w-16 h-2 bg-gray-200 rounded-full">
                  <div class="h-2 bg-green-500 rounded-full" style="width: 80%" />
                </div>
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium">
                    最后登录
                  </p>
                  <p class="text-xs text-muted-foreground">
                    {{ user?.last_login_at ? formatRelativeTime(user.last_login_at) : '从未登录' }}
                  </p>
                </div>
                <CheckCircle class="h-5 w-5 text-green-500" />
              </div>
            </div>
          </div>
        </Card>

        <!-- 安全提醒 -->
        <Card>
          <div class="px-6">
            <h3 class="text-lg font-semibold mb-4 flex items-center">
              <AlertTriangle class=" h-5 w-5 text-orange-500" />
              安全提醒
            </h3>
            <div class="space-y-2 text-sm text-muted-foreground">
              <p>• 定期更新密码，建议每3个月更换一次</p>
              <p>• 使用强密码，包含大小写字母、数字和特殊字符</p>
              <p>• 不要在公共设备上保存登录信息</p>
              <p>• 发现异常登录记录请及时联系管理员</p>
            </div>
          </div>
        </Card>
      </div>
    </div>

    <!-- 添加 Passkey 对话框 -->
    <AlertDialog v-model:open="showAddCredentialDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>添加新的 Passkey</AlertDialogTitle>
          <AlertDialogDescription>
            为这个 Passkey 设置一个名称，方便您识别不同的设备或浏览器。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div class="py-4">
          <Input
            v-model="newCredentialName"
            placeholder="例如：MacBook Pro 指纹"
            :disabled="addingCredential"
          />
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel :disabled="addingCredential">
            取消
          </AlertDialogCancel>
          <AlertDialogAction
            :disabled="!newCredentialName.trim() || addingCredential"
            @click="addCredential"
          >
            <Loader2 v-if="addingCredential" class="h-4 w-4 animate-spin" />
            {{ addingCredential ? '添加中...' : '继续' }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
