<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { jwtDecode } from 'jwt-decode' // Necesario para leer el rol

const authStore = useAuthStore()

// --- 0. SEGURIDAD Y PERMISOS ---
const isAdmin = ref(false)
const currentUserId = ref<number | null>(null)

async function checkUserPermissions() {
  const token = authStore.token
  if (!token) return

  try {
    const decoded: any = jwtDecode(token)
    // Intentamos obtener el ID
    const userIdFromToken = decoded.id || decoded.user_id || (decoded.MapClaims ? decoded.MapClaims.user_id : null)
    
    if (userIdFromToken) {
        currentUserId.value = Number(userIdFromToken)
        // Consultamos rol real en BD
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
    console.error("Error verificando permisos", e)
  }
}

// --- 1. INTERFAZ ---
interface Song {
  id: number;
  name: string;
  author: string;
  key: string;
  bpm: number;
  timeSignature: string;
  duration: string;
  structure: string[]; 
  hasSequence: string;
  hasChart: string;
  hasScore: string;
  youtubeUrl: string;
  sequenceUrl: string;
  chartUrl: string;
  scoreUrl: string;
  voiceUrl: string;
  guitarUrl: string;
  pianoUrl: string;
  drumsUrl: string;
  bassUrl: string;
}

const allSongs = ref<Song[]>([])
const isLoading = ref(true)
const errorMsg = ref('')

const musicalKeys = ['C', 'Cm', 'C#', 'Db', 'D', 'Dm', 'D#', 'Eb', 'E', 'Em', 'F', 'Fm', 'F#', 'Gb', 'G', 'Gm', 'G#', 'Ab', 'A', 'Am', 'A#', 'Bb', 'B', 'Bm']
const timeSignatures = ['4/4', '6/8', '3/4', '2/4', '12/8']
const structureOptions = ['Intro', 'Verso', 'Pre-Coro', 'Coro', 'Puente', 'Solo', 'Instrumental', 'Final', 'Outro']

// Helper URLs
function getLocalPath(folder: string, prefix: string, id: number, ext: string) {
  return `/files/${folder}/${prefix}${id}.${ext}`
}

// --- 2. GET SONGS ---
async function getSongs() {
  isLoading.value = true
  errorMsg.value = ''
  try {
    const response = await fetch('http://localhost:8080/songs', {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` }
    })
    if (!response.ok) throw new Error('Error al cargar la biblioteca')
    const data = await response.json()

    allSongs.value = data.map((s: any) => ({
      id: s.id,
      name: s.name,
      author: s.author,
      key: s.song_key,
      bpm: s.bpm,
      timeSignature: s.time_signature,
      duration: s.duration,
      structure: s.structure ? s.structure.split(',') : [],
      hasSequence: s.has_sequence,
      hasChart: s.has_chart,
      hasScore: s.has_score,
      youtubeUrl: s.youtube_url,
      sequenceUrl: s.sequence_url || s.has_sequence || '', 
      chartUrl:    s.chart_url    || s.has_chart    || '',
      scoreUrl:    s.score_url    || s.has_score    || '',
      voiceUrl: s.voice_url || '',
      guitarUrl: s.guitar_url || '',
      pianoUrl: s.piano_url || '',
      drumsUrl: s.drums_url || '',
      bassUrl: s.bass_url || ''
    }))
  } catch (e: any) {
    console.error(e)
    errorMsg.value = 'No se pudieron cargar las canciones.'
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  await checkUserPermissions() // 1. Verificar Rol
  await getSongs()             // 2. Cargar canciones
})

// --- 3. FILTROS ---
const searchQuery = ref('')
const filteredSongs = computed(() => {
  if (!searchQuery.value) return allSongs.value
  const term = searchQuery.value.toLowerCase()
  return allSongs.value.filter(song => 
    song.name.toLowerCase().includes(term) || song.author.toLowerCase().includes(term)
  )
})

// --- 4. MODAL VISUALIZACIÓN ---
const isModalOpen = ref(false)
const selectedSongTitle = ref('')
const selectedSongStructure = ref<string[]>([])

function openStructureModal(song: Song) {
  if (!song.structure || song.structure.length === 0) return
  selectedSongTitle.value = song.name
  selectedSongStructure.value = song.structure
  isModalOpen.value = true
}
function closeModal() { isModalOpen.value = false }

// --- 5. ELIMINAR (SOLO ADMIN) ---
async function deleteSong(id: number) {
  if(confirm('¿Estás seguro de eliminar esta canción permanentemente?')) {
    try {
      const response = await fetch(`http://localhost:8080/songs/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (!response.ok) throw new Error('Error al eliminar')
      allSongs.value = allSongs.value.filter(s => s.id !== id)
    } catch (e) {
      alert('No se pudo eliminar la canción')
    }
  }
}

// --- 6. CREAR Y EDITAR (SOLO ADMIN) ---
const isCreateModalOpen = ref(false)
const isSaving = ref(false)
const createError = ref('')
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const structureBuilder = ref<string[]>([])

const newSongForm = reactive({
  name: '', author: '', song_key: '', bpm: 120, time_signature: '4/4', duration: '',
  youtube_url: '', sequence_url: '', chart_url: '', score_url: '',
  voice_url: '', guitar_url: '', piano_url: '', drums_url: '', bass_url: ''
})

function resetForm() {
  Object.assign(newSongForm, {
    name: '', author: '', song_key: '', bpm: 120, time_signature: '4/4', duration: '', youtube_url: '',
    sequence_url: '', chart_url: '', score_url: '',
    voice_url: '', guitar_url: '', piano_url: '', drums_url: '', bass_url: ''
  })
  structureBuilder.value = [] 
  createError.value = ''
}

// ABRIR PARA CREAR
function openCreateModal() {
  isEditing.value = false
  editingId.value = null
  resetForm()
  isCreateModalOpen.value = true
}

// ABRIR PARA EDITAR
async function editSong(id: number) {
  isEditing.value = true
  editingId.value = id
  resetForm()
  isCreateModalOpen.value = true 

  try {
      const song = allSongs.value.find(s => s.id === id)
      if (song) {
          newSongForm.name = song.name
          newSongForm.author = song.author
          newSongForm.song_key = song.key
          newSongForm.bpm = song.bpm
          newSongForm.time_signature = song.timeSignature
          newSongForm.duration = song.duration
          newSongForm.youtube_url = song.youtubeUrl
          newSongForm.sequence_url = song.sequenceUrl
          newSongForm.chart_url = song.chartUrl
          newSongForm.score_url = song.scoreUrl
          newSongForm.voice_url = song.voiceUrl
          newSongForm.guitar_url = song.guitarUrl
          newSongForm.piano_url = song.pianoUrl
          newSongForm.drums_url = song.drumsUrl
          newSongForm.bass_url = song.bassUrl
          structureBuilder.value = [...song.structure]
      }
  } catch (e) { console.error("Error cargando datos", e) }
}

function closeCreateModal() { isCreateModalOpen.value = false }

function addStructurePart(part: string) { structureBuilder.value.push(part) }
function removeStructurePart(index: number) { structureBuilder.value.splice(index, 1) }

// Manejo de Archivos (Nombre)
function handleFileName(event: Event, fieldKey: keyof typeof newSongForm) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    // @ts-ignore
    newSongForm[fieldKey] = target.files[0].name 
  }
}
function handleResourceFile(event: Event, fieldKey: 'sequence_url' | 'chart_url' | 'score_url') {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    newSongForm[fieldKey] = target.files[0].name
  }
}

