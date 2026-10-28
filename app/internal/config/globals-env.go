package config

//lint:file-ignore ST1003 goal is to ignore underscore issues

import (
	"net"
	"net/url"
	"os"
	"roommates/internal/validators"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EnvironmentVariableKey string

// REQUIRED env key for postgreSQL database
const (
	EVK_DatabaseHost     EnvironmentVariableKey = "DATABASE_HOST"
	EVK_DatabasePort     EnvironmentVariableKey = "DATABASE_PORT"
	EVK_DatabaseName     EnvironmentVariableKey = "POSTGRES_DB"
	EVK_DatabaseUser     EnvironmentVariableKey = "POSTGRES_USER"
	EVK_DatabasePassword EnvironmentVariableKey = "POSTGRES_PASSWORD"
)

// REQUIRED env key for server
const (
	EVK_ServerPort EnvironmentVariableKey = "SERVER_PORT"

	// hash key https://github.com/gorilla/securecookie#examples
	// generate one with
	//  - openssl rand -hex 32
	//  - securecookie.GenerateRandomKey(64)
	EVK_SessionHashKey EnvironmentVariableKey = "SESSION_HASH_KEY"
	// block key https://github.com/gorilla/securecookie#examples
	// generate one with
	//  - openssl rand -hex 16
	//  - securecookie.GenerateRandomKey(32)
	EVK_SessionBlockKey EnvironmentVariableKey = "SESSION_BLOCK_KEY"
)

// env key for server
const (
	// to rotate out with [EVK_SessionHashKey]
	EVK_SessionHashKeyPrevious EnvironmentVariableKey = "SESSION_HASH_KEY_PREVIOUS"
	// to rotate out with [EVK_SessionBlockKey]
	EVK_SessionBlockKeyPrevious EnvironmentVariableKey = "SESSION_BLOCK_KEY_PREVIOUS"
)

// env keys for email, required in release mode, see [GetAppBaseURL]
const (
	// links in emails are built from this, never from the request Host header
	// (it could point a password reset link at another site), e.g. https://example.com
	EVK_AppBaseURL EnvironmentVariableKey = "APP_BASE_URL"
	// if unset in development, emails are written to the log instead
	EVK_SMTPHost EnvironmentVariableKey = "SMTP_HOST"
	// STARTTLS port (587), implicit TLS (465) is not supported
	EVK_SMTPPort     EnvironmentVariableKey = "SMTP_PORT"
	EVK_SMTPUsername EnvironmentVariableKey = "SMTP_USERNAME"
	EVK_SMTPPassword EnvironmentVariableKey = "SMTP_PASSWORD"
	// e.g. Roommates <noreply@example.com>
	EVK_SMTPFrom EnvironmentVariableKey = "SMTP_FROM"
)

// GetAppBaseURL is [EVK_AppBaseURL] without a trailing slash,
// in development it defaults to https://localhost:[EVK_ServerPort]
func GetAppBaseURL() string {
	if IsDevelopment() {
		return strings.TrimSuffix(EVK_AppBaseURL.Get("https://localhost:"+EVK_ServerPort.Get("8080")), "/")
	}
	return strings.TrimSuffix(EVK_AppBaseURL.MustGet(), "/")
}

func (c EnvironmentVariableKey) String() string {
	return string(c)
}

// Get is a wrapper of [GetEnv]
func (c EnvironmentVariableKey) Get(fallback string) string {
	return GetEnv(c, fallback)
}

// MustGet is a wrapper of [MustGetEnv]
func (c EnvironmentVariableKey) MustGet() string {
	return MustGetEnv(c)
}

// -----------------------------------------------------------------------------

// IsDevelopment checks if [gin.Mode] == [gin.DebugMode]
func IsDevelopment() bool {
	return gin.Mode() == gin.DebugMode
}

func GetEnv(key EnvironmentVariableKey, fallback string) string {
	if value, ok := os.LookupEnv(key.String()); ok {
		return value
	}
	return fallback
}

// MustGetEnv will get env value, panics if it does not exist
func MustGetEnv(key EnvironmentVariableKey) string {
	keyAsString := key.String()
	if value, ok := os.LookupEnv(keyAsString); ok {
		return value
	}

	panic("environment variable " + strconv.Quote(keyAsString) + " needs to be set")
}

// GetDatabaseURL constructs a database url using environment
//
//	sslmode=disable
//
// a URL so that values with spaces, quotes etc. are escaped
func GetDatabaseURL() string {
	//exhaustruct:ignore
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			MustGetEnv(EVK_DatabaseUser),
			MustGetEnv(EVK_DatabasePassword),
		),
		Host:     net.JoinHostPort(MustGetEnv(EVK_DatabaseHost), MustGetEnv(EVK_DatabasePort)),
		Path:     MustGetEnv(EVK_DatabaseName),
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

// GetSessionKeyPairs must get current hash/block pair
// panics on any error regarding them
//
// if there are any errors regarding prev hash/block pair
// returns only current pair and error of previous pair
// otherwise returns both pairs
//
// see [utils.ValidateCookieHashBlockPair]
func GetSessionKeyPairs() ([][]byte, error) {
	hash := []byte(MustGetEnv(EVK_SessionHashKey))
	block := []byte(MustGetEnv(EVK_SessionBlockKey))
	err := validators.ValidateCookieHashBlockPair(hash, block)
	if err != nil {
		panic(err)
	}
	pairs := [][]byte{hash, block}

	prevHash := []byte(EVK_SessionHashKeyPrevious.Get(""))
	prevBlock := []byte(EVK_SessionBlockKeyPrevious.Get(""))
	// previous pair is optional, only needed while rotating keys
	if len(prevHash) == 0 && len(prevBlock) == 0 {
		return pairs, nil
	}
	err = validators.ValidateCookieHashBlockPair(prevHash, prevBlock)
	if err != nil {
		return pairs, err
	}
	pairs = append(pairs, prevHash, prevBlock)
	return pairs, nil
}
