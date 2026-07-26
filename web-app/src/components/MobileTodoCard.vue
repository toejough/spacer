<template>
  <div :class="['card', stateClass]" role="listitem" :aria-label="`todo ${todo.title}`">
    <div>
      <div class="title">{{ todo.title }}</div>
      <div class="small">{{ subtitle }}</div>
    </div>
    <div style="flex:1"></div>
    <div class="actions">
      <button v-if="isOpen" class="btn btn-done" @click="$emit('done', todo)">Done</button>
      <button v-if="isOpen" class="btn btn-abandon" @click="$emit('abandon', todo)">Abandon</button>
      <button v-if="!isOpen" class="btn btn-reopen" @click="$emit('reopen', todo)">Reopen</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue'
interface Todo { id: string; title: string; status?: 'open'|'done'|'abandoned'; completed_at?: string|null; archived_at?: string|null }
const props = defineProps<{ todo: Todo }>()
const isOpen = computed(()=> !props.todo.status || props.todo.status === 'open')
const stateClass = computed(()=> {
  if(props.todo.status === 'done') return 'done'
  if(props.todo.status === 'abandoned') return 'abandoned'
  return 'open'
})
const subtitle = computed(()=>{
  if(props.todo.status === 'done') return 'completed'
  if(props.todo.status === 'abandoned') return 'abandoned'
  return ''
})
</script>

<style scoped>
.card{display:flex;align-items:center;gap:12px;background:#fff;border-radius:12px;padding:12px;border:1px solid #e6edf3}
.card.open{background:#fff;border-color:#e6edf3;color:#0f172a}
.card.done{background:linear-gradient(90deg,#ecfdf5,#ffffff);border-color:#bbf7d0;color:#064e3b}
.card.abandoned{background:#f8fafc;border-color:#e6e7eb;color:#52525b;opacity:0.95}
.title{font-weight:600}
.small{font-size:12px;color:#64748b}
.actions{display:flex;gap:8px}
.btn{padding:8px 10px;border-radius:10px;border:0}
.btn-done{background:#10b981;color:#fff}
.btn-abandon{background:#fee2e2;color:#991b1b}
.btn-reopen{background:#2563eb;color:#fff}
</style>
