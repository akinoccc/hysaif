<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import PermissionButton from '@/components/common/permission/PermissionButton.vue'
import { Button } from '@/components/ui/button'

interface ActionButton {
  text: string
  icon?: any
  disabled?: boolean
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  permission?: { resource: string, action: string }
  onClick?: () => void
  to?: string
}

const props = defineProps<{
  title: string
  description?: string
  createButton?: {
    text: string
    to?: string
    permission?: { resource: string, action: string }
    onClick?: () => void
  }
  actions?: ActionButton[]
}>()

const router = useRouter()

function handleCreateClick() {
  if (props.createButton?.onClick) {
    props.createButton.onClick()
  }
  else if (props.createButton?.to) {
    router.push(props.createButton.to)
  }
}

function handleActionClick(action: ActionButton) {
  if (action.onClick) {
    action.onClick()
  }
  else if (action.to) {
    router.push(action.to)
  }
}
</script>

<template>
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">
        {{ title }}
      </h1>
      <p v-if="description" class="text-muted-foreground">
        {{ description }}
      </p>
    </div>
    <div v-if="createButton || actions" class="flex items-center space-x-2">
      <!-- 自定义操作按钮 -->
      <template v-if="actions">
        <template v-for="action in actions" :key="action.text">
          <PermissionButton
            v-if="action.permission"
            :variant="action.variant || 'outline'"
            :permission="action.permission"
            @click="handleActionClick(action)"
          >
            <component :is="action.icon" v-if="action.icon" class="h-4 w-4" />
            {{ action.text }}
          </PermissionButton>
          <Button
            v-else
            :disabled="action.disabled"
            :variant="action.variant || 'outline'"
            @click="handleActionClick(action)"
          >
            <component :is="action.icon" v-if="action.icon" class="h-4 w-4" />
            {{ action.text }}
          </Button>
        </template>
      </template>

      <!-- 创建按钮 -->
      <PermissionButton
        v-if="createButton"
        :permission="createButton.permission!"
        @click="handleCreateClick"
      >
        <Plus class="h-4 w-4" />
        {{ createButton.text }}
      </PermissionButton>
    </div>
  </div>
</template>
