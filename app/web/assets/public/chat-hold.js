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
  document.addEventListener("contextmenu", function (e) {
    if (held && e.target.closest("[data-slot=message]")) e.preventDefault();
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
// a live refresh keeps the reader's place, sending goes to the newest message.
// in the column-reverse log scrollTop counts from the bottom, 0 is the newest message
(function () {
  var keep = null;
  document.addEventListener("htmx:beforeSwap", function (e) {
    var id = e.detail.target.id;
    // a response that is not swapped (204) leaves nothing to restore
    keep = id === "chat" && e.detail.shouldSwap ? 0 : null;
    if (id !== "chat-messages" || !e.detail.shouldSwap) return;
    var log = e.detail.target.querySelector("[role=log]");
    keep = log ? log.scrollTop : 0;
  });
  document.addEventListener("htmx:afterSwap", function () {
    if (keep === null) return;
    var log = document.querySelector("#chat-messages [role=log]");
    if (log) log.scrollTop = keep;
    keep = null;
  });
})();
