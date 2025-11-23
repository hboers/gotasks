import { createRouter, createWebHistory } from "vue-router";

import LoginView from "./views/LoginView.vue";
import RegisterView from "./views/RegisterView.vue";
import TodoView from "./views/TodoView.vue";

const routes = [
  { path: "/", component: LoginView },
  { path: "/register", component: RegisterView },
  { path: "/todos", component: TodoView },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
