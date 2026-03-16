<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { db, getDeck, getDueCards, updateCardReview, type Card } from "../db";
import { sm2 } from "../sm2";
import { type ViewState } from "../view-state";
import { AlertTriangle, Search } from "lucide-vue-next";

const route = useRoute();
const router = useRouter();
const deckId = Number(route.params.deckId);

const loadState = ref<ViewState<Card[]>>({ status: "loading" });
const currentIndex = ref(0);
const flipped = ref(false);
const done = ref(false);

const cards = computed(() =>
  loadState.value.status === "loaded" ? loadState.value.data : []
);
const current = computed(() => cards.value[currentIndex.value]);

async function load() {
  try {
    const deck = await getDeck(db, deckId);
    if (!deck) {
      loadState.value = { status: "not-found" };
      return;
    }
    const dueCards = await getDueCards(db, deckId);
    if (dueCards.length === 0) {
      loadState.value = { status: "empty" };
      done.value = true;
    } else {
      loadState.value = { status: "loaded", data: dueCards };
    }
  } catch (err) {
    loadState.value = {
      status: "error",
      message: err instanceof Error ? err.message : String(err),
    };
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

  await updateCardReview(db, card.id, result);

  flipped.value = false;
  if (currentIndex.value < cards.value.length - 1) {
    currentIndex.value++;
  } else {
    done.value = true;
  }
}

onMounted(load);
</script>

<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="loadState.status === 'loading'" data-testid="loading-skeleton">
      <div class="h-4 w-24 bg-surface rounded-sm animate-skeleton mb-2" />
      <div class="bg-surface rounded-md min-h-[200px] animate-skeleton mb-4" />
      <div class="flex justify-center">
        <div class="h-10 w-32 bg-surface rounded-sm animate-skeleton" />
      </div>
    </div>

    <!-- Not found -->
    <div
      v-else-if="loadState.status === 'not-found'"
      data-testid="not-found"
      class="text-center py-12"
    >
      <Search class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Deck not found</p>
      <p class="text-text-faint text-sm mb-4">This deck may have been deleted</p>
      <router-link to="/" class="text-primary text-sm">&larr; Back to decks</router-link>
    </div>

    <!-- Error state -->
    <div
      v-else-if="loadState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <AlertTriangle class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm mb-4">{{ loadState.message }}</p>
      <router-link :to="`/deck/${deckId}`" class="text-primary text-sm">
        &larr; Back to deck
      </router-link>
    </div>

    <!-- Review complete / empty -->
    <div v-else-if="done" data-testid="review-done">
      <h2 class="text-lg font-semibold mb-2">Review Complete</h2>
      <p class="text-text-muted mb-4">
        Reviewed {{ cards.length }} card{{ cards.length === 1 ? "" : "s" }}.
      </p>
      <button
        @click="router.push(`/deck/${deckId}`)"
        class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
      >
        Back to Deck
      </button>
    </div>

    <!-- Active review -->
    <div v-else-if="current">
      <p class="text-sm text-text-muted mb-2" data-testid="review-progress">
        Card {{ currentIndex + 1 }} of {{ cards.length }}
      </p>

      <div class="bg-surface rounded-md shadow-card p-6 min-h-[200px] flex items-center justify-center">
        <div class="text-center">
          <p class="text-xl" data-testid="card-front">{{ current.front }}</p>
          <p
            v-if="flipped"
            class="mt-4 text-text-muted border-t border-border pt-4"
            data-testid="card-back"
          >
            {{ current.back }}
          </p>
        </div>
      </div>

      <div v-if="!flipped" class="mt-4 text-center">
        <button
          @click="flip"
          class="bg-accent hover:bg-accent-hover text-white px-6 py-2 rounded-sm"
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
          class="px-4 py-2 rounded-sm border border-border hover:bg-surface-raised"
          :data-testid="`rate-${q}`"
        >
          {{ q }}
        </button>
      </div>
    </div>
  </div>
</template>
