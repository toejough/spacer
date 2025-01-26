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
  // define any undefined fields
  if (note.id === undefined) { note.id = uid() }
  if (note.content === undefined) { note.content = "" }
  if (note.flashcards === undefined) { note.flashcards = [] as flashcard[] }
  if (note.parentNoteID === undefined) { note.parentNoteID = "" }
  if (note.subnoteIDs === undefined) { note.subnoteIDs = [] as string[] }
  // true up the parent ID's based on subnote ID's
  note.subnoteIDs.forEach(id => {
    // get the subnote & set parent ID
    const subnote = notes.value.find(value => value.id == id)
    if (subnote != undefined) { subnote.parentNoteID = note.id }
  })
});

// get the top level id's
const topLevelNoteIDs = notes.value.filter(note => note.parentNoteID == "").map(note => note.id)

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
