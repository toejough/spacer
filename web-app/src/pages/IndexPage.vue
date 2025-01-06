<template>
  <q-page>
    <q-tabs v-model="tabs">
      <q-tab name="notes" label="notes" />
      <!-- <q-tab name="flashcards" label="flashcards" /> -->
      <q-tab name="flashcards" label="flashcards" @click="checkCards" />
    </q-tabs>
    <q-tab-panels v-model="tabs">
      <q-tab-panel name="notes">
        <!-- When a note is clicked to edit, show that note's parent list at the top level, with the parents listed above -->
        <!-- to start with, just a focus button to do that? -->
        <NoteList v-model:notes="notes" v-model:flashcards="flashcards" v-model:listIDs="topLevelNoteIDs" />
      </q-tab-panel>
      <q-tab-panel name="flashcards">
        <FlashcardList v-model:flashcards="flashcards" v-model:notes="notes" />
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
import { uid } from 'quasar';

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
notes.value.forEach(note => {
  if (note.id === undefined) { note.id = uid() }
  if (note.content === undefined) { note.content = "" }
  if (note.flashcards === undefined) { note.flashcards = [] as flashcard[] }
  if (note.subnoteIDs === undefined) { note.subnoteIDs = [] as string[] }
});
const defaultIDs = notes.value.map(n => n.id)
const topLevelNoteIDs = useStorage("topLevelNoteIDs", defaultIDs)
console.dir(notes)
console.dir(topLevelNoteIDs)
const subIDsForID = (id: string): string[] => {
  const noteIndex = notes.value.findIndex(n => n.id == id)
  return [id, ...notes.value[noteIndex]?.subnoteIDs.flatMap(subIDsForID) ?? [] as string[]]
};
const allNoteIDs = notes.value.map(n => n.id)
const usedNoteIDs = topLevelNoteIDs.value.flatMap(subIDsForID)
console.log(allNoteIDs)
console.log(usedNoteIDs)
const missingIDs = allNoteIDs.filter(id => !usedNoteIDs.includes(id))
console.log(missingIDs)
const flashcards = useStorage("flashcards", [] as flashcard[])
topLevelNoteIDs.value.push(...missingIDs)
const extraIDs = topLevelNoteIDs.value.filter(id => !defaultIDs.includes(id))
console.log(extraIDs)
extraIDs.forEach(id => {
  const index = topLevelNoteIDs.value.indexOf(id)
  topLevelNoteIDs.value.splice(index, 1)
});
</script>

<style lang="sass">
</style>
