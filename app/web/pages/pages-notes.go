package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// hard wraps, a newline typed in the textarea is a line break
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithHardWraps()),
)

// noteHTML renders the markdown of a note, without html.WithUnsafe goldmark
// drops raw HTML and javascript: style links, so members cannot inject scripts
func noteHTML(src string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		return markdown.Convert([]byte(src), w)
	})
}
