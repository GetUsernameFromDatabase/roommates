package config

//lint:file-ignore ST1003 goal is to ignore underscore issues

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"roommates/internal/validators"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// EnvironmentVariableKey is an env key the app reads. An empty value counts as unset,
// like compose's ${KEY:-default}. The app container gets only the keys listed under
// app.environment in docker-compose.yml, so add a new key there too
type EnvironmentVariableKey string

// REQUIRED env key for postgreSQL database
const (
	EVK_DatabaseHost     EnvironmentVariableKey = "DATABASE_HOST"
	EVK_DatabasePort     EnvironmentVariableKey = "DATABASE_PORT"
	EVK_DatabaseName     EnvironmentVariableKey = "POSTGRES_DB"
	EVK_DatabaseUser     EnvironmentVariableKey = "POSTGRES_USER"
	EVK_DatabasePassword EnvironmentVariableKey = "POSTGRES_PASSWORD"
	// pgxpool pool_max_conns, defaults to 10
	EVK_DatabaseMaxConns EnvironmentVariableKey = "DATABASE_MAX_CONNS"
)

// REQUIRED env key for server
const (
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
	// defaults to 8080, the port inside the app container; `task dev` sets 8090
	EVK_ServerPort EnvironmentVariableKey = "SERVER_PORT"
	// to rotate out with [EVK_SessionHashKey]
	EVK_SessionHashKeyPrevious EnvironmentVariableKey = "SESSION_HASH_KEY_PREVIOUS"
	// to rotate out with [EVK_SessionBlockKey]
	EVK_SessionBlockKeyPrevious EnvironmentVariableKey = "SESSION_BLOCK_KEY_PREVIOUS"
	// cloudflare or flyio: TLS ends at that platform, the server speaks plain
	// HTTP and takes the client IP from its header; unset serves HTTPS with
	// ./certificates. Only safe while the server port is not reachable
	// publicly, since anyone could forge the header
	EVK_TrustedPlatform EnvironmentVariableKey = "TRUSTED_PLATFORM"
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

// REQUIRED env keys for file storage, an S3 API (RustFS in docker compose)
const (
	// e.g. http://127.0.0.1:9000, buckets are addressed path-style
	EVK_S3Endpoint EnvironmentVariableKey = "S3_ENDPOINT"
	EVK_S3Bucket   EnvironmentVariableKey = "S3_BUCKET"
	// the app's own key, limited to the bucket, made by deploy/rustfs/app-key.sh;
	// the root keys S3_ACCESS_KEY / S3_SECRET_KEY are for rustfs and backup only
	EVK_AppS3AccessKey EnvironmentVariableKey = "APP_S3_ACCESS_KEY"
	EVK_AppS3SecretKey EnvironmentVariableKey = "APP_S3_SECRET_KEY"
	// defaults to us-east-1, RustFS ignores it
	EVK_S3Region EnvironmentVariableKey = "S3_REGION"
)

// env keys for web push, the keys are required in release, see [push.New]
const (
	// generate both with `go run ./tools/vapid-keys`
	EVK_VAPIDPublicKey  EnvironmentVariableKey = "VAPID_PUBLIC_KEY"
	EVK_VAPIDPrivateKey EnvironmentVariableKey = "VAPID_PRIVATE_KEY"
	// how push services reach the operator, an email address or https URL,
	// defaults to [EVK_AppBaseURL]
	EVK_VAPIDSubject EnvironmentVariableKey = "VAPID_SUBJECT"
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

// GetEnv returns the env value, or fallback when it is unset or empty
func GetEnv(key EnvironmentVariableKey, fallback string) string {
	if value := os.Getenv(key.String()); value != "" {
		return value
	}
	return fallback
}

// MustGetEnv returns the env value, panics when it is unset or empty
func MustGetEnv(key EnvironmentVariableKey) string {
	if value := os.Getenv(key.String()); value != "" {
		return value
	}
	panic("environment variable " + strconv.Quote(key.String()) + " needs to be set")
}

// CheckPublicEnv catches a public deployment ([EVK_TrustedPlatform] set, release mode)
// still running on the development values docker-compose.yml falls back to:
// email links to a loopback address, or mail sent to the dev mailpit
func CheckPublicEnv() error {
	if IsDevelopment() || EVK_TrustedPlatform.Get("") == "" {
		return nil
	}
	raw := EVK_AppBaseURL.Get("")
	base, err := url.Parse(raw)
	if err != nil || base.Scheme != "https" || base.Hostname() == "" ||
		base.Hostname() == "localhost" || net.ParseIP(base.Hostname()).IsLoopback() {
		return fmt.Errorf("%s must be the public https address when %s is set, got %q",
			EVK_AppBaseURL, EVK_TrustedPlatform, raw)
	}
	if EVK_SMTPHost.Get("") == "mailpit" {
		return fmt.Errorf("%s is the development mailpit, set a real provider when %s is set",
			EVK_SMTPHost, EVK_TrustedPlatform)
	}
	return nil
}

// GetDatabaseURL constructs a database url using environment
//
//	sslmode=disable
//	pool_max_conns=[EVK_DatabaseMaxConns]
//
// a URL so that values with spaces, quotes etc. are escaped
func GetDatabaseURL() string {
	query := url.Values{}
	query.Set("sslmode", "disable")
	query.Set("pool_max_conns", EVK_DatabaseMaxConns.Get("10"))
	//exhaustruct:ignore
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			MustGetEnv(EVK_DatabaseUser),
			MustGetEnv(EVK_DatabasePassword),
		),
		Host:     net.JoinHostPort(MustGetEnv(EVK_DatabaseHost), MustGetEnv(EVK_DatabasePort)),
		Path:     MustGetEnv(EVK_DatabaseName),
		RawQuery: query.Encode(),
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
