package services

import (
	"bytes"
	"image"
	"io"
	g "roommates/internal/config"
	"testing"
)

func TestAvatars(t *testing.T) {
	s := newTestService(t)
	store := s.store.(*fakeStore)
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	set := func(body []byte) error {
		return s.SetAvatar(ctx, memberM.UserID, bytes.NewReader(body), int64(len(body)))
	}

	wantErr(t, "not an image", set([]byte("hello")), g.ErrorInvalidInput)
	wantErr(t, "image that does not decode", set([]byte("\x89PNG\r\n\x1a\n rest")), g.ErrorInvalidInput)
	_, err := s.OpenAvatar(ctx, memberM.UserID, memberM.UserID)
	wantErr(t, "no picture yet", err, g.ErrorNotFound)

	if err := set(realPNG(t)); err != nil {
		t.Fatal(err)
	}
	first, err := s.GetUser(ctx, memberM.UserID)
	if err != nil || first.AvatarKey == nil {
		t.Fatalf("user %+v, %v", first, err)
	}
	for _, viewer := range []Membership{memberM, ownerM} {
		body, err := s.OpenAvatar(ctx, viewer.UserID, memberM.UserID)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(body)
		if _, format, err := image.DecodeConfig(bytes.NewReader(b)); err != nil || format != "jpeg" {
			t.Fatalf("picture is %q, %v", format, err)
		}
	}
	_, err = s.OpenAvatar(ctx, strangerM.UserID, memberM.UserID)
	wantErr(t, "stranger opens picture", err, g.ErrorNotFound)

	// a new picture and removing it each purge the one before
	if err := set(realPNG(t)); err != nil {
		t.Fatal(err)
	}
	second, err := s.GetUser(ctx, memberM.UserID)
	if err != nil || second.AvatarKey == nil || *second.AvatarKey == *first.AvatarKey {
		t.Fatalf("user %+v, %v", second, err)
	}
	if err := s.RemoveAvatar(ctx, memberM.UserID); err != nil {
		t.Fatal(err)
	}
	if err := s.PurgeFiles(ctx); err != nil {
		t.Fatal(err)
	}
	if store.has(*first.AvatarKey) || store.has(*second.AvatarKey) {
		t.Fatal("old pictures were not purged")
	}
	_, err = s.OpenAvatar(ctx, ownerM.UserID, memberM.UserID)
	wantErr(t, "removed picture", err, g.ErrorNotFound)
}
