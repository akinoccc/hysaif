<script setup lang="ts">
import {
  Bell,
  Braces,
  Coins,
  FileText,
  Key,
  KeyRound,
  LayoutDashboard,
  Lock,
  Moon,
  Shield,
  Sun,
  Terminal,
  Users,
} from 'lucide-vue-next'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { type MenuItemData, permissionAPI } from '@/api/permission'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { useTheme } from '@/composables/useTheme'
import { useAuthStore } from '@/stores/auth'
import NavUser from './NavUser.vue'

const route = useRoute()
const authStore = useAuthStore()
const sidebarOpen = ref(true)
const menuItems = ref<MenuItemData[]>([])
const loading = ref(false)

const { isDark, toggleTheme } = useTheme()

// 图标映射
const iconMap = {
  LayoutDashboard,
  Users,
  Shield,
  FileText,
  Key,
  KeyRound,
  Terminal,
  Lock,
  Coins,
  Braces,
  Bell,
}

function isMenuActive(menuPath: string, currentPath: string): boolean {
  if (menuPath === currentPath) {
    return true
  }

  // 检查是否为子路径
  if (currentPath.startsWith(menuPath) && menuPath !== '/') {
    return true
  }

  return false
}

// 获取用户可访问的菜单
async function fetchUserMenus() {
  try {
    loading.value = true
    const response = await permissionAPI.getUserAccessibleMenus()
    if (response.data?.menus) {
      menuItems.value = response.data.menus.sort((a, b) => a.order - b.order)
    }
  }
  catch (error) {
    console.error('获取菜单失败:', error)
    menuItems.value = []
  }
  finally {
    loading.value = false
  }
}

// 获取图标组件
function getIconComponent(iconName: string) {
  return iconMap[iconName as keyof typeof iconMap] || Shield
}

// 组件挂载时获取菜单
onMounted(() => {
  if (authStore.isAuthenticated) {
    fetchUserMenus()
  }
})

// 监听认证状态变化
watch(() => authStore.isAuthenticated, () => {
  if (authStore.isAuthenticated && menuItems.value.length === 0) {
    fetchUserMenus()
  }
  else if (!authStore.isAuthenticated) {
    menuItems.value = []
  }
})
</script>

<template>
  <SidebarProvider :open="sidebarOpen" @update:open="sidebarOpen = $event">
    <!-- 侧边栏 -->
    <Sidebar side="left" variant="sidebar" collapsible="offcanvas" class="theme-transition">
      <!-- 侧边栏头部 -->
      <SidebarHeader class="py-4">
        <div class="flex items-center justify-between px-4">
          <router-link to="/" class="flex items-center space-x-2">
            <!-- <Shield class="h-6 w-6 text-primary" /> -->
            <img src="/logo.svg" alt="HySAIF Logo" class="h-8 w-8 text-primary">
            <span class="font-bold text-lg text-sidebar-foreground">Hysaif SIMS</span>
          </router-link>
          <!-- 主题切换按钮 -->
          <Button
            variant="ghost"
            size="icon"
            class="h-10 w-10 text-sidebar-foreground hover:bg-sidebar-accent cursor-pointer"
            @click="toggleTheme"
          >
            <Sun v-if="isDark" class="h-4 w-4" />
            <Moon v-else class="h-4 w-4" />
          </Button>
        </div>
      </SidebarHeader>

      <Separator />

      <!-- 导航菜单 -->
      <SidebarContent class="px-2 py-4">
        <SidebarMenu>
          <SidebarMenuItem
            v-for="item in menuItems" :key="item.path"
          >
            <SidebarMenuSkeleton />
            <SidebarMenuButton
              class="h-fit py-0 theme-transition"
              :is-active="isMenuActive(item.path, route.path)"
            >
              <RouterLink :to="item.path" class="w-full">
                <div class="flex items-center gap-2 px-4 py-2.5">
                  <component :is="getIconComponent(item.icon)" class="h-4 w-4" />
                  <span>{{ item.title }}</span>
                </div>
              </RouterLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarContent>

      <!-- 用户信息区域 -->
      <SidebarFooter>
        <Separator />
        <NavUser />
      </SidebarFooter>
    </Sidebar>

    <SidebarTrigger class="absolute top-0 left-0 z-1 md:relative" />

    <!-- 主要内容区域 -->
    <SidebarInset class="theme-transition">
      <!-- 主要内容 -->
      <div class="p-6">
        <router-view />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
