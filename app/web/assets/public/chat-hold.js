// chat message interactions and scroll:
// a long press on a chat message shows its controls (data-open), hover does it on screens with a mouse
(function () {
  var HOLD_MS = 500;
  var MOVE_PX = 10;
  var timer = null;
  var start = null;
  // the contextmenu of a press that opened the controls, Android shows its own menu otherwise
  var held = false;

  function cancel() {
    clearTimeout(timer);
    timer = null;
  }

  function close(except) {
    document.querySelectorAll("[data-slot=message][data-open]").forEach(function (m) {
      if (m !== except) m.removeAttribute("data-open");
    });
  }

  document.addEventListener("pointerdown", function (e) {
    var m = e.target.closest("[data-slot=message]");
    if (!m || !m.hasAttribute("data-open")) close(null);
    held = false;
    if (!m || e.pointerType === "mouse") return;
    start = { x: e.clientX, y: e.clientY };
    cancel();
    timer = setTimeout(function () {
      timer = null;
      held = true;
      close(m);
      m.setAttribute("data-open", "");
    }, HOLD_MS);
  });
  document.addEventListener("pointermove", function (e) {
    if (timer && Math.hypot(e.clientX - start.x, e.clientY - start.y) > MOVE_PX) cancel();
  });
  ["pointerup", "pointercancel"].forEach(function (type) {
    document.addEventListener(type, cancel);
  });
  document.addEventListener("scroll", cancel, true);
  // on touch screens the long press that opens the browser's menu (images, links) opens the controls instead,
  // Android fires it before or with the timer
  document.addEventListener("contextmenu", function (e) {
    var m = e.target.closest("[data-slot=message]");
    if (!m || (!held && !matchMedia("(pointer: coarse)").matches)) return;
    e.preventDefault();
    cancel();
    held = true;
    close(m);
    m.setAttribute("data-open", "");
  });

  document.addEventListener("click", function (e) {
    var copy = e.target.closest("[data-chat-copy]");
    if (!copy) return;
    var m = copy.closest("[data-slot=message]");
    var body = m && m.querySelector("[data-chat-body]");
    if (body && navigator.clipboard) navigator.clipboard.writeText(body.textContent.trim());
    close(null);
  });
})();

// every swap renders a new chat log and browsers disagree where it starts, so the scroll is set here:
// only sending (the composer form, a direct child of #chat) goes to the newest message, everything else
// (a live refresh, reply, cancelling it, deleting) keeps the reader's place.
// in the column-reverse log scrollTop counts from the bottom, 0 is the newest message
(function () {
  var keep = null;
  document.addEventListener("htmx:beforeSwap", function (e) {
    var log = document.querySelector("#chat-messages [role=log]");
    // a response that is not swapped (204) leaves nothing to restore
    if (!log || !e.detail.shouldSwap) {
      keep = null;
      return;
    }
    var cfg = e.detail.requestConfig;
    var sending = cfg && cfg.verb === "post" && cfg.elt.matches && cfg.elt.matches("#chat > form");
    keep = sending ? 0 : log.scrollTop;
  });
  document.addEventListener("htmx:afterSwap", function () {
    if (keep === null) return;
    var log = document.querySelector("#chat-messages [role=log]");
    if (log) log.scrollTop = keep;
    keep = null;
  });
})();
