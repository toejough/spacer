<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import { db, getDeck, getDeckCards, getDueCards, createCard as createCardDb } from "../db";
import { useLiveViewState } from "../view-state";
import { useLiveQuery } from "../use-live-query";
import { Layers, Search, AlertTriangle } from "lucide-vue-next";

const route = useRoute();
const deckId = Number(route.params.id);

// Primary state: deck existence (not-found vs loaded)
const deckState = useLiveViewState(
  () => getDeck(db, deckId),
  () => false,
  { notFoundWhen: (d) => d === undefined }
);

// Secondary state: card list (empty vs loaded)
const cardsState = useLiveViewState(
  () => getDeckCards(db, deckId),
  (c) => c.length === 0
);

// Due cards — just a count for the review button, no view state needed
const dueCards = useLiveQuery(() => getDueCards(db, deckId), []);
const dueCount = computed(() => dueCards.value.length);

const front = ref("");
const back = ref("");

async function addCard() {
  const f = front.value.trim();
  const b = back.value.trim();
  if (!f || !b) return;
  await createCardDb(db, deckId, f, b);
  front.value = "";
  back.value = "";
}
</script>

<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="deckState.status === 'loading'" data-testid="loading-skeleton">
      <div class="flex items-center justify-between mb-4">
        <div class="h-6 w-32 bg-surface rounded-sm animate-skeleton" />
        <div class="h-9 w-24 bg-surface rounded-sm animate-skeleton" />
      </div>
      <div class="space-y-2 mb-4">
        <div class="h-10 bg-surface rounded-sm animate-skeleton" />
        <div class="h-10 bg-surface rounded-sm animate-skeleton" />
        <div class="h-9 w-24 bg-surface rounded-sm animate-skeleton" />
      </div>
      <div class="space-y-2">
        <div v-for="i in 2" :key="i" class="h-16 bg-surface rounded-md animate-skeleton" />
      </div>
    </div>

    <!-- Not found -->
    <div
      v-else-if="deckState.status === 'not-found'"
      data-testid="not-found"
      class="text-center py-12"
    >
      <Search class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Deck not found</p>
      <p class="text-text-faint text-sm mb-4">This deck may have been deleted</p>
      <router-link to="/" class="text-primary text-sm">&larr; Back to decks</router-link>
    </div>

    <!-- Error -->
    <div
      v-else-if="deckState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <AlertTriangle class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm">{{ deckState.message }}</p>
    </div>

    <!-- Deck loaded -->
    <div v-else-if="deckState.status === 'loaded'">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold">{{ deckState.data.name }}</h2>
        <router-link
          :to="`/review/${deckId}`"
          v-if="dueCount > 0"
          class="bg-accent hover:bg-accent-hover text-white px-4 py-2 rounded-sm text-sm"
          data-testid="start-review-btn"
        >
          Review ({{ dueCount }} due)
        </router-link>
      </div>

      <form @submit.prevent="addCard" class="space-y-2 mb-4">
        <input
          v-model="front"
          placeholder="Front (question)"
          class="w-full border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
          data-testid="card-front-input"
        />
        <input
          v-model="back"
          placeholder="Back (answer)"
          class="w-full border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
          data-testid="card-back-input"
        />
        <button
          type="submit"
          class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
          data-testid="add-card-btn"
        >
          Add Card
        </button>
      </form>

      <!-- Empty cards -->
      <div
        v-if="cardsState.status === 'empty'"
        data-testid="no-cards"
        class="text-center py-8"
      >
        <Layers class="mx-auto mb-2 text-text-faint" :size="36" />
        <p class="text-text-muted font-medium">No cards yet</p>
        <p class="text-text-faint text-sm">Add cards using the form above</p>
      </div>

      <!-- Card list -->
      <ul v-else-if="cardsState.status === 'loaded'" class="space-y-2">
        <li
          v-for="card in cardsState.data"
          :key="card.id"
          class="p-3 bg-surface rounded-md shadow-card"
          :data-testid="`card-${card.id}`"
        >
          <div class="font-medium">{{ card.front }}</div>
          <div class="text-text-muted text-sm">{{ card.back }}</div>
        </li>
      </ul>

      <router-link to="/" class="inline-block mt-4 text-primary text-sm">
        &larr; Back to decks
      </router-link>
    </div>
  </div>
</template>
