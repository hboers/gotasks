import { createApp } from 'vue'
import App from './App.vue'
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min.js";
import "bootstrap-icons/font/bootstrap-icons.css";

import { router } from "./router"
import { loadUser } from "../stores/user"
loadUser()
createApp(App)
    .use(router)
    .mount('#app')
