<template>
  <q-page>
    <q-tabs v-model="tabs">
      <q-tab name="notes" label="notes" />
      <q-tab name="flashcards" label="flashcards" />
    </q-tabs>
    <q-tab-panels v-model="tabs">
      <q-tab-panel name="notes">
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
          <draggable :list="draggableNotes" item-key="id" animation=200 handle=".handle">
            <template #item="{ element }">
              <q-item>
                <q-item-section>
                  <q-card>
                    <q-card-section horizontal class="flex justify-between items-center"
                      v-if="draggableClicked != element.id">
                      <q-card-section>
                        <q-icon name="drag_indicator" class="handle" />
                      </q-card-section>
                      <q-card-section @click="editorOpenedOnNote(element.id)" class="flex col">
                        <div v-sanitize:inline="element.content" />
                      </q-card-section>
                      <q-card-actions>
                        <q-btn @click="removeDraggable(element.id)" round dense flat icon="remove" />
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
                        <q-editor v-model="element.content" min-height="5rem" class="col" :toolbar="[]" />
                      </q-card-section>
                      <q-separator />
                      <q-card-section horizontal v-for="flashcard in element.flashcards" :key="flashcard">
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
            </template>
          </draggable>
        </q-list>
      </q-tab-panel>
      <q-tab-panel name="flashcards">
        <q-list>
          <div v-for="card in flashcards" :key="card.id">
            <q-item>
              <q-item-section>
                <q-card>
                  <q-card-section>
                    <div v-sanitize="card.prompt">
                    </div>
                  </q-card-section>
                </q-card>
              </q-item-section>
            </q-item>
          </div>
        </q-list>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useStorage } from '@vueuse/core'
import { vOnClickOutside } from '@vueuse/components'
import draggable from "vuedraggable";
import { uid } from 'quasar';

// Tabs
const tabs = ref("notes")

// Flashcards: data
// TODO: forget the set logic stuff - going to have to implement it manually instead of with the Set type, because we won't be able to compare other state such as last time we practiced, etc.
// or maybe that's... results? results for a flashcard? Interesting...
type flashcard = {
  id: string
  answer: string
  prompt: string
  noteID: string
};
const flashcards = useStorage("flashcards", [] as flashcard[])

// Notes: data
type draggableNote = {
  id: string;
  content: string;
  flashcards: flashcard[];
};
const draggableNotes = useStorage("draggableNotes", [] as draggableNote[])

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
  if (index !== -1) {
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
    return { prompt: prompt, answer: answer, noteID: note.id, id: note.id + prompt + answer } as flashcard
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

  // TODO: add flashcard functionality: prompt, answer, remembered, forgot
  // TODO: add spaced repetition logic
  // TODO: add section for due vs not
  // TODO: add notifications
  // before doing anything, true up the notes:
  //   if the number of bolded segments != the number of flashcards:
  //     delete the old flashcards
  //     create a new flashcard for each bolded item
  // if the selected text is completely unbolded:
  //   bold the text
  //   insert a new flashcard at the index of the newly bolded text, in the list of bolded text snippets
  // if the selected text is completely bolded:
  //   unbold the text
  //   delete the flashcard at teh old bolded text's index, in the list of bolded text snippets
  // if the selected text overlaps only one bolded item
  //   bold the unbolded part
  //   update the flashcard to contain all the bolded text
  // else (selected text overlaps more than one bolded item)
  //   bold the unbolded part
  //   update the first flashcard to contain all the bolded text
  //   delete the other flashcards
  // todo: add grabbing
};

</script>

<style lang="sass">
.handle
  cursor: grab
.button-style
  background-color: inherit
  border-style: none
  color: $primary
</style>
