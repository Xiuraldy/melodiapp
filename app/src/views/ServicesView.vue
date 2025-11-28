<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ref, reactive, onMounted, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { jwtDecode } from 'jwt-decode';

const router = useRouter();
const authStore = useAuthStore();

// --- 0. SEGURIDAD Y DATOS DEL USUARIO ---
const isAdmin = ref(false)
const currentUserId = ref<number | null>(null)

async function checkUserPermissions() {
  const token = authStore.token
  if (!token) return

  try {
    const decoded: any = jwtDecode(token)
    // Obtener ID
    const userIdFromToken = decoded.MapClaims?.user_id || decoded.user_id || (decoded.sub ? Number(decoded.sub) : null)
    
    if (userIdFromToken) {
        currentUserId.value = Number(userIdFromToken)
        // Consultar rol real en BD
        const response = await fetch(`http://localhost:8080/users/${currentUserId.value}`, {
            headers: { 'Authorization': `Bearer ${authStore.token}` }
        })
        
        if (response.ok) {
            const userData = await response.json()
            const role = userData.role || userData.Role || 'general'
            isAdmin.value = role.toLowerCase() === 'admin' || role.toLowerCase() === 'administrador'
        }
    }
  } catch (e) {
    console.error("Error verificando permisos:", e)
  }
}

// --- 1. ESTADO ---
const programations = ref<any[]>([])
const availableSongs = ref<any[]>([]) 
const availableUsers = ref<any[]>([]) 
const isLoading = ref(true)

const AVAILABLE_OUTFITS = [
  { id: 1, name: "Blanco y Negro", filename: "1 Blanco y Negro.pdf" },
  { id: 2, name: "Azul Celeste", filename: "2 Azul Celeste.pdf" },
  { id: 3, name: "Beige y Blanco", filename: "3 Beige y Blanco.pdf" },
  { id: 4, name: "Militar", filename: "4 Militar.pdf" },
  { id: 5, name: "Lila", filename: "5 Lila.pdf" },
  { id: 6, name: "Denim", filename: "6 Denim.pdf" },
  { id: 7, name: "Mostaza", filename: "7 Mostaza.pdf" },
  { id: 8, name: "Palo de Rosa", filename: "8 Palo de Rosa.pdf" },
  { id: 9, name: "Azul Noche", filename: "9 Azul Noche.pdf" },
  { id: 10, name: "Rojo", filename: "10 Rojo.pdf" },
  { id: 11, name: "Vino Tinto", filename: "11 Vino Tinto.pdf" },
];

const isAddOutfitsModalOpen = ref(false)
const selectedOutfitIds = ref<number[]>([])
const isAddingOutfits = ref(false)
const targetServiceId = ref<number | null>(null) 

