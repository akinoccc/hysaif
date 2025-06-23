<script setup lang="ts" generic="TData, TValue">
import type { ColumnDef } from '@tanstack/vue-table'
import type { Component } from 'vue'
import {
  FlexRender,
  getCoreRowModel,
  getPaginationRowModel,
  useVueTable,
} from '@tanstack/vue-table'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { computed } from 'vue'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

interface DataTableProps {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  loading?: boolean
  total?: number
  pageSize?: number
  totalPages?: number
  pageCount?: number
  emptyIcon?: Component
  emptyTitle?: string
  emptyDescription?: string
  showPagination?: boolean
}

const props = withDefaults(defineProps<DataTableProps>(), {
  loading: false,
  pageSize: 10,
  showPagination: true,
  emptyTitle: '暂无数据',
  emptyDescription: '没有找到相关数据',
})

const emit = defineEmits<{
  pageChange: [page: number]
  pageSizeChange: [size: number]
  viewItem: [item: TData]
  viewItemWithAccess: [item: TData]
  editItem: [item: TData]
  deleteItem: [item: TData]
  requestAccess: [item: TData]
}>()

const currentPage = defineModel('currentPage', { default: 1 })

// 监听自定义事件（用于secret-list的兼容性）
if (typeof window !== 'undefined') {
  window.addEventListener('view-item', ((e: CustomEvent) => {
    emit('viewItem', e.detail)
  }) as EventListener)

  window.addEventListener('edit-item', ((e: CustomEvent) => {
    emit('editItem', e.detail)
  }) as EventListener)

  window.addEventListener('delete-item', ((e: CustomEvent) => {
    emit('deleteItem', e.detail)
  }) as EventListener)

  window.addEventListener('request-access', ((e: CustomEvent) => {
    emit('requestAccess', e.detail)
  }) as EventListener)

  window.addEventListener('view-item-with-access', ((e: CustomEvent) => {
    emit('viewItemWithAccess', e.detail)
  }) as EventListener)
}

// 计算总页数
const computedPageCount = computed(() => {
  if (props.pageCount)
    return props.pageCount
  if (props.totalPages)
    return props.totalPages
  if (props.total && props.pageSize)
    return Math.ceil(props.total / props.pageSize)
  return 0
})

const table = useVueTable({
  get data() { return props.data },
  get columns() { return props.columns },
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  manualPagination: true,
  pageCount: computedPageCount.value,
  state: {
    pagination: computed(() => ({
      pageIndex: (currentPage.value || 1) - 1, // TanStack Table 是0索引的
      pageSize: props.pageSize || 10,
    })).value,
  },
})

// 页码变化处理
function handlePageChange(page: number) {
  emit('pageChange', page)
}
</script>

<template>
  <div>
    <div class="border rounded-md">
      <Table>
        <TableHeader>
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead v-for="header in headerGroup.headers" :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <!-- 加载状态 -->
          <template v-if="loading">
            <TableRow>
              <TableCell :colspan="columns.length" class="h-24 text-center">
                <div class="flex items-center justify-center space-x-2">
                  <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-primary" />
                  <span class="text-muted-foreground">加载中...</span>
                </div>
              </TableCell>
            </TableRow>
          </template>
          <!-- 有数据时显示 -->
          <template v-else-if="table.getRowModel().rows?.length">
            <TableRow v-for="row in table.getRowModel().rows" :key="row.id" class="hover:bg-muted/50 transition-colors">
              <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>
          <!-- 空状态 -->
          <template v-else>
            <TableRow>
              <TableCell :colspan="columns.length" class="h-32">
                <!-- 自定义空状态 -->
                <div v-if="emptyIcon || emptyTitle !== '暂无数据' || emptyDescription !== '没有找到相关数据'" class="flex flex-col items-center justify-center space-y-3 text-center">
                  <component :is="emptyIcon" v-if="emptyIcon" class="h-12 w-12 text-muted-foreground/50" />
                  <div class="space-y-1">
                    <h3 class="text-lg font-medium text-muted-foreground">
                      {{ emptyTitle }}
                    </h3>
                    <p class="text-sm text-muted-foreground/80">
                      {{ emptyDescription }}
                    </p>
                  </div>
                </div>
                <!-- 默认空状态 -->
                <div v-else class="text-center">
                  <TableEmpty />
                </div>
              </TableCell>
            </TableRow>
          </template>
        </TableBody>
      </Table>
    </div>

    <!-- 分页 -->
    <div v-if="showPagination && computedPageCount > 1" class="flex justify-center mt-4">
      <Pagination
        v-model:page="currentPage"
        :total="total || (computedPageCount * pageSize)"
        :items-per-page="pageSize"
        :sibling-count="1"
        :show-edges="true"
        @update:page="handlePageChange"
      >
        <PaginationContent>
          <PaginationPrevious>
            <ChevronLeft />
            <span>上一页</span>
          </PaginationPrevious>
          <template v-for="(item, index) in computedPageCount" :key="index">
            <PaginationItem
              v-if="item <= computedPageCount"
              :value="item"
              :is-active="item === currentPage"
            >
              {{ item }}
            </PaginationItem>
          </template>
          <PaginationNext>
            <span>下一页</span>
            <ChevronRight />
          </PaginationNext>
        </PaginationContent>
      </Pagination>
    </div>
  </div>
</template>
