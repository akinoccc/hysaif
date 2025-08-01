<script setup lang="ts">
import type { SecretItem } from '@/api/types'
import { AlertCircle, Clock, Info, Loader2, RefreshCw } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { toast } from 'vue-sonner'
import { secretItemAPI } from '@/api/secret'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { getCategoryByKey } from '@/constants'
import { formatDate } from '@/utils/date'

defineProps<{
  secretDataTitle: string
  secretDataIcon?: any
  errorText?: string
}>()

const emit = defineEmits<{
  (e: 'loadSuccess', item: SecretItem): void
  (e: 'loadError', error: any): void
}>()

const route = useRoute()

const item = ref<SecretItem>()
const loading = ref(false)

// 方法
async function loadItem() {
  try {
    loading.value = true
    const itemId = route.params.id as string
    const useAccessCheck = route.query.access === 'true'

    // 根据URL参数决定使用哪个API
    const response = useAccessCheck
      ? await secretItemAPI.getItemWithAccess(itemId)
      : await secretItemAPI.getItem(itemId)

    item.value = response as SecretItem
    emit('loadSuccess', item.value)
  }
  catch (error) {
    console.error('加载数据失败:', error)
    emit('loadError', error)
  }
  finally {
    loading.value = false
  }
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    toast.success('复制成功')
    // 这里可以添加成功提示
  }
  catch (error) {
    console.error('复制失败:', error)
  }
}

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
  if (item.expires_at && item.expires_at < (Date.now()) + 30 * 24 * 60 * 60) {
    return '即将过期'
  }
  return '正常'
}

// 组件挂载时加载数据
onMounted(() => {
  loadItem()
})

// 暴露方法给父组件
defineExpose({
  item,
  loadItem,
  copyToClipboard,
})
</script>

<template>
  <div v-if="item" class="grid gap-6 lg:grid-cols-2">
    <!-- 基础信息 -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-1">
          <Info class=" h-5 w-5" />
          基础信息
        </CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4">
          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">名称</Label>
            <p class="text-lg font-medium">
              {{ item.name }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">描述</Label>
            <p class="text-sm">
              {{ item.description || '无' }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">分类</Label>
            <div class="flex items-center mt-1">
              <span
                class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary/10 text-primary"
              >
                {{ getCategoryByKey(item.category)?.label }}
              </span>
            </div>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">环境</Label>
            <p class="text-sm">
              {{ item.environment || '无' }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">标签</Label>
            <div class="flex flex-wrap gap-2">
              <template v-if="item.tags && item.tags.length > 0">
                <span
                  v-for="tag in item.tags" :key="tag"
                  class="inline-flex items-center px-3 py-1 rounded-md text-xs font-medium bg-muted text-muted-foreground"
                >
                  {{ tag }}
                </span>
              </template>
              <template v-else>
                <span class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-muted text-muted-foreground">
                  无
                </span>
              </template>
            </div>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">状态</Label>
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium mt-1" :class="[
                getStatusClass(item),
              ]"
            >
              {{ getStatusText(item) }}
            </span>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">过期时间</Label>
            <p class="text-sm">
              {{ item.expires_at ? formatDate(item.expires_at) : '无' }}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- 密钥信息 -->
    <Card>
      <CardHeader class="flex items-center justify-between">
        <CardTitle class="flex items-center gap-1">
          <component :is="secretDataIcon" class=" h-5 w-5" />
          {{ secretDataTitle }}
        </CardTitle>
        <CardAction>
          <slot name="secret-actions" />
        </CardAction>
      </CardHeader>
      <CardContent class="space-y-4">
        <slot
          name="secret-data"
          :item="item"
        />
      </CardContent>
    </Card>

    <!-- 系统信息 -->
    <Card class="lg:col-span-2">
      <CardHeader>
        <CardTitle class="flex items-center gap-1">
          <Clock class=" h-5 w-5" />
          系统信息
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid gap-4 md:grid-cols-2">
          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">创建时间</Label>
            <p class="text-sm mt-1">
              {{ formatDate(item.created_at) }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">更新时间</Label>
            <p class="text-sm mt-1">
              {{ formatDate(item.updated_at) }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">创建者</Label>
            <p class="text-sm mt-1">
              {{ item.creator.name }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-sm font-medium text-muted-foreground">更新者</Label>
            <p class="text-sm mt-1">
              {{ item.updater.name }}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>

  <!-- 加载状态 -->
  <div v-else-if="loading" class="flex justify-center py-12">
    <Loader2 class="h-8 w-8 animate-spin" />
  </div>

  <!-- 错误状态 -->
  <div v-else class="text-center py-12">
    <AlertCircle class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
    <h3 class="text-lg font-medium mb-2">
      加载失败
    </h3>
    <p class="text-muted-foreground mb-4">
      无法加载{{ errorText }}
    </p>
    <PermissionButton
      :permission="{ resource: 'secret', action: 'read' }"
      @click="loadItem"
    >
      <RefreshCw class=" h-4 w-4" />
      重试
    </PermissionButton>
  </div>
</template>