// --- 2. GET SERVICES ---
async function getServices() {
  isLoading.value = true
  try {
    const response = await fetch('http://localhost:8080/services', {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` }
    })

    if (!response.ok) throw new Error('Error al cargar servicios')
    const data = await response.json()

    programations.value = data.map((service: any) => {
      const dateObj = formatDateForCard(service.start_time)
      return {
        id: service.id,
        day: dateObj.day,
        month: dateObj.month,
        event: service.name,
        songsList: service.songs || [],
        usersList: service.users || [], 
        outfitsList: service.outfits || []
      }
    })
  } catch (e) { console.error(e) } 
  finally { isLoading.value = false }
}

// --- COMPUTED: FILTRO ---
const visibleProgramations = computed(() => {
  if (isAdmin.value) return programations.value
  return programations.value.filter(service => {
    if (!service.usersList || service.usersList.length === 0) return false
    return service.usersList.some((u: any) => u.id === currentUserId.value)
  })
})

// --- HELPER: MOSTRAR BOTONES DE RESPUESTA ---
function shouldShowResponseButtons(service: any) {
    // Si es admin, no responde (gestiona)
    if (isAdmin.value) return false 
    if (!currentUserId.value || !service.usersList) return false
    
    // Buscar mi participación
    const me = service.usersList.find((u: any) => u.id === currentUserId.value)
    
    // Mostrar si existe Y está pendiente
    // (Asumiendo que status vacío o 'pending' significa pendiente)
    return me && (me.status === 'pending' || !me.status)
}

// --- LÓGICA RESPONDER A SERVICIO ---
async function respondToService(serviceId: number, status: 'accepted' | 'rejected', event: Event) {
  event.stopPropagation() 
  if (!currentUserId.value) return

  const actionText = status === 'accepted' ? 'Aceptar' : 'Declinar'
  if(!confirm(`¿${actionText} participación en este servicio?`)) return

  try {
      // PATCH http://localhost:8080/services/:id/users/:userId/status
      const url = `http://localhost:8080/services/${serviceId}/users/${currentUserId.value}/status`
      
      const response = await fetch(url, {
          method: 'PATCH', 
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${authStore.token}` 
          },
          body: JSON.stringify({ status: status })
      })

      if (!response.ok) throw new Error('Error al actualizar estado')
      
      alert(`Servicio ${status === 'accepted' ? 'Aceptado' : 'Declinado'}`)
      await getServices() // Recargar UI

  } catch (e) {
      console.error(e)
      alert('Error al enviar respuesta')
  }
}

// --- 3. GET SONGS & USERS ---
async function getAvailableSongs() {
  try {
    const response = await fetch('http://localhost:8080/songs', {
      method: 'GET', headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (response.ok) availableSongs.value = await response.json()
  } catch (e) { console.error(e) }
}

async function getAvailableUsers() {
  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET', headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (response.ok) availableUsers.value = await response.json()
  } catch (e) { console.error(e) }
}

function formatDateForCard(isoString: string) {
  const date = new Date(isoString)
  const month = date.toLocaleString('es-ES', { month: 'short' }).toUpperCase().replace('.', '')
  const day = date.getDate().toString()
  return { day, month }
}

onMounted(async () => {
  await checkUserPermissions() 
  await getServices()
  if (isAdmin.value) {
    getAvailableSongs()
    getAvailableUsers()
  }
})

function redirectProgramation(id: number) {
  router.push({ name: 'service', params: { id: id } })
}

// --- 4. CREAR SERVICIO ---
const isCreateModalOpen = ref(false)
const isSaving = ref(false)
const createError = ref('')
const serviceOptions = ['Servicio Domingo', 'Culto de Jóvenes', 'Culto de Mujeres', 'Culto de Caballeros', 'Culto de Camino Kids', 'Ayuno', 'Otro...']
const formState = reactive({ nameSelection: '', customName: '', startDate: '', endDate: '' })

function openCreateModal() {
  Object.assign(formState, { nameSelection: '', customName: '', startDate: '', endDate: '' })
  createError.value = ''
  isCreateModalOpen.value = true
}
function closeCreateModal() { isCreateModalOpen.value = false }

async function saveNewService() {
  isSaving.value = true
  createError.value = ''
  try {
    let finalName = formState.nameSelection
    if (finalName === 'Otro...') {
      if (!formState.customName) throw new Error('Escribe el nombre del evento')
      finalName = formState.customName
    }
    if (!formState.startDate || !formState.endDate) throw new Error('Faltan fechas')

    const payload = {
      name: finalName,
      start_time: `${formState.startDate}:00Z`, 
      end_time: `${formState.endDate}:00Z`
    }
    const response = await fetch('http://localhost:8080/services', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
      body: JSON.stringify(payload)
    })
    if (!response.ok) throw new Error('Error al crear')
    await getServices() 
    closeCreateModal()
  } catch (e: any) { createError.value = e.message } 
  finally { isSaving.value = false }
}

// --- 5. ELIMINAR SERVICIO ---
async function deleteService(id: number, event: Event) {
  event.stopPropagation() 
  if(confirm('¿Estás seguro de eliminar este servicio?')) {
    try {
      const response = await fetch(`http://localhost:8080/services/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (!response.ok) throw new Error('No se pudo eliminar')
      await getServices()
    } catch (e) { console.error(e); alert('Error al eliminar') }
  }
}

// --- 6. AGREGAR CANCIONES ---
const isAddSongsModalOpen = ref(false)
const selectedSongIds = ref<number[]>([]) 
const isAddingSongs = ref(false)
const songSearchQuery = ref('')

const filteredAvailableSongs = computed(() => {
  if (!songSearchQuery.value) return availableSongs.value
  const term = songSearchQuery.value.toLowerCase()
  return availableSongs.value.filter(song => 
    song.name.toLowerCase().includes(term) || (song.author && song.author.toLowerCase().includes(term))
  )
})

function openAddSongsModal(serviceId: number, event: Event) {
  event.stopPropagation()
  targetServiceId.value = serviceId
  const service = programations.value.find(p => p.id === serviceId)
  if (service && service.songsList) { selectedSongIds.value = service.songsList.map((s: any) => s.id) } 
  else { selectedSongIds.value = [] }
  songSearchQuery.value = ''
  isAddSongsModalOpen.value = true
}
function closeAddSongsModal() { isAddSongsModalOpen.value = false }

async function saveAddedSongs() {
  isAddingSongs.value = true
  try {
    const payload = { song_ids: selectedSongIds.value }
    const response = await fetch(`http://localhost:8080/services/${targetServiceId.value}/songs`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
      body: JSON.stringify(payload)
    })
    if (!response.ok) throw new Error('Error')
    alert('Canciones actualizadas')
    closeAddSongsModal()
    await getServices() 
  } catch (e) { alert('Error al guardar') } 
  finally { isAddingSongs.value = false }
}

// --- 7. AGREGAR EQUIPO ---
const isAddUsersModalOpen = ref(false)
const selectedUserIds = ref<number[]>([]) 
const isAddingUsers = ref(false)
const userSearchQuery = ref('')

const filteredAvailableUsers = computed(() => {
  if (!userSearchQuery.value) return availableUsers.value
  const term = userSearchQuery.value.toLowerCase()
  return availableUsers.value.filter(user => 
    (user.username && user.username.toLowerCase().includes(term)) || 
    (user.lastname && user.lastname.toLowerCase().includes(term))
  )
})

function openAddUsersModal(serviceId: number, event: Event) {
  event.stopPropagation()
  targetServiceId.value = serviceId
  const service = programations.value.find(p => p.id === serviceId)
  if (service && service.usersList) { selectedUserIds.value = service.usersList.map((u: any) => u.id) } 
  else { selectedUserIds.value = [] }
  userSearchQuery.value = ''
  isAddUsersModalOpen.value = true
}
function closeAddUsersModal() { isAddUsersModalOpen.value = false }

async function saveAddedUsers() {
  isAddingUsers.value = true
  try {
    const payload = { user_ids: selectedUserIds.value }
    const response = await fetch(`http://localhost:8080/services/${targetServiceId.value}/users`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
      body: JSON.stringify(payload)
    })
    if (!response.ok) throw new Error('Error')
    alert('Equipo actualizado')
    closeAddUsersModal()
    await getServices() 
  } catch (e) { alert('Error al guardar') } 
  finally { isAddingUsers.value = false }
}

// --- 8. OUTFITS ---
function openAddOutfitsModal(serviceId: number, event: Event) {
  event.stopPropagation()
  targetServiceId.value = serviceId
  const service = programations.value.find(p => p.id === serviceId)
  if (service && service.outfitsList) { selectedOutfitIds.value = service.outfitsList.map((o: any) => o.id || o.outfit_id) } 
  else { selectedOutfitIds.value = [] }
  isAddOutfitsModalOpen.value = true
}
function closeAddOutfitsModal() { isAddOutfitsModalOpen.value = false }

async function saveAddedOutfits() {
  isAddingOutfits.value = true
  try {
    const payload = { outfit_ids: selectedOutfitIds.value }
    const response = await fetch(`http://localhost:8080/services/${targetServiceId.value}/outfits`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
      body: JSON.stringify(payload)
    })
    if (!response.ok) throw new Error('Error')
    alert('Paleta asignada')
    closeAddOutfitsModal()
    await getServices() 
  } catch (e) { alert('Error al guardar') } 
  finally { isAddingOutfits.value = false }
}

