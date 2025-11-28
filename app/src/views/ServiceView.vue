<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router' 
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

// --- 1. CONSTANTES Y CONFIGURACIÓN ---
// Lista maestra para saber qué archivo es cada ID
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

// --- 2. INTERFACES ---
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
  voiceUrl: string;
  guitarUrl: string;
  pianoUrl: string;
  drumsUrl: string;
  bassUrl: string;
  sequenceUrl: string;
  chartUrl: string;
  scoreUrl: string;
}

interface Musician {
  id: number;
  name: string;
  role: string;
  secondaryRole: string;
  photo: string;
  status: string;
}

interface EventInfo {
  date: string;
  title: string;
}

// Nueva interfaz para los outfits vinculados
interface LinkedOutfit {
  id: number;
  name: string;
  filename: string;
}

// --- 3. ESTADO ---
const currentEventInfo = ref<EventInfo>({ date: 'Cargando...', title: '' }) 
const songs = ref<Song[]>([]) 
const musicians = ref<Musician[]>([]) 
const linkedOutfits = ref<LinkedOutfit[]>([]) // <--- AQUÍ GUARDAMOS LOS OUTFITS DEL SERVICIO
const isLoading = ref(true)
const errorMsg = ref('')

// Estados para Modales
const isModalOpen = ref(false) // Modal estructura
const isOutfitModalOpen = ref(false) // Modal outfits

// Helper para convertir links de archivos
function getFileLink(fileName: string, folder: string) {
  if (!fileName || fileName === '#') return '#'
  if (fileName.startsWith('http')) return fileName
  return `http://localhost:8080/files/${folder}/${fileName}`
}

