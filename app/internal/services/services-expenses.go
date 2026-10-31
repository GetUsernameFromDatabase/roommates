package services

import (
	"bytes"
	"cmp"
	"context"
	"regexp"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/utils"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxExpenseTitleLength = 100
	// 10 000 000.00 in the household currency
	MaxAmountCents = 1_000_000_000
	DateLayout     = "2006-01-02"
)

var amountPattern = regexp.MustCompile(`^(\d{1,10})(?:[.,](\d{1,2}))?$`)

// ParseCents reads "12", "12.5" or "12,50" as minor units,
// [g.ErrorInvalidInput] unless it is 0 < amount <= [MaxAmountCents]
func ParseCents(s string) (int64, error) {
	m := amountPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, g.ErrorInvalidInput
	}
	whole, _ := strconv.ParseInt(m[1], 10, 64)
	frac, _ := strconv.ParseInt((m[2] + "00")[:2], 10, 64)
	cents := whole*100 + frac
	if cents <= 0 || cents > MaxAmountCents {
		return 0, g.ErrorInvalidInput
	}
	return cents, nil
}

// FormatCentsInput is cents as a form value that [ParseCents] reads back
func FormatCentsInput(cents int64) string {
	return strconv.FormatInt(cents/100, 10) + "." + strconv.FormatInt(cents%100+100, 10)[1:]
}

// splitEqual divides total between participants, the leftover cents go to the payer,
// or one each to the first participants when the payer is not one of them,
// web/assets/public/expense-split.js shows the same sums on the form
func splitEqual(total int64, participants []pgtype.UUID, payer pgtype.UUID) []int64 {
	n := int64(len(participants))
	shares := make([]int64, n)
	left := total % n
	payerAt := slices.Index(participants, payer)
	for i := range shares {
		shares[i] = total / n
		if payerAt < 0 && int64(i) < left {
			shares[i]++
		}
	}
	if payerAt >= 0 {
		shares[payerAt] += left
	}
	return shares
}

// ExpenseInput is what the expense form sends, amounts as typed
type ExpenseInput struct {
	Title      string
	Amount     string
	Category   string
	IncurredOn string
	// Custom splits by Shares, otherwise equally between Members
	Custom  bool
	Members []pgtype.UUID
	Shares  map[pgtype.UUID]string
	// ListID links the list's bought items to the expense, only ItemID when that is set too
	ListID pgtype.UUID
	ItemID pgtype.UUID
}

type expenseShare struct {
	userID pgtype.UUID
	cents  int64
}

// clean checks in against the active members, [g.ErrorInvalidInput] if anything is off
func (in ExpenseInput) clean(payer pgtype.UUID, members []dbqueries.ListActiveMembersRow) (dbqueries.InsertExpenseParams, []expenseShare, error) {
	var p dbqueries.InsertExpenseParams
	title, okTitle := trimToLength(in.Title, 1, MaxExpenseTitleLength)
	total, errAmount := ParseCents(in.Amount)
	date, errDate := time.Parse(DateLayout, strings.TrimSpace(in.IncurredOn))
	if !okTitle || errAmount != nil || errDate != nil || !slices.Contains(g.ExpenseCategories, in.Category) {
		return p, nil, g.ErrorInvalidInput
	}
	active := func(id pgtype.UUID) bool {
		return slices.ContainsFunc(members, func(m dbqueries.ListActiveMembersRow) bool { return m.UserID == id })
	}

	var shares []expenseShare
	if in.Custom {
		var sum int64
		for _, m := range members {
			text := strings.TrimSpace(in.Shares[m.UserID])
			if text == "" {
				continue
			}
			cents, err := ParseCents(text)
			if err != nil {
				return p, nil, err
			}
			sum += cents
			shares = append(shares, expenseShare{userID: m.UserID, cents: cents})
		}
		if sum != total {
			return p, nil, g.ErrorInvalidInput
		}
	} else {
		var ids []pgtype.UUID
		for _, id := range in.Members {
			if !active(id) || slices.Contains(ids, id) {
				return p, nil, g.ErrorInvalidInput
			}
			ids = append(ids, id)
		}
		for i, cents := range splitEqual(total, ids, payer) {
			shares = append(shares, expenseShare{userID: ids[i], cents: cents})
		}
	}
	if len(shares) == 0 {
		return p, nil, g.ErrorInvalidInput
	}
	p = dbqueries.InsertExpenseParams{
		HouseholdID: pgtype.UUID{},
		PayerID:     payer,
		Category:    in.Category,
		Title:       title,
		AmountCents: total,
		IncurredOn:  pgtype.Date{Time: date, InfinityModifier: pgtype.Finite, Valid: true},
	}
	return p, shares, nil
}

