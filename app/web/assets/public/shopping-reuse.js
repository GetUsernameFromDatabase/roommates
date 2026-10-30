// a Buy again button fills the pad's blank line and asks for the quantity, the old one shows
// as the placeholder, and after adding from the quantity field the next item can be typed right away
// htmx runs this file again when a swap renders its component again
if (!window.shoppingReuseLoaded) {
  window.shoppingReuseLoaded = true;

  document.addEventListener("click", function (e) {
    var btn = e.target.closest("[data-reuse-text]");
    if (!btn) return;
    var text = document.getElementById("item-text");
    var qty = document.getElementById("item-qty");
    if (!text || !qty) return;
    text.value = btn.dataset.reuseText;
    qty.value = "";
    if (btn.dataset.reuseQty) qty.placeholder = btn.dataset.reuseQty;
    qty.focus();
  });

  // only the add, a flatmate's change arriving while typing must not move the focus
  document.addEventListener("htmx:afterSettle", function (e) {
    var request = e.detail.requestConfig;
    if (!request || request.verb !== "post" || !request.path.endsWith("/items")) return;
    var text = document.getElementById("item-text");
    if (text && document.activeElement && document.activeElement.id === "item-qty") text.focus();
  });
}
