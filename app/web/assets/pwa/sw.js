// service worker, served at /sw.js so its scope is the whole app
//
// shows push notifications and the offline page, nothing else is cached

// shortcut: the offline page keeps the theme and language of the visit that installed
// the worker, bump the version to refresh it for everyone
const CACHE = "offline-v1";
const OFFLINE = "/offline";
const CSS = "/assets/output.css";

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE).then((c) => c.addAll([OFFLINE, CSS])),
  );
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))),
      )
      .then(() => self.clients.claim()),
  );
});

self.addEventListener("fetch", (event) => {
  const req = event.request;
  const navigate = req.mode === "navigate";
  // pages and the offline page's styles, live updates and downloads go straight to the network
  if (!navigate && new URL(req.url).pathname !== CSS) {
    return;
  }
  event.respondWith(
    fetch(req).catch(() => caches.match(navigate ? OFFLINE : req)),
  );
});

// a message is {title, body, url, tag}, see internal/push
self.addEventListener("push", (event) => {
  if (!event.data) {
    return;
  }
  const msg = event.data.json();
  event.waitUntil(
    self.clients
      .matchAll({ type: "window", includeUncontrolled: true })
      .then((windows) => {
        // the page is open and live already, browsers allow skipping while it has focus
        const seen = windows.some(
          (w) => w.focused && new URL(w.url).pathname === msg.url,
        );
        if (seen) {
          return;
        }
        return self.registration.showNotification(msg.title, {
          body: msg.body,
          tag: msg.tag,
          icon: "/icons/android-chrome-192x192.png",
          data: { url: msg.url },
        });
      }),
  );
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const url = event.notification.data.url;
  event.waitUntil(
    self.clients
      .matchAll({ type: "window", includeUncontrolled: true })
      .then((windows) => {
        const open = windows.find((w) => new URL(w.url).pathname === url);
        return open ? open.focus() : self.clients.openWindow(url);
      }),
  );
});
