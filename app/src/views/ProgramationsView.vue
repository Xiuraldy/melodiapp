<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ref } from 'vue';

// --- 1. TUS DATOS QUEMADOS EXACTOS (Canciones del 24 Dic) ---
const songs = ref([
  { 
    id: 1, 
    name: 'La Gloria de Dios', 
    author: 'Evaluna', 
    key: 'C', 
    bpm: 72,
    timeSignature: '4/4', 
    duration: '4:25',
    structure: ['Intro', 'Verso 1', 'Coro', 'Verso 2', 'Coro', 'Puente', 'Coro', 'Final'],
    hasSequence: true, hasChart: true, hasScore: false, youtubeUrl: 'https://youtube.com',
    voiceUrl: '#', guitarUrl: '#', pianoUrl: '#', drumsUrl: '', bassUrl: ''
  },
  { 
    id: 2, 
    name: 'Jehova', 
    author: 'Barcelona Worship', 
    key: 'Em', 
    bpm: 130, 
    timeSignature: '6/8', 
    duration: '5:10',
    structure: ['Intro', 'Verso', 'Coro', 'Puente', 'Outro'],
    hasSequence: true, hasChart: true, hasScore: true, youtubeUrl: 'https://youtube.com',
    voiceUrl: '#', guitarUrl: '#', pianoUrl: '', drumsUrl: '#', bassUrl: '#'
  },
  { 
    id: 3, 
    name: 'Rompimiento', 
    author: 'New Wine', 
    key: 'Bb', 
    bpm: 145, 
    timeSignature: '4/4', 
    duration: '3:50',
    structure: [],
    hasSequence: false, hasChart: true, hasScore: false, youtubeUrl: '', 
    voiceUrl: '#', guitarUrl: '', pianoUrl: '', drumsUrl: '#', bassUrl: ''
  }
])

// --- 2. ARRAY DE TARJETAS (PROGRAMACIONES) ---
const programations = ref([
  {
    id: 1,
    day: '24',
    month: 'DIC',
    event: 'Culto General',
    role: 'Cantante',
    // ASIGNAMOS TUS DATOS AQUÍ 👇
    songsList: songs.value 
  },
  {
    id: 2,
    day: '30',
    month: 'ENE',
    event: 'Culto General',
    role: 'Pro. en Vivo',
    // Datos dummy simples para la segunda tarjeta (para que no se vea vacía)
    songsList: [
      { id: 4, name: 'Danzando' }, 
      { id: 5, name: 'Vida Tu Me Das' }, 
      { id: 6, name: 'Yaweh' }
    ]
  }
])

const router = useRouter();

function redirectProgramation(id: number) {
  // AHORA: Usamos el nombre de la ruta y le pasamos el parámetro
  router.push({ name: 'programation', params: { id: id } })
}
</script>

<template>
  <main>
    <div class="programations-container">
      
      <div 
        v-for="prog in programations" 
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
            <div class="role-badge">{{ prog.role }}</div>
            <h3 class="event-title">{{ prog.event }}</h3>
          </div>
        </div>

        <div class="divider"></div>

        <div class="card-body">
          <div class="songs-label">Repertorio:</div>
          <ul class="song-list">
            <li v-for="song in prog.songsList" :key="song.id">
              <span class="music-icon">🎵</span> {{ song.name }}
            </li>
          </ul>
        </div>

        <div class="card-footer">
          <span>Ver detalle</span>
          <span class="arrow">→</span>
        </div>

      </div>

    </div>
  </main>
</template>

<style scoped>
/* Contenedor Principal */
.programations-container {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  padding: 50px;
  gap: 30px;
}

/* TARJETA BASE */
.card-programation {
  background-color: white;
  width: 320px;
  border-radius: 16px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.card-programation:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
  border-color: var(--color-secundary, #2c3e50);
}

/* --- CABECERA --- */
.card-header {
  display: flex;
  padding: 20px;
  gap: 15px;
  align-items: center;
  background: linear-gradient(to bottom, #ffffff, #fafafa);
}

/* Calendario */
.calendar-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: white;
  border: 2px solid var(--color-secundary, #2c3e50);
  border-radius: 12px;
  min-width: 60px;
  height: 65px;
}

.calendar-box .month {
  font-size: 0.7rem;
  background-color: var(--color-secundary, #2c3e50);
  color: white;
  width: 100%;
  text-align: center;
  padding: 2px 0;
  font-weight: 700;
}

.calendar-box .day {
  font-size: 1.8rem;
  font-weight: 800;
  color: #333;
  line-height: 1;
  margin-top: 2px;
}

/* Info */
.header-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.role-badge {
  font-size: 0.75rem;
  background-color: var(--color-tertiary, #e67e22);
  color: white;
  padding: 4px 10px;
  border-radius: 20px;
  font-weight: 600;
  margin-bottom: 5px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.event-title {
  margin: 0;
  font-size: 1.1rem;
  color: #333;
  font-weight: 700;
  line-height: 1.2;
}

/* Divider */
.divider {
  height: 1px;
  background: #eee;
  margin: 0 20px;
}

/* --- CUERPO --- */
.card-body {
  padding: 20px;
  flex-grow: 1;
}

.songs-label {
  font-size: 0.75rem;
  color: #999;
  font-weight: 700;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.song-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.song-list li {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: #555;
  margin-bottom: 6px;
  padding: 8px 10px;
  background-color: #f8fafc;
  border-radius: 8px;
  transition: background 0.2s;
  border: 1px solid transparent;
}

.song-list li:hover {
  background-color: #fff;
  border-color: #eee;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
}

.music-icon {
  margin-right: 10px;
  font-size: 0.8rem;
  filter: grayscale(100%);
  opacity: 0.6;
}

/* --- PIE --- */
.card-footer {
  padding: 15px 20px;
  background-color: #fcfcfc;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-secundary, #2c3e50);
}

.arrow {
  transition: transform 0.2s;
}

.card-programation:hover .arrow {
  transform: translateX(5px);
}

/* Responsive */
@media (max-width: 600px) {
  .programations-container {
    justify-content: center;
  }
}
</style>