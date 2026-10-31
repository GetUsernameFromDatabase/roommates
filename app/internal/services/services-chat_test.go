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

func TestChatFileRefs(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := pgtype.UUID{}
	loose := FileTarget{NoteID: none, ListID: none, ShoppingItemID: none, ExpenseID: none, MessageID: none}

	lease, err := s.UploadFile(ctx, ownerM, loose, "lease.txt", strings.NewReader("lease"), 5)
	if err != nil {
		t.Fatal(err)
	}
	noteID, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Mine", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	onNote := loose
	onNote.NoteID = noteID
	secret, err := s.UploadFile(ctx, ownerM, onNote, "secret.txt", strings.NewReader("secret"), 6)
	if err != nil {
		t.Fatal(err)
	}

	wantErr(t, "reference to a private note's file", s.PostMessage(ctx, ownerM, "see", none, nil, []pgtype.UUID{secret}), g.ErrorInvalidInput)
	wantErr(t, "reference across households", s.PostMessage(ctx, strangerM, "see", none, nil, []pgtype.UUID{lease}), g.ErrorInvalidInput)
	wantErr(t, "too many references", s.PostMessage(ctx, ownerM, "", none, nil, make([]pgtype.UUID, MaxMessageFiles+1)), g.ErrorInvalidInput)
	if c, _ := s.ListChat(ctx, ownerM, none); len(c.Messages) != 0 {
		t.Fatalf("refused references left messages: %+v", c.Messages)
	}
	picks, err := s.ListReferableFiles(ctx, ownerM, "")
	if err != nil || len(picks) != 1 || picks[0].ID != lease {
		t.Fatalf("picks = %+v, %v", picks, err)
	}

	// a message of only a reference is sent, the member sees the file
	if err := s.PostMessage(ctx, ownerM, "", none, nil, []pgtype.UUID{lease}); err != nil {
		t.Fatal(err)
	}
	c, err := s.ListChat(ctx, memberM, none)
	if err != nil || len(c.Messages) != 1 || len(c.Messages[0].Files) != 1 || c.Messages[0].Files[0].ID != lease {
		t.Fatalf("chat = %+v, %v", c, err)
	}

	// the file stays when the message goes
	if err := s.DeleteMessage(ctx, ownerM, c.Messages[0].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetFile(ctx, memberM, lease); err != nil {
		t.Fatal(err)
	}

	// the reference goes when the file goes
	if err := s.PostMessage(ctx, ownerM, "the lease", none, nil, []pgtype.UUID{lease}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteFile(ctx, ownerM, lease); err != nil {
		t.Fatal(err)
	}
	c, err = s.ListChat(ctx, memberM, none)
	if err != nil || len(c.Messages) != 1 || len(c.Messages[0].Files) != 0 {
		t.Fatalf("after deleting the file: %+v, %v", c, err)
	}

	// a note made private later hides a file only on it, one also in the chat stays visible
	shared, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Ours", BodyMD: "", Private: false})
	if err != nil {
		t.Fatal(err)
	}
	onNote.NoteID = shared
	bill, err := s.UploadFile(ctx, ownerM, onNote, "bill.txt", strings.NewReader("bill"), 4)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.PostMessage(ctx, ownerM, "the bill", none, nil, []pgtype.UUID{bill}); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateNote(ctx, ownerM, shared, NoteInput{Title: "Ours", BodyMD: "", Private: true}); err != nil {
		t.Fatal(err)
	}
	seen := func(m Membership) int {
		t.Helper()
		c, err := s.ListChat(ctx, m, none)
		if err != nil {
			t.Fatal(err)
		}
		return len(c.Messages[0].Files)
	}
	if owner, member := seen(ownerM), seen(memberM); owner != 1 || member != 1 {
		t.Fatalf("after going private the owner sees %d files, the member %d", owner, member)
	}
	onNote.NoteID = shared
	hidden, err := s.UploadFile(ctx, ownerM, onNote, "hidden.txt", strings.NewReader("hidden"), 6)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.GetFile(ctx, memberM, hidden)
	wantErr(t, "member opens a file only on a private note", err, g.ErrorNotFound)
}

func TestChat(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := pgtype.UUID{}

	wantErr(t, "empty message", s.PostMessage(ctx, memberM, " \r\n ", none, nil, nil), g.ErrorInvalidInput)
	wantErr(t, "too long", s.PostMessage(ctx, memberM, strings.Repeat("a", MaxMessageLength+1), none, nil, nil), g.ErrorInvalidInput)
	// a line break counts once, as in the textarea's maxlength
	if err := s.PostMessage(ctx, memberM, strings.Repeat("a\r\n", MaxMessageLength/2-1)+"a", none, nil, nil); err != nil {
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

	photos := []Upload{
		{Name: "cat.png", Body: bytes.NewReader([]byte("\x89PNG\r\n\x1a\n")), Size: 8},
		{Name: "dog.txt", Body: strings.NewReader("woof"), Size: 4},
	}
	if err := s.PostMessage(ctx, ownerM, "", none, photos, nil); err != nil {
		t.Fatal(err)
	}
	// a failed upload takes its message with it
	tooLarge := []Upload{
		{Name: "ok.txt", Body: strings.NewReader("x"), Size: 1},
		{Name: "big.bin", Body: bytes.NewReader([]byte("x")), Size: MaxFileSize + 1},
	}
	wantErr(t, "too large file", s.PostMessage(ctx, ownerM, "big one", none, tooLarge, nil), g.ErrorFileTooLarge)
	wantErr(t, "too many files", s.PostMessage(ctx, ownerM, "", none, make([]Upload, MaxMessageFiles+1), nil), g.ErrorInvalidInput)

	c, err := s.ListChat(ctx, ownerM, none)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Messages) != 2 || c.HasOlder || len(c.Messages[0].Files) != 2 {
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
	if files, _ := s.ListFiles(ctx, ownerM, FileFilter{}); len(files) != 0 {
		t.Fatalf("files of a deleted message are left: %+v", files)
	}

	// pages of ChatPageSize, newest first
	for i := range ChatPageSize {
		if err := s.PostMessage(ctx, memberM, string(rune('a'+i%26)), none, nil, nil); err != nil {
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

// a reply quotes a message of its own household until that message is deleted
func TestChatReply(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := pgtype.UUID{}

	if err := s.PostMessage(ctx, ownerM, "Who took my milk?", none, nil, nil); err != nil {
		t.Fatal(err)
	}
	c, err := s.ListChat(ctx, ownerM, none)
	if err != nil {
		t.Fatal(err)
	}
	question := c.Messages[0].ID
	wantErr(t, "stranger replies", s.PostMessage(ctx, strangerM, "not me", question, nil, nil), g.ErrorNotFound)
	_, err = s.GetReplyPreview(ctx, strangerM, question)
	wantErr(t, "stranger previews", err, g.ErrorNotFound)
	if err := s.PostMessage(ctx, memberM, "Sorry", question, nil, nil); err != nil {
		t.Fatal(err)
	}
	c, err = s.ListChat(ctx, ownerM, none)
	if err != nil {
		t.Fatal(err)
	}
	reply := c.Messages[0]
	if reply.ReplyToID != question || reply.ReplyAuthorName == nil || reply.ReplyBody != "Who took my milk?" {
		t.Fatalf("reply %+v", reply)
	}

	if err := s.DeleteMessage(ctx, ownerM, question); err != nil {
		t.Fatal(err)
	}
	c, err = s.ListChat(ctx, ownerM, none)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Messages) != 1 || c.Messages[0].ReplyToID.Valid || c.Messages[0].ReplyAuthorName != nil {
		t.Fatalf("reply of a deleted message %+v", c.Messages)
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

	if err := s.PostMessage(ctx, memberM, "hi", pgtype.UUID{}, nil, nil); err != nil {
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
		ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
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
