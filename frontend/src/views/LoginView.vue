<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-4">

        <div class="card shadow-sm">
          <div class="card-body">
            <h1 class="h4 mb-3 text-center">Login</h1>

            <!-- Fehlermeldung -->
            <div
              v-if="error"
              class="alert alert-danger py-2"
              role="alert"
            >
              {{ error }}
            </div>

            <form @submit.prevent="doLogin">
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
                  autocomplete="current-password"
                />
              </div>

              <button
                type="submit"
                class="btn btn-primary w-100"
                :disabled="loading"
              >
                <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                Anmelden
              </button>
            </form>

            <div class="mt-3 text-center">
              <small>
                Noch kein Konto?
                <router-link to="/register">Registrieren</router-link>
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
import { login } from "../api";

const email = ref("");
const password = ref("");

async function doLogin() {
  await login(email.value, password.value);
  location.href = "/"; // navigation
}
</script>
