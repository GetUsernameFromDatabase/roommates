// registers the service worker (web/assets/pwa/sw.js) on every page, and drives
// the notification switch [data-push] and the logout form [data-push-logout]
(() => {
  if (!("serviceWorker" in navigator)) {
    return;
  }
  const ready = navigator.serviceWorker.register("/sw.js").then(() => navigator.serviceWorker.ready);
  const supported = "PushManager" in window && "Notification" in window;

  const post = (path, body) =>
    fetch(path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

  // applicationServerKey takes bytes, the server gives base64url
  const keyBytes = (key) => {
    const b64 = key.replace(/-/g, "+").replace(/_/g, "/");
    return Uint8Array.from(atob(b64 + "=".repeat((4 - (b64.length % 4)) % 4)), (c) => c.charCodeAt(0));
  };

  const unsubscribe = async () => {
    const reg = await ready;
    const sub = await reg.pushManager.getSubscription();
    if (sub) {
      await post("/push/unsubscribe", { endpoint: sub.endpoint });
      await sub.unsubscribe();
    }
  };

  const setup = async (box) => {
    const key = box.dataset.vapidKey;
    const button = box.querySelector("button");
    const hint = box.querySelector("[data-push-hint]");
    const show = (text) => {
      hint.textContent = text;
      hint.hidden = !text;
    };
    if (!key) {
      return;
    }
    box.hidden = false;
    if (!supported) {
      // iOS only has push in an app added to the home screen
      button.hidden = true;
      show(box.dataset.install);
      return;
    }
    const reg = await ready;
    const render = async () => {
      const on = (await reg.pushManager.getSubscription()) !== null;
      button.textContent = on ? box.dataset.disable : box.dataset.enable;
      button.setAttribute("aria-pressed", String(on));
      show(Notification.permission === "denied" ? box.dataset.denied : "");
    };
    button.addEventListener("click", async () => {
      button.disabled = true;
      try {
        if (await reg.pushManager.getSubscription()) {
          await unsubscribe();
        } else if ((await Notification.requestPermission()) === "granted") {
          const sub = await reg.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: keyBytes(key),
          });
          const res = await post("/push/subscribe", sub.toJSON());
          if (!res.ok) {
            await sub.unsubscribe();
          }
        }
      } finally {
        button.disabled = false;
        await render();
      }
    });
    await render();
  };

  htmx.onLoad((root) => {
    root.querySelectorAll("[data-push]").forEach(setup);
    // the next person signing in on this device must not get these notifications
    root.querySelectorAll("form[data-push-logout]").forEach((form) => {
      form.addEventListener("submit", (event) => {
        if (!supported || form.dataset.done) {
          return;
        }
        event.preventDefault();
        form.dataset.done = "1";
        // a stuck service worker must not keep anyone signed in
        const timeout = new Promise((resolve) => setTimeout(resolve, 3000));
        Promise.race([unsubscribe(), timeout]).finally(() => form.submit());
      });
    });
  });
})();
