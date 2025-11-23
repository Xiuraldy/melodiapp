<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

// --- 1. INTERFACES ---
interface User {
  id: number;
  username: string;
  lastname: string;
  role: string;      
  subRole: string;   
  email: string;
  celphone: string;
  photo: string;
}

// --- 2. VARIABLES Y STORE ---
const authStore = useAuthStore()
const allUsers = ref<User[]>([])
const isLoading = ref(true)
const errorMsg = ref('')

// --- 3. CARGAR TODOS LOS USUARIOS (GET) ---
async function getUsers() {
  isLoading.value = true
  errorMsg.value = ''
  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await response.json()
    
    allUsers.value = data.map((u: any) => ({
      id: u.ID || u.id,
      username: u.Username || u.username,
      lastname: u.Lastname || u.lastname || '',
      email: u.Email || u.email,
      celphone: u.Celphone || u.celphone || '',
      role: u.Role || u.role || 'General',
      subRole: u.SubRole || u.subRole || 'Ninguno',
      photo: '/logo.png'
    }))
  } catch (e: any) {
    console.error(e)
    errorMsg.value = 'Error cargando usuarios.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => { getUsers() })

// --- 4. FILTROS Y HELPERS ---
const searchQuery = ref('')
const filteredUsers = computed(() => {
  if (!searchQuery.value) return allUsers.value
  const term = searchQuery.value.toLowerCase()
  return allUsers.value.filter(user => 
    user.username.toLowerCase().includes(term) || 
    user.lastname.toLowerCase().includes(term) ||
    user.subRole.toLowerCase().includes(term)
  )
})

function getRoleEmoji(subRole: string) {
  if (!subRole) return '👤'
  const r = subRole.toLowerCase()
  if (r.includes('cantante') || r.includes('voz')) return '🎤'
  if (r.includes('guitarra')) return '🎸'
  if (r.includes('batería') || r.includes('bateria')) return '🥁'
  if (r.includes('piano') || r.includes('teclado')) return '🎹'
  if (r.includes('bajo')) return '🎸'
  if (r.includes('sonido') || r.includes('audio')) return '🎚'
  return '🎵'
}

function deleteUser(id: number) {
  if(confirm('¿Eliminar usuario localmente?')) {
    allUsers.value = allUsers.value.filter(u => u.id !== id)
  }
}

// Listas de opciones
const mainRoles = ['Administrador', 'General']
const instrumentRoles = ['Cantante', 'Guitarrista', 'Guitarrista Eléctrico', 'Baterista', 'Pianista', 'Bajista', 'Saxofonista', 'Sonidista', 'Líder de Alabanza', 'Ninguno']

// --- 5. LÓGICA DE CREACIÓN (POST) ---
const isCreateModalOpen = ref(false)
const isSaving = ref(false)
const createError = ref('')

const initialForm = { username: '', lastname: '', email: '', password: '', celphone: '', role: '', subRole: '' }
const newUserForm = reactive({ ...initialForm })

function openCreateModal() { Object.assign(newUserForm, initialForm); createError.value = ''; isCreateModalOpen.value = true }
function closeCreateModal() { isCreateModalOpen.value = false }

async function saveNewUser() {
  isSaving.value = true
  createError.value = ''
  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
      body: JSON.stringify(newUserForm)
    })
    const data = await response.json()
    if (!response.ok) throw new Error(data.error || 'Error al crear')
    
    await getUsers() 
    closeCreateModal()
  } catch (e: any) {
    createError.value = e.message
  } finally {
    isSaving.value = false
  }
}

// --- 6. LÓGICA DE EDICIÓN (PUT + GET/:id) ---
const isEditModalOpen = ref(false)
const editError = ref('')
const editingUserId = ref<number | null>(null)

const editUserForm = reactive({ 
  username: '', lastname: '', email: '', password: '', celphone: '', role: '', subRole: '' 
})

