package services

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/routes"
	"roommates/internal/utils"

	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// dataExportPendingFor is how long a pending export is retried before the purge drops it
const dataExportPendingFor = 24 * time.Hour

// DataExportStatus is the user's newest export, ID is not Valid when there is none
// (or it expired), ReadyAt is zero while it is being built
type DataExportStatus struct {
	ID        pgtype.UUID
	ReadyAt   time.Time
	ExpiresAt time.Time
}

// Pending is true while the export is being built
func (d DataExportStatus) Pending() bool { return d.ID.Valid && d.ReadyAt.IsZero() }

// Ready is true while the export can be downloaded
func (d DataExportStatus) Ready() bool { return d.ID.Valid && !d.ReadyAt.IsZero() }

// RequestDataExport queues a ZIP of the user's data, nothing while one is pending,
// otherwise their older exports are deleted; lang is the language of the email
func (s *DefaultService) RequestDataExport(ctx context.Context, userID pgtype.UUID, lang string) error {
	return s.db.RequestDataExport(ctx, dbqueries.RequestDataExportParams{UserID: userID, Lang: lang})
}

// GetDataExportStatus is the user's newest export for the settings page
func (s *DefaultService) GetDataExportStatus(ctx context.Context, userID pgtype.UUID) (DataExportStatus, error) {
	none := DataExportStatus{ID: pgtype.UUID{}, ReadyAt: time.Time{}, ExpiresAt: time.Time{}}
	row, err := s.db.GetLatestDataExport(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return none, nil
	}
	if err != nil {
		return none, err
	}
	if !row.ReadyAt.Valid {
		return DataExportStatus{ID: row.ID, ReadyAt: time.Time{}, ExpiresAt: time.Time{}}, nil
	}
	expires := row.ReadyAt.Time.Add(g.DataExportValidFor)
	if time.Now().After(expires) {
		return none, nil
	}
	return DataExportStatus{ID: row.ID, ReadyAt: row.ReadyAt.Time, ExpiresAt: expires}, nil
}

// OpenDataExport is the user's own ready, unexpired export, [g.ErrorNotFound] otherwise;
// the caller closes the contents
func (s *DefaultService) OpenDataExport(ctx context.Context, userID, exportID pgtype.UUID) (io.ReadCloser, int64, time.Time, error) {
	row, err := s.db.GetReadyDataExport(ctx, dbqueries.GetReadyDataExportParams{
		ID:         exportID,
		UserID:     userID,
		ReadyAfter: timestamptz(time.Now().Add(-g.DataExportValidFor)),
	})
	if err != nil {
		return nil, 0, time.Time{}, noRowsAs(err, g.ErrorNotFound)
	}
	if row.ObjectKey == nil || row.SizeBytes == nil {
		return nil, 0, time.Time{}, g.ErrorNotFound
	}
	r, err := s.store.Get(ctx, *row.ObjectKey)
	return r, *row.SizeBytes, row.ReadyAt.Time, err
}

// PurgeExpiredExports deletes expired exports and pending ones that kept failing,
// the trigger queues their objects, run it before [DefaultService.PurgeFiles]
func (s *DefaultService) PurgeExpiredExports(ctx context.Context) error {
	now := time.Now()
	return s.db.DeleteExpiredDataExports(ctx, dbqueries.DeleteExpiredDataExportsParams{
		ReadyBefore:   timestamptz(now.Add(-g.DataExportValidFor)),
		PendingBefore: timestamptz(now.Add(-dataExportPendingFor)),
	})
}

