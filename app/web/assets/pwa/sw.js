// service worker, served at routes.ServiceWorker (/sw.js) so its scope is the whole app
//
// shows push notifications and the offline page, nothing else is cached

// self.ROUTES, the paths of internal/routes; the one path a worker must know before it has them
importScripts("/routes.js");

// shortcut: the offline page keeps the theme and language of the visit that installed
// the worker, bump the version to refresh it for everyone
const CACHE = "offline-v1";
const OFFLINE = self.ROUTES.offline;
const CSS = self.ROUTES.css;

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
    // pages ask for the CSS with ?v=, the cache holds it without
    fetch(req).catch(() => caches.match(navigate ? OFFLINE : req, { ignoreSearch: true })),
  );
});

// whether a window shows url: its query only counts when url has one (a calendar reminder's ?on=),
// a page's own search (chat, expenses) still shows that page
const shows = (w, url) => {
  const a = new URL(w.url);
  const b = new URL(url, self.location.origin);
  return a.pathname === b.pathname && (!b.search || a.search === b.search);
};

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
        const seen = windows.some((w) => w.focused && shows(w, msg.url));
        if (seen) {
          return;
        }
        return self.registration.showNotification(msg.title, {
          body: msg.body,
          tag: msg.tag,
          icon: self.ROUTES.icon,
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
        const open = windows.find((w) => shows(w, url));
        return open ? open.focus() : self.clients.openWindow(url);
      }),
  );
});
