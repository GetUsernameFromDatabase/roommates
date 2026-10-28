package services

import (
	"errors"
	"roommates/internal/db/dbqueries"
	g "roommates/internal/config"
	"roommates/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (s *DefaultService) RegisterUser(ctx *gin.Context, user dbqueries.InsertUserParams) (*pgtype.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil { // this should never be an issue
		log.Error().Err(err).Msg("error processing password")
		return nil, err
	}

	returnValue, err := s.db.InsertUser(ctx, dbqueries.InsertUserParams{
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		switch err := err.(type) {
		case *pgconn.PgError:
			// https://www.postgresql.org/docs/current/errcodes-appendix.html
			if err.Code == "23505" {
				return nil, g.ErrorAccountAlreadyExists
			}
		}
		return nil, err
	}

	return &returnValue, nil
}

func (s *DefaultService) ShouldUserBeSignedIn(ctx *gin.Context, req models.LoginRequest) (*dbqueries.GetUserCredentialsRow, error) {
	credsInDB, err := s.db.GetUserCredentials(ctx, req.Email)
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