// --- 4. CARGAR SERVICIO ---
async function getServiceDetail() {
  const id = route.params.id 
  isLoading.value = true
  errorMsg.value = ''
  
  try {
    const response = await fetch(`http://localhost:8080/services/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      }
    })

    if (!response.ok) throw new Error('No se pudo cargar el servicio')

    const data = await response.json()

    // 1. Info Evento
    if (data.service) {
        const dateObj = new Date(data.service.start_time)
        const formattedDate = dateObj.toLocaleDateString('es-ES', { day: 'numeric', month: 'long' })
        const finalDate = formattedDate.charAt(0).toUpperCase() + formattedDate.slice(1)

        currentEventInfo.value = {
            date: finalDate, 
            title: data.service.name || 'Servicio'
        }
    }

    // 2. Canciones
    if (data.songs) {
        songs.value = data.songs.map((s: any) => ({
            id: s.id,
            name: s.name,
            author: s.author || 'Autor desconocido',
            key: s.song_key || '-',
            bpm: s.bpm || 0,
            timeSignature: s.time_signature || '4/4',
            duration: s.duration || '0:00',
            structure: s.structure ? s.structure.split(',') : [],
            hasSequence: s.has_sequence,
            hasChart: s.has_chart,
            hasScore: s.has_score,
            youtubeUrl: s.youtube_url,
            sequenceUrl: s.sequence_url || s.has_sequence || '',
            chartUrl: s.chart_url || s.has_chart || '',
            scoreUrl: s.score_url || s.has_score || '',
            voiceUrl: s.voice_url || '',
            guitarUrl: s.guitar_url || '',
            pianoUrl: s.piano_url || '',
            drumsUrl: s.drums_url || '',
            bassUrl: s.bass_url || ''
        }))
    }

    // 3. Músicos
    if (data.users) {
        musicians.value = data.users.map((u: any) => ({
            id: u.id,
            name: u.username,
            role: u.role,
            secondaryRole: u.secondary_role || '', 
            photo: u.profile_picture_url || '/logo.png',
            status: u.status || 'pending'
        }))
    }

    // 4. OUTFITS (NUEVO)
    if (data.outfits) {
        // Mapeamos los IDs que vienen del backend con nuestra lista maestra
        linkedOutfits.value = data.outfits.map((o: any) => {
            const found = AVAILABLE_OUTFITS.find(avail => avail.id === o.outfit_id);
            return found ? found : null;
        }).filter((item: any) => item !== null);
    }

  } catch (e: any) {
    console.error(e)
    errorMsg.value = 'Error al cargar la información.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => { getServiceDetail() })

// --- UI HELPERS ---
const selectedSongTitle = ref('')
const selectedSongStructure = ref<string[]>([])

const totalDuration = computed(() => {
  let totalSeconds = 0
  songs.value.forEach(song => {
    if (song.duration) {
      const parts = song.duration.split(':')
      if(parts.length === 2) {
         totalSeconds += (parseInt(parts[0]) * 60) + parseInt(parts[1])
      }
    }
  })
  const finalMinutes = Math.floor(totalSeconds / 60)
  const finalSeconds = totalSeconds % 60
  return `${finalMinutes}:${finalSeconds.toString().padStart(2, '0')}`
})

// Modal Estructura
function openStructureModal(song: Song) {
  if (!song.structure || song.structure.length === 0) return
  selectedSongTitle.value = song.name
  selectedSongStructure.value = song.structure
  isModalOpen.value = true
}
function closeModal() { isModalOpen.value = false }

// Modal Outfits
function openOutfitModal() {
    isOutfitModalOpen.value = true
}
function closeOutfitModal() {
    isOutfitModalOpen.value = false
}

function getRoleEmoji(role: string) { 
  if (!role) return '👤'
  const r = role.toLowerCase(); 
  if (r.includes('cantante') || r.includes('voz')) return '🎤'; 
  if (r.includes('guitarra')) return '🎸'; 
  if (r.includes('batería') || r.includes('bateria')) return '🥁'; 
  if (r.includes('piano') || r.includes('teclado')) return '🎹'; 
  if (r.includes('bajo')) return '🎸'; 
  if (r.includes('sonido')) return '🎚';
  return '🎵'; 
}

function getStatusIcon(status: string) { 
  if (status === 'accepted') return '✅'; 
  if (status === 'rejected') return '❌'; 
  return '❔'; 
}
</script>

<template>
  <main>
    <div class="programation">
      
      <div v-if="isLoading" class="loading-msg">Cargando detalles...</div>
      <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>

      <div v-if="!isLoading && !errorMsg">
          
          <div class="header-row">
            <div class="badge date">{{ currentEventInfo.date }}</div>
            <div class="badge event">{{ currentEventInfo.title }}</div>
          </div>
          
          <div class="team-songs">
            
            <div class="card songs">
                <div class="card-header-flex">
                    <h3>Repertorio</h3>
                    <div class="total-badge">{{ totalDuration }}</div>
                </div>
                <table class="custom-table songs-table">
                    <thead>
                        <tr>
                            <th class="center-text">Tono</th>
                            <th>Canción</th>
                            <th>Recursos</th>
                            <th>Tracks</th>
                            <th class="center-text">Estruc.</th>
                            <th class="center-text">Bpm</th>
                            <th class="center-text">Cifra</th>
                            <th class="center-text">Dur.</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="songs.length === 0">
                            <td colspan="8" class="empty-row">No hay canciones en este servicio.</td>
                        </tr>
                        <tr v-for="song in songs" :key="song.id">
                            <td class="center-text"><div class="key-badge">{{ song.key }}</div></td>
                            <td>
                                <div class="info-stack">
                                    <span class="primary-text">{{ song.name }}</span>
                                    <span class="secondary-text">{{ song.author }}</span>
                                </div>
                            </td>
                            <td>
                                <div class="resource-actions">
                                    <a v-if="song.youtubeUrl" :href="song.youtubeUrl" target="_blank" class="btn-icon youtube">▶</a>
                                    <a v-if="song.sequenceUrl" :href="getFileLink(song.sequenceUrl, 'sequences')" target="_blank" class="btn-icon sequence">🎧</a>
                                    <a v-if="song.chartUrl" :href="getFileLink(song.chartUrl, 'charts')" target="_blank" class="btn-icon chart">📄</a>
                                    <a v-if="song.scoreUrl" :href="getFileLink(song.scoreUrl, 'scores')" target="_blank" class="btn-icon score">🎼</a>
                                </div>
                            </td>
                            <td>
                                <div class="resource-actions">
                                    <a v-if="song.voiceUrl" :href="getFileLink(song.voiceUrl, 'tracks')" target="_blank" class="btn-icon voice">🎤</a>
                                    <a v-if="song.pianoUrl" :href="getFileLink(song.pianoUrl, 'tracks')" target="_blank" class="btn-icon piano">🎹</a>
                                    <a v-if="song.guitarUrl" :href="getFileLink(song.guitarUrl, 'tracks')" target="_blank" class="btn-icon guitar">🎸</a>
                                    <a v-if="song.bassUrl" :href="getFileLink(song.bassUrl, 'tracks')" target="_blank" class="btn-icon bass">🎸</a>
                                    <a v-if="song.drumsUrl" :href="getFileLink(song.drumsUrl, 'tracks')" target="_blank" class="btn-icon drums">🥁</a>
                                </div>
                            </td>
                            <td class="center-text">
                                <button v-if="song.structure.length > 0" @click="openStructureModal(song)" class="btn-structure">Ver</button>
                                <span v-else class="text-muted">-</span>
                            </td>
                            <td class="center-text bpm-cell">{{ song.bpm }}</td>
                            <td class="center-text time-sig-cell">{{ song.timeSignature }}</td>
                            <td class="center-text duration-cell">{{ song.duration }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div class="card team">
                <div class="card-header-flex">
                    <h3>Equipo</h3>
                    <button class="btn-outfit" @click="openOutfitModal" :disabled="linkedOutfits.length === 0">
                        👕 Paleta ({{ linkedOutfits.length }})
                    </button>
                </div>
                <table class="custom-table">
                    <tbody>
                        <tr v-if="musicians.length === 0">
                            <td colspan="2" class="empty-row">No hay músicos asignados.</td>
                        </tr>
                        <tr v-for="musician in musicians" :key="musician.id">
                            <td style="width: 60px;">
                                <div class="avatar-wrapper">
                                    <img :src="musician.photo" alt="foto" class="avatar" :class="`status-${musician.status}`" />
                                    <div class="status-icon">{{ getStatusIcon(musician.status) }}</div>
                                </div>
                            </td>
                            <td>
                                <div class="info-stack">
                                  <div class="roles-tags-container">
                                    <span class="role-tag main">
                                        {{ musician.role }}
                                    </span>
                                  </div>

                                    <span class="primary-text">{{ musician.name }}</span>
                                    
                                    <div class="roles-tags-container">
                                        <template v-if="musician.secondaryRole">
                                            <span 
                                                v-for="subRole in musician.secondaryRole.split(',')" 
                                                :key="subRole" 
                                                class="role-tag sub"
                                            >
                                                {{ getRoleEmoji(subRole) }} {{ subRole.trim() }}
                                            </span>
                                        </template>
                                    </div>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

          </div>
      </div>

      <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
            <div class="modal-content">
                <button class="close-btn" @click="closeModal">×</button>
                <div class="modal-header">
                    <h3>Estructura</h3>
                    <p>{{ selectedSongTitle }}</p>
                </div>
                <div class="structure-timeline">
                    <div v-for="(part, index) in selectedSongStructure" :key="index" class="structure-item">
                        <div class="dot"></div>
                        <div class="part-name">{{ part }}</div>
                    </div>
                </div>
            </div>
      </div>

      <div v-if="isOutfitModalOpen" class="modal-overlay" @click.self="closeOutfitModal">
            <div class="modal-content modal-lg">
                <button class="close-btn" @click="closeOutfitModal">×</button>
                <div class="modal-header">
                    <h3>Paleta de Colores</h3>
                    <p>{{ currentEventInfo.title }}</p>
                </div>
                
                <div class="outfits-container">
                    <div v-if="linkedOutfits.length === 0" class="empty-msg">
                        No hay estilos asignados a este servicio.
                    </div>
                    
                    <div v-for="outfit in linkedOutfits" :key="outfit.id" class="outfit-preview-card">
                        <h4>{{ outfit.name }}</h4>
                        <div class="pdf-wrapper">
                             <iframe 
                                :src="`/files/outfits/${outfit.filename}#view=FitH&toolbar=0`"
                                title="PDF Preview"
                                width="100%"
                                height="300px"
                             ></iframe>
                        </div>
                        <div class="actions">
                            <a :href="`/files/outfits/${outfit.filename}`" target="_blank" class="btn-download">
                                Abrir en pestaña nueva ↗
                            </a>
                        </div>
                    </div>
                </div>
            </div>
      </div>

    </div>
  </main>
</template>

<style scoped>
/* ESTILOS BASE (Mismos que tenías) */
.programation { display: flex; flex-direction: column; padding: 40px; gap: 20px; }
.header-row { display: flex; justify-content: flex-start; gap: 0; margin-bottom: 20px; }
.badge { padding: 10px 20px; font-weight: 700; color: white; text-transform: uppercase; }
.badge.date { background-color: var(--color-secundary, #2c3e50); border-radius: 20px 0 0 20px; }
.badge.event { background-color: var(--color-tertiary, #e67e22); border-radius: 0 20px 20px 0; }
.team-songs { display: flex; align-items: flex-start; gap: 20px; flex-wrap: wrap; }
.card { background: white; border: 1px solid #eee; padding: 20px; border-radius: 16px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
.team { flex: 1; width: -webkit-fill-available; }
.songs { flex: 3; min-width: 800px; }
.total-badge { background-color: #f3f4f6; color: #4b5563; font-weight: 700; padding: 6px 12px; border-radius: 8px; font-family: monospace; font-size: 0.9rem; border: 1px solid #e5e7eb; display: flex; align-items: center; gap: 5px; }
.card-header-flex { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px; margin-bottom: 15px; width: 100%; }
.card h3 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.2rem; }

/* Botón modificado */
.btn-outfit { background-color: #f3f4f6; border: 1px solid #e5e7eb; color: #4b5563; padding: 6px 12px; border-radius: 8px; font-size: 0.8rem; font-weight: 600; cursor: pointer; transition: all 0.2s; display: flex; gap: 5px; align-items: center; }
.btn-outfit:hover:not(:disabled) { background-color: #e5e7eb; color: #1f2937; }
.btn-outfit:disabled { opacity: 0.5; cursor: not-allowed; }

.custom-table { width: 100%; border-collapse: collapse; }
.custom-table th { text-align: left; font-size: 0.75rem; text-transform: uppercase; color: #999; padding: 10px; font-weight: 600; }
.custom-table td { padding: 12px 8px; border-bottom: 1px dashed #eee; vertical-align: middle; }
.center-text { text-align: center; }
.empty-row { text-align: center; color: #999; font-style: italic; padding: 20px; }
.info-stack { display: flex; flex-direction: column; }
.primary-text { font-weight: 700; color: #333; font-size: 0.95rem; margin-top: 6px; margin-bottom: 6px;}
.roles-tags-container { display: flex; flex-wrap: wrap; gap: 6px; }
.role-tag { padding: 3px 8px; border-radius: 12px; font-size: 0.75rem; font-weight: 600; display: inline-flex; align-items: center; gap: 3px; white-space: nowrap; }
.role-tag.main { background-color: #e0f2fe; color: #0369a1; border: 1px solid #bae6fd; text-transform: uppercase; font-size: 0.65rem; letter-spacing: 0.5px; }
.role-tag.sub { background-color: #f3f4f6; color: #4b5563; border: 1px solid #e5e7eb; font-style: italic; }
.bpm-cell { font-family: monospace; font-weight: 600; color: #555; }
.time-sig-cell { font-family: monospace; color: #666; font-weight: 600; background: #f9fafb; padding: 4px; border-radius: 4px; font-size: 0.85rem;}
.duration-cell { font-family: monospace; color: #333; }
.avatar-wrapper { position: relative; width: 45px; height: 45px; }
.avatar { width: 100%; height: 100%; border-radius: 50%; object-fit: cover; border: 3px solid transparent; }
.avatar.status-confirmed { border-color: #22c55e; }
.avatar.status-declined { border-color: #ef4444; }
.avatar.status-pending { border-color: #d1d5db; }
.status-icon { position: absolute; bottom: -5px; right: -5px; font-size: 0.8rem; background: white; width: 18px; height: 18px; display: flex; justify-content: center; align-items: center; box-shadow: 0 1px 3px rgba(0,0,0,0.2); }
.key-badge { background-color: #f3f4f6; color: #333; font-weight: 800; width: 35px; height: 35px; display: flex; align-items: center; justify-content: center; border-radius: 8px; font-size: 0.9rem; border: 1px solid #e5e7eb; margin: 0 auto; }
.resource-actions { display: flex; gap: 6px; align-items: center; flex-wrap: wrap; }
.btn-icon { display: flex; justify-content: center; align-items: center; width: 32px; height: 32px; border-radius: 8px; border: none; cursor: pointer; font-size: 1.1rem; transition: transform 0.2s; background-color: transparent; text-decoration: none;}
.btn-icon:hover { transform: translateY(-2px); }
.btn-icon.youtube { color: #ef4444; background-color: #fef2f2; border: 1px solid #fee2e2; }
.btn-icon.sequence { color: #8b5cf6; background-color: #f5f3ff; border: 1px solid #ddd6fe; }
.btn-icon.score { color: #1f2937; background-color: #f9fafb; border: 1px solid #e5e7eb; }
.btn-icon.chart { color: #3b82f6; background-color: #eff6ff; border: 1px solid #dbeafe; }
.btn-icon.voice { color: #db2777; background-color: #fdf2f8; border: 1px solid #fbcfe8; }
.btn-icon.guitar { color: #ea580c; background-color: #fff7ed; border: 1px solid #ffedd5; }
.btn-icon.piano { color: #0d9488; background-color: #f0fdfa; border: 1px solid #ccfbf1; }
.btn-icon.drums { color: #ca8a04; background-color: #fefce8; border: 1px solid #fef08a; }
.btn-icon.bass { color: #4f46e5; background-color: #eef2ff; border: 1px solid #e0e7ff; }
.btn-structure { background-color: var(--color-secundary, #2c3e50); color: white; border: none; padding: 5px 12px; border-radius: 15px; font-size: 0.75rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-structure:hover { background-color: #1a252f; }
.text-muted { color: #ccc; font-size: 1.2rem; }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.modal-content { background: white; padding: 30px; border-radius: 20px; width: 90%; max-width: 400px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); position: relative; animation: slideUp 0.3s ease-out; max-height: 90vh; overflow-y: auto;}
/* Nuevo modificador para el modal grande */
.modal-content.modal-lg { max-width: 700px; }

@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
.close-btn { position: absolute; top: 15px; right: 20px; background: none; border: none; font-size: 2rem; color: #999; cursor: pointer; z-index: 10; }
.close-btn:hover { color: #333; }
.modal-header { text-align: center; margin-bottom: 20px; }
.modal-header h3 { margin: 0; color: var(--color-secundary, #2c3e50); }
.modal-header p { margin: 5px 0 0; color: #666; font-size: 0.9rem; }
.structure-timeline { display: flex; flex-direction: column; gap: 0; border-left: 2px solid #e5e7eb; margin-left: 20px; padding: 10px 0; }
.structure-item { position: relative; padding-left: 20px; padding-bottom: 20px; }
.structure-item:last-child { padding-bottom: 0; }
.dot { width: 12px; height: 12px; background-color: var(--color-tertiary, #e67e22); border-radius: 50%; position: absolute; left: -7px; top: 5px; border: 2px solid white; box-shadow: 0 0 0 1px var(--color-tertiary, #e67e22); }
.part-name { font-weight: 600; color: #333; font-size: 1rem; }
.loading-msg { text-align: center; color: #888; padding: 20px; font-style: italic; }
.error-msg { background-color: #fee2e2; color: #b91c1c; padding: 10px; border-radius: 8px; text-align: center; margin-bottom: 15px; }

/* ESTILOS PARA LOS OUTFITS */
.outfits-container { display: flex; flex-direction: column; gap: 20px; }
.outfit-preview-card { border: 1px solid #eee; padding: 15px; border-radius: 12px; background-color: #f9fafb; }
.outfit-preview-card h4 { margin: 0 0 10px 0; color: #333; }
.pdf-wrapper { width: 100%; border: 1px solid #ddd; background: #525659; /* Color de fondo típico de visor PDF */ }
.actions { margin-top: 10px; text-align: right; }
.btn-download { font-size: 0.85rem; color: #2563eb; text-decoration: none; font-weight: 600; }
.btn-download:hover { text-decoration: underline; }
.empty-msg { text-align: center; color: #999; font-style: italic; }

@media (max-width: 900px) { .team-songs { flex-direction: column; } .team, .songs { width: 100%; } }
</style>