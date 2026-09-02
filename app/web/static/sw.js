// Kill switch. A tablet can have the old PWA service worker installed. This
// file replaces it at the same path, removes the caches, and unregisters.
self.addEventListener("install", () => {
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      const cacheNames = await caches.keys();
      await Promise.all(cacheNames.map((name) => caches.delete(name)));
      await self.registration.unregister();

      const windowClients = await self.clients.matchAll({ type: "window" });
      for (const windowClient of windowClients) {
        windowClient.navigate(windowClient.url);
      }
    })()
  );
});
