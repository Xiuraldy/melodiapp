<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const activeLogin = ref(true)

// --- LOGIN ---
const loginInputs = reactive({ email: '', password: '' })
const loginError = ref('') 
const isLoggingIn = ref(false)

async function signIn() {
  isLoggingIn.value = true
  loginError.value = ''
  try {
    const response = await fetch('http://localhost:8080/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginInputs),
    })
    const data = await response.json()
    if (!response.ok) throw new Error(data.error || 'Error en credenciales')
    authStore.setSession(data.token)
    router.push('/services')
  } catch (e: any) {
    loginError.value = e.message
  } finally {
    isLoggingIn.value = false
  }
}

// --- REGISTRO ---
const registerInputs = reactive({
  username: '',
  lastname: '',
  email: '',
  password: '',
  celphone: ''
})

// 1. CAMBIO: Array para múltiples selecciones
const selectedInstruments = ref<string[]>([]) 

const registerError = ref('') 
const isRegistering = ref(false)

const instrumentOptions = [
  'Cantante', 'Guitarrista', 'Guitarrista Eléctrico', 
  'Pianista', 'Saxofonista', 'Baterista', 'Bajista', 'Sonidista'
]

async function register() {
  isRegistering.value = true
  registerError.value = ''

  try {
    // 2. CAMBIO: Convertir el Array a String separado por comas
    const rolesString = selectedInstruments.value.join(',')

    const payload = {
      username: registerInputs.username,
      email: registerInputs.email,
      password: registerInputs.password,
      role: "general",
      celphone: registerInputs.celphone,
      lastname: registerInputs.lastname,
      secondary_role: rolesString // Enviamos el string combinado
    }

    console.log("Enviando body:", payload)

    const response = await fetch('http://localhost:8080/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    const data = await response.json()

    if (!response.ok) throw new Error(data.error || 'Error al registrarse')

    authStore.setSession(data.token)
    router.push('/services')

  } catch (e: any) {
    registerError.value = e.message
  } finally {
    isRegistering.value = false
  }
}
</script>

<template>
  <main class="auth-container">
    
    <div class="image-panel" :class="{ 'slide-right': !activeLogin }">
      <img v-if="activeLogin" src="/auth/group 2.jpg" alt="Login Image" class="auth-img fade-in">
      <img v-else src="/auth/group.jpg" alt="Register Image" class="auth-img fade-in">
      <div class="overlay">
        <h2>{{ activeLogin ? 'Bienvenido de nuevo' : 'Únete al equipo' }}</h2>
        <p>{{ activeLogin ? 'Sigue gestionando tus eventos.' : 'Comienza a organizar tus participaciones.' }}</p>
      </div>
    </div>

    <div class="form-panel">
      
      <div v-if="activeLogin" class="form-wrapper fade-in">
        <h1>Iniciar Sesión</h1>
        <p class="subtitle">Ingresa tus credenciales para continuar</p>
        <form @submit.prevent="signIn">
          <div class="input-group">
            <input v-model="loginInputs.email" type="email" placeholder="Correo Electrónico" required />
          </div>
          <div class="input-group">
            <input v-model="loginInputs.password" type="password" placeholder="Contraseña" required />
          </div>
          <div v-if="loginError" class="error-banner">{{ loginError }}</div>
          <button type="submit" :disabled="isLoggingIn">
            {{ isLoggingIn ? 'Entrando...' : 'Iniciar Sesión' }}
          </button>
        </form>
        <p class="switch-text">¿No tienes usuario? <a @click="activeLogin = false">Regístrate aquí</a></p>
      </div>
  
      <div v-else class="form-wrapper fade-in">
        <h1>Crear Cuenta</h1>
        <p class="subtitle">Completa tus datos para registrarte</p>
        
        <form @submit.prevent="register">
          <div class="row">
            <input v-model="registerInputs.username" type="text" placeholder="Nombre" required />
            <input v-model="registerInputs.lastname" type="text" placeholder="Apellido" required />
          </div>
          
          <div class="input-group">
            <input v-model="registerInputs.email" type="email" placeholder="Correo Electrónico" required />
          </div>

          <div class="input-group">
             <input v-model="registerInputs.celphone" type="tel" placeholder="Celular" />
          </div>

          <div class="roles-section">
            <label class="roles-label">Soy... (Selecciona uno o varios)</label>
            <div class="checkbox-grid">
              <label 
                v-for="inst in instrumentOptions" 
                :key="inst" 
                class="checkbox-item"
                :class="{ 'checked': selectedInstruments.includes(inst.toLowerCase()) }"
              >
                <input 
                  type="checkbox" 
                  :value="inst.toLowerCase()" 
                  v-model="selectedInstruments" 
                />
                <span>{{ inst }}</span>
              </label>
            </div>
          </div>

          <div class="input-group">
            <input v-model="registerInputs.password" type="password" placeholder="Contraseña" required />
          </div>

          <div v-if="registerError" class="error-banner">{{ registerError }}</div>

          <button type="submit" :disabled="isRegistering">
            {{ isRegistering ? 'Creando...' : 'Registrarse' }}
          </button>
        </form>

        <p class="switch-text">¿Ya tienes usuario? <a @click="activeLogin = true">Inicia Sesión</a></p>
      </div>

    </div>
  </main>
</template>

<style scoped>
/* LAYOUT PRINCIPAL */
.auth-container { display: flex; height: -webkit-fill-available; width: 98vw; background-color: #f8fafc; overflow: hidden; }
.image-panel { width: 50%; position: relative; transition: all 0.5s ease-in-out; overflow: hidden; }
.auth-img { width: 100%; height: 100%; object-fit: cover; filter: brightness(0.7); }
.overlay { position: absolute; bottom: 10%; left: 10%; color: white; text-shadow: 0 2px 10px rgba(0,0,0,0.5); z-index: 2; }
.overlay h2 { font-size: 2.5rem; margin-bottom: 10px; font-weight: 700; }
.overlay p { font-size: 1.1rem; max-width: 80%; }
.form-panel { width: 50%; display: flex; align-items: center; justify-content: center; background-color: white; }
.form-wrapper { width: 100%; max-width: 450px; padding: 40px; display: flex; flex-direction: column; overflow-y: auto; max-height: 80vh; } /* Agregué scroll por si la lista es larga */

h1 { color: #1e293b; font-size: 2rem; margin-bottom: 5px; font-weight: 700; }
.subtitle { color: #64748b; margin-bottom: 20px; font-size: 0.95rem; }

form { display: flex; flex-direction: column; gap: 15px; }
.input-group { width: 100%; }
.row { display: flex; gap: 15px; }
.row input { flex: 1; }

input {
  width: 100%; padding: 12px 15px; border-radius: 8px; border: 1px solid #e2e8f0;
  font-size: 0.95rem; outline: none; background-color: #f8fafc; transition: all 0.2s; font-family: inherit;
}
input:focus { border-color: #10b981; background-color: white; box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1); }

/* --- ESTILOS DE LA CHECKLIST --- */
.roles-section { margin-top: 5px; }
.roles-label { font-size: 0.85rem; color: #64748b; font-weight: 600; margin-bottom: 8px; display: block; }

.checkbox-grid {
  display: grid;
  grid-template-columns: 1fr 1fr; /* Dos columnas */
  gap: 10px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #f8fafc;
  font-size: 0.85rem;
  color: #475569;
}

.checkbox-item:hover { background-color: #f1f5f9; }

/* Estilo cuando está seleccionado */
.checkbox-item.checked {
  background-color: #ecfdf5; /* Fondo verde muy claro */
  border-color: #10b981;     /* Borde verde */
  color: #065f46;            /* Texto verde oscuro */
  font-weight: 600;
}

.checkbox-item input { display: none; } /* Ocultamos el checkbox nativo feo */

button { margin-top: 10px; padding: 12px; border-radius: 8px; border: none; background-color: #10b981; color: white; font-weight: 600; font-size: 1rem; cursor: pointer; transition: background 0.2s, transform 0.1s; }
button:hover { background-color: #059669; }
button:active { transform: scale(0.98); }
button:disabled { opacity: 0.7; cursor: not-allowed; }

.switch-text { margin-top: 20px; text-align: center; color: #64748b; font-size: 0.9rem; }
.switch-text a { color: #10b981; font-weight: 600; cursor: pointer; text-decoration: none; }
.switch-text a:hover { text-decoration: underline; }
.error-banner { background-color: #fef2f2; color: #ef4444; padding: 10px; border-radius: 6px; font-size: 0.85rem; text-align: center; border: 1px solid #fee2e2; }
.fade-in { animation: fadeIn 0.5s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

@media (max-width: 768px) {
  .auth-container { flex-direction: column; }
  .image-panel { display: none; }
  .form-panel { width: 100%; height: 100%; }
  .form-wrapper { padding: 20px; }
}
</style>