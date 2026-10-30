// the calendar keeps a Date, its data-date-input names the hidden input that sends it
// as yyyy-MM-dd and "<id>-text" is what the trigger shows
(function () {
  if (window.datePickerLoaded) return;
  window.datePickerLoaded = true;

  function pad(n) {
    return String(n).padStart(2, "0");
  }

  document.addEventListener("calendar-select", function (e) {
    var id = e.target.getAttribute("data-date-input");
    if (!id) return;
    var day = e.detail.selected;
    // local parts: toISOString is the day before east of UTC
    var value = day.getFullYear() + "-" + pad(day.getMonth() + 1) + "-" + pad(day.getDate());
    document.getElementById(id).value = value;
    document.getElementById(id + "-text").textContent = value;
    window.templ.popover.closeNearest(e.target);
  });
})();