// CreateExpense records an expense the actor paid, see [ExpenseInput.clean] for errors
func (s *DefaultService) CreateExpense(ctx context.Context, actor Membership, in ExpenseInput) (pgtype.UUID, error) {
	var expenseID pgtype.UUID
	var p dbqueries.InsertExpenseParams
	var shares []expenseShare
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		members, err := q.ListActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		p, shares, err = in.clean(actor.UserID, members)
		if err != nil {
			return err
		}
		p.HouseholdID = actor.HouseholdID
		if expenseID, err = q.InsertExpense(ctx, p); err != nil {
			return err
		}
		for _, sh := range shares {
			if err := q.InsertExpenseShare(ctx, dbqueries.InsertExpenseShareParams{
				ExpenseID:   expenseID,
				UserID:      sh.userID,
				AmountCents: sh.cents,
			}); err != nil {
				return err
			}
		}
		if !in.ListID.Valid {
			return nil
		}
		return q.LinkShoppingItems(ctx, dbqueries.LinkShoppingItemsParams{
			ExpenseID:   expenseID,
			HouseholdID: actor.HouseholdID,
			ListID:      in.ListID,
			ItemID:      in.ItemID,
		})
	})
	if err != nil {
		return expenseID, err
	}
	owed := map[pgtype.UUID]int64{}
	ids := make([]pgtype.UUID, 0, len(shares))
	for _, sh := range shares {
		if sh.userID != actor.UserID {
			owed[sh.userID] = sh.cents
			ids = append(ids, sh.userID)
		}
	}
	s.notifyAsync(ctx, actor, ids, func(ctx context.Context, userID pgtype.UUID, payer, household string) push.Message {
		return push.Message{
			Title: utils.T(ctx, locales.LK_PushExpenseTitle, household),
			Body:  utils.T(ctx, locales.LK_PushExpenseBody, payer, p.Title, utils.FormatMoney(ctx, owed[userID], actor.Currency)),
			URL:   householdPath(actor.HouseholdID, "expenses", expenseID.String()),
			Tag:   "expense-" + expenseID.String(),
		}
	})
	return expenseID, nil
}

// ListMembers returns the active members, who an expense can be split between
func (s *DefaultService) ListMembers(ctx context.Context, actor Membership) ([]dbqueries.ListActiveMembersRow, error) {
	return s.db.ListActiveMembers(ctx, actor.HouseholdID)
}

// ExpenseShowAll is the [ExpenseFilter] Show that lists paid expenses too
const ExpenseShowAll = "all"

// ExpenseFilter narrows the expense list, Query is a part of the title in any case,
// only those someone still owes a part of are shown unless Show is [ExpenseShowAll]
type ExpenseFilter struct {
	Query string
	Show  string
}

// Active says whether f differs from the default list
func (f ExpenseFilter) Active() bool {
	return strings.TrimSpace(f.Query) != "" || f.Show == ExpenseShowAll
}

type ExpensesOverview struct {
	Debts   []Debt
	Pending []dbqueries.ListPendingSettlementsRow
	// Expenses match the filter, HasExpenses says whether the household has any at all
	Expenses    []dbqueries.ListExpensesRow
	HasExpenses bool
}

