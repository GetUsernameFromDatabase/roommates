package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Main = zerolog.New(os.Stdout).With().Timestamp().Logger()
var MigrationLoggger = Main.With().Str("component", "migration").Logger()
var ControllerLoggger = Main.With().Str("component", "controller").Logger()
var RedisLoggger = Main.With().Str("component", "controller").Logger()

// Initializes zerolog as the project logger
// replaces standard log with zerolog
func InitLogger(isDebugMode bool) {
	// https://github.com/rs/zerolog?tab=readme-ov-file#getting-started
	buildInfo, _ := debug.ReadBuildInfo()

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if isDebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.SetFlags(0)
	log.SetOutput(Main)

	Main.Info().Str("GoVersion", buildInfo.GoVersion).Msg("logger initialized")
}
