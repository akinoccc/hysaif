<script setup lang="ts">
import { BarChart3, Edit, FileText, Key, RefreshCw, Shield, Users } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import permissionAPI from '@/api/permission'
import { PermissionButton, PermissionWrapper } from '@/components'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import RolePermissionSheet from './RolePermissionSheet.vue'

// 权限模块信息
const permissionModules = [
  {
    id: 'user',
    name: '用户管理',
    description: '管理系统用户账户、角色分配等',
    icon: Users,
    color: 'bg-blue-500',
    permissions: [
      { action: 'create', label: '创建用户', description: '允许创建新的用户账户' },
      { action: 'read', label: '查看用户', description: '允许查看用户列表和详情' },
      { action: 'update', label: '更新用户', description: '允许修改用户信息和状态' },
      { action: 'delete', label: '删除用户', description: '允许删除用户账户' },
    ],
  },
  {
    id: 'secret',
    name: '敏感信息',
    description: '管理各类敏感信息和密钥',
    icon: Key,
    color: 'bg-green-500',
    permissions: [
      { action: 'create', label: '创建信息', description: '允许创建新的敏感信息' },
      { action: 'read', label: '查看信息', description: '允许查看敏感信息列表' },
      { action: 'update', label: '更新信息', description: '允许修改敏感信息内容' },
      { action: 'delete', label: '删除信息', description: '允许删除敏感信息' },
      { action: 'request', label: '请求访问', description: '允许请求访问敏感信息' },
      { action: 'temp', label: '临时访问', description: '允许临时访问敏感信息' },
    ],
  },
  {
    id: 'policy',
    name: '策略管理',
    description: '管理权限策略和访问控制',
    icon: Shield,
    color: 'bg-purple-500',
    permissions: [
      { action: 'create', label: '创建策略', description: '允许创建新的权限策略' },
      { action: 'read', label: '查看策略', description: '允许查看权限策略列表' },
      { action: 'update', label: '更新策略', description: '允许修改权限策略' },
      { action: 'delete', label: '删除策略', description: '允许删除权限策略' },
    ],
  },
  {
    id: 'audit',
    name: '审计日志',
    description: '查看系统操作审计记录',
    icon: FileText,
    color: 'bg-orange-500',
    permissions: [
      { action: 'read', label: '查看日志', description: '允许查看系统审计日志' },
    ],
  },
  {
    id: 'access_request',
    name: '访问申请',
    description: '管理敏感信息访问请求',
    icon: Shield,
    color: 'bg-purple-500',
    permissions: [
      { action: 'create', label: '创建申请', description: '允许创建新的访问申请' },
      { action: 'read', label: '查看申请', description: '允许查看访问申请列表' },
      { action: 'approve', label: '批准申请', description: '允许批准访问申请' },
      { action: 'reject', label: '拒绝申请', description: '允许拒绝访问申请' },
      { action: 'cancel', label: '取消申请', description: '允许取消访问申请' },
    ],
  },
  {
    id: 'dashboard',
    name: '仪表盘',
    description: '查看系统概览和统计信息',
    icon: BarChart3,
    color: 'bg-indigo-500',
    permissions: [
      { action: 'read', label: '查看仪表盘', description: '允许访问系统仪表盘' },
    ],
  },
]

// 角色信息
const roleInfo = [
  {
    value: 'super_admin',
    label: '超级管理员',
    description: '拥有系统所有权限，可以管理所有功能模块',
    color: 'bg-red-500',
  },
  {
    value: 'sec_mgr',
    label: '安全管理员',
    description: '负责安全策略管理和敏感信息监控',
    color: 'bg-yellow-500',
  },
  {
    value: 'dev',
    label: '开发人员',
    description: '可以访问开发相关的敏感信息',
    color: 'bg-blue-500',
  },
  {
    value: 'auditor',
    label: '审计员',
    description: '只能查看审计日志和系统状态',
    color: 'bg-gray-500',
  },
]

// 角色默认权限（不可删除）
const defaultRolePermissions: Record<string, Record<string, string[]>> = {
  super_admin: {
    // 超级管理员默认拥有所有权限，通过通配符实现
  },
  sec_mgr: {
    audit: ['read'],
    access_request: ['read', 'approve', 'reject', 'cancel'],
    policy: ['read', 'update'],
    secret: ['read', 'create', 'update', 'delete'],
    dashboard: ['read'],
  },
  dev: {
    access_request: ['create', 'read', 'cancel'],
    secret: ['request', 'read'],
    dashboard: ['read'],
  },
  auditor: {
    audit: ['read'],
    dashboard: ['read'],
  },
}

// 角色权限数据
const rolePermissionsData = ref<Record<string, Record<string, string[]>>>({})
const loading = ref(false)

