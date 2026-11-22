package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"os"
	g "roommates/internal/config"
	"roommates/internal/db"
	"roommates/internal/db/dbqueries"
	"roommates/internal/mailer"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newTestService migrates the database named in TEST_POSTGRES_DB to the latest
// version, rows are left behind (users get random emails)
func newTestService(t *testing.T) *DefaultService {
	t.Helper()
	testDB := os.Getenv("TEST_POSTGRES_DB")
	if testDB == "" {
		t.Skip("TEST_POSTGRES_DB not set")
	}
	t.Setenv(g.EVK_DatabaseName.String(), testDB)

	pool, err := pgxpool.New(t.Context(), g.GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	db.LockTestDB(t, pool)
	db.MigrateToLatest(pool)
	return NewDefaultService(pool, &fakeMailer{sent: make(chan mailer.Message, 10)}, newFakeStore(), newFakePusher())
}

// fakeMailer collects emails, [DefaultService] sends them from a goroutine
type fakeMailer struct {
	sent chan mailer.Message
}

func (m *fakeMailer) Send(msg mailer.Message) error {
	m.sent <- msg
	return nil
}

// nextEmail waits for the next email sent by s
func nextEmail(t *testing.T, s *DefaultService) mailer.Message {
	t.Helper()
	select {
	case msg := <-s.mailer.(*fakeMailer).sent:
		return msg
	case <-time.After(5 * time.Second):
		t.Fatal("no email was sent")
		return mailer.Message{}
	}
}

func noEmail(t *testing.T, s *DefaultService) {
	t.Helper()
	select {
	case msg := <-s.mailer.(*fakeMailer).sent:
		t.Fatalf("unexpected email %q to %s", msg.Subject, msg.To)
	case <-time.After(100 * time.Millisecond):
	}
}

// newVerifiedTestUser is [newTestUser] with a confirmed email, invites need it
func newVerifiedTestUser(t *testing.T, s *DefaultService) pgtype.UUID {
	t.Helper()
	id := newTestUser(t, s)
	if err := s.db.SetEmailValidated(t.Context(), id); err != nil {
		t.Fatal(err)
	}
	return id
}

func newTestUser(t *testing.T, s *DefaultService) pgtype.UUID {
	t.Helper()
	id, err := s.db.InsertUser(t.Context(), dbqueries.InsertUserParams{
		Email:    strings.ToLower(rand.Text()) + "@example.com",
		Username: "test",
		Password: "not a hash",
		Language: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// newTestHousehold is an owner, a member who joined by invite,
// and a stranger who owns another household
func newTestHousehold(t *testing.T, s *DefaultService) (ownerM, memberM, strangerM Membership) {
	t.Helper()
	ctx := t.Context()
	owner, member, stranger := newVerifiedTestUser(t, s), newTestUser(t, s), newTestUser(t, s)

	householdID, err := s.CreateHousehold(ctx, owner, "Flat", "")
	if err != nil {
		t.Fatal(err)
	}
	ownerM = mustMembership(t, s, householdID, owner)
	if err := s.CreateInvite(ctx, ownerM); err != nil {
		t.Fatal(err)
	}
	overview, err := s.GetHouseholdOverview(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.JoinWithInvite(ctx, overview.Invites[0].Code, member); err != nil {
		t.Fatal(err)
	}
	memberM = mustMembership(t, s, householdID, member)

	otherID, err := s.CreateHousehold(ctx, stranger, "Other", "")
	if err != nil {
		t.Fatal(err)
	}
	strangerM = mustMembership(t, s, otherID, stranger)
	return ownerM, memberM, strangerM
}

func mustMembership(t *testing.T, s *DefaultService, householdID, userID pgtype.UUID) Membership {
	t.Helper()
	m, err := s.GetMembership(t.Context(), householdID, userID)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// lockHouseholdInTx is a household change in flight, holding the lock until tx ends,
// pid is its connection for [waitBlockedBy]
func lockHouseholdInTx(t *testing.T, s *DefaultService, householdID pgtype.UUID) (tx pgx.Tx, pid int32) {
	t.Helper()
	tx, err := s.pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { tx.Rollback(context.Background()) }) //nolint:errcheck // no-op after commit
	if err := s.db.WithTx(tx).LockHousehold(t.Context(), householdID); err != nil {
		t.Fatal(err)
	}
	if err := tx.QueryRow(t.Context(), "SELECT pg_backend_pid()").Scan(&pid); err != nil {
		t.Fatal(err)
	}
	return tx, pid
}

// waitBlockedBy waits until another connection waits for a lock pid holds
func waitBlockedBy(t *testing.T, s *DefaultService, pid int32) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for blocked := false; !blocked; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("nothing waited for the lock")
		}
		if err := s.pool.QueryRow(t.Context(),
			"SELECT EXISTS (SELECT 1 FROM pg_stat_activity WHERE $1::int = ANY(pg_blocking_pids(pid)))", pid,
		).Scan(&blocked); err != nil {
			t.Fatal(err)
		}
	}
}

// a join waits for a revoke in flight and then finds the invite gone
func TestJoinWaitsForRevoke(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, _, strangerM := newTestHousehold(t, s)
	overview, err := s.GetHouseholdOverview(ctx, ownerM)
	if err != nil || len(overview.Invites) == 0 {
		t.Fatalf("invites = %+v, %v", overview.Invites, err)
	}

	tx, pid := lockHouseholdInTx(t, s, ownerM.HouseholdID)
	joined := make(chan error, 1)
	go func() {
		_, err := s.JoinWithInvite(ctx, overview.Invites[0].Code, strangerM.UserID)
		joined <- err
	}()
	waitBlockedBy(t, s, pid)
	if err := s.db.WithTx(tx).DeleteInvitesOfHousehold(ctx, ownerM.HouseholdID); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	wantErr(t, "join with an invite revoked meanwhile", <-joined, g.ErrorInviteInvalid)
	_, err = s.GetMembership(ctx, ownerM.HouseholdID, strangerM.UserID)
	wantErr(t, "membership after the refused join", err, g.ErrorNotFound)
}

func wantErr(t *testing.T, step string, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("%s: got error %v, want %v", step, got, want)
	}
}

func TestHouseholdLifecycle(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	owner, member, stranger := newVerifiedTestUser(t, s), newVerifiedTestUser(t, s), newTestUser(t, s)

	householdID, err := s.CreateHousehold(ctx, owner, "  Flat 1 ", "Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	ownerM := mustMembership(t, s, householdID, owner)
	if ownerM.Role != g.RoleOwner {
		t.Fatalf("creator role is %q", ownerM.Role)
	}
	// the creator's browser zone is kept, an unknown one leaves the default
	if h, err := s.GetHouseholdSettings(ctx, ownerM); err != nil || h.TimeZone != "Asia/Tokyo" {
		t.Fatalf("time zone %q %v", h.TimeZone, err)
	}
	otherID, err := s.CreateHousehold(ctx, owner, "Flat 2", "Mars/Olympus")
	if err != nil {
		t.Fatal(err)
	}
	if h, err := s.GetHouseholdSettings(ctx, mustMembership(t, s, otherID, owner)); err != nil || h.TimeZone != "Europe/Tallinn" {
		t.Fatalf("default time zone %q %v", h.TimeZone, err)
	}

	// isolation: a non-member cannot see the household
	_, err = s.GetMembership(ctx, householdID, stranger)
	wantErr(t, "stranger membership", err, g.ErrorNotFound)

	// invites
	if err := s.CreateInvite(ctx, ownerM); err != nil {
		t.Fatal(err)
	}
	overview, err := s.GetHouseholdOverview(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Household.Name != "Flat 1" || len(overview.Invites) != 1 {
		t.Fatalf("unexpected overview %+v", overview)
	}
	code := overview.Invites[0].Code

	_, err = s.JoinWithInvite(ctx, "WRONGCODE0", member)
	wantErr(t, "join with wrong code", err, g.ErrorInviteInvalid)

	// lowercase with spaces, as someone might type it
	joinedID, err := s.JoinWithInvite(ctx, "  "+strings.ToLower(code)+" ", member)
	if err != nil || joinedID != householdID {
		t.Fatalf("join: %v %v", joinedID, err)
	}
	memberM := mustMembership(t, s, householdID, member)

	// joining again keeps the owner an owner
	if _, err := s.JoinWithInvite(ctx, code, owner); err != nil {
		t.Fatal(err)
	}
	if mustMembership(t, s, householdID, owner).Role != g.RoleOwner {
		t.Fatal("owner was demoted by joining again")
	}

	// members cannot manage
	wantErr(t, "member creates invite", s.CreateInvite(ctx, memberM), g.ErrorForbidden)
	wantErr(t, "member removes owner", s.RemoveMember(ctx, memberM, owner), g.ErrorForbidden)
	overview, err = s.GetHouseholdOverview(ctx, memberM)
	if err != nil || overview.Invites != nil {
		t.Fatalf("member sees invites: %v %v", overview.Invites, err)
	}

	// owners cannot be removed, strangers are not found
	wantErr(t, "owner removes self", s.RemoveMember(ctx, ownerM, owner), g.ErrorForbidden)
	wantErr(t, "owner removes stranger", s.RemoveMember(ctx, ownerM, stranger), g.ErrorNotFound)

	// the only owner cannot leave while others remain
	wantErr(t, "last owner leaves", s.Leave(ctx, ownerM), g.ErrorLastOwner)

	// after promoting, the first owner can leave
	if err := s.PromoteMember(ctx, ownerM, member); err != nil {
		t.Fatal(err)
	}
	if err := s.Leave(ctx, ownerM); err != nil {
		t.Fatal(err)
	}
	_, err = s.GetMembership(ctx, householdID, owner)
	wantErr(t, "membership after leaving", err, g.ErrorNotFound)

	// a former member can rejoin, as a member
	if _, err := s.JoinWithInvite(ctx, code, owner); err != nil {
		t.Fatal(err)
	}
	if mustMembership(t, s, householdID, owner).Role != g.RoleMember {
		t.Fatal("rejoined former owner should be a member")
	}

	// removing revokes the invites, so the removed member cannot rejoin with one
	newOwnerM := mustMembership(t, s, householdID, member)
	if err := s.RemoveMember(ctx, newOwnerM, owner); err != nil {
		t.Fatal(err)
	}
	_, err = s.JoinWithInvite(ctx, code, owner)
	wantErr(t, "rejoin after removal", err, g.ErrorInviteInvalid)

	// the last member leaving revokes the invites, nobody could manage a joiner
	if err := s.CreateInvite(ctx, newOwnerM); err != nil {
		t.Fatal(err)
	}
	overview, err = s.GetHouseholdOverview(ctx, newOwnerM)
	if err != nil {
		t.Fatal(err)
	}
	lastCode := overview.Invites[0].Code
	if err := s.Leave(ctx, newOwnerM); err != nil {
		t.Fatal(err)
	}
	_, err = s.JoinWithInvite(ctx, lastCode, stranger)
	wantErr(t, "join empty household", err, g.ErrorInviteInvalid)
}

func TestDeleteInviteOnlyInOwnHousehold(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerA, ownerB := newVerifiedTestUser(t, s), newVerifiedTestUser(t, s)

	householdA, err := s.CreateHousehold(ctx, ownerA, "A", "")
	if err != nil {
		t.Fatal(err)
	}
	householdB, err := s.CreateHousehold(ctx, ownerB, "B", "")
	if err != nil {
		t.Fatal(err)
	}
	mA := mustMembership(t, s, householdA, ownerA)
	mB := mustMembership(t, s, householdB, ownerB)
	if err := s.CreateInvite(ctx, mA); err != nil {
		t.Fatal(err)
	}
	overview, err := s.GetHouseholdOverview(ctx, mA)
	if err != nil {
		t.Fatal(err)
	}

	// owner of B tries to delete A's invite by id
	if err := s.DeleteInvite(ctx, mB, overview.Invites[0].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetInvite(ctx, overview.Invites[0].Code); err != nil {
		t.Fatalf("invite of A was deleted by owner of B: %v", err)
	}
}

func TestNewInviteCode(t *testing.T) {
	for range 100 {
		code := newInviteCode()
		if len(code) != inviteCodeLength || NormalizeInviteCode(code) != code {
			t.Fatalf("bad code %q", code)
		}
	}
}

func TestNormalizeInviteCode(t *testing.T) {
	cases := map[string]string{
		"PZMPRD1KKT":   "PZMPRD1KKT",
		" pzmprd1kkt ": "PZMPRD1KKT",
		"https://localhost:8090/invite?code=PZMPRD1KKT": "PZMPRD1KKT",
		" /invite?code=pzmprd1kkt":                      "PZMPRD1KKT",
		"":                                              "",
	}
	for in, want := range cases {
		if got := NormalizeInviteCode(in); got != want {
			t.Errorf("NormalizeInviteCode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUpdateHousehold(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)

	const tz = "Europe/Tallinn"
	wantErr(t, "member updates", s.UpdateHousehold(ctx, memberM, "Mine", "EUR", tz, 9), g.ErrorForbidden)
	wantErr(t, "empty name", s.UpdateHousehold(ctx, ownerM, "  ", "EUR", tz, 9), g.ErrorInvalidInput)
	wantErr(t, "long name", s.UpdateHousehold(ctx, ownerM, strings.Repeat("x", MaxHouseholdNameLength+1), "EUR", tz, 9), g.ErrorInvalidInput)
	wantErr(t, "unknown currency", s.UpdateHousehold(ctx, ownerM, "Flat", "XYZ", tz, 9), g.ErrorInvalidInput)
	for _, bad := range []string{"", "Mars/Olympus", "posix/Europe/Tallinn", "EST5EDT"} {
		err := s.UpdateHousehold(ctx, ownerM, "Flat", "EUR", bad, 9)
		var fe FieldError
		if !errors.As(err, &fe) || fe.Field != FieldTimeZone {
			t.Errorf("time zone %q: %v is not a time zone FieldError", bad, err)
		}
	}
	for _, bad := range []int{-1, 24} {
		err := s.UpdateHousehold(ctx, ownerM, "Flat", "EUR", tz, bad)
		var fe FieldError
		if !errors.As(err, &fe) || fe.Field != FieldReminderHour {
			t.Errorf("hour %d: %v is not a reminder hour FieldError", bad, err)
		}
	}

	// no expenses yet, so the currency can change
	if err := s.UpdateHousehold(ctx, ownerM, " Home ", "USD", "America/New_York", 20); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetHouseholdSettings(ctx, memberM)
	if err != nil || got.Name != "Home" || got.Currency != "USD" || !got.Settled || got.TimeZone != "America/New_York" || got.ReminderHour != 20 {
		t.Fatalf("settings %+v %v", got, err)
	}

	if _, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
		Title: "Rent", Amount: "10", Category: "rent", IncurredOn: "2026-10-01",
		Custom: false, Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil,
		ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}); err != nil {
		t.Fatal(err)
	}
	err = s.UpdateHousehold(ctx, ownerM, "Home", "EUR", tz, 9)
	wantErr(t, "currency with a debt", err, g.ErrorInvalidInput)
	var fe FieldError
	if !errors.As(err, &fe) || fe.Field != FieldCurrency {
		t.Fatalf("error %v is not a currency FieldError", err)
	}
	if got, err = s.GetHouseholdSettings(ctx, ownerM); err != nil || got.Settled || got.Currency != "USD" {
		t.Fatalf("settings with a debt %+v %v", got, err)
	}
	// the same currency is no change, a rename is fine
	if err := s.UpdateHousehold(ctx, ownerM, "Flat", "USD", tz, 9); err != nil {
		t.Fatal(err)
	}

	// the member's same expense nets the pair to zero with both shares open, the change closes them
	if _, err := s.CreateExpense(ctx, memberM, ExpenseInput{
		Title: "Power", Amount: "10", Category: "utilities", IncurredOn: "2026-10-02",
		Custom: false, Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil,
		ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateHousehold(ctx, ownerM, "Flat", "EUR", tz, 9); err != nil {
		t.Fatalf("currency when even: %v", err)
	}
	rows, err := s.db.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{HouseholdID: ownerM.HouseholdID, UserID: pgtype.UUID{}})
	if err != nil || len(rows) != 0 {
		t.Fatalf("open items after the change %v %v", rows, err)
	}
}

func TestTransferOwnershipAndStepDown(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)

	wantErr(t, "member transfers", s.TransferOwnership(ctx, memberM, ownerM.UserID), g.ErrorForbidden)
	wantErr(t, "transfer to a stranger", s.TransferOwnership(ctx, ownerM, strangerM.UserID), g.ErrorNotFound)
	wantErr(t, "only owner steps down", s.StepDown(ctx, ownerM), g.ErrorLastOwner)

	if err := s.TransferOwnership(ctx, ownerM, memberM.UserID); err != nil {
		t.Fatal(err)
	}
	formerOwner := mustMembership(t, s, ownerM.HouseholdID, ownerM.UserID)
	newOwner := mustMembership(t, s, ownerM.HouseholdID, memberM.UserID)
	if formerOwner.Role != g.RoleMember || newOwner.Role != g.RoleOwner {
		t.Fatalf("roles after transfer: former %s, new %s", formerOwner.Role, newOwner.Role)
	}
	// a request of the owner's that ran at the same time (a double submit) loaded the role before the transfer
	wantErr(t, "stale owner promotes itself back", s.PromoteMember(ctx, ownerM, ownerM.UserID), g.ErrorForbidden)
	wantErr(t, "member steps down", s.StepDown(ctx, formerOwner), g.ErrorForbidden)
	wantErr(t, "new only owner steps down", s.StepDown(ctx, newOwner), g.ErrorLastOwner)

	if err := s.PromoteMember(ctx, newOwner, formerOwner.UserID); err != nil {
		t.Fatal(err)
	}
	if err := s.StepDown(ctx, newOwner); err != nil {
		t.Fatalf("step down with two owners: %v", err)
	}
	if m := mustMembership(t, s, ownerM.HouseholdID, memberM.UserID); m.Role != g.RoleMember {
		t.Fatalf("role after stepping down %s", m.Role)
	}
}

func TestDeleteAndRestoreHousehold(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	if err := s.CreateInvite(ctx, ownerM); err != nil {
		t.Fatal(err)
	}
	// an open debt does not stop the deletion
	if _, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
		Title: "Rent", Amount: "10", Category: "rent", IncurredOn: "2026-10-01",
		Custom: false, Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil,
		ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}); err != nil {
		t.Fatal(err)
	}
	if settings, err := s.GetHouseholdSettings(ctx, ownerM); err != nil || len(settings.Debts) != 1 {
		t.Fatalf("settings debts %+v, %v", settings.Debts, err)
	}

	wantErr(t, "member deletes", s.DeleteHousehold(ctx, memberM, "Flat"), g.ErrorForbidden)
	err := s.DeleteHousehold(ctx, ownerM, "flat")
	var fe FieldError
	if !errors.As(err, &fe) || fe.Field != FieldName {
		t.Fatalf("wrong name: error %v is not a name FieldError", err)
	}
	if m := mustMembership(t, s, ownerM.HouseholdID, ownerM.UserID); m.DeletedAt.Valid {
		t.Fatal("deleted after a wrong name")
	}

	if err := s.DeleteHousehold(ctx, ownerM, " Flat "); err != nil {
		t.Fatal(err)
	}
	if m := mustMembership(t, s, ownerM.HouseholdID, memberM.UserID); !m.DeletedAt.Valid {
		t.Fatal("deleted_at not set")
	}
	invites, err := s.db.ListActiveInvites(ctx, ownerM.HouseholdID)
	if err != nil || len(invites) != 0 {
		t.Fatalf("invites after delete %v, %v", invites, err)
	}

	wantErr(t, "member restores", s.RestoreHousehold(ctx, memberM.UserID, ownerM.HouseholdID), g.ErrorForbidden)
	wantErr(t, "stranger restores", s.RestoreHousehold(ctx, strangerM.UserID, ownerM.HouseholdID), g.ErrorNotFound)
	if err := s.RestoreHousehold(ctx, ownerM.UserID, ownerM.HouseholdID); err != nil {
		t.Fatal(err)
	}
	if m := mustMembership(t, s, ownerM.HouseholdID, memberM.UserID); m.DeletedAt.Valid {
		t.Fatal("deleted_at still set after restore")
	}
}

func TestPurgeDeletedHouseholds(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	oldM, _, _ := newTestHousehold(t, s)
	recentM, _, _ := newTestHousehold(t, s)
	none := FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}, EventID: pgtype.UUID{}}
	body := []byte("\x89PNG\r\n\x1a\n rest of the image")
	fileID, err := s.UploadFile(ctx, oldM, none, "photo.png", bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	var key string
	if err := s.pool.QueryRow(ctx, "SELECT object_key FROM files WHERE id = $1", fileID).Scan(&key); err != nil {
		t.Fatal(err)
	}
	for _, m := range []Membership{oldM, recentM} {
		if err := s.DeleteHousehold(ctx, m, "Flat"); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.pool.Exec(ctx, "UPDATE households SET deleted_at = $2 WHERE id = $1",
		oldM.HouseholdID, time.Now().Add(-g.HouseholdPurgeAfter-time.Hour)); err != nil {
		t.Fatal(err)
	}
	// rows two cascade steps down that name a member (a share, a rota) must not block the purge
	if _, err := s.pool.Exec(ctx, `
		WITH e AS (
			INSERT INTO expenses (household_id, payer_id, category, title, amount_cents, currency, incurred_on)
			VALUES ($1, $2, 'other', 'Rent', 100, 'EUR', CURRENT_DATE) RETURNING id
		), sh AS (
			INSERT INTO expense_shares (expense_id, household_id, user_id, amount_cents) SELECT id, $1, $2, 100 FROM e
		), c AS (
			INSERT INTO chores (household_id, title, rotation, weekdays, next_due, assignee_id, created_by)
			VALUES ($1, 'Bins', 'fixed', 1, CURRENT_DATE, $2, $2) RETURNING id
		)
		INSERT INTO chore_members (chore_id, household_id, user_id, position) SELECT id, $1, $2, 0 FROM c`,
		oldM.HouseholdID, oldM.UserID); err != nil {
		t.Fatal(err)
	}

	if err := s.PurgeDeletedHouseholds(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.GetHousehold(ctx, oldM.HouseholdID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("old household after purge: %v", err)
	}
	if _, err := s.db.GetHousehold(ctx, recentM.HouseholdID); err != nil {
		t.Fatalf("recently deleted household purged: %v", err)
	}
	var queued int
	if err := s.pool.QueryRow(ctx, "SELECT count(*) FROM file_purges WHERE object_key = $1", key).Scan(&queued); err != nil || queued != 1 {
		t.Fatalf("queued %d, %v", queued, err)
	}
}
