<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// --- 1. INTERFAZ ---
interface User {
  id: number;
  username: string;
  lastname: string;
  email: string;
  celphone: string;
  role: string;
  secondary_role: string;
  photo: string;
  status: string;
}

// --- 2. ESTADO ---
const user = ref<User | null>(null)
const isLoading = ref(true)
const errorMsg = ref('')

// --- HELPER PARA CORREGIR URL DE FOTO ---
function getPhotoUrl(path: string) {
  if (!path) return '/logo.png'; // Si no hay foto, usa el logo
  if (path.startsWith('http')) return path; // Si ya es absoluta (ej: Google), déjala
  // Si es relativa (/files/...), pégale el dominio del backend
  return `http://localhost:8080${path}`;
}

// --- 3. GET ME ---
async function getMe() {
  isLoading.value = true
  errorMsg.value = ''
  try {
    const response = await fetch('http://localhost:8080/users/me', {
      method: 'GET',
      headers: { 
        'Content-Type': 'application/json', 
        'Authorization': `Bearer ${authStore.token}` 
      }
    })
    if (!response.ok) throw new Error('Error al cargar tu perfil')
    const data = await response.json()

    user.value = {
      id: data.id,
      username: data.username,
      lastname: data.lastname || '',
      email: data.email,
      celphone: data.celphone || 'No registrado',
      role: data.role || 'general',
      secondary_role: data.secondary_role || '',
      
      // AQUÍ USAMOS LA FUNCIÓN CORRECTORA 👇
      photo: getPhotoUrl(data.profile_picture_url),
      
      status: data.status || 'active'
    }
  } catch (e: any) {
    console.error(e)
    errorMsg.value = 'No se pudo cargar la información del perfil.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => { getMe() })

// --- 4. HELPERS VISUALES ---
function getMainRoleBadgeClass(role: string) {
  return (role.toLowerCase() === 'admin' || role.toLowerCase() === 'administrador') ? 'admin' : 'general'
}

// --- 5. LÓGICA EDICIÓN ---
const isEditModalOpen = ref(false)
const isSaving = ref(false)
const editError = ref('')
const selectedInstruments = ref<string[]>([])
const instrumentOptions = ['Cantante', 'Guitarrista', 'Baterista', 'Pianista', 'Bajista', 'Sonidista', 'Líder', 'Producción']
const selectedPhotoFile = ref<File | null>(null)
const photoPreviewUrl = ref<string | null>(null)

const editForm = reactive({
  username: '', lastname: '', email: '', password: '', celphone: '', secondary_role: ''
})

function openEditModal() {
  if (!user.value) return
  
  editForm.username = user.value.username
  editForm.lastname = user.value.lastname
  editForm.email = user.value.email
  editForm.celphone = user.value.celphone === 'No registrado' ? '' : user.value.celphone
  
  // IMPORTANTE: Resetear la contraseña a vacío siempre que abres el modal
  editForm.password = '' 

  // --- LÓGICA CORREGIDA PARA CHECKBOXES ---
  // 1. Obtenemos los roles del usuario en minúsculas y sin espacios
  const userRoles = user.value.secondary_role 
    ? user.value.secondary_role.split(',').map(r => r.trim().toLowerCase()) 
    : []

  // 2. Filtramos las opciones disponibles. Si la opción (ej: "Baterista") 
  // coincide con alguno de los roles del usuario (ej: "baterista"), la seleccionamos.
  selectedInstruments.value = instrumentOptions.filter(option => 
    userRoles.includes(option.toLowerCase())
  )

  selectedPhotoFile.value = null
  photoPreviewUrl.value = null 

  editError.value = ''
  isEditModalOpen.value = true
}

function closeEditModal() { isEditModalOpen.value = false }

function handlePhotoChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    selectedPhotoFile.value = file
    photoPreviewUrl.value = URL.createObjectURL(file)
  }
}

async function saveProfile() {
  if (!user.value) return
  isSaving.value = true
  editError.value = ''

  try {
    const formData = new FormData()
    formData.append('username', editForm.username)
    formData.append('lastname', editForm.lastname)
    formData.append('email', editForm.email)
    formData.append('celphone', editForm.celphone)
    formData.append('role', user.value.role) 
    formData.append('secondary_role', selectedInstruments.value.join(','))

    // --- LÓGICA CORREGIDA PASSWORD ---
    // Solo enviamos 'password' si NO está vacío y NO son solo espacios
    if (editForm.password && editForm.password.trim().length > 0) {
        formData.append('password', editForm.password)
    }

    if (selectedPhotoFile.value) {
      formData.append('file', selectedPhotoFile.value)
    }

    const response = await fetch(`http://localhost:8080/users/${user.value.id}`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: formData
    })

    if (!response.ok) {
       const data = await response.json().catch(()=>({}))
       throw new Error(data.error || 'Error al actualizar perfil')
    }

    await getMe() 
    closeEditModal()
    // Opcional: Limpiar el formulario de password por seguridad tras guardar
    editForm.password = '' 
    alert('Perfil actualizado correctamente')

  } catch (e: any) {
    editError.value = e.message
  } finally {
    isSaving.value = false
  }
}

