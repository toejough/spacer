<template>
  <q-page>
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
      <template v-for="(note, index) in notes.toReversed()" :key="note">
        <q-item>
          <q-item-section>
            <q-card>
              <q-card-section v-if="openEditorFor != index" horizontal class="flex justify-between">
                <q-card-section @click="openEditor(index)">
                  <span v-sanitize.inline="note"> </span>
                </q-card-section>
                <q-card-actions>
                  <q-btn @click=" remove(index)" round dense flat icon="remove" />
                </q-card-actions>
              </q-card-section>
              <q-card-section v-else horizontal>
                <q-editor v-model="editorContent" min-height="5rem" class="col" v-on-click-outside="closeEditor" />
              </q-card-section>
            </q-card>
          </q-item-section>
        </q-item>
      </template>
    </q-list>
    <q-list>
      <draggable :list="draggableNotes" item-key="id" animation=200>
        <template #item="{ element }">
          <q-item>
            <q-item-section>
              <q-card>
                <q-card-section horizontal class="flex justify-between" v-if="draggableClicked != element.id">
                  <q-card-section @click="draggableClicked = element.id" v-sanitize:inline="element.content" />
                  <q-card-actions>
                    <q-btn @click="removeDraggable(element.id)" round dense flat icon="remove" />
                  </q-card-actions>
                </q-card-section>
                <q-card-section v-else horizontal>
                  <q-editor v-model="element.content" min-height="5rem" class="col"
                    v-on-click-outside="closeDraggableEditor" />
                </q-card-section>
              </q-card>
            </q-item-section>
          </q-item>
        </template>
      </draggable>
    </q-list>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useStorage } from '@vueuse/core'
// TODO: can I do this with just on-blur?
import { vOnClickOutside } from '@vueuse/components'
import draggable from "vuedraggable";
import { uid } from 'quasar';

const draggableClicked = ref("")
const notes = useStorage("notes", <string[]>[])
const newItem = ref("")
const update = () => {
  draggableNotes.value.unshift({ id: uid(), content: newItem.value })
  newItem.value = ""
};
const remove = (index: number) => {
  const reversed = notes.value.length - index - 1
  notes.value = notes.value.filter((_, i) => { return reversed != i })
};
const openEditorFor = ref(-1)
const editorContent = ref("")
const openEditor = (index: number) => {
  openEditorFor.value = index
  editorContent.value = notes.value.toReversed()[index] || ""
};
const closeEditor = () => {
  notes.value[notes.value.length - openEditorFor.value - 1] = editorContent.value
  openEditorFor.value = -1
};
const closeDraggableEditor = () => {
  draggableClicked.value = ""
};
type draggableNote = {
  id: string;
  content: string;
};
const draggableNotes = useStorage("draggableNotes", [] as draggableNote[])
const removeDraggable = (id: string) => {
  const index = draggableNotes.value.findIndex((item) => item.id === id);
  if (index !== -1) {
    draggableNotes.value.splice(index, 1);
  }
}
</script>
