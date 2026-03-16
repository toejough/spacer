<template>
  <div class="min-h-screen bg-bg text-text">
    <header class="bg-primary text-white p-4">
      <router-link to="/" class="text-xl font-bold">Spacer</router-link>
    </header>
    <main class="max-w-lg mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>
