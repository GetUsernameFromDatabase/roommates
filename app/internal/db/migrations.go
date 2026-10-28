package db

// will panic on most errors in this file
// as processes here are considered critical for the function of this app

//lint:file-ignore ST1000 no need to describe this as dbqueries

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"regexp"
	"roommates/internal/logger"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type migrationFiles struct {
	up   string // [migrationFS] path
	down string // [migrationFS] path
}

type migrations struct {
	currentVersion int // taken from db
	maxVersion     int // taken from files
	// key is the version of the migration 0-
	migrationFileMap map[int]*migrationFiles
	db               *pgxpool.Pool
}

// consts used for sql
const (
	TableMigrations = "migration_version"
	ColumnVersion   = "version"
	// this just for more information should it make something easier in the future
	ColumnFilename = "file_name"
	// to know when migration happened
	ColumnMigratedAt = "migrated_at"
)

const (
	// migration error end; used for error messages
	mee = "\nNote: migrations must start from 1 and increment by 1"
)

var log = logger.Main.With().Str("package", "migration").Logger()

//go:embed migrations/*.sql
var migrationFS embed.FS

// File Name Info regular expression, used to get information from the file name
//
// Regex group info:
//  0. full match, in this case the entire filename
//  1. migration version
//  2. migration name
//  3. migration type (up or down)
var fniRegex = regexp.MustCompile(`^(\d{5})_(.+)\.(up|down)\.sql$`)
var ctx = context.Background()

func MigrateToLatest(db *pgxpool.Pool) {
	m := NewMigrations(db)
	m.MigrateTo(m.maxVersion)
}

func NewMigrations(db *pgxpool.Pool) *migrations {
	mVersion := getMigrationVersion(db)

	m := migrations{
		currentVersion:   mVersion,
		maxVersion:       mVersion, // changed by validateMigrations
		migrationFileMap: make(map[int]*migrationFiles),
		db:               db,
	}
	m.loadMigrationFiles()
	return &m
}

func SelectMigrationVersion(db *pgxpool.Pool) int {
	migrationVersion := 0

	// note: only one row is expected to be
	query := fmt.Sprintf(
		`SELECT %s FROM %s LIMIT 1;`,
		ColumnVersion, TableMigrations,
	)
	err := db.QueryRow(ctx, query).Scan(&migrationVersion)
	if errors.Is(err, sql.ErrNoRows) {
		return 0
	} else if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("error getting migration version")
	}

	log.Info().Msgf("current migration version is %v", migrationVersion)
	return migrationVersion
}

// ensure that one migration row exists since we will only be updating it
func ensureMigrationRow(db *pgxpool.Pool) {
	query := fmt.Sprintf(
		`INSERT INTO %[1]s (%[2]s, %[3]s) SELECT 'NO_FILE', 0
		WHERE NOT EXISTS (SELECT 1 FROM %[1]s);`,
		TableMigrations, ColumnFilename, ColumnVersion,
	)
	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("error initializing migration")
	}
}

// ensures existence of the migration table
func ensureMigrationTable(db *pgxpool.Pool) {
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY,
			%s VARCHAR(255) NOT NULL,
			%s TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		TableMigrations, ColumnVersion, ColumnFilename, ColumnMigratedAt,
	)
	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("error creating migration table")
	}
	ensureMigrationRow(db)
}

func updateMigrationVersion(db pgx.Tx, migrationVersion int, file string) {
	query := fmt.Sprintf(`UPDATE %s SET
	%s = $1,
	%s = $2,
	%s = $3;`,
		TableMigrations, ColumnVersion, ColumnFilename, ColumnMigratedAt,
	)
	_, err := db.Exec(ctx, query, migrationVersion, file, time.Now())
	if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("error inserting migration version")
	}
}

// Ensures the existence of migration table and gets current version
func getMigrationVersion(db *pgxpool.Pool) int {
	ensureMigrationTable(db)
	return SelectMigrationVersion(db)
}

// --- migrations struct functions ---

