package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// will get env value, panics if it does not exist
func MustGetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("environment variable " + strconv.Quote(key) + "needs to be set")
}

// constructs a database url using environment
func GetDatabaseURL() string {
	host := MustGetEnv("DATABASE_HOST")
	port := MustGetEnv("DATABASE_PORT")
	dbname := MustGetEnv("POSTGRES_DB")
	user := MustGetEnv("POSTGRES_USER")
	password := MustGetEnv("POSTGRES_PASSWORD")

	return fmt.Sprintf("host=%s port=%s dbname=%s "+
		"user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)
}
