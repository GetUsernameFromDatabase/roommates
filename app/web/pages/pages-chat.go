package pages

import (
	g "roommates/internal/config"
	"roommates/internal/services"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// runGap is how far apart messages of one author can be and still read as one run
const runGap = 5 * time.Minute

// sameRun is true when b follows a without its author's name and time being shown again
func sameRun(a, b services.Message) bool {
	d := a.CreatedAt.Time.Sub(b.CreatedAt.Time).Abs()
	return a.AuthorID == b.AuthorID && d < runGap
}

// startsRun and endsRun are about messages ms, newest first, so the older neighbour of i is i+1
func startsRun(ms []services.Message, i int) bool {
	return i+1 >= len(ms) || !sameRun(ms[i], ms[i+1])
}

func endsRun(ms []services.Message, i int) bool {
	return i == 0 || !sameRun(ms[i], ms[i-1])
}

// chatReplySwap makes a link to url (the chat, ?reply= answers a message) swap only the composer's
// reply banner, so the typed text and picked files stay
func chatReplySwap(url templ.SafeURL) templ.Attributes {
	return templ.Attributes{
		"hx-get":    string(url),
		"hx-target": "#chat-reply",
		"hx-select": "#chat-reply",
		"hx-swap":   "outerHTML",
	}
}

// splitFiles is the images of a message (the files with a thumbnail) and its other files, each in order
func splitFiles(fs []services.MessageFile) (images, others []services.MessageFile) {
	for _, f := range fs {
		if f.HasThumb {
			images = append(images, f)
		} else {
			others = append(others, f)
		}
	}
	return images, others
}

// avatarURL is where a user's picture is, the key's random end makes a new picture a new URL
func avatarURL(userID pgtype.UUID, key string) string {
	return g.RouteUsers + "/" + userID.String() + "/avatar?v=" + key[strings.LastIndexByte(key, '/')+1:]
}

// initial is the first letter of a name for its avatar
func initial(name string) string {
	r, _ := utf8.DecodeRuneInString(strings.TrimSpace(name))
	if r == utf8.RuneError {
		return "?"
	}
	return string(unicode.ToUpper(r))
}