async function openEditModal(id: number) {
  isEditModalOpen.value = true
  editError.value = ''
  editingUserId.value = id
  
  Object.assign(editUserForm, { username: 'Cargando...', lastname: '', email: '', celphone: '', role: '', subRole: '' })

  try {
    const response = await fetch(`http://localhost:8080/users/${id}`, {
      method: 'GET',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}` 
      }
    })

    if (response.status === 401) {
      alert('Tu sesión ha expirado.')
      authStore.clearSession()
      return
    }

    if (!response.ok) throw new Error('No se pudo cargar el usuario')

    const data = await response.json()

    editUserForm.username = data.Username || data.username
    editUserForm.lastname = data.Lastname || data.lastname || ''
    editUserForm.email = data.Email || data.email
    editUserForm.celphone = data.Celphone || data.celphone || ''
    editUserForm.role = data.Role || data.role
    editUserForm.subRole = data.SubRole || data.subRole || ''
    editUserForm.password = '' 

  } catch (e: any) {
    console.error(e)
    editError.value = 'Error de conexión: ' + e.message
  }
}

function closeEditModal() {
  isEditModalOpen.value = false
  editingUserId.value = null
}

// --- AQUÍ ESTÁ LA CORRECCIÓN IMPORTANTE ---
async function saveEditedUser() {
  if (!editingUserId.value) return
  isSaving.value = true
  editError.value = ''

  try {
    // 1. CREAMOS UNA COPIA DEL FORMULARIO
    // Usamos 'any' temporalmente para poder borrar propiedades sin que TS se queje
    const payload: any = { ...editUserForm }

    // 2. SI LA CONTRASEÑA ESTÁ VACÍA, LA ELIMINAMOS DEL ENVÍO
    // Esto evita que se envíe "password": "" y el backend lance error 400
    if (!payload.password || payload.password.trim() === '') {
      delete payload.password
    }

    const response = await fetch(`http://localhost:8080/users/${editingUserId.value}`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json', 
        'Authorization': `Bearer ${authStore.token}` 
      },
      // 3. ENVIAMOS EL PAYLOAD LIMPIO
      body: JSON.stringify(payload)
    })

    const data = await response.json()
    
    if (!response.ok) {
      // Si el backend devuelve error, mostramos el mensaje exacto
      throw new Error(data.error || 'Error al editar')
    }

    await getUsers()
    closeEditModal()

  } catch (e: any) {
    console.error(e)
    editError.value = e.message
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <main>
    <div class="programation">
      
      <div class="view-header-row">
        <div>
           <h1>Usuarios del Sistema</h1>
           <p class="subtitle">Administra permisos y músicos</p>
        </div>
        <div class="header-actions">
           <div class="search-box">
              <span class="search-icon">🔍</span>
              <input v-model="searchQuery" type="text" placeholder="Buscar por nombre o instrumento..." />
           </div>
           <button class="btn-create" @click="openCreateModal">
             + Nuevo Usuario
           </button>
        </div>
      </div>

      <div v-if="errorMsg" class="error-banner">⚠ {{ errorMsg }}</div>

      <div class="card users-card">
         <div class="card-header-flex">
             <h3>Miembros ({{ filteredUsers.length }})</h3>
             <span v-if="isLoading" class="loading-text">Cargando...</span>
         </div>

         <table class="custom-table">
            <thead>
               <tr>
                  <th style="width: 70px;">Foto</th>
                  <th>Nombre Completo</th>
                  <th>Rol / Instrumento</th>
                  <th>Contacto</th>
                  <th class="center-text">Acciones</th>
               </tr>
            </thead>
            <tbody>
               <tr v-for="user in filteredUsers" :key="user.id">
                  <td><div class="avatar-wrapper"><img :src="user.photo" alt="avatar" class="avatar" /></div></td>
                  <td><div class="info-stack"><span class="primary-text">{{ user.username }} {{ user.lastname }}</span></div></td>
                  <td>
                     <div class="info-stack">
                        <div class="role-badge-wrapper">
                           <span :class="['role-badge', user.role === 'Administrador' ? 'admin' : 'general']">
                              {{ user.role }}
                           </span>
                        </div>
                        <span class="secondary-text subrole-text">
                           {{ getRoleEmoji(user.subRole) }} {{ user.subRole }}
                        </span>
                     </div>
                  </td>
                  <td><div class="info-stack"><span class="contact-text email">✉ {{ user.email }}</span><span class="contact-text phone">📱 {{ user.celphone }}</span></div></td>
                  <td class="center-text">
                     <div class="actions-row">
                        <button @click="openEditModal(user.id)" class="btn-action edit">✎</button>
                        <button @click="deleteUser(user.id)" class="btn-action delete">🗑</button>
                     </div>
                  </td>
               </tr>
            </tbody>
         </table>
      </div>

      <div v-if="isCreateModalOpen" class="modal-overlay" @click.self="closeCreateModal">
        <div class="modal-content form-modal">
          <button class="close-btn" @click="closeCreateModal">×</button>
          <div class="modal-header"><h3>Nuevo Usuario</h3><p>Registra un nuevo miembro</p></div>
          <div v-if="createError" class="modal-error">{{ createError }}</div>

          <form @submit.prevent="saveNewUser" class="user-form">
            <div class="form-row">
              <div class="form-group"><label>Nombre *</label><input v-model="newUserForm.username" type="text" required /></div>
              <div class="form-group"><label>Apellido</label><input v-model="newUserForm.lastname" type="text" /></div>
            </div>
            <div class="form-row">
              <div class="form-group"><label>Email *</label><input v-model="newUserForm.email" type="email" required /></div>
              <div class="form-group"><label>Celular</label><input v-model="newUserForm.celphone" type="tel" /></div>
            </div>
            <div class="form-row">
               <div class="form-group"><label>Contraseña *</label><input v-model="newUserForm.password" type="password" required /></div>
               <div class="form-group"><label>Permisos *</label><select v-model="newUserForm.role" required><option v-for="r in mainRoles" :value="r">{{r}}</option></select></div>
            </div>
            <div class="form-row">
               <div class="form-group full-width"><label>Instrumento *</label><select v-model="newUserForm.subRole" required><option v-for="i in instrumentRoles" :value="i">{{i}}</option></select></div>
            </div>
            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closeCreateModal">Cancelar</button>
              <button type="submit" class="btn-save" :disabled="isSaving">{{ isSaving ? '...' : 'Crear' }}</button>
            </div>
          </form>
        </div>
      </div>

      <div v-if="isEditModalOpen" class="modal-overlay" @click.self="closeEditModal">
        <div class="modal-content form-modal">
          <button class="close-btn" @click="closeEditModal">×</button>
          
          <div class="modal-header">
            <h3>Editar Usuario</h3>
            <p>Modifica los datos del miembro</p>
          </div>

          <div v-if="editError" class="modal-error">{{ editError }}</div>

          <form @submit.prevent="saveEditedUser" class="user-form">
            <div class="form-row">
              <div class="form-group">
                <label>Nombre *</label>
                <input v-model="editUserForm.username" type="text" required />
              </div>
              <div class="form-group">
                <label>Apellido</label>
                <input v-model="editUserForm.lastname" type="text" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Email *</label>
                <input v-model="editUserForm.email" type="email" required />
              </div>
              <div class="form-group">
                <label>Celular</label>
                <input v-model="editUserForm.celphone" type="tel" />
              </div>
            </div>

            <div class="form-row">
               <div class="form-group">
                <label>Contraseña</label>
                <input v-model="editUserForm.password" type="password" placeholder="Dejar vacía para mantener" />
              </div>
              <div class="form-group">
                <label>Permisos *</label>
                <select v-model="editUserForm.role" required>
                  <option v-for="role in mainRoles" :key="role" :value="role">{{ role }}</option>
                </select>
              </div>
            </div>
            
            <div class="form-row">
               <div class="form-group full-width">
                <label>Instrumento *</label>
                <select v-model="editUserForm.subRole" required>
                  <option v-for="inst in instrumentRoles" :key="inst" :value="inst">{{ inst }}</option>
                </select>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closeEditModal" :disabled="isSaving">Cancelar</button>
              <button type="submit" class="btn-save" :disabled="isSaving">
                {{ isSaving ? 'Guardando...' : 'Actualizar Usuario' }}
              </button>
            </div>
          </form>
        </div>
      </div>

    </div>
  </main>
</template>

<style scoped>
/* === ESTILOS COMPARTIDOS === */
.programation { display: flex; flex-direction: column; padding: 40px; gap: 20px; }
.view-header-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; flex-wrap: wrap; gap: 20px; }
h1 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.5rem; }
.subtitle { margin: 0; color: #888; font-size: 0.9rem; }
.header-actions { display: flex; gap: 15px; align-items: center; }
.search-box { position: relative; display: flex; align-items: center; }
.search-box input { padding: 8px 10px 8px 35px; border-radius: 8px; border: 1px solid #e5e7eb; font-size: 0.9rem; outline: none; width: 250px; transition: border 0.2s; }
.search-box input:focus { border-color: var(--color-secundary, #2c3e50); }
.search-icon { position: absolute; left: 10px; font-size: 0.8rem; opacity: 0.6; }
.btn-create { background-color: #10b981; color: white; border: none; padding: 8px 16px; border-radius: 8px; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-create:hover { background-color: #059669; }
.card { background: white; border: 1px solid #eee; padding: 20px; border-radius: 16px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); width: 100%; overflow-x: auto;}
.card-header-flex { border-bottom: 2px solid #f0f0f0; padding-bottom: 10px; margin-bottom: 15px; display: flex; justify-content: space-between; }
.card h3 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.2rem; }
.custom-table { width: 100%; border-collapse: collapse; }
.custom-table th { text-align: left; font-size: 0.75rem; text-transform: uppercase; color: #999; padding: 10px; font-weight: 600; }
.custom-table td { padding: 12px 15px; border-bottom: 1px dashed #eee; vertical-align: middle; }
.center-text { text-align: center; }
.info-stack { display: flex; flex-direction: column; justify-content: center; }
.primary-text { font-weight: 700; color: #333; font-size: 0.95rem; }
.secondary-text { font-size: 0.8rem; color: #888; margin-top: 2px; }
.contact-text { font-size: 0.85rem; color: #555; margin-bottom: 2px; }
.contact-text.email { color: var(--color-secundary, #2c3e50); font-weight: 500; }
.avatar-wrapper { position: relative; width: 45px; height: 45px; }
.avatar { width: 100%; height: 100%; border-radius: 50%; object-fit: cover; border: 1px solid #eee; }
.actions-row { display: flex; justify-content: center; gap: 8px; }
.btn-action { border: none; width: 32px; height: 32px; border-radius: 6px; cursor: pointer; font-size: 1rem; transition: transform 0.1s; display: flex; justify-content: center; align-items: center;}
.btn-action:hover { transform: scale(1.1); }
.btn-action.edit { background-color: #f0fdfa; color: #0d9488; border: 1px solid #ccfbf1; }
.btn-action.delete { background-color: #fef2f2; color: #ef4444; border: 1px solid #fee2e2; }
.error-banner { background-color: #fee2e2; color: #b91c1c; padding: 10px; border-radius: 8px; margin-bottom: 15px; text-align: center; }
.loading-text { color: #888; font-size: 0.9rem; font-style: italic; }

.role-badge-wrapper { display: flex; }
.role-badge { padding: 4px 10px; border-radius: 20px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
.role-badge.admin { background-color: #f3e8ff; color: #7e22ce; border: 1px solid #d8b4fe; }
.role-badge.general { background-color: #eff6ff; color: #1e40af; border: 1px solid #dbeafe; }
.subrole-text { font-style: italic; color: #555; font-weight: 500; display: flex; align-items: center; gap: 5px; }

/* MODAL */
.modal-content.form-modal { max-width: 600px; width: 95%; }
.user-form { display: flex; flex-direction: column; gap: 15px; margin-top: 10px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.full-width { grid-column: span 2; } 
.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-group label { font-size: 0.85rem; font-weight: 600; color: #4b5563; }
.form-group input, .form-group select { padding: 10px; border-radius: 8px; border: 1px solid #d1d5db; font-size: 0.95rem; outline: none; background-color: #f9fafb; font-family: inherit; }
.form-group input:focus, .form-group select:focus { border-color: var(--color-secundary, #2c3e50); background-color: white; box-shadow: 0 0 0 3px rgba(44, 62, 80, 0.1); }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 15px; border-top: 1px solid #eee; padding-top: 20px; }
.btn-cancel { background: white; border: 1px solid #d1d5db; color: #4b5563; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.btn-save { background: var(--color-secundary, #2c3e50); border: none; color: white; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.btn-save:disabled { opacity: 0.7; cursor: not-allowed; }
.modal-error { color: #ef4444; background-color: #fef2f2; padding: 10px; border-radius: 8px; font-size: 0.9rem; border: 1px solid #fee2e2; margin-bottom: 10px; text-align: center; }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.modal-content { background: white; padding: 30px; border-radius: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); position: relative; animation: slideUp 0.3s ease-out; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
.close-btn { position: absolute; top: 15px; right: 20px; background: none; border: none; font-size: 2rem; color: #999; cursor: pointer; }
.modal-header { text-align: center; margin-bottom: 20px; }
.modal-header h3 { margin: 0; color: var(--color-secundary, #2c3e50); }
.modal-header p { margin: 5px 0 0; color: #666; font-size: 0.9rem; }
@media (max-width: 600px) { .form-row { grid-template-columns: 1fr; gap: 10px; } .full-width { grid-column: span 1; } }
</style>