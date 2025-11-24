<template>
  <div class="container mt-3">
    <h1>Todos</h1>

    <ul>
      <li v-for="t in todos" :key="t.id">
        {{ t.title }}
      </li>
    </ul>

    <form @submit.prevent="addTodo">
      <input v-model="newTodo" placeholder="New todo" />
      <button>Add</button>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { fetchTodos } from "../api";

const router = useRouter();
const todos = ref([]);
const error = ref("");
const newTodo = ref("");   

onMounted(async () => {
  try {
    todos.value = await fetchTodos();
  } catch (e) {
    console.log("fetchTodos error:", e, "status:", e.status); // debug
    if (e.status === 401) {
      await router.push("/login");
      return;
    }
    error.value = e.message || "Failed to load todos";
  }
});
</script>