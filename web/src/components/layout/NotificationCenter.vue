<script setup lang="ts">
import { Bell } from 'lucide-vue-next'
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

function toggleNotificationPanel() {
  router.push('/notifications')
}

onMounted(() => {
  // 初始加载未读通知数量
  notificationStore.loadUnreadCount()

  // 开始定期轮询
  notificationStore.startPolling(30000)
})

onUnmounted(() => {
  // 停止轮询
  notificationStore.stopPolling()
})
</script>

<template>
  <div class="notification-center">
    <!-- 通知铃铛图标 -->
    <div class="relative" @click.stop="toggleNotificationPanel">
      <Bell class="w-6 h-6" />
      <!-- 未读通知徽章 -->
      <span
        v-if="notificationStore.unreadCount > 0"
        class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
      >
        {{ notificationStore.unreadCount > 99 ? '99+' : notificationStore.unreadCount }}
      </span>
    </div>
  </div>
</template>
