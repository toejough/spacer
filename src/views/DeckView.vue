<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { db, type Card, type Deck } from "../db";

const route = useRoute();
const deckId = Number(route.params.id);

const deck = ref<Deck | undefined>();
const cards = ref<Card[]>([]);
const front = ref("");
const back = ref("");

async function load() {
  deck.value = await db.decks.get(deckId);
  cards.value = await db.cards.where("deckId").equals(deckId).toArray();
}

async function addCard() {
  const f = front.value.trim();
  const b = back.value.trim();
  if (!f || !b) return;
  await db.cards.add({
    deckId,
    front: f,
    back: b,
    easeFactor: 2.5,
    interval: 0,
    repetitions: 0,
    nextReview: new Date(),
  } as Card);
  front.value = "";
  back.value = "";
  await load();
}

function dueCount() {
  const now = new Date();
  return cards.value.filter((c) => c.nextReview <= now).length;
}

onMounted(load);
</script>

<template>
  <div v-if="deck">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold">{{ deck.name }}</h2>
      <router-link
        :to="`/review/${deckId}`"
        v-if="dueCount() > 0"
        class="bg-indigo-600 text-white px-4 py-2 rounded text-sm"
        data-testid="start-review-btn"
      >
        Review ({{ dueCount() }} due)
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
