// the note editor: the formatted note is edited in place and turndown writes it back as markdown
// into the form's textarea on every change, so a note nobody touched keeps its markdown as typed;
// the words on top set what is typed next or the line the cursor is on, a bubble at a selection
// changes the selected text, the Markdown switch edits the textarea itself and switching back
// asks the server to format it
// htmx runs this file again when a swap renders its component again
if (!window.noteEditorLoaded) {
  window.noteEditorLoaded = true;

  var turndown = new TurndownService({
    headingStyle: "atx",
    bulletListMarker: "-",
    emDelimiter: "_",
    codeBlockStyle: "fenced",
    // the server renders with hard wraps, a plain newline is a line break
    br: "",
  });
  turndown.use(turndownPluginGfm.gfm);

  // the last selection in the editor, a tap on a tool or the link field takes focus from it
  var noteSelection = null;

  var escapeHTML = function (s) {
    return s.replace(/[&<>"]/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c];
    });
  };

  // the link around the selection, if there is one
  var selectedLink = function (editor) {
    var node = noteSelection && noteSelection.commonAncestorContainer;
    var a = node && (node.nodeType === 1 ? node : node.parentElement).closest("a");
    return a && editor.contains(a) ? a : null;
  };

  var isOn = function (editor, b) {
    var cmd = b.dataset.cmd;
    if (cmd === "link") return !!selectedLink(editor);
    if (cmd === "formatBlock") return document.queryCommandValue("formatBlock").toLowerCase() === b.dataset.value;
    return document.queryCommandState(cmd);
  };

  // the bubble or the link panel above the selection, below it on touch where the phone's own
  // menu sits above, kept inside the paper
  var placeBubble = function (bubble, rect) {
    var paper = bubble.offsetParent;
    var box = paper.getBoundingClientRect();
    var below = window.matchMedia("(pointer: coarse)").matches;
    var left = rect.left + rect.width / 2 - box.left - bubble.offsetWidth / 2;
    var top = rect.top - box.top - bubble.offsetHeight - 8;
    // on the first line it would leave the paper, where a folded corner clips it
    if (below || top < 0) top = rect.bottom - box.top + 8;
    bubble.style.left = Math.max(4, Math.min(left, paper.clientWidth - bubble.offsetWidth - 4)) + "px";
    bubble.style.top = top + "px";
  };

  // lights the tools whose style is on where the cursor is
  var refreshPressed = function (editor) {
    document.querySelectorAll("#note-tools [aria-pressed], #note-bubble [aria-pressed]").forEach(function (b) {
      b.setAttribute("aria-pressed", String(isOn(editor, b)));
    });
  };

  document.addEventListener("selectionchange", function () {
    var editor = document.getElementById("note-editor");
    var bubble = document.getElementById("note-bubble");
    var linkRow = document.getElementById("note-link");
    var s = document.getSelection();
    if (!editor || editor.hidden) return;
    if (!s.rangeCount || !editor.contains(s.anchorNode)) {
      // the selection left the note, the bubble goes too unless it is being used
      if (!bubble.contains(document.activeElement) && !linkRow.contains(document.activeElement)) {
        bubble.hidden = true;
      }
      return;
    }
    noteSelection = s.getRangeAt(0);
    refreshPressed(editor);
    if (!linkRow.hidden) return;
    bubble.hidden = s.isCollapsed;
    if (!s.isCollapsed) placeBubble(bubble, noteSelection.getBoundingClientRect());
  });

  htmx.onLoad(function () {
    var editor = document.getElementById("note-editor");
    if (!editor || editor.dataset.ready) return;
    editor.dataset.ready = "1";
    var raw = document.getElementById("note-body");
    var hint = document.getElementById("note-body-hint");
    var tools = document.getElementById("note-tools");
    var bubble = document.getElementById("note-bubble");
    var linkRow = document.getElementById("note-link");
    var linkURL = document.getElementById("note-link-url");
    var linkText = document.getElementById("note-link-text");
    var linkError = document.getElementById("note-link-error");
    var linkAdd = linkRow.querySelector("[data-link-add]");
    var linkRemove = linkRow.querySelector("[data-link-remove]");
    // the link being changed, null while adding one
    var editing = null;
    var failed = document.getElementById("note-render-failed");
    var markdown = tools.querySelector('[data-slot="switch"]');
    var commands = document.querySelectorAll("#note-tools [data-cmd], #note-bubble [data-cmd]");
    noteSelection = null;

    // Enter starts a paragraph, Shift+Enter a line break
    document.execCommand("defaultParagraphSeparator", false, "p");

    // the switch is driven from here, it only flips once the formatted view is back
    var setMarkdown = function (on) {
      markdown.setAttribute("aria-checked", String(on));
      [markdown, markdown.firstElementChild].forEach(function (el) {
        el.toggleAttribute("data-checked", on);
        el.toggleAttribute("data-unchecked", !on);
      });
      markdown.nextElementSibling.checked = on;
    };

    var show = function (formatted) {
      editor.hidden = !formatted;
      raw.hidden = formatted;
      hint.hidden = formatted;
      closeBubble();
      tools.querySelectorAll("[data-cmd]").forEach(function (b) {
        b.disabled = !formatted;
      });
      setMarkdown(!formatted);
    };

    var restore = function () {
      editor.focus();
      if (!noteSelection) return;
      var s = document.getSelection();
      s.removeAllRanges();
      s.addRange(noteSelection);
    };

    var clearLinkError = function () {
      linkError.hidden = true;
      linkURL.removeAttribute("aria-invalid");
    };

    var closeBubble = function () {
      clearLinkError();
      linkURL.value = "";
      linkText.value = "";
      linkRow.hidden = true;
      bubble.hidden = true;
    };

    // selects all of the link being changed, so it is replaced or unlinked whole
    var selectEditing = function () {
      restore();
      var range = document.createRange();
      range.selectNode(editing);
      document.getSelection().removeAllRanges();
      document.getSelection().addRange(range);
    };

    // the text and address fields in the bubble, at the selection or the cursor, filled from
    // the link the cursor is in or the selected text
    var openLink = function () {
      var rect = noteSelection ? noteSelection.getBoundingClientRect() : editor.getBoundingClientRect();
      // a cursor in an empty line has no box, the line's element has one
      if (!rect.height && noteSelection) {
        var node = noteSelection.startContainer;
        rect = (node.nodeType === 1 ? node : node.parentElement).getBoundingClientRect();
      }
      editing = selectedLink(editor);
      linkText.value = editing ? editing.textContent : noteSelection ? noteSelection.toString() : "";
      linkURL.value = editing ? editing.getAttribute("href") : "";
      linkAdd.textContent = editing ? linkAdd.dataset.edit : linkAdd.dataset.new;
      linkRemove.hidden = !editing;
      bubble.hidden = true;
      linkRow.hidden = false;
      placeBubble(linkRow, rect);
      (linkText.value ? linkURL : linkText).focus();
    };

    var addLink = function () {
      var url = linkURL.value.trim();
      // "example.com" means the web page, not a path on this site
      // "localhost:8090" has no scheme either, only "scheme://", mailto: and tel: do
      if (url && !/^([a-z][a-z0-9+.-]*:\/\/|mailto:|tel:)/i.test(url)) url = "https://" + url;
      linkURL.value = url;
      // empty, or still not an address
      if (!linkURL.checkValidity()) {
        linkError.hidden = false;
        linkURL.setAttribute("aria-invalid", "true");
        linkURL.focus();
        return;
      }
      var text = linkText.value.trim() || url;
      closeBubble();
      if (editing) selectEditing();
      else restore();
      var s = document.getSelection();
      if (!s.isCollapsed && s.toString() === text) {
        // the same words keep their bold or italic
        document.execCommand("createLink", false, url);
      } else {
        document.execCommand("insertHTML", false, '<a href="' + escapeHTML(url) + '">' + escapeHTML(text) + "</a>");
      }
    };

    var run = function (b) {
      var cmd = b.dataset.cmd;
      if (cmd === "link") {
        openLink();
        return;
      }
      restore();
      if (cmd === "formatBlock") {
        // a second tap turns the heading or quote back into a paragraph
        document.execCommand("formatBlock", false, isOn(editor, b) ? "p" : b.dataset.value);
      } else {
        document.execCommand(cmd);
      }
      // a style switched on for the next text moves no cursor, so no selectionchange says so
      refreshPressed(editor);
    };

    // the formatted note is written back as markdown only when it is needed, a long note
    // would otherwise be converted on every key
    var dirty = false;
    var writeBack = function () {
      if (!dirty) return;
      raw.value = turndown.turndown(editor.innerHTML);
      dirty = false;
    };
    editor.addEventListener("input", function () {
      dirty = true;
    });
    // htmx reads the form before a submit listener would run
    document.getElementById("note-form").addEventListener("htmx:configRequest", function (e) {
      if (editor.hidden) return;
      writeBack();
      e.detail.parameters.body = raw.value;
    });

    // a click on a link opens its text and address to change, a drag across it only selects;
    // the link itself is never followed from here
    editor.addEventListener("click", function (e) {
      var a = e.target.closest("a");
      if (!a || !editor.contains(a)) return;
      e.preventDefault();
      var s = document.getSelection();
      if (!s.rangeCount || !s.isCollapsed) return;
      noteSelection = s.getRangeAt(0);
      openLink();
    });

    // a web page's styles and pictures stay behind, only its text comes along
    editor.addEventListener("paste", function (e) {
      e.preventDefault();
      document.execCommand("insertText", false, e.clipboardData.getData("text/plain"));
    });

    commands.forEach(function (b) {
      // a tap keeps the selection and the phone keyboard in the editor
      b.addEventListener("pointerdown", function (e) {
        e.preventDefault();
      });
      b.addEventListener("click", function () {
        run(b);
      });
    });

    linkAdd.addEventListener("click", addLink);
    linkRemove.addEventListener("click", function () {
      closeBubble();
      selectEditing();
      document.execCommand("unlink");
    });
    linkRow.querySelector("[data-link-cancel]").addEventListener("click", function () {
      closeBubble();
      restore();
    });
    linkURL.addEventListener("input", clearLinkError);
    // Enter in these would send the note
    [linkText, linkURL].forEach(function (field) {
      field.addEventListener("keydown", function (e) {
        if (e.key === "Enter") {
          e.preventDefault();
          addLink();
        }
      });
    });
    [bubble, linkRow].forEach(function (el) {
      el.addEventListener("keydown", function (e) {
        if (e.key === "Escape") {
          closeBubble();
          restore();
        }
      });
    });
    // a click anywhere else leaves the link as it was, the click itself goes where it was aimed
    document.addEventListener("pointerdown", function (e) {
      if (linkRow.isConnected && !linkRow.hidden && !linkRow.contains(e.target)) closeBubble();
    });

    markdown.addEventListener("switch-change", function (e) {
      e.preventDefault();
      if (e.detail.checked) {
        writeBack();
        show(false);
        raw.focus();
        return;
      }
      fetch(editor.dataset.render, { method: "POST", body: new URLSearchParams({ body: raw.value }) })
        .then(function (r) {
          if (!r.ok) throw new Error(r.status);
          return r.text();
        })
        .then(function (html) {
          editor.innerHTML = html;
          dirty = false;
          failed.hidden = true;
          show(true);
          editor.focus();
        })
        .catch(function () {
          failed.hidden = false;
        });
    });

    tools.hidden = false;
    show(true);
  });
}
