package models

import (
	"strings"
	"unicode"
)

// ibanLengths is the IBAN length per country, the SEPA area and the common others
var ibanLengths = map[string]int{
	"AD": 24, "AT": 20, "BE": 16, "BG": 22, "CH": 21, "CY": 28, "CZ": 24, "DE": 22,
	"DK": 18, "EE": 20, "ES": 24, "FI": 18, "FR": 27, "GB": 22, "GI": 23, "GR": 27,
	"HR": 21, "HU": 28, "IE": 22, "IS": 26, "IT": 27, "LI": 21, "LT": 20, "LU": 20,
	"LV": 21, "MC": 27, "MT": 31, "NL": 18, "NO": 15, "PL": 28, "PT": 25, "RO": 24,
	"SE": 24, "SI": 19, "SK": 24, "SM": 27, "VA": 22,
}

// NormalizeIBAN drops spaces and upper-cases
func NormalizeIBAN(s string) string {
	return strings.ToUpper(strings.Join(strings.FieldsFunc(s, unicode.IsSpace), ""))
}

// FormatIBAN is the normalized IBAN in groups of four
func FormatIBAN(s string) string {
	s = NormalizeIBAN(s)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && i%4 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ValidIBAN checks s, already normalized: a known country with its length, only A-Z and 0-9,
// and the ISO 7064 mod 97 checksum
func ValidIBAN(s string) bool {
	if len(s) < 4 || ibanLengths[s[:2]] != len(s) {
		return false
	}
	rem := 0
	// moving the first four characters to the end is the same as reading them last
	for _, r := range s[4:] + s[:4] {
		switch {
		case r >= '0' && r <= '9':
			rem = (rem*10 + int(r-'0')) % 97
		case r >= 'A' && r <= 'Z':
			rem = (rem*100 + int(r-'A') + 10) % 97
		default:
			return false
		}
	}
	return rem == 1
}
