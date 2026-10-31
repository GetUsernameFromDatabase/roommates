package services

import (
	"context"
	"errors"
	"io"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/utils"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxMessageLength = 4000
	// MaxMessageFiles is how many files one message can carry
	MaxMessageFiles = 10
	// ChatPageSize is how many messages the chat shows, older ones are a link away
	ChatPageSize = 50
)

// Upload is a file sent with a form
type Upload struct {
	Name string
	Body io.ReadSeeker
	Size int64
}

type MessageFile = dbqueries.ListMessageFilesRow

// Message is a chat message with its files
type Message struct {
	dbqueries.ListMessagesRow
	Files []MessageFile
	// LastFiles is how many files go with the message when it is deleted
	LastFiles int
}

// Chat is a page of messages, newest first
type Chat struct {
	Messages []Message
	// HasOlder is true when messages older than the last one exist
	HasOlder bool
}

// ReplyPreview is the quote of a message being answered
type ReplyPreview = dbqueries.GetReplyPreviewRow

// GetReplyPreview returns [g.ErrorNotFound] if the message is not in the actor's household
func (s *DefaultService) GetReplyPreview(ctx context.Context, actor Membership, messageID pgtype.UUID) (ReplyPreview, error) {
	p, err := s.db.GetReplyPreview(ctx, dbqueries.GetReplyPreviewParams{ID: messageID, HouseholdID: actor.HouseholdID})
	return p, noRowsAs(err, g.ErrorNotFound)
}

// PostMessage sends body with up to [MaxMessageFiles] new files and references to existing ones
// as an answer to replyTo when it is valid,
// [g.ErrorInvalidInput] when body, files and refs are all missing, there are too many of them,
// body is longer than [MaxMessageLength] or a reference is to a file the household cannot see,
// [g.ErrorNotFound] when replyTo is not (or no longer) in the actor's household,
// file errors are those of [DefaultService.UploadFile]
func (s *DefaultService) PostMessage(ctx context.Context, actor Membership, body string, replyTo pgtype.UUID, files []Upload, refs []pgtype.UUID) error {
	// forms send CRLF but maxlength counts a line break once
	body = strings.TrimSpace(strings.ReplaceAll(body, "\r\n", "\n"))
	if (body == "" && len(files)+len(refs) == 0) || len(files)+len(refs) > MaxMessageFiles || utf8.RuneCountInString(body) > MaxMessageLength {
		return g.ErrorInvalidInput
	}
	var messageID pgtype.UUID
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		var err error
		messageID, err = q.InsertMessage(ctx, dbqueries.InsertMessageParams{
			HouseholdID: actor.HouseholdID,
			AuthorID:    actor.UserID,
			Body:        body,
			ReplyToID:   replyTo,
		})
		if err != nil || len(refs) == 0 {
			return err
		}
		n, err := q.InsertMessageFileRefs(ctx, dbqueries.InsertMessageFileRefsParams{
			MessageID:   messageID,
			FileIds:     refs,
			HouseholdID: actor.HouseholdID,
		})
		if err == nil && n != int64(len(refs)) {
			return g.ErrorInvalidInput
		}
		return err
	})
	if err != nil {
		var pgErr *pgconn.PgError
		// the foreign key keeps replies in the household, the answered message may also just be deleted
		if replyTo.Valid && errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return g.ErrorNotFound
		}
		return err
	}
	for _, file := range files {
		t := FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: messageID}
		if _, err := s.UploadFile(ctx, actor, t, file.Name, file.Body, file.Size); err != nil {
			// the form shows the text again, sending it twice would duplicate it,
			// the files already stored go with the message
			if errDelete := s.DeleteMessage(context.WithoutCancel(ctx), actor, messageID); errDelete != nil {
				log.Error().Err(errDelete).Msg("could not delete message of a failed upload")
			}
			return err
		}
	}
	s.notifyChat(ctx, actor, body)
	return nil
}

// notifyChat tells the other members about a new message
func (s *DefaultService) notifyChat(ctx context.Context, actor Membership, body string) {
	members, err := s.db.ListActiveMembers(ctx, actor.HouseholdID)
	if err != nil {
		// the message is sent, only the notifications are missing
		log.Error().Err(err).Msg("could not list members for chat push")
		return
	}
	ids := make([]pgtype.UUID, 0, len(members))
	for _, m := range members {
		if m.UserID != actor.UserID {
			ids = append(ids, m.UserID)
		}
	}
	s.notifyAsync(ctx, actor, ids, func(ctx context.Context, _ pgtype.UUID, author, household string) push.Message {
		text := pushPreview(body)
		if text == "" {
			text = utils.T(ctx, locales.LK_PushChatFile)
		}
		return push.Message{
			Title: utils.T(ctx, locales.LK_PushChatTitle, author, household),
			Body:  text,
			URL:   householdPath(actor.HouseholdID, "chat"),
			// a new message replaces the notification of the previous one
			Tag: "chat-" + actor.HouseholdID.String(),
		}
	})
}

// ListChat returns [ChatPageSize] messages older than before, or the newest when before
// is not valid, and then also marks the chat read
func (s *DefaultService) ListChat(ctx context.Context, actor Membership, before pgtype.UUID) (Chat, error) {
	rows, err := s.db.ListMessages(ctx, dbqueries.ListMessagesParams{
		HouseholdID: actor.HouseholdID,
		Before:      before,
		MaxCount:    ChatPageSize + 1,
	})
	if err != nil {
		return Chat{}, err
	}
	c := Chat{Messages: make([]Message, 0, len(rows)), HasOlder: len(rows) > ChatPageSize}
	ids := make([]pgtype.UUID, 0, len(rows))
	for _, r := range rows[:min(len(rows), ChatPageSize)] {
		c.Messages = append(c.Messages, Message{ListMessagesRow: r, Files: nil, LastFiles: 0})
		ids = append(ids, r.ID)
	}
	files, err := s.db.ListMessageFiles(ctx, dbqueries.ListMessageFilesParams{
		HouseholdID: actor.HouseholdID,
		MessageIds:  ids,
		UserID:      actor.UserID,
	})
	if err != nil {
		return Chat{}, err
	}
	byMessage := make(map[pgtype.UUID][]MessageFile, len(files))
	for _, f := range files {
		byMessage[f.MessageID] = append(byMessage[f.MessageID], f)
	}
	for i := range c.Messages {
		c.Messages[i].Files = byMessage[c.Messages[i].ID]
		for _, f := range c.Messages[i].Files {
			if f.LastRef {
				c.Messages[i].LastFiles++
			}
		}
	}
	if !before.Valid {
		if err := s.db.MarkChatRead(ctx, dbqueries.MarkChatReadParams{HouseholdID: actor.HouseholdID, UserID: actor.UserID}); err != nil {
			return Chat{}, err
		}
		err = s.db.NotifyChatRead(ctx, actor.HouseholdID)
	}
	return c, err
}

// DeleteMessage returns [g.ErrorNotFound] unless the actor wrote the message
func (s *DefaultService) DeleteMessage(ctx context.Context, actor Membership, messageID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteMessage(ctx, dbqueries.DeleteMessageParams{
		ID:          messageID,
		HouseholdID: actor.HouseholdID,
		AuthorID:    actor.UserID,
	}))
}
