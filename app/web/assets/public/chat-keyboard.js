// an on-screen keyboard covers the bottom of the page without resizing it (iOS, and Android by default),
// <html data-keyboard> and --visible-height let the chat fit above it so the text being typed stays in view
(function () {
  var vv = window.visualViewport;
  if (!vv) return;
  // smaller differences are browser toolbars, pinch zoom changes vv.scale instead
  var KEYBOARD_PX = 120;

  function update() {
    var root = document.documentElement;
    var open = vv.scale === 1 && root.clientHeight - vv.height > KEYBOARD_PX;
    root.style.setProperty("--visible-height", vv.height + "px");
    if (open === root.hasAttribute("data-keyboard")) return;
    root.toggleAttribute("data-keyboard", open);
    // the browser pans the page to the focused field, the chat is laid out from the top
    if (open && document.getElementById("chat")) window.scrollTo(0, 0);
  }
  vv.addEventListener("resize", update);
  update();
})();

// with a mouse (a computer) Enter sends the chat message and Shift+Enter is a new line,
// on touch screens Enter stays a new line, the keyboard has no Shift+Enter
document.addEventListener("keydown", function (e) {
  if (e.key !== "Enter" || e.shiftKey || e.isComposing || !matchMedia("(pointer: fine)").matches) return;
  var body = e.target.closest("#chat textarea[name=body]");
  if (!body) return;
  e.preventDefault();
  if (!body.form.querySelector("[data-chat-send]").disabled) body.form.requestSubmit();
});

// Send is disabled while the message box is blank and no file is picked
function chatSendState() {
  var form = document.getElementById("chat-composer");
  if (!form) return;
  var body = form.querySelector("textarea[name=body]");
  var files = form.querySelector("[data-chat-attach]");
  form.querySelector("[data-chat-send]").disabled = !body.value.trim() && !files.files.length;
}
document.addEventListener("input", function (e) {
  if (e.target.matches("#chat textarea[name=body]")) chatSendState();
});
document.addEventListener("htmx:afterSwap", chatSendState);
document.addEventListener("DOMContentLoaded", chatSendState);
chatSendState();

// with a mouse the message box keeps the cursor after sending, on touch screens the keyboard stays closed
document.addEventListener("htmx:afterSwap", function (e) {
  if (e.detail.target.id !== "chat" || !matchMedia("(pointer: fine)").matches) return;
  var body = document.querySelector("#chat textarea[name=body]");
  if (body) body.focus();
});

// answering a message goes straight to typing the answer
document.addEventListener("htmx:afterSwap", function (e) {
  if (e.detail.target.id !== "chat-reply" || !document.querySelector("#chat-reply [name=reply_to]")) return;
  var body = document.querySelector("#chat textarea[name=body]");
  if (body) body.focus();
});

// on touch screens the paperclip sends the chosen files right away as a message of their own, the text
// typed so far is left out (a disabled field is not sent) and put back into the new message box after the swap;
// with a mouse they wait for Send, shown above the message box, and picking again adds to them
var chatDraft = null;
var chatFiles = [];
var chatURLs = [];

function chatPreviews(input) {
  var dt = new DataTransfer();
  chatFiles.forEach(function (f) {
    dt.items.add(f);
  });
  input.files = dt.files;
  chatSendState();
  chatURLs.forEach(URL.revokeObjectURL);
  chatURLs = [];
  var box = input.form.querySelector("[data-chat-previews]");
  box.replaceChildren();
  box.hidden = !chatFiles.length;
  chatFiles.forEach(function (f, i) {
    var item = document.createElement("div");
    item.className = "relative flex h-16 items-center rounded-lg bg-muted";
    if (f.type.startsWith("image/")) {
      var img = document.createElement("img");
      var url = URL.createObjectURL(f);
      chatURLs.push(url);
      img.src = url;
      img.alt = f.name;
      img.className = "size-16 rounded-lg object-cover";
      item.append(img);
    } else {
      var name = document.createElement("span");
      name.textContent = f.name;
      name.className = "max-w-40 truncate px-3 pr-8 text-sm";
      item.append(name);
    }
    var remove = document.createElement("button");
    remove.type = "button";
    remove.textContent = "✕";
    remove.setAttribute("aria-label", box.dataset.remove + " " + f.name);
    remove.className = "absolute -top-2 -right-2 flex size-6 items-center justify-center rounded-full bg-foreground text-xs text-background";
    remove.addEventListener("click", function () {
      chatFiles.splice(i, 1);
      chatPreviews(input);
    });
    item.append(remove);
    box.append(item);
  });
}

document.addEventListener("change", function (e) {
  if (!e.target.matches("#chat [data-chat-attach]")) return;
  if (matchMedia("(pointer: fine)").matches) {
    chatFiles = chatFiles.concat(Array.from(e.target.files));
    chatPreviews(e.target);
    return;
  }
  if (!e.target.files.length) return;
  var body = e.target.form.querySelector("textarea[name=body]");
  chatDraft = body.value;
  body.disabled = true;
  e.target.form.requestSubmit();
});
// the new form has no files
document.addEventListener("htmx:afterSwap", function (e) {
  if (e.detail.target.id !== "chat") return;
  chatFiles = [];
  chatURLs.forEach(URL.revokeObjectURL);
  chatURLs = [];
});
// after the swap, or a failed request that left the old box disabled
document.addEventListener("htmx:afterRequest", function (e) {
  if (chatDraft === null || !e.detail.elt.closest || !e.detail.elt.closest("#chat")) return;
  var body = document.querySelector("#chat textarea[name=body]");
  if (body) {
    body.disabled = false;
    if (!body.value) body.value = chatDraft;
  }
  chatDraft = null;
});
