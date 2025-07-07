import { createPinia } from 'pinia'
import piniaPersistPlugin from 'pinia-plugin-persistedstate'
import { createApp } from 'vue'
import App from './App.vue'
import { setupPermissionDirectives } from './directives/permission'
import router from './router'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

pinia.use(piniaPersistPlugin)

app.use(pinia)
app.use(router)

// 注册权限指令
setupPermissionDirectives(app)

// 应用初始化时检查用户状态并初始化权限
async function initializeApp() {
  const { useAuthStore } = await import('./stores/auth')
  const { usePermissionStore } = await import('./stores/permission')

  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()

  // 如果用户已登录，初始化权限缓存
  if (authStore.isAuthenticated && authStore.user) {
    try {
      await permissionStore.initializePermissions()
    }
    catch (error) {
      // 提示用户初始化权限失败
      console.error('应用启动时权限缓存初始化失败:', error)
    }
  }
}

app.mount('#app')

// 应用挂载后初始化权限
initializeApp()
