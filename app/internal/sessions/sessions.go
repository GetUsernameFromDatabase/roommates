// Package sessions makes handling sessions easier
// https://github.com/gin-contrib/sessions
package sessions

//lint:file-ignore ST1003 goal is to ignore underscore issues

import (
	"database/sql"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/logger"
	"roommates/internal/validators"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
)

var log = logger.Main.With().Str("package", "sessions").Logger()

const SessionNamePostgreSQL = "session"

const (
	KeyUserID   = "user_id"
	KeyUsername = "username"
	KeyEmail    = "email"
)

// -----------------------------------------------------------------------------

// SessionInfo values all [Session] should have
//
//exhaustruct:enforce
type SessionInfo struct {
	UserID   string `validate:"required"`
	Email    string `validate:"required"`
	Username string `validate:"required"`
}

type Session struct {
	// exposing [sessions.Session]
	Session sessions.Session
	SessionInfo
	// will also be used to check if is logged in
	isSessionInfoValid bool
}

// FromGinContext gets [Session] from [gin.Context]
//

func FromGinContext(ctx *gin.Context) Session {
	s := sessions.Default(ctx)
	str := func(key string) string { v, _ := s.Get(key).(string); return v }

	sessionInfo := SessionInfo{
		UserID:   str(KeyUserID),
		Email:    str(KeyEmail),
		Username: str(KeyUsername),
	}

	err := validators.Main.Struct(sessionInfo)
	return Session{
		Session:            s,
		SessionInfo:        sessionInfo,
		isSessionInfoValid: err == nil,
	}
}

func (s *Session) Login(sessionInfo SessionInfo) {
	err := validators.Main.Struct(sessionInfo)
	if err != nil {
		validators.LogValidationErrors(log, err)
		panic(err)
	}
	s.SessionInfo = sessionInfo

	// I don't like this forces a check that all
	// fields from [SessionInfo] get used
	var _ = SessionInfo{"", "", ""}
	s.Session.Set(KeyUserID, sessionInfo.UserID)
	s.Session.Set(KeyUsername, sessionInfo.Username)
	s.Session.Set(KeyEmail, sessionInfo.Email)
	s.Save()
}

// Save is wrapper of [sessions.Session.Save] | [sessions.Session].Save
//
//	the latter is if vscode "editor.links"==true
//
// panics on it's err (intended)
func (s *Session) Save() {
	err := s.Session.Save()
	if err != nil {
		panic(err)
	}
}

func (s *Session) IsLoggedIn() bool {
	return s.isSessionInfoValid
}

// -----------------------------------------------------------------------------

// NewPostgreSessionStoreHandler https://gin-gonic.com/en/docs/middleware/session-management/#_top
func NewPostgreSessionStoreHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	// https://github.com/jackc/pgx/issues/1065#issuecomment-2991745812
	sqlDB := sql.OpenDB(pgxstdlib.GetPoolConnector(dbpool))
	// FIX FOR securecookie: error - caused by:
	//  gob: type not registered for interface: config.GinSessionKey
	// gob.Register(g.GinSessionKey("")) NO LONGER USED will keep here tho

	pairs, err := g.GetSessionKeyPairs()
	if err != nil {
		log.Warn().Err(err).Str("status", "ignored").Msg("error with previous pair")
	}
	// https://github.com/antonlindstrom/pgstore/blob/master/examples/sessions.go#L14
	pgStore, err := postgres.NewStore(sqlDB, pairs...)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing postgres session store")
	}

	// https://gin-gonic.com/en/docs/middleware/session-management/#session-options
	//exhaustruct:ignore
	pgStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   604800, // 1 week
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// https://github.com/gin-contrib/sessions#postgresql
	return sessions.Sessions(SessionNamePostgreSQL, pgStore)
}
