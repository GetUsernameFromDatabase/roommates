package services

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	g "roommates/internal/config"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

// fakeStore keeps objects in memory
type fakeStore struct {
	mu      sync.Mutex
	objects map[string][]byte
}

func newFakeStore() *fakeStore {
	return &fakeStore{mu: sync.Mutex{}, objects: map[string][]byte{}}
}

func (f *fakeStore) Put(_ context.Context, key string, body io.ReadSeeker, _ int64, _ string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objects[key] = b
	return nil
}

func (f *fakeStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return io.NopCloser(bytes.NewReader(f.objects[key])), nil
}

func (f *fakeStore) Delete(_ context.Context, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.objects, key)
	return nil
}

func (f *fakeStore) has(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.objects[key]
	return ok
}

// realPNG is a small image that decodes, so it gets a thumbnail
func realPNG(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 300, 100))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestCleanFileName(t *testing.T) {
	for in, want := range map[string]string{
		`C:\Users\me\bill.pdf`: "bill.pdf",
		"../../etc/passwd":     "passwd",
		" a\r\nb.txt ":         "ab.txt",
	} {
		if got, err := cleanFileName(in); err != nil || got != want {
			t.Errorf("cleanFileName(%q) = %q, %v, want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "dir/", "\x00"} {
		if _, err := cleanFileName(in); err == nil {
			t.Errorf("cleanFileName(%q) should fail", in)
		}
	}
}

func TestFiles(t *testing.T) {
	s := newTestService(t)
	store := s.store.(*fakeStore)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}}
	png := []byte("\x89PNG\r\n\x1a\n rest of the image")
	upload := func(actor Membership, target FileTarget, name string, body []byte) (pgtype.UUID, error) {
		return s.UploadFile(ctx, actor, target, name, bytes.NewReader(body), int64(len(body)))
	}

	_, err := upload(memberM, none, "empty.txt", nil)
	wantErr(t, "empty file", err, g.ErrorInvalidInput)
	_, err = s.UploadFile(ctx, memberM, none, "big.bin", bytes.NewReader(png), MaxFileSize+1)
	wantErr(t, "too large", err, g.ErrorFileTooLarge)

	loose, err := upload(memberM, none, "photo.png", png)
	if err != nil {
		t.Fatal(err)
	}
	meta, body, err := s.OpenFile(ctx, ownerM, loose, false)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(body)
	if meta.Mime != "image/png" || meta.Name != "photo.png" || !bytes.Equal(got, png) {
		t.Fatalf("opened %+v %q", meta, got)
	}
	// the fake PNG does not decode, so it is uploaded without a thumbnail
	_, _, err = s.OpenFile(ctx, ownerM, loose, true)
	wantErr(t, "thumbnail of a broken image", err, g.ErrorNotFound)
	_, _, err = s.OpenFile(ctx, strangerM, loose, false)
	wantErr(t, "stranger opens file", err, g.ErrorNotFound)
	wantErr(t, "owner deletes member's file", s.DeleteFile(ctx, ownerM, loose), g.ErrorNotFound)

	// a file on a private note is seen only by the note's author
	private, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Lease", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	onPrivate := FileTarget{NoteID: private, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}}
	_, err = upload(memberM, onPrivate, "x.txt", []byte("x"))
	wantErr(t, "member uploads to private note", err, g.ErrorNotFound)
	secret, err := upload(ownerM, onPrivate, "lease.png", realPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	_, thumb, err := s.OpenFile(ctx, ownerM, secret, true)
	if err != nil {
		t.Fatal(err)
	}
	thumbBytes, _ := io.ReadAll(thumb)
	if _, format, err := image.DecodeConfig(bytes.NewReader(thumbBytes)); err != nil || format != "jpeg" {
		t.Fatalf("thumbnail is %q, %v", format, err)
	}
	_, _, err = s.OpenFile(ctx, memberM, secret, true)
	wantErr(t, "member opens private thumbnail", err, g.ErrorNotFound)
	_, _, err = s.OpenFile(ctx, memberM, secret, false)
	wantErr(t, "member opens private file", err, g.ErrorNotFound)
	memberFiles, err := s.ListFiles(ctx, memberM)
	if err != nil {
		t.Fatal(err)
	}
	if len(memberFiles) != 1 || memberFiles[0].ID != loose {
		t.Fatalf("member lists %+v", memberFiles)
	}
	ownerFiles, err := s.ListFiles(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	if len(ownerFiles) != 2 || ownerFiles[0].ID != secret || ownerFiles[0].NoteTitle == nil ||
		!ownerFiles[0].HasThumb || ownerFiles[1].HasThumb {
		t.Fatalf("owner lists %+v", ownerFiles)
	}

	// files attach to an item of the household's list only
	listID, err := s.CreateShoppingList(ctx, ownerM, "Food", "groceries")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AddShoppingItem(ctx, ownerM, listID, "Cheese", ""); err != nil {
		t.Fatal(err)
	}
	list, err := s.GetShoppingList(ctx, ownerM, listID)
	if err != nil {
		t.Fatal(err)
	}
	onItem := FileTarget{NoteID: pgtype.UUID{}, ListID: listID, ShoppingItemID: list.Items[0].ID, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}}
	_, err = upload(strangerM, onItem, "x.txt", []byte("x"))
	wantErr(t, "stranger uploads to item", err, g.ErrorNotFound)
	if _, err := upload(memberM, onItem, "cheese.png", png); err != nil {
		t.Fatal(err)
	}
	attached, err := s.AttachedFiles(ctx, ownerM, onItem)
	if err != nil {
		t.Fatal(err)
	}
	if len(attached) != 1 || attached[0].Name != "cheese.png" {
		t.Fatalf("item files %+v", attached)
	}
	list, err = s.GetShoppingList(ctx, ownerM, listID)
	if err != nil {
		t.Fatal(err)
	}
	if list.Items[0].FileCount != 1 {
		t.Fatalf("item file count %d", list.Items[0].FileCount)
	}

	// deleting the note deletes its files, and the purge removes their objects
	var key, thumbKey string
	if err := s.pool.QueryRow(ctx, "SELECT object_key, thumb_key FROM files WHERE id = $1", secret).Scan(&key, &thumbKey); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteNote(ctx, ownerM, private); err != nil {
		t.Fatal(err)
	}
	_, _, err = s.OpenFile(ctx, ownerM, secret, false)
	wantErr(t, "file of deleted note", err, g.ErrorNotFound)
	if !store.has(key) || !store.has(thumbKey) {
		t.Fatal("objects removed before the purge")
	}
	if err := s.PurgeFiles(ctx); err != nil {
		t.Fatal(err)
	}
	if store.has(key) || store.has(thumbKey) {
		t.Fatal("objects left after the purge")
	}

	if err := s.DeleteFile(ctx, memberM, loose); err != nil {
		t.Fatal(err)
	}
}
