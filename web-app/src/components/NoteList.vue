<template>
  <div v-if="focused === undefined">
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
      <Sortable :list="listNotes" item-key="id" :options="{ animation: '500', handle: '.handle', group: 'notes' }"
        @add="onAdd" @remove="onRemove" @update="onUpdate">
        <template #item="{ element: note }">
          <TransitionGroup name="drag">
            <q-item :key=note.id :data-note-id=note.id>
              <q-item-section>
                <q-card>
                  <q-card-section horizontal class="flex justify-between items-center"
                    v-if="draggableClicked != note.id">
                    <q-card-section>
                      <q-icon name="drag_indicator" class="handle" />
                    </q-card-section>
                    <q-card-section @click="editorOpenedOnNote(note.id)" class="flex col">
                      <div v-sanitize:inline="note.content" />
                    </q-card-section>
                    <q-card-actions>
                      <q-btn v-if="note.subnoteIDs.length == 0" @click="removeDraggable(note.id)" round dense flat
                        icon="remove" />
                      <div v-else> ({{ note.subnoteIDs.length }} subnotes) </div>
                    </q-card-actions>
                  </q-card-section>
                  <div v-else v-on-click-outside="closeDraggableEditor">
                    <q-card-section horizontal>
                      <q-card-section class="col">
                        <button @click="toggleFlashCard" class="button-style">
                          <q-icon name="flash_on" />
                          Toggle flashcard with BOLD
                        </button>
                      </q-card-section>
                      <q-card-section>
                        <q-btn label="Focus" @click="focused = note" />
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
                    <q-separator />
                    <q-card-section>
                      <q-card-section horizontal>
                        <q-card-section class="col">
                          <div class="subnotes">Subnotes</div>
                        </q-card-section>
                      </q-card-section>
                      <NoteList v-model:notes="notes" v-model:flashcards="flashcards"
                        v-model:listIDs="note.subnoteIDs" />
                    </q-card-section>
                  </div>
                </q-card>
              </q-item-section>
            </q-item>
          </TransitionGroup>
        </template>
      </Sortable>
    </q-list>
  </div>
  <div v-else>
    <q-list>
      <q-item>
        <q-item-section>
          <q-card>
            <q-card-section horizontal>
              <q-card-section class="col">
                <div v-sanitize:inline="focused.content" />
              </q-card-section>
              <q-card-section>
                <q-btn label="Unfocus" @click="focused = undefined" />
              </q-card-section>
            </q-card-section>
          </q-card>
        </q-item-section>
      </q-item>
    </q-list>
    <NoteList v-model:notes="notes" v-model:flashcards="flashcards" v-model:listIDs="focused.subnoteIDs" />
  </div>
</template>

<script lang="ts">
export type draggableNote = {
  id: string;
  content: string;
  flashcards: flashcard[];
  subnoteIDs: string[];
};
</script>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Ref } from 'vue'
import { vOnClickOutside } from '@vueuse/components'
import { uid } from 'quasar';
import { Sortable } from "sortablejs-vue3";
import type { flashcard } from './FlashcardList.vue'
import type { SortableEvent } from "sortablejs";

const focused: Ref<undefined | draggableNote> = ref(undefined)

