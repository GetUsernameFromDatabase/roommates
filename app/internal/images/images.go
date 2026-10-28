// Package images makes small previews of uploaded pictures
package images

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	_ "image/gif" // registers the decoder for image.Decode
	"image/jpeg"
	_ "image/png" // registers the decoder for image.Decode
	"io"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // registers the decoder for image.Decode
)

// maxPixels is checked before decoding, a small file can claim a huge size
// and decoding it would take all memory; a var so tests can lower it
var maxPixels int64 = 50_000_000

var ErrTooLarge = errors.New("image has too many pixels")

// Thumbnail decodes a PNG, JPEG, GIF (first frame) or WebP and returns a JPEG whose
// longest side is at most maxSide, turned upright by the JPEG's EXIF orientation
func Thumbnail(r io.ReadSeeker, maxSide int) ([]byte, error) {
	cfg, format, err := image.DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > maxPixels {
		return nil, ErrTooLarge
	}
	orientation := 1
	if format == "jpeg" {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		orientation = jpegOrientation(r)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	src, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	b := src.Bounds()
	w, h := fit(b.Dx(), b.Dy(), maxSide)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	// white under transparent pixels, JPEG has no alpha and would show them black
	draw.Draw(dst, dst.Bounds(), image.White, image.Point{}, draw.Src)
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)

	var buf bytes.Buffer
	//exhaustruct:ignore
	if err := jpeg.Encode(&buf, orient(dst, orientation), &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// fit scales w x h down so the longest side is maxSide, never up
func fit(w, h, maxSide int) (int, int) {
	longest := max(w, h)
	if longest <= maxSide {
		return w, h
	}
	return max(1, w*maxSide/longest), max(1, h*maxSide/longest)
}

// orient turns src by an EXIF orientation (1..8) so it displays upright
func orient(src *image.RGBA, orientation int) *image.RGBA {
	if orientation < 2 || orientation > 8 {
		return src
	}
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dw, dh := w, h
	if orientation >= 5 {
		dw, dh = h, w
	}
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	for y := range dh {
		for x := range dw {
			// the source pixel shown at x, y
			var sx, sy int
			switch orientation {
			case 2: // mirrored
				sx, sy = w-1-x, y
			case 3: // upside down
				sx, sy = w-1-x, h-1-y
			case 4: // mirrored upside down
				sx, sy = x, h-1-y
			case 5: // mirrored, turned 90° counterclockwise
				sx, sy = y, x
			case 6: // turned 90° counterclockwise, needs 90° clockwise
				sx, sy = y, h-1-x
			case 7: // mirrored, turned 90° clockwise
				sx, sy = w-1-y, h-1-x
			case 8: // turned 90° clockwise, needs 90° counterclockwise
				sx, sy = w-1-y, x
			}
			dst.SetRGBA(x, y, src.RGBAAt(sx, sy))
		}
	}
	return dst
}

// jpegOrientation returns the EXIF orientation of a JPEG, 1 (upright) if it has none
// or anything is unexpected
func jpegOrientation(r io.Reader) int {
	br := bufio.NewReader(r)
	var soi [2]byte
	if _, err := io.ReadFull(br, soi[:]); err != nil || soi != [2]byte{0xFF, 0xD8} {
		return 1
	}
	for {
		// a segment is 0xFF, a marker byte, a 2 byte length that counts itself, the payload
		var head [4]byte
		if _, err := io.ReadFull(br, head[:]); err != nil || head[0] != 0xFF {
			return 1
		}
		marker := head[1]
		// start of the image data, EXIF always comes before it
		if marker == 0xDA || marker == 0xD9 {
			return 1
		}
		size := int(binary.BigEndian.Uint16(head[2:])) - 2
		if size < 0 {
			return 1
		}
		if marker != 0xE1 {
			if _, err := br.Discard(size); err != nil {
				return 1
			}
			continue
		}
		data := make([]byte, size)
		if _, err := io.ReadFull(br, data); err != nil {
			return 1
		}
		if exif, ok := bytes.CutPrefix(data, []byte("Exif\x00\x00")); ok {
			return exifOrientation(exif)
		}
	}
}

// exifOrientation reads tag 0x0112 from the first directory of a TIFF structure
func exifOrientation(t []byte) int {
	if len(t) < 8 {
		return 1
	}
	var order binary.ByteOrder
	switch string(t[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}
	ifd := int64(order.Uint32(t[4:]))
	if ifd+2 > int64(len(t)) {
		return 1
	}
	entries := int64(order.Uint16(t[ifd:]))
	for i := range entries {
		// an entry is tag, type, count (4 bytes), value (4 bytes)
		e := ifd + 2 + 12*i
		if e+12 > int64(len(t)) {
			return 1
		}
		if order.Uint16(t[e:]) == 0x0112 {
			// a SHORT value sits in the first 2 bytes of the value field
			if v := int(order.Uint16(t[e+8:])); v >= 1 && v <= 8 {
				return v
			}
			return 1
		}
	}
	return 1
}
