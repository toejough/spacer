<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import { db, getDeck, getDeckCards, getDueCards, createCard as createCardDb } from "../db";
import { useLiveQuery } from "../use-live-query";

const route = useRoute();
const deckId = Number(route.params.id);

const deck = useLiveQuery(() => getDeck(db, deckId), undefined);
const cards = useLiveQuery(() => getDeckCards(db, deckId), []);
const dueCards = useLiveQuery(() => getDueCards(db, deckId), []);
const front = ref("");
const back = ref("");

const dueCount = computed(() => dueCards.value.length);

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
  <div v-if="deck">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold">{{ deck.name }}</h2>
      <router-link
        :to="`/review/${deckId}`"
        v-if="dueCount > 0"
        class="bg-indigo-600 text-white px-4 py-2 rounded text-sm"
        data-testid="start-review-btn"
      >
        Review ({{ dueCount }} due)
      </router-link>
    </div>

    <form @submit.prevent="addCard" class="space-y-2 mb-4">
      <input
        v-model="front"
        placeholder="Front (question)"
        class="w-full border rounded px-3 py-2"
        data-testid="card-front-input"
      />
      <input
        v-model="back"
        placeholder="Back (answer)"
        class="w-full border rounded px-3 py-2"
        data-testid="card-back-input"
      />
      <button
        type="submit"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
        data-testid="add-card-btn"
      >
        Add Card
      </button>
    </form>

    <p v-if="cards.length === 0" class="text-gray-500" data-testid="no-cards">
      No cards yet. Add one above.
    </p>

    <ul class="space-y-2">
      <li
        v-for="card in cards"
        :key="card.id"
        class="p-3 bg-white rounded shadow"
        :data-testid="`card-${card.id}`"
      >
        <div class="font-medium">{{ card.front }}</div>
        <div class="text-gray-500 text-sm">{{ card.back }}</div>
      </li>
    </ul>

    <router-link to="/" class="inline-block mt-4 text-indigo-600 text-sm">
      &larr; Back to decks
    </router-link>
  </div>
</template>
