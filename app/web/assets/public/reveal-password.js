// a [data-reveal] toggle button shows the password of the input it names as text,
// and hides it again, its icon follows aria-pressed
document.addEventListener("click", function (e) {
  var btn = e.target.closest("[data-reveal]");
  if (!btn) return;
  var input = document.getElementById(btn.dataset.reveal);
  var show = input.type === "password";
  input.type = show ? "text" : "password";
  btn.setAttribute("aria-pressed", show);
});
