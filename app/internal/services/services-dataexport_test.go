package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	g "roommates/internal/config"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestZipName(t *testing.T) {
	for in, want := range map[string]string{"../../etc": "_.._etc", "..": "file", `a\b:c`: "a_b_c", " Flat ": "Flat"} {
		if got := zipName(in); got != want {
			t.Errorf("zipName(%q) = %q, want %q", in, got, want)
		}
	}
	taken := map[string]bool{}
	for _, want := range []string{"a.txt", "a-2.txt", "a-3.txt"} {
		if got := uniqueZipName("a.txt", taken); got != want {
			t.Errorf("uniqueZipName = %q, want %q", got, want)
		}
	}
}

// buildExport runs the builder until userID's export is ready, other tests' pending ones go first
// and may fail on purpose
func buildExport(t *testing.T, s *DefaultService, userID pgtype.UUID) DataExportStatus {
	t.Helper()
	for range 50 {
		if err := s.BuildPendingExports(t.Context()); err != nil {
			t.Log(err)
		}
		st, err := s.GetDataExportStatus(t.Context(), userID)
		if err != nil {
			t.Fatal(err)
		}
		if st.Ready() {
			return st
		}
	}
	t.Fatal("export never became ready")
	return DataExportStatus{ID: pgtype.UUID{}, ReadyAt: time.Time{}, ExpiresAt: time.Time{}}
}

func TestDataExport(t *testing.T) {
	s := newTestService(t)
	ctx := localeCtx(t, "en")
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	none := pgtype.UUID{}

	in := ExpenseInput{
		Title: "Rent", Amount: "10", Category: "rent", IncurredOn: "2026-10-01", Custom: false,
		Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}
	if _, err := s.CreateExpense(ctx, ownerM, in); err != nil {
		t.Fatal(err)
	}
	if err := s.PostMessage(ctx, memberM, "my own words", none, nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := s.PostMessage(ctx, ownerM, "the owner's words", none, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateNote(ctx, memberM, NoteInput{Title: "Secret diary", BodyMD: "", Private: true}); err != nil {
		t.Fatal(err)
	}
	target := FileTarget{NoteID: none, ListID: none, ShoppingItemID: none, ExpenseID: none, MessageID: none, EventID: none}
	body := []byte("the lease text")
	if _, err := s.UploadFile(ctx, memberM, target, "lease.txt", bytes.NewReader(body), int64(len(body))); err != nil {
		t.Fatal(err)
	}

	if err := s.RequestDataExport(ctx, memberM.UserID, "en"); err != nil {
		t.Fatal(err)
	}
	pending, err := s.GetDataExportStatus(ctx, memberM.UserID)
	if err != nil || !pending.Pending() {
		t.Fatalf("status after request = %+v, %v, want pending", pending, err)
	}
	if err := s.RequestDataExport(ctx, memberM.UserID, "en"); err != nil {
		t.Fatal(err)
	}
	if again, _ := s.GetDataExportStatus(ctx, memberM.UserID); again.ID != pending.ID {
		t.Fatal("a second request while pending made a new export")
	}

	st := buildExport(t, s, memberM.UserID)
	if st.ID != pending.ID {
		t.Fatal("the built export is not the requested one")
	}
	r, size, _, err := s.OpenDataExport(ctx, memberM.UserID, st.ID)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(r)
	if err != nil || int64(len(data)) != size {
		t.Fatalf("read %d bytes of %d: %v", len(data), size, err)
	}
	z, err := zip.NewReader(bytes.NewReader(data), size)
	if err != nil {
		t.Fatal(err)
	}
	entries := map[string]string{}
	for _, f := range z.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(rc)
		rc.Close() //nolint:errcheck,gosec // read only
		entries[f.Name] = string(b)
	}
	if !strings.Contains(entries["account.json"], `"household": "Flat"`) {
		t.Errorf("account.json lacks the membership:\n%s", entries["account.json"])
	}
	prefix := "Flat-" + ownerM.HouseholdID.String()[:8]
	household := entries["households/"+prefix+".json"]
	for _, want := range []string{`"title": "Rent"`, "my own words", "Secret diary", `"name": "lease.txt"`} {
		if !strings.Contains(household, want) {
			t.Errorf("household JSON lacks %q:\n%s", want, household)
		}
	}
	if strings.Contains(household, "the owner's words") {
		t.Error("household JSON has another member's message")
	}
	if got := entries["files/"+prefix+"/lease.txt"]; got != string(body) {
		t.Errorf("uploaded file = %q, want %q", got, body)
	}

	if _, _, _, err := s.OpenDataExport(ctx, strangerM.UserID, st.ID); !errors.Is(err, g.ErrorNotFound) {
		t.Errorf("another user opening the export: %v, want not found", err)
	}

	// expire it: ready long ago
	if _, err := s.pool.Exec(ctx, "UPDATE data_exports SET ready_at = $2 WHERE id = $1",
		st.ID, time.Now().Add(-g.DataExportValidFor-time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.OpenDataExport(ctx, memberM.UserID, st.ID); !errors.Is(err, g.ErrorNotFound) {
		t.Errorf("opening an expired export: %v, want not found", err)
	}
	if err := s.PurgeExpiredExports(ctx); err != nil {
		t.Fatal(err)
	}
	var left, queued int
	if err := s.pool.QueryRow(ctx, "SELECT count(*) FROM data_exports WHERE id = $1", st.ID).Scan(&left); err != nil {
		t.Fatal(err)
	}
	if err := s.pool.QueryRow(ctx, "SELECT count(*) FROM file_purges WHERE object_key = $1",
		"exports/"+st.ID.String()+".zip").Scan(&queued); err != nil {
		t.Fatal(err)
	}
	if left != 0 || queued != 1 {
		t.Errorf("after the purge %d rows left and %d objects queued, want 0 and 1", left, queued)
	}
}

func TestDataExportFailureDoesNotBlock(t *testing.T) {
	s := newTestService(t)
	ctx := localeCtx(t, "en")
	ownerM, memberM, _ := newTestHousehold(t, s)
	// a file whose object cannot be read, so the owner's export can never be built
	missing := "missing/" + ownerM.UserID.String()
	s.store.(*fakeStore).failKey = missing
	if _, err := s.pool.Exec(ctx, `INSERT INTO files (household_id, uploader_id, object_key, name, mime, size_bytes)
		VALUES ($1, $2, $3, 'gone.txt', 'text/plain', 1)`, ownerM.HouseholdID, ownerM.UserID, missing); err != nil {
		t.Fatal(err)
	}
	if err := s.RequestDataExport(ctx, ownerM.UserID, "en"); err != nil {
		t.Fatal(err)
	}
	if err := s.RequestDataExport(ctx, memberM.UserID, "en"); err != nil {
		t.Fatal(err)
	}
	// the older, failing export goes to the back instead of being picked every run
	buildExport(t, s, memberM.UserID)
	if st, err := s.GetDataExportStatus(ctx, ownerM.UserID); err != nil || !st.Pending() {
		t.Fatalf("failing export = %+v, %v, want still pending", st, err)
	}
}