// 角色权限编辑
const showPermissionSheet = ref(false)
const selectedRole = ref<typeof roleInfo[0] | undefined>()

// 获取角色权限
async function loadRolePermissions(role: string) {
  try {
    const response = await permissionAPI.getPermissionsForRole(role)
    const permissions = response.data?.permissions || []

    // 将权限数组转换为按模块分组的格式
    const groupedPermissions: Record<string, string[]> = {}

    permissions.forEach((permission: string[]) => {
      if (permission.length >= 3) {
        const [, resource, action] = permission
        if (!groupedPermissions[resource]) {
          groupedPermissions[resource] = []
        }
        groupedPermissions[resource].push(action)
      }
    })

    rolePermissionsData.value[role] = groupedPermissions
  }
  catch (error) {
    console.error(`获取角色 ${role} 权限失败:`, error)
    toast({
      title: '获取权限失败',
      description: `无法获取角色 ${role} 的权限信息`,
      variant: 'destructive',
    })
  }
}

// 加载所有角色权限
async function loadAllRolePermissions() {
  loading.value = true
  try {
    await Promise.all(roleInfo.map(role => loadRolePermissions(role.value)))
  }
  finally {
    loading.value = false
  }
}

// 打开权限编辑弹窗
async function openPermissionSheet(role: typeof roleInfo[0]) {
  selectedRole.value = role
  // 确保加载了该角色的权限数据
  if (!rolePermissionsData.value[role.value]) {
    await loadRolePermissions(role.value)
  }
  showPermissionSheet.value = true
}

// 保存权限
async function saveRolePermissions(roleValue: string, permissions: Record<string, string[]>) {
  try {
    await permissionAPI.updateRolePermissions(roleValue, permissions)

    // 更新本地数据
    rolePermissionsData.value[roleValue] = permissions

    toast({
      title: '权限更新成功',
      description: `角色 ${roleValue} 的权限已更新`,
    })
  }
  catch (error) {
    console.error('保存权限失败:', error)
    toast({
      title: '权限更新失败',
      description: '请稍后重试',
      variant: 'destructive',
    })
  }
}

// 检查是否为超级管理员通配符权限
function isSuperAdminWildcard(roleValue: string) {
  const rolePermissions = rolePermissionsData.value[roleValue]
  if (!rolePermissions)
    return false

  // 检查是否有通配符权限 (*, *, *)
  return Object.values(rolePermissions).some(permissions =>
    permissions.includes('*'),
  )
}

// 获取角色在模块中的权限数量
function getRolePermissionCount(roleValue: string, moduleId: string) {
  // 超级管理员通配符权限处理
  if (roleValue === 'super_admin' && isSuperAdminWildcard(roleValue)) {
    const module = permissionModules.find(m => m.id === moduleId)
    return module ? module.permissions.length : 0
  }

  const rolePermissions = rolePermissionsData.value[roleValue]
  if (!rolePermissions || !rolePermissions[moduleId]) {
    return 0
  }
  return rolePermissions[moduleId].length
}

// 检查角色是否有模块权限
function hasModulePermission(roleValue: string, moduleId: string) {
  // 超级管理员通配符权限处理
  if (roleValue === 'super_admin' && isSuperAdminWildcard(roleValue)) {
    return true
  }

  return getRolePermissionCount(roleValue, moduleId) > 0
}

onMounted(() => {
  loadAllRolePermissions()
})
</script>

