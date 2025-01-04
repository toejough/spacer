<template>
  <q-list>
    <q-item>
      <q-item-section>
        <q-card>
          <q-input filled v-model="newItem" @keyup.enter="update" placeholder="Enter a new note here">
            <template v-slot:append>
              <q-btn @click="update" round dense flat icon="add" />
            </template></q-input>
        </q-card>
      </q-item-section>
    </q-item>
    <Sortable :list="draggableNotes" item-key="id" :options="{ animation: '500', handle: '.handle' }">
      <template #item="{ element: note }">
        <TransitionGroup name="drag">
          <q-item :key=note.id>
            <q-item-section>
              <q-card>
                <q-card-section horizontal class="flex justify-between items-center" v-if="draggableClicked != note.id">
                  <q-card-section>
                    <q-icon name="drag_indicator" class="handle" />
                  </q-card-section>
                  <q-card-section @click="editorOpenedOnNote(note.id)" class="flex col">
                    <div v-sanitize:inline="note.content" />
                  </q-card-section>
                  <q-card-actions>
                    <q-btn @click="removeDraggable(note.id)" round dense flat icon="remove" />
                  </q-card-actions>
                </q-card-section>
                <div v-else v-on-click-outside="closeDraggableEditor">
                  <q-card-section horizontal>
                    <q-card-section>
                      <button @click="toggleFlashCard" class="button-style">
                        <q-icon name="flash_on" />
                        Toggle flashcard with BOLD
                      </button>
                    </q-card-section>
                  </q-card-section>
                  <q-separator />
                  <q-card-section horizontal>
                    <q-editor v-model="note.content" min-height="5rem" class="col" :toolbar="[]" />
                  </q-card-section>
                  <q-separator />
                  <q-card-section horizontal v-for="flashcard in note.flashcards" :key="flashcard">
                    <q-card-section class="col">
                      <div v-sanitize:inline="flashcard.prompt" />
                    </q-card-section>
                    <q-card-section>
                      (<span v-sanitize:inline="flashcard.answer" />)
                    </q-card-section>
                  </q-card-section>
                </div>
              </q-card>
            </q-item-section>
          </q-item>
        </TransitionGroup>
      </template>
    </Sortable>
  </q-list>
</template>

<script lang="ts">
export type draggableNote = {
  id: string;
  content: string;
  flashcards: flashcard[];
};
</script>

<script setup lang="ts">
import { ref } from 'vue'
import { vOnClickOutside } from '@vueuse/components'
import { uid } from 'quasar';
import { Sortable } from "sortablejs-vue3";
import type { flashcard } from './FlashcardList.vue'

const draggableNotes = defineModel<draggableNote[]>('notes', { required: true })
const flashcards = defineModel<flashcard[]>('flashcards', { required: true })
// Notes: data
// Notes: Add/remove note
const newItem = ref("")
const update = () => {
  draggableNotes.value.unshift({
    id: uid(), content: newItem.value, flashcards: [] as flashcard[]
  })
  newItem.value = ""
};
const removeDraggable = (id: string) => {
  const index = draggableNotes.value.findIndex((item) => item.id === id);
  if (index !== -1 && draggableNotes.value[index] != undefined) {
    removeCardsFrom(flashcards.value, draggableNotes.value[index].flashcards)
    draggableNotes.value.splice(index, 1);
  }
}

// Notes: Open/close editor
const draggableClicked = ref("")
const closeDraggableEditor = () => {
  draggableClicked.value = ""
};
const editorOpenedOnNote = (noteId: string) => {
  // get the note
  draggableClicked.value = noteId
  const index = draggableNotes.value.findIndex((item) => item.id === noteId);
  const note = draggableNotes.value[index];
  // fix the flashcards, as necessary
  if (note != null) {
    ensureCardsForNote(note)
  }
};

const ensureCardsForNote = (note: draggableNote) => {
  // identify the correct flashcards
  const regexp = /<b>(.*?)<\/b>/g
  const array = [...note.content.matchAll(regexp)];
  note.flashcards = array.map((value) => {
    const input = value.input;
    const index = value.index;
    const answer = value[1] || "";
    const beginning = input.slice(0, index + 3)
    const blank = "_".repeat(answer.length)
    const end = input.slice(index + 3 + answer.length)
    const prompt = beginning + blank + end
    return { prompt: prompt, answer: answer, noteID: note.id, id: note.id + prompt + answer, show: false, due: new Date(), fibDays: 0, forgetfulness: 1 } as flashcard
  })
  // if any of these flashcards are not present in the overall list, add them with new id's.
  const newCards = diffCards(note.flashcards, flashcards.value)
  flashcards.value.push(...newCards)
  console.log("adding new cards: ")
  console.dir(newCards)
  // if any of the flashcards in the overall list that are pointed at this note don't match, delete them.
  const globalCardsForThisNote = flashcards.value.filter(card => {
    return card.noteID == note.id
  })
  console.log("global cards for this note: ")
  console.dir(globalCardsForThisNote)
  console.log("fset: ")
  const deadCards = diffCards(globalCardsForThisNote, note.flashcards)
  removeCardsFrom(flashcards.value, deadCards)
  // TODO: make adding to the list preserve uniqueness
  reduceToUnique(flashcards.value)
};

const diffCards = (base: flashcard[], other: flashcard[]): flashcard[] => {
  return base.filter(card => {
    return !other.map(ocard => { return ocard.id }).includes(card.id)
  })
};

const removeCardsFrom = (from: flashcard[], toRemove: flashcard[]) => {
  toRemove.forEach(element => {
    const index = from.map(card => { return card.id }).indexOf(element.id)
    if (index > -1) {
      from.splice(index, 1)
    }
  });
};

const reduceToUnique = (cards: flashcard[]) => {
  // for each card
  for (let i = 0; i < cards.length; i++) {
    // get current card
    const current = cards[i]
    if (current === undefined) { break }
    // search the rest of the list for a match
    let remainingIndex = cards.slice(i + 1).map(card => { return card.id }).indexOf(current.id)
    // if found, remove it
    let found = remainingIndex >= 0
    while (found) {
      cards.splice(i + 1 + remainingIndex, 1)
      remainingIndex = cards.slice(i + 1).map(card => { return card.id }).indexOf(current.id)
      // if found, remove it
      found = remainingIndex >= 0
    }
  }
};

// Flashcard: toggle
const toggleFlashCard = () => {
  // TODO: replace this with the example here: https://jsfiddle.net/y9qzejmf/1/
  // const selection = document.getSelection()
  document.execCommand('bold')

  const index = draggableNotes.value.findIndex((item) => item.id === draggableClicked.value);
  const note = draggableNotes.value[index];
  if (note != null) {
    ensureCardsForNote(note)
  }

};

</script>

<style lang="sass">
.handle
  cursor: grab
.button-style
  background-color: inherit
  border-style: none
  color: $primary
.drag-move
  transition: all 1s cubic-bezier(0.55, 0, 0.1, 1)
</style>
