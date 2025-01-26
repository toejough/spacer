<template>
  <q-card>
    <q-card-section v-if="hasAncestors(getFocusedNote())">
      Parent NoteCards:
      <q-list v-for="note in getAncestorsOf(getFocusedNote())" :key="note.id">
        <q-separator />
        <div v-sanitize:inline="note.content" />
      </q-list>
    </q-card-section>
  </q-card>
</template>

<script lang="ts">
</script>

<script setup lang="ts">
import { useNoteCardStore } from 'src/stores/noteCards';
import { storeToRefs } from 'pinia';
import type { draggableNote } from "../stores/noteCards.ts"


// TODO: organize this file
// TODO: when clicking ancestor card, select parent

const draggableClicked = storeToRefs(useNoteCardStore()).clicked
const notes = storeToRefs(useNoteCardStore()).noteCards

const hasAncestors = (note: draggableNote | undefined) => {
  if (note === undefined) { return false }
  if (note.parentNoteID == "") { return false }
  return true
};

const getFocusedNote = () => {
  if (draggableClicked.value == "") { return undefined }
  return notes.value.find(note => note.id == draggableClicked.value)
};

const getAncestorsOf = (note: draggableNote | undefined) => {
  const ancestors = [] as draggableNote[]
  while (note?.parentNoteID != "") {
    note = notes.value.find(n => n.id == note?.parentNoteID)
    if (note !== undefined) {
      ancestors.unshift(note)
    }
  }
  return ancestors
};

</script>

<style lang="sass">
</style>
