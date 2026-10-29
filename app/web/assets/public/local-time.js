// shows <time data-local> in the viewer's time zone and language, the server only knows UTC
htmx.onLoad(function (root) {
  root.querySelectorAll("time[data-local]").forEach(function (t) {
    t.textContent = new Date(t.dateTime).toLocaleString(document.documentElement.lang || undefined, {
      dateStyle: "short",
      timeStyle: "short",
    });
  });
});
