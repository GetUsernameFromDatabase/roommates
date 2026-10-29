// [data-appearance] form: picking a theme or mode shows it on the page right away,
// only Save keeps it, leaving the page without saving goes back to the saved one
// htmx runs this file again when a swap renders its component again
if (!window.appearancePreviewLoaded) {
  window.appearancePreviewLoaded = true;
  htmx.onLoad((root) => {
    root.querySelectorAll("[data-appearance]").forEach((form) => {
      form.addEventListener("change", (e) => {
        if (e.target.name === "theme" || e.target.name === "mode") {
          document.documentElement.dataset[e.target.name] = e.target.value;
        }
      });
    });
  });
}
