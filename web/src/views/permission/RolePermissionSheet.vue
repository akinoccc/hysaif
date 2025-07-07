<script setup lang="ts">
import { Save, Shield, X } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Separator } from '@/components/ui/separator'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'

interface Permission {
  action: string
  label: string
  description: string
}

interface PermissionModule {
  id: string
  name: string
  description: string
  permissions: Permission[]
}

interface Role {
  value: string
  label: string
  description: string
  color: string
}

interface Props {
  open?: boolean
  role?: Role
  modules: PermissionModule[]
  initialPermissions: Record<string, string[]>
  defaultPermissions?: Record<string, string[]>
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'save', roleValue: string, permissions: Record<string, string[]>): void
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
  defaultPermissions: () => ({}),
})

const emit = defineEmits<Emits>()

const rolePermissions = ref<Record<string, string[]>>({})

function initializePermissions() {
  if (!props.role)
    return

  const permissions: Record<string, string[]> = {}
  // 使用传入的实际权限数据，如果没有则初始化为空
  props.modules.forEach((module) => {
    const moduleId = module.id
    const initialPerms = props.initialPermissions?.[moduleId] || []
    const defaultPerms = props.defaultPermissions?.[moduleId] || []

    // 合并初始权限和默认权限，确保默认权限始终存在
    const allPerms = new Set([...initialPerms, ...defaultPerms])
    permissions[moduleId] = Array.from(allPerms)
  })
  rolePermissions.value = permissions
}

// 监听角色和权限数据变化
watch(() => props.role, initializePermissions, { immediate: true })
watch(() => props.initialPermissions, initializePermissions, { immediate: true })
watch(() => props.defaultPermissions, initializePermissions, { immediate: true })

// 检查是否为超级管理员通配符权限
function isSuperAdminWildcard() {
  if (!props.role || props.role.value !== 'super_admin')
    return false
  // 检查是否有通配符权限 (*, *, *)
  return Object.values(props.initialPermissions).some(permissions =>
    permissions.includes('*'),
  )
}

// 检查权限是否为默认权限（不可删除）
function isDefaultPermission(moduleId: string, action: string) {
  const defaultPerms = props.defaultPermissions?.[moduleId] || []
  return defaultPerms.includes(action)
}

// 检查权限是否被选中
function isPermissionChecked(moduleId: string, action: string) {
  // 超级管理员通配符权限处理
  if (isSuperAdminWildcard()) {
    return true
  }

  return rolePermissions.value[moduleId]?.includes(action) || false
}

// 切换权限状态
function togglePermission(moduleId: string, action: string) {
  // 超级管理员通配符权限不允许修改
  if (isSuperAdminWildcard()) {
    return
  }

  // 默认权限不允许删除
  if (isDefaultPermission(moduleId, action)) {
    return
  }

  if (!rolePermissions.value[moduleId]) {
    rolePermissions.value[moduleId] = []
  }

  const permissions = rolePermissions.value[moduleId]
  const index = permissions.indexOf(action)

  if (index > -1) {
    permissions.splice(index, 1)
  }
  else {
    permissions.push(action)
  }
}

// 检查模块是否全选
function isModuleAllSelected(module: PermissionModule) {
  // 超级管理员通配符权限处理
  if (isSuperAdminWildcard()) {
    return true
  }

  const modulePermissions = rolePermissions.value[module.id] || []
  return module.permissions.every(p => modulePermissions.includes(p.action))
}

// 检查模块是否部分选中
function isModuleIndeterminate(module: PermissionModule) {
  // 超级管理员通配符权限处理
  if (isSuperAdminWildcard()) {
    return false // 全选状态下不显示部分选中
  }

  const modulePermissions = rolePermissions.value[module.id] || []
  const selectedCount = module.permissions.filter(p => modulePermissions.includes(p.action)).length
  return selectedCount > 0 && selectedCount < module.permissions.length
}

// 切换模块全选状态
function toggleModuleAll(module: PermissionModule) {
  // 超级管理员通配符权限不允许修改
  if (isSuperAdminWildcard()) {
    return
  }

  const defaultPerms = props.defaultPermissions?.[module.id] || []

  if (isModuleAllSelected(module)) {
    // 取消全选时保留默认权限
    rolePermissions.value[module.id] = [...defaultPerms]
  }
  else {
    // 全选时包含所有权限
    rolePermissions.value[module.id] = module.permissions.map(p => p.action)
  }
}

