<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

interface User {
  ID: number; 
  username: string;
  lastname: string;
  email: string;
  celphone: string;
  role: string;
}

const authStore = useAuthStore()
const getUsersError = ref('') 

const users = ref<User[]>([]) 

async function getUsers() {
  getUsersError.value = ''

  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}` 
      },
    })

    const data = await response.json()

    if (data.error) {
      getUsersError.value = data.error 
      return
    }
    
    users.value = data

  } catch (error) {
    console.error(error)
    getUsersError.value = 'Error de conexión'
  }
}

onMounted(() => {
  getUsers()
})
</script>

<template>
  <main>
    <div class="error" v-if="getUsersError">⚠︎ Error: {{ getUsersError }}</div>
    <button>Crear Usuario</button>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Nombre</th>
            <th>Rol</th> <th>Correo Electrónico</th>
            <th>Celular</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.ID"> 
            <td>{{ user.username }}</td>
            <td>{{ user.role }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.celphone }}</td>
            <td v-if="authStore.isAdmin" class="actions-cell">
                <button class="btn-icon edit" title="Editar">✎</button>
                <button class="btn-icon delete" title="Eliminar">🗑</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<style scoped>

.table-container {
  margin: 2rem auto;
  width: 90%;
  border-radius: 12px;
  overflow: hidden; 
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.05);
  font-family: 'Poppins', sans-serif;
}

table {
  width: 100%;
  border-collapse: collapse;
  background-color: #ffffff;
  text-align: left;
}

thead {
  background-color: var(--color-tertiary, #4f46e5); /* Usa tu color o un índigo por defecto */
  color: #ffffff;
}

th {
  padding: 16px 20px;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.85rem;
  letter-spacing: 0.5px;
}

td {
  padding: 14px 20px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
  font-size: 0.95rem;
}

tbody tr:nth-child(even) {
  background-color: #f9fafb;
}

tbody tr:hover {
  background-color: #eef2ff; /* Un color suave relacionado con tu terciario */
  transform: scale(1.005); /* Un micro efecto de zoom */
  transition: all 0.2s ease-in-out;
  cursor: default;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
}

td:nth-child(2) {
  font-weight: 600;
  color: var(--color-tertiary, #4f46e5);
}

.error {
  background-color: #fee2e2;
  color: #ef4444;
  padding: 1rem;
  border-radius: 8px;
  margin: 1rem auto;
  width: 90%;
  text-align: center;
  border: 1px solid #fca5a5;
}

@media (max-width: 768px) {
  .table-container {
    width: 100%;
    border-radius: 0;
    box-shadow: none;
    overflow-x: auto; 
  }
  
  table {
    min-width: 600px; 
  }
}
</style>