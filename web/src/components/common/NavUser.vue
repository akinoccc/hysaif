<script setup lang="ts">
import { BadgeCheck, Bell, ChevronsUpDown, LogOut } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from '@/components/ui/sidebar'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import NotificationCenter from '../NotificationCenter.vue'
import { Badge } from '../ui/badge'

const authStore = useAuthStore()
const { isMobile } = useSidebar()
const router = useRouter()
const notificationStore = useNotificationStore()

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

async function gotoNotifications() {
  router.push('/notifications')
}

function gotoProfile() {
  router.push('/profile')
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar class="h-8 w-8 rounded-lg">
              <AvatarImage v-if="authStore.user?.avatar" :src="authStore.user?.avatar" :alt="authStore.user?.name" />
              <AvatarFallback class="rounded-lg">
                CN
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ authStore.user?.name }}</span>
              <span class="truncate text-xs">{{ authStore.user?.email }}</span>
            </div>
            <NotificationCenter class="mr-1" />
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage v-if="authStore.user?.avatar" :src="authStore.user?.avatar" :alt="authStore.user?.name" />
                <AvatarFallback class="rounded-lg">
                  CN
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ authStore.user?.name }}</span>
                <span class="truncate text-xs">{{ authStore.user?.email }}</span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem @click="gotoProfile">
              <BadgeCheck />
              Account
            </DropdownMenuItem>
            <DropdownMenuItem @click="gotoNotifications">
              <Bell />
              Notifications
              <Badge v-if="notificationStore.unreadCount > 0" variant="destructive">
                {{ notificationStore.unreadCount }}
              </Badge>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleLogout">
            <LogOut />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
