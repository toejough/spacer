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
type flashcard = {
  answer: string
  prompt: string
};
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
  draggableNotes.value.unshift({ id: uid(), content: newItem.value, flashcards: [] })
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
  draggableClicked.value = noteId
  const index = draggableNotes.value.findIndex((item) => item.id === noteId);
  const note = draggableNotes.value[index];
  if (note != null) {
    const regexp = /<b>(.*?)<\/b>/g
    const array = [...note.content.matchAll(regexp)];
    const flashcards = array.map((value) => {
      const input = value.input;
      const index = value.index;
      const answer = value[1] || "";
      const beginning = input.slice(0, index + 3)
      const blank = "_".repeat(answer.length)
      const end = input.slice(index + 3 + answer.length)
      const prompt = beginning + blank + end
      return { prompt: prompt, answer: answer } as flashcard
    })
    note.flashcards = flashcards
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
    const regexp = /<b>(.*?)<\/b>/g
    const array = [...note.content.matchAll(regexp)];
    const flashcards = array.map((value) => {
      const input = value.input;
      const index = value.index;
      const answer = value[1] || "";
      const beginning = input.slice(0, index + 3)
      const blank = "_".repeat(answer.length)
      const end = input.slice(index + 3 + answer.length)
      const prompt = beginning + blank + end
      return { prompt: prompt, answer: answer } as flashcard
    })
    note.flashcards = flashcards
    // TODO: convert answer/prompt to flashcard
    // TODO: just add answers/prompts to a flashcards list
    // TODO: add a flashcards tab
    // TODO: add flashcard functionality: prompt, answer, remembered, forgot
    // TODO: add spaced repetition logic
    // TODO: add section for due vs not
    // TODO: add notifications
  }

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
