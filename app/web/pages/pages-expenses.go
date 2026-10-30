package pages

import (
	"context"
	"net/url"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// newExpenseURL opens the expense form with title filled in
func newExpenseURL(householdID, listID, itemID pgtype.UUID) templ.SafeURL {
	q := url.Values{"list": {listID.String()}}
	if itemID.Valid {
		q.Set("item", itemID.String())
	}
	return householdURL(householdID, "expenses", "new") + templ.SafeURL("?"+q.Encode())
}

// formatDay is a date column, it has no time zone
func formatDay(d pgtype.Date) string {
	return d.Time.Format(time.DateOnly)
}

var categoryKeys = map[string]locales.LK{
	"groceries": locales.LK_ExpensesCategoryGroceries,
	"utilities": locales.LK_ExpensesCategoryUtilities,
	"rent":      locales.LK_ExpensesCategoryRent,
	"household": locales.LK_ExpensesCategoryHousehold,
	"other":     locales.LK_ExpensesCategoryOther,
}

// CategoryName is the translated expenses.category
func CategoryName(ctx context.Context, category string) string {
	return utils.T(ctx, categoryKeys[category])
}

// pairPhrase says something "from" did to "to", with "you" for the actor: youFrom takes
// the other's name when the actor is from, youTo when the actor is to, both takes both names
func pairPhrase(ctx context.Context, actor services.Membership, fromID pgtype.UUID, fromName string,
	toID pgtype.UUID, toName string, both, youFrom, youTo locales.LK) string {
	switch actor.UserID {
	case fromID:
		return utils.T(ctx, youFrom, toName)
	case toID:
		return utils.T(ctx, youTo, fromName)
	default:
		return utils.T(ctx, both, fromName, toName)
	}
}

// paidByLine is "Paid by <payer> on <day>", or by you
func paidByLine(ctx context.Context, actor services.Membership, payerID pgtype.UUID, payerName string, day pgtype.Date) string {
	if payerID == actor.UserID {
		return utils.T(ctx, locales.LK_ExpensesPaidByYou, formatDay(day))
	}
	return utils.T(ctx, locales.LK_ExpensesPaidBy, payerName, formatDay(day))
}

// personName is name, or "You" for the actor
func personName(ctx context.Context, actor services.Membership, id pgtype.UUID, name string) string {
	if id == actor.UserID {
		return utils.T(ctx, locales.LK_ExpensesYou)
	}
	return name
}

// formatSignedMoney is [utils.FormatMoney] with a + or a minus sign
func formatSignedMoney(ctx context.Context, cents int64, code string) string {
	if cents < 0 {
		return "−" + utils.FormatMoney(ctx, -cents, code)
	}
	return "+" + utils.FormatMoney(ctx, cents, code)
}

// itemClass colours a breakdown amount green when it is in the actor's favour and red
// when it is against them, a debt between two others stays uncoloured
func itemClass(actor services.Membership, d services.Debt, cents int64) string {
	switch {
	case actor.UserID != d.From && actor.UserID != d.To:
		return ""
	case (cents > 0) == (actor.UserID == d.From):
		return "text-against"
	default:
		return "text-favour"
	}
}
