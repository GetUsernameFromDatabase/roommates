// a [data-jump-nav] shows on a page about two screens long, scrolls its links' sections in just
// under the bars above them and marks the link of the section being read; htmx runs this file again when a swap renders the nav again
if (!window.jumpNavLoaded) {
  window.jumpNavLoaded = true;
  (function () {
    // where a section's top lands once the nav is stuck: under the nav and a little space
    function line(nav) {
      return parseFloat(getComputedStyle(nav).top) + nav.offsetHeight + 12;
    }

    function section(a) {
      return document.getElementById(a.hash.slice(1));
    }

    function update() {
      var nav = document.querySelector("[data-jump-nav]");
      if (!nav) return;
      var root = document.documentElement;
      var links = Array.prototype.slice.call(nav.querySelectorAll("a"));
      var visible = links.filter(function (a) {
        var s = section(a);
        // a live refresh can empty a section, its link goes with it
        a.parentElement.hidden = !s || s.offsetHeight === 0;
        return !a.parentElement.hidden;
      });
      // a page that fits in about two screens is quick to scroll by hand; the shown nav adds exactly
      // its height (its -my-4 cancels the gap it adds), so taking it off keeps it from flapping;
      // clientHeight, not innerHeight: that changes as a phone's address bar slides away mid-scroll
      var height = root.scrollHeight - (nav.hidden ? 0 : nav.offsetHeight);
      nav.hidden = height < 2 * root.clientHeight || visible.length < 2;
      if (nav.hidden) return;
      var y = line(nav) + 1;
      var current = null;
      var currentTop = -Infinity;
      visible.forEach(function (a) {
        var top = section(a).getBoundingClientRect().top;
        // sections side by side start level, the first of them stays current
        if (!current || (top <= y && top > currentTop)) {
          current = a;
          currentTop = top;
        }
      });
      // the last sections may be too short to ever reach the line, at the bottom the last one is
      // current, or the one tapped while it is in view
      if (innerHeight + scrollY >= root.scrollHeight - 2) {
        var tapped = visible.indexOf(pinned) >= 0 ? section(pinned).getBoundingClientRect().top : -1;
        current = tapped >= 0 && tapped < innerHeight ? pinned : visible[visible.length - 1];
      }
      links.forEach(function (a) {
        if (a === current) a.setAttribute("aria-current", "true");
        else a.removeAttribute("aria-current");
      });
      if (!current) return;
      // the row's own scroll, scrollIntoView could move the page as well
      var row = current.closest("ul");
      var li = current.parentElement;
      if (li.offsetLeft < row.scrollLeft || li.offsetLeft + li.offsetWidth > row.scrollLeft + row.clientWidth) {
        row.scrollLeft = li.offsetLeft - (row.clientWidth - li.offsetWidth) / 2;
      }
    }

    // the link last tapped, until the user scrolls by hand
    var pinned = null;
    var queued = false;
    function queue() {
      if (queued) return;
      queued = true;
      requestAnimationFrame(function () {
        queued = false;
        update();
      });
    }

    document.addEventListener("click", function (e) {
      var a = e.target.closest("[data-jump-nav] a");
      if (!a) return;
      var s = section(a);
      if (!s) return;
      e.preventDefault();
      pinned = a;
      var nav = a.closest("[data-jump-nav]");
      var reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;
      scrollTo({ top: s.getBoundingClientRect().top + scrollY - line(nav), behavior: reduce ? "auto" : "smooth" });
      // keyboard and screen reader users carry on from the section
      if (!s.hasAttribute("tabindex")) s.tabIndex = -1;
      s.focus({ preventScroll: true });
      // at the bottom the page cannot scroll, no scroll event would mark the tapped link
      queue();
    });
    addEventListener("scroll", queue, { passive: true });
    ["wheel", "touchstart"].forEach(function (type) {
      addEventListener(type, function () { pinned = null; }, { passive: true });
    });
    addEventListener("resize", queue);
    document.addEventListener("htmx:afterSettle", queue);
    // the script comes before the sections, until they are parsed their links would hide
    if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", queue);
    else queue();
  })();
}
