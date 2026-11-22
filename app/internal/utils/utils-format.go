package utils

import (
	"context"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// defaultDateLayout is the short date of Estonian and any other language, 19.11.2026
const defaultDateLayout = "02.01.2006"

// dateLayouts is the short date of each language that differs from [defaultDateLayout]
var dateLayouts = map[string]string{"en": "2 Jan 2006"}

// FormatDate is the calendar date of t in the request language, "19.11.2026" or "19 Nov 2026"
func FormatDate(ctx context.Context, t time.Time) string {
	if layout, ok := dateLayouts[GetLocaleString(ctx)]; ok {
		return t.Format(layout)
	}
	return t.Format(defaultDateLayout)
}

func printer(ctx context.Context) *message.Printer {
	return message.NewPrinter(language.Make(GetLocaleString(ctx)))
}

// FormatSize is a byte count in B, KB or MB with the request language's decimal separator
func FormatSize(ctx context.Context, n int64) string {
	p := printer(ctx)
	switch {
	case n < 1<<10:
		return p.Sprint(number.Decimal(n)) + " B"
	case n < 1<<20:
		return p.Sprint(number.Decimal(float64(n)/(1<<10), number.Scale(0))) + " KB"
	default:
		return p.Sprint(number.Decimal(float64(n)/(1<<20), number.Scale(1))) + " MB"
	}
}
