<template>
  <nav>
  <router-link to="/">Dashboard</router-link> |
  <router-link to="/login">Login</router-link> |
  <router-link to="/register">Register</router-link> |
  <router-link to="/todos">Todos</router-link> |
  <span v-if="user">
          {{ user.email }}
        </span>
        <span v-else>
          Not logged in
  </span>
  </nav>

    
  <router-view />
</template>
<script setup>
import { computed } from "vue";
import { currentUser } from "../stores/user";

const user = computed(() => currentUser.value);

</script>