// 保存权限
function handleSave() {
  if (props.role) {
    emit('save', props.role.value, rolePermissions.value)
  }
  handleClose()
}

// 关闭弹窗
function handleClose() {
  emit('update:open', false)
}

// 计算标题
const sheetTitle = computed(() => {
  return props.role ? `编辑 ${props.role.label} 权限` : '编辑角色权限'
})

// 计算描述
const sheetDescription = computed(() => {
  return props.role ? `配置 ${props.role.label} 在各个模块中的访问权限` : '配置角色权限'
})
</script>

<template>
  <Sheet :open="open" @update:open="handleClose">
    <SheetContent class="sm:max-w-[600px] overflow-y-auto px-6 py-4 gap-6">
      <SheetHeader class="p-0">
        <div class="flex items-center space-x-3">
          <div v-if="role" class="p-2 rounded-lg text-white" :class="[role.color]">
            <Shield class="h-5 w-5" />
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

      <Separator />

      <!-- 超级管理员通配符权限提示 -->
      <div v-if="isSuperAdminWildcard()" class="p-4 bg-amber-50 border border-amber-200 rounded-lg">
        <div class="flex items-center space-x-2">
          <Shield class="h-4 w-4 text-amber-600" />
          <span class="text-sm font-medium text-amber-800">超级管理员通配符权限</span>
        </div>
        <p class="text-xs text-amber-700 mt-1">
          当前角色拥有通配符权限（*），自动拥有所有模块的全部权限，无需手动配置。
        </p>
      </div>

      <!-- 权限配置 -->
      <div class="space-y-6">
        <Card v-for="module in modules" :key="module.id" class="border-2">
          <CardHeader class="pb-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <CardTitle class="text-base">
                  {{ module.name }}
                </CardTitle>
                <Badge variant="outline" class="text-xs">
                  {{ isSuperAdminWildcard() ? module.permissions.length : (rolePermissions[module.id]?.length || 0) }}/{{ module.permissions.length }}
                </Badge>
              </div>
              <Checkbox
                :model-value="isModuleAllSelected(module)"
                :indeterminate="isModuleIndeterminate(module)"
                :disabled="isSuperAdminWildcard()"
                @update:model-value="toggleModuleAll(module)"
              />
            </div>
            <p class="text-sm text-muted-foreground">
              {{ module.description }}
            </p>
          </CardHeader>
          <CardContent class="pt-0">
            <div class="grid gap-3 sm:grid-cols-2">
              <div
                v-for="permission in module.permissions"
                :key="permission.action"
                class="flex items-start space-x-3 p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
                :class="{
                  'border-amber-200 bg-amber-50': isDefaultPermission(module.id, permission.action),
                }"
              >
                <Checkbox
                  :model-value="isPermissionChecked(module.id, permission.action)"
                  :disabled="isSuperAdminWildcard() || isDefaultPermission(module.id, permission.action)"
                  class="mt-0.5"
                  @update:model-value="togglePermission(module.id, permission.action)"
                />
                <div class="space-y-1 flex-1">
                  <div class="flex items-center space-x-2">
                    <Badge variant="secondary" class="text-xs">
                      {{ permission.action }}
                    </Badge>
                    <span class="text-sm font-medium">{{ permission.label }}</span>
                  </div>
                  <p class="text-xs text-muted-foreground leading-relaxed">
                    {{ permission.description }}
                  </p>
                  <Badge
                    v-if="isDefaultPermission(module.id, permission.action)"
                    variant="outline"
                    class="text-xs text-amber-700 border-amber-300 bg-amber-100"
                  >
                    默认权限
                  </Badge>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-end space-x-3 mt-8 pt-6 border-t">
        <Button variant="outline" @click="handleClose">
          <X class="h-4 w-4" />
          取消
        </Button>
        <PermissionButton
          :permission="{ resource: 'policy', action: 'update' }"
          @click="handleSave"
        >
          <Save class="h-4 w-4" />
          保存权限
        </PermissionButton>
      </div>
    </SheetContent>
  </Sheet>
</template>
