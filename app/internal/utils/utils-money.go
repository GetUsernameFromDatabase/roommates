package utils

import (
	"context"

	"github.com/invopop/ctxi18n"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// FormatMoney is cents in the ISO 4217 code with the request language's separators, "12,50 €"
func FormatMoney(ctx context.Context, cents int64, code string) string {
	unit, err := currency.ParseISO(code)
	if err != nil {
		unit = currency.EUR
	}
	tag := language.Make(ctxi18n.Locale(ctx).Code().String())
	p := message.NewPrinter(tag)
	return p.Sprint(number.Decimal(float64(cents)/100, number.Scale(2))) + " " + p.Sprint(currency.Symbol(unit))
}
