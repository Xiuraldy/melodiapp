<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { apiCall } from '@/services/utils'
import { onBeforeMount, ref, watch } from 'vue'
import { jwtDecode } from 'jwt-decode'

const authStore = useAuthStore()
const isAdmin = ref(false)

// --- LÓGICA DE ROLES ---
async function checkAdminPermission() {
  if (!authStore.token) {
    isAdmin.value = false
    return
  }
  try {
    const decoded: any = jwtDecode(authStore.token)
    const userId = decoded.id || decoded.user_id || (decoded.MapClaims ? decoded.MapClaims.user_id : null)
    if (userId) {
      const response = await fetch(`http://localhost:8080/users/${userId}`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (response.ok) {
        const userData = await response.json()
        const role = userData.role || userData.Role || ''
        isAdmin.value = role.toLowerCase() === 'admin' || role.toLowerCase() === 'administrador'
      }
    }
  } catch (error) {
    console.error('Error verificando admin:', error)
    isAdmin.value = false
  }
}

onBeforeMount(async () => {
  authStore.init()
  await checkAdminPermission()
})

watch(() => authStore.token, () => {
  checkAdminPermission()
})

async function logout() {
  try {
    await apiCall('/auth/logout', { method: 'POST' })
    authStore.clearSession()
    isAdmin.value = false
  } catch (error) {
    authStore.clearSession()
    isAdmin.value = false
  }
}
</script>

<template>
  <div class="app-layout">
    
    <header class="main-header">
      <div class="logo-container">
        <img src="/logo.png" alt="MelodiApp" class="app-logo"/>
      </div>
      
      <nav class="main-nav">
        <div class="nav-links" v-if="authStore.isLoggedIn">
          <RouterLink to="/services">Servicios</RouterLink>
          <RouterLink to="/songs">Canciones</RouterLink>
          <RouterLink v-if="isAdmin" to="/users">Usuarios</RouterLink>
          <RouterLink to="/me">Mi Perfil</RouterLink>
        </div>
        
        <button v-if="authStore.isLoggedIn" @click="logout" class="logout-btn">
          Salir
        </button>
      </nav>
    </header>

    <div class="content-wrapper">
        <RouterView />
    </div>

    <footer>
      <h3>© Proyecto SENA - MelodiApp</h3> 
    </footer>

  </div>
</template>

<style scoped>
/* --- LAYOUT DE PANTALLA COMPLETA --- */
.app-layout {
  display: grid;
  /* AUTO: Altura del header según su contenido/padding
     1FR: El resto del espacio disponible (aprox 70-80%)
     AUTO: Altura del footer según su contenido 
  */
  grid-template-rows: auto 1fr auto; 
  height: 100vh; /* Altura exacta de la ventana */
  width: 100vw;
  overflow: hidden; /* Evita scroll en la ventana principal */
  background-color: #f8fafc;
}

/* --- HEADER --- */
.main-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0rem 2rem; /* Padding controla la altura visual */
  background-color: white;
  border-bottom: 1px solid #eee; 
  z-index: 10;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  height: 10vh;
}

/* Logo */
.logo-container { display: flex; height: -webkit-fill-available; align-items: center; }
.app-logo {
  height: -webkit-fill-available;
  width: auto;  /* Mantiene la proporción correcta */
  object-fit: contain; /* Asegura que la imagen se vea nítida */
  transition: transform 0.3s ease; /* Efecto suave al pasar el mouse */
}

.app-logo:hover {
  transform: scale(1.05); /* Efecto de "zoom" sutil */
}

/* Navegación */
.main-nav { display: flex; align-items: center; gap: 25px; }
.nav-links { display: flex; align-items: center; gap: 15px; }

.main-nav a {
  text-decoration: none;
  color: #64748b;
  font-weight: 500;
  font-size: 0.9rem;
  padding: 5px 0;
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.main-nav a:hover { color: var(--color-secundary); background-color: transparent; }
.main-nav a.router-link-exact-active {
  color: var(--color-secundary);
  border-bottom-color: var(--color-secundary);
  font-weight: 600;
}

/* Botón Salir */
.logout-btn {
  cursor: pointer;
  background: none;
  border: 1px solid #fee2e2;
  border-radius: 6px;
  color: #ef4444; 
  font-weight: 600;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}
.logout-btn:hover { background-color: #fef2f2; border-color: #ef4444; }

/* --- CONTENIDO CENTRAL --- */
.content-wrapper {
    /* Esto habilita el scroll SOLO en el cuerpo si el contenido es largo */
    overflow-y: auto; 
    width: 100%;
    position: relative;
}

/* --- FOOTER --- */
footer {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: white; 
  border-top: 1px solid #eee;
  color: #94a3b8;
  font-weight: 500;
  font-size: 0.75rem;
  padding: 15px; /* Controla la altura del footer */
  z-index: 10;
}

/* --- RESPONSIVE --- */
@media (max-width: 768px) {
  .main-header {
    padding: .2rem;
    gap: 10px;
  }
  .main-nav {
    width: 100%;
    justify-content: flex-end;
  }
  .nav-links { gap: 10px; font-size: 0.8rem; }
}
</style>