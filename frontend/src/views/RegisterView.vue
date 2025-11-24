<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-4">

        <div class="card shadow-sm">
          <div class="card-body">
            <h1 class="h4 mb-3 text-center">Registrieren</h1>

            <!-- Fehlermeldung -->
            <div
              v-if="error"
              class="alert alert-danger py-2"
              role="alert"
            >
              {{ error }}
            </div>

            <!-- Erfolgsmeldung (optional) -->
            <div
              v-if="success"
              class="alert alert-success py-2"
              role="alert"
            >
              {{ success }}
            </div>

            <form @submit.prevent="doRegister">
              <div class="mb-3">
                <label for="name" class="form-label">Name</label>
                <input
                  id="name"
                  v-model="name"
                  type="text"
                  class="form-control"
                  placeholder="Ihr Name"
                  required
                  autocomplete="name"
                />
              </div>

              <div class="mb-3">
                <label for="email" class="form-label">E-Mail</label>
                <input
                  id="email"
                  v-model="email"
                  type="email"
                  class="form-control"
                  placeholder="name@example.com"
                  required
                  autocomplete="email"
                />
              </div>

              <div class="mb-3">
                <label for="password" class="form-label">Passwort</label>
                <input
                  id="password"
                  v-model="password"
                  type="password"
                  class="form-control"
                  placeholder="Passwort"
                  required
                  autocomplete="new-password"
                />
              </div>

              <button
                type="submit"
                class="btn btn-primary w-100"
                :disabled="loading"
              >
                <span
                  v-if="loading"
                  class="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                Konto anlegen
              </button>
            </form>

            <div class="mt-3 text-center">
              <small>
                Schon ein Konto?
                <router-link to="/login">Zum Login</router-link>
              </small>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
// Passe den Import ggf. an deinen api.js Namen an:
import { register } from "../api";

const router = useRouter();

const name = ref("");
const email = ref("");
const password = ref("");
const error = ref("");
const success = ref("");
const loading = ref(false);

async function doRegister() {
  error.value = "";
  success.value = "";
  loading.value = true;

  try {
    await register(email.value, name.value, password.value);

    success.value = "Registrierung erfolgreich. Du kannst dich jetzt einloggen.";
    // Kurze Pause optional, direktes Redirect ist meist praktischer:
    await router.push("/login");
  } catch (e) {
    error.value = e?.message || "Registrierung fehlgeschlagen";
  } finally {
    loading.value = false;
  }
}
</script>