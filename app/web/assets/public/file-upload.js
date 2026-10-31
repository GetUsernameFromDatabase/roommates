// a picked file is uploaded right away, its button is off until the page reloads
document.addEventListener("change", function (e) {
  if (!e.target.matches("input[data-upload-on-pick]") || !e.target.files.length) return;
  e.target.form.requestSubmit();
  e.target.closest("label").inert = true;
});

// the back button can restore the page as it was left, mid upload
addEventListener("pageshow", function () {
  document.querySelectorAll("label:has(input[data-upload-on-pick])").forEach(function (label) {
    label.inert = false;
  });
});
