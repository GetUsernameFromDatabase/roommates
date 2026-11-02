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
  // a region whose page is gone (deleted, or the user was removed) gets the 404 page, which has no
  // copy of it and would swap in nothing, so its message fills the page's <main> instead
  document.addEventListener("htmx:beforeSwap", function (e) {
    var main = document.querySelector("main");
    if (e.detail.xhr.status !== 404 || !main || !e.detail.requestConfig.elt.hasAttribute("data-live")) return;
    e.detail.target = main;
    e.detail.selectOverride = "main > *";
    e.detail.swapOverride = "innerHTML";
  });
})();
