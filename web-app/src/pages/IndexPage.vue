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
import FlashcardList from '../components/FlashcardList.vue'
import type { flashcard } from '../components/FlashcardList.vue'
import { uid } from 'quasar';
import { useNoteCardStore } from 'src/stores/noteCards';
import { storeToRefs } from 'pinia';

// STORE
const store = useNoteCardStore();

// NOTECARDS
// TODO: push this data into the noteList component
const notes = storeToRefs(store).noteCards

// Sanitizing noteCard data
// TODO: push this into a store getter
notes.value.forEach(note => {
  if (note.id === undefined) { note.id = uid() }
  if (note.content === undefined) { note.content = "" }
  if (note.flashcards === undefined) { note.flashcards = [] as flashcard[] }
  if (note.subnoteIDs === undefined) { note.subnoteIDs = [] as string[] }
});

// top level notes
// TODO: these should just be all the notes with no parents. Add parent id's and then calculate this.
// TODO: put these in the store?
const defaultIDs = notes.value.map(n => n.id)
const topLevelNoteIDs = useStorage("topLevelNoteIDs", defaultIDs)

// find notecards that are in the list but are not present in the tree of top-level ID's down.
const allNoteIDs = notes.value.map(n => n.id)
const subIDsForID = (id: string): string[] => {
  const noteIndex = notes.value.findIndex(n => n.id == id)
  return [id, ...notes.value[noteIndex]?.subnoteIDs.flatMap(subIDsForID) ?? [] as string[]]
};
const usedNoteIDs = topLevelNoteIDs.value.flatMap(subIDsForID)
const missingIDs = allNoteIDs.filter(id => !usedNoteIDs.includes(id))
// push the notecards that are not in the tree into the top level
topLevelNoteIDs.value.push(...missingIDs)
// find the top level id's that are not in the defaultID's list
const extraIDs = topLevelNoteIDs.value.filter(id => !defaultIDs.includes(id))
// remove them from the top level?? isn't this just undoing what we just did? all to true up the top level?
// TODO: replace all of this nonsense with finding the notes with no parents & using them as the top level list.
extraIDs.forEach(id => {
  const index = topLevelNoteIDs.value.indexOf(id)
  topLevelNoteIDs.value.splice(index, 1)
});

// Tabs
const tabs = ref("notes")

// Top level Flashcards
// TODO: put these into pinia
// TODO: split out flashcards & results & due dates? these are different concepts...
// TODO: push this data into the flashcard list component
const flashcards = useStorage("flashcards", [] as flashcard[])

const checkCards = () => {
  console.log("clicked flashcards")
  console.dir(flashcards)
  // check all the cards are for notes that still exist
  flashcards.value = flashcards.value.filter(element => {
    const index = notes.value.findIndex(note => { return note.id == element.noteID })
    return index >= 0
  })
  // sanitize remaining flashcard data
  // TODO: push this into a store getter
  flashcards.value.forEach(e => {
    e.show = e.show === undefined ? false : e.show
    e.due = e.due === undefined ? new Date() : e.due
    e.fibDays = e.fibDays === undefined ? 0 : e.fibDays
    e.forgetfulness = e.forgetfulness === undefined ? 1 : e.forgetfulness
  });
  // sort cards by due date
  flashcards.value.sort((a: flashcard, b: flashcard): number => { return (new Date(a.due)).getTime() - (new Date(b.due)).getTime() })
};

</script>

<style lang="sass">
</style>
