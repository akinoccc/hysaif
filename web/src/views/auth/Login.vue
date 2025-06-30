<script setup lang="ts">
import { AlertCircle, Building2, Fingerprint, Loader2, Shield } from 'lucide-vue-next'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { webauthnAPI } from '@/api'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '@/stores/auth'
import { getCredential, isWebAuthnSupported } from '@/utils/webauthn'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const passkeyLoading = ref(false)
const error = ref('')
const webauthnSupported = ref(false)

onMounted(() => {
  webauthnSupported.value = isWebAuthnSupported()
})

async function handleLogin() {
  if (!form.email || !form.password) {
    error.value = '请输入用户名和密码'
    return
  }

  loading.value = true
  error.value = ''

  const result = await authStore.login(form.email, form.password)

  if (result.success) {
    router.push('/dashboard')
  }
  else {
    error.value = result.message || '登录失败'
  }

  loading.value = false
}

async function handlePasskeyLogin() {
  if (!form.email) {
    error.value = '请输入邮箱地址'
    return
  }

  passkeyLoading.value = true
  error.value = ''

  try {
    // 开始 WebAuthn 登录
    const options = await webauthnAPI.beginLogin(form.email)

    // 获取凭证
    const credential = await getCredential(options)

    // 完成登录
    const data = await webauthnAPI.finishLogin(form.email, credential)

    // 保存登录状态
    authStore.setAuth(data.token, data.user)

    router.push('/dashboard')
  }
  catch (err: any) {
    console.error('Passkey login error:', err)
    error.value = err.response?.data?.error || err.message || 'Passkey 登录失败'
  }
  finally {
    passkeyLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-background overflow-hidden">
    <div class="absolute inset-0">
      <!-- 几何图案背景 -->
      <div class="absolute inset-0 bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-900 dark:to-slate-700">
        <!-- 网格图案 -->
        <div
          class="absolute inset-0 opacity-[0.02] dark:opacity-[0.05]"
          style="background-image: linear-gradient(45deg, hsl(var(--foreground)) 1px, transparent 1px), linear-gradient(-45deg, hsl(var(--foreground)) 1px, transparent 1px); background-size: 60px 60px;"
        />
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="relative z-10 min-h-screen flex">
      <!-- 左侧品牌展示区 -->
      <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-slate-800 to-slate-900 dark:from-slate-900 dark:to-black items-center justify-center p-12 relative">
        <!-- 右侧斜角装饰 -->
        <div class="absolute top-0 right-0 w-32 h-full bg-gradient-to-l from-blue-400 to-transparent opacity-10 transform skew-x-12" />
        <div class="absolute top-0 right-4 w-24 h-full bg-gradient-to-l from-blue-500/20 to-transparent transform skew-x-12" />

        <!-- 几何装饰元素 -->
        <div class="absolute top-20 right-12 w-16 h-16 border border-blue-500/10 rounded-lg transform rotate-12" />
        <div class="absolute bottom-32 right-8 w-12 h-12 border border-blue-500/10 rounded-full" />
        <div class="absolute top-1/3 right-20 w-8 h-8 border-2 border-blue-500/10 transform rotate-45" />

        <div class="max-w-md text-center text-white relative z-10">
          <div class="mb-8 inline-flex h-20 w-20 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-blue-700 shadow-2xl shadow-blue-500/30">
            <Building2 class="h-12 w-12 text-white" />
          </div>
          <h1 class="mb-6 text-4xl font-bold leading-tight">
            企业级安全管理
          </h1>
          <p class="text-xl text-slate-300 dark:text-slate-400 leading-relaxed mb-8">
            保护您的关键业务数据，确保企业信息安全合规
          </p>

          <div class="space-y-6 text-left">
            <div class="flex items-center space-x-4 p-4 bg-white/5 dark:bg-white/10 rounded-xl backdrop-blur-sm border border-white/10">
              <div class="h-3 w-3 rounded-full bg-gradient-to-r from-destructive to-destructive shadow-lg shadow-destructive/50" />
              <span class="text-slate-200 dark:text-slate-300 font-medium">安全身份认证</span>
            </div>
            <div class="flex items-center space-x-4 p-4 bg-white/5 dark:bg-white/10 rounded-xl backdrop-blur-sm border border-white/10">
              <div class="h-3 w-3 rounded-full bg-gradient-to-r from-success to-success shadow-lg shadow-success/50" />
              <span class="text-slate-200 dark:text-slate-300 font-medium">权限精细管控</span>
            </div>
            <div class="flex items-center space-x-4 p-4 bg-white/5 dark:bg-white/10 rounded-xl backdrop-blur-sm border border-white/10">
              <div class="h-3 w-3 rounded-full bg-gradient-to-r from-info to-info shadow-lg shadow-info/50" />
              <span class="text-slate-200 dark:text-slate-300 font-medium">操作审计追踪</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录区域 -->
      <div class="flex-1 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative">
        <!-- 左侧装饰元素 -->
        <div class="hidden lg:block absolute left-0 top-0 w-24 h-full">
          <div class="absolute left-0 top-1/4 w-16 h-32 bg-gradient-to-r from-blue-500/10 to-transparent rounded-r-3xl" />
          <div class="absolute left-4 top-1/2 w-12 h-24 bg-gradient-to-r from-slate-400/10 to-transparent rounded-r-2xl" />
          <div class="absolute left-2 bottom-1/4 w-14 h-16 bg-gradient-to-r from-blue-600/15 to-transparent rounded-r-xl" />
        </div>

        <div class="w-full max-w-md relative z-10">
          <!-- 移动端顶部logo -->
          <div class="lg:hidden text-center mb-8">
            <div class="mx-auto h-16 w-16 flex items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-blue-700 shadow-lg shadow-blue-500/30 mb-4">
              <Building2 class="h-10 w-10 text-white" />
            </div>
            <h2 class="text-2xl font-bold text-foreground">
              企业敏感信息管理系统
            </h2>
          </div>

          <!-- 登录卡片 -->
          <Card class="bg-card/95 backdrop-blur-sm shadow-2xl border-0 relative overflow-hidden theme-transition">
            <!-- 卡片顶部装饰 -->
            <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary via-primary to-primary" />

            <CardHeader class="space-y-1 pb-6 relative">
              <CardTitle class="text-2xl font-semibold text-card-foreground text-center">
                安全登录
              </CardTitle>
              <CardDescription class="text-center text-muted-foreground">
                请使用您的企业账户登录系统
              </CardDescription>
            </CardHeader>

            <CardContent class="space-y-6">
              <Form @submit="handleLogin">
                <div class="space-y-4">
                  <FormField v-slot="{ componentField }" name="email">
                    <FormItem>
                      <FormLabel class="text-sm font-medium text-card-foreground">
                        企业邮箱
                      </FormLabel>
                      <FormControl>
                        <Input
                          v-bind="componentField"
                          v-model="form.email"
                          type="text"
                          required
                          class="h-11 border-input focus:border-ring focus:ring-ring/20 bg-background/80 theme-transition"
                          placeholder="请输入您的企业邮箱"
                        />
                      </FormControl>
                      <FormMessage class="text-destructive" />
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="password">
                    <FormItem>
                      <FormLabel class="text-sm font-medium text-card-foreground">
                        登录密码
                      </FormLabel>
                      <FormControl>
                        <Input
                          v-bind="componentField"
                          v-model="form.password"
                          type="password"
                          required
                          class="h-11 border-input focus:border-ring focus:ring-ring/20 bg-background/80 theme-transition"
                          placeholder="请输入您的登录密码"
                        />
                      </FormControl>
                      <FormMessage class="text-destructive" />
                    </FormItem>
                  </FormField>

                  <Alert v-if="error" variant="destructive" class="bg-destructive/10 border-destructive/20 text-destructive">
                    <AlertCircle class="h-4 w-4" />
                    <AlertDescription>
                      {{ error }}
                    </AlertDescription>
                  </Alert>

                  <Button
                    type="submit"
                    :disabled="loading"
                    class="w-full h-11 bg-gradient-to-r from-primary to-primary hover:from-primary/90 hover:to-primary/90 text-primary-foreground font-medium shadow-lg shadow-primary/25 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
                    {{ loading ? '登录中...' : '登录系统' }}
                  </Button>

                  <!-- Passkey 登录按钮 -->
                  <div v-if="webauthnSupported && form.email" class="relative">
                    <div class="absolute inset-0 flex items-center">
                      <span class="w-full border-t border-border" />
                    </div>
                    <div class="relative flex justify-center text-xs uppercase">
                      <span class="bg-background px-2 text-muted-foreground">或</span>
                    </div>
                  </div>

                  <Button
                    v-if="webauthnSupported && form.email"
                    type="button"
                    variant="outline"
                    :disabled="passkeyLoading"
                    class="w-full h-11"
                    @click="handlePasskeyLogin"
                  >
                    <Loader2 v-if="passkeyLoading" class="mr-2 h-4 w-4 animate-spin" />
                    <Fingerprint v-else class="mr-2 h-4 w-4" />
                    {{ passkeyLoading ? '验证中...' : '使用 Passkey 登录' }}
                  </Button>
                </div>
              </Form>

              <!-- 安全提示 -->
              <div class="mt-6 p-4 gradient-bg-primary rounded-lg border border-border/60">
                <div class="flex items-start space-x-3">
                  <Shield class="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                  <div class="text-sm text-card-foreground">
                    <p class="font-medium text-card-foreground mb-1">
                      安全提醒
                    </p>
                    <p class="text-muted-foreground">
                      请确保您在安全的网络环境中登录，不要在公共设备上保存登录信息。
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 底部信息 -->
          <div class="mt-8 text-center">
            <p class="text-sm text-muted-foreground">
              © 2024 Hysaif. 保留所有权利.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
