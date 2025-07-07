<script setup lang="ts">
import type { CreateUserRequest, UpdateUserRequest, User } from '@/api/types'
import { toTypedSchema } from '@vee-validate/zod'
import {
  Activity,
  Calendar,
  Clock,
  Eye,
  EyeOff,
  Key,
  Mail,
  MapPin,
  Save,
  Shield,
  UserCheck,
  User as UserIcon,
  X,
} from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import * as z from 'zod'
import { userAPI } from '@/api/user'
import { PermissionButton } from '@/components'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { formatDateTime } from '@/utils/date'

const props = defineProps<{
  open: boolean
  mode: 'create' | 'edit' | 'view'
  user?: User | null
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'saved': []
}>()

// 角色选项
const roleOptions = [
  { value: 'super_admin', label: '超级管理员' },
  { value: 'sec_mgr', label: '安全管理员' },
  { value: 'dev', label: '开发人员' },
  { value: 'auditor', label: '审计员' },
  { value: 'bot', label: '服务账号' },
]

// 状态选项
const statusOptions = [
  { value: 'active', label: '活跃' },
  { value: 'disabled', label: '禁用' },
  { value: 'locked', label: '锁定' },
  { value: 'expired', label: '过期' },
]

// 角色显示名称映射
const roleLabels: Record<string, string> = {
  super_admin: '超级管理员',
  sec_mgr: '安全管理员',
  dev: '开发人员',
  auditor: '审计员',
  bot: '服务账号',
}

// 状态显示名称映射
const statusLabels: Record<string, string> = {
  active: '活跃',
  disabled: '禁用',
  locked: '锁定',
  expired: '过期',
}

const loading = ref(false)
const showPassword = ref(false)

const isReadonly = computed(() => props.mode === 'view')
const isCreate = computed(() => props.mode === 'create')
const isEdit = computed(() => props.mode === 'edit')

const sheetTitle = computed(() => {
  switch (props.mode) {
    case 'create':
      return '新建用户'
    case 'edit':
      return '编辑用户'
    case 'view':
      return '用户详情'
    default:
      return '用户信息'
  }
})

const sheetDescription = computed(() => {
  switch (props.mode) {
    case 'create':
      return '创建一个新的用户账户'
    case 'edit':
      return '修改用户信息和权限设置'
    case 'view':
      return '查看用户详细信息'
    default:
      return ''
  }
})

// 表单验证模式
const createSchema = toTypedSchema(z.object({
  name: z.string().min(1, '姓名不能为空'),
  email: z.string().email('请输入有效的邮箱地址'),
  password: z.string().min(8, '密码至少8位'),
  status: z.string().min(1, '请选择状态'),
  role: z.string().min(1, '请选择角色'),
  permissions: z.array(z.string()).optional(),
}))

const editSchema = toTypedSchema(z.object({
  name: z.string().min(1, '姓名不能为空'),
  email: z.string().email('请输入有效的邮箱地址'),
  role: z.string().min(1, '请选择角色'),
  status: z.string().min(1, '请选择状态'),
  password: z.string().optional(),
  permissions: z.array(z.string()).optional(),
}))

const formSchema = computed(() => {
  return isCreate.value ? createSchema : editSchema
})

// 表单
const form = useForm({
  validationSchema: formSchema.value,
})

// 重置表单
function resetForm() {
  if (isCreate.value) {
    form.resetForm({
      values: {
        name: '',
        email: '',
        password: '',
        role: 'dev',
        status: 'active',
        permissions: [],
      },
    })
  }
  else if (props.user) {
    form.resetForm({
      values: {
        name: props.user.name || '',
        email: props.user.email || '',
        role: props.user.role || 'dev',
        status: props.user.status || 'active',
        password: '',
        permissions: props.user.permissions || [],
      },
    })
  }
}

// 提交表单
const onSubmit = form.handleSubmit(async (values) => {
  try {
    loading.value = true

    if (isCreate.value) {
      const createData: CreateUserRequest = {
        name: values.name,
        email: values.email,
        password: values.password!,
        role: values.role,
        permissions: values.permissions,
      }
      await userAPI.createUser(createData)
    }
    else if (isEdit.value && props.user) {
      const updateData: UpdateUserRequest = {
        name: values.name,
        email: values.email,
        role: values.role,
        status: values.status || 'active',
        permissions: values.permissions,
      }

      // 只有在密码不为空时才包含密码字段
      if (values.password && values.password.trim()) {
        updateData.password = values.password
      }

      await userAPI.updateUser(props.user.id, updateData)
    }

    emit('saved')
  }
  catch (error) {
    console.error('保存用户失败:', error)
  }
  finally {
    loading.value = false
  }
})

// 关闭Sheet
function handleClose() {
  emit('update:open', false)
}

