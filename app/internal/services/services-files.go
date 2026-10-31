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
	MaxFileSize       = 25 << 20
	MaxFileNameLength = 200
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
	FilePlaces = []string{"note", "item", "expense", "message", "none"}
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
// a chat message, or nothing
type FileTarget struct {
	NoteID         pgtype.UUID
	ListID         pgtype.UUID
	ShoppingItemID pgtype.UUID
	ExpenseID      pgtype.UUID
	MessageID      pgtype.UUID
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

// checkFileTarget returns [g.ErrorNotFound] unless the actor sees the note, item or expense,
// or wrote the message
func (s *DefaultService) checkFileTarget(ctx context.Context, actor Membership, t FileTarget) error {
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
	if size > MaxFileSize {
		return pgtype.UUID{}, g.ErrorFileTooLarge
	}
	name, err := cleanFileName(name)
	if err != nil || size <= 0 {
		return pgtype.UUID{}, g.ErrorInvalidInput
	}
	if err := s.checkFileTarget(ctx, actor, t); err != nil {
		return pgtype.UUID{}, err
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(body, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return pgtype.UUID{}, err
	}
	mime := http.DetectContentType(head[:n])
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return pgtype.UUID{}, err
	}

	var thumb []byte
	if ImageMimes[mime] {
		// a broken image is still uploaded, only without a thumbnail
		if thumb, err = images.Thumbnail(body, thumbnailSide); err != nil {
			log.Warn().Err(err).Str("mime", mime).Msg("could not make a thumbnail")
		}
		if _, err := body.Seek(0, io.SeekStart); err != nil {
			return pgtype.UUID{}, err
		}
	}

	key := "households/" + actor.HouseholdID.String() + "/" + rand.Text()
	stored := []string{}
	// the row is missing, so nothing would ever purge the objects
	deleteStored := func() {
		for _, k := range stored {
			if err := s.store.Delete(context.WithoutCancel(ctx), k); err != nil {
				log.Error().Err(err).Str("key", k).Msg("could not delete object of a failed upload")
			}
		}
	}
	if err := s.store.Put(ctx, key, body, size, mime); err != nil {
		return pgtype.UUID{}, err
	}
	stored = append(stored, key)
	var thumbKey *string
	if thumb != nil {
		k := key + "-thumb"
		if err := s.store.Put(ctx, k, bytes.NewReader(thumb), int64(len(thumb)), "image/jpeg"); err != nil {
			deleteStored()
			return pgtype.UUID{}, err
		}
		stored = append(stored, k)
		thumbKey = &k
	}
	fileID, err := s.db.InsertFile(ctx, dbqueries.InsertFileParams{
		HouseholdID:    actor.HouseholdID,
		UploaderID:     actor.UserID,
		ObjectKey:      key,
		ThumbKey:       thumbKey,
		Name:           name,
		Mime:           mime,
		SizeBytes:      size,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
		ExpenseID:      t.ExpenseID,
		MessageID:      t.MessageID,
	})
	if err != nil {
		deleteStored()
	}
	return fileID, err
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
	rows, err := s.db.ListFiles(ctx, dbqueries.ListFilesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
		SharedOnly:  true,
		Q:           strings.TrimSpace(q),
		Kind:        "",
		Place:       "",
	})
	return rows[:min(len(rows), MaxReferablePicks)], err
}

// AttachedFiles returns the files of a note, item or expense the actor has already loaded
func (s *DefaultService) AttachedFiles(ctx context.Context, actor Membership, t FileTarget) ([]AttachedFile, error) {
	return s.db.ListAttachedFiles(ctx, dbqueries.ListAttachedFilesParams{
		HouseholdID:    actor.HouseholdID,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
		ExpenseID:      t.ExpenseID,
		MessageID:      t.MessageID,
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

// DeleteFile returns [g.ErrorNotFound] unless the actor uploaded the file,
// the object is removed later by [DefaultService.PurgeFiles]
func (s *DefaultService) DeleteFile(ctx context.Context, actor Membership, fileID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteFile(ctx, dbqueries.DeleteFileParams{
		ID:          fileID,
		HouseholdID: actor.HouseholdID,
		UploaderID:  actor.UserID,
	}))
}

// PurgeFiles removes the objects of deleted files from storage,
// a trigger queues them, so files deleted by a cascade are included
func (s *DefaultService) PurgeFiles(ctx context.Context) error {
	for {
		keys, err := s.db.ListFilePurges(ctx, 100)
		if err != nil || len(keys) == 0 {
			return err
		}
		for _, key := range keys {
			if err := s.store.Delete(ctx, key); err != nil {
				return err
			}
			if err := s.db.DeleteFilePurge(ctx, key); err != nil {
				return err
			}
		}
	}
}

// RunFilePurger calls [DefaultService.PurgeFiles] every interval until ctx is done
func (s *DefaultService) RunFilePurger(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if err := s.PurgeFiles(ctx); err != nil {
			log.Error().Err(err).Msg("could not purge deleted files")
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
