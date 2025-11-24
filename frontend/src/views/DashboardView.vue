<template>
  <div>
    <h1>DashBoard</h1>

    <router-link to="/logout">Logout</router-link>

  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { currentUser } from "../../stores/user";

console.log(currentUser.value?.email);
</script>