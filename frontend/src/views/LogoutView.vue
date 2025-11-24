<template>
  <div class="logout">
    <p>Logging out...</p>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { currentUser } from "../../stores/user";

const router = useRouter();

onMounted(async () => {
  console.log("[LogoutView] mounted, starting logout");

  try {
    await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });
    console.log("[LogoutView] logout request finished");
  } catch (e) {
    console.error("[LogoutView] logout failed", e);
  }

  // clear frontend state
  currentUser.value = null;
  console.log("[LogoutView] currentUser cleared, redirecting to /login");

  // use named route if possible
  await router.push("/login");
});
</script>