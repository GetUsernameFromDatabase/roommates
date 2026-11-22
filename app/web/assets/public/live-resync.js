// changes made while a page was not listening are not sent again, so the live regions refetch
// (live-resync, see sseTrigger): after the event stream reconnects, when the page comes back
// from the history or back-forward cache, and when a hidden tab or app is shown again
(function () {
  function resync() {
    if (document.body) htmx.trigger(document.body, "live-resync");
  }
  var opened = false;
  document.addEventListener("htmx:sseOpen", function () {
    // the first open is the page load, which is fresh
    if (opened) resync();
    opened = true;
  });
  document.addEventListener("htmx:historyRestore", resync);
  window.addEventListener("pageshow", function (e) {
    if (e.persisted) resync();
  });
  document.addEventListener("visibilitychange", function () {
    if (document.visibilityState === "visible") resync();
  });
  // an action of this tab that moves it to another page (its own delete) is answered before its
  // live event comes back, so the page it leaves refetches what it just deleted and gets a 404
  var leaving = false;
  document.addEventListener("htmx:afterRequest", function (e) {
    var xhr = e.detail.xhr;
    if (xhr && (xhr.getResponseHeader("HX-Location") || xhr.getResponseHeader("HX-Redirect"))) leaving = true;
  });
  // the move is over once a page settles: the new one in <body>, or an error page in <main>; or
  // once it could not be sent
  document.addEventListener("htmx:afterSettle", function (e) {
    if (e.detail.target === document.body || e.detail.target.tagName === "MAIN") leaving = false;
  });
  document.addEventListener("htmx:sendError", function () {
    leaving = false;
  });
  // a region whose page is gone (deleted, or the user was removed) gets the 404 page, which has no
  // copy of it and would swap in nothing, so its message fills the page's <main> instead
  document.addEventListener("htmx:beforeSwap", function (e) {
    var main = document.querySelector("main");
    var elt = e.detail.requestConfig.elt;
    if (e.detail.xhr.status !== 404 || !main || !elt.hasAttribute("data-live")) return;
    // the tab is moving on, or has moved on, from the page this region was on
    if (leaving || !elt.isConnected) {
      e.detail.shouldSwap = false;
      return;
    }
    e.detail.target = main;
    e.detail.selectOverride = "main > *";
    e.detail.swapOverride = "innerHTML";
  });

  // a form that is sending cannot be sent again: htmx queues a second submit (Enter included)
  // and sends it after the first, so stop it before htmx sees it
  document.addEventListener(
    "submit",
    function (e) {
      if (e.target.classList && e.target.classList.contains("htmx-request")) {
        e.preventDefault();
        e.stopImmediatePropagation();
      }
    },
    true
  );

  // a request that never got an answer (offline, timeout) or an error page not from the app says so
  // in a toast, 4xx answers are pages htmx swaps in on purpose (form errors, the 404 above), the
  // toast is `templ toast` in layout.templ without its icon, its classes come on <body>
  function requestFailed() {
    var text = document.body && document.body.dataset.requestFailed;
    if (!text) return;
    var old = document.getElementById("request-failed-toast");
    if (old) old.remove();
    var toast = document.createElement("div");
    toast.id = "request-failed-toast";
    toast.setAttribute("role", "alert");
    toast.className = document.body.dataset.toastClass;
    var span = document.createElement("span");
    span.className = "min-w-0 wrap-break-word";
    span.textContent = text;
    toast.appendChild(span);
    document.body.appendChild(toast);
    setTimeout(function () {
      toast.remove();
    }, 5500);
  }
  document.addEventListener("htmx:sendError", requestFailed);
  document.addEventListener("htmx:timeout", requestFailed);
  // htmx swaps every error (responseHandling): the app's own error page retargets <main>
  // (swapMainOnError), anything else (a proxy's 502 page) would replace the whole app
  document.addEventListener("htmx:beforeSwap", function (e) {
    var xhr = e.detail.xhr;
    if (xhr.status < 500 || xhr.getResponseHeader("HX-Retarget")) return;
    e.detail.shouldSwap = false;
    requestFailed();
  });

  // after htmx moved to another page (a boosted form, HX-Location) focus is lost with the old one:
  // put it on <main> so a keyboard or screen reader user starts at the new content, unless htmx
  // gave focus back to an element by id; stay() does not push a URL, so it keeps its place
  // (htmx pushes the URL before it swaps)
  var moved = false;
  document.addEventListener("htmx:pushedIntoHistory", function () {
    moved = true;
  });
  document.addEventListener("htmx:afterSettle", function () {
    if (!moved) return;
    moved = false;
    var main = document.querySelector("main");
    var a = document.activeElement;
    if (main && (!a || a === document.body)) main.focus({ preventScroll: true });
  });
})();
