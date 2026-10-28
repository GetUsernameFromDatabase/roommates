package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/images"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxFileSize       = 25 << 20
	MaxFileNameLength = 200
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

// FileTarget is what an upload attaches to: a note, an item of a list, or nothing
type FileTarget struct {
	NoteID         pgtype.UUID
	ListID         pgtype.UUID
	ShoppingItemID pgtype.UUID
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

// checkFileTarget returns [g.ErrorNotFound] unless the actor sees the note or item
func (s *DefaultService) checkFileTarget(ctx context.Context, actor Membership, t FileTarget) error {
	if t.NoteID.Valid {
		_, err := s.GetNote(ctx, actor, t.NoteID)
		return err
	}
	if t.ShoppingItemID.Valid {
		_, err := s.GetShoppingItem(ctx, actor, t.ListID, t.ShoppingItemID)
		return err
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
	})
	if err != nil {
		deleteStored()
	}
	return fileID, err
}

// ListFiles returns the household's files except those on others' private notes
func (s *DefaultService) ListFiles(ctx context.Context, actor Membership) ([]dbqueries.ListFilesRow, error) {
	return s.db.ListFiles(ctx, dbqueries.ListFilesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
}

// AttachedFiles returns the files of a note or item the actor has already loaded
func (s *DefaultService) AttachedFiles(ctx context.Context, actor Membership, t FileTarget) ([]AttachedFile, error) {
	return s.db.ListAttachedFiles(ctx, dbqueries.ListAttachedFilesParams{
		HouseholdID:    actor.HouseholdID,
		NoteID:         t.NoteID,
		ShoppingItemID: t.ShoppingItemID,
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