func (s *DefaultService) GetExpensesOverview(ctx context.Context, actor Membership, filter ExpenseFilter) (ExpensesOverview, error) {
	var o ExpensesOverview
	var err error
	if o.Debts, err = s.Debts(ctx, actor); err != nil {
		return o, err
	}
	if o.Pending, err = s.ListPendingSettlements(ctx, actor); err != nil {
		return o, err
	}
	if o.HasExpenses, err = s.db.HouseholdHasExpenses(ctx, actor.HouseholdID); err != nil {
		return o, err
	}
	o.Expenses, err = s.db.ListExpenses(ctx, dbqueries.ListExpensesParams{
		HouseholdID: actor.HouseholdID,
		Q:           strings.TrimSpace(filter.Query),
		OpenOnly:    filter.Show != ExpenseShowAll,
	})
	return o, err
}

type Expense struct {
	Expense dbqueries.GetExpenseRow
	Shares  []dbqueries.ListExpenseSharesRow
	Files   []AttachedFile
}

// GetExpense returns [g.ErrorNotFound] if the expense is not in the actor's household
func (s *DefaultService) GetExpense(ctx context.Context, actor Membership, expenseID pgtype.UUID) (Expense, error) {
	var e Expense
	var err error
	e.Expense, err = s.db.GetExpense(ctx, dbqueries.GetExpenseParams{ID: expenseID, HouseholdID: actor.HouseholdID})
	if err != nil {
		return e, noRowsAs(err, g.ErrorNotFound)
	}
	if e.Shares, err = s.db.ListExpenseShares(ctx, expenseID); err != nil {
		return e, err
	}
	e.Files, err = s.AttachedFiles(ctx, actor, FileTarget{
		NoteID:         pgtype.UUID{},
		ListID:         pgtype.UUID{},
		ShoppingItemID: pgtype.UUID{},
		ExpenseID:      expenseID,
		MessageID:      pgtype.UUID{},
	})
	return e, err
}

// CopyInput fills the expense form from e for today, as an equal split
// when e was split equally, so a new amount can be typed in
func (e Expense) CopyInput() ExpenseInput {
	in := ExpenseInput{
		Title:      e.Expense.Title,
		Amount:     FormatCentsInput(e.Expense.AmountCents),
		Category:   e.Expense.Category,
		IncurredOn: time.Now().Format(DateLayout),
		Custom:     false,
		Members:    nil,
		Shares:     map[pgtype.UUID]string{},
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
	}
	var cents []int64
	for _, sh := range e.Shares {
		in.Members = append(in.Members, sh.UserID)
		in.Shares[sh.UserID] = FormatCentsInput(sh.AmountCents)
		cents = append(cents, sh.AmountCents)
	}
	in.Custom = !slices.Equal(cents, splitEqual(e.Expense.AmountCents, in.Members, e.Expense.PayerID))
	return in
}

// DeleteExpense returns [g.ErrorNotFound] unless the actor paid the expense
// and everyone sharing it is still a member
func (s *DefaultService) DeleteExpense(ctx context.Context, actor Membership, expenseID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteExpense(ctx, dbqueries.DeleteExpenseParams{
		ID:          expenseID,
		HouseholdID: actor.HouseholdID,
		PayerID:     actor.UserID,
	}))
}

// DeleteExpenseKeepingFiles is [DefaultService.DeleteExpense] that leaves the files only the
// expense holds in Files, a refused delete keeps nothing
func (s *DefaultService) DeleteExpenseKeepingFiles(ctx context.Context, actor Membership, expenseID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		err := q.KeepExpenseFiles(ctx, dbqueries.KeepExpenseFilesParams{ExpenseID: expenseID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return err
		}
		return rowsAffected(q.DeleteExpense(ctx, dbqueries.DeleteExpenseParams{
			ID:          expenseID,
			HouseholdID: actor.HouseholdID,
			PayerID:     actor.UserID,
		}))
	})
}

