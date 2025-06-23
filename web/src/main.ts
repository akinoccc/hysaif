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

app.mount('#app')
