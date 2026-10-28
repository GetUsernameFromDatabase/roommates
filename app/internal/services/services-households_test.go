package services

import (
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
	return NewDefaultService(pool, &fakeMailer{sent: make(chan mailer.Message, 10)})
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

	householdID, err := s.CreateHousehold(ctx, owner, "Flat")
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

	otherID, err := s.CreateHousehold(ctx, stranger, "Other")
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

	householdID, err := s.CreateHousehold(ctx, owner, "  Flat 1 ")
	if err != nil {
		t.Fatal(err)
	}
	ownerM := mustMembership(t, s, householdID, owner)
	if ownerM.Role != g.RoleOwner {
		t.Fatalf("creator role is %q", ownerM.Role)
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

	householdA, err := s.CreateHousehold(ctx, ownerA, "A")
	if err != nil {
		t.Fatal(err)
	}
	householdB, err := s.CreateHousehold(ctx, ownerB, "B")
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
