<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { jwtDecode } from 'jwt-decode'

const authStore = useAuthStore()

// --- 1. INTERFACES ---
interface User {
  id: number;
  username: string;
  lastname: string;
  email: string;
  celphone: string;
  role: string;
  secondary_role: string;
  photo: string; 
}

// --- 2. ESTADO ---
const allUsers = ref<User[]>([])
const isLoading = ref(true)
const errorMsg = ref('')

// Opciones
const mainRoles = ['admin', 'general']
const instrumentOptions = [
  'Cantante', 'Guitarrista', 'Guitarrista Eléctrico', 
  'Pianista', 'Saxofonista', 'Baterista', 'Bajista', 'Sonidista', 'Producción'
]

// --- HELPER URL FOTOS ---
function getPhotoUrl(path: string) {
  if (!path) return '/logo.png'; 
  if (path.startsWith('http')) return path; 
  if (path.startsWith('blob:')) return path; // Permitir previsualizaciones locales
  return `http://localhost:8080${path}`;
}

// --- 3. CARGAR USUARIOS ---
async function getUsers() {
  isLoading.value = true
  errorMsg.value = ''
  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    
    if (!response.ok) throw new Error('Error cargando usuarios')
    
    const data = await response.json()
    
    allUsers.value = data.map((u: any) => ({
      id: u.id,
      username: u.username,
      lastname: u.lastname || '',
      email: u.email,
      celphone: u.celphone || '',
      role: u.role || 'general',
      secondary_role: u.secondary_role || '',
      photo: getPhotoUrl(u.profile_picture_url) 
    }))

  } catch (e: any) {
    console.error(e)
    errorMsg.value = 'No se pudo conectar con el servidor.'
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
    user.secondary_role.toLowerCase().includes(term)
  )
})

function getRoleEmoji(secRole: string) {
  if (!secRole) return '👤'
  const r = secRole.toLowerCase()
  if (r.includes('cantante')) return '🎤'
  if (r.includes('guitarra')) return '🎸'
  if (r.includes('batería') || r.includes('baterista')) return '🥁'
  if (r.includes('piano')) return '🎹'
  if (r.includes('bajo')) return '🎸'
  if (r.includes('sonido')) return '🎚'
  return '🎵'
}

function getMainRoleBadgeClass(role: string) {
  return (role.toLowerCase() === 'admin' || role.toLowerCase() === 'administrador') ? 'admin' : 'general'
}

