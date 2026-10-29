// a [data-copy] button copies its value and says so in its label for a moment,
// without the clipboard (plain http) the buttons are hidden by .no-clipboard
// htmx runs this file again when a swap renders its component again
if (!window.copyButtonLoaded) {
  window.copyButtonLoaded = true;
  if (!navigator.clipboard) document.documentElement.classList.add("no-clipboard");

  document.addEventListener("click", function (e) {
    var btn = e.target.closest("[data-copy]");
    if (!btn || !navigator.clipboard) return;
    // the label of the first click, a second click within the moment would keep "copied"
    if (!btn.dataset.label) btn.dataset.label = btn.textContent;
    navigator.clipboard.writeText(btn.dataset.copy).then(function () {
      btn.textContent = btn.dataset.copied;
      setTimeout(function () {
        btn.textContent = btn.dataset.label;
      }, 2000);
    });
  });
}
