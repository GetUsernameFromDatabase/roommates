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
	var id pgtype.UUID
	err = s.inTx(ctx, func(q *dbqueries.Queries) error {
		var err error
		id, err = q.InsertNote(ctx, dbqueries.InsertNoteParams{
			HouseholdID: actor.HouseholdID,
			AuthorID:    actor.UserID,
			Title:       in.Title,
			BodyMd:      in.BodyMD,
			Visibility:  visibility,
		})
		if err != nil {
			return err
		}
		return logNote(ctx, q, actor, g.ActNoteCreated, id, in.Title, visibility)
	})
	return id, err
}

// logNote logs a note, but not a private one, its title is for its author only
func logNote(ctx context.Context, q *dbqueries.Queries, actor Membership, kind string, id pgtype.UUID, title, visibility string) error {
	if visibility == g.NoteVisibilityPrivate {
		return nil
	}
	return logActivity(ctx, q, actor, ActivityEntry{Kind: kind, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: title, AmountCents: nil, Currency: ""})
}

// ListNotes returns the household's notes and the actor's private ones, only those with q
// in the title or body (in any case) when q is not blank
func (s *DefaultService) ListNotes(ctx context.Context, actor Membership, q string) ([]dbqueries.ListNotesRow, error) {
	return s.db.ListNotes(ctx, dbqueries.ListNotesParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
		Q:           strings.TrimSpace(q),
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
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := rowsAffected(q.UpdateNote(ctx, dbqueries.UpdateNoteParams{
			Title:       in.Title,
			BodyMd:      in.BodyMD,
			UserID:      actor.UserID,
			Visibility:  visibility,
			ID:          noteID,
			HouseholdID: actor.HouseholdID,
		})); err != nil {
			return err
		}
		// a member cannot make a note private, so the note's own visibility is the one that counts
		n, err := q.GetNote(ctx, dbqueries.GetNoteParams{ID: noteID, HouseholdID: actor.HouseholdID, UserID: actor.UserID})
		if err != nil {
			return err
		}
		return logNote(ctx, q, actor, g.ActNoteEdited, noteID, in.Title, n.Visibility)
	})
}

func (s *DefaultService) SetNotePinned(ctx context.Context, actor Membership, noteID pgtype.UUID, pinned bool) error {
	return rowsAffected(s.db.SetNotePinned(ctx, dbqueries.SetNotePinnedParams{
		Pinned:      pinned,
		ID:          noteID,
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	}))
}

// DeleteNote returns [g.ErrorNotFound] unless the actor wrote the note or is an owner and the note is shared
func (s *DefaultService) DeleteNote(ctx context.Context, actor Membership, noteID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		n, err := q.GetNote(ctx, dbqueries.GetNoteParams{ID: noteID, HouseholdID: actor.HouseholdID, UserID: actor.UserID})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if err := rowsAffected(q.DeleteNote(ctx, dbqueries.DeleteNoteParams{
			ID:          noteID,
			HouseholdID: actor.HouseholdID,
			AuthorID:    actor.UserID,
			IsOwner:     actor.Role == g.RoleOwner,
		})); err != nil {
			return err
		}
		return logNote(ctx, q, actor, g.ActNoteDeleted, noteID, n.Title, n.Visibility)
	})
}