function getRoleEmoji(secRole: string) {
  if (!secRole) return '👤'
  const r = secRole.toLowerCase()
  if (r.includes('cantante')) return '🎤'
  if (r.includes('guitarra')) return '🎸'
  if (r.includes('batería')) return '🥁'
  if (r.includes('piano')) return '🎹'
  if (r.includes('bajo')) return '🎸'
  if (r.includes('sonido')) return '🎚'
  return '🎵'
}
</script>

<template>
  <main class="profile-container">
    
    <div v-if="isLoading" class="loading-state">Cargando perfil...</div>
    <div v-if="errorMsg" class="error-banner">⚠ {{ errorMsg }}</div>

    <div v-if="user && !isLoading" class="profile-card fade-in">
      <div class="card-header-bg"></div>

      <div class="card-content">
        <div class="profile-header">
          <div class="avatar-container">
            <img :src="user.photo" alt="Foto de perfil" class="big-avatar" />
          </div>
          <h1 class="user-fullname">{{ user.username }} {{ user.lastname }}</h1>
          <p class="user-email">✉ {{ user.email }} | ☏ {{ user.celphone }}</p>
          
          <div class="roles-container">
             <span class="role-badge" :class="getMainRoleBadgeClass(user.role)">
                {{ user.secondary_role }}
             </span>
          </div>
        </div>

        <div class="card-footer">
           <button class="btn-edit" @click="openEditModal">Editar Perfil</button>
        </div>
      </div>
    </div>

    <div v-if="isEditModalOpen" class="modal-overlay" @click.self="closeEditModal">
      <div class="modal-content form-modal">
        <button class="close-btn" @click="closeEditModal">×</button>
        
        <div class="modal-header">
          <h3>Editar mi Perfil</h3>
          <p>Actualiza tu información personal</p>
        </div>

        <div v-if="editError" class="modal-error">{{ editError }}</div>

        <form @submit.prevent="saveProfile" class="user-form">
          
          <div class="photo-upload-section">
             <div class="preview-circle">
               <img :src="photoPreviewUrl || (user ? user.photo : '/logo.png')" alt="Preview" />
             </div>
             <label class="btn-upload">
               📷 Cambiar Foto
               <input type="file" @change="handlePhotoChange" accept="image/*" hidden />
             </label>
          </div>

          <div class="form-row">
            <div class="form-group"><label>Nombre</label><input v-model="editForm.username" type="text" required /></div>
            <div class="form-group"><label>Apellido</label><input v-model="editForm.lastname" type="text" /></div>
          </div>

          <div class="form-row">
             <div class="form-group"><label>Email</label><input v-model="editForm.email" type="email" required /></div>
             <div class="form-group"><label>Celular</label><input v-model="editForm.celphone" type="tel" /></div>
          </div>
          
          <div class="form-row">
             <div class="form-group full-width">
               <label>Nueva Contraseña (Opcional)</label>
               <input v-model="editForm.password" type="password" placeholder="Dejar vacío para no cambiar" />
             </div>
          </div>

          <div class="roles-section">
              <label class="roles-label">Mis Roles Musicales</label>
              <div class="checkbox-grid">
                <label v-for="inst in instrumentOptions" :key="inst" class="checkbox-item" :class="{ 'checked': selectedInstruments.includes(inst) }">
                  <input type="checkbox" :value="inst" v-model="selectedInstruments" />
                  <span>{{ inst }}</span>
                </label>
              </div>
          </div>

          <div class="form-actions">
             <button type="button" class="btn-cancel" @click="closeEditModal">Cancelar</button>
             <button type="submit" class="btn-save" :disabled="isSaving">
                {{ isSaving ? 'Guardando...' : 'Guardar Cambios' }}
             </button>
          </div>
        </form>
      </div>
    </div>

  </main>
</template>