function getRoleEmoji(role: string) {
  if (!role) return '👤'
  const r = role.toLowerCase()
  if (r.includes('cantante')) return '🎤'
  if (r.includes('guitarra')) return '🎸'
  if (r.includes('batería')) return '🥁'
  if (r.includes('piano')) return '🎹'
  return '🎵'
}
</script>

<template>
  <main>
    <div class="programation-view">
    
      <div class="header-actions">
          <button v-if="isAdmin" class="btn-create" @click="openCreateModal">
            + Nuevo Servicio
          </button>
      </div>

      <div v-if="isLoading" class="loading-state">Cargando servicios...</div>

      <div class="programations-container" v-else>
        <div v-if="visibleProgramations.length === 0" class="empty-msg">
           {{ isAdmin ? 'No hay servicios programados.' : 'No tienes servicios asignados.' }}
        </div>

        <div 
          v-for="prog in visibleProgramations" 
          :key="prog.id" 
          class="card-programation" 
          @click="redirectProgramation(prog.id)"
        >
          <div class="card-header">
            <div class="calendar-box">
              <span class="month">{{ prog.month }}</span>
              <span class="day">{{ prog.day }}</span>
            </div>
            <div class="header-info">
              <h3 class="event-title">{{ prog.event }}</h3>
              
              <div v-if="!isAdmin && currentUserId" class="my-status-badge">
                 <span v-for="u in prog.usersList" :key="u.id">
                    <span v-if="u.id === currentUserId">
                       <span v-if="u.status === 'accepted'" class="accepted">✅ Asistiré</span>
                       <span v-else-if="u.status === 'rejected'" class="rejected">❌ No asistiré</span>
                       <span v-else class="pending">🕒 Pendiente</span>
                    </span>
                 </span>
              </div>

            </div>
            
            <button v-if="isAdmin" class="btn-quick delete-btn" title="Eliminar" @click="deleteService(prog.id, $event)">
               ❌
            </button>
          </div>
          
          <div class="divider"></div>
          
          <div class="card-body">
            <div class="songs-label">Repertorio:</div>
            <ul class="song-list" v-if="prog.songsList && prog.songsList.length > 0">
              <li v-for="song in prog.songsList" :key="song.id">
                <span class="music-icon">🎵</span> {{ song.name }}
              </li>
            </ul>
            <div v-else class="no-songs">Aún no hay canciones asignadas.</div>

            <div v-if="shouldShowResponseButtons(prog)" class="response-actions" @click.stop>
                <button class="btn-response accept" @click="respondToService(prog.id, 'accepted', $event)">Aceptar</button>
                <button class="btn-response decline" @click="respondToService(prog.id, 'rejected', $event)">Declinar</button>
            </div>
          </div>
          
          <div class="card-footer">
            <div class="quick-actions" v-if="isAdmin">
                <button class="btn-quick" title="Canciones" @click="openAddSongsModal(prog.id, $event)">🎵</button>
                <button class="btn-quick" title="Equipo" @click="openAddUsersModal(prog.id, $event)">👤</button>
                <button class="btn-quick" title="Paleta" @click="openAddOutfitsModal(prog.id, $event)">🎨</button>
            </div>
            <div v-else class="quick-actions"></div>
            <span class="arrow">Ver Detalle →</span>
          </div>
        </div>
      </div>

      <div v-if="isCreateModalOpen" class="modal-overlay" @click.self="closeCreateModal">
        <div class="modal-content form-modal">
          <button class="close-btn" @click="closeCreateModal">×</button>
          <div class="modal-header"><h3>Crear Servicio</h3></div>
          <form @submit.prevent="saveNewService" class="user-form">
             <div class="form-row"><div class="form-group full-width"><label>Tipo</label><select v-model="formState.nameSelection" required><option v-for="opt in serviceOptions" :key="opt" :value="opt">{{ opt }}</option></select></div></div>
             <div class="form-row" v-if="formState.nameSelection === 'Otro...'"><div class="form-group full-width"><label>Nombre</label><input v-model="formState.customName" type="text" required /></div></div>
             <div class="form-row"><div class="form-group"><label>Inicio</label><input v-model="formState.startDate" type="datetime-local" required /></div><div class="form-group"><label>Fin</label><input v-model="formState.endDate" type="datetime-local" required /></div></div>
             <div class="form-actions"><button type="button" class="btn-cancel" @click="closeCreateModal">Cancelar</button><button type="submit" class="btn-save" :disabled="isSaving">Crear</button></div>
          </form>
        </div>
      </div>

      <div v-if="isAddSongsModalOpen" class="modal-overlay" @click.self="closeAddSongsModal">
        <div class="modal-content form-modal">
          <div class="modal-header"><h3>Repertorio</h3></div>
          <div class="modal-search-box"><input v-model="songSearchQuery" type="text" placeholder="Buscar..." /></div>
          <div class="song-selection-list">
             <label v-for="song in filteredAvailableSongs" :key="song.id" class="song-checkbox-item" :class="{ 'selected': selectedSongIds.includes(song.id) }"><input type="checkbox" :value="song.id" v-model="selectedSongIds" /><div class="song-info-check"><span class="song-name">{{ song.name }}</span></div></label>
          </div>
          <div class="form-actions"><button class="btn-cancel" @click="closeAddSongsModal">Cancelar</button><button class="btn-save" @click="saveAddedSongs" :disabled="isAddingSongs">Guardar</button></div>
        </div>
      </div>

      <div v-if="isAddUsersModalOpen" class="modal-overlay" @click.self="closeAddUsersModal">
        <div class="modal-content form-modal">
          <div class="modal-header"><h3>Equipo</h3></div>
          <div class="modal-search-box"><input v-model="userSearchQuery" type="text" placeholder="Buscar..." /></div>
          <div class="song-selection-list">
             <label v-for="user in filteredAvailableUsers" :key="user.id" class="song-checkbox-item" :class="{ 'selected': selectedUserIds.includes(user.id) }"><input type="checkbox" :value="user.id" v-model="selectedUserIds" /><div class="song-info-check"><span class="song-name">{{ user.username }}</span></div></label>
          </div>
          <div class="form-actions"><button class="btn-cancel" @click="closeAddUsersModal">Cancelar</button><button class="btn-save" @click="saveAddedUsers" :disabled="isAddingUsers">Guardar</button></div>
        </div>
      </div>

      <div v-if="isAddOutfitsModalOpen" class="modal-overlay" @click.self="closeAddOutfitsModal">
        <div class="modal-content form-modal">
           <div class="modal-header"><h3>Paleta</h3></div>
           <div class="song-selection-list">
              <label v-for="outfit in AVAILABLE_OUTFITS" :key="outfit.id" class="song-checkbox-item" :class="{ 'selected': selectedOutfitIds.includes(outfit.id) }"><input type="checkbox" :value="outfit.id" v-model="selectedOutfitIds" /><span class="song-name">{{ outfit.name }}</span></label>
           </div>
           <div class="form-actions"><button class="btn-cancel" @click="closeAddOutfitsModal">Cancelar</button><button class="btn-save" @click="saveAddedOutfits" :disabled="isAddingOutfits">Guardar</button></div>
        </div>
      </div>

    </div>
  </main>
