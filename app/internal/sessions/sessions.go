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
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	gorillasessions "github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
)

var log = logger.Main.With().Str("package", "sessions").Logger()

const SessionNamePostgreSQL = "session"

// https://gin-gonic.com/en/docs/middleware/session-management/#session-options
//
//exhaustruct:ignore
var sessionOptions = sessions.Options{
	Path:     "/",
	MaxAge:   604800, // 1 week
	HttpOnly: true,
	Secure:   true,
	SameSite: http.SameSiteStrictMode,
}

// set by [NewPostgreSessionStoreHandler], same db as the session store
var sessionDB *sql.DB

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

// Login stores sessionInfo under a new session id, so an id that existed
// before signing in (possibly planted by someone else) never becomes authenticated
func (s *Session) Login(sessionInfo SessionInfo) {
	err := validators.Main.Struct(sessionInfo)
	if err != nil {
		validators.LogValidationErrors(log, err)
		panic(err)
	}
	s.regenerateID()
	s.SessionInfo = sessionInfo
	s.isSessionInfoValid = true

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

// Logout deletes the session from the database and expires its cookie
func (s *Session) Logout() {
	s.Session.Clear()
	expired := sessionOptions
	expired.MaxAge = -1
	s.Session.Options(expired)
	s.Save()
	s.SessionInfo = SessionInfo{"", "", ""}
	s.isSessionInfoValid = false
}

// the gin-contrib session wrapper does not expose the gorilla session in its interface
type gorillaSessionGetter interface {
	Session() *gorillasessions.Session
}

// deletes the stored session (if any) and makes the next [Session.Save]
// create a new row with a new id, see pgstore Save
func (s *Session) regenerateID() {
	getter, ok := s.Session.(gorillaSessionGetter)
	if !ok {
		panic("session does not expose gorilla session, cannot regenerate id")
	}
	gs := getter.Session()
	if gs.ID != "" {
		// same query as pgstore destroy
		if _, err := sessionDB.Exec("DELETE FROM http_sessions WHERE key = $1", gs.ID); err != nil {
			panic(err)
		}
	}
	s.Session.Clear()
	gs.ID = ""
	gs.IsNew = true
}

// -----------------------------------------------------------------------------

// how often expired rows are deleted from http_sessions
const expiredSessionCleanupInterval = time.Hour

// promoted from the [pgstore.PGStore] that the gin-contrib postgres store embeds
//
// [pgstore.PGStore]: https://pkg.go.dev/github.com/antonlindstrom/pgstore#PGStore.Cleanup
type sessionCleaner interface {
	Cleanup(interval time.Duration) (chan<- struct{}, <-chan struct{})
}

// pgstore only deletes a session row on logout, rows of sessions that
// simply expired stay until the cleanup removes them
//
// runs for the lifetime of the app, so the stop channels are not kept;
// errors are logged through the standard log (see [logger.InitLogger])
func startExpiredSessionCleanup(store postgres.Store) {
	cleaner, ok := store.(sessionCleaner)
	if !ok {
		log.Fatal().Msg("session store has no Cleanup, expired sessions would never be deleted")
	}
	cleaner.Cleanup(expiredSessionCleanupInterval)
}

// NewPostgreSessionStoreHandler https://gin-gonic.com/en/docs/middleware/session-management/#_top
func NewPostgreSessionStoreHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	// https://github.com/jackc/pgx/issues/1065#issuecomment-2991745812
	sessionDB = sql.OpenDB(pgxstdlib.GetPoolConnector(dbpool))
	// FIX FOR securecookie: error - caused by:
	//  gob: type not registered for interface: config.GinSessionKey
	// gob.Register(g.GinSessionKey("")) NO LONGER USED will keep here tho

	pairs, err := g.GetSessionKeyPairs()
	if err != nil {
		log.Warn().Err(err).Str("status", "ignored").Msg("error with previous pair")
	}
	// https://github.com/antonlindstrom/pgstore/blob/master/examples/sessions.go#L14
	pgStore, err := postgres.NewStore(sessionDB, pairs...)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing postgres session store")
	}

	pgStore.Options(sessionOptions)
	startExpiredSessionCleanup(pgStore)

	// https://github.com/gin-contrib/sessions#postgresql
	return sessions.Sessions(SessionNamePostgreSQL, pgStore)
}