const notes = defineModel<draggableNote[]>('notes', { required: true })
const flashcards = defineModel<flashcard[]>('flashcards', { required: true })
const listIDs = defineModel<string[]>('listIDs', { required: true })
const listNotes = computed(() => {
  return listIDs.value.map(id => {
    const noteIndex = notes.value.findIndex(n => n.id == id)
    const note = notes.value[noteIndex]
    if (note !== undefined) {
      if (note.id === undefined) { note.id = uid() }
      if (note.content === undefined) { note.content = "" }
      if (note.flashcards === undefined) { note.flashcards = [] as flashcard[] }
      if (note.subnoteIDs === undefined) { note.subnoteIDs = [] as string[] }
    }
    return note
  }).filter(n => n !== undefined)
})
// Notes: data
// Notes: Add/remove note
const newItem = ref("")
const update = () => {
  const newNote = {
    id: uid(), content: newItem.value, flashcards: [] as flashcard[], subnoteIDs: [] as string[]
  }
  notes.value.unshift(newNote)
  listIDs.value.unshift(newNote.id)
  newItem.value = ""
};
const removeDraggable = (id: string) => {
  {
    const index = notes.value.findIndex((item) => item.id === id);
    if (index !== -1 && notes.value[index] != undefined) {
      removeCardsFrom(flashcards.value, notes.value[index].flashcards)
      notes.value.splice(index, 1);
    }
  }
  {
    const index = listIDs.value.findIndex((item) => item === id);
    if (index !== -1 && listIDs.value[index] != undefined) {
      listIDs.value.splice(index, 1);
    }
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
  const index = notes.value.findIndex((item) => item.id === noteId);
  const note = notes.value[index];
  // fix the flashcards, as necessary
  if (note != null) {
    ensureCardsForNote(note)
    if (note.subnoteIDs === undefined) {
      note.subnoteIDs = [] as string[]
    }
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
  console.log("dead cards: ")
  const deadCards = diffCards(globalCardsForThisNote, note.flashcards)
  console.dir(deadCards)
  removeCardsFrom(flashcards.value, deadCards)
  // TODO: make adding to the list preserve uniqueness
  console.log("before unique:")
  console.dir(flashcards)
  reduceToUnique(flashcards.value)
  console.log("after unique:")
  console.dir(flashcards)
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

  const index = notes.value.findIndex((item) => item.id === draggableClicked.value);
  const note = notes.value[index];
  if (note != null) {
    ensureCardsForNote(note)
  }

};

// type noteMoveEvent = {
//   item: draggableNote  // dragged HTMLElement
//   to: draggableNote[]    // target list
//   from: draggableNote[]  // previous list
//   oldIndex: number // element's old index within old parent
//   newIndex: number // element's new index within new parent
// }

const onRemove = (evt: SortableEvent) => {
  const item = evt.item;  // dragged HTMLElement
  const prevIndex = evt.oldIndex as number;  // element's old index within old parent

  console.log("removing...")
  console.log(item.dataset.noteId)
  console.dir(listNotes)
  // draggableNotes.value.splice(prevIndex, 1)
  listIDs.value.splice(prevIndex, 1)
  const noteIndex = notes.value.findIndex(n => n.id == item.dataset.noteId)
  notes.value.splice(noteIndex, 1)
  item.remove();
  console.log("...removed")
  console.dir(listNotes)
};

const onAdd = (evt: SortableEvent) => {
  const item = evt.item;  // dragged HTMLElement
  const newIndex = evt.newIndex as number;

  console.log("adding...")
  console.log(item.dataset.noteId)
  const noteId = item.dataset.noteId || ""
  const noteIndex = notes.value.findIndex(n => n.id == noteId)
  if (noteIndex < 0) {
    console.log("missing note! could not find note for id:")
    console.log(noteId)
    return
  }
  console.dir(item.dataset)
  console.dir(listNotes)
  // draggableNotes.value.splice(newIndex, 0, note)
  const note = notes.value[noteIndex]
  if (note === undefined) {
    console.log("undefined note at index")
    console.log(noteIndex)
    return
  }
  listIDs.value.splice(newIndex, 0, noteId)
  notes.value.push(note)
  console.log("...added")
  console.dir(listNotes)
};

const onUpdate = (evt: SortableEvent) => {
  const item = evt.item;  // dragged HTMLElement
  const prevIndex = evt.oldIndex as number;  // element's old index within old parent

  console.log("removing...")
  console.log(item.dataset.noteId)
  console.dir(listNotes)
  // draggableNotes.value.splice(prevIndex, 1)
  listIDs.value.splice(prevIndex, 1)
  const noteIndex = notes.value.findIndex(n => n.id == item.dataset.noteId)
  const note = notes.value[noteIndex]
  if (note === undefined) {
    console.log("undefined note at index")
    console.log(noteIndex)
    return
  }
  console.log("...removed")
  console.dir(listNotes)
  const newIndex = evt.newIndex as number;

  console.log("adding...")
  console.log(item.dataset.noteId)
  const noteId = item.dataset.noteId || ""
  if (noteIndex < 0) {
    console.log("missing note! could not find note for id:")
    console.log(noteId)
    return
  }
  console.dir(item.dataset)
  console.dir(listNotes)
  // draggableNotes.value.splice(newIndex, 0, note)
  listIDs.value.splice(newIndex, 0, noteId)
  console.log("...added")
  console.dir(listNotes)
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
.subnotes
  color: $grey-6
  font-weight: bold
</style>
