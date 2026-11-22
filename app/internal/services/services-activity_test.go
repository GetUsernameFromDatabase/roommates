package services

import (
	"bytes"
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// a cursor goes through a ?before= link unchanged, to the microsecond postgres keeps
func TestCursorRoundTrip(t *testing.T) {
	id, _ := utils.ParseUUID("0193a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b")
	c := Cursor{At: pgtype.Timestamptz{Time: time.UnixMicro(1732371234567890), Valid: true, InfinityModifier: pgtype.Finite}, ID: id}
	got, err := ParseCursor(c.String())
	if err != nil || !got.At.Time.Equal(c.At.Time) || !got.At.Valid || got.ID != c.ID {
		t.Fatalf("%q -> %+v, %v", c.String(), got, err)
	}
	if got, err := ParseCursor(""); err != nil || got != (Cursor{}) {
		t.Fatalf("empty -> %+v, %v", got, err)
	}
	for _, bad := range []string{"x", "123", "123-", "-" + routes.ShortUUID(id), "12a-" + routes.ShortUUID(id)} {
		if _, err := ParseCursor(bad); err == nil {
			t.Errorf("%q parsed", bad)
		}
	}
}

func TestActivity(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	cents := int64(65000)

	log := func(actor Membership, e ActivityEntry) {
		t.Helper()
		if err := s.inTx(ctx, func(q *dbqueries.Queries) error { return logActivity(ctx, q, actor, e) }); err != nil {
			t.Fatal(err)
		}
	}
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		return logActivity(ctx, q, ownerM, ActivityEntry{Kind: "nope", SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: "", AmountCents: nil, Currency: ""})
	})
	if err == nil {
		t.Fatal("unknown kind was logged")
	}

	// round trip with names and amount
	log(ownerM, ActivityEntry{Kind: g.ActExpenseAdded, SubjectID: pgtype.UUID{}, TargetUserID: memberM.UserID, Title: "Rent", AmountCents: &cents, Currency: "EUR"})
	rows, err := s.ListActivity(ctx, ownerM, Cursor{})
	// newest first, after the member's join from newTestHousehold
	if err != nil || len(rows) != 2 || rows[1].Kind != g.ActMemberJoined {
		t.Fatalf("rows %+v, %v", rows, err)
	}
	r := rows[0]
	if r.Kind != g.ActExpenseAdded || r.Title != "Rent" || r.AmountCents == nil || *r.AmountCents != cents ||
		r.Currency == nil || *r.Currency != "EUR" || r.ActorName == nil || r.TargetName == nil || r.ExpenseExists {
		t.Fatalf("row %+v", r)
	}
	if rows, err = s.ListActivity(ctx, strangerM, Cursor{}); err != nil || len(rows) != 0 {
		t.Fatalf("stranger sees %+v, %v", rows, err)
	}

	// pagination: 60 more rows, a page is 50 and before continues after the last one;
	// rows 50 and 51 have the same time, as the rows of one transaction do
	at := time.Now()
	for i := range 60 {
		mins := i + 1
		if mins > 48 {
			mins--
		}
		if _, err := s.pool.Exec(ctx, `INSERT INTO activity (household_id, actor_id, grp, kind, title, created_at)
			VALUES ($1, $2, 'notes', $3, 'n', $4)`,
			ownerM.HouseholdID, ownerM.UserID, g.ActNoteEdited, at.Add(-time.Duration(mins)*time.Minute)); err != nil {
			t.Fatal(err)
		}
	}
	page1, err := s.ListActivity(ctx, ownerM, Cursor{})
	if err != nil || len(page1) != ActivityPageSize {
		t.Fatalf("page 1: %d rows, %v", len(page1), err)
	}
	last := page1[len(page1)-1]
	page2, err := s.ListActivity(ctx, ownerM, Cursor{At: last.CreatedAt, ID: last.ID})
	// 60 + the expense + the member's join
	if err != nil || len(page2) != 12 {
		t.Fatalf("page 2: %d rows, %v", len(page2), err)
	}
	if !page2[0].CreatedAt.Time.Equal(last.CreatedAt.Time) || page2[0].ID == last.ID {
		t.Fatalf("page 2 starts at %v, not the row that shares the time of page 1's last", page2[0].CreatedAt.Time)
	}

	// a private note's activity is for its author
	noteID, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Diary", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	log(ownerM, ActivityEntry{Kind: g.ActNoteCreated, SubjectID: noteID, TargetUserID: pgtype.UUID{}, Title: "Diary", AmountCents: nil, Currency: ""})
	has := func(actor Membership) bool {
		t.Helper()
		rows, err := s.ListActivity(ctx, actor, Cursor{})
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range rows {
			if r.Kind == g.ActNoteCreated {
				return r.NoteExists
			}
		}
		return false
	}
	if !has(ownerM) || has(memberM) {
		t.Fatal("private note activity must be shown to its author only")
	}

	// purge: shopping is kept 30 days, money a year
	old := time.Now().Add(-40 * 24 * time.Hour)
	for _, kind := range []string{g.ActItemBought, g.ActExpenseAdded} {
		if _, err := s.pool.Exec(ctx, `INSERT INTO activity (household_id, actor_id, grp, kind, title, created_at)
			VALUES ($1, $2, $3, $4, 'old', $5)`, ownerM.HouseholdID, ownerM.UserID, g.ActivityGroup[kind], kind, old); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.PurgeActivity(ctx); err != nil {
		t.Fatal(err)
	}
	var shopping, money int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE grp = 'shopping'), count(*) FILTER (WHERE grp = 'money' AND title = 'old')
		FROM activity WHERE household_id = $1`, ownerM.HouseholdID).Scan(&shopping, &money); err != nil {
		t.Fatal(err)
	}
	if shopping != 0 || money != 1 {
		t.Fatalf("after purge %d shopping, %d old money rows", shopping, money)
	}
}

// loggedKinds is the kinds of the household's activity the actor sees, with their titles
func loggedKinds(t *testing.T, s *DefaultService, actor Membership) map[string][]string {
	t.Helper()
	rows, err := s.ListActivity(context.Background(), actor, Cursor{})
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string][]string{}
	for _, r := range rows {
		kinds[r.Kind] = append(kinds[r.Kind], r.Title)
	}
	return kinds
}

func wantKinds(t *testing.T, s *DefaultService, actor Membership, want ...string) map[string][]string {
	t.Helper()
	kinds := loggedKinds(t, s, actor)
	for _, k := range want {
		if len(kinds[k]) == 0 {
			t.Errorf("no %q activity, have %v", k, kinds)
		}
	}
	return kinds
}

func TestActivityMoney(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	in := ExpenseInput{
		Title: "Rent", Amount: "10", Category: "rent", IncurredOn: "2026-10-01", Custom: false,
		Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}
	id, err := s.CreateExpense(ctx, ownerM, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Title = "Rent 2"
	if err := s.UpdateExpense(ctx, ownerM, id, in); err != nil {
		t.Fatal(err)
	}
	if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, false); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cash", ""); err != nil {
		t.Fatal(err)
	}
	pending, err := s.ListPendingSettlements(ctx, ownerM)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending %+v, %v", pending, err)
	}
	if err := s.ConfirmSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	// a cancel by the sender is not a rejection, the receiver's reject is
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cash", ""); err != nil {
		t.Fatal(err)
	}
	pending, _ = s.ListPendingSettlements(ctx, ownerM)
	if err := s.DeleteSettlement(ctx, memberM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	if kinds := loggedKinds(t, s, ownerM); len(kinds[g.ActSettlementDeclined]) != 0 {
		t.Fatal("a cancel was logged as a rejection")
	}
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cash", ""); err != nil {
		t.Fatal(err)
	}
	pending, _ = s.ListPendingSettlements(ctx, ownerM)
	if err := s.DeleteSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteExpense(ctx, ownerM, id); err != nil {
		t.Fatal(err)
	}
	wantKinds(t, s, ownerM, g.ActExpenseAdded, g.ActExpenseEdited, g.ActSharePaid, g.ActShareUnpaid,
		g.ActSettlementRequested, g.ActSettlementConfirmed, g.ActSettlementDeclined, g.ActExpenseDeleted)
}

func TestActivityShopping(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	listID, err := s.CreateShoppingList(ctx, ownerM, "Groceries", "groceries")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.RenameShoppingList(ctx, memberM, listID, "Food"); err != nil {
		t.Fatal(err)
	}
	if err := s.AddShoppingItem(ctx, memberM, listID, "milk", ""); err != nil {
		t.Fatal(err)
	}
	l, err := s.GetShoppingList(ctx, ownerM, listID)
	if err != nil || len(l.Items) != 1 {
		t.Fatalf("list %+v, %v", l, err)
	}
	// unchecking is not logged
	if err := s.SetShoppingItemChecked(ctx, memberM, listID, l.Items[0].ID, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, memberM, listID, l.Items[0].ID, false); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, memberM, listID, l.Items[0].ID, true); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearCheckedShoppingItems(ctx, ownerM, listID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteShoppingList(ctx, ownerM, listID); err != nil {
		t.Fatal(err)
	}
	kinds := wantKinds(t, s, ownerM, g.ActListCreated, g.ActListRenamed, g.ActItemBought, g.ActItemsCleared, g.ActListDeleted)
	if len(kinds[g.ActItemBought]) != 2 {
		t.Fatalf("bought logged %d times, want 2", len(kinds[g.ActItemBought]))
	}
}

func TestActivityNotes(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	noteID, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Wifi", BodyMD: "", Private: false})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateNote(ctx, memberM, noteID, NoteInput{Title: "Wi-Fi", BodyMD: "x", Private: false}); err != nil {
		t.Fatal(err)
	}
	// a private note leaves no trace, not even for its author
	private, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Diary", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateNote(ctx, ownerM, private, NoteInput{Title: "Diary 2", BodyMD: "", Private: true}); err != nil {
		t.Fatal(err)
	}
	body := []byte("hello")
	fileID, err := s.UploadFile(ctx, ownerM, FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}, EventID: pgtype.UUID{}},
		"hello.txt", bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteFile(ctx, ownerM, fileID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteNote(ctx, ownerM, noteID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteNote(ctx, ownerM, private); err != nil {
		t.Fatal(err)
	}
	kinds := wantKinds(t, s, ownerM, g.ActNoteCreated, g.ActNoteEdited, g.ActNoteDeleted, g.ActFileUploaded, g.ActFileRemoved)
	for _, titles := range kinds {
		for _, title := range titles {
			if strings.HasPrefix(title, "Diary") {
				t.Fatalf("private note logged as %q", title)
			}
		}
	}
}

func TestActivityMembers(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	if err := s.UpdateHousehold(ctx, ownerM, "Flat 2", "USD", "Europe/Tallinn", 9); err != nil {
		t.Fatal(err)
	}
	if err := s.PromoteMember(ctx, ownerM, memberM.UserID); err != nil {
		t.Fatal(err)
	}
	if err := s.StepDown(ctx, ownerM); err != nil {
		t.Fatal(err)
	}
	newOwner := memberM
	newOwner.Role = g.RoleOwner
	if err := s.TransferOwnership(ctx, newOwner, ownerM.UserID); err != nil {
		t.Fatal(err)
	}
	if err := s.Leave(ctx, memberM); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteHousehold(ctx, ownerM, "Flat 2"); err != nil {
		t.Fatal(err)
	}
	if err := s.RestoreHousehold(ctx, ownerM.UserID, ownerM.HouseholdID); err != nil {
		t.Fatal(err)
	}
	wantKinds(t, s, ownerM, g.ActMemberJoined, g.ActHouseholdRenamed, g.ActCurrencyChanged, g.ActMemberPromoted, g.ActSteppedDown,
		g.ActOwnershipTransferred, g.ActMemberLeft, g.ActHouseholdDeleted, g.ActHouseholdRestored)
}
