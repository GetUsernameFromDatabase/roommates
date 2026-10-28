package globals

//lint:file-ignore ST1003 goal is to ignore underscore issues

import (
	"fmt"
	"os"
	"roommates/validators"
	"strconv"

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
	// used to send emails via resend
	EVK_ResendApiKey EnvironmentVariableKey = "RESEND_API_KEY"

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
func GetDatabaseURL() string {
	host := MustGetEnv(EVK_DatabaseHost)
	port := MustGetEnv(EVK_DatabasePort)
	dbname := MustGetEnv(EVK_DatabaseName)
	user := MustGetEnv(EVK_DatabaseUser)
	password := MustGetEnv(EVK_DatabasePassword)

	return fmt.Sprintf("host=%s port=%s dbname=%s "+
		"user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)
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
	err = validators.ValidateCookieHashBlockPair(prevHash, prevBlock)
	if err != nil {
		return pairs, err
	}
	pairs = append(pairs, prevHash, prevBlock)
	return pairs, nil
}
