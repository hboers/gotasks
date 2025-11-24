import { createRouter, createWebHistory } from "vue-router";

import LoginView from "./views/LoginView.vue";
import LogoutView from "./views/LogoutView.vue";
import RegisterView from "./views/RegisterView.vue";
import DashboardView from "./views/DashboardView.vue";
import TodoView from "./views/TodoView.vue";

const routes = [
  { path: "/login",  component: LoginView },
  { path: "/logout",  component: LogoutView },
  { path: "/todos",  component: TodoView },
  { path: "/register",  component: RegisterView },
  { path: "/", component: DashboardView },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
