import { createApp } from 'vue'
import App from './App.vue'

import { router } from "./router"
import { loadUser } from "../stores/user"
loadUser()
createApp(App)
    .use(router)
    .mount('#app')
