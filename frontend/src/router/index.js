import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../views/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  const isLoggedIn = !!localStorage.getItem('authToken')

  if (to.meta.requiresAuth && !isLoggedIn) {
    // 如果目标路由需要认证但用户未登录
    next({ name: 'Login' })
  } else if (to.name === 'Login' && isLoggedIn) {
    // 如果用户已登录，但试图访问登录页，则重定向到首页
    next({ name: 'Dashboard' })
  } else {
    // 其他情况正常放行
    next()
  }
})

export default router