// SetExpenseSharePaid lets the actor who paid the expense mark the share of userID, who must
// still be a member, paid or not, this spends or gives back what userID has paid them,
// the rest of the share counts as paid outside the app
//
// the debtor marks their own share paid only when confirmed payments cover all of it,
// anything else is [g.ErrorNotFound]
func (s *DefaultService) SetExpenseSharePaid(ctx context.Context, actor Membership, expenseID, userID pgtype.UUID, paid bool) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		sh, err := q.GetShareToMark(ctx, dbqueries.GetShareToMarkParams{
			ExpenseID:   expenseID,
			UserID:      userID,
			HouseholdID: actor.HouseholdID,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		debtor := userID == actor.UserID && sh.PayerID != actor.UserID
		if (sh.PayerID != actor.UserID && !debtor) || (debtor && !paid) {
			return g.ErrorNotFound
		}
		if paid {
			if sh.PaidAt.Valid {
				return nil
			}
			uncovered, err := payShare(ctx, q, actor.HouseholdID, expenseID, userID, sh.PayerID, sh.AmountCents)
			if err != nil {
				return err
			}
			if debtor && uncovered > 0 {
				return g.ErrorNotFound
			}
			return nil
		}
		if err := rowsAffected(q.SetExpenseSharePaid(ctx, dbqueries.SetExpenseSharePaidParams{
			Paid:        false,
			ExpenseID:   expenseID,
			UserID:      userID,
			HouseholdID: actor.HouseholdID,
			PayerID:     sh.PayerID,
		})); err != nil {
			return err
		}
		return q.DeleteSettlementShares(ctx, dbqueries.DeleteSettlementSharesParams{ExpenseID: expenseID, UserID: userID})
	})
}

// payShare marks an unpaid share of cents paid and spends userID's confirmed payments to
// payerID on it, oldest first, it returns how much of the share they did not cover
func payShare(ctx context.Context, q *dbqueries.Queries, householdID, expenseID, userID, payerID pgtype.UUID, cents int64) (int64, error) {
	if err := q.LockSettlementCredit(ctx, dbqueries.LockSettlementCreditParams{
		HouseholdID: householdID,
		FromUser:    userID,
		ToUser:      payerID,
	}); err != nil {
		return cents, err
	}
	credits, err := q.ListSettlementCredit(ctx, dbqueries.ListSettlementCreditParams{
		HouseholdID: householdID,
		FromUser:    userID,
		ToUser:      payerID,
	})
	if err != nil {
		return cents, err
	}
	for _, c := range credits {
		spend := min(cents, c.LeftCents)
		if spend <= 0 {
			continue
		}
		if err := q.InsertSettlementShare(ctx, dbqueries.InsertSettlementShareParams{
			SettlementID: c.ID,
			ExpenseID:    expenseID,
			UserID:       userID,
			Cents:        spend,
		}); err != nil {
			return cents, err
		}
		cents -= spend
	}
	return cents, rowsAffected(q.SetExpenseSharePaid(ctx, dbqueries.SetExpenseSharePaidParams{
		Paid:        true,
		ExpenseID:   expenseID,
		UserID:      userID,
		HouseholdID: householdID,
		PayerID:     payerID,
	}))
}