</template>

<style scoped>
/* ESTILOS IDÉNTICOS + NUEVOS PARA BOTONES RESPUESTA */
.programation-view { padding: 40px; display: flex; flex-direction: column; gap: 20px; }
.header-actions { display: flex; gap: 15px; align-items: center; justify-content: center;}
.btn-create { background-color: #10b981; color: white; border: none; padding: 8px 16px; border-radius: 8px; font-weight: 600; font-size: 0.85rem; cursor: pointer; transition: background 0.2s; display: flex; align-items: center; gap: 5px; }
.btn-create:hover { background-color: #059669; }
.programations-container { display: flex; justify-content: center; flex-wrap: wrap; gap: 50px; margin-top: 20px; }
.loading-state, .empty-msg { width: 100%; text-align: center; color: #999; padding: 40px; font-style: italic; }

.card-programation { background-color: white; width: 320px; border-radius: 16px; box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05); border: 1px solid #f0f0f0; cursor: pointer; transition: all 0.3s ease; overflow: hidden; display: flex; flex-direction: column; }
.card-programation:hover { transform: translateY(-5px); box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1); border-color: var(--color-secundary, #2c3e50); }
.card-header { display: flex; padding: 20px; gap: 15px; align-items: center; background: linear-gradient(to bottom, #ffffff, #fafafa); }
.calendar-box { display: flex; flex-direction: column; align-items: center; justify-content: center; background-color: white; border: 2px solid var(--color-secundary, #2c3e50); border-radius: 12px; min-width: 60px; height: 65px; box-shadow: 3px 3px 0px var(--color-tertiary, #e67e22); }
.calendar-box .month { font-size: 0.7rem; background-color: var(--color-secundary, #2c3e50); color: white; width: 100%; text-align: center; padding: 2px 0; font-weight: 700; }
.calendar-box .day { font-size: 1.8rem; font-weight: 800; color: #333; line-height: 1; margin-top: 2px; }
.header-info { display: flex; flex-direction: column; align-items: flex-start; flex-grow: 1; }
.role-badge { font-size: 0.75rem; background-color: var(--color-tertiary, #e67e22); color: white; padding: 4px 10px; border-radius: 20px; font-weight: 600; margin-bottom: 5px; text-transform: uppercase; letter-spacing: 0.5px; }
.event-title { margin: 0; font-size: 1.1rem; color: #333; font-weight: 700; line-height: 1.2; }
.divider { height: 1px; background: #eee; margin: 0 20px; }
.card-body { padding: 20px; flex-grow: 1; }
.songs-label { font-size: 0.75rem; color: #999; font-weight: 700; margin-bottom: 10px; text-transform: uppercase; letter-spacing: 1px; }
.song-list { list-style: none; padding: 0; margin: 0; }
.song-list li { display: flex; align-items: center; font-size: 0.9rem; color: #555; margin-bottom: 6px; padding: 8px 10px; background-color: #f8fafc; border-radius: 8px; transition: background 0.2s; border: 1px solid transparent; }
.music-icon { margin-right: 10px; font-size: 0.8rem; filter: grayscale(100%); opacity: 0.6; }
.no-songs { color: #bbb; font-size: 0.9rem; font-style: italic; padding: 10px; text-align: center; }
.card-footer { padding: 15px 20px; background-color: #fcfcfc; border-top: 1px solid #eee; display: flex; justify-content: space-between; align-items: center; font-size: 0.85rem; font-weight: 600; color: var(--color-secundary, #2c3e50); }
.quick-actions { display: flex; gap: 8px; }
.btn-quick { background-color: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; padding: 6px 10px; border-radius: 6px; font-size: 1rem; cursor: pointer; transition: all 0.2s; }
.btn-quick:hover { background-color: #dbeafe; transform: translateY(-2px); }
.btn-quick.delete-btn { color: #ef4444; border-color: #fecaca; background-color: #fef2f2; }
.arrow { transition: transform 0.2s; }
.card-programation:hover .arrow { transform: translateX(5px); }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.modal-content.form-modal { background: white; padding: 30px; border-radius: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); position: relative; animation: slideUp 0.3s ease-out; width: fit-content; min-width: 320px; max-width: -webkit-fill-available; height: auto; max-height: 90vh; overflow-y: auto; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
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
.modal-search-box { margin-bottom: 10px; }
.modal-search-box input { width: 100%; padding: 10px; border-radius: 8px; border: 1px solid #dbeafe; outline: none; background-color: #f8fafc; }
.modal-search-box input:focus { border-color: #10b981; background: #fff; }
.song-selection-list { display: flex; flex-direction: column; gap: 8px; margin-top: 5px; max-height: 300px; overflow-y: auto; border: 1px solid #eee; padding: 10px; border-radius: 8px; }
.song-checkbox-item, .checkbox-item { display: flex; align-items: center; gap: 12px; padding: 10px; border: 1px solid #f0f0f0; border-radius: 8px; cursor: pointer; transition: all 0.2s; }
.song-checkbox-item:hover, .checkbox-item:hover { background-color: #f9fafb; }
.song-checkbox-item.selected, .checkbox-item.selected { background-color: #ecfdf5; border-color: #10b981; }
.song-checkbox-item input, .checkbox-item input { width: 18px; height: 18px; accent-color: #10b981; }
.song-info-check, .item-info { display: flex; flex-direction: column; width: 100%;}
.song-name, .item-name { font-weight: 700; font-size: 0.9rem; color: #333; }
.song-author, .item-sub { font-size: 0.8rem; color: #666; }
.mini-avatar { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; border: 1px solid #eee; }

/* --- ESTILOS NUEVOS: BOTONES RESPUESTA --- */
.response-actions {
    display: flex;
    gap: 10px;
    margin-top: 15px;
    border-top: 1px dashed #eee;
    padding-top: 10px;
}

.btn-response {
    flex: 1;
    padding: 8px;
    border-radius: 6px;
    font-weight: 600;
    font-size: 0.85rem;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
}

.btn-response.accept {
    background-color: #dcfce7;
    color: #166534;
    border: 1px solid #bbf7d0;
}
.btn-response.accept:hover { background-color: #bbf7d0; }

.btn-response.decline {
    background-color: #fee2e2;
    color: #991b1b;
    border: 1px solid #fecaca;
}
.btn-response.decline:hover { background-color: #fecaca; }

.my-status-badge {
    margin-top: 5px;
    font-size: 0.75rem;
    font-weight: 600;
}
.my-status-badge .accepted { color: #166534; }
.my-status-badge .rejected { color: #991b1b; }
.my-status-badge .pending { color: #d97706; }

@media (max-width: 600px) {
  .form-row { grid-template-columns: 1fr; gap: 10px; }
  .full-width { grid-column: span 1; }
  .programations-container { justify-content: center; }
  .modal-content.form-modal { width: 95%; }
}
</style>