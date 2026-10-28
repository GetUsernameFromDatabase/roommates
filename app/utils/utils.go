// Package utils a collection of helper functions used throughout the codebase
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	_ "github.com/a-h/templ" // needed for ReplaceParam doc
)

// RandomIDForHTML for HTML element
func RandomIDForHTML(prefix string) string {
	b := make([]byte, 3) // 3 bytes = 6 hex chars
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return prefix + "-" + hex.EncodeToString(b)
}

// ReplaceParam replaces url param with desired value, will only replace the first occurrence
//
// does not sanitize the url (use [templ.URL] should it be needed ex: hx-get)
// [templ url-attributes]
//
// no need to add : in front of the param
//
//	Example: (url="example/:id/example", param="id", value="bb")
//	 Return: "example/bb/example"
//
// [templ url-attributes]: https://templ.guide/syntax-and-usage/attributes/#url-attributes
func ReplaceParam(url, param, value string) string {
	return strings.Replace(url, ":"+param, value, 1)
}
