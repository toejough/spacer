<template>
  <Sortable :list="flashcards" item-key="id" :options="{ animation: '500', handle: '.handle' }">
    <template #item="{ element: card }">
      <TransitionGroup name="drag">
        <q-item :key=card.id>
          <q-item-section>
            <q-card>
              <q-card-section horizontal class="dueDate">
                <q-card-section>
                  <q-icon name="drag_indicator" class="handle" />
                </q-card-section>
                <q-card-section>
                  <div>Practice Again On/After: {{ date.formatDate(card.due, "YYYY-MM-DD") }}</div>
                </q-card-section>
              </q-card-section>
              <q-card-section v-if="parentNotesForCard(card).length > 0">
                <q-list v-for="note in parentNotesForCard(card)" :key="note.id">
                  <div v-sanitize:inline="note.content" />
                  <q-separator />
                </q-list>
              </q-card-section>
              <div v-if="!card.show">
                <q-card-section>
                  <div v-sanitize="card.prompt"></div>
                </q-card-section>
                <q-card-actions>
                  <q-btn label="Show Answer" @click="card.show = true" />
                </q-card-actions>
              </div>
              <div v-else>
                <q-card-section>
                  <div v-sanitize="card.prompt"></div>
                </q-card-section>
                <q-card-section class="answer">
                  <div v-sanitize="card.answer"></div>
                </q-card-section>
                <q-card-actions>
                  <q-btn label="Remembered" @click="rememberedCard(card)" />
                  <q-btn label="Forgot" @click="forgotCard(card)" />
                </q-card-actions>
              </div>
            </q-card>
          </q-item-section>
        </q-item>
      </TransitionGroup>
    </template>
  </Sortable>
</template>

<script lang="ts">
export type flashcard = {
  id: string
  answer: string
  prompt: string
  noteID: string
  show: boolean
  due: Date
  fibDays: number
  forgetfulness: number
};
// Flashcards: data
// TODO: split out flashcards & results & due dates? these are different concepts...
</script>

<script setup lang="ts">
// todo - onclick outside of card, collapse card again?
import { Sortable } from "sortablejs-vue3";
import { date } from 'quasar'
import type { draggableNote } from '../stores/noteCards.ts'


const flashcards = defineModel<flashcard[]>("flashcards", { required: true })
const notes = defineModel<draggableNote[]>("notes", { required: true })

const parentNotesForCard = (card: flashcard): draggableNote[] => {
  const parentNotes = [] as draggableNote[]
  let noteID = card.noteID
  let noteIndex = notes.value.findIndex(n => n.id == noteID)
  let parentNote = notes.value[noteIndex]
  if (noteIndex >= 0 && parentNote !== undefined) {
    noteID = parentNote.id
    noteIndex = notes.value.findIndex(n => n.subnoteIDs.includes(noteID))
    parentNote = notes.value[noteIndex]
  }
  while (noteIndex >= 0 && parentNote !== undefined) {
    noteID = parentNote.id
    parentNotes.unshift(parentNote)
    noteIndex = notes.value.findIndex(n => n.subnoteIDs.includes(noteID))
    parentNote = notes.value[noteIndex]
  }
  return parentNotes
};
const rememberedCard = (card: flashcard) => {
  card.show = false
  card.fibDays = nextFib(card.fibDays)
  card.due = new Date((new Date()).getTime() + card.fibDays / card.forgetfulness * 1000 * 60 * 60 * 24)
  console.log("new due date:")
  console.log(card.due)
  console.dir(card)
  flashcards.value.sort((a: flashcard, b: flashcard): number => { return (new Date(a.due)).getTime() - (new Date(b.due)).getTime() })
  console.dir(flashcards)
  // TODO: set up reminders
};

const forgotCard = (card: flashcard) => {
  card.show = false
  card.due = new Date()
  card.forgetfulness++
  card.fibDays = 0
  console.log("new forgetfulness:")
  console.log(card.forgetfulness)
  console.dir(card)
  flashcards.value.sort((a: flashcard, b: flashcard): number => { return (new Date(a.due)).getTime() - (new Date(b.due)).getTime() })
  // TODO: set up reminders
};

const nextFib = (currentNum: number): number => {
  let current = 0
  let next = 1
  while (current <= currentNum) {
    const nextNext = current + next
    current = next
    next = nextNext
  }
  return current
};

</script>

<style lang="sass">
.handle
  cursor: grab
.dueDate
  background-color: $grey-3
.answer
  background-color: $cyan-3
.drag-move
  transition: all 1s cubic-bezier(0.55, 0, 0.1, 1)
</style>
