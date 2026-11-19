// the calendar keeps a Date, its data-date-input names the hidden input that sends it
// as yyyy-MM-dd and "<id>-text" is what the trigger shows
(function () {
  if (window.datePickerLoaded) return;
  window.datePickerLoaded = true;

  function pad(n) {
    return String(n).padStart(2, "0");
  }

  // local parts: toISOString is the day before east of UTC
  function isoDay(day) {
    return day.getFullYear() + "-" + pad(day.getMonth() + 1) + "-" + pad(day.getDate());
  }

  // the short date of the page language, as utils.FormatDate on the server
  function shownDay(value) {
    var p = value.split("-");
    var day = new Date(+p[0], +p[1] - 1, +p[2]);
    if (document.documentElement.lang === "en") {
      return day.toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric" });
    }
    return pad(day.getDate()) + "." + pad(day.getMonth() + 1) + "." + day.getFullYear();
  }

  function setDay(id, value) {
    document.getElementById(id).value = value;
    document.getElementById(id + "-text").textContent = shownDay(value);
  }

  // a new expense comes without a day, today is the browser's, not the server's
  // the whole document: the popover holding the calendar may sit outside the swapped part
  function fill() {
    document.querySelectorAll("[data-date-input]").forEach(function (cal) {
      var id = cal.getAttribute("data-date-input");
      var input = document.getElementById(id);
      if (input) setDay(id, input.value || isoDay(new Date()));
    });
  }

  document.addEventListener("htmx:load", fill);
  document.addEventListener("DOMContentLoaded", fill);

  document.addEventListener("calendar-select", function (e) {
    var id = e.target.getAttribute("data-date-input");
    if (!id) return;
    setDay(id, isoDay(e.detail.selected));
    window.templ.popover.closeNearest(e.target);
  });
})();
