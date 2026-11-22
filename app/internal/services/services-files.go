package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"path"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/images"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxFileSize = 25 << 20
	// MaxUploadBody is the largest request body of an upload: a file and room for multipart headers and other fields
	MaxUploadBody     = MaxFileSize + 1<<20
	MaxFileNameLength = 200
	// http.DetectContentType looks at the first 512 bytes only
	sniffLength = 512
	// files removed from storage per query when purging
	purgeBatch = 100
	// MaxReferablePicks is how many files the chat's picker shows, typing narrows them down
	MaxReferablePicks = 20
	// longest side of a thumbnail, sharp at 48px on a 3x screen
	thumbnailSide = 256
)

// ImageMimes are the sniffed types [images.Thumbnail] decodes and pages show as images
var ImageMimes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

// what [FileFilter] can narrow files down to, the place "none" is not attached to anything
var (
	FileKinds  = []string{"image", "pdf", "other"}
	FilePlaces = []string{"note", "item", "expense", "message", "event", "none"}
)

// FileFilter narrows a file list, empty fields do not, Query is a part of the name in any case
type FileFilter struct {
	Query string
	Kind  string
	Place string
}

// Active says whether f narrows anything down
func (f FileFilter) Active() bool {
	return strings.TrimSpace(f.Query) != "" || slices.Contains(FileKinds, f.Kind) || slices.Contains(FilePlaces, f.Place)
}

// FileTarget is what an upload attaches to: a note, an item of a list, an expense,
// a chat message, a calendar event, or nothing
type FileTarget struct {
	NoteID         pgtype.UUID
	ListID         pgtype.UUID
	ShoppingItemID pgtype.UUID
	ExpenseID      pgtype.UUID
	MessageID      pgtype.UUID
	EventID        pgtype.UUID
}

type AttachedFile = dbqueries.ListAttachedFilesRow

// FileMeta is a file without its contents, with what it is attached to
type FileMeta = dbqueries.GetFileRow

