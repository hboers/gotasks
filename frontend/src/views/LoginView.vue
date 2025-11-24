<template>
  <div>
    <h1>Login</h1>

    <form @submit.prevent="doLogin">
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Login</button>
    </form>

    <router-link to="/register">Create an account</router-link>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { login } from "../api";

const email = ref("");
const password = ref("");

async function doLogin() {
  await login(email.value, password.value);
  location.href = "/"; // navigation
}
</script>
