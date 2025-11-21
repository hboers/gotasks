<template>
  <main class="app">
    <h1>Go + Vue Todo App</h1>

    <form @submit.prevent="addTodo" class="todo-form">
      <input
        v-model="newTitle"
        type="text"
        placeholder="New todo..."
        required
      />
      <button type="submit">Add</button>
    </form>

    <p v-if="loading">Loading…</p>
    <p v-if="error" class="error">{{ error }}</p>

    <ul class="todo-list">
      <li
        v-for="todo in todos"
        :key="todo.id"
        :class="{ done: todo.done }"
      >
        <label>
          <input
            type="checkbox"
            :checked="todo.done"
            @change="toggleDone(todo)"
          />
          {{ todo.title }}
        </label>
        <button @click="removeTodo(todo)">✕</button>
      </li>
    </ul>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchTodos, createTodo, updateTodo, deleteTodo } from './api'

const todos = ref([])
const newTitle = ref('')
const loading = ref(false)
const error = ref('')

async function loadTodos() {
  loading.value = true
  error.value = ''
  try {
    todos.value = await fetchTodos()
  } catch (e) {
    error.value = e.message || 'Failed to load todos'
  } finally {
    loading.value = false
  }
}

async function addTodo() {
  if (!newTitle.value.trim()) return
  try {
    const todo = await createTodo(newTitle.value.trim())
    todos.value.push(todo)
    newTitle.value = ''
  } catch (e) {
    error.value = e.message || 'Failed to create todo'
  }
}

async function toggleDone(todo) {
  try {
    const updated = await updateTodo(todo.id, { done: !todo.done })
    const idx = todos.value.findIndex(t => t.id === todo.id)
    if (idx !== -1) todos.value[idx] = updated
  } catch (e) {
    error.value = e.message || 'Failed to update todo'
  }
}

async function removeTodo(todo) {
  try {
    await deleteTodo(todo.id)
    todos.value = todos.value.filter(t => t.id !== todo.id)
  } catch (e) {
    error.value = e.message || 'Failed to delete todo'
  }
}

onMounted(loadTodos)
</script>

<style scoped>
.app {
  max-width: 480px;
  margin: 2rem auto;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
}
.todo-form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}
.todo-form input {
  flex: 1;
}
.todo-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.todo-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: .4rem .2rem;
  border-bottom: 1px solid #ddd;
}
.todo-list li.done label {
  text-decoration: line-through;
  opacity: 0.6;
}
button {
  cursor: pointer;
}
.error {
  color: red;
}
</style>
