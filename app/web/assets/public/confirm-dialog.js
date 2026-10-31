// hx-confirm opens the app's dialog, see confirmDialog: the first line is its title, the
// second its stamp, the rest its text, and data-confirm-label on the form names its button
// when that is not the default Delete; without the dialog htmx asks with the browser's own box
var confirmDialogRequest = null;
var confirmDialogForm = null;
document.addEventListener("htmx:confirm", function (e) {
  var title = document.querySelector("[data-confirm-title]");
  if (!e.detail.question || !title || !window.templ || !window.templ.dialog) return;
  e.preventDefault();
  var lines = e.detail.question.split("\n");
  title.textContent = lines[0];
  document.querySelector("[data-confirm-stamp]").textContent = lines[1] || "";
  document.querySelector("[data-confirm-text]").textContent = lines.slice(2).join("\n");
  var ok = document.querySelector("[data-confirm-ok]");
  ok.textContent = e.detail.elt.getAttribute("data-confirm-label") || ok.getAttribute("data-confirm-ok");
  confirmDialogRequest = e.detail.issueRequest;
  confirmDialogForm = e.detail.elt;
  window.templ.dialog.open(title.closest("[role=alertdialog]"));
});
document.addEventListener("click", function (e) {
  if (!e.target.closest("[data-confirm-ok]")) return;
  var issueRequest = confirmDialogRequest;
  confirmDialogRequest = null;
  confirmDialogForm = null;
  if (issueRequest) issueRequest(true);
});
document.addEventListener("dialog-close", function () {
  confirmDialogRequest = null;
  confirmDialogForm = null;
});
// htmx drops a request whose form is gone, so a live update waits until the dialog is answered
document.addEventListener("htmx:beforeSwap", function (e) {
  if (confirmDialogForm && e.detail.target.contains(confirmDialogForm)) e.detail.shouldSwap = false;
});
