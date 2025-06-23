import type { ColumnDef } from '@tanstack/vue-table'
import type { AccessRequest } from '@/api/types'
import { Ban, Check, X } from 'lucide-vue-next'
import { h } from 'vue'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import { Badge } from '@/components/ui/badge'
import { formatDate } from '@/lib/utils'

function getStatusVariant(status: string) {
  switch (status) {
    case 'pending':
      return 'secondary'
    case 'approved':
      return 'default'
    case 'rejected':
      return 'destructive'
    case 'expired':
      return 'outline'
    case 'revoked':
      return 'outline'
    default:
      return 'secondary'
  }
}

function getStatusText(status: string) {
  switch (status) {
    case 'pending':
      return '待审批'
    case 'approved':
      return '已批准'
    case 'rejected':
      return '已拒绝'
    case 'expired':
      return '已过期'
    case 'revoked':
      return '已撤销'
    default:
      return status
  }
}

export const accessRequestColumns: ColumnDef<AccessRequest>[] = [
  {
    accessorKey: 'applicant',
    header: '申请人',
    cell: ({ row }) => {
      const request = row.original
      return h('div', {}, [
        h('p', { class: 'font-medium' }, request.applicant?.name || '未知用户'),
        h('p', { class: 'text-sm text-muted-foreground' }, request.applicant?.email || ''),
      ])
    },
  },
  {
    accessorKey: 'secret_item',
    header: '密钥项',
    cell: ({ row }) => {
      const request = row.original
      return h('div', {}, [
        h('p', { class: 'font-medium' }, request.secret_item?.name || '未知密钥'),
        h('p', { class: 'text-sm text-muted-foreground' }, request.secret_item?.type || ''),
      ])
    },
  },
  {
    accessorKey: 'reason',
    header: '申请理由',
    cell: ({ row }) => {
      const reason = row.getValue('reason') as string
      return h('p', {
        class: 'text-sm truncate max-w-xs',
        title: reason,
      }, reason || '无')
    },
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => {
      const status = row.getValue('status') as string
      return h(Badge, {
        variant: getStatusVariant(status),
      }, () => getStatusText(status))
    },
  },
  {
    accessorKey: 'note',
    header: '审批备注',
    cell: ({ row }) => {
      const note = row.getValue('note') as string
      return h('p', { class: 'text-sm text-muted-foreground' }, note || '-')
    },
  },
  {
    accessorKey: 'reject_reason',
    header: '拒绝理由',
    cell: ({ row }) => {
      const rejectReason = row.getValue('reject_reason') as string
      return h('p', { class: 'text-sm text-muted-foreground' }, rejectReason || '-')
    },
  },
  {
    accessorKey: 'created_at',
    header: '申请时间',
    cell: ({ row }) => {
      const date = row.getValue('created_at') as string
      return h('span', { class: 'text-sm text-muted-foreground' }, formatDate(date))
    },
  },
  {
    id: 'validity',
    header: '有效期',
    cell: ({ row }) => {
      const request = row.original
      if (request.status === 'approved' && request.valid_from && request.valid_until) {
        return h('div', { class: 'text-sm text-muted-foreground' }, [
          h('p', {}, formatDate(request.valid_from)),
          h('p', {}, `至 ${formatDate(request.valid_until)}`),
        ])
      }
      return h('span', { class: 'text-sm text-muted-foreground' }, '-')
    },
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => {
      const request = row.original
      const actions = []

      if (request.status === 'pending') {
        // 批准按钮
        actions.push(
          h(PermissionButton, {
            permission: { resource: 'access_request', action: 'approve' },
            size: 'sm',
            onClick: () => {
              window.dispatchEvent(new CustomEvent('approve-request', { detail: request }))
            },
          }, () => [
            h(Check, { class: 'h-3 w-3' }),
            '批准',
          ]),
        )

        // 拒绝按钮
        actions.push(
          h(PermissionButton, {
            permission: { resource: 'access_request', action: 'approve' },
            variant: 'outline',
            size: 'sm',
            onClick: () => {
              window.dispatchEvent(new CustomEvent('reject-request', { detail: request }))
            },
          }, () => [
            h(X, { class: 'h-3 w-3' }),
            '拒绝',
          ]),
        )
      }
      else if (request.status === 'approved') {
        // 作废按钮
        actions.push(
          h(PermissionButton, {
            permission: { resource: 'access_request', action: 'approve' },
            variant: 'destructive',
            size: 'sm',
            onClick: () => {
              window.dispatchEvent(new CustomEvent('revoke-request', { detail: request }))
            },
          }, () => [
            h(Ban, { class: 'h-3 w-3' }),
            '作废',
          ]),
        )
      }

      return h('div', { class: 'flex items-center space-x-2' }, actions)
    },
  },
]
