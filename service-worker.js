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

self.addEventListener('fetch', event => {
  event.respondWith(
    fetch(event.request)
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
});
