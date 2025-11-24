<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="btn-group">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <router-link class="nav-link" to="/">Dashboard</router-link> 
        </li>
        <li class="nav-item" v-if="user">
          <router-link class="nav-link" to="/login">Login</router-link> 
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/register">Register</router-link> 
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/todos">Todos</router-link> 
        </li>
      </ul>
        <ul class="navbar-nav ms-auto mb-2 mb-lg-0">

          <!-- Not logged in -->
          <template v-if="!user">
            <li class="nav-item">
              <router-link class="nav-link" to="/login">Login</router-link>
            </li>

            <li class="nav-item">
              <router-link class="nav-link" to="/register">Register</router-link>
            </li>
          </template>

          <!-- Logged in -->
          <template v-else>
            <li class="nav-item">
              <span class="navbar-text me-3">
                <i class="bi bi-person-circle"></i>
                {{ user.email }}
              </span>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/logout">Logout</router-link>
            </li>
          </template>
        </ul>
    </div>
  <span class="float-end" v-if="user">
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