// 监听props变化
watch(
  () => [props.open, props.user, props.mode],
  () => {
    if (props.open) {
      resetForm()
      showPassword.value = false
    }
  },
  { immediate: true },
)
</script>

<template>
  <Sheet :open="open" @update:open="handleClose">
    <SheetContent class="sm:max-w-[640px] overflow-y-auto px-6 py-4">
      <SheetHeader class="p-0">
        <div class="flex items-center gap-3">
          <div class="p-2 bg-primary/10 rounded-lg">
            <UserIcon class="h-8 w-8 text-primary" />
          </div>
          <div>
            <SheetTitle class="text-lg">
              {{ sheetTitle }}
            </SheetTitle>
            <SheetDescription class="text-sm text-muted-foreground mt-1">
              {{ sheetDescription }}
            </SheetDescription>
          </div>
        </div>
      </SheetHeader>
      <div class="space-y-6">
        <!-- 查看模式 -->
        <div v-if="isReadonly && user" class="space-y-6">
          <!-- 用户概览卡片 -->
          <Card class="border-0 shadow-sm bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/20 dark:to-indigo-950/20">
            <CardContent class="px-6 py-0">
              <div class="flex items-start gap-4">
                <div class="p-3 bg-white dark:bg-gray-800 rounded-full shadow-sm">
                  <UserIcon class="h-8 w-8 text-primary" />
                </div>
                <div class="flex-1 space-y-3">
                  <div>
                    <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100">
                      {{ user.name || '未命名用户' }}
                    </h2>
                    <p class="text-sm text-muted-foreground flex items-center gap-1 mt-1">
                      <Mail class="h-3 w-3" />
                      {{ user.email || '无邮箱' }}
                    </p>
                  </div>
                  <div class="flex items-center gap-3">
                    <Badge
                      :variant="user.status === 'active' ? 'default' : user.status === 'disabled' ? 'destructive' : 'secondary'"
                      class="flex items-center gap-1"
                    >
                      <Activity class="h-3 w-3" />
                      {{ statusLabels[user.status || ''] || user.status || '未知' }}
                    </Badge>
                    <Badge variant="outline" class="flex items-center gap-1">
                      <Shield class="h-3 w-3" />
                      {{ roleLabels[user.role || ''] || user.role || '未知角色' }}
                    </Badge>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 详细信息网格 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 登录信息 -->
            <Card class="gap-2">
              <CardHeader>
                <CardTitle class="text-base flex items-center gap-2">
                  <Clock class="h-4 w-4 text-blue-500" />
                  登录信息
                </CardTitle>
              </CardHeader>
              <CardContent class="space-y-4">
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">最后登录时间</Label>
                  <div class="flex items-center gap-2 p-3 bg-muted/50 rounded-lg">
                    <Calendar class="h-4 w-4 text-muted-foreground" />
                    <span class="text-sm font-medium">
                      {{ user.last_login_at ? formatDateTime(user.last_login_at) : '从未登录' }}
                    </span>
                  </div>
                </div>
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">最后登录IP</Label>
                  <div class="flex items-center gap-2 p-3 bg-muted/50 rounded-lg">
                    <MapPin class="h-4 w-4 text-muted-foreground" />
                    <code class="text-sm font-mono">
                      {{ user.last_login_ip || '无记录' }}
                    </code>
                  </div>
                </div>
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">登录失败次数</Label>
                  <div class="flex items-center gap-2 p-3 bg-muted/50 rounded-lg">
                    <Key class="h-4 w-4 text-muted-foreground" />
                    <span class="text-sm font-medium">
                      {{ user.failed_attempts || 0 }} 次
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- 权限信息 -->
            <Card class="gap-2">
              <CardHeader class="pb-3">
                <CardTitle class="text-base flex items-center gap-2">
                  <Shield class="h-4 w-4 text-green-500" />
                  权限信息
                </CardTitle>
              </CardHeader>
              <CardContent class="space-y-4">
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">特殊权限</Label>
                  <div class="p-3 bg-muted/50 rounded-lg min-h-[60px] flex items-center">
                    <div v-if="user.permissions && user.permissions.length > 0" class="flex flex-wrap gap-2">
                      <Badge v-for="permission in user.permissions" :key="permission" variant="outline" class="text-xs">
                        {{ permission }}
                      </Badge>
                    </div>
                    <p v-else class="text-sm text-muted-foreground italic">
                      无特殊权限
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- 元数据 -->
          <Card class="gap-2">
            <CardHeader class="pb-3">
              <CardTitle class="text-base flex items-center gap-2">
                <UserCheck class="h-4 w-4 text-purple-500" />
                元数据
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">创建时间</Label>
                  <p class="text-sm p-3 bg-muted/50 rounded-lg">
                    {{ formatDateTime(user.created_at) }}
                  </p>
                </div>
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">更新时间</Label>
                  <p class="text-sm p-3 bg-muted/50 rounded-lg">
                    {{ formatDateTime(user.updated_at) }}
                  </p>
                </div>
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">创建者</Label>
                  <p class="text-sm p-3 bg-muted/50 rounded-lg">
                    {{ user.creator?.name || '系统' }}
                  </p>
                </div>
                <div class="space-y-2">
                  <Label class="text-xs font-medium text-muted-foreground uppercase tracking-wide">更新者</Label>
                  <p class="text-sm p-3 bg-muted/50 rounded-lg">
                    {{ user.updater?.name || '无' }}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- 编辑/创建模式 -->
        <form v-else class="space-y-6" @submit="onSubmit">
          <!-- 基本信息卡片 -->
          <Card class="gap-4">
            <CardHeader>
              <CardTitle class="text-base flex items-center gap-2">
                <UserIcon class="h-4 w-4 text-blue-500" />
                基本信息
              </CardTitle>
            </CardHeader>
            <CardContent class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 名称 -->
                <FormField v-slot="{ componentField }" name="name">
                  <FormItem>
                    <FormLabel class="flex items-center gap-2">
                      <UserIcon class="h-3 w-3" />
                      名称 <span class="text-red-500">*</span>
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        placeholder="输入用户名称"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                </FormField>

                <!-- 邮箱 -->
                <FormField v-slot="{ componentField }" name="email">
                  <FormItem>
                    <FormLabel class="flex items-center gap-2">
                      <Mail class="h-3 w-3" />
                      邮箱 <span class="text-red-500">*</span>
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        type="email"
                        placeholder="输入邮箱地址"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                </FormField>
              </div>

              <!-- 密码 -->
              <FormField v-slot="{ componentField }" name="password">
                <FormItem>
                  <FormLabel class="flex items-center gap-2">
                    <Key class="h-3 w-3" />
                    密码
                    <span v-if="isCreate" class="text-red-500">*</span>
                    <span v-else class="text-muted-foreground text-xs">(留空则不修改)</span>
                  </FormLabel>
                  <FormControl>
                    <div class="relative">
                      <Input
                        v-bind="componentField"
                        :type="showPassword ? 'text' : 'password'"
                        :placeholder="isCreate ? '输入密码' : '输入新密码'"
                        class="pr-10"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                        @click="showPassword = !showPassword"
                      >
                        <Eye v-if="!showPassword" class="h-4 w-4" />
                        <EyeOff v-else class="h-4 w-4" />
                      </Button>
                    </div>
                  </FormControl>
                  <FormDescription v-if="isCreate" class="text-xs">
                    密码至少8位，建议包含大小写字母、数字和特殊字符
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              </FormField>
            </CardContent>
          </Card>

          <!-- 权限设置卡片 -->
          <Card class="gap-4">
            <CardHeader>
              <CardTitle class="text-base flex items-center gap-2">
                <Shield class="h-4 w-4 text-green-500" />
                权限设置
              </CardTitle>
            </CardHeader>
            <CardContent class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 角色 -->
                <FormField v-slot="{ componentField }" name="role">
                  <FormItem>
                    <FormLabel class="flex items-center gap-2">
                      <Shield class="h-3 w-3" />
                      角色 <span class="text-red-500">*</span>
                    </FormLabel>
                    <Select v-bind="componentField">
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="选择角色" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem v-for="option in roleOptions" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FormDescription class="text-xs">
                      角色决定用户的基础权限范围
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                </FormField>

                <!-- 状态 -->
                <FormField v-slot="{ componentField }" name="status">
                  <FormItem>
                    <FormLabel class="flex items-center gap-2">
                      <Activity class="h-3 w-3" />
                      状态 <span class="text-red-500">*</span>
                    </FormLabel>
                    <Select v-bind="componentField">
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="选择状态" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem v-for="option in statusOptions" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FormDescription class="text-xs">
                      状态决定用户的可用性
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                </FormField>
              </div>
            </CardContent>
          </Card>

          <!-- 操作按钮 -->
          <div class="flex justify-end gap-3 pt-6">
            <Button type="button" variant="outline" size="lg" @click="handleClose">
              <X class="h-4 w-4" />
              取消
            </Button>
            <PermissionButton
              :permission="{ resource: 'user', action: isCreate ? 'create' : 'update' }"
              type="submit"
              size="lg"
              :disabled="loading"
              class="min-w-[100px]"
            >
              <Save class="h-4 w-4" />
              {{ loading ? '保存中...' : '保存' }}
            </PermissionButton>
          </div>
        </form>
      </div>
    </SheetContent>
  </Sheet>
</template>
