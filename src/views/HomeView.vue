<script setup lang="ts">
import { ref } from "vue";
import { db, getAllDecks, createDeck as createDeckDb } from "../db";
import { useLiveQuery } from "../use-live-query";

const decks = useLiveQuery(() => getAllDecks(db), []);
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
        class="flex-1 border rounded px-3 py-2"
        data-testid="deck-name-input"
      />
      <button
        type="submit"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
        data-testid="create-deck-btn"
      >
        Create
      </button>
    </form>

    <p v-if="decks.length === 0" class="text-gray-500" data-testid="empty-state">
      No decks yet. Create one above.
    </p>

    <ul class="space-y-2">
      <li v-for="deck in decks" :key="deck.id">
        <router-link
          :to="`/deck/${deck.id}`"
          class="block p-3 bg-white rounded shadow hover:shadow-md"
          :data-testid="`deck-${deck.id}`"
        >
          {{ deck.name }}
        </router-link>
      </li>
    </ul>
  </div>
</template>
