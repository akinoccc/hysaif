<script setup lang="ts">
import type { SecretItemHistory } from '@/api/types'
import { formatDate } from 'date-fns'
import {
  Eye,
  GitCompare,
  History,
  Loader2,
  RefreshCw,
  Undo2,
} from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { secretItemAPI } from '@/api/secret'
import { PermissionButton } from '@/components/common/permission'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Pagination } from '@/components/ui/pagination'
import RestoreConfirmDialog from './RestoreConfirmDialog.vue'
import VersionCompareDialog from './VersionCompareDialog.vue'
import VersionDetailDialog from './VersionDetailDialog.vue'

const props = defineProps<{
  itemId: string | number
}>()

const emit = defineEmits<{
  (e: 'restored'): void
}>()

// 响应式数据
const loading = ref(false)
const comparing = ref(false)
const restoring = ref(false)
const historyList = ref<SecretItemHistory[]>([])
const selectedVersions = ref<number[]>([])
const selectedHistory = ref<SecretItemHistory | null>(null)
const showCompareDialog = ref(false)
const showDetailDialog = ref(false)
const showRestoreDialog = ref(false)

// 分页状态
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
})

const sortedVersions = computed(() => {
  return [...selectedVersions.value].sort((a, b) => a - b)
})

// 加载历史记录
async function loadHistory() {
  try {
    loading.value = true
    const response = await secretItemAPI.getItemHistory(props.itemId, {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
    })

    historyList.value = response.data || []
    pagination.value.total = response.pagination?.total || 0
  }
  catch (error) {
    console.error('加载历史记录失败:', error)
    toast.error('加载历史记录失败')
  }
  finally {
    loading.value = false
  }
}

// 切换版本选择
function toggleVersionSelection(version: number) {
  const index = selectedVersions.value.indexOf(version)
  if (index > -1) {
    selectedVersions.value.splice(index, 1)
  }
  else if (selectedVersions.value.length < 2) {
    selectedVersions.value.push(version)
  }
}

// 对比版本
async function compareVersions() {
  if (selectedVersions.value.length !== 2) {
    toast.error('请选择两个版本进行对比')
    return
  }

  showCompareDialog.value = true
}

// 查看版本详情
function viewVersion(history: SecretItemHistory) {
  selectedHistory.value = history
  showDetailDialog.value = true
}

// 恢复版本
function restoreVersion(history: SecretItemHistory) {
  selectedHistory.value = history
  showRestoreDialog.value = true
}

// 确认恢复
async function confirmRestore(reason?: string) {
  if (!selectedHistory.value) {
    return
  }

  try {
    restoring.value = true
    await secretItemAPI.restoreItemFromHistory(props.itemId, {
      version: selectedHistory.value.version,
      reason: reason || `恢复到版本 ${selectedHistory.value.version}`,
    })

    toast.success('恢复成功')
    showRestoreDialog.value = false
    selectedHistory.value = null

    // 重新加载历史记录
    await loadHistory()

    // 通知父组件刷新数据
    emit('restored')
  }
  catch (error) {
    console.error('恢复失败:', error)
    toast.error('恢复失败')
  }
  finally {
    restoring.value = false
  }
}

// 监听 itemId 变化
watch(() => props.itemId, () => {
  if (props.itemId) {
    selectedVersions.value = []
    pagination.value.page = 1
    loadHistory()
  }
}, { immediate: true })

onMounted(() => {
  if (props.itemId) {
    loadHistory()
  }
})
</script>

<template>
  <div class="space-y-6">
    <!-- 历史记录标题 -->
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <!-- <History class="h-5 w-5 text-muted-foreground" /> -->
        <h3 class="text-lg font-semibold">
          历史版本
        </h3>
      </div>
      <div class="flex items-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="loading"
          @click="loadHistory"
        >
          <RefreshCw
            class="h-4 w-4"
            :class="{ 'animate-spin': loading }"
          />
          刷新
        </Button>
        <Button
          v-if="selectedVersions.length === 2"
          variant="outline"
          size="sm"
          :disabled="comparing"
          @click="compareVersions"
        >
          <GitCompare class="h-4 w-4" />
          对比版本
        </Button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-8">
      <Loader2 class="h-6 w-6 animate-spin" />
    </div>

    <!-- 历史记录列表 -->
    <div v-else-if="historyList.length > 0" class="space-y-4">
      <div class="text-sm text-muted-foreground">
        选择两个版本进行对比，或点击恢复按钮回溯到指定版本
      </div>

      <div class="space-y-2">
        <div
          v-for="(history, index) in historyList"
          :key="history.id"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
          :class="{
            'bg-muted/30': selectedVersions.includes(history.version),
            'border-primary': selectedVersions.includes(history.version),
          }"
        >
          <div class="flex items-center space-x-4">
            <Checkbox
              :model-value="selectedVersions.includes(history.version)"
              :disabled="selectedVersions.length >= 2 && !selectedVersions.includes(history.version)"
              @update:model-value="toggleVersionSelection(history.version)"
            />
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-2">
                <Badge variant="secondary" class="text-xs">
                  v{{ history.version }}
                </Badge>
                <span v-if="index === 0" class="text-xs text-green-600 font-medium">
                  (当前版本)
                </span>
              </div>
              <p class="text-sm text-muted-foreground mt-1">
                {{ history.change_reason || '无修改说明' }}
              </p>
              <div class="flex items-center space-x-4 mt-2 text-xs text-muted-foreground">
                <span>{{ formatDate(new Date(history.created_at), 'yyyy-MM-dd HH:mm:ss') }}</span>
                <span v-if="history.created_by">{{ history.created_by.name }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              title="查看此版本详情"
              @click="viewVersion(history)"
            >
              <Eye class="h-4 w-4" />
            </Button>
            <PermissionButton
              v-if="index > 0"
              :permission="{ resource: 'secret', action: 'update' }"
              variant="ghost"
              size="sm"
              :disabled="restoring"
              title="恢复到此版本"
              @click="restoreVersion(history)"
            >
              <Undo2 class="h-4 w-4" />
            </PermissionButton>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.total > pagination.pageSize" class="flex justify-center">
        <Pagination
          v-model:page="pagination.page"
          :items-per-page="pagination.pageSize"
          :total="pagination.total"
          @update:page="loadHistory"
        />
      </div>
    </div>

    <!-- 无数据状态 -->
    <div v-else class="text-center py-8">
      <History class="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
      <h3 class="text-lg font-medium mb-2">
        暂无历史记录
      </h3>
      <p class="text-muted-foreground">
        此密钥尚未有修改记录
      </p>
    </div>

    <!-- 版本对比对话框 -->
    <VersionCompareDialog
      v-if="showCompareDialog"
      v-model:open="showCompareDialog"
      :item-id="itemId"
      :version1="sortedVersions[0]"
      :version2="sortedVersions[1]"
      :comparing="comparing"
      @close="showCompareDialog = false"
    />

    <!-- 版本详情对话框 -->
    <VersionDetailDialog
      v-if="showDetailDialog"
      v-model:open="showDetailDialog"
      :history="selectedHistory"
      @close="showDetailDialog = false"
    />

    <!-- 恢复确认对话框 -->
    <RestoreConfirmDialog
      v-if="showRestoreDialog"
      v-model:open="showRestoreDialog"
      :history="selectedHistory"
      :restoring="restoring"
      @confirm="confirmRestore"
      @close="showRestoreDialog = false"
    />
  </div>
</template>
