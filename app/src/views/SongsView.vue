<script setup lang="ts">
import { ref, computed } from 'vue'

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
  hasSequence: boolean;
  hasChart: boolean;
  hasScore: boolean;
  youtubeUrl: string;
  // En la vista admin no mostramos los links de tracks individuales, 
  // pero los mantenemos en el objeto por si al editar los necesitas.
  voiceUrl: string;
  guitarUrl: string;
  pianoUrl: string;
  drumsUrl: string;
  bassUrl: string;
}

// --- 2. DATOS (Simulación de BD) ---
const allSongs = ref<Song[]>([
  { 
    id: 1, name: 'La Gloria de Dios', author: 'Evaluna', key: 'C', bpm: 72, timeSignature: '4/4', duration: '4:25',
    structure: ['Intro', 'Verso 1', 'Coro', 'Verso 2', 'Coro', 'Puente', 'Coro', 'Final'],
    hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: 'https://youtube.com', voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 2, name: 'Jehova', author: 'Barcelona Worship', key: 'Em', bpm: 130, timeSignature: '6/8', duration: '5:10',
    structure: ['Intro', 'Verso', 'Coro', 'Puente', 'Outro'],
    hasSequence: true, hasChart: true, hasScore: true, youtubeUrl: 'https://youtube.com', voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 3, name: 'Rompimiento', author: 'New Wine', key: 'Bb', bpm: 145, timeSignature: '4/4', duration: '3:50',
    structure: [],
    hasSequence: false, hasChart: true, hasScore: false, youtubeUrl: '', voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 4, name: 'Danzando', author: 'Gateway', key: 'G', bpm: 110, timeSignature: '4/4', duration: '4:00',
    structure: ['Intro', 'Coro', 'Puente'],
    hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: '', voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 5, name: 'Vida Tu Me Das', author: 'Hillsong', key: 'D', bpm: 76, timeSignature: '4/4', duration: '5:30',
    structure: [],
    hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: '', voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 6, name: 'Yaweh', author: 'Christine DClario', key: 'A', bpm: 68, timeSignature: '4/4', duration: '6:00',
    structure: [],
    hasSequence: false, hasChart: true, hasScore: false, youtubeUrl: '', voiceUrl: '', guitarUrl: '', pianoUrl: '', drumsUrl: '#', bassUrl: ''
  }
])

// --- 3. LÓGICA DE BÚSQUEDA ---
const searchQuery = ref('')

const filteredSongs = computed(() => {
  if (!searchQuery.value) return allSongs.value
  const term = searchQuery.value.toLowerCase()
  return allSongs.value.filter(song => 
    song.name.toLowerCase().includes(term) || 
    song.author.toLowerCase().includes(term)
  )
})

// --- 4. MODAL Y ACCIONES ---
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

// Acciones de Admin
function editSong(id: number) {
  console.log('Ir a editar canción:', id)
  // router.push(`/songs/edit/${id}`)
}
function deleteSong(id: number) {
  if(confirm('¿Eliminar esta canción?')) {
    allSongs.value = allSongs.value.filter(s => s.id !== id)
  }
}
</script>

<template>
  <main>
    <div class="programation"> <div class="view-header-row">
        <div>
           <h1>Biblioteca</h1>
           <p class="subtitle">Repertorio general</p>
        </div>
        <div class="header-actions">
           <div class="search-box">
              <span class="search-icon">🔍</span>
              <input v-model="searchQuery" type="text" placeholder="Buscar canción..." />
           </div>
           <button class="btn-create">
             + Nueva Canción
           </button>
        </div>
      </div>

      <div class="card songs">
         <div class="card-header-flex">
             <h3>Todas las Canciones ({{ filteredSongs.length }})</h3>
         </div>

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
                  <th class="center-text">Acciones</th> </tr>
            </thead>
            <tbody>
               <tr v-for="song in filteredSongs" :key="song.id">
                  
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
                        <span v-if="song.hasSequence" class="btn-icon sequence" title="Tiene Secuencia">🎧</span>
                        <span v-if="song.hasScore" class="btn-icon score" title="Tiene Cifrado">📄</span>
                        <span v-if="song.hasChart" class="btn-icon chart" title="Tiene Partitura">🎼</span>
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

                  <td class="center-text">
                     <div class="actions-row">
                        <button @click="editSong(song.id)" class="btn-action edit" title="Editar">✎</button>
                        <button @click="deleteSong(song.id)" class="btn-action delete" title="Eliminar">🗑</button>
                     </div>
                  </td>

               </tr>
            </tbody>
         </table>
         
         <div v-if="filteredSongs.length === 0" class="empty-state">
            No se encontraron canciones.
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
/* REUTILIZANDO TUS ESTILOS EXACTOS */
.programation { display: flex; flex-direction: column; padding: 40px; gap: 20px; }

