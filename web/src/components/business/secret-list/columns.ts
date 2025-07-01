import type { ColumnDef } from '@tanstack/vue-table'
import type { SecretItem } from '@/api/types'
import { Edit, Eye, Key, Trash2 } from 'lucide-vue-next'
import { h } from 'vue'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import { getCategoryByKey, SECRET_ITEM_TYPE } from '@/constants'
import { formatDate } from '@/lib/utils'
import { usePermissionStore } from '@/stores/permission'

function getStatusClass(item: SecretItem) {
  if (item.expires_at && item.expires_at < Date.now()) {
    return 'status-error'
  }
  if (item.expires_at && item.expires_at < (Date.now()) + 30 * 24 * 60 * 60) {
    return 'status-warning'
  }
  return 'status-active'
}

function getStatusText(item: SecretItem) {
  if (item.expires_at && item.expires_at < Date.now()) {
    return '已过期'
  }
  if (item.expires_at && item.expires_at < (Date.now() + 30 * 24 * 60 * 60)) {
    return '即将过期'
  }
  return '正常'
}

// 通用列定义生成函数
export function generateColumns(type: string): ColumnDef<SecretItem>[] {
  // 基础列定义（所有类型通用）
  const baseColumns: ColumnDef<SecretItem>[] = [
    // 名称列
    {
      accessorKey: 'name',
      header: '名称',
      cell: ({ row }) => {
        const item = row.original
        return h('div', { class: 'flex items-center space-x-3' }, [
          h('div', {}, [
            h('p', { class: 'font-medium' }, item.name),
            h('p', { class: 'text-sm text-muted-foreground' }, item.description),
          ]),
        ])
      },
    },
    // 分类列
    {
      accessorKey: 'category',
      header: '分类',
      cell: ({ row }) => {
        const category = row.getValue('category') as string
        const label = getCategoryByKey(category)?.label
        return h('span', {
          class: `
            inline-flex
            items-center
            px-2.5
            py-0.5
            rounded-full
            text-xs
            font-medium
            ${label ? 'bg-primary/10 text-primary' : 'bg-muted text-muted-foreground'}
          `,
        }, label || '未分类')
      },
    },
    // 状态列
    {
      id: 'status',
      header: '状态',
      cell: ({ row }) => {
        const item = row.original
        return h('span', {
          class: `inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusClass(item)}`,
        }, getStatusText(item))
      },
    },
    // 创建者列
    {
      accessorKey: 'creator',
      header: '创建者',
      cell: ({ row }) => {
        const item = row.original
        return h('span', { class: 'text-sm text-muted-foreground' }, item.creator?.name)
      },
    },
    // 更新者列
    {
      accessorKey: 'updater',
      header: '更新者',
      cell: ({ row }) => {
        const item = row.original
        return h('span', { class: 'text-sm text-muted-foreground' }, item.updater?.name)
      },
    },
    // 创建时间列
    {
      accessorKey: 'created_at',
      header: '创建时间',
      cell: ({ row }) => {
        const date = row.getValue('created_at') as string
        return h('span', { class: 'text-sm text-muted-foreground' }, formatDate(date))
      },
    },
    // 更新时间列
    {
      accessorKey: 'updated_at',
      header: '更新时间',
      cell: ({ row }) => {
        const date = row.getValue('updated_at') as string
        return h('span', { class: 'text-sm text-muted-foreground' }, formatDate(date))
      },
    },
    // 操作列
    {
      id: 'actions',
      header: '操作',
      cell: async ({ row }) => {
        const item = row.original
        const permissionStore = usePermissionStore()
        const hasDirectAccess = await permissionStore.hasPermission('secret', 'read')
        const hasApprovedAccess = item.has_approved_access // 后端需要提供此字段

        const actions = []

        if (hasDirectAccess) {
          // 有直接访问权限的用户显示查看按钮
          actions.push(
            h(PermissionButton, {
              permission: { resource: 'secret', action: 'read' },
              variant: 'ghost',
              size: 'sm',
              onClick: () => {
                window.dispatchEvent(new CustomEvent('view-item', { detail: item }))
              },
            }, () => h(Eye, { class: 'h-4 w-4' })),
          )
        }
        else if (hasApprovedAccess) {
          // 有已批准访问申请的用户显示查看按钮（使用getItemWithAccess）
          actions.push(
            h('button', {
              class: 'inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9',
              title: '查看密钥（通过访问申请）',
              onClick: () => {
                window.dispatchEvent(new CustomEvent('view-item-with-access', { detail: item }))
              },
            }, h(Eye, { class: 'h-4 w-4 text-green-600' })),
          )
        }
        else {
          // 没有直接访问权限且无已批准申请的用户显示申请访问按钮
          actions.push(
            h('button', {
              class: 'inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9',
              title: '申请访问权限',
              onClick: () => {
                window.dispatchEvent(new CustomEvent('request-access', { detail: item }))
              },
            }, h(Key, { class: 'h-4 w-4 text-blue-600' })),
          )
        }

        // 编辑按钮（需要权限）
        actions.push(
          h(PermissionButton, {
            permission: { resource: 'secret', action: 'update' },
            variant: 'ghost',
            size: 'sm',
            onClick: () => {
              window.dispatchEvent(new CustomEvent('edit-item', { detail: item }))
            },
          }, () => h(Edit, { class: 'h-4 w-4' })),
        )

        // 删除按钮（需要权限）
        actions.push(
          h(PermissionButton, {
            permission: { resource: 'secret', action: 'delete' },
            variant: 'ghost',
            size: 'sm',
            class: 'text-destructive hover:text-destructive',
            onClick: () => {
              window.dispatchEvent(new CustomEvent('delete-item', { detail: item }))
            },
          }, () => h(Trash2, { class: 'h-4 w-4' })),
        )

        return h('div', { class: 'flex items-center space-x-2' }, actions)
      },
    },
  ]

  // 根据不同类型添加特定列
  switch (type) {
    case SECRET_ITEM_TYPE.ApiKey:
      // API密钥特定列（如果有的话）
      return baseColumns

    case SECRET_ITEM_TYPE.Password:
      // 密码特定列
      return [
        ...baseColumns.slice(0, 2), // 名称和分类列
        ...baseColumns.slice(2), // 其余列
      ]

    case SECRET_ITEM_TYPE.SshKey:
      // SSH密钥特定列
      return [
        ...baseColumns.slice(0, 2), // 名称和分类列
        ...baseColumns.slice(2), // 其余列
      ]

    default:
      return baseColumns
  }
}

// 导出特定类型的列定义（为了向后兼容）
export const apiKeyColumns = generateColumns(SECRET_ITEM_TYPE.ApiKey)
export const passwordColumns = generateColumns(SECRET_ITEM_TYPE.Password)
export const sshKeyColumns = generateColumns(SECRET_ITEM_TYPE.SshKey)
export const tokenColumns = generateColumns(SECRET_ITEM_TYPE.Token)
export const accessKeyColumns = generateColumns(SECRET_ITEM_TYPE.AccessKey)
export const customColumns = generateColumns(SECRET_ITEM_TYPE.Custom)
