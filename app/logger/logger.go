// Package logger this projects logger
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

// .Caller() https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log

var Main = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

// TODO: filelogger
// fileLogger := &lumberjack.Logger{
// Filename: "app.log", // Log file path
// MaxSize: 5, // Max size in MB before rotation
// MaxBackups: 3, // Max number of old log files to keep
// MaxAge: 28, // Max age in days to keep a log file
// Compress: true, // Compress rotated files
// }

// InitLogger initializes zerolog logger
//
// replaces standard log with zerolog
func InitLogger(isDebugMode bool) {
	// https://github.com/rs/zerolog?tab=readme-ov-file#getting-started
	buildInfo, _ := debug.ReadBuildInfo()

	zerolog.TimeFieldFormat = time.RFC3339
	// https://github.com/rs/zerolog?tab=readme-ov-file#error-logging-with-stacktrace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// Lshortfile https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
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