async function saveNewSong() {
  isSaving.value = true
  createError.value = ''
  try {
    if (!newSongForm.name || !newSongForm.song_key) throw new Error('Faltan campos obligatorios')

    const payload = {
      name: newSongForm.name,
      author: newSongForm.author,
      song_key: newSongForm.song_key,
      bpm: Number(newSongForm.bpm),
      time_signature: newSongForm.time_signature,
      duration: newSongForm.duration,
      structure: structureBuilder.value.join(','), 
      youtube_url: newSongForm.youtube_url,
      sequence_url: newSongForm.sequence_url,
      chart_url: newSongForm.chart_url,
      score_url: newSongForm.score_url,
      voice_url: newSongForm.voice_url,
      guitar_url: newSongForm.guitar_url,
      piano_url: newSongForm.piano_url,
      drums_url: newSongForm.drums_url,
      bass_url: newSongForm.bass_url,
      has_sequence: newSongForm.sequence_url,
      has_chart: newSongForm.chart_url,
      has_score: newSongForm.score_url
    }

    const url = isEditing.value ? `http://localhost:8080/songs/${editingId.value}` : 'http://localhost:8080/songs'
    const method = isEditing.value ? 'PUT' : 'POST'

    const response = await fetch(url, {
      method: method,
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` }, 
      body: JSON.stringify(payload)
    })

    if (!response.ok) throw new Error('Error al guardar')
    
    await getSongs()
    closeCreateModal()
    alert(isEditing.value ? 'Canción actualizada' : 'Canción creada')
  } catch (e: any) { createError.value = e.message } 
  finally { isSaving.value = false }
}
</script>

<template>
  <main class="view-container">
    <div class="view-header">
      <div><h1>Biblioteca</h1><p class="subtitle">Repertorio general</p></div>
      <div class="header-actions">
        <div class="search-bar-container"><span class="search-icon">🔍</span><input v-model="searchQuery" type="text" placeholder="Buscar canción..." class="search-input" /></div>
        <button v-if="isAdmin" class="btn-create" @click="openCreateModal">+ Nueva Canción</button>
      </div>
    </div>

    <div v-if="errorMsg" class="error-banner">⚠ {{ errorMsg }}</div>

    <div class="card">
      <div class="card-header-flex"><h3>Todas las Canciones ({{ filteredSongs.length }})</h3><span v-if="isLoading" class="loading-text">Cargando...</span></div>
      <table class="custom-table songs-table">
        <thead>
          <tr>
            <th class="center-text">Tono</th>
            <th>Canción</th>
            <th>Recursos</th>
            <th>Tracks</th>
            <th class="center-text">Estruc.</th>
            <th class="center-text">BPM</th>
            <th class="center-text">Cifra</th>
            <th class="center-text">Dur.</th>
            <th v-if="isAdmin" class="center-text">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="song in filteredSongs" :key="song.id">
            <td class="center-text"><div class="key-badge">{{ song.key }}</div></td>
            <td><div class="info-stack"><span class="primary-text">{{ song.name }}</span><span class="secondary-text">{{ song.author }}</span></div></td>
            <td><div class="resource-actions"><a v-if="song.youtubeUrl" :href="song.youtubeUrl" target="_blank" class="btn-icon youtube">▶</a><a v-if="song.sequenceUrl" :href="getLocalPath('sequences', 'sequence', song.id, 'mp3')" target="_blank" class="btn-icon sequence" title="Secuencia">🎧</a><a v-if="song.chartUrl" :href="getLocalPath('charts', 'chart', song.id, 'pdf')" target="_blank" class="btn-icon chart" title="Cifrado">📄</a><a v-if="song.scoreUrl" :href="getLocalPath('scores', 'score', song.id, 'pdf')" target="_blank" class="btn-icon score" title="Partitura">🎼</a></div></td>
            <td><div class="resource-actions"><a v-if="song.voiceUrl" :href="getLocalPath('tracks', 'voice', song.id, 'mp3')" target="_blank" class="btn-icon voice" title="Voz">🎤</a><a v-if="song.pianoUrl" :href="getLocalPath('tracks', 'piano', song.id, 'mp3')" target="_blank" class="btn-icon piano" title="Piano">🎹</a><a v-if="song.guitarUrl" :href="getLocalPath('tracks', 'guitar', song.id, 'mp3')" target="_blank" class="btn-icon guitar" title="Guitarra">🎸</a><a v-if="song.bassUrl" :href="getLocalPath('tracks', 'bass', song.id, 'mp3')" target="_blank" class="btn-icon bass" title="Bajo">🎸</a><a v-if="song.drumsUrl" :href="getLocalPath('tracks', 'drums', song.id, 'mp3')" target="_blank" class="btn-icon drums" title="Batería">🥁</a></div></td>
            <td class="center-text"><button v-if="song.structure && song.structure.length > 0" @click="openStructureModal(song)" class="btn-structure">Ver</button><span v-else class="text-muted">-</span></td>
            <td class="center-text bpm-cell">{{ song.bpm }}</td>
            <td class="center-text time-sig-cell">{{ song.timeSignature }}</td>
            <td class="center-text duration-cell">{{ song.duration }}</td>
            
            <td v-if="isAdmin" class="center-text">
                <div class="actions-row">
                    <button @click="editSong(song.id)" class="btn-action edit">✎</button>
                    <button @click="deleteSong(song.id)" class="btn-action delete">🗑</button>
                </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content"><button class="close-btn" @click="closeModal">×</button><div class="modal-header"><h3>Estructura</h3><p>{{ selectedSongTitle }}</p></div><div class="structure-timeline"><div v-for="(part, index) in selectedSongStructure" :key="index" class="structure-item"><div class="dot"></div><div class="part-name">{{ part }}</div></div></div></div>
    </div>

    <div v-if="isCreateModalOpen && isAdmin" class="modal-overlay" @click.self="closeCreateModal">
      <div class="modal-content form-modal">
        <button class="close-btn" @click="closeCreateModal">×</button>
        <div class="modal-header"><h3>{{ isEditing ? 'Editar Canción' : 'Nueva Canción' }}</h3><p>{{ isEditing ? 'Modifica los detalles' : 'Añade un tema al repertorio' }}</p></div>
        <div v-if="createError" class="modal-error">{{ createError }}</div>
        <form @submit.prevent="saveNewSong" class="user-form">
          <div class="form-row"><div class="form-group"><label>Nombre *</label><input v-model="newSongForm.name" type="text" required /></div><div class="form-group"><label>Artista</label><input v-model="newSongForm.author" type="text" /></div></div>
          <div class="form-row"><div class="form-group"><label>Tono *</label><select v-model="newSongForm.song_key" required><option value="" disabled selected>...</option><option v-for="k in musicalKeys" :key="k" :value="k">{{ k }}</option></select></div><div class="form-group"><label>BPM</label><input v-model="newSongForm.bpm" type="number" /></div><div class="form-group"><label>Compás</label><select v-model="newSongForm.time_signature"><option v-for="ts in timeSignatures" :key="ts" :value="ts">{{ ts }}</option></select></div><div class="form-group"><label>Duración</label><input v-model="newSongForm.duration" type="text" /></div></div>
          <div class="structure-builder-section"><label class="form-label">Estructura</label><div class="structure-buttons"><button type="button" v-for="part in structureOptions" :key="part" class="badge-btn" @click="addStructurePart(part)">+ {{ part }}</button></div><div class="structure-preview"><div v-if="structureBuilder.length === 0" class="preview-placeholder">Selecciona partes...</div><div v-else class="preview-chips"><div v-for="(part, index) in structureBuilder" :key="index" class="struct-chip">{{ part }} <span class="remove-x" @click="removeStructurePart(index)">×</span></div></div></div></div>
          <div class="form-row"><div class="form-group full-width"><label>Youtube URL</label><input v-model="newSongForm.youtube_url" type="url" /></div></div>
          <div class="form-section-title">Archivos</div>
          <div class="form-row"><div class="form-group"><label>Secuencia</label><input type="file" @change="handleResourceFile($event, 'sequence_url')" /></div><div class="form-group"><label>Cifrado</label><input type="file" @change="handleResourceFile($event, 'chart_url')" /></div><div class="form-group"><label>Partitura</label><input type="file" @change="handleResourceFile($event, 'score_url')" /></div></div>
          <div class="form-section-title">Tracks</div>
          <div class="form-row"><div class="form-group"><label>Voz</label><input type="file" @change="handleFileName($event, 'voice_url')" /></div><div class="form-group"><label>Piano</label><input type="file" @change="handleFileName($event, 'piano_url')" /></div></div>
          <div class="form-row"><div class="form-group"><label>Guitarra</label><input type="file" @change="handleFileName($event, 'guitar_url')" /></div><div class="form-group"><label>Bajo</label><input type="file" @change="handleFileName($event, 'bass_url')" /></div></div>
          <div class="form-row"><div class="form-group"><label>Batería</label><input type="file" @change="handleFileName($event, 'drums_url')" /></div></div>
          <div class="form-actions"><button type="button" class="btn-cancel" @click="closeCreateModal">Cancelar</button><button type="submit" class="btn-save" :disabled="isSaving">Guardar</button></div>
        </form>
      </div>
    </div>
  </main>
</template>

<style scoped>
/* ESTILOS (Iguales a tu versión funcional) */
.view-container { padding: 40px; display: flex; flex-direction: column; gap: 20px; }
.view-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; flex-wrap: wrap; gap: 20px;}
h1 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.5rem; }
.subtitle { margin: 0; color: #888; font-size: 0.9rem; }
.header-actions { display: flex; gap: 15px; align-items: center; }
.search-bar-container { position: relative; width: 250px; }
.search-input { width: 100%; padding: 8px 10px 8px 35px; border-radius: 8px; border: 1px solid #e5e7eb; font-size: 0.9rem; outline: none; transition: border 0.2s; }
.search-input:focus { border-color: var(--color-secundary, #2c3e50); }
.search-icon { position: absolute; left: 10px; top: 50%; transform: translateY(-50%); opacity: 0.5; font-size: 0.8rem;}
.btn-create { background-color: #10b981; color: white; border: none; padding: 8px 16px; border-radius: 8px; font-weight: 600; cursor: pointer; font-size: 0.9rem; transition: background 0.2s; }
.btn-create:hover { background-color: #059669; }
.card { background: white; border: 1px solid #eee; padding: 20px; border-radius: 16px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); overflow-x: auto; }
.card-header-flex { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px; margin-bottom: 15px; }
.card h3 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.2rem; }
.custom-table { width: 100%; border-collapse: collapse; }
.custom-table th { text-align: left; font-size: 0.75rem; text-transform: uppercase; color: #999; padding: 10px; font-weight: 600; white-space: nowrap; }
.custom-table td { padding: 12px 8px; border-bottom: 1px dashed #eee; vertical-align: middle; }
.center-text { text-align: center; }
.info-stack { display: flex; flex-direction: column; }
.primary-text { font-weight: 700; color: #333; font-size: 0.95rem; }
.secondary-text { font-size: 0.8rem; color: #888; margin-top: 2px; }
.key-badge { background-color: #f3f4f6; color: #333; font-weight: 800; width: 35px; height: 35px; display: flex; align-items: center; justify-content: center; border-radius: 8px; font-size: 0.9rem; border: 1px solid #e5e7eb; margin: 0 auto; }
.bpm-cell { font-family: monospace; font-weight: 600; color: #555; }
.time-sig-cell { font-family: monospace; color: #666; font-weight: 600; background: #f9fafb; padding: 4px; border-radius: 4px; font-size: 0.85rem;}
.duration-cell { font-family: monospace; color: #333; }
.resource-actions { display: flex; gap: 6px; align-items: center; flex-wrap: wrap; }
.btn-icon { display: flex; justify-content: center; align-items: center; width: 32px; height: 32px; border-radius: 8px; border: none; font-size: 1.1rem; background-color: transparent; text-decoration: none; cursor: default;}
.btn-icon.youtube { color: #ef4444; background-color: #fef2f2; border: 1px solid #fee2e2; cursor: pointer; transition: transform 0.2s; }
.btn-icon.youtube:hover { transform: translateY(-2px); }
.btn-icon.sequence { color: #8b5cf6; background-color: #f5f3ff; border: 1px solid #ddd6fe; cursor: pointer;}
.btn-icon.score { color: #1f2937; background-color: #f9fafb; border: 1px solid #e5e7eb; cursor: pointer;}
.btn-icon.chart { color: #3b82f6; background-color: #eff6ff; border: 1px solid #dbeafe; cursor: pointer;}
.btn-icon.voice { color: #db2777; background-color: #fdf2f8; border: 1px solid #fbcfe8; cursor: pointer; }
.btn-icon.guitar { color: #ea580c; background-color: #fff7ed; border: 1px solid #ffedd5; cursor: pointer; }
.btn-icon.piano { color: #0d9488; background-color: #f0fdfa; border: 1px solid #ccfbf1; cursor: pointer; }
.btn-icon.drums { color: #ca8a04; background-color: #fefce8; border: 1px solid #fef08a; cursor: pointer; }
.btn-icon.bass { color: #4f46e5; background-color: #eef2ff; border: 1px solid #e0e7ff; cursor: pointer; }
.btn-icon:hover { transform: translateY(-2px); }
.actions-row { display: flex; justify-content: center; gap: 8px; }
.btn-action { border: none; width: 32px; height: 32px; border-radius: 6px; cursor: pointer; font-size: 1rem; transition: transform 0.1s; display: flex; justify-content: center; align-items: center;}
.btn-action:hover { transform: scale(1.1); }
.btn-action.edit { background-color: #eff6ff; color: #3b82f6; border: 1px solid #dbeafe; }
.btn-action.delete { background-color: #fef2f2; color: #ef4444; border: 1px solid #fee2e2; }
.btn-structure { background-color: var(--color-secundary, #2c3e50); color: white; border: none; padding: 5px 12px; border-radius: 15px; font-size: 0.75rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-structure:hover { background-color: #1a252f; }
.text-muted { color: #ccc; font-size: 1.2rem; }
.empty-state { text-align: center; padding: 40px; color: #999; font-style: italic; }
.error-banner { background-color: #fee2e2; color: #b91c1c; padding: 10px; border-radius: 8px; margin-bottom: 15px; text-align: center; }
.loading-text { color: #888; font-size: 0.9rem; font-style: italic; }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.modal-content.form-modal { max-width: 650px; width: 95%; max-height: 90vh; overflow-y: auto; }
.modal-content { background: white; padding: 30px; border-radius: 20px; width: 90%; max-width: 400px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); position: relative; animation: slideUp 0.3s ease-out; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
.close-btn { position: absolute; top: 15px; right: 20px; background: none; border: none; font-size: 2rem; color: #999; cursor: pointer; }
.close-btn:hover { color: #333; }
.modal-header { text-align: center; margin-bottom: 20px; }
.modal-header h3 { margin: 0; color: var(--color-secundary, #2c3e50); }
.modal-header p { margin: 5px 0 0; color: #666; font-size: 0.9rem; }
.structure-timeline { display: flex; flex-direction: column; gap: 0; border-left: 2px solid #e5e7eb; margin-left: 20px; padding: 10px 0; }
.structure-item { position: relative; padding-left: 20px; padding-bottom: 20px; }
.structure-item:last-child { padding-bottom: 0; }
.dot { width: 12px; height: 12px; background-color: var(--color-tertiary, #e67e22); border-radius: 50%; position: absolute; left: -7px; top: 5px; border: 2px solid white; box-shadow: 0 0 0 1px var(--color-tertiary, #e67e22); }
.part-name { font-weight: 600; color: #333; font-size: 1rem; }
.structure-builder-section { border: 1px solid #e2e8f0; padding: 15px; border-radius: 10px; margin-bottom: 15px; background: #fcfcfc; }
.form-label { font-size: 0.85rem; font-weight: 600; color: #4b5563; margin-bottom: 8px; display: block; }
.structure-buttons { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 15px; }
.badge-btn { background: white; border: 1px solid #cbd5e1; padding: 5px 10px; border-radius: 20px; cursor: pointer; font-size: 0.8rem; color: #475569; transition: all 0.2s; }
.badge-btn:hover { background: #e2e8f0; border-color: #94a3b8; }
.structure-preview { min-height: 40px; border: 2px dashed #e2e8f0; border-radius: 8px; padding: 10px; background: white; display: flex; align-items: center; }
.preview-placeholder { color: #cbd5e1; font-style: italic; font-size: 0.85rem; width: 100%; text-align: center; }
.preview-chips { display: flex; flex-wrap: wrap; gap: 8px; }
.struct-chip { background-color: var(--color-secundary, #2c3e50); color: white; padding: 4px 10px; border-radius: 6px; font-size: 0.85rem; display: flex; align-items: center; gap: 6px; }
.remove-x { cursor: pointer; font-weight: bold; opacity: 0.7; }
.remove-x:hover { opacity: 1; }
.form-section-title { font-weight: 700; color: var(--color-secundary, #2c3e50); border-bottom: 1px solid #eee; padding-bottom: 5px; margin-top: 15px; font-size: 0.9rem; }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 15px; border-top: 1px solid #eee; padding-top: 20px; }
.resources-checks { display: flex; gap: 20px; align-items: center; flex-wrap: wrap; }
.checkbox-item { display: flex; align-items: center; gap: 8px; font-size: 0.9rem; color: #4b5563; cursor: pointer; }

@media (max-width: 600px) {
  .form-row { grid-template-columns: 1fr; gap: 10px; }
  .full-width { grid-column: span 1; }
}
</style>