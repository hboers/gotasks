<template>
  <div class="container mt-4">
    <h1 class="h3 mb-3">Todos</h1>

    <!-- Neues Todo -->
    <form class="mb-3" @submit.prevent="addTodo">
      <div class="input-group">
        <input
          v-model="newTodo"
          type="text"
          class="form-control"
          placeholder="Neues Todo…"
          required
        />
        <button class="btn btn-primary" type="submit">
          Hinzufügen
        </button>
      </div>
    </form>

    <!-- Fehler -->
    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <!-- Liste -->
    <ul class="list-group">
      <li
        v-for="t in todos"
        :key="t.id"
        class="list-group-item d-flex justify-content-between align-items-center"
      >
        <span>
          <input
            class="form-check-input me-2"
            type="checkbox"
            :checked="t.done"
            @change="toggleTodo(t)"
          />
          {{ t.title }}
        </span>
      <button
      class="btn btn-sm btn-outline-danger"
      @click="onDelete(t)"
    >
      <i class="bi bi-trash"></i>
    </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { fetchTodos, createTodo, updateTodo, deleteTodo } from "../api"

const todos = ref([])
const newTodo = ref("")
const error = ref("")

onMounted(async () => {
  try {
    todos.value = await fetchTodos()
  } catch (e) {
    error.value = e.message || "Fehler beim Laden der Todos"
  }
})

// 🔥 HIER: fehlendes addTodo
async function addTodo() {
  if (!newTodo.value.trim()) {
    return
  }

  try {
    // neues Todo im Backend anlegen
    const created = await createTodo(newTodo.value.trim())

    // lokal zur Liste hinzufügen
    todos.value.push(created)

    // Eingabefeld leeren
    newTodo.value = ""
  } catch (e) {
    error.value = e.message || "Fehler beim Anlegen des Todos"
  }
}

async function toggleTodo(todo) {
  updateTodo(todo.id, !todo.done) 
  todo.done = !todo.done
}
async function onDelete(todo) {
  if (!confirm(`Todo "${todo.title}" wirklich löschen?`)) {
    return;
  }

  try {
    await deleteTodo(todo.id);

    // Sofort aus der Liste entfernen
    todos.value = todos.value.filter(t => t.id !== todo.id);
  } catch (e) {
    error.value = e.message || "Fehler beim Löschen";
  }
}
</script>