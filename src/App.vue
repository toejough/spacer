<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm">
      <div class="max-w-3xl mx-auto px-4 py-3">
        <router-link to="/" class="text-2xl font-bold text-gray-900">Spacer</router-link>
      </div>
    </header>
    <main class="max-w-3xl mx-auto">
      <router-view />
    </main>
  </div>
</template>
