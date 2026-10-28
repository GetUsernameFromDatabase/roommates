package services

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxNoteTitleLength = 100
	MaxNoteBodyLength  = 20000
)

type Note = dbqueries.GetNoteRow

// NoteInput is what the note form sends, Private is ignored when someone else edits the note
type NoteInput struct {
	Title   string
	BodyMD  string
	Private bool
}

// clean returns [g.ErrorInvalidInput] unless the title has 1..[MaxNoteTitleLength] runes
// and the body at most [MaxNoteBodyLength]
func (in NoteInput) clean() (NoteInput, string, error) {
	title, ok := trimToLength(in.Title, 1, MaxNoteTitleLength)
	// forms send CRLF but maxlength counts a line break once, so count and store LF;
	// trailing spaces end a markdown line break, only drop trailing newlines
	body := strings.TrimRight(strings.ReplaceAll(in.BodyMD, "\r\n", "\n"), "\n")
	if !ok || utf8.RuneCountInString(body) > MaxNoteBodyLength {
		return in, "", g.ErrorInvalidInput
	}
	visibility := g.NoteVisibilityHousehold
	if in.Private {
		visibility = g.NoteVisibilityPrivate
	}
	return NoteInput{Title: title, BodyMD: body, Private: in.Private}, visibility, nil
}

func (s *DefaultService) CreateNote(ctx context.Context, actor Membership, in NoteInput) (pgtype.UUID, error) {
	in, visibility, err := in.clean()
	if err != nil {
		return pgtype.UUID{}, err
	}
	return s.db.InsertNote(ctx, dbqueries.InsertNoteParams{
		HouseholdID: actor.HouseholdID,
		AuthorID:    actor.UserID,
		Title:       in.Title,
		BodyMd:      in.BodyMD,
		Visibility:  visibility,
	})
}

// ListNotes returns the household's notes and the actor's private ones
func (s *DefaultService) ListNotes(ctx context.Context, actor Membership) ([]dbqueries.ListNotesRow, error) {
	return s.db.ListNotes(ctx, dbqueries.ListNotesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
}

// GetNote returns [g.ErrorNotFound] if the note is not in the actor's household or is someone else's private one
func (s *DefaultService) GetNote(ctx context.Context, actor Membership, noteID pgtype.UUID) (Note, error) {
	n, err := s.db.GetNote(ctx, dbqueries.GetNoteParams{
		ID:          noteID,
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
	return n, noRowsAs(err, g.ErrorNotFound)
}

// UpdateNote overwrites the note, last write wins
func (s *DefaultService) UpdateNote(ctx context.Context, actor Membership, noteID pgtype.UUID, in NoteInput) error {
	in, visibility, err := in.clean()
	if err != nil {
		return err
	}
	return rowsAffected(s.db.UpdateNote(ctx, dbqueries.UpdateNoteParams{
		Title:       in.Title,
		BodyMd:      in.BodyMD,
		UserID:      actor.UserID,
		Visibility:  visibility,
		ID:          noteID,
		HouseholdID: actor.HouseholdID,
	}))
}

func (s *DefaultService) SetNotePinned(ctx context.Context, actor Membership, noteID pgtype.UUID, pinned bool) error {
	return rowsAffected(s.db.SetNotePinned(ctx, dbqueries.SetNotePinnedParams{
		Pinned:      pinned,
		ID:          noteID,
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	}))
}

// DeleteNote returns [g.ErrorNotFound] unless the actor wrote the note
func (s *DefaultService) DeleteNote(ctx context.Context, actor Membership, noteID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteNote(ctx, dbqueries.DeleteNoteParams{
		ID:          noteID,
		HouseholdID: actor.HouseholdID,
		AuthorID:    actor.UserID,
	}))
}