// BuildPendingExports builds the pending export tried longest ago, stores it and emails its link;
// on failure it stays pending, behind the others, and is tried again on a later run
//
// shortcut: one export per run inside the one-minute cleanup loop, so a huge export delays the
// other cleanups and a queue drains at one a minute; move to its own goroutine if exports grow big
func (s *DefaultService) BuildPendingExports(ctx context.Context) error {
	// claimed in its own statement, so the try counts even when the build fails
	row, err := s.db.ClaimPendingDataExport(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	lctx, err := ctxi18n.WithLocale(ctx, row.Lang)
	if err != nil {
		lctx = ctx
	}
	f, err := os.CreateTemp("", "roommates-export-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name()) //nolint:errcheck // a temp file
	defer f.Close()           //nolint:errcheck // closed or failed already
	if err := s.writeDataExport(lctx, row.UserID, f); err != nil {
		return err
	}
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	key := "exports/" + row.ID.String() + ".zip"
	if err := s.store.Put(ctx, key, f, size, "application/zip"); err != nil {
		return err
	}
	readyAt, err := s.db.SetDataExportReady(ctx, dbqueries.SetDataExportReadyParams{ObjectKey: &key, SizeBytes: &size, ID: row.ID})
	if errors.Is(err, pgx.ErrNoRows) {
		// purged or erased while it was built, nothing points to the object any more
		return s.store.Delete(ctx, key)
	}
	if err != nil {
		return err
	}
	user, err := s.db.GetUser(ctx, row.UserID)
	if err != nil {
		return err
	}
	link := g.GetAppBaseURL() + routes.URL(routes.DataExportDownload, row.ID)
	s.sendAsync(mailer.Message{
		To:      user.Email,
		Subject: utils.T(lctx, locales.LK_MailDataExportSubject),
		Body: utils.T(lctx, locales.LK_MailDataExportBody,
			user.Username, utils.FormatDate(lctx, readyAt.Time.Add(g.DataExportValidFor)), link),
	})
	return nil
}

// writeDataExport writes the ZIP of what userID wrote or took part in to w, other people
// only appear by name; ctx gives the language of "Former member"
func (s *DefaultService) writeDataExport(ctx context.Context, userID pgtype.UUID, w io.Writer) error {
	account, err := s.db.GetDataExportAccount(ctx, userID)
	if err != nil {
		return err
	}
	memberships, err := s.db.ListDataExportMemberships(ctx, userID)
	if err != nil {
		return err
	}
	z := zip.NewWriter(w)
	if err := zipJSON(z, "account.json", map[string]any{"profile": account, "memberships": memberships}); err != nil {
		return err
	}
	for _, m := range memberships {
		hid := m.HouseholdID.String()
		dir := zipName(m.Household) + "-" + hid[:8]
		if err := s.writeHouseholdExport(ctx, z, dir, m.HouseholdID, userID); err != nil {
			return err
		}
	}
	return z.Close()
}

// exportFile is an uploaded file in the household JSON, path is where it is in the ZIP
type exportFile struct {
	Name      string             `json:"name"`
	Mime      string             `json:"mime"`
	SizeBytes int64              `json:"size_bytes"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Path      string             `json:"path"`
}

// writeHouseholdExport adds households/<dir>.json and the user's files under files/<dir>/
func (s *DefaultService) writeHouseholdExport(ctx context.Context, z *zip.Writer, dir string, householdID, userID pgtype.UUID) error {
	expenses, err := s.db.ListDataExportExpenses(ctx, dbqueries.ListDataExportExpensesParams{UserID: userID, HouseholdID: householdID})
	if err != nil {
		return err
	}
	for i := range expenses {
		expenses[i].PayerName = utils.MemberName(ctx, expenses[i].PayerName)
	}
	settlements, err := s.db.ListDataExportSettlements(ctx, dbqueries.ListDataExportSettlementsParams{UserID: userID, HouseholdID: householdID})
	if err != nil {
		return err
	}
	for i := range settlements {
		settlements[i].CounterpartyName = utils.MemberName(ctx, settlements[i].CounterpartyName)
	}
	messages, err := s.db.ListDataExportMessages(ctx, dbqueries.ListDataExportMessagesParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	notes, err := s.db.ListDataExportNotes(ctx, dbqueries.ListDataExportNotesParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	shopping, err := s.db.ListDataExportShoppingItems(ctx, dbqueries.ListDataExportShoppingItemsParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	chores, err := s.db.ListDataExportChoreLog(ctx, dbqueries.ListDataExportChoreLogParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	events, err := s.db.ListDataExportEvents(ctx, dbqueries.ListDataExportEventsParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	rows, err := s.db.ListDataExportFiles(ctx, dbqueries.ListDataExportFilesParams{HouseholdID: householdID, UserID: userID})
	if err != nil {
		return err
	}
	files := make([]exportFile, len(rows))
	taken := map[string]bool{}
	for i, f := range rows {
		path := "files/" + dir + "/" + uniqueZipName(zipName(f.Name), taken)
		files[i] = exportFile{Name: f.Name, Mime: f.Mime, SizeBytes: f.SizeBytes, CreatedAt: f.CreatedAt, Path: path}
	}
	if err := zipJSON(z, "households/"+dir+".json", map[string]any{
		"expenses":       expenses,
		"settlements":    settlements,
		"messages":       messages,
		"notes":          notes,
		"shopping_items": shopping,
		"chore_log":      chores,
		"events":         events,
		"files":          files,
	}); err != nil {
		return err
	}
	for i, f := range rows {
		if err := s.zipObject(ctx, z, files[i].Path, f.ObjectKey); err != nil {
			return err
		}
	}
	return nil
}

func (s *DefaultService) zipObject(ctx context.Context, z *zip.Writer, name, key string) error {
	r, err := s.store.Get(ctx, key)
	if err != nil {
		return err
	}
	defer r.Close() //nolint:errcheck // read only
	w, err := z.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

func zipJSON(z *zip.Writer, name string, v any) error {
	w, err := z.Create(name)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// zipName is s as one ZIP path segment: no separators, control characters or leading dots
func zipName(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || unicode.IsControl(r) {
			return '_'
		}
		return r
	}, s)
	s = strings.TrimLeft(strings.TrimSpace(s), ".")
	if s == "" {
		return "file"
	}
	return s
}

// uniqueZipName is name, or name-2, name-3 … before its extension when taken already
func uniqueZipName(name string, taken map[string]bool) string {
	stem, ext := SplitFileName(name)
	out := name
	for n := 2; taken[strings.ToLower(out)]; n++ {
		out = stem + "-" + strconv.Itoa(n) + ext
	}
	taken[strings.ToLower(out)] = true
	return out
}
