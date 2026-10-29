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

	"github.com/jackc/pgx/v5/pgtype"
)

// longest side of a picture, sharp at 40px on a 3x screen with room to spare
const avatarSide = 256

// SetAvatar makes the image in size bytes of body the user's picture, stored as a small JPEG
// so the original and its EXIF (location) are never kept, the old picture is purged
//
// returns [g.ErrorFileTooLarge] above [MaxFileSize], [g.ErrorInvalidInput] unless it is an image that decodes
func (s *DefaultService) SetAvatar(ctx context.Context, userID pgtype.UUID, body io.ReadSeeker, size int64) error {
	if size > MaxFileSize {
		return g.ErrorFileTooLarge
	}
	head := make([]byte, 512)
	n, err := io.ReadFull(body, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return g.ErrorInvalidInput
	}
	if !ImageMimes[http.DetectContentType(head[:n])] {
		return g.ErrorInvalidInput
	}
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return err
	}
	jpeg, err := images.Thumbnail(body, avatarSide)
	if err != nil {
		return g.ErrorInvalidInput
	}

	key := "users/" + userID.String() + "/" + rand.Text()
	if err := s.store.Put(ctx, key, bytes.NewReader(jpeg), int64(len(jpeg)), "image/jpeg"); err != nil {
		return err
	}
	if err := s.db.SetUserAvatar(ctx, dbqueries.SetUserAvatarParams{ID: userID, AvatarKey: &key}); err != nil {
		// the row does not have it, so nothing would ever purge the object
		if err := s.store.Delete(context.WithoutCancel(ctx), key); err != nil {
			log.Error().Err(err).Str("key", key).Msg("could not delete object of a failed avatar upload")
		}
		return err
	}
	return nil
}

// RemoveAvatar takes the user's picture away, the initial shows again
func (s *DefaultService) RemoveAvatar(ctx context.Context, userID pgtype.UUID) error {
	return s.db.SetUserAvatar(ctx, dbqueries.SetUserAvatarParams{ID: userID, AvatarKey: nil})
}

// OpenAvatar returns the JPEG of userID's picture, [g.ErrorNotFound] if there is none
// or the viewer neither is them nor lives with them; the caller closes it
func (s *DefaultService) OpenAvatar(ctx context.Context, viewerID, userID pgtype.UUID) (io.ReadCloser, error) {
	key, err := s.db.GetUserAvatarKey(ctx, dbqueries.GetUserAvatarKeyParams{UserID: userID, ViewerID: viewerID})
	if err != nil {
		return nil, noRowsAs(err, g.ErrorNotFound)
	}
	return s.store.Get(ctx, key)
}
