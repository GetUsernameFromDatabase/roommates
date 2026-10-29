package services

import (
	"bytes"
	"context"
	g "roommates/internal/config"
	"roommates/internal/events"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestChat(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := pgtype.UUID{}

	wantErr(t, "empty message", s.PostMessage(ctx, memberM, " \r\n ", nil), g.ErrorInvalidInput)
	wantErr(t, "too long", s.PostMessage(ctx, memberM, strings.Repeat("a", MaxMessageLength+1), nil), g.ErrorInvalidInput)
	// a line break counts once, as in the textarea's maxlength
	if err := s.PostMessage(ctx, memberM, strings.Repeat("a\r\n", MaxMessageLength/2-1)+"a", nil); err != nil {
		t.Fatal(err)
	}

	unread := func(m Membership) int64 {
		t.Helper()
		return mustMembership(t, s, m.HouseholdID, m.UserID).UnreadMessages
	}
	if got := unread(ownerM); got != 1 {
		t.Fatalf("owner has %d unread, want 1", got)
	}
	if got := unread(memberM); got != 0 {
		t.Fatalf("own message counts as unread: %d", got)
	}

	photo := &Upload{Name: "cat.png", Body: bytes.NewReader([]byte("\x89PNG\r\n\x1a\n")), Size: 8}
	if err := s.PostMessage(ctx, ownerM, "", photo); err != nil {
		t.Fatal(err)
	}
	// a failed upload takes its message with it
	tooLarge := &Upload{Name: "big.bin", Body: bytes.NewReader([]byte("x")), Size: MaxFileSize + 1}
	wantErr(t, "too large file", s.PostMessage(ctx, ownerM, "big one", tooLarge), g.ErrorFileTooLarge)

	c, err := s.ListChat(ctx, ownerM, none)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Messages) != 2 || c.HasOlder || len(c.Messages[0].Files) != 1 || c.Messages[0].Files[0].Name != "cat.png" {
		t.Fatalf("chat %+v", c)
	}
	if got := unread(ownerM); got != 0 {
		t.Fatalf("owner has %d unread after reading", got)
	}
	other, err := s.ListChat(ctx, strangerM, none)
	if err != nil || len(other.Messages) != 0 {
		t.Fatalf("stranger sees %+v, %v", other, err)
	}

	// only the author attaches to or deletes a message
	photoMessage := c.Messages[0].ID
	onMessage := FileTarget{NoteID: none, ListID: none, ShoppingItemID: none, ExpenseID: none, MessageID: photoMessage}
	_, err = s.UploadFile(ctx, memberM, onMessage, "x.txt", strings.NewReader("x"), 1)
	wantErr(t, "member attaches to owner's message", err, g.ErrorNotFound)
	wantErr(t, "member deletes", s.DeleteMessage(ctx, memberM, photoMessage), g.ErrorNotFound)
	wantErr(t, "stranger deletes", s.DeleteMessage(ctx, strangerM, photoMessage), g.ErrorNotFound)
	if err := s.DeleteMessage(ctx, ownerM, photoMessage); err != nil {
		t.Fatal(err)
	}
	if files, _ := s.ListFiles(ctx, ownerM); len(files) != 0 {
		t.Fatalf("files of a deleted message are left: %+v", files)
	}

	// pages of ChatPageSize, newest first
	for i := range ChatPageSize {
		if err := s.PostMessage(ctx, memberM, string(rune('a'+i%26)), nil); err != nil {
			t.Fatal(err)
		}
	}
	first, err := s.ListChat(ctx, memberM, none)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Messages) != ChatPageSize || !first.HasOlder {
		t.Fatalf("first page has %d messages, older %v", len(first.Messages), first.HasOlder)
	}
	older, err := s.ListChat(ctx, memberM, first.Messages[ChatPageSize-1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(older.Messages) != 1 || older.HasOlder || older.Messages[0].Body[0] != 'a' || len(older.Messages[0].Body) < 100 {
		t.Fatalf("older page %+v", older)
	}
}

// the triggers tell the household what changed once the transaction commits
func TestHouseholdEvents(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "LISTEN "+events.Channel); err != nil {
		t.Fatal(err)
	}
	defer conn.Exec(context.WithoutCancel(ctx), "UNLISTEN *") //nolint:errcheck // test cleanup
	next := func(step, want string) {
		t.Helper()
		waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		n, err := conn.Conn().WaitForNotification(waitCtx)
		if err != nil {
			t.Fatalf("%s: %v", step, err)
		}
		if want = ownerM.HouseholdID.String() + " " + want; n.Payload != want {
			t.Fatalf("%s: got %q, want %q", step, n.Payload, want)
		}
	}

	if err := s.PostMessage(ctx, memberM, "hi", nil); err != nil {
		t.Fatal(err)
	}
	next("message", events.TopicChat)

	// reading the chat tells only the badges, not the chat, which would read again
	if _, err := s.ListChat(ctx, ownerM, pgtype.UUID{}); err != nil {
		t.Fatal(err)
	}
	next("read", events.TopicRead)
	expenseID, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
		Title: "Pizza", Amount: "10", Category: "groceries", IncurredOn: "2026-10-01",
		Custom: false, Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	// the shares' rows notify too, Postgres folds equal payloads of one transaction
	next("expense", events.TopicExpenses)
	if err := s.SetExpenseSharePaid(ctx, ownerM, expenseID, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	next("share paid", events.TopicExpenses)
}
