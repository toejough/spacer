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
import { vOnClickOutside } from '@vueuse/components'
import draggable from "vuedraggable";
import { uid } from 'quasar';

// Notes: data
type draggableNote = {
  id: string;
  content: string;
};
const draggableNotes = useStorage("draggableNotes", [] as draggableNote[])

// Notes: Add/remove note
const newItem = ref("")
const update = () => {
  draggableNotes.value.unshift({ id: uid(), content: newItem.value })
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
</script>
