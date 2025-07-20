import { createRouter, createWebHistory } from 'vue-router'
import Zhuye from '../views/zhuye.vue'
import Podlist from '../views/podlist.vue'
import Bushu from '../views/bushu.vue'
import Login from '../views/login.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: Zhuye,
    meta: { requiresAuth: true }
  },
  {
    path: '/zhuye',
    name: 'zhuye',
    component: Zhuye,
    meta: { requiresAuth: true }
  },
  {
    path: '/podlist',
    name: 'podlist',
    component: Podlist,
    meta: { requiresAuth: true }
  },
  {
    path: '/bushu',
    name: 'bushu',
    component: Bushu,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})


// 🔐 添加全局路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    // 没有 token 且访问需要登录的页面，跳转登录页
    next('/login')
  } else if (to.path === '/login' && token) {
    // 已登录还访问登录页，重定向主页
    next('/zhuye')
  } else {
    next() // 正常放行
  }
})

export default router
