<template>
  <div>
    <h1>Todos</h1>

    <router-link to="/">Logout</router-link>

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
import { fetchTodos, createTodo } from "../api";

const todos = ref([]);
const newTodo = ref("");

onMounted(async () => {
  todos.value = await fetchTodos();
});

async function addTodo() {
  const todo = await createTodo(newTodo.value);
  todos.value.push(todo);
  newTodo.value = "";
}
</script>
