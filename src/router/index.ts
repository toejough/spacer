import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    {
      path: '/deck/:id',
      name: 'deck',
      component: () => import('../views/DeckView.vue'),
    },
    {
      path: '/review/:deckId?',
      name: 'review',
      component: () => import('../views/ReviewView.vue'),
    },
  ],
})

export default router
