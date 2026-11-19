package services

import (
	"crypto/rand"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"strings"
	"testing"

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
