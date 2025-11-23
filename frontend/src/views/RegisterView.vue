<template>
  <div>
    <h1>Register</h1>

    <form @submit.prevent="doRegister">
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="name" type="text" placeholder="Name" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Register</button>
    </form>

    <router-link to="/">Back to login</router-link>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { register } from "../api";

const email = ref("");
const name = ref("");
const password = ref("");

async function doRegister() {
  await register(email.value, name.value, password.value);
  alert("Registered!");
  location.href = "/";
}
</script>
