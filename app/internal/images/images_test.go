package images

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func testImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 0, A: 255})
		}
	}
	return img
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// jpegWithOrientation is a JPEG with an APP1 EXIF segment right after the start marker
func jpegWithOrientation(t *testing.T, img image.Image, orientation uint16) []byte {
	t.Helper()
	var buf bytes.Buffer
	//exhaustruct:ignore
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	be := binary.BigEndian
	exif := []byte("Exif\x00\x00MM\x00\x2a")
	exif = be.AppendUint32(exif, 8) // first directory right after the header
	exif = be.AppendUint16(exif, 1) // one entry
	exif = be.AppendUint16(exif, 0x0112)
	exif = be.AppendUint16(exif, 3) // SHORT
	exif = be.AppendUint32(exif, 1)
	exif = be.AppendUint16(exif, orientation)
	exif = append(exif, 0, 0, 0, 0, 0, 0) // value padding, no next directory
	segment := be.AppendUint16([]byte{0xFF, 0xE1}, uint16(len(exif)+2))
	segment = append(segment, exif...)
	b := buf.Bytes()
	return append(append(append([]byte{}, b[:2]...), segment...), b[2:]...)
}

func thumbSize(t *testing.T, data []byte, maxSide int) (int, int) {
	t.Helper()
	thumb, err := Thumbnail(bytes.NewReader(data), maxSide)
	if err != nil {
		t.Fatal(err)
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(thumb))
	if err != nil || format != "jpeg" {
		t.Fatalf("thumbnail is %q, %v", format, err)
	}
	return cfg.Width, cfg.Height
}

func TestThumbnail(t *testing.T) {
	if w, h := thumbSize(t, encodePNG(t, testImage(600, 300)), 256); w != 256 || h != 128 {
		t.Errorf("png thumbnail %dx%d, want 256x128", w, h)
	}
	// small images are not scaled up
	if w, h := thumbSize(t, jpegWithOrientation(t, testImage(40, 20), 1), 256); w != 40 || h != 20 {
		t.Errorf("small jpeg thumbnail %dx%d, want 40x20", w, h)
	}
	// orientation 6 is a portrait photo stored on its side
	if w, h := thumbSize(t, jpegWithOrientation(t, testImage(40, 20), 6), 256); w != 20 || h != 40 {
		t.Errorf("rotated jpeg thumbnail %dx%d, want 20x40", w, h)
	}

	if _, err := Thumbnail(bytes.NewReader([]byte("not an image")), 256); err == nil {
		t.Error("garbage should fail")
	}

	old := maxPixels
	maxPixels = 100
	defer func() { maxPixels = old }()
	if _, err := Thumbnail(bytes.NewReader(encodePNG(t, testImage(20, 20))), 256); !errors.Is(err, ErrTooLarge) {
		t.Errorf("too many pixels: %v", err)
	}
}

func TestOrient(t *testing.T) {
	// a b
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	a, b := color.RGBA{R: 1, G: 0, B: 0, A: 255}, color.RGBA{R: 2, G: 0, B: 0, A: 255}
	src.SetRGBA(0, 0, a)
	src.SetRGBA(1, 0, b)
	for orientation, want := range map[int][]color.RGBA{
		1: {a, b}, // rows of the result, read left to right, top to bottom
		2: {b, a},
		3: {b, a},
		4: {a, b},
		5: {a, b},
		6: {a, b},
		7: {b, a},
		8: {b, a},
	} {
		dst := orient(src, orientation)
		var got []color.RGBA
		for y := range dst.Bounds().Dy() {
			for x := range dst.Bounds().Dx() {
				got = append(got, dst.RGBAAt(x, y))
			}
		}
		if wantVertical := orientation >= 5; wantVertical != (dst.Bounds().Dy() == 2) {
			t.Errorf("orientation %d: size %v", orientation, dst.Bounds())
		}
		if got[0] != want[0] || got[1] != want[1] {
			t.Errorf("orientation %d: got %v, want %v", orientation, got, want)
		}
	}
}