<template>
  <div class="container py-6 space-y-6">
    <div class="flex items-center justify-between">
      <div class="space-y-2">
        <h1 class="text-3xl font-bold tracking-tight">
          权限管理
        </h1>
        <p class="text-muted-foreground">
          管理系统的权限策略、角色配置和访问控制
        </p>
      </div>
      <PermissionButton
        :permission="{ resource: 'permissions', action: 'read' }"
        variant="outline"
        :disabled="loading"
        @click="loadAllRolePermissions"
      >
        <RefreshCw class="h-4 w-4" :class="[{ 'animate-spin': loading }]" />
        刷新权限
      </PermissionButton>
    </div>

    <Tabs default-value="roles" class="space-y-6">
      <TabsList class="grid w-full grid-cols-2 bg-muted-foreground/10">
        <TabsTrigger value="roles">
          角色管理
        </TabsTrigger>
        <TabsTrigger value="modules">
          模块权限
        </TabsTrigger>
      </TabsList>

      <!-- 角色说明标签页 -->
      <TabsContent value="roles" class="space-y-6">
        <PermissionWrapper resource="policy" action="read">
          <div class="grid gap-4 md:grid-cols-2">
            <Card v-for="role in roleInfo" :key="role.value" class="relative overflow-hidden">
              <CardHeader class="pb-3">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-3">
                    <div class="p-2 rounded-lg text-white" :class="[role.color]">
                      <Shield class="h-5 w-5" />
                    </div>
                    <div>
                      <CardTitle class="text-lg">
                        {{ role.label }}
                      </CardTitle>
                      <Badge variant="secondary" class="text-xs">
                        {{ role.value }}
                      </Badge>
                    </div>
                  </div>
                  <PermissionButton
                    :permission="{ resource: 'policy', action: 'update' }"
                    variant="outline"
                    size="sm"
                    @click="openPermissionSheet(role)"
                  >
                    <Edit class="h-4 w-4 mr-1" />
                    编辑权限
                  </PermissionButton>
                </div>
              </CardHeader>
              <CardContent>
                <p class="text-sm text-muted-foreground">
                  {{ role.description }}
                </p>
              </CardContent>
            </Card>
          </div>
          <template #fallback>
            <Card>
              <CardContent class="flex items-center justify-center py-12">
                <div class="text-center space-y-2">
                  <Shield class="h-12 w-12 mx-auto text-muted-foreground" />
                  <h3 class="text-lg font-medium">
                    权限不足
                  </h3>
                  <p class="text-muted-foreground">
                    您没有权限查看角色管理功能
                  </p>
                </div>
              </CardContent>
            </Card>
          </template>
        </PermissionWrapper>

        <!-- 权限矩阵说明 -->
        <Card>
          <CardHeader>
            <CardTitle>权限矩阵说明</CardTitle>
            <CardDescription>
              不同角色在各个模块中的权限分配情况
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div class="overflow-x-auto">
              <table class="w-full border-collapse">
                <thead>
                  <tr class="border-b">
                    <th class="text-left p-2 font-medium">
                      角色 \ 模块
                    </th>
                    <th v-for="module in permissionModules" :key="module.id" class="text-center p-2 font-medium">
                      {{ module.name }}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="role in roleInfo" :key="role.value" class="border-b">
                    <td class="p-2">
                      <div class="flex items-center space-x-2">
                        <div class="w-3 h-3 rounded-full" :class="[role.color]" />
                        <span class="font-medium">{{ role.label }}</span>
                      </div>
                    </td>
                    <td v-for="module in permissionModules" :key="module.id" class="text-center p-2">
                      <div v-if="hasModulePermission(role.value, module.id)" class="flex flex-col items-center space-y-1">
                        <Badge
                          :variant="getRolePermissionCount(role.value, module.id) === module.permissions.length ? 'default' : 'secondary'"
                          class="text-xs"
                        >
                          {{ getRolePermissionCount(role.value, module.id) === module.permissions.length ? '全部' : '部分' }}
                        </Badge>
                        <span class="text-xs text-muted-foreground">
                          {{ getRolePermissionCount(role.value, module.id) }}/{{ module.permissions.length }}
                        </span>
                      </div>
                      <span v-else class="text-muted-foreground text-xs">-</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- 模块权限标签页 -->
      <TabsContent value="modules" class="space-y-6">
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card v-for="module in permissionModules" :key="module.id" class="relative overflow-hidden">
            <CardHeader class="pb-3">
              <div class="flex items-center space-x-3">
                <div class="p-2 rounded-lg text-white" :class="[module.color]">
                  <component :is="module.icon" class="h-5 w-5" />
                </div>
                <div>
                  <CardTitle class="text-lg">
                    {{ module.name }}
                  </CardTitle>
                  <CardDescription class="text-sm">
                    {{ module.description }}
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent class="space-y-3">
              <div class="space-y-2">
                <h4 class="text-sm font-medium text-muted-foreground">
                  可用权限
                </h4>
                <div class="space-y-2">
                  <div
                    v-for="permission in module.permissions"
                    :key="permission.action"
                    class="flex items-start justify-between p-2 rounded-lg bg-muted/50"
                  >
                    <div class="space-y-1">
                      <div class="flex items-center space-x-2">
                        <Badge variant="outline" class="text-xs">
                          {{ permission.action }}
                        </Badge>
                        <span class="text-sm font-medium">{{ permission.label }}</span>
                      </div>
                      <p class="text-xs text-muted-foreground">
                        {{ permission.description }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </TabsContent>
    </Tabs>
  </div>
  <!-- 角色权限编辑弹窗 -->
  <RolePermissionSheet
    v-model:open="showPermissionSheet"
    :role="selectedRole"
    :modules="permissionModules"
    :initial-permissions="selectedRole ? rolePermissionsData[selectedRole.value] : {}"
    :default-permissions="selectedRole ? defaultRolePermissions[selectedRole.value] : {}"
    @save="saveRolePermissions"
  />
</template>