// cleanFileName drops any path and control characters,
// [g.ErrorInvalidInput] unless 1..[MaxFileNameLength] runes are left
func cleanFileName(name string) (string, error) {
	name = name[strings.LastIndexAny(name, `/\`)+1:]
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	name, ok := trimToLength(name, 1, MaxFileNameLength)
	if !ok {
		return "", g.ErrorInvalidInput
	}
	return name, nil
}

// SplitFileName cuts the extension off name, a name that is all extension like .bashrc has none,
// and neither has one whose last dot is inside a sentence, like "Report v1.2 final"
func SplitFileName(name string) (stem, ext string) {
	ext = path.Ext(name)
	if ext == name || strings.Contains(ext, " ") {
		return name, ""
	}
	return strings.TrimSuffix(name, ext), ext
}

// checkFileTarget returns [g.ErrorNotFound] unless the actor sees the note, item, expense or
// event, or wrote the message, and [g.ErrorForbidden] unless they may change the event
func (s *DefaultService) checkFileTarget(ctx context.Context, actor Membership, t FileTarget) error {
	if t.EventID.Valid {
		e, err := s.db.GetEvent(ctx, dbqueries.GetEventParams{ID: t.EventID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		return canChangeEvent(actor, e)
	}
	if t.NoteID.Valid {
		_, err := s.GetNote(ctx, actor, t.NoteID)
		return err
	}
	if t.ShoppingItemID.Valid {
		_, err := s.GetShoppingItem(ctx, actor, t.ListID, t.ShoppingItemID)
		return err
	}
	if t.ExpenseID.Valid {
		_, err := s.db.GetExpense(ctx, dbqueries.GetExpenseParams{ID: t.ExpenseID, HouseholdID: actor.HouseholdID})
		return noRowsAs(err, g.ErrorNotFound)
	}
	if t.MessageID.Valid {
		author, err := s.db.GetMessageAuthor(ctx, dbqueries.GetMessageAuthorParams{ID: t.MessageID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if author != actor.UserID {
			return g.ErrorNotFound
		}
	}
	return nil
}

// UploadFile stores size bytes of body as name, attached to t, the type is sniffed from the contents
//
// returns [g.ErrorFileTooLarge] above [MaxFileSize], [g.ErrorInvalidInput] for an empty file or a bad name
func (s *DefaultService) UploadFile(ctx context.Context, actor Membership, t FileTarget,
	name string, body io.ReadSeeker, size int64) (pgtype.UUID, error) {
	if err := s.checkFileTarget(ctx, actor, t); err != nil {
		return pgtype.UUID{}, err
	}
	f, err := s.storeUpload(ctx, actor, Upload{Name: name, Body: body, Size: size})
	if err != nil {
		return pgtype.UUID{}, err
	}
	var fileID pgtype.UUID
	err = s.inTx(ctx, func(q *dbqueries.Queries) error {
		fileID, err = insertStoredFile(ctx, q, actor, t, f)
		return err
	})
	if err != nil {
		s.discardStored(ctx, f)
	}
	return fileID, err
}

// storedFile is an upload whose contents are in the store but which has no row yet
type storedFile struct {
	name, mime, key string
	thumbKey        *string
	size            int64
}

// storeUpload checks u and puts its contents and thumbnail in the store, errors as
// [DefaultService.UploadFile]; without a row nothing would purge them, so a caller that
// does not insert it calls [DefaultService.discardStored]
func (s *DefaultService) storeUpload(ctx context.Context, actor Membership, u Upload) (storedFile, error) {
	if u.Size > MaxFileSize {
		return storedFile{}, g.ErrorFileTooLarge
	}
	name, err := cleanFileName(u.Name)
	if err != nil || u.Size <= 0 {
		return storedFile{}, g.ErrorInvalidInput
	}

	head := make([]byte, sniffLength)
	n, err := io.ReadFull(u.Body, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return storedFile{}, err
	}
	mime := http.DetectContentType(head[:n])
	if _, err := u.Body.Seek(0, io.SeekStart); err != nil {
		return storedFile{}, err
	}

	var thumb []byte
	if ImageMimes[mime] {
		// a broken image is still uploaded, only without a thumbnail
		if thumb, err = images.Thumbnail(u.Body, thumbnailSide); err != nil {
			log.Warn().Err(err).Str("mime", mime).Msg("could not make a thumbnail")
		}
		if _, err := u.Body.Seek(0, io.SeekStart); err != nil {
			return storedFile{}, err
		}
	}

	f := storedFile{name: name, mime: mime, key: "households/" + actor.HouseholdID.String() + "/" + rand.Text(), thumbKey: nil, size: u.Size}
	if err := s.store.Put(ctx, f.key, u.Body, f.size, mime); err != nil {
		return storedFile{}, err
	}
	if thumb != nil {
		k := f.key + "-thumb"
		if err := s.store.Put(ctx, k, bytes.NewReader(thumb), int64(len(thumb)), "image/jpeg"); err != nil {
			s.discardStored(ctx, f)
			return storedFile{}, err
		}
		f.thumbKey = &k
	}
	return f, nil
}

// discardStored deletes the objects of files whose rows were not inserted
func (s *DefaultService) discardStored(ctx context.Context, files ...storedFile) {
	for _, f := range files {
		keys := []string{f.key}
		if f.thumbKey != nil {
			keys = append(keys, *f.thumbKey)
		}
		for _, k := range keys {
			if err := s.store.Delete(context.WithoutCancel(ctx), k); err != nil {
				log.Error().Err(err).Str("key", k).Msg("could not delete object of a failed upload")
			}
		}
	}
}

// insertStoredFile adds the row of f attached to t, which the caller has checked
func insertStoredFile(ctx context.Context, q *dbqueries.Queries, actor Membership, t FileTarget, f storedFile) (pgtype.UUID, error) {
	// a file on nothing is kept, without a reference nothing else would hold it
	kept := !t.NoteID.Valid && !t.ShoppingItemID.Valid && !t.ExpenseID.Valid && !t.MessageID.Valid && !t.EventID.Valid
	fileID, err := q.InsertFile(ctx, dbqueries.InsertFileParams{
		HouseholdID: actor.HouseholdID,
		UploaderID:  actor.UserID,
		ObjectKey:   f.key,
		ThumbKey:    f.thumbKey,
		Name:        f.name,
		Mime:        f.mime,
		SizeBytes:   f.size,
		Kept:        kept,
	})
	if err != nil {
		return fileID, err
	}
	if !kept {
		if err := q.InsertFileRef(ctx, dbqueries.InsertFileRefParams{
			FileID:         fileID,
			HouseholdID:    actor.HouseholdID,
			NoteID:         t.NoteID,
			ShoppingItemID: t.ShoppingItemID,
			ExpenseID:      t.ExpenseID,
			MessageID:      t.MessageID,
			EventID:        t.EventID,
		}); err != nil {
			return fileID, err
		}
	}
	// chat files are chat's business, a private note's files are for its author
	if t.MessageID.Valid {
		return fileID, nil
	}
	visibility := g.NoteVisibilityHousehold
	if t.NoteID.Valid {
		n, err := q.GetNote(ctx, dbqueries.GetNoteParams{ID: t.NoteID, HouseholdID: actor.HouseholdID, UserID: actor.UserID})
		if err != nil {
			return fileID, noRowsAs(err, g.ErrorNotFound)
		}
		visibility = n.Visibility
	}
	return fileID, logNote(ctx, q, actor, g.ActFileUploaded, fileID, f.name, visibility)
}

// ListFiles returns the household's files that match filter, except those on others' private notes
func (s *DefaultService) ListFiles(ctx context.Context, actor Membership, filter FileFilter) ([]dbqueries.ListFilesRow, error) {
	if !slices.Contains(FileKinds, filter.Kind) {
		filter.Kind = ""
	}
	if !slices.Contains(FilePlaces, filter.Place) {
		filter.Place = ""
	}
	return s.db.ListFiles(ctx, dbqueries.ListFilesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
		SharedOnly:  false,
		Q:           strings.TrimSpace(filter.Query),
		Kind:        filter.Kind,
		Place:       filter.Place,
	})
}

// ListReferableFiles returns up to [MaxReferablePicks] files whose name has q that a chat
// message can refer to: all but those on private notes, the others could not open them
func (s *DefaultService) ListReferableFiles(ctx context.Context, actor Membership, q string) ([]dbqueries.ListFilesRow, error) {
	return s.filePicks(ctx, actor, q, nil)
}

// filePicks returns up to [MaxReferablePicks] files whose name has q that all members can see,
// except those in skip
func (s *DefaultService) filePicks(ctx context.Context, actor Membership, q string, skip []AttachedFile) ([]dbqueries.ListFilesRow, error) {
	rows, err := s.db.ListFiles(ctx, dbqueries.ListFilesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
		SharedOnly:  true,
		Q:           strings.TrimSpace(q),
		Kind:        "",
		Place:       "",
	})
	rows = slices.DeleteFunc(rows, func(r dbqueries.ListFilesRow) bool {
		return slices.ContainsFunc(skip, func(a AttachedFile) bool { return a.ID == r.ID })
	})
	return rows[:min(len(rows), MaxReferablePicks)], err
}

// ListAttachableFiles returns up to [MaxReferablePicks] files whose name has q that are not on t
// yet: all but those on private notes, the others could not open them
func (s *DefaultService) ListAttachableFiles(ctx context.Context, actor Membership, t FileTarget, q string) ([]dbqueries.ListFilesRow, error) {
	if err := s.checkFileTarget(ctx, actor, t); err != nil {
		return nil, err
	}
	attached, err := s.AttachedFiles(ctx, actor, t)
	if err != nil {
		return nil, err
	}
	return s.filePicks(ctx, actor, q, attached)
}

// checkFileRefTarget is [DefaultService.checkFileTarget] for attaching and removing, which needs
// a note, item, expense or event, a message gets its files when it is posted
func (s *DefaultService) checkFileRefTarget(ctx context.Context, actor Membership, t FileTarget) error {
	if !t.NoteID.Valid && !t.ShoppingItemID.Valid && !t.ExpenseID.Valid && !t.EventID.Valid {
		return g.ErrorInvalidInput
	}
	return s.checkFileTarget(ctx, actor, t)
}

// AttachFile shows an existing file on t too, a file that cannot go there is skipped, see
// [dbqueries.Queries.InsertFileRef]; returns [g.ErrorNotFound] if the actor may not see t
func (s *DefaultService) AttachFile(ctx context.Context, actor Membership, t FileTarget, fileID pgtype.UUID) error {
	if err := s.checkFileRefTarget(ctx, actor, t); err != nil {
		return err
	}
	return s.db.InsertFileRef(ctx, dbqueries.InsertFileRefParams{
		FileID:         fileID,
		HouseholdID:    actor.HouseholdID,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
		ExpenseID:      t.ExpenseID,
		MessageID:      t.MessageID,
		EventID:        t.EventID,
	})
}

// KeepFile makes the file stay after its last reference goes, there is no way back;
// returns [g.ErrorNotFound] if the actor may not see the file
func (s *DefaultService) KeepFile(ctx context.Context, actor Membership, fileID pgtype.UUID) error {
	if _, err := s.GetFile(ctx, actor, fileID); err != nil {
		return err
	}
	return rowsAffected(s.db.KeepFile(ctx, dbqueries.KeepFileParams{ID: fileID, HouseholdID: actor.HouseholdID}))
}

// KeepAndDetachFile takes a file off t and leaves it in Files, kept first so that a failure
// cannot lose it; returns [g.ErrorNotFound] if the actor may not see t or the file
func (s *DefaultService) KeepAndDetachFile(ctx context.Context, actor Membership, t FileTarget, fileID pgtype.UUID) error {
	if err := s.checkFileRefTarget(ctx, actor, t); err != nil {
		return err
	}
	if err := s.KeepFile(ctx, actor, fileID); err != nil {
		return err
	}
	return s.DetachFile(ctx, actor, t, fileID)
}

// DetachFile takes a file off t, a trigger deletes it with its last reference unless it is kept;
// returns [g.ErrorNotFound] if the actor may not see t
func (s *DefaultService) DetachFile(ctx context.Context, actor Membership, t FileTarget, fileID pgtype.UUID) error {
	if err := s.checkFileRefTarget(ctx, actor, t); err != nil {
		return err
	}
	return s.db.DeleteFileRef(ctx, dbqueries.DeleteFileRefParams{
		FileID:         fileID,
		HouseholdID:    actor.HouseholdID,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
		ExpenseID:      t.ExpenseID,
		EventID:        t.EventID,
	})
}

// AttachedFiles returns the files of a note, item, expense or event the actor has already loaded
func (s *DefaultService) AttachedFiles(ctx context.Context, actor Membership, t FileTarget) ([]AttachedFile, error) {
	return s.db.ListAttachedFiles(ctx, dbqueries.ListAttachedFilesParams{
		HouseholdID:    actor.HouseholdID,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
		ExpenseID:      t.ExpenseID,
		MessageID:      t.MessageID,
		EventID:        t.EventID,
		UserID:         actor.UserID,
	})
}

// GetFile returns [g.ErrorNotFound] if the actor may not see the file
func (s *DefaultService) GetFile(ctx context.Context, actor Membership, fileID pgtype.UUID) (FileMeta, error) {
	f, err := s.db.GetFile(ctx, dbqueries.GetFileParams{
		ID:          fileID,
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
	return f, noRowsAs(err, g.ErrorNotFound)
}

// OpenFile returns [g.ErrorNotFound] if the actor may not see the file, or with thumb if it has
// no thumbnail; the caller closes the contents
func (s *DefaultService) OpenFile(ctx context.Context, actor Membership, fileID pgtype.UUID, thumb bool) (FileMeta, io.ReadCloser, error) {
	f, err := s.GetFile(ctx, actor, fileID)
	if err != nil {
		return f, nil, err
	}
	key := f.ObjectKey
	if thumb {
		if f.ThumbKey == nil {
			return f, nil, g.ErrorNotFound
		}
		key = *f.ThumbKey
	}
	r, err := s.store.Get(ctx, key)
	return f, r, err
}

// RenameFile gives the file a new name before its extension, which stays, the type is
// the sniffed mime and does not change either
//
// returns [g.ErrorNotFound] unless the actor uploaded the file, [g.ErrorInvalidInput] for a bad name
func (s *DefaultService) RenameFile(ctx context.Context, actor Membership, fileID pgtype.UUID, stem string) error {
	f, err := s.GetFile(ctx, actor, fileID)
	if err != nil {
		return err
	}
	// checked alone too, or a blank stem would leave only the extension
	stem, err = cleanFileName(stem)
	if err != nil {
		return err
	}
	_, ext := SplitFileName(f.Name)
	name, err := cleanFileName(stem + ext)
	if err != nil {
		return err
	}
	return rowsAffected(s.db.RenameFile(ctx, dbqueries.RenameFileParams{
		Name:        name,
		ID:          fileID,
		HouseholdID: actor.HouseholdID,
		UploaderID:  actor.UserID,
	}))
}

// DeleteFile returns [g.ErrorNotFound] unless the actor uploaded the file or is an owner who can see it,
// the object is removed later by [DefaultService.PurgeFiles]
func (s *DefaultService) DeleteFile(ctx context.Context, actor Membership, fileID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		f, err := q.GetFile(ctx, dbqueries.GetFileParams{UserID: actor.UserID, ID: fileID, HouseholdID: actor.HouseholdID})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if err := rowsAffected(q.DeleteFile(ctx, dbqueries.DeleteFileParams{
			ID:          fileID,
			HouseholdID: actor.HouseholdID,
			UploaderID:  actor.UserID,
			IsOwner:     actor.Role == g.RoleOwner,
		})); err != nil {
			return err
		}
		// only a file everyone sees, and not one that lives in the chat alone
		if !f.Shared || (f.MessageID.Valid && !f.NoteID.Valid && !f.ShoppingItemID.Valid && !f.ExpenseID.Valid && !f.EventID.Valid) {
			return nil
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActFileRemoved, SubjectID: fileID, TargetUserID: pgtype.UUID{}, Title: f.Name, AmountCents: nil, Currency: ""})
	})
}

// PurgeFiles removes the objects of deleted files from storage,
// a trigger queues them, so files deleted by a cascade are included
func (s *DefaultService) PurgeFiles(ctx context.Context) error {
	// failed keys stay queued for the next run at the front of the queue, skip past them
	var failed int32
	for {
		keys, err := s.db.ListFilePurges(ctx, dbqueries.ListFilePurgesParams{Max: purgeBatch, Skip: failed})
		if err != nil || len(keys) == 0 {
			return err
		}
		for _, key := range keys {
			if err := s.store.Delete(ctx, key); err != nil {
				log.Error().Err(err).Str("key", key).Msg("could not delete object of a deleted file")
				failed++
				continue
			}
			if err := s.db.DeleteFilePurge(ctx, key); err != nil {
				return err
			}
		}
	}
}

// RunCleanup purges deleted files and expired email tokens every interval until ctx is done
func (s *DefaultService) RunCleanup(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		// before PurgeFiles, so an erased avatar, a purged household's files and
		// expired exports go in the same run
		if err := s.EraseDeletedAccounts(ctx); err != nil {
			log.Error().Err(err).Msg("could not erase deleted accounts")
		}
		if err := s.PurgeDeletedHouseholds(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge deleted households")
		}
		if err := s.PurgeExpiredExports(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge expired data exports")
		}
		if err := s.PurgeFiles(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge deleted files")
		}
		if err := s.db.DeleteExpiredEmailTokens(ctx); err != nil {
			log.Error().Err(err).Msg("could not delete expired email tokens")
		}
		if err := s.PurgeActivity(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge activity")
		}
		if err := s.NotifyRecurringDue(ctx); err != nil {
			log.Error().Err(err).Msg("could not notify due recurring expenses")
		}
		if err := s.NotifyChoresDue(ctx); err != nil {
			log.Error().Err(err).Msg("could not notify due chores")
		}
		if err := s.NotifyCalendarDue(ctx); err != nil {
			log.Error().Err(err).Msg("could not notify due event reminders")
		}
		if err := s.db.PurgeEventReminders(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge event reminders")
		}
		// last, it can take a while
		if err := s.BuildPendingExports(ctx); err != nil {
			log.Error().Err(err).Msg("could not build a data export")
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
