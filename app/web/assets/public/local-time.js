// shows <time data-local> in the viewer's time zone and language, the server only knows UTC,
// without today's date and this year
// htmx runs this file again when a swap renders its component again
if (!window.localTimeLoaded) {
  window.localTimeLoaded = true;
  htmx.onLoad(function (root) {
    var lang = document.documentElement.lang || undefined;
    var now = new Date();
    root.querySelectorAll("time[data-local]").forEach(function (t) {
      var d = new Date(t.dateTime);
      var options = { timeStyle: "short" };
      if (d.toDateString() !== now.toDateString()) {
        options = { day: "numeric", month: "short", hour: "2-digit", minute: "2-digit" };
        if (d.getFullYear() !== now.getFullYear()) options.year = "numeric";
      }
      t.textContent = d.toLocaleString(lang, options);
    });
  });
}
