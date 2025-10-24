package models

import (
	"roommates/db/dbqueries"
	l "roommates/locales"
	"roommates/utils"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type Note struct {
	ModelBase
	// not visual
	ID int32 `form:"id"`
	// mostly unused, id from the uri is prioritised
	HouseID string `form:"house_id"`
	// taken from authenticated user when made
	MakerID pgtype.UUID

	Title   string `form:"title"`
	Content string `form:"content"`
	// not stored
	HouseName string `form:"house_name"`
}

func NewNote(note dbqueries.SelectNoteRow) Note {
	return Note{
		ModelBase: ModelBase{Initial: true},
		ID:        note.NoteID,
		HouseID:   note.HouseID.String(),
		MakerID:   note.MakerID,
		Title:     note.Title,
		Content:   note.Content,
		HouseName: note.HouseName,
	}
}

func NewNoteOnlyHouse(house dbqueries.House) Note {
	return Note{
		ModelBase: ModelBase{Initial: true},
		HouseID:   house.ID.String(),
		HouseName: house.Name,
	}
}

func (m *Note) ValidateTitle() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Title == "" {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsNameErrorEmpty})
		return msgs
	}

	charProblems := utils.ValidateString(m.Title, utils.RuneValidationRules{
		LettersAllowed:       true,
		DigitsAllowed:        true,
		MaxConsecutiveSpaces: 1,
	})
	msgs = append(msgs, StringValidationMessages(charProblems)...)
	return msgs
}

func (m *Note) ValidateContent() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Content == "" {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsContentErrorEmpty})
		return msgs
	}

	return msgs
}

func (m *Note) GetIDString() string {
	return strconv.Itoa(int(m.ID))
}

func (m *Note) GetValidators() []Validator {
	return []Validator{
		m.ValidateTitle,
	}
}

func (m *Note) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// checks if the form is valid and sets the Initial to false
func (m *Note) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
