package services

import (
	"errors"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// NormalizeEmail emails are stored lowercase (enforced by the users table)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ~250ms per hash, see the auth rate limit in routes.go
const bcryptCost = 12

func (s *DefaultService) RegisterUser(ctx *gin.Context, user dbqueries.InsertUserParams) (*pgtype.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil { // only happens if password is longer than 72 bytes
		log.Error().Err(err).Msg("error processing password")
		return nil, err
	}

	returnValue, err := s.db.InsertUser(ctx, dbqueries.InsertUserParams{
		Email:    NormalizeEmail(user.Email),
		Username: user.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, g.ErrorAccountAlreadyExists
		}
		return nil, err
	}

	return &returnValue, nil
}

func (s *DefaultService) ShouldUserBeSignedIn(ctx *gin.Context, req models.LoginRequest) (*dbqueries.GetUserCredentialsRow, error) {
	credsInDB, err := s.db.GetUserCredentials(ctx, NormalizeEmail(req.Email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, g.ErrorInvalidCredential
		}
		return nil, err
	}

	// no need to check email as it's used to get stored credentials
	if err = bcrypt.CompareHashAndPassword([]byte(credsInDB.Password),
		[]byte(req.Password)); err != nil {
		return nil, g.ErrorInvalidCredential
	}

	return &credsInDB, nil
}
