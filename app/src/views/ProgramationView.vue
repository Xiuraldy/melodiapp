<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// --- 1. DEFINIR INTERFACES (El "molde" de tus datos) ---
interface Song {
  id: number;
  name: string;
  author: string;
  key: string;
  bpm: number;
  timeSignature: string;
  duration: string;
  structure: string[];
  hasSequence: boolean;
  hasChart: boolean;
  hasScore: boolean;
  youtubeUrl: string;
  voiceUrl: string;
  guitarUrl: string;
  pianoUrl: string;
  drumsUrl: string;
  bassUrl: string;
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

// Estructura completa de un evento en la "base de datos"
interface EventData {
  info: EventInfo;
  musicians: Musician[];
  songs: Song[];
}

// --- 2. VARIABLES REACTIVAS TIPADAS ---
// Aquí le decimos a Vue: "Este array empieza vacío, pero contendrá objetos tipo Song"
const currentEventInfo = ref<EventInfo>({ date: '', title: '' }) 
const songs = ref<Song[]>([]) 
const musicians = ref<Musician[]>([]) 

// --- 3. BASE DE DATOS TIPADA ---
// Record<number, EventData> significa: Las llaves son números, los valores son EventData
const eventsDatabase: Record<number, EventData> = {
  1: {
    info: { date: 'Diciembre 24', title: 'Evento: Navidad' },
    musicians: [
      { id: 1, name: 'Valentina Carrillo', role: 'Cantante', secondaryRole: '', photo: '/logo.png', status: 'confirmed' },
      { id: 2, name: 'Juan Perez', role: 'Guitarrista', secondaryRole: 'Director Musical', photo: '/logo.png', status: 'pending' },
      { id: 3, name: 'Ana Gomez', role: 'Baterista', secondaryRole: '', photo: '/logo.png', status: 'declined' }
    ],
    songs: [
        { 
            id: 1, name: 'La Gloria de Dios', author: 'Evaluna', key: 'C', bpm: 72, timeSignature: '4/4', duration: '4:25',
            structure: ['Intro', 'Verso 1', 'Coro', 'Verso 2', 'Coro', 'Puente', 'Coro', 'Final'],
            hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: 'https://youtube.com',
            voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
        },
        { 
            id: 2, name: 'Jehova', author: 'Barcelona Worship', key: 'Em', bpm: 130, timeSignature: '6/8', duration: '5:10',
            structure: ['Intro', 'Verso', 'Coro', 'Puente', 'Outro'],
            hasSequence: true, hasChart: true, hasScore: true, youtubeUrl: 'https://youtube.com',
            voiceUrl: '#', guitarUrl: '#', pianoUrl: '', drumsUrl: '#', bassUrl: '#'
        },
        { 
            id: 3, name: 'Rompimiento', author: 'New Wine', key: 'Bb', bpm: 145, timeSignature: '4/4', duration: '3:50',
            structure: [],
            hasSequence: false, hasChart: true, hasScore: false, youtubeUrl: '', 
            voiceUrl: '#', guitarUrl: '', pianoUrl: '', drumsUrl: '#', bassUrl: ''
        }
    ]
  },
  2: {
    info: { date: 'Enero 30', title: 'Culto General' },
    musicians: [
      { id: 1, name: 'Pedro Bajista', role: 'Bajista', secondaryRole: '', photo: '/logo.png', status: 'confirmed' },
      { id: 2, name: 'Juan Perez', role: 'Guitarrista', secondaryRole: '', photo: '/logo.png', status: 'confirmed' }
    ],
    songs: [
        { 
            id: 4, name: 'Danzando', author: 'Gateway', key: 'G', bpm: 110, timeSignature: '4/4', duration: '4:00',
            structure: ['Intro', 'Coro', 'Puente'],
            hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: '',
            voiceUrl: '#', guitarUrl: '#', pianoUrl: '', drumsUrl: '', bassUrl: ''
        },
        { 
            id: 5, name: 'Vida Tu Me Das', author: 'Hillsong', key: 'D', bpm: 76, timeSignature: '4/4', duration: '5:30',
            structure: [],
            hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: '',
            voiceUrl: '#', guitarUrl: '#', pianoUrl: '', drumsUrl: '', bassUrl: ''
        },
        { 
            id: 6, name: 'Yaweh', author: 'Christine DClario', key: 'A', bpm: 68, timeSignature: '4/4', duration: '6:00',
            structure: [],
            hasSequence: false, hasChart: true, hasScore: false, youtubeUrl: '',
            voiceUrl: '', guitarUrl: '', pianoUrl: '', drumsUrl: '', bassUrl: ''
        }
    ]
  }
}

// --- FUNCIÓN PARA CARGAR DATOS ---
onMounted(() => {
    // CORRECCIÓN CLAVE: Convertimos el ID de la ruta a Número explícitamente
    const id = Number(route.params.id)
    
    // Buscamos en nuestra Base de datos
    const eventData = eventsDatabase[id]

    if (eventData) {
        currentEventInfo.value = eventData.info
        musicians.value = eventData.musicians
        songs.value = eventData.songs
    } else {
        console.error("Evento no encontrado para el ID:", id)
    }
})

// --- LÓGICA DE UI (Sin cambios, solo tipado inferido) ---
const isModalOpen = ref(false)
const selectedSongTitle = ref('')
const selectedSongStructure = ref<string[]>([])

const totalDuration = computed(() => {
  let totalSeconds = 0
  songs.value.forEach(song => {
    if (song.duration) {
      const parts = song.duration.split(':')
      const minutes = parseInt(parts[0])
      const seconds = parseInt(parts[1])
      if (!isNaN(minutes) && !isNaN(seconds)) {
         totalSeconds += (minutes * 60) + seconds
      }
    }
  })
  const finalMinutes = Math.floor(totalSeconds / 60)
  const finalSeconds = totalSeconds % 60
  return `${finalMinutes}:${finalSeconds.toString().padStart(2, '0')}`
})

// Helpers
function openStructureModal(song: Song) { // Tipamos el argumento como Song
  if (!song.structure || song.structure.length === 0) return
  selectedSongTitle.value = song.name
  selectedSongStructure.value = song.structure
  isModalOpen.value = true
}

function closeModal() { isModalOpen.value = false }

function getRoleEmoji(role: string) { 
  const r = role.toLowerCase(); 
  if (r.includes('cantante') || r.includes('voz')) return '🎤'; 
  if (r.includes('guitarra')) return '🎸'; 
  if (r.includes('batería') || r.includes('bateria')) return '🥁'; 
  if (r.includes('piano') || r.includes('teclado')) return '🎹'; 
  if (r.includes('bajo')) return '🎸'; 
  return '🎵'; 
}

function getStatusIcon(status: string) { 
  if (status === 'confirmed') return '✅'; 
  if (status === 'declined') return '❌'; 
  return '❔'; 
}
</script>

<template>
  <main>
    <div class="programation">
      
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
                         <th>Estruc.</th>
                         <th>Bpm</th>
                         <th>Cifra</th>
                         <th>Dur.</th>
                         </tr>
                 </thead>
                 <tbody>
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
                                <a v-if="song.youtubeUrl" :href="song.youtubeUrl" target="_blank" class="btn-icon youtube" title="Ver Youtube">▶</a>
                                <button v-if="song.hasSequence" class="btn-icon sequence" title="Secuencia">🎧</button>
                                <button v-if="song.hasScore" class="btn-icon score" title="Cifrado">📄</button>
                                <button v-if="song.hasChart" class="btn-icon chart" title="Partitura">🎼</button>
                            </div>
                        </td>
                        <td>
                          <div class="resource-actions">
                            <a v-if="song.voiceUrl" :href="song.voiceUrl" class="btn-icon voice" title="Voces">🎤</a>
                            <a v-if="song.pianoUrl" :href="song.pianoUrl" class="btn-icon piano" title="Track Piano">🎹</a>
                            <a v-if="song.guitarUrl" :href="song.guitarUrl" class="btn-icon guitar" title="Track Guitarra">🎸</a>
                            <a v-if="song.bassUrl" :href="song.bassUrl" class="btn-icon bass" title="Track Bajo">🎸</a>
                            <a v-if="song.drumsUrl" :href="song.drumsUrl" class="btn-icon drums" title="Track Bateria">🥁</a>
                          </div>
                        </td>
                         <td class="center-text">
                            <button v-if="song.structure && song.structure.length > 0" @click="openStructureModal(song)" class="btn-structure">Ver</button>
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
                <button class="btn-outfit">👕 Paleta de Colores</button>
             </div>
             <table class="custom-table">
                 <tbody>
                     <tr v-for="musician in musicians" :key="musician.id">
                         <td style="width: 60px;">
                             <div class="avatar-wrapper">
                                 <img :src="musician.photo" alt="foto" class="avatar" :class="`status-${musician.status}`" />
                                 <div class="status-icon">{{ getStatusIcon(musician.status) }}</div>
                             </div>
                         </td>
                         <td>
                             <div class="info-stack">
                                 <span class="primary-text">{{ musician.name }}</span>
                                 <span class="secondary-text">
                                     {{ getRoleEmoji(musician.role) }} {{ musician.role }}
                                     <span v-if="musician.secondaryRole"> • {{ musician.secondaryRole }}</span>
                                 </span>
                             </div>
                         </td>
                     </tr>
                 </tbody>
             </table>
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
    </div>
  </main>
</template>

<style scoped>
/* COPIA EXACTAMENTE TUS ESTILOS DE LA RESPUESTA ANTERIOR AQUÍ */
/* No hace falta cambiarlos, funcionarán igual */
.programation { display: flex; flex-direction: column; padding: 40px; gap: 20px; }
.header-row { display: flex; justify-content: flex-start; gap: 0; }
.badge { padding: 10px 20px; font-weight: 700; color: white; text-transform: uppercase; }
.badge.date { background-color: var(--color-secundary, #2c3e50); border-radius: 20px 0 0 20px; }
.badge.event { background-color: var(--color-tertiary, #e67e22); border-radius: 0 20px 20px 0; }
.team-songs { display: flex; align-items: flex-start; gap: 20px; flex-wrap: wrap; }
/* ... resto de estilos ... */
.card { background: white; border: 1px solid #eee; padding: 20px; border-radius: 16px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
.team { flex: 1; width: -webkit-fill-available; }
.songs { flex: 3; min-width: 800px; }
.total-badge { background-color: #f3f4f6; color: #4b5563; font-weight: 700; padding: 6px 12px; border-radius: 8px; font-family: monospace; font-size: 0.9rem; border: 1px solid #e5e7eb; display: flex; align-items: center; gap: 5px; }
.card-header-flex { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px; margin-bottom: 15px; width: 100%; }
.card h3 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.2rem; }
.btn-outfit { background-color: #f3f4f6; border: 1px solid #e5e7eb; color: #4b5563; padding: 6px 12px; border-radius: 8px; font-size: 0.8rem; font-weight: 600; cursor: pointer; transition: all 0.2s; display: flex; gap: 5px; align-items: center; }
.btn-outfit:hover { background-color: #e5e7eb; color: #1f2937; }
.custom-table { width: 100%; border-collapse: collapse; }
.custom-table th { text-align: left; font-size: 0.75rem; text-transform: uppercase; color: #999; padding: 10px; font-weight: 600; }
.custom-table td { padding: 12px 8px; border-bottom: 1px dashed #eee; vertical-align: middle; }
.center-text { text-align: center; }
.info-stack { display: flex; flex-direction: column; }
.primary-text { font-weight: 700; color: #333; font-size: 0.95rem; }
.secondary-text { font-size: 0.8rem; color: #888; margin-top: 2px; }
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
@media (max-width: 900px) { .team-songs { flex-direction: column; } .team, .songs { width: 100%; } }
</style>