// closePaidDebt marks the shares paid once from no longer owes to anything, after a payment
// from from to to: when the two are exactly even all shares both ways, when from overpaid
// only those of from's shares the payments cover, marking more would change the balance
func closePaidDebt(ctx context.Context, q *dbqueries.Queries, householdID, from, to pgtype.UUID) error {
	rows, err := q.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{HouseholdID: householdID, UserID: from})
	if err != nil {
		return err
	}
	var net, credit int64
	var owed, owing []dbqueries.ListDebtItemsRow
	for _, r := range rows {
		switch {
		case r.FromUser == from && r.ToUser == to:
			net += r.AmountCents
			if r.ExpenseID.Valid {
				owed = append(owed, r)
			}
		case r.FromUser == to && r.ToUser == from:
			net -= r.AmountCents
			if r.ExpenseID.Valid {
				owing = append(owing, r)
			} else {
				credit += r.AmountCents
			}
		}
	}
	if net > 0 {
		return nil
	}
	if net == 0 {
		owed = append(owed, owing...)
	} else {
		covered := 0
		for covered < len(owed) && owed[covered].AmountCents <= credit {
			credit -= owed[covered].AmountCents
			covered++
		}
		owed = owed[:covered]
	}
	// all shares before any settlement, marking one locks them in that order too,
	// the list was read without a lock, someone can have marked a share since
	unpaid := make([]dbqueries.ListDebtItemsRow, 0, len(owed))
	for _, r := range owed {
		sh, err := q.GetShareToMark(ctx, dbqueries.GetShareToMarkParams{
			ExpenseID:   r.ExpenseID,
			UserID:      r.FromUser,
			HouseholdID: householdID,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if !sh.PaidAt.Valid {
			unpaid = append(unpaid, r)
		}
	}
	for _, r := range unpaid {
		if _, err := payShare(ctx, q, householdID, r.ExpenseID, r.FromUser, r.ToUser, r.AmountCents); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

// Debt is From owing To, after the debts both ways are netted, Items add up to Cents
type Debt struct {
	From     pgtype.UUID
	FromName string
	To       pgtype.UUID
	ToName   string
	Cents    int64
	Items    []DebtItem
}

// DebtItem is an unpaid share (ExpenseID set) or a settlement, Row.ToUser paid either,
// Cents is positive when it adds to the debt it is listed under
type DebtItem struct {
	Row   dbqueries.ListDebtItemsRow
	Cents int64
}

// Credit is what is left of the payments that lower d, money confirmed but not yet spent on a share
func (d Debt) Credit() int64 {
	var credit int64
	for _, it := range d.Items {
		if !it.Row.ExpenseID.Valid && it.Cents < 0 {
			credit -= it.Cents
		}
	}
	return credit
}

// Covered says whether it is a share of d the [Debt.Credit] pays for in full
func (d Debt) Covered(it DebtItem) bool {
	return it.Row.ExpenseID.Valid && it.Cents > 0 && it.Cents <= d.Credit()
}

// flipped is d the other way round
func (d Debt) flipped() Debt {
	items := make([]DebtItem, len(d.Items))
	for i, it := range d.Items {
		items[i] = DebtItem{Row: it.Row, Cents: -it.Cents}
	}
	return Debt{From: d.To, FromName: d.ToName, To: d.From, ToName: d.FromName, Cents: -d.Cents, Items: items}
}

// netDebts nets the items of each pair of users and drops settled pairs
func netDebts(rows []dbqueries.ListDebtItemsRow) []Debt {
	byPair := map[[2]pgtype.UUID]*Debt{}
	var order [][2]pgtype.UUID
	for _, r := range rows {
		from, fromName, to, toName, cents := r.FromUser, r.FromName, r.ToUser, r.ToName, r.AmountCents
		// one key per pair whatever the direction
		if bytes.Compare(from.Bytes[:], to.Bytes[:]) > 0 {
			from, fromName, to, toName, cents = to, toName, from, fromName, -cents
		}
		key := [2]pgtype.UUID{from, to}
		d, ok := byPair[key]
		if !ok {
			d = &Debt{From: from, FromName: fromName, To: to, ToName: toName, Cents: 0, Items: nil}
			byPair[key] = d
			order = append(order, key)
		}
		d.Cents += cents
		d.Items = append(d.Items, DebtItem{Row: r, Cents: cents})
	}
	var debts []Debt
	for _, key := range order {
		d := *byPair[key]
		if d.Cents < 0 {
			d = d.flipped()
		}
		if d.Cents != 0 {
			debts = append(debts, d)
		}
	}
	slices.SortFunc(debts, func(a, b Debt) int {
		return cmp.Or(strings.Compare(a.FromName, b.FromName), strings.Compare(a.ToName, b.ToName))
	})
	return debts
}

// Debts returns who owes whom in the actor's household
func (s *DefaultService) Debts(ctx context.Context, actor Membership) ([]Debt, error) {
	rows, err := s.db.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{
		HouseholdID: actor.HouseholdID,
		UserID:      pgtype.UUID{},
	})
	return netDebts(rows), err
}

// checkSettled returns [g.ErrorBalanceNotZero] while userID owes or is owed anything,
// or has a pending settlement, confirming it later would leave a debt nobody can settle
func checkSettled(ctx context.Context, q *dbqueries.Queries, householdID, userID pgtype.UUID) error {
	rows, err := q.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	pending, err := q.ListPendingSettlements(ctx, dbqueries.ListPendingSettlementsParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	if len(netDebts(rows)) > 0 || len(pending) > 0 {
		return g.ErrorBalanceNotZero
	}
	return nil
}

// -----------------------------------------------------------------------------

// CreateSettlement records that the actor paid amount to toUserID, it counts once they confirm,
// [g.ErrorInvalidInput] for a bad amount or a receiver who is not another active member
func (s *DefaultService) CreateSettlement(ctx context.Context, actor Membership, toUserID pgtype.UUID, amount string) error {
	cents, err := ParseCents(amount)
	if err != nil || toUserID == actor.UserID {
		return g.ErrorInvalidInput
	}
	if _, err := s.db.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
		HouseholdID: actor.HouseholdID,
		UserID:      toUserID,
	}); err != nil {
		return noRowsAs(err, g.ErrorInvalidInput)
	}
	if err := s.db.InsertSettlement(ctx, dbqueries.InsertSettlementParams{
		HouseholdID: actor.HouseholdID,
		FromUser:    actor.UserID,
		ToUser:      toUserID,
		AmountCents: cents,
	}); err != nil {
		return err
	}
	s.notifyAsync(ctx, actor, []pgtype.UUID{toUserID}, func(ctx context.Context, _ pgtype.UUID, from, household string) push.Message {
		return push.Message{
			Title: utils.T(ctx, locales.LK_PushSettlementTitle, from),
			Body:  utils.T(ctx, locales.LK_PushSettlementBody, utils.FormatMoney(ctx, cents, actor.Currency), household),
			URL:   householdPath(actor.HouseholdID, "expenses"),
			Tag:   "settlement-" + actor.HouseholdID.String() + "-" + actor.UserID.String(),
		}
	})
	return nil
}

// ListPendingSettlements returns the unconfirmed settlements the actor sent or has to confirm
func (s *DefaultService) ListPendingSettlements(ctx context.Context, actor Membership) ([]dbqueries.ListPendingSettlementsRow, error) {
	return s.db.ListPendingSettlements(ctx, dbqueries.ListPendingSettlementsParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
}

// ConfirmSettlement returns [g.ErrorNotFound] unless the actor receives the pending settlement,
// a payment that clears the debt marks its shares paid, see [closePaidDebt]
func (s *DefaultService) ConfirmSettlement(ctx context.Context, actor Membership, settlementID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		st, err := q.ConfirmSettlement(ctx, dbqueries.ConfirmSettlementParams{
			ID:          settlementID,
			HouseholdID: actor.HouseholdID,
			ToUser:      actor.UserID,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		return closePaidDebt(ctx, q, actor.HouseholdID, st.FromUser, st.ToUser)
	})
}

// DeleteSettlement lets the sender cancel or the receiver reject a pending settlement
func (s *DefaultService) DeleteSettlement(ctx context.Context, actor Membership, settlementID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteSettlement(ctx, dbqueries.DeleteSettlementParams{
		ID:          settlementID,
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	}))
}