<style scoped>
/* CONTENEDOR PRINCIPAL */
.profile-container { display: flex; justify-content: center; align-items: center; min-height: 80vh; padding: 40px 20px; flex-direction: column; }
.profile-card { background: white; width: 100%; max-width: 500px; border-radius: 24px; box-shadow: 0 20px 40px rgba(0,0,0,0.08); overflow: hidden; position: relative; border: 1px solid #f0f0f0; }
.card-header-bg { height: 100px; background: linear-gradient(135deg, var(--color-secundary, #10b981), #34d399); width: 100%; }
.card-content { padding: 0 40px 40px 40px; margin-top: -50px; display: flex; flex-direction: column; align-items: center; }
.avatar-container { position: relative; margin-bottom: 15px; }
.big-avatar { width: 100px; height: 100px; border-radius: 50%; object-fit: cover; border: 4px solid white; box-shadow: 0 5px 15px rgba(0,0,0,0.1); background-color: #f9fafb; }
.profile-header { text-align: center; width: 100%; }
.user-fullname { margin: 0; font-size: 1.5rem; color: #1e293b; font-weight: 700; }
.user-email { margin: 5px 0 15px 0; color: #64748b; font-size: 0.95rem; }
.roles-container { margin-bottom: 15px; display: flex; flex-direction: column; align-items: center; gap: 5px; }
.role-badge { padding: 5px 15px; border-radius: 20px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; }
.role-badge.admin { background-color: #f3e8ff; color: #7e22ce; border: 1px solid #d8b4fe; }
.role-badge.general { background-color: #eff6ff; color: #1e40af; border: 1px solid #dbeafe; }
.secondary-roles-text { font-size: 0.85rem; color: #555; font-style: italic; }
.card-footer { margin-top: 30px; width: 100%; }
.btn-edit { width: 100%; padding: 12px; border-radius: 12px; border: 1px solid #e2e8f0; background-color: white; color: #475569; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.btn-edit:hover { background-color: #f8fafc; border-color: #cbd5e1; color: #1e293b; }
.loading-state { color: #94a3b8; font-style: italic; }
.error-banner { background-color: #fee2e2; color: #ef4444; padding: 10px; border-radius: 8px; margin-bottom: 20px; }
.fade-in { animation: fadeIn 0.5s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

/* ESTILOS MODAL Y FORMULARIO */
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.modal-content.form-modal { background: white; padding: 30px; border-radius: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); position: relative; animation: slideUp 0.3s ease-out; width: 95%; max-width: 500px; max-height: 90vh; overflow-y: auto; }
.close-btn { position: absolute; top: 15px; right: 20px; background: none; border: none; font-size: 2rem; color: #999; cursor: pointer; }
.modal-header { text-align: center; margin-bottom: 20px; }
.modal-header h3 { margin: 0; color: var(--color-secundary, #2c3e50); }
.modal-header p { margin: 5px 0 0; color: #666; font-size: 0.9rem; }
.modal-error { color: #ef4444; background-color: #fef2f2; padding: 10px; border-radius: 8px; font-size: 0.9rem; border: 1px solid #fee2e2; margin-bottom: 10px; text-align: center; }
.user-form { display: flex; flex-direction: column; gap: 15px; margin-top: 10px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.full-width { grid-column: span 2; } 
.form-group { display: flex; flex-direction: column; gap: 6px; width: 100%;} 
.form-group label { font-size: 0.85rem; font-weight: 600; color: #4b5563; }
.form-group input, .form-group select { padding: 10px; border-radius: 8px; border: 1px solid #d1d5db; font-size: 0.95rem; outline: none; background-color: #f9fafb; font-family: inherit; width: 100%; box-sizing: border-box;}
.form-group input:focus, .form-group select:focus { border-color: var(--color-secundary, #2c3e50); background-color: white; box-shadow: 0 0 0 3px rgba(44, 62, 80, 0.1); }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 15px; border-top: 1px solid #eee; padding-top: 20px; }
.btn-cancel { background: white; border: 1px solid #d1d5db; color: #4b5563; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.btn-save { background: var(--color-secundary, #2c3e50); border: none; color: white; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.btn-save:disabled { opacity: 0.7; cursor: not-allowed; }

/* --- ESTILOS SUBIDA DE FOTO --- */
.photo-upload-section { display: flex; flex-direction: column; align-items: center; margin-bottom: 20px; gap: 10px; }
.preview-circle { width: 100px; height: 100px; border-radius: 50%; overflow: hidden; border: 3px solid #f3f4f6; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.preview-circle img { width: 100%; height: 100%; object-fit: cover; }
.btn-upload { font-size: 0.85rem; color: var(--color-secundary, #10b981); cursor: pointer; font-weight: 600; padding: 5px 10px; border: 1px solid var(--color-secundary, #10b981); border-radius: 20px; transition: all 0.2s; }
.btn-upload:hover { background-color: #ecfdf5; }

/* Checklist */
.roles-section { margin-top: 5px; }
.roles-label { font-size: 0.85rem; color: #64748b; font-weight: 600; margin-bottom: 8px; display: block; }
.checkbox-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.checkbox-item { display: flex; align-items: center; gap: 8px; padding: 8px 12px; border: 1px solid #e2e8f0; border-radius: 8px; cursor: pointer; transition: all 0.2s; background-color: #f8fafc; font-size: 0.85rem; color: #475569; }
.checkbox-item:hover { background-color: #f1f5f9; }
.checkbox-item.checked { background-color: #ecfdf5; border-color: #10b981; color: #065f46; font-weight: 600; }
.checkbox-item input { display: none; }

@media (max-width: 600px) {
  .form-row { grid-template-columns: 1fr; gap: 10px; }
  .full-width { grid-column: span 1; }
}
</style>