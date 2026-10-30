// [data-avatar-crop] in the picture form: after choosing a file it shows the whole picture with
// the circle housemates will see clear and the rest in shadow, drag (or arrow keys) moves the circle,
// the slider, wheel or a pinch resizes it,
// and Upload sends only that square; without JS the whole picture is sent and centred
// htmx runs this file again when a swap renders its component again
if (!window.avatarCropLoaded) {
  window.avatarCropLoaded = true;
  (() => {
    // the server shrinks it to 256px, twice that is enough for any screen
    const OUT = 512;

    const setup = (box) => {
      const form = box.closest("form");
      const input = form.querySelector('input[type="file"]');
      const [view, ...previews] = box.querySelectorAll("canvas");
      const zoom = box.querySelector('input[type="range"]');
      const img = new Image();
      // the shown square of the picture, in its pixels: centre and side
      let cx = 0;
      let cy = 0;
      let side = 0;
      let url = "";

      const clamp = () => {
        const short = Math.min(img.naturalWidth, img.naturalHeight);
        side = short / Number(zoom.value);
        cx = Math.min(Math.max(cx, side / 2), img.naturalWidth - side / 2);
        cy = Math.min(Math.max(cy, side / 2), img.naturalHeight - side / 2);
      };

      const paint = (canvas, size) => {
        const c = canvas.getContext("2d");
        canvas.width = size;
        canvas.height = size;
        // white under transparent pixels, like the server does for the JPEG
        c.fillStyle = "#fff";
        c.fillRect(0, 0, size, size);
        c.imageSmoothingQuality = "high";
        c.drawImage(img, cx - side / 2, cy - side / 2, side, side, 0, 0, size, size);
      };

      // shortcut: the stage is sized when a picture loads, not on window resize; redraw on resize if it matters
      const fitStage = () => {
        const s = Math.min(Math.min(view.parentElement.clientWidth, 384) / img.naturalWidth, 320 / img.naturalHeight);
        view.style.width = `${Math.round(img.naturalWidth * s)}px`;
        view.style.height = `${Math.round(img.naturalHeight * s)}px`;
        view.width = Math.round(view.clientWidth * devicePixelRatio);
        view.height = Math.round(view.clientHeight * devicePixelRatio);
      };

      // the whole picture, the circle that is kept clear and the rest in shadow
      const drawStage = () => {
        const c = view.getContext("2d");
        const k = view.width / img.naturalWidth;
        // white under transparent pixels like the upload, and so the last shadow does not show through
        c.fillStyle = "#fff";
        c.fillRect(0, 0, view.width, view.height);
        c.drawImage(img, 0, 0, view.width, view.height);
        c.beginPath();
        c.rect(0, 0, view.width, view.height);
        c.arc(cx * k, cy * k, (side / 2) * k, 0, 2 * Math.PI);
        c.fillStyle = "rgb(0 0 0 / 55%)";
        c.fill("evenodd");
        c.beginPath();
        c.arc(cx * k, cy * k, (side / 2) * k, 0, 2 * Math.PI);
        c.strokeStyle = "#fff";
        c.lineWidth = 2 * devicePixelRatio;
        c.stroke();
      };

      const draw = () => {
        clamp();
        drawStage();
        previews.forEach((canvas) => paint(canvas, Math.round(canvas.clientWidth * devicePixelRatio)));
      };

      // the circle follows the pointer, screen pixels to picture pixels
      const move = (dx, dy) => {
        const k = img.naturalWidth / view.clientWidth;
        cx += dx * k;
        cy += dy * k;
        draw();
      };

      input.addEventListener("change", () => {
        URL.revokeObjectURL(url);
        const file = input.files[0];
        box.hidden = true;
        if (!file) {
          return;
        }
        url = URL.createObjectURL(file);
        img.src = url;
      });

      // the Cancel button, it clears the file but a reset fires no change
      form.addEventListener("reset", () => {
        URL.revokeObjectURL(url);
        // a picture still loading must not open the box again
        img.removeAttribute("src");
        box.hidden = true;
      });

      img.addEventListener("load", () => {
        zoom.value = "1";
        cx = img.naturalWidth / 2;
        cy = img.naturalHeight / 2;
        box.hidden = false;
        fitStage();
        draw();
      });
      // not an image the browser can show, the server says what is wrong
      img.addEventListener("error", () => {
        box.hidden = true;
      });

      zoom.addEventListener("input", draw);

      view.addEventListener("wheel", (e) => {
        e.preventDefault();
        // by how far it scrolled, a trackpad sends many small events, a mouse notch is about 100
        zoom.value = String(Number(zoom.value) - e.deltaY * 0.002);
        draw();
      });

      // the fingers (or mouse) on the view by pointerId, one drags the circle,
      // two pinch it smaller or spread it bigger and move it with their midpoint
      const pointers = new Map();
      view.addEventListener("pointerdown", (e) => {
        view.setPointerCapture(e.pointerId);
        pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });
      });
      view.addEventListener("pointermove", (e) => {
        const [a, b] = [...pointers.keys()];
        if (e.pointerId !== a && e.pointerId !== b) {
          return;
        }
        const prev = pointers.get(e.pointerId);
        const next = { x: e.clientX, y: e.clientY };
        pointers.set(e.pointerId, next);
        if (b === undefined) {
          move(next.x - prev.x, next.y - prev.y);
          return;
        }
        const other = pointers.get(e.pointerId === a ? b : a);
        const before = Math.hypot(prev.x - other.x, prev.y - other.y);
        const after = Math.hypot(next.x - other.x, next.y - other.y);
        if (before > 0 && after > 0) {
          // the range input keeps it within 1 to 4
          zoom.value = String((Number(zoom.value) * before) / after);
        }
        // the midpoint moves half of what this finger moved
        move((next.x - prev.x) / 2, (next.y - prev.y) / 2);
      });
      // after every up or cancel, and when the browser takes the pointer away without either
      view.addEventListener("lostpointercapture", (e) => {
        pointers.delete(e.pointerId);
      });

      view.addEventListener("keydown", (e) => {
        const step = e.shiftKey ? 40 : 10;
        const d = { ArrowLeft: [-step, 0], ArrowRight: [step, 0], ArrowUp: [0, -step], ArrowDown: [0, step] }[e.key];
        if (d) {
          e.preventDefault();
          move(d[0], d[1]);
        }
      });

      form.addEventListener("submit", (e) => {
        if (box.hidden) {
          return;
        }
        e.preventDefault();
        const out = document.createElement("canvas");
        paint(out, OUT);
        out.toBlob(
          (blob) => {
            // without a blob (out of memory) or DataTransfer the whole picture goes as it was chosen
            if (blob) {
              try {
                const files = new DataTransfer();
                files.items.add(new File([blob], "avatar.jpg", { type: "image/jpeg" }));
                input.files = files.files;
              } catch {}
            }
            form.submit();
          },
          "image/jpeg",
          0.9,
        );
      });
    };

    htmx.onLoad((root) => {
      root.querySelectorAll("[data-avatar-crop]").forEach(setup);
    });
  })();
}
