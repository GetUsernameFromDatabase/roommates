// [data-appearance] form: picking a theme or mode shows it on the page right away,
// the form saves it itself with hx-post, so the page and the browser bar are not rendered again
// htmx runs this file again when a swap renders its component again
if (!window.appearancePreviewLoaded) {
  window.appearancePreviewLoaded = true;
  htmx.onLoad((root) => {
    root.querySelectorAll("[data-appearance]").forEach((form) => {
      form.addEventListener("change", (e) => {
        if (e.target.name === "theme" || e.target.name === "mode") {
          document.documentElement.dataset[e.target.name] = e.target.value;
          // shortcut: in system mode the bar misses an OS dark switch until the next full load, fine for one page view
          const bg = getComputedStyle(document.body).backgroundColor;
          document.querySelectorAll('meta[name="theme-color"]').forEach((m) => (m.content = bg));
        }
      });
    });
  });
}
