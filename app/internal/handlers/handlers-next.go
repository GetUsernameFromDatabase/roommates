package handlers

import (
	"strings"
	"unicode"
)

// safeNextPath returns next if it is a path on this site, otherwise "/"
//
// "//evil.com" and "/\evil.com" are treated by browsers as other hosts,
// browsers also drop tabs and newlines, so "/<tab>/evil.com" becomes "//evil.com"
func safeNextPath(next string) string {
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") ||
		strings.HasPrefix(next, "/\\") || strings.IndexFunc(next, unicode.IsControl) >= 0 {
		return "/"
	}
	return next
}
