// a [data-dirty-save] form shows its [data-save] button only while its input differs from what
// was loaded, Escape goes back to that; htmx runs this file again when a swap renders its component again
if (!window.dirtySaveLoaded) {
  window.dirtySaveLoaded = true;
  (function () {
    function update(input) {
      input.form.querySelector("[data-save]").hidden = input.value.trim() === input.defaultValue;
    }
    document.addEventListener("input", function (e) {
      if (e.target.matches("[data-dirty-save] input")) update(e.target);
    });
    document.addEventListener("keydown", function (e) {
      if (e.key !== "Escape" || !e.target.matches("[data-dirty-save] input")) return;
      e.target.value = e.target.defaultValue;
      update(e.target);
      e.target.blur();
    });
  })();
}
