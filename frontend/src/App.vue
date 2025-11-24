<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    
    <router-link class="navbar-brand" to="/">TodoApp</router-link>
    
    <button
      class="navbar-toggler"
      type="button"
      data-bs-toggle="collapse"
      data-bs-target="#navbarSupportedContent"
    >
      <span class="navbar-toggler-icon"></span>
    </button>
      
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <template v-if="user">
        <li class="nav-item">
          <router-link class="nav-link" to="/todos">Todos</router-link> 
        </li>
        </template>
        <template v-else>
        <li class="nav-item" v-if="user">
          <router-link class="nav-link" to="/login">Login</router-link> 
        </li>
        </template>
      </ul>
      <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
        <template v-if="user">
          <li class="nav-item">
            <span class="nav-link disabled">
              <i class="bi bi-person-circle"></i>
              {{ user.email }}
            </span>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/logout">Logout</router-link>
          </li>
        </template>
      </ul>
  </nav>
  <router-view />
</template>
<script setup>
  import { computed } from "vue";
  import { currentUser } from "../stores/user";
  const user = computed(() => currentUser.value);
</script>
