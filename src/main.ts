import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./main.css";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: () => import("./views/HomeView.vue") },
    { path: "/deck/:id", component: () => import("./views/DeckView.vue") },
    { path: "/review/:deckId", component: () => import("./views/ReviewView.vue") },
  ],
});

createApp(App).use(router).mount("#app");
