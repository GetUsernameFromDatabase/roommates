package services

import (
	"testing"

	g "roommates/internal/config"
)

func TestPaymentDetails(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)

	if err := s.SetPaymentDetails(ctx, memberM.UserID, "Mari Maasikas", "ee38 2200 2210 2014 5685"); err != nil {
		t.Fatal(err)
	}
	d, err := s.GetPaymentDetails(ctx, memberM.UserID)
	if err != nil || d != (PaymentDetails{Name: "Mari Maasikas", IBAN: "EE382200221020145685"}) {
		t.Fatalf("details %+v, %v", d, err)
	}
	for _, tc := range []struct{ name, iban, field string }{
		{"Mari", "EE382200221020145686", FieldIBAN},
		{"Mari", "", FieldIBAN},
		{"", "EE382200221020145685", FieldPayeeName},
	} {
		err := s.SetPaymentDetails(ctx, memberM.UserID, tc.name, tc.iban)
		wantErr(t, tc.iban, err, g.ErrorInvalidInput)
		if fe, ok := err.(FieldError); !ok || fe.Field != tc.field {
			t.Errorf("%q %q: error %v, want field %s", tc.name, tc.iban, err, tc.field)
		}
	}

	// shown only after the member switches it on
	if m, err := s.ListPayeeDetails(ctx, ownerM); err != nil || len(m) != 0 {
		t.Fatalf("before: %v, %v", m, err)
	}
	if err := s.SetShowPaymentDetails(ctx, memberM, true); err != nil {
		t.Fatal(err)
	}
	if show, err := s.GetShowPaymentDetails(ctx, memberM); err != nil || !show {
		t.Fatalf("show %v, %v", show, err)
	}
	if m, err := s.ListPayeeDetails(ctx, ownerM); err != nil || m[memberM.UserID].IBAN == "" || len(m) != 1 {
		t.Fatalf("after: %v, %v", m, err)
	}

	// both empty clears
	if err := s.SetPaymentDetails(ctx, memberM.UserID, " ", " "); err != nil {
		t.Fatal(err)
	}
	if d, _ = s.GetPaymentDetails(ctx, memberM.UserID); d != (PaymentDetails{}) {
		t.Fatalf("not cleared: %+v", d)
	}
	if m, _ := s.ListPayeeDetails(ctx, ownerM); len(m) != 0 {
		t.Fatalf("cleared but listed: %v", m)
	}
}
