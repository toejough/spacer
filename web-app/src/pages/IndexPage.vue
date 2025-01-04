<template>
  <q-page>
    <q-tabs v-model="tabs">
      <q-tab name="notes" label="notes" />
      <q-tab name="flashcards" label="flashcards" @click="checkCards" />
    </q-tabs>
    <q-tab-panels v-model="tabs">
      <q-tab-panel name="notes">
        <NoteList v-model:notes="notes" v-model:flashcards="flashcards" />
      </q-tab-panel>
      <q-tab-panel name="flashcards">
        <FlashcardList v-model="flashcards" />
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useStorage } from '@vueuse/core'
import NoteList from '../components/NoteList.vue'
import type { draggableNote } from '../components/NoteList.vue'
import FlashcardList from '../components/FlashcardList.vue'
import type { flashcard } from '../components/FlashcardList.vue'

// dragging with sortablejs
// Tabs
const tabs = ref("notes")

// Flashcards: data
// TODO: split out flashcards & results & due dates? these are different concepts...

const checkCards = () => {
  console.log("clicked flashcards")
  console.dir(flashcards)
  // check all the cards are for notes that still exist
  flashcards.value = flashcards.value.filter(element => {
    const index = notes.value.findIndex(note => { return note.id == element.noteID })
    return index >= 0
  })
  flashcards.value.forEach(e => {
    e.show = e.show === undefined ? false : e.show
    e.due = e.due === undefined ? new Date() : e.due
    e.fibDays = e.fibDays === undefined ? 0 : e.fibDays
    e.forgetfulness = e.forgetfulness === undefined ? 1 : e.forgetfulness
  });
  flashcards.value.sort((a: flashcard, b: flashcard): number => { return (new Date(a.due)).getTime() - (new Date(b.due)).getTime() })
};

const notes = useStorage("draggableNotes", [] as draggableNote[])
const flashcards = useStorage("flashcards", [] as flashcard[])
</script>

<style lang="sass">
</style>
