package utils

import (
	"context"

	"golang.org/x/text/currency"
	"golang.org/x/text/number"
)

// FormatMoney is cents in the ISO 4217 code with the request language's separators,
// "€12.50" in English, "12,50 €" in Estonian, see utils-format.go
func FormatMoney(ctx context.Context, cents int64, code string) string {
	unit, err := currency.ParseISO(code)
	if err != nil {
		unit = currency.EUR
	}
	p := printer(ctx)
	amount, symbol := p.Sprint(number.Decimal(float64(cents)/100, number.Scale(2))), p.Sprint(currency.Symbol(unit))
	if GetLocaleString(ctx) == "en" {
		return symbol + amount
	}
	return amount + " " + symbol
}
