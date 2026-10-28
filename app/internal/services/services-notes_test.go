package services

import (
	"context"
	g "roommates/internal/config"
	"strings"
	"testing"
)

func TestNotes(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	var err error

	_, err = s.CreateNote(ctx, ownerM, NoteInput{Title: "  ", BodyMD: "", Private: false})
	wantErr(t, "empty title", err, g.ErrorInvalidInput)
	_, err = s.CreateNote(ctx, ownerM, NoteInput{Title: "x", BodyMD: strings.Repeat("x", MaxNoteBodyLength+1), Private: false})
	wantErr(t, "long body", err, g.ErrorInvalidInput)

	shared, err := s.CreateNote(ctx, ownerM, NoteInput{Title: " Wifi ", BodyMD: "**pass**\n\n", Private: false})
	if err != nil {
		t.Fatal(err)
	}
	private, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Diary", BodyMD: "", Private: true})
	if err != nil {
		t.Fatal(err)
	}

	// the member sees only the shared note, the stranger nothing
	notes, err := s.ListNotes(ctx, memberM)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 || notes[0].ID != shared {
		t.Fatalf("member lists %+v", notes)
	}
	_, err = s.GetNote(ctx, memberM, private)
	wantErr(t, "member reads private note", err, g.ErrorNotFound)
	wantErr(t, "member edits private note", s.UpdateNote(ctx, memberM, private, NoteInput{Title: "x", BodyMD: "", Private: false}), g.ErrorNotFound)
	wantErr(t, "member pins private note", s.SetNotePinned(ctx, memberM, private, true), g.ErrorNotFound)
	_, err = s.GetNote(ctx, strangerM, shared)
	wantErr(t, "stranger reads note", err, g.ErrorNotFound)
	wantErr(t, "stranger edits note", s.UpdateNote(ctx, strangerM, shared, NoteInput{Title: "x", BodyMD: "", Private: false}), g.ErrorNotFound)
	wantErr(t, "stranger deletes note", s.DeleteNote(ctx, strangerM, shared), g.ErrorNotFound)
	wantErr(t, "member deletes owner's note", s.DeleteNote(ctx, memberM, shared), g.ErrorNotFound)

	// the member edits the shared note, but cannot make it private
	if err := s.UpdateNote(ctx, memberM, shared, NoteInput{Title: "Wi-Fi", BodyMD: "new", Private: true}); err != nil {
		t.Fatal(err)
	}
	n, err := s.GetNote(ctx, ownerM, shared)
	if err != nil {
		t.Fatal(err)
	}
	if n.Title != "Wi-Fi" || n.BodyMd != "new" || n.Visibility != g.NoteVisibilityHousehold {
		t.Fatalf("after member edit %+v", n)
	}

	// pinned notes come first
	if err := s.SetNotePinned(ctx, ownerM, private, true); err != nil {
		t.Fatal(err)
	}
	notes, err = s.ListNotes(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 2 || notes[0].ID != private || !notes[0].Pinned {
		t.Fatalf("owner lists %+v", notes)
	}

	if err := s.DeleteNote(ctx, ownerM, shared); err != nil {
		t.Fatal(err)
	}
	_, err = s.GetNote(ctx, ownerM, shared)
	wantErr(t, "deleted note", err, g.ErrorNotFound)
}
