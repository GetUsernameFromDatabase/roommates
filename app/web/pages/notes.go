package pages

import (
	"roommates/internal/services"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// NoteForm is the note editor, a zero NoteID is a new note
//
// IsAuthor shows the private switch, only the author changes who sees a note
type NoteForm struct {
	NoteID       pgtype.UUID
	Input        services.NoteInput
	IsAuthor     bool
	ErrorMessage string
}

var noteEditorScript = templ.NewOnceHandle()
