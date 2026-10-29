// an open <details data-menu> closes on a click outside it and on Escape
document.addEventListener("click", function (e) {
  document.querySelectorAll("details[data-menu][open]").forEach(function (d) {
    if (!d.contains(e.target)) d.open = false;
  });
});
document.addEventListener("keydown", function (e) {
  if (e.key !== "Escape") return;
  document.querySelectorAll("details[data-menu][open]").forEach(function (d) {
    if (d.contains(document.activeElement)) d.querySelector("summary").focus();
    d.open = false;
  });
});
