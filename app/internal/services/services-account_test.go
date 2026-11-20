package services

import (
	"crypto/rand"
	"errors"
	"net/http"
	"net/http/httptest"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/models"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// setTestPassword gives a test user a real password hash
func setTestPassword(t *testing.T, s *DefaultService, userID pgtype.UUID, password string) {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.db.UpdateUserPassword(t.Context(), dbqueries.UpdateUserPasswordParams{ID: userID, Password: string(hashed)}); err != nil {
		t.Fatal(err)
	}
}

func addTestSession(t *testing.T, s *DefaultService, userID pgtype.UUID) {
	t.Helper()
	if _, err := s.pool.Exec(t.Context(), "INSERT INTO http_sessions (key, data, user_id) VALUES ($1, '', $2)",
		"key-"+rand.Text(), userID); err != nil {
		t.Fatal(err)
	}
}

func newTestEmail() string {
	return strings.ToLower(rand.Text()) + "@example.com"
}

func TestChangePassword(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	userID, otherID := newTestUser(t, s), newTestUser(t, s)
	setTestPassword(t, s, userID, "Old-pass1!")
	addTestSession(t, s, userID)
	addTestSession(t, s, userID)
	addTestSession(t, s, otherID)

	wantErr(t, "wrong current", s.ChangePassword(ctx, userID, "wrong", "New-pass1!"), g.ErrorInvalidCredential)
	if n := sessionsOf(t, s, userID); n != 2 {
		t.Fatalf("a refused change should sign nobody out, %d sessions left", n)
	}

	if err := s.ChangePassword(ctx, userID, "Old-pass1!", "New-pass1!"); err != nil {
		t.Fatal(err)
	}
	if err := s.checkPassword(ctx, userID, "New-pass1!"); err != nil {
		t.Fatalf("new password does not work: %v", err)
	}
	// the handler signs the current session in again
	if n := sessionsOf(t, s, userID); n != 0 {
		t.Fatalf("a change should sign the user out everywhere, %d sessions left", n)
	}
	if n := sessionsOf(t, s, otherID); n != 1 {
		t.Fatalf("other users should stay signed in, %d sessions left", n)
	}
}

func TestChangeEmail(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	userID, otherID := newTestUser(t, s), newTestUser(t, s)
	setTestPassword(t, s, userID, "Pass-word1!")
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	other, err := s.GetUser(ctx, otherID)
	if err != nil {
		t.Fatal(err)
	}
	newEmail := newTestEmail()

	wantErr(t, "wrong password", s.RequestEmailChange(ctx, userID, "wrong", newEmail), g.ErrorInvalidCredential)
	wantErr(t, "same email", s.RequestEmailChange(ctx, userID, "Pass-word1!", strings.ToUpper(user.Email)), ErrorEmailUnchanged)
	wantErr(t, "taken email", s.RequestEmailChange(ctx, userID, "Pass-word1!", other.Email), g.ErrorAccountAlreadyExists)
	noEmail(t, s)

	// a reset token is not an email change token
	if err := s.RequestPasswordReset(ctx, user.Email); err != nil {
		t.Fatal(err)
	}
	resetToken := tokenOf(t, nextEmail(t, s).Body)
	wantErr(t, "reset token used for change", s.ConfirmEmailChange(ctx, resetToken), g.ErrorTokenInvalid)

	if err := s.RequestEmailChange(ctx, userID, "Pass-word1!", " "+strings.ToUpper(newEmail)); err != nil {
		t.Fatal(err)
	}
	msg := nextEmail(t, s)
	if msg.To != newEmail || !strings.Contains(msg.Body, g.RouteConfirmEmail+"?token=") {
		t.Fatalf("unexpected confirmation email %+v", msg)
	}
	token := tokenOf(t, msg.Body)
	if u, _ := s.GetUser(ctx, userID); u.Email != user.Email {
		t.Fatal("the email should only change once the link is opened")
	}

	if err := s.ConfirmEmailChange(ctx, token); err != nil {
		t.Fatal(err)
	}
	changed, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if changed.Email != newEmail || !changed.EmailValidated {
		t.Fatalf("email not changed and confirmed: %+v", changed)
	}
	notice := nextEmail(t, s)
	if notice.To != user.Email || !strings.Contains(notice.Body, newEmail) {
		t.Fatalf("the old address should be told the new one, got %+v", notice)
	}
	wantErr(t, "token reused", s.ConfirmEmailChange(ctx, token), g.ErrorTokenInvalid)
	// the reset link went to the old address
	wantErr(t, "reset after change", s.CheckPasswordResetToken(ctx, resetToken), g.ErrorTokenInvalid)
}

// countRows runs a SELECT count(*) query
func countRows(t *testing.T, s *DefaultService, query string, args ...any) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(t.Context(), query, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func isDeleted(t *testing.T, s *DefaultService, userID pgtype.UUID) bool {
	t.Helper()
	return countRows(t, s, "SELECT count(*) FROM users WHERE id = $1 AND deleted_at IS NOT NULL", userID) == 1
}

func wantBlocked(t *testing.T, step string, err error, lastOwner, unsettled []string) {
	t.Helper()
	var blocked AccountDeleteBlocked
	if !errors.As(err, &blocked) || !slices.Equal(blocked.LastOwner, lastOwner) || !slices.Equal(blocked.Unsettled, unsettled) {
		t.Fatalf("%s: got %#v, want last owner %v and unsettled %v", step, err, lastOwner, unsettled)
	}
}

func TestDeleteAccountBlocked(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	ownerM, memberM, _ := newTestHousehold(t, s)
	setTestPassword(t, s, ownerM.UserID, "Pass-word1!")
	setTestPassword(t, s, memberM.UserID, "Pass-word1!")

	wantErr(t, "wrong password", s.DeleteAccount(ctx, memberM.UserID, "wrong"), g.ErrorInvalidCredential)
	wantBlocked(t, "last owner", s.DeleteAccount(ctx, ownerM.UserID, "Pass-word1!"), []string{"Flat"}, nil)

	expense, _, _, _ := settlementHelpers(t, s, ownerM, memberM)
	expense(ownerM, memberM, "10.00")
	wantBlocked(t, "open balance", s.DeleteAccount(ctx, memberM.UserID, "Pass-word1!"), nil, []string{"Flat"})
	// one try names every reason, the owner is owed the 10.00 too
	wantBlocked(t, "last owner and open balance", s.DeleteAccount(ctx, ownerM.UserID, "Pass-word1!"), []string{"Flat"}, []string{"Flat"})

	// a settlement waiting to be confirmed blocks without any debt
	owner2M, member2M, _ := newTestHousehold(t, s)
	setTestPassword(t, s, member2M.UserID, "Pass-word1!")
	if err := s.CreateSettlement(ctx, member2M, owner2M.UserID, "5.00", "transfer", ""); err != nil {
		t.Fatal(err)
	}
	wantBlocked(t, "pending settlement", s.DeleteAccount(ctx, member2M.UserID, "Pass-word1!"), nil, []string{"Flat"})

	// a blocked delete changes nothing
	for _, m := range []Membership{ownerM, memberM, member2M} {
		mustMembership(t, s, m.HouseholdID, m.UserID)
		if isDeleted(t, s, m.UserID) {
			t.Fatal("a blocked delete marked the account deleted")
		}
	}
	noEmail(t, s)
}

func TestDeleteAndRestoreAccount(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	ownerM, memberM, _ := newTestHousehold(t, s)
	userID := memberM.UserID
	setTestPassword(t, s, userID, "Pass-word1!")
	addTestSession(t, s, userID)
	subscribeTestDevice(t, s, userID, "en")
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	mustCreateChore(t, s, ownerM, everyDayChore("Bins", g.RotationFixed, userID, ownerM.UserID, userID))
	in := eventInput("Party", time.Now())
	in.Members = []pgtype.UUID{ownerM.UserID, userID}
	mustCreateEvent(t, s, ownerM, in)
	soloID, err := s.CreateHousehold(ctx, userID, "Solo", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteAccount(ctx, userID, "Pass-word1!"); err != nil {
		t.Fatal(err)
	}
	if !isDeleted(t, s, userID) {
		t.Fatal("account not marked deleted")
	}
	// the shared household is left at once
	_, err = s.GetMembership(ctx, memberM.HouseholdID, userID)
	wantErr(t, "membership of the shared household", err, g.ErrorNotFound)
	for _, table := range []string{"event_members", "chore_members"} {
		if n := countRows(t, s, "SELECT count(*) FROM "+table+" WHERE household_id = $1 AND user_id = $2", memberM.HouseholdID, userID); n != 0 {
			t.Fatalf("%d %s rows left", n, table)
		}
	}
	// the solo household is deleted but can still be restored
	if solo := mustMembership(t, s, soloID, userID); !solo.DeletedAt.Valid {
		t.Fatal("solo household not deleted")
	}
	if n := sessionsOf(t, s, userID); n != 0 {
		t.Fatalf("%d sessions left", n)
	}
	if n := countRows(t, s, "SELECT count(*) FROM push_subscriptions WHERE user_id = $1", userID); n != 0 {
		t.Fatalf("%d push subscriptions left", n)
	}
	if msg := nextEmail(t, s); msg.To != user.Email {
		t.Fatalf("deleted email went to %s", msg.To)
	}

	if err := s.RestoreAccount(ctx, userID); err != nil {
		t.Fatal(err)
	}
	if isDeleted(t, s, userID) {
		t.Fatal("restore left the account deleted")
	}
}

func TestEraseDeletedAccounts(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	ownerM, memberM, _ := newTestHousehold(t, s)
	userID := memberM.UserID
	setTestPassword(t, s, userID, "Pass-word1!")
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	private, err := s.CreateNote(ctx, memberM, NoteInput{Title: "Diary", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	shared, err := s.CreateNote(ctx, memberM, NoteInput{Title: "Wifi", BodyMD: "", Private: false})
	if err != nil {
		t.Fatal(err)
	}
	expense, _, _, _ := settlementHelpers(t, s, ownerM, memberM)
	expenseID := expense(memberM, ownerM, "10.00")
	avatarKey := "avatars/" + rand.Text()
	if err := s.db.SetUserAvatar(ctx, dbqueries.SetUserAvatarParams{ID: userID, AvatarKey: &avatarKey}); err != nil {
		t.Fatal(err)
	}
	addTestSession(t, s, userID)
	if _, err := s.pool.Exec(ctx, `UPDATE users SET payee_name = 'Me', payee_iban = 'EE00', language = 'en',
		email_validated = TRUE, deleted_at = now() - interval '31 days' WHERE id = $1`, userID); err != nil {
		t.Fatal(err)
	}
	// still within the grace period
	recentID := newTestUser(t, s)
	if _, err := s.pool.Exec(ctx, "UPDATE users SET deleted_at = now() - interval '29 days' WHERE id = $1", recentID); err != nil {
		t.Fatal(err)
	}

	if err := s.EraseDeletedAccounts(ctx); err != nil {
		t.Fatal(err)
	}
	var username, email, password string
	var avatar, payeeName, payeeIBAN, language *string
	var validated bool
	if err := s.pool.QueryRow(ctx, `SELECT username, email, password, avatar_key, payee_name, payee_iban, language, email_validated
		FROM users WHERE id = $1`, userID).Scan(&username, &email, &password, &avatar, &payeeName, &payeeIBAN, &language, &validated); err != nil {
		t.Fatal(err)
	}
	if username != "" || email != userID.String()+"@erased.invalid" || password != "" || avatar != nil ||
		payeeName != nil || payeeIBAN != nil || language != nil || validated {
		t.Fatalf("not wiped: %q %q %q %v %v %v %v %v", username, email, password, avatar, payeeName, payeeIBAN, language, validated)
	}
	if recent, err := s.GetUser(ctx, recentID); err != nil || recent.Username != "test" {
		t.Fatalf("account within the grace period erased: %+v, %v", recent, err)
	}
	_, err = s.GetMembership(ctx, memberM.HouseholdID, userID)
	wantErr(t, "membership after erase", err, g.ErrorNotFound)
	if n := countRows(t, s, "SELECT count(*) FROM notes WHERE id = $1", private); n != 0 {
		t.Fatal("private note kept")
	}
	if n := countRows(t, s, "SELECT count(*) FROM notes WHERE id = $1", shared); n != 1 {
		t.Fatal("household note removed")
	}
	if _, err := s.GetExpense(ctx, ownerM, expenseID); err != nil {
		t.Fatalf("expense removed: %v", err)
	}
	if n := countRows(t, s, "SELECT count(*) FROM file_purges WHERE object_key = $1", avatarKey); n != 1 {
		t.Fatal("avatar purge not queued")
	}
	if n := sessionsOf(t, s, userID); n != 0 {
		t.Fatalf("%d sessions left", n)
	}
	if msg := nextEmail(t, s); msg.To != user.Email {
		t.Fatalf("erased email went to %s", msg.To)
	}

	// an erased user cannot sign in, by the old email or the new one
	gctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	gctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	for _, e := range []string{user.Email, email} {
		for _, p := range []string{"Pass-word1!", ""} {
			_, err := s.ShouldUserBeSignedIn(gctx, models.LoginRequest{Email: e, Password: p})
			wantErr(t, "sign in as erased "+e, err, g.ErrorInvalidCredential)
		}
	}
	// a sign-in whose password check beat the erase must not revive the row
	wantErr(t, "restore erased", s.RestoreAccount(ctx, userID), g.ErrorInvalidCredential)
	if !isDeleted(t, s, userID) {
		t.Fatal("restore revived an erased account")
	}

	// erased once
	if err := s.EraseDeletedAccounts(ctx); err != nil {
		t.Fatal(err)
	}
	noEmail(t, s)
}

func TestEmailChangeTakenAtConfirm(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	userID, otherID := newTestUser(t, s), newTestUser(t, s)
	setTestPassword(t, s, userID, "Pass-word1!")
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	newEmail := newTestEmail()

	if err := s.RequestEmailChange(ctx, userID, "Pass-word1!", newEmail); err != nil {
		t.Fatal(err)
	}
	token := tokenOf(t, nextEmail(t, s).Body)
	// someone else took the address after the link was sent
	if err := s.db.UpdateUserEmail(ctx, dbqueries.UpdateUserEmailParams{ID: otherID, Email: newEmail}); err != nil {
		t.Fatal(err)
	}
	wantErr(t, "taken at confirm", s.ConfirmEmailChange(ctx, token), g.ErrorAccountAlreadyExists)
	if u, _ := s.GetUser(ctx, userID); u.Email != user.Email {
		t.Fatalf("email changed to a taken one: %q", u.Email)
	}
	noEmail(t, s)
}
