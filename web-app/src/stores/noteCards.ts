import { defineStore, acceptHMRUpdate } from 'pinia';
import { useStorage } from '@vueuse/core'
import type { flashcard } from '../components/FlashcardList.vue'

export type draggableNote = {
  id: string;
  content: string;
  flashcards: flashcard[];
  parentNoteID: string;
  subnoteIDs: string[];
};

export const useNoteCardStore = defineStore('myStore', {
  state: () => ({
    noteCards: useStorage("draggableNotes", [] as draggableNote[]),
  }),
  getters: {},
  actions: {}
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNoteCardStore, import.meta.hot));
}
