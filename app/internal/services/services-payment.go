package services

import (
	"context"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
)

// Payment form fields a [FieldError] names
const (
	FieldPayeeName = "name"
	FieldIBAN      = "iban"
)

const MaxPayeeNameLength = 70

// PaymentDetails is where a user wants to be paid, both empty when not set
type PaymentDetails struct {
	Name string
	// the normalized IBAN, without spaces
	IBAN string
}

func (s *DefaultService) GetPaymentDetails(ctx context.Context, userID pgtype.UUID) (PaymentDetails, error) {
	row, err := s.db.GetPaymentDetails(ctx, userID)
	if err != nil {
		return PaymentDetails{}, noRowsAs(err, g.ErrorNotFound)
	}
	return payeeDetails(row.PayeeName, row.PayeeIban), nil
}

func payeeDetails(name, iban *string) PaymentDetails {
	var d PaymentDetails
	if name != nil {
		d.Name = *name
	}
	if iban != nil {
		d.IBAN = *iban
	}
	return d
}

// SetPaymentDetails saves the account holder and the IBAN, both empty clears them,
// a [FieldError] names what is wrong
func (s *DefaultService) SetPaymentDetails(ctx context.Context, userID pgtype.UUID, name, iban string) error {
	name, okName := trimToLength(name, 0, MaxPayeeNameLength)
	iban = models.NormalizeIBAN(iban)
	switch {
	case !okName || (iban != "" && name == ""):
		return FieldError{Field: FieldPayeeName}
	case iban != "" && !models.ValidIBAN(iban), iban == "" && name != "":
		return FieldError{Field: FieldIBAN}
	}
	p := dbqueries.SetPaymentDetailsParams{Name: nil, Iban: nil, UserID: userID}
	if iban != "" {
		p.Name, p.Iban = &name, &iban
	}
	return s.db.SetPaymentDetails(ctx, p)
}

// GetShowPaymentDetails is whether the actor's details are shown in their household
func (s *DefaultService) GetShowPaymentDetails(ctx context.Context, actor Membership) (bool, error) {
	show, err := s.db.GetShowPaymentDetails(ctx, dbqueries.GetShowPaymentDetailsParams{
		HouseholdID: actor.HouseholdID, UserID: actor.UserID,
	})
	return show, noRowsAs(err, g.ErrorNotFound)
}

func (s *DefaultService) SetShowPaymentDetails(ctx context.Context, actor Membership, show bool) error {
	return s.db.SetShowPaymentDetails(ctx, dbqueries.SetShowPaymentDetailsParams{
		Show: show, HouseholdID: actor.HouseholdID, UserID: actor.UserID,
	})
}

// ListPayeeDetails is the details of the active members who show them in the actor's household
func (s *DefaultService) ListPayeeDetails(ctx context.Context, actor Membership) (map[pgtype.UUID]PaymentDetails, error) {
	rows, err := s.db.ListPayeeDetails(ctx, actor.HouseholdID)
	if err != nil {
		return nil, err
	}
	out := make(map[pgtype.UUID]PaymentDetails, len(rows))
	for _, r := range rows {
		out[r.ID] = payeeDetails(r.PayeeName, r.PayeeIban)
	}
	return out, nil
}