func (m *migrations) MigrateTo(desiredVersion int) {
	if desiredVersion < 0 || desiredVersion > m.maxVersion {
		log.Fatal().Int("min", 0).Int("max", m.maxVersion).
			Int("desiredVersion", desiredVersion).
			Msg("desired migration version must be between min and max")
	}

	if desiredVersion == m.currentVersion {
		log.Info().Int("version", desiredVersion).
			Msg("currently at the desired migration version")
		return
	}

	// version N is reached by running N.up from N-1
	// and left by running N.down back to N-1
	for m.currentVersion < desiredVersion {
		next := m.currentVersion + 1
		m.runMigrationStep(m.migrationFileMap[next].up, next)
	}
	for m.currentVersion > desiredVersion {
		m.runMigrationStep(m.migrationFileMap[m.currentVersion].down, m.currentVersion-1)
	}
}

// runs the migration file and records newVersion in a single transaction,
// so a failing file leaves neither the schema nor the version half-updated
//
// file is used with [migrationFS.ReadFile]
func (m *migrations) runMigrationStep(file string, newVersion int) {
	fContent, err := migrationFS.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).
			Str("fileName", file).
			Str("from", "migrationFS").
			Msg("error reading file")
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		log.Fatal().Err(err).Str("fileName", file).
			Msg("could not begin migration transaction")
	}
	// no-op after a successful commit
	defer func() { _ = tx.Rollback(ctx) }()

	// no arguments so pgx uses the simple protocol, which allows multiple statements
	if _, err = tx.Exec(ctx, string(fContent)); err != nil {
		log.Fatal().Err(err).
			Str("fileName", file).
			Msg("failed to run migration file")
	}
	updateMigrationVersion(tx, newVersion, file)
	if err = tx.Commit(ctx); err != nil {
		log.Fatal().Err(err).Str("fileName", file).
			Msg("could not commit migration transaction")
	}

	m.currentVersion = newVersion
	log.Info().
		Int("migrationVersion", m.currentVersion).
		Str("file", file).
		Msg("successful migration")
}

// Reads migration files from [migrationFS]
func (m *migrations) loadMigrationFiles() {
	dir := "migrations"
	files, err := migrationFS.ReadDir(dir)
	if err != nil {
		log.Fatal().Err(err).Msgf("error reading %s directory", dir)
	}

	for _, file := range files {
		// the dir is required to read from [migrationFS] later
		// file.Name() just gets base name, seems like file.Info.Name() also does the same
		// do not use [filepath.Join] as that uses os specific seperator
		joinedName := dir + "/" + file.Name()
		m.modifyAddMigration(joinedName)
	}

	err = m.validateMigrations()
	if err != nil {
		log.Fatal().Err(err).Msg("error during migration validation")
	}
}

// Validates migrations and sets max migration version:
//  1. migrations must start from 1 and increment by 1 (growing order)
//  2. migrations must have up sql and down sql
//
// len(map) keys checked as 1..len means no key can be outside that range
func (m *migrations) validateMigrations() error {
	migrationLength := len(m.migrationFileMap)

	for i := 1; i <= migrationLength; i++ {
		migration, ok := m.migrationFileMap[i]
		if !ok {
			return fmt.Errorf("migration version %v missing"+mee, i)
		}

		// if the entry exists then there must be at least one field defined
		if migration.up == "" {
			return fmt.Errorf("migration up file missing for version %v", i)
		} else if migration.down == "" {
			return fmt.Errorf("migration down file missing for version %v", i)
		}
	}

	m.maxVersion = migrationLength
	return nil
}

// Modifies or adds a migration
//
// Uses regular expression to extract information from file name
func (m *migrations) modifyAddMigration(migrationFilePath string) {
	baseName := filepath.Base(migrationFilePath)
	matches := fniRegex.FindStringSubmatch(baseName)
	if matches == nil {
		log.Warn().Str("filename", baseName).
			Msg("could not get migration info; ignoring it")
		return
	}
	mVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal().Err(err).
			Str("versionFromRegex", matches[1]).
			Str("filename", baseName).
			Msg("could not get version from filename")
	}

	isUpSQL := matches[3] == "up"

	migration, ok := m.migrationFileMap[mVersion]
	if !ok {
		migration = &migrationFiles{}
		m.migrationFileMap[mVersion] = migration
	}

	if isUpSQL {
		migration.up = migrationFilePath
	} else {
		migration.down = migrationFilePath
	}
}
