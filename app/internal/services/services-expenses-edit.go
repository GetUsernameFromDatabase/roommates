package services

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Why the amount and split of an expense cannot be edited, see [Expense.AmountLock]
const (
	// the expense is in another currency than the household now
	AmountLockCurrency = "currency"
	// someone sharing it has left, their settled balance cannot change
	AmountLockLeft = "left"
)

// AmountLock says why the amount and split of e cannot change while the household is in
// currency, "" when they can
func (e Expense) AmountLock(currency string) string {
	switch {
	case e.Expense.Currency != currency:
		return AmountLockCurrency
	case slices.ContainsFunc(e.Shares, func(sh dbqueries.ListExpenseSharesRow) bool { return !sh.Active }):
		return AmountLockLeft
	}
	return ""
}

// EditInput fills the expense form with e as it is, see [Expense.CopyInput]
func (e Expense) EditInput(order []pgtype.UUID) ExpenseInput {
	in := e.CopyInput(order)
	in.IncurredOn = e.Expense.IncurredOn.Time.Format(time.DateOnly)
	return in
}

// UpdateExpense lets the payer change an expense, [g.ErrorNotFound] for anyone else, errors of
// the input as [ExpenseInput.clean]; the amount and split change only without an [Expense.AmountLock],
// with one an amount other than the old is a [FieldError] and the split is left as it is
//
// a paid share that changes is opened again: what was paid outside the app becomes a confirmed
// payment from_edit, the payments spent on it are given back, and it is paid again only when
// all that covers the new amount; each balance then moves by exactly the change of the share
func (s *DefaultService) UpdateExpense(ctx context.Context, actor Membership, expenseID pgtype.UUID, in ExpenseInput) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// a member leaving or the currency changing waits for the edit
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		ex, err := q.GetExpenseToEdit(ctx, dbqueries.GetExpenseToEditParams{ID: expenseID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if ex.PayerID != actor.UserID {
			return g.ErrorNotFound
		}
		old, err := q.LockExpenseShares(ctx, dbqueries.LockExpenseSharesParams{ExpenseID: expenseID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return err
		}
		title, date, err := in.cleanDetails()
		if err != nil {
			return err
		}
		p := dbqueries.UpdateExpenseParams{
			ID:          expenseID,
			HouseholdID: actor.HouseholdID,
			Category:    in.Category,
			Title:       title,
			AmountCents: ex.AmountCents,
			IncurredOn:  date,
		}
		locked := ex.Currency != ex.HouseholdCurrency ||
			slices.ContainsFunc(old, func(sh dbqueries.LockExpenseSharesRow) bool { return !sh.Active })
		if locked {
			// the form disables the amount, so it is sent only by hand
			if cents, err := ParseCents(in.Amount); in.Amount != "" && (err != nil || cents != ex.AmountCents) {
				return FieldError{Field: FieldAmount}
			}
			if err := q.UpdateExpense(ctx, p); err != nil {
				return err
			}
			return logExpenseEdit(ctx, q, actor, p)
		}
		members, err := q.ListActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		cleaned, shares, err := in.clean(actor.UserID, members)
		if err != nil {
			return err
		}
		p.AmountCents = cleaned.AmountCents
		if err := q.UpdateExpense(ctx, p); err != nil {
			return err
		}
		if err := editShares(ctx, q, actor.HouseholdID, expenseID, actor.UserID, ex.Currency, old, shares); err != nil {
			return err
		}
		return logExpenseEdit(ctx, q, actor, p)
	})
}

func logExpenseEdit(ctx context.Context, q *dbqueries.Queries, actor Membership, p dbqueries.UpdateExpenseParams) error {
	cents := p.AmountCents
	return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActExpenseEdited, SubjectID: p.ID, TargetUserID: pgtype.UUID{}, Title: p.Title, AmountCents: &cents, Currency: actor.Currency})
}

// editShares turns the locked old shares of an expense into shares, see [DefaultService.UpdateExpense]
func editShares(ctx context.Context, q *dbqueries.Queries, householdID, expenseID, payerID pgtype.UUID, currency string,
	old []dbqueries.LockExpenseSharesRow, shares []expenseShare) error {
	added := map[pgtype.UUID]int64{}
	for _, sh := range shares {
		added[sh.userID] = sh.cents
	}
	for _, o := range old {
		cents, kept := added[o.UserID]
		delete(added, o.UserID)
		if kept && cents == o.AmountCents {
			continue
		}
		// the payer's own share is owed to nobody
		reopened := o.PaidAt.Valid && o.UserID != payerID
		if reopened {
			if outside := o.AmountCents - o.SpentCents; outside > 0 {
				if err := q.InsertEditSettlement(ctx, dbqueries.InsertEditSettlementParams{
					HouseholdID: householdID,
					FromUser:    o.UserID,
					ToUser:      payerID,
					AmountCents: outside,
					Currency:    currency,
					ExpenseID:   expenseID,
				}); err != nil {
					return err
				}
			}
			if err := q.DeleteSettlementShares(ctx, dbqueries.DeleteSettlementSharesParams{ExpenseID: expenseID, UserID: o.UserID}); err != nil {
				return err
			}
			if err := rowsAffected(q.SetExpenseSharePaid(ctx, dbqueries.SetExpenseSharePaidParams{
				Paid:        false,
				ExpenseID:   expenseID,
				UserID:      o.UserID,
				HouseholdID: householdID,
				PayerID:     payerID,
			})); err != nil {
				return err
			}
		}
		if !kept {
			if err := q.DeleteExpenseShare(ctx, dbqueries.DeleteExpenseShareParams{ExpenseID: expenseID, UserID: o.UserID, HouseholdID: householdID}); err != nil {
				return err
			}
			continue
		}
		if err := q.UpdateExpenseShare(ctx, dbqueries.UpdateExpenseShareParams{
			AmountCents: cents,
			ExpenseID:   expenseID,
			UserID:      o.UserID,
			HouseholdID: householdID,
		}); err != nil {
			return err
		}
		if reopened {
			if err := payIfCovered(ctx, q, householdID, expenseID, o.UserID, payerID, cents); err != nil {
				return err
			}
		}
	}
	for _, sh := range shares {
		if _, ok := added[sh.userID]; !ok {
			continue
		}
		if err := q.InsertExpenseShare(ctx, dbqueries.InsertExpenseShareParams{
			ExpenseID:   expenseID,
			UserID:      sh.userID,
			AmountCents: sh.cents,
		}); err != nil {
			return err
		}
	}
	return nil
}

// payIfCovered marks an unpaid share paid when userID's confirmed payments to payerID cover all
// of it, [payShare] would count an uncovered rest as paid outside the app
func payIfCovered(ctx context.Context, q *dbqueries.Queries, householdID, expenseID, userID, payerID pgtype.UUID, cents int64) error {
	if err := q.LockSettlementCredit(ctx, dbqueries.LockSettlementCreditParams{HouseholdID: householdID, FromUser: userID, ToUser: payerID}); err != nil {
		return err
	}
	credits, err := q.ListSettlementCredit(ctx, dbqueries.ListSettlementCreditParams{HouseholdID: householdID, FromUser: userID, ToUser: payerID})
	if err != nil {
		return err
	}
	var credit int64
	for _, c := range credits {
		credit += max(c.LeftCents, 0)
	}
	if credit < cents {
		return nil
	}
	_, err = payShare(ctx, q, householdID, expenseID, userID, payerID, cents)
	return err
}
