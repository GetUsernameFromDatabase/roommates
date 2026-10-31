// hx-confirm opens the app's dialog, see confirmDialog: the first line is its title, the
// second its stamp, the rest its text, and data-confirm-label on the form names its button
// when that is not the default Delete; data-confirm-keep names a second button that sends
// keep=1 with the form; without the dialog htmx asks with the browser's own box
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
  var keep = document.querySelector("button[data-confirm-keep]");
  keep.textContent = e.detail.elt.getAttribute("data-confirm-keep");
  keep.parentElement.hidden = !keep.textContent;
  confirmDialogRequest = e.detail.issueRequest;
  confirmDialogForm = e.detail.elt;
  window.templ.dialog.open(title.closest("[role=alertdialog]"));
});
document.addEventListener("click", function (e) {
  var keep = e.target.closest("button[data-confirm-keep]");
  if (!keep && !e.target.closest("[data-confirm-ok]")) return;
  var field = document.createElement("input");
  field.type = "hidden";
  field.name = "keep";
  field.value = "1";
  if (keep && confirmDialogForm) confirmDialogForm.append(field);
  var issueRequest = confirmDialogRequest;
  confirmDialogRequest = null;
  confirmDialogForm = null;
  if (issueRequest) issueRequest(true);
  // htmx has read the form by now, after a failed request it stays and a plain Remove must not keep
  field.remove();
});
document.addEventListener("dialog-close", function () {
  confirmDialogRequest = null;
  confirmDialogForm = null;
});
// htmx drops a request whose form is gone, so a live update waits until the dialog is answered
document.addEventListener("htmx:beforeSwap", function (e) {
  if (confirmDialogForm && e.detail.target.contains(confirmDialogForm)) e.detail.shouldSwap = false;
});
