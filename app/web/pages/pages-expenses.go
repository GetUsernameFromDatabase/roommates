package pages

import (
	"context"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// amountValue shows a form amount with the decimal comma outside English, like
// [utils.FormatMoney]; [services.ParseCents] reads either separator back
func amountValue(ctx context.Context, s string) string {
	if utils.GetLocaleString(ctx) == "en" {
		return s
	}
	return strings.Replace(s, ".", ",", 1)
}

// newExpenseURL opens the expense form with title filled in
func newExpenseURL(householdID, listID, itemID pgtype.UUID) templ.SafeURL {
	q := url.Values{"list": {listID.String()}}
	if itemID.Valid {
		q.Set("item", itemID.String())
	}
	return householdURL(householdID, "expenses", "new") + templ.SafeURL("?"+q.Encode())
}

// formatDay is a date column, it has no time zone
func formatDay(ctx context.Context, d pgtype.Date) string {
	return utils.FormatDate(ctx, d.Time)
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

var methodKeys = map[string]locales.LK{
	"transfer": locales.LK_ExpensesMethodTransfer,
	"cash":     locales.LK_ExpensesMethodCash,
	"card":     locales.LK_ExpensesMethodCard,
	"other":    locales.LK_ExpensesMethodOther,
}

// settlementLine is how a settlement was paid and its note, "" for one an edit made
func settlementLine(ctx context.Context, method *string, note string) string {
	var parts []string
	if method != nil {
		parts = append(parts, utils.T(ctx, methodKeys[*method]))
	}
	if note != "" {
		parts = append(parts, note)
	}
	return strings.Join(parts, " · ")
}

// repeatNever is the repeat select's value for an expense that does not repeat
const repeatNever = "never"

var intervalKeys = map[string]locales.LK{
	repeatNever:     locales.LK_RecurringNever,
	g.IntervalWeek:  locales.LK_RecurringWeek,
	g.IntervalMonth: locales.LK_RecurringMonth,
	g.IntervalYear:  locales.LK_RecurringYear,
}

// expenseFrame is what the forms of [expenseForm] differ in
type expenseFrame struct {
	Action, Back      templ.SafeURL
	BackText, Heading string
	Submit            locales.LK
	// Lock is the [services.Expense.AmountLock] of an expense being edited
	Lock string
	// PayerID gets the leftover cents of an equal split, see expense-split.js
	PayerID   pgtype.UUID
	DateLabel locales.LK
	// Repeats are the repeat select's options, none hides it, Repeat is the chosen one or the first
	Repeats []string
	Repeat  string
}

// canEditRecurring says whether the actor can change schedule r, see services.UpdateRecurring
func canEditRecurring(actor services.Membership, r services.RecurringRow) bool {
	return r.PayerID == actor.UserID || actor.Role == g.RoleOwner
}

// SettleSlip is a failed "I paid" shown again in the slip of To, with what was typed
type SettleSlip struct {
	To, Amount, Method, Note, Error string
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
		return utils.T(ctx, locales.LK_ExpensesPaidByYou, formatDay(ctx, day))
	}
	return utils.T(ctx, locales.LK_ExpensesPaidBy, payerName, formatDay(ctx, day))
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

// yoursFirst puts the debts the actor is part of before the others', each in their order
func yoursFirst(debts []services.Debt, actor services.Membership) []services.Debt {
	others := func(d services.Debt) bool { return d.From != actor.UserID && d.To != actor.UserID }
	sorted := slices.Clone(debts)
	slices.SortStableFunc(sorted, func(a, b services.Debt) int {
		return utils.IfElse(others(a), 1, 0) - utils.IfElse(others(b), 1, 0)
	})
	return sorted
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

// iPaidAmount fills the "I paid" form with what is left of the debt after the payments
// waiting to be confirmed, so accepting it cannot pay twice, empty when nothing is left
func iPaidAmount(d services.Debt, pending []dbqueries.ListPendingSettlementsRow) string {
	left := d.Cents
	for _, p := range pending {
		if p.FromUser == d.From && p.ToUser == d.To {
			left -= p.AmountCents
		}
	}
	if left <= 0 {
		return ""
	}
	return services.FormatCentsInput(left)
}
