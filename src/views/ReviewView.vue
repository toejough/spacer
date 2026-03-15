<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { db, type Card } from "../db";
import { sm2 } from "../sm2";

const route = useRoute();
const router = useRouter();
const deckId = Number(route.params.deckId);

const dueCards = ref<Card[]>([]);
const currentIndex = ref(0);
const flipped = ref(false);
const done = ref(false);

const current = computed(() => dueCards.value[currentIndex.value]);

async function load() {
  const now = new Date();
  const all = await db.cards.where("deckId").equals(deckId).toArray();
  dueCards.value = all.filter((c) => c.nextReview <= now);
  if (dueCards.value.length === 0) {
    done.value = true;
  }
}

function flip() {
  flipped.value = true;
}

async function rate(quality: number) {
  const card = current.value;
  if (!card) return;

  const result = sm2(
    {
      easeFactor: card.easeFactor,
      interval: card.interval,
      repetitions: card.repetitions,
    },
    quality
  );

  await db.cards.update(card.id, {
    easeFactor: result.easeFactor,
    interval: result.interval,
    repetitions: result.repetitions,
    nextReview: result.nextReview,
  });

  flipped.value = false;
  if (currentIndex.value < dueCards.value.length - 1) {
    currentIndex.value++;
  } else {
    done.value = true;
  }
}

onMounted(load);
</script>

<template>
  <div>
    <div v-if="done" data-testid="review-done">
      <h2 class="text-lg font-semibold mb-2">Review Complete</h2>
      <p class="text-gray-600 mb-4">
        Reviewed {{ dueCards.length }} card{{ dueCards.length === 1 ? "" : "s" }}.
      </p>
      <button
        @click="router.push(`/deck/${deckId}`)"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
      >
        Back to Deck
      </button>
    </div>

    <div v-else-if="current">
      <p class="text-sm text-gray-500 mb-2" data-testid="review-progress">
        Card {{ currentIndex + 1 }} of {{ dueCards.length }}
      </p>

      <div class="bg-white rounded shadow p-6 min-h-[200px] flex items-center justify-center">
        <div class="text-center">
          <p class="text-xl" data-testid="card-front">{{ current.front }}</p>
          <p
            v-if="flipped"
            class="mt-4 text-gray-600 border-t pt-4"
            data-testid="card-back"
          >
            {{ current.back }}
          </p>
        </div>
      </div>

      <div v-if="!flipped" class="mt-4 text-center">
        <button
          @click="flip"
          class="bg-indigo-600 text-white px-6 py-2 rounded"
          data-testid="flip-btn"
        >
          Show Answer
        </button>
      </div>

      <div v-else class="mt-4 flex justify-center gap-2">
        <button
          v-for="q in [1, 2, 3, 4, 5]"
          :key="q"
          @click="rate(q)"
          class="px-4 py-2 rounded border hover:bg-gray-100"
          :data-testid="`rate-${q}`"
        >
          {{ q }}
        </button>
      </div>
    </div>
  </div>
</template>
