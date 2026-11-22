package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

// CSV export formats of the expenses
const (
	ExportWide = "wide"
	ExportLong = "long"
)

// ExportExpenses is the expenses as CSV in the request language, an owner gets all of them,
// a member those they paid or have a share of, format is [ExportWide] or [ExportLong]
func (s *DefaultService) ExportExpenses(ctx context.Context, actor Membership, format string) ([]byte, error) {
	if format != ExportWide && format != ExportLong {
		return nil, g.ErrorInvalidInput
	}
	scope := pgtype.UUID{}
	if actor.Role != g.RoleOwner {
		scope = actor.UserID
	}
	expenses, err := s.db.ListExportExpenses(ctx, dbqueries.ListExportExpensesParams{HouseholdID: actor.HouseholdID, UserID: scope})
	if err != nil {
		return nil, err
	}
	shares, err := s.db.ListExportShares(ctx, dbqueries.ListExportSharesParams{HouseholdID: actor.HouseholdID, UserID: scope})
	if err != nil {
		return nil, err
	}
	return writeCSV(ctx, expenseRows(ctx, format, expenses, shares))
}

// ExportBalances is who is owed what as CSV in the request language: an owner gets every active member's
// net balance, a member what each person they have open debts with owes them
func (s *DefaultService) ExportBalances(ctx context.Context, actor Membership) ([]byte, error) {
	owner := actor.Role == g.RoleOwner
	scope := actor.UserID
	if owner {
		scope = pgtype.UUID{}
	}
	items, err := s.db.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{HouseholdID: actor.HouseholdID, UserID: scope})
	if err != nil {
		return nil, err
	}
	var members []dbqueries.ListActiveMembersRow
	if owner {
		if members, err = s.db.ListActiveMembers(ctx, actor.HouseholdID); err != nil {
			return nil, err
		}
	}
	return writeCSV(ctx, balanceRows(ctx, actor, items, members))
}

// csvComma is ';' where the decimal separator is a comma, as spreadsheets read it
func csvComma(ctx context.Context) rune {
	if utils.GetLocaleString(ctx) == "en" {
		return ','
	}
	return ';'
}

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// writeCSV starts with a BOM, which makes Excel read UTF-8
func writeCSV(ctx context.Context, rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(utf8BOM)
	w := csv.NewWriter(&buf)
	w.Comma = csvComma(ctx)
	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// csvAmount is cents as a plain number with the request language's decimal separator
func csvAmount(ctx context.Context, cents int64) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}
	s := fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
	if utils.GetLocaleString(ctx) == "en" {
		return s
	}
	return strings.Replace(s, ".", ",", 1)
}

// csvText keeps a spreadsheet from running user text as a formula
func csvText(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}

func expenseRows(ctx context.Context, format string, expenses []dbqueries.ListExportExpensesRow, shares []dbqueries.ListExportSharesRow) [][]string {
	t := func(k locales.LK) string { return utils.T(ctx, k) }
	head := []string{t(locales.LK_ExportDate), t(locales.LK_ExportTitle), t(locales.LK_ExportCategory), t(locales.LK_ExportPayer)}
	byExpense := map[pgtype.UUID][]dbqueries.ListExportSharesRow{}
	for _, sh := range shares {
		byExpense[sh.ExpenseID] = append(byExpense[sh.ExpenseID], sh)
	}
	first := func(e dbqueries.ListExportExpensesRow) []string {
		return []string{
			e.IncurredOn.Time.Format(time.DateOnly), csvText(e.Title), csvText(utils.T(ctx, locales.LK("expenses.category-"+e.Category))), csvText(utils.MemberName(ctx, e.PayerName)),
		}
	}

	if format == ExportLong {
		rows := [][]string{append(head, t(locales.LK_ExportTotal), t(locales.LK_ExportCurrency),
			t(locales.LK_ExportMember), t(locales.LK_ExportShare), t(locales.LK_ExportPaid))}
		for _, e := range expenses {
			for _, sh := range byExpense[e.ID] {
				paid := t(locales.LK_ExportNo)
				if sh.Paid {
					paid = t(locales.LK_ExportYes)
				}
				rows = append(rows, append(first(e), csvAmount(ctx, e.AmountCents), e.Currency,
					csvText(utils.MemberName(ctx, sh.Username)), csvAmount(ctx, sh.AmountCents), paid))
			}
		}
		return rows
	}

	// shares come ordered by name, a column per member in that order
	var columns []pgtype.UUID
	header := append(head, t(locales.LK_ExportTotal), t(locales.LK_ExportCurrency))
	for _, sh := range shares {
		if !slices.Contains(columns, sh.UserID) {
			columns = append(columns, sh.UserID)
			header = append(header, csvText(utils.MemberName(ctx, sh.Username)))
		}
	}
	rows := [][]string{header}
	for _, e := range expenses {
		row := append(first(e), csvAmount(ctx, e.AmountCents), e.Currency)
		for _, c := range columns {
			cell := ""
			for _, sh := range byExpense[e.ID] {
				if sh.UserID == c {
					cell = csvAmount(ctx, sh.AmountCents)
				}
			}
			row = append(row, cell)
		}
		rows = append(rows, row)
	}
	return rows
}

func balanceRows(ctx context.Context, actor Membership, items []dbqueries.ListDebtItemsRow, members []dbqueries.ListActiveMembersRow) [][]string {
	owner := actor.Role == g.RoleOwner
	net := map[pgtype.UUID]int64{}
	names := map[pgtype.UUID]string{}
	for _, m := range members {
		names[m.UserID] = m.Username
	}
	for _, it := range items {
		// from owes to
		names[it.FromUser], names[it.ToUser] = it.FromName, it.ToName
		switch {
		case owner:
			net[it.FromUser] -= it.AmountCents
			net[it.ToUser] += it.AmountCents
		case it.ToUser == actor.UserID:
			net[it.FromUser] += it.AmountCents
		case it.FromUser == actor.UserID:
			net[it.ToUser] -= it.AmountCents
		}
	}
	// the owner lists every active member, also at zero
	ids := make([]pgtype.UUID, 0, len(net))
	if owner {
		for _, m := range members {
			ids = append(ids, m.UserID)
		}
	} else {
		for id, cents := range net {
			if cents != 0 {
				ids = append(ids, id)
			}
		}
	}
	slices.SortFunc(ids, func(a, b pgtype.UUID) int {
		if c := strings.Compare(strings.ToLower(names[a]), strings.ToLower(names[b])); c != 0 {
			return c
		}
		return bytes.Compare(a.Bytes[:], b.Bytes[:])
	})

	who, amount := locales.LK_ExportCounterpart, locales.LK_ExportOwedToYou
	if owner {
		who, amount = locales.LK_ExportMember, locales.LK_ExportBalance
	}
	rows := [][]string{{utils.T(ctx, who), utils.T(ctx, amount), utils.T(ctx, locales.LK_ExportCurrency)}}
	for _, id := range ids {
		rows = append(rows, []string{csvText(utils.MemberName(ctx, names[id])), csvAmount(ctx, net[id]), actor.Currency})
	}
	return rows
}
