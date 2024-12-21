<template>
  <q-page>
    <q-list>
      <q-item>
        <q-item-section>
          <q-card>
            <q-input v-model="newItem" @keyup.enter="update" placeholder="Enter a new note here">
              <template v-slot:append>
                <q-btn @click="update" round dense flat icon="add" />
              </template></q-input>
          </q-card>
        </q-item-section>
      </q-item>
      <template v-for="(note, index) in notes" :key="note">
        <q-item>
          <q-item-section>
            <q-card>
              <q-card-section v-if="openEditorFor != index" horizontal class="flex justify-between">
                <q-card-section @click="openEditor(index)">
                  {{ note }}
                </q-card-section>
                <q-card-actions>
                  <q-btn @click="remove(index)" round dense flat icon="remove" />
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
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useStorage } from '@vueuse/core'
import { vOnClickOutside } from '@vueuse/components'

const notes = useStorage("notes", <string[]>[])
const newItem = ref("")
const update = () => {
  notes.value.push(newItem.value)
  newItem.value = ""
};
const remove = (index: number) => {
  notes.value = notes.value.filter((_, i) => { return index != i })
};
const openEditorFor = ref(-1)
const editorContent = ref("")
const openEditor = (index: number) => {
  openEditorFor.value = index
  editorContent.value = notes.value[index] || ""
};
const closeEditor = () => {
  notes.value[openEditorFor.value] = editorContent.value
  openEditorFor.value = -1
};
</script>
