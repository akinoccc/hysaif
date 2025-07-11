<script setup lang="ts">
import {
  ArrowLeft,
  Edit,
  History,
} from 'lucide-vue-next'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SecretHistory } from '@/components/business/secret-history'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import DetailInfo from './DetailInfo.vue'

const props = defineProps<{
  title: string
  description: string
  secretType: string
  secretDataTitle: string
  secretDataIcon?: any
  errorText?: string
  editPath?: string
}>()

const router = useRouter()
const route = useRoute()

const infoRef = ref<InstanceType<typeof DetailInfo>>()

// 处理历史记录恢复
function handleHistoryRestore() {
  // 重新加载当前项目数据
  infoRef.value?.loadItem()
}

function onEdit() {
  const path = props.editPath || `/${props.secretType}/${route.params.id}/edit`
  router.push(path)
}

// 暴露方法给父组件
defineExpose({
  item: infoRef.value?.item,
  loadItem: infoRef.value?.loadItem,
  copyToClipboard: infoRef.value?.copyToClipboard,
})
</script>

<template>
  <div class="space-y-6">
    <!-- 页面标题和操作 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">
          {{ title }}
        </h1>
        <p class="text-muted-foreground">
          {{ description }}
        </p>
      </div>
      <div class="flex items-center space-x-2">
        <!-- 自定义头部操作插槽 -->
        <slot
          name="header-actions"
          :item="infoRef?.item"
        />
        <Button variant="outline" @click="router.back()">
          <ArrowLeft class=" h-4 w-4" />
          返回
        </Button>
        <PermissionButton
          :permission="{ resource: 'secret', action: 'update' }"
          @click="onEdit"
        >
          <Edit class=" h-4 w-4" />
          编辑
        </PermissionButton>
      </div>
    </div>

    <!-- 详情内容 -->
    <div>
      <Tabs default-value="details" class="w-full space-y-2">
        <TabsList class="grid w-full grid-cols-2">
          <TabsTrigger value="details">
            详情信息
          </TabsTrigger>
          <TabsTrigger value="history">
            <History class="h-4 w-4 mr-2" />
            修改历史
          </TabsTrigger>
        </TabsList>

        <TabsContent value="details" class="space-y-6">
          <DetailInfo
            ref="infoRef"
            :secret-data-title="secretDataTitle"
            :secret-data-icon="secretDataIcon"
            :error-text="errorText"
          >
            <template #secret-data="{ item }">
              <slot
                name="secret-data"
                :item="item"
              />
            </template>
          </DetailInfo>
        </TabsContent>

        <TabsContent value="history" class="space-y-6">
          <SecretHistory
            :item-id="route.params.id as string"
            @restored="handleHistoryRestore"
          />
        </TabsContent>
      </Tabs>
    </div>
  </div>
</template>
