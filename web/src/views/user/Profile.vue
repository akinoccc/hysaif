<script setup lang="ts">
import {
  AlertTriangle,
  Calendar,
  CheckCircle,
  Eye,
  EyeOff,
  History,
  Key,
  Loader2,
  LogIn,
  Mail,
  RefreshCw,
  Save,
  Shield,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { userAPI } from '@/api'
import PasskeyManager from '@/components/business/passkey-manager/PasskeyManager.vue'
import { Button } from '@/components/ui/button'

import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { formatDate, formatRelativeTime, generatePassword, getRoleDisplayName, validatePasswordStrength } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const profileSubmitting = ref(false)
const passwordSubmitting = ref(false)
const loadingHistory = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

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
    // const response = await userAPI.getLoginHistory()
    // loginHistory.value = response.data || []
  }
  catch (error) {
    console.error('Failed to load login history:', error)
  }
  finally {
    loadingHistory.value = false
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
          <div class="p-6">
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
          <div class="p-6">
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
        <PasskeyManager />

        <!-- 登录历史 -->
        <Card>
          <div class="p-6">
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
              <span class="ml-2 text-sm">加载中...</span>
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
          <div class="p-6">
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
          <div class="p-6">
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
          <div class="p-6">
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
  </div>
</template>
