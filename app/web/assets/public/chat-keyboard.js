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
  body.form.requestSubmit();
});

// with a mouse the message box keeps the cursor after sending, on touch screens the keyboard stays closed
document.addEventListener("htmx:afterSwap", function (e) {
  if (e.detail.target.id !== "chat" || !matchMedia("(pointer: fine)").matches) return;
  var body = document.querySelector("#chat textarea[name=body]");
  if (body) body.focus();
});
