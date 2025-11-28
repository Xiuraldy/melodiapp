import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import { useAuthStore } from '@/stores/auth'
import SongsView from '@/views/SongsView.vue'
import UsersView from '@/views/UsersView.vue'
import ServicesView from '../views/ServicesView.vue'
import ServiceView from '@/views/ServiceView.vue'
import Me from '@/views/Me.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'auth',
      component: AuthView
    },
    {
      path: '/services',
      name: 'services',
      component: ServicesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/service/:id',
      name: 'service',
      component: ServiceView
    },
    {
      path: '/songs',
      name: 'songs',
      component: SongsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
      meta: { requiresAuth: true }
    },
    {
      path: '/me',
      name: 'me',
      component: Me,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.name === 'auth' && authStore.isLoggedIn) {
    next({ name: 'services' })
  } else if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    console.log('Bloqueado por falta de permisos')
    next({ name: 'auth' })
  } else {
    next()
  }
})

export default router
