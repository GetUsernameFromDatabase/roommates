// shows <time data-local> in the viewer's time zone and language, the server only knows UTC,
// without today's date and this year; data-today and data-yesterday ("today at %s") put the
// time into a phrase on those days instead
// htmx runs this file again when a swap renders its component again
if (!window.localTimeLoaded) {
  window.localTimeLoaded = true;
  htmx.onLoad(function (root) {
    var lang = document.documentElement.lang || undefined;
    var now = new Date();
    var yesterday = new Date(now);
    yesterday.setDate(now.getDate() - 1);
    root.querySelectorAll("time[data-local]").forEach(function (t) {
      var d = new Date(t.dateTime);
      var day = d.toDateString();
      var phrase = day === now.toDateString() ? t.dataset.today : day === yesterday.toDateString() ? t.dataset.yesterday : "";
      if (phrase) {
        t.textContent = phrase.replace("%s", d.toLocaleString(lang, { timeStyle: "short" }));
        return;
      }
      var options = { timeStyle: "short" };
      if (day !== now.toDateString()) {
        options = { day: "numeric", month: "short", hour: "2-digit", minute: "2-digit" };
        if (d.getFullYear() !== now.getFullYear()) options.year = "numeric";
      }
      t.textContent = d.toLocaleString(lang, options);
    });
  });
}