/* Estilos Header Propio de esta vista */
.view-header-row { 
  display: flex; justify-content: space-between; align-items: center; 
  margin-bottom: 10px; flex-wrap: wrap; gap: 20px;
}
h1 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.5rem; }
.subtitle { margin: 0; color: #888; font-size: 0.9rem; }

.header-actions { display: flex; gap: 15px; align-items: center; }

/* Search Box Estilizado */
.search-box {
  position: relative; display: flex; align-items: center;
}
.search-box input {
  padding: 8px 10px 8px 35px;
  border-radius: 8px; border: 1px solid #e5e7eb;
  font-size: 0.9rem; outline: none; width: 250px;
  transition: border 0.2s;
}
.search-box input:focus { border-color: var(--color-secundary, #2c3e50); }
.search-icon { position: absolute; left: 10px; font-size: 0.8rem; opacity: 0.6; }

/* Botón Crear */
.btn-create {
  background-color: #10b981; color: white; border: none;
  padding: 8px 16px; border-radius: 8px; font-weight: 600;
  cursor: pointer; transition: background 0.2s; display: flex; align-items: center; gap: 5px;
}
.btn-create:hover { background-color: #059669; }



/* TARJETA Y TABLA (TUS ESTILOS) */
.card { background: white; border: 1px solid #eee; padding: 20px; border-radius: 16px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); width: 100%; }
.card-header-flex { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px; margin-bottom: 15px; width: 100%; }
.card h3 { margin: 0; color: var(--color-secundary, #2c3e50); font-size: 1.2rem; }

.custom-table { width: 100%; border-collapse: collapse; }
.custom-table th { text-align: left; font-size: 0.75rem; text-transform: uppercase; color: #999; padding: 10px; font-weight: 600; }
.custom-table td { padding: 12px 8px; border-bottom: 1px dashed #eee; vertical-align: middle; }
.center-text { text-align: center; }

/* CELDAS */
.info-stack { display: flex; flex-direction: column; }
.primary-text { font-weight: 700; color: #333; font-size: 0.95rem; }
.secondary-text { font-size: 0.8rem; color: #888; margin-top: 2px; }
.key-badge { background-color: #f3f4f6; color: #333; font-weight: 800; width: 35px; height: 35px; display: flex; align-items: center; justify-content: center; border-radius: 8px; font-size: 0.9rem; border: 1px solid #e5e7eb; margin: 0 auto; }
.bpm-cell { font-family: monospace; font-weight: 600; color: #555; }
.time-sig-cell { font-family: monospace; color: #666; font-weight: 600; background: #f9fafb; padding: 4px; border-radius: 4px; font-size: 0.85rem;}
.duration-cell { font-family: monospace; color: #333; }

/* RECURSOS (Iconos) */
.resource-actions { display: flex; gap: 6px; align-items: center; justify-content: flex-start; flex-wrap: wrap; }
.btn-icon { display: flex; justify-content: center; align-items: center; width: 32px; height: 32px; border-radius: 8px; border: none; cursor: default; font-size: 1.1rem; background-color: transparent; text-decoration: none;}
/* Youtube si es link */
.btn-icon.youtube { color: #ef4444; background-color: #fef2f2; border: 1px solid #fee2e2; cursor: pointer; transition: transform 0.2s; }
.btn-icon.youtube:hover { transform: translateY(-2px); }
.btn-icon.sequence { color: #8b5cf6; background-color: #f5f3ff; border: 1px solid #ddd6fe; }
.btn-icon.score { color: #1f2937; background-color: #f9fafb; border: 1px solid #e5e7eb; }
.btn-icon.chart { color: #3b82f6; background-color: #eff6ff; border: 1px solid #dbeafe; }

.btn-icon.voice { color: #db2777; background-color: #fdf2f8; border: 1px solid #fbcfe8; }
.btn-icon.guitar { color: #ea580c; background-color: #fff7ed; border: 1px solid #ffedd5; }
.btn-icon.piano { color: #0d9488; background-color: #f0fdfa; border: 1px solid #ccfbf1; }
.btn-icon.drums { color: #ca8a04; background-color: #fefce8; border: 1px solid #fef08a; }
.btn-icon.bass { color: #4f46e5; background-color: #eef2ff; border: 1px solid #e0e7ff; }

/* ESTRUCTURA */
.btn-structure { background-color: var(--color-secundary, #2c3e50); color: white; border: none; padding: 5px 12px; border-radius: 15px; font-size: 0.75rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-structure:hover { background-color: #1a252f; }
.text-muted { color: #ccc; font-size: 1.2rem; }

/* --- NUEVOS ESTILOS: ACCIONES --- */
.actions-row { display: flex; justify-content: center; gap: 8px; }
.btn-action { border: none; width: 32px; height: 32px; border-radius: 6px; cursor: pointer; font-size: 1rem; transition: transform 0.1s; display: flex; justify-content: center; align-items: center;}
.btn-action:hover { transform: scale(1.1); }
.btn-action.edit { background-color: #eff6ff; color: #3b82f6; border: 1px solid #dbeafe; } /* Azul */
.btn-action.delete { background-color: #fef2f2; color: #ef4444; border: 1px solid #fee2e2; } /* Rojo */

/* Empty State */
.empty-state { text-align: center; padding: 40px; color: #999; font-style: italic; }

/* MODAL */
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

@media (max-width: 768px) {
  .view-header-row { flex-direction: column; align-items: flex-start; }
  .search-box input { width: 100%; }
  .header-actions { width: 100%; flex-direction: column; align-items: stretch; }
}
</style>