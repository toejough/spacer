<script setup lang="ts">
import { ref } from "vue";
import { db, getAllDecks, createDeck as createDeckDb } from "../db";
import { useLiveViewState } from "../view-state";
import { Library } from "lucide-vue-next";

const decksState = useLiveViewState(() => getAllDecks(db), (d) => d.length === 0);
const newDeckName = ref("");

async function createDeck() {
  const name = newDeckName.value.trim();
  if (!name) return;
  await createDeckDb(db, name);
  newDeckName.value = "";
}
</script>

<template>
  <div>
    <h2 class="text-lg font-semibold mb-4">Your Decks</h2>

    <form @submit.prevent="createDeck" class="flex gap-2 mb-4">
      <input
        v-model="newDeckName"
        placeholder="New deck name"
        class="flex-1 border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
        data-testid="deck-name-input"
      />
      <button
        type="submit"
        class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
        data-testid="create-deck-btn"
      >
        Create
      </button>
    </form>

    <!-- Loading skeleton -->
    <div v-if="decksState.status === 'loading'" data-testid="loading-skeleton" class="space-y-2">
      <div v-for="i in 3" :key="i" class="h-12 bg-surface rounded-md animate-skeleton" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="decksState.status === 'empty'"
      data-testid="empty-state"
      class="text-center py-12"
    >
      <Library class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">No decks yet</p>
      <p class="text-text-faint text-sm">Create your first deck to start studying</p>
    </div>

    <!-- Error state -->
    <div
      v-else-if="decksState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm">{{ decksState.message }}</p>
    </div>

    <!-- Loaded state -->
    <ul v-else-if="decksState.status === 'loaded'" class="space-y-2">
      <li v-for="deck in decksState.data" :key="deck.id">
        <router-link
          :to="`/deck/${deck.id}`"
          class="block p-3 bg-surface rounded-md shadow-card hover:shadow-md border-l-[3px] border-secondary"
          :data-testid="`deck-${deck.id}`"
        >
          {{ deck.name }}
        </router-link>
      </li>
    </ul>
  </div>
</template>