// --- 5. ELIMINAR ---
async function deleteUser(id: number) {
  if(confirm('¿Estás seguro de eliminar este usuario?')) {
    try {
      const response = await fetch(`http://localhost:8080/users/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (!response.ok) throw new Error('Error al eliminar')
      allUsers.value = allUsers.value.filter(u => u.id !== id)
    } catch (e) {
      alert('No se pudo eliminar el usuario')
    }
  }
}

// --- 6. CREAR Y EDITAR (MODIFICADO PARA FOTOS) ---
const isModalOpen = ref(false)
const isSaving = ref(false)
const modalError = ref('')
const isEditing = ref(false)
const editingId = ref<number | null>(null)

// Variables para manejo de archivo
const selectedFile = ref<File | null>(null)
const previewUrl = ref<string | null>(null)

const userForm = reactive({
  username: '', lastname: '', email: '', password: '', celphone: '', role: 'general'
})
const selectedInstruments = ref<string[]>([])

// ABRIR PARA CREAR
function openCreateModal() {
  isEditing.value = false
  editingId.value = null
  Object.assign(userForm, { username: '', lastname: '', email: '', password: '', celphone: '', role: 'general' })
  selectedInstruments.value = [] 
  
  // Resetear foto
  selectedFile.value = null
  previewUrl.value = '/logo.png' // Imagen por defecto

  modalError.value = ''
  isModalOpen.value = true
}

// ABRIR PARA EDITAR
async function openEditModal(id: number) {
  isEditing.value = true
  editingId.value = id
  modalError.value = ''
  
  // Resetear foto antes de cargar
  selectedFile.value = null
  previewUrl.value = null // Se llenará con la data del usuario
  
  isModalOpen.value = true

  Object.assign(userForm, { username: 'Cargando...', lastname: '', email: '', password: '', celphone: '', role: '' })
  selectedInstruments.value = []

  try {
    const response = await fetch(`http://localhost:8080/users/${id}`, {
       headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    
    if (response.status === 401) {
       authStore.clearSession()
       return
    }
    if (!response.ok) throw new Error('Error cargando datos')
    
    const data = await response.json()

    userForm.username = data.username
    userForm.lastname = data.lastname
    userForm.email = data.email
    userForm.celphone = data.celphone
    userForm.role = data.role
    userForm.password = '' 
    
    // Cargar foto actual en el preview
    previewUrl.value = getPhotoUrl(data.profile_picture_url)

    if (data.secondary_role) {
        selectedInstruments.value = data.secondary_role.split(',').map((s: string) => s.trim())
    }

  } catch (e) {
    console.error(e)
    modalError.value = 'No se pudo cargar la información'
  }
}

function closeModal() { isModalOpen.value = false }

// NUEVO: Manejar selección de archivo
function onFileSelected(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    selectedFile.value = file;
    // Crear URL temporal para ver la foto seleccionada
    previewUrl.value = URL.createObjectURL(file);
  }
}

// GUARDAR (MODIFICADO PARA INCLUIR FILE)
async function saveUser() {
  isSaving.value = true
  modalError.value = ''

  try {
    if (!userForm.username || !userForm.email) throw new Error('Nombre y Email obligatorios')
    if (!isEditing.value && !userForm.password) throw new Error('Contraseña obligatoria para nuevos usuarios')

    const url = isEditing.value 
      ? `http://localhost:8080/users/${editingId.value}`
      : 'http://localhost:8080/users'
    
    const method = isEditing.value ? 'PUT' : 'POST'

    // Usamos FormData para enviar texto + archivo binario
    const formData = new FormData()
    formData.append('username', userForm.username)
    formData.append('lastname', userForm.lastname || '')
    formData.append('email', userForm.email)
    formData.append('celphone', userForm.celphone || '')
    formData.append('role', userForm.role)
    formData.append('secondary_role', selectedInstruments.value.join(','))

    // Password solo si aplica
    if (userForm.password && userForm.password.trim() !== '') {
        formData.append('password', userForm.password)
    }

    // ADJUNTAR LA FOTO SI SE SELECCIONÓ UNA NUEVA
    if (selectedFile.value) {
        formData.append('file', selectedFile.value) // El backend debe esperar "file"
    }

    const response = await fetch(url, {
      method: method,
      headers: { 
        'Authorization': `Bearer ${authStore.token}`
        // IMPORTANTE: No pongas 'Content-Type': 'multipart/form-data' manual, fetch lo hace solo.
      },
      body: formData 
    })

    if (!response.ok) {
       const errData = await response.json().catch(() => ({}))
       throw new Error(errData.error || 'Error al guardar cambios')
    }

    await getUsers() 
    closeModal()
    alert(isEditing.value ? 'Usuario actualizado con éxito' : 'Usuario creado con éxito')
    
  } catch (e: any) {
    modalError.value = e.message
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
              <input v-model="searchQuery" type="text" placeholder="Buscar por nombre o rol..." />
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
                  <th>Rol / Instrumentos</th>
                  <th>Contacto</th>
                  <th class="center-text">Acciones</th>
               </tr>
            </thead>
            <tbody>
               <tr v-for="user in filteredUsers" :key="user.id">
                  <td>
                     <div class="avatar-wrapper">
                        <img :src="user.photo" alt="avatar" class="avatar" />
                     </div>
                  </td>
                  <td>
                     <div class="info-stack">
                        <span class="primary-text">{{ user.username }} {{ user.lastname }}</span>
                        <span class="secondary-text user-id">ID: {{ user.id }}</span>
                     </div>
                  </td>
                  <td>
                     <div class="info-stack">
                        <div class="role-badge-wrapper">
                           <span class="role-badge" :class="getMainRoleBadgeClass(user.role)">
                              {{ user.role }}
                           </span>
                        </div>
                        <span v-if="user.secondary_role" class="secondary-text subrole-text">
                           {{ getRoleEmoji(user.secondary_role) }} {{ user.secondary_role }}
                        </span>
                        <span v-else class="text-muted">-</span>
                     </div>
                  </td>
                  <td>
                     <div class="info-stack">
                        <span class="contact-text email">✉ {{ user.email }}</span>
                        <span class="contact-text phone">📱 {{ user.celphone }}</span>
                     </div>
                  </td>
                  <td class="center-text">
                     <div class="actions-row">
                        <button @click="openEditModal(user.id)" class="btn-action edit">✎</button>
                        <button @click="deleteUser(user.id)" class="btn-action delete">🗑</button>
                     </div>
                  </td>
               </tr>
            </tbody>
         </table>
         <div v-if="!isLoading && filteredUsers.length === 0" class="empty-state">
            No se encontraron usuarios.
         </div>
      </div>

      <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content form-modal">
          <button class="close-btn" @click="closeModal">×</button>
          
          <div class="modal-header">
            <h3>{{ isEditing ? 'Editar Usuario' : 'Nuevo Usuario' }}</h3>
            <p>{{ isEditing ? 'Modifica los datos y foto de perfil' : 'Registra un nuevo miembro' }}</p>
          </div>

          <div v-if="modalError" class="modal-error">{{ modalError }}</div>

          <form @submit.prevent="saveUser" class="user-form">
            
            <div class="photo-upload-section">
                <div class="preview-container">
                    <img :src="previewUrl || '/logo.png'" alt="Previsualización" class="modal-avatar">
                </div>
                <div class="upload-controls">
                    <label class="btn-file">
                        📷 Cambiar Foto
                        <input type="file" accept="image/*" @change="onFileSelected" />
                    </label>
                    <span v-if="selectedFile" class="file-name">{{ selectedFile.name }}</span>
                </div>
            </div>
            <div class="form-row">
              <div class="form-group"><label>Nombre *</label><input v-model="userForm.username" type="text" required /></div>
              <div class="form-group"><label>Apellido</label><input v-model="userForm.lastname" type="text" /></div>
            </div>
            <div class="form-row">
              <div class="form-group"><label>Email *</label><input v-model="userForm.email" type="email" required /></div>
              <div class="form-group"><label>Celular</label><input v-model="userForm.celphone" type="tel" /></div>
            </div>
            <div class="form-row">
               <div class="form-group">
                <label>Contraseña <span v-if="!isEditing">*</span></label>
                <input 
                  v-model="userForm.password" 
                  type="password" 
                  :required="!isEditing" 
                  :placeholder="isEditing ? 'Dejar vacía para no cambiar' : 'Mínimo 6 caracteres'" 
                />
              </div>
              <div class="form-group">
                <label>Permisos *</label>
                <select v-model="userForm.role" required>
                  <option v-for="r in mainRoles" :key="r" :value="r">{{ r }}</option>
                </select>
              </div>
            </div>
            
            <div class="roles-section">
              <label class="roles-label">Roles Musicales (Selecciona varios)</label>
              <div class="checkbox-grid">
                <label v-for="inst in instrumentOptions" :key="inst" class="checkbox-item" :class="{ 'checked': selectedInstruments.includes(inst) }">
                  <input type="checkbox" :value="inst" v-model="selectedInstruments" />
                  <span>{{ inst }}</span>
                </label>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closeModal">Cancelar</button>
              <button type="submit" class="btn-save" :disabled="isSaving">
                {{ isSaving ? 'Guardando...' : (isEditing ? 'Actualizar' : 'Crear') }}
              </button>
            </div>
          </form>
        </div>
      </div>

    </div>
  </main>
</template>

<style scoped>
/* ESTILOS ORIGINALES + NUEVOS PARA LA FOTO EN MODAL */
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
.text-muted { color: #ddd; font-size: 0.8rem; }
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
.user-id { font-family: monospace; color: #aaa; font-size: 0.7rem; }
.modal-content.form-modal { max-width: 600px; width: 95%; max-height: 90vh; overflow-y: auto; }
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
.roles-section { margin-top: 5px; }
.roles-label { font-size: 0.85rem; color: #64748b; font-weight: 600; margin-bottom: 8px; display: block; }
.checkbox-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.checkbox-item { display: flex; align-items: center; gap: 8px; padding: 8px 12px; border: 1px solid #e2e8f0; border-radius: 8px; cursor: pointer; transition: all 0.2s; background-color: #f8fafc; font-size: 0.85rem; color: #475569; }
.checkbox-item:hover { background-color: #f1f5f9; }
.checkbox-item.checked { background-color: #ecfdf5; border-color: #10b981; color: #065f46; font-weight: 600; }
.checkbox-item input { display: none; }

/* ESTILOS NUEVOS PARA LA CARGA DE FOTO */
.photo-upload-section { display: flex; align-items: center; gap: 20px; margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px dashed #eee; }
.preview-container { width: 80px; height: 80px; border-radius: 50%; overflow: hidden; border: 2px solid #eee; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
.modal-avatar { width: 100%; height: 100%; object-fit: cover; }
.upload-controls { display: flex; flex-direction: column; gap: 5px; }
.btn-file { background: #f3f4f6; border: 1px solid #d1d5db; padding: 8px 12px; border-radius: 6px; cursor: pointer; font-size: 0.85rem; color: #374151; font-weight: 600; text-align: center; transition: background 0.2s; }
.btn-file:hover { background: #e5e7eb; }
.btn-file input { display: none; }
.file-name { font-size: 0.75rem; color: #6b7280; font-style: italic; max-width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

@media (max-width: 600px) { .form-row { grid-template-columns: 1fr; gap: 10px; } .full-width { grid-column: span 1; } }
</style>