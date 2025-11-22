import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { jwtDecode } from 'jwt-decode'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const token = ref(sessionStorage.getItem('token') || '')
  const userRole = ref('')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userRole.value === 'admin')

  function setSession(tokenStr: string) {
    if (!tokenStr) return
    token.value = tokenStr
    sessionStorage.setItem('token', tokenStr)

    try {
      const payload: any = jwtDecode(tokenStr)
      userRole.value = payload.role || payload.Role || ''
    } catch (e) {
      console.error('Error decodificando token:', e)
      clearSession()
    }
  }

  function clearSession() {
    token.value = ''
    userRole.value = ''
    sessionStorage.removeItem('token')
    router.push('/')
  }

  function init() {
    const t = sessionStorage.getItem('token')
    if (t) {
      setSession(t)
    }
  }

  return {
    token,
    isLoggedIn,
    isAdmin,
    userRole,
    setSession,
    init,
    clearSession
  }
})
