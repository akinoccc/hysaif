import type { User } from '@/api/types'
import { Edit, Eye, Trash2 } from 'lucide-vue-next'
import { h } from 'vue'
import { PermissionButton } from '@/components'
import { Badge } from '@/components/ui/badge'

import { formatDateTime } from '@/utils/date'

export interface UserColumnActions {
  onView: (user: User) => void
  onEdit: (user: User) => void
  onDelete: (user: User) => void
}

// 角色标签样式映射
const roleVariants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  super_admin: 'destructive',
  sec_mgr: 'default',
  dev: 'secondary',
  auditor: 'outline',
  bot: 'outline',
}

// 角色显示名称映射
const roleLabels: Record<string, string> = {
  super_admin: '超级管理员',
  sec_mgr: '安全管理员',
  dev: '开发人员',
  auditor: '审计员',
  bot: '服务账号',
}

// 状态标签样式映射
const statusVariants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  active: 'default',
  disabled: 'secondary',
  locked: 'destructive',
  expired: 'outline',
}

// 状态显示名称映射
const statusLabels: Record<string, string> = {
  active: '活跃',
  disabled: '禁用',
  locked: '锁定',
  expired: '过期',
}

export function generateColumns(actions: UserColumnActions) {
  return [
    {
      accessorKey: 'name',
      header: '名称',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        return h('div', {}, user.name || '-')
      },
    },
    {
      accessorKey: 'email',
      header: '邮箱',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        return h('div', { class: 'text-muted-foreground' }, user.email || '-')
      },
    },
    {
      accessorKey: 'role',
      header: '角色',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        const variant = roleVariants[user.role || ''] || 'outline'
        const label = roleLabels[user.role || ''] || user.role || '-'
        return h(Badge, { variant }, () => label)
      },
    },
    {
      accessorKey: 'status',
      header: '状态',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        const variant = statusVariants[user.status || ''] || 'outline'
        const label = statusLabels[user.status || ''] || user.status || '-'
        return h(Badge, { variant }, () => label)
      },
    },
    {
      accessorKey: 'last_login_at',
      header: '最后登录',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        if (!user.last_login_at) {
          return h('div', { class: 'text-muted-foreground' }, '从未登录')
        }
        return h('div', { class: 'text-sm' }, formatDateTime(user.last_login_at))
      },
    },
    {
      accessorKey: 'created_at',
      header: '创建时间',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        return h('div', { class: 'text-sm text-muted-foreground' }, formatDateTime(user.created_at))
      },
    },
    {
      id: 'actions',
      header: '操作',
      cell: ({ row }: { row: { original: User } }) => {
        const user = row.original
        return h(
          'div',
          {},
          {
            default: () => [
              h(
                PermissionButton,
                {
                  permission: { resource: 'user', action: 'read' },
                  variant: 'ghost',
                  size: 'sm',
                  onClick: () => actions.onView(user),
                },
                {
                  default: () => [
                    h(Eye, { class: ' h-4 w-4' }),
                  ],
                },
              ),
              h(
                PermissionButton,
                {
                  permission: { resource: 'user', action: 'update' },
                  variant: 'ghost',
                  size: 'sm',
                  onClick: () => actions.onEdit(user),
                },
                {
                  default: () => [
                    h(Edit, { class: ' h-4 w-4' }),
                  ],
                },
              ),
              h(
                PermissionButton,
                {
                  permission: { resource: 'user', action: 'delete' },
                  variant: 'ghost',
                  size: 'sm',
                  class: 'text-destructive hover:text-destructive',
                  onClick: () => actions.onDelete(user),
                },
                {
                  default: () => [
                    h(Trash2, { class: ' h-4 w-4' }),
                  ],
                },
              ),
            ],
          },
        )
      },
    },
  ]
}
