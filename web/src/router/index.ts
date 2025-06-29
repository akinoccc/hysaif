import { Bell, LayoutDashboard, Shield, UserCheck, Users } from 'lucide-vue-next'
import { createRouter, createWebHistory } from 'vue-router'
import { SECRET_ITEM_TYPE, SECRET_ITEM_TYPE_MAP } from '@/constants'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/components/common/Layout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: '/dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: '仪表板',
              icon: LayoutDashboard,
              order: 1,
              showInMenu: true,
            },
          },
        },
        {
          path: '/users',
          name: 'UserList',
          component: () => import('@/views/user/UserList.vue'),
          meta: {
            requiresAuth: true,
            roles: ['super_admin', 'sec_mgr'],
            menu: {
              title: '用户管理',
              icon: Users,
              order: 1,
              showInMenu: true,
            },
          },
        },
        {
          path: '/permission',
          name: 'PermissionManagement',
          component: () => import('@/views/permission/Permission.vue'),
          meta: {
            requiresAuth: true,
            roles: ['super_admin', 'sec_mgr'],
            menu: {
              title: '权限管理',
              icon: Shield,
              order: 3,
              showInMenu: true,
            },
          },
        },
        {
          path: '/access_requests',
          name: 'AccessRequests',
          component: () => import('@/views/access_request/AccessRequestsList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: '访问申请',
              icon: UserCheck,
              order: 4,
              showInMenu: true,
            },
          },
        },
        {
          path: '/notifications',
          name: 'Notifications',
          component: () => import('@/views/NotificationView.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: '通知中心',
              icon: Bell,
              order: 5,
              showInMenu: true,
            },
          },
        },
        {
          path: '/api_key',
          name: 'ApiKeyList',
          component: () => import('@/views/api_key/ApiKeyList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.ApiKey].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.ApiKey].icon,
              showInMenu: true,
            },
          },
        },
        {
          path: '/api_key/create',
          name: 'ApiKeyCreate',
          component: () => import('@/views/api_key/ApiKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/api_key/:id',
          name: 'ApiKeyDetail',
          component: () => import('@/views/api_key/ApiKeyDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/api_key/:id/edit',
          name: 'ApiKeyEdit',
          component: () => import('@/views/api_key/ApiKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/access_key',
          name: 'AccessKeyList',
          component: () => import('@/views/access_key/AccessKeyList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.AccessKey].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.AccessKey].icon,
              order: 3,
              showInMenu: true,
            },
          },
        },
        {
          path: '/access_key/create',
          name: 'AccessKeyCreate',
          component: () => import('@/views/access_key/AccessKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/access_key/:id',
          name: 'AccessKeyDetail',
          component: () => import('@/views/access_key/AccessKeyDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/access_key/:id/edit',
          name: 'AccessKeyEdit',
          component: () => import('@/views/access_key/AccessKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/ssh_key',
          name: 'SshKeyList',
          component: () => import('@/views/ssh_key/SshKeyList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.SshKey].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.SshKey].icon,
              order: 4,
              showInMenu: true,
            },
          },
        },
        {
          path: '/ssh_key/create',
          name: 'SshKeyCreate',
          component: () => import('@/views/ssh_key/SshKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/ssh_key/:id',
          name: 'SshKeyDetail',
          component: () => import('@/views/ssh_key/SshKeyDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/ssh_key/:id/edit',
          name: 'SshKeyEdit',
          component: () => import('@/views/ssh_key/SshKeyForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/password',
          name: 'PasswordList',
          component: () => import('@/views/password/PasswordList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Password].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Password].icon,
              order: 5,
              showInMenu: true,
            },
          },
        },
        {
          path: '/password/create',
          name: 'PasswordCreate',
          component: () => import('@/views/password/PasswordForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/password/:id',
          name: 'PasswordDetail',
          component: () => import('@/views/password/PasswordDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/password/:id/edit',
          name: 'PasswordEdit',
          component: () => import('@/views/password/PasswordForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/token',
          name: 'TokenList',
          component: () => import('@/views/token/TokenList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Token].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Token].icon,
              order: 7,
              showInMenu: true,
            },
          },
        },
        {
          path: '/token/create',
          name: 'TokenCreate',
          component: () => import('@/views/token/TokenForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/token/:id',
          name: 'TokenDetail',
          component: () => import('@/views/token/TokenDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/token/:id/edit',
          name: 'TokenEdit',
          component: () => import('@/views/token/TokenForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/custom',
          name: 'CustomList',
          component: () => import('@/views/custom/CustomList.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Custom].label,
              icon: SECRET_ITEM_TYPE_MAP[SECRET_ITEM_TYPE.Custom].icon,
              order: 8,
              showInMenu: true,
            },
          },
        },
        {
          path: '/custom/create',
          name: 'CustomCreate',
          component: () => import('@/views/custom/CustomForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/custom/:id',
          name: 'CustomDetail',
          component: () => import('@/views/custom/CustomDetail.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/custom/:id/edit',
          name: 'CustomEdit',
          component: () => import('@/views/custom/CustomForm.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/profile',
          name: 'Profile',
          component: () => import('@/views/user/Profile.vue'),
          meta: {
            requiresAuth: true,
            menu: {
              title: '个人设置',
              icon: 'Settings',
              order: 5,
              showInMenu: false, // 不在主菜单中显示，在用户菜单中
            },
          },
        },
      ],
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false },
    },
  ],
})

// 路由守卫
router.beforeEach(async (to, _, next) => {
  const authStore = useAuthStore()

  // 检查是否需要认证
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  // 已登录用户访问登录页面，重定向到仪表板
  if (to.name === 'Login' && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  // 如果用户已登录，进行基本角色权限检查
  if (authStore.isAuthenticated && authStore.user) {
    // 检查角色权限（基于路由meta配置）
    if (to.meta.roles && authStore.user) {
      const userRole = authStore.user.role
      const allowedRoles = to.meta.roles as string[]
      if (!allowedRoles.includes(userRole || '')) {
        next('/dashboard')
        return
      }
    }
  }

  next()
})

export default router
