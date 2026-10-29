package pages

import (
	"roommates/internal/services"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
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

// initial is the first letter of a name for its avatar
func initial(name string) string {
	r, _ := utf8.DecodeRuneInString(strings.TrimSpace(name))
	if r == utf8.RuneError {
		return "?"
	}
	return string(unicode.ToUpper(r))
}
