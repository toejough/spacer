// Service worker version
const SW_VERSION = '0.1.0.2025-07-26.14';

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open('flashcards-v1').then(cache => {
      return cache.addAll([
        '/',
        '/index.html',
        '/styles.css',
        '/script.js',
        '/manifest.webmanifest',
        '/icon-192.png',
        '/icon-512.png'
      ]);
    })
  );
  self.skipWaiting();
});

self.addEventListener('activate', event => {
  event.waitUntil(self.clients.claim());
});

// On iOS, service workers are less aggressive about updates. Force update check on every launch.
self.addEventListener('fetch', event => {
  // Only handle GET requests
  if (event.request.method !== 'GET') return;

  // Try network first, fallback to cache
  event.respondWith(
    fetch(event.request, { cache: 'reload' })
      .then(response => {
        // Clone and store in cache
        const responseClone = response.clone();
        caches.open('flashcards-v1').then(cache => {
          cache.put(event.request, responseClone);
        });
        return response;
      })
      .catch(() => caches.match(event.request))
  );

  // For iOS: update cache in the background for next launch
  event.waitUntil(
    caches.open('flashcards-v1').then(cache => {
      return fetch(event.request, { cache: 'reload' })
        .then(response => {
          cache.put(event.request, response.clone());
        })
        .catch(() => {});
    })
  );
});
