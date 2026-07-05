// frontend/src/router/index.ts
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  // 公开文档视图
  {
    path: '/:sitePath/:pagePath(.*)?',
    name: 'public-doc',
    component: () => import('@/views/PublicDoc.vue'),
  },
  // 后台
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('@/views/admin/Login.vue'),
  },
  {
    path: '/admin',
    component: () => import('@/views/admin/Layout.vue'),
    beforeEnter: async (to, from, next) => {
      const auth = useAuthStore()
      if (!auth.user) await auth.fetchMe()
      if (!auth.user) {
        next('/admin/login')
        return
      }
      next()
    },
    children: [
      { path: '', redirect: '/admin/sites' },
      { path: 'sites', name: 'admin-sites', component: () => import('@/views/admin/SiteList.vue') },
      { path: 'sites/new', name: 'admin-site-new', component: () => import('@/views/admin/SiteEdit.vue') },
      { path: 'sites/:id', name: 'admin-site-edit', component: () => import('@/views/admin/SiteEdit.vue') },
      { path: 'sites/:id/pages', name: 'admin-pages', component: () => import('@/views/admin/PageTree.vue') },
      { path: 'sites/:siteId/pages/:id', name: 'admin-page-edit', component: () => import('@/views/admin/PageEdit.vue') },
      {
        path: 'users',
        name: 'admin-users',
        component: () => import('@/views/admin/UserList.vue'),
        beforeEnter: (to, from, next) => {
          const auth = useAuthStore()
          if (!auth.isSuper()) {
            next('/admin/sites')
            return
          }
          next()
        },
      },
    ],
  },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/views/NotFound.vue') },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0 }
  },
})
