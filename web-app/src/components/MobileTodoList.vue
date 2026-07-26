<template>
  <div>
    <h3>Todos (mobile demo)</h3>
    <div role="list">
      <MobileTodoCard v-for="t in todos" :key="t.id" :todo="t" @done="onDone" @abandon="onAbandon" @reopen="onReopen" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MobileTodoCard from './MobileTodoCard.vue'
import { uid } from 'quasar'

const todos = ref([{
  id: '1', title: 'Write project spec', status: 'open'
},{ id: '2', title: 'Refactor API', status: 'open' },
{ id: '3', title: 'Publish release notes', status: 'done' },
{ id: '4', title: 'Experimental feature X', status: 'abandoned' }
])

function onDone(todo){
  const t = todos.value.find(x=>x.id===todo.id);
  if(t){ t.status='done'; t.completed_at = (new Date()).toISOString(); }
}
function onAbandon(todo){
  const t = todos.value.find(x=>x.id===todo.id);
  if(t){ t.status='abandoned'; t.archived_at = (new Date()).toISOString(); }
}
function onReopen(todo){
  const t = todos.value.find(x=>x.id===todo.id);
  if(t){ t.status='open'; t.completed_at=null; t.archived_at=null }
}
</script>

<style scoped>
h3{margin:8px 0}
</style>
