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

// ensure that at least one migration row exists since we will only be updating it
func ensureMigrationRow(db *pgxpool.Pool) {
	insertMigrationRow := func() {
		query := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2);`,
			TableMigrations, ColumnFilename, ColumnVersion,
		)
		_, err := db.Exec(ctx, query, "NO_FILE", 0)
		if err != nil {
			log.Fatal().Err(err).
				Str("query", query).
				Msg("error initializing migration")
		}
	}

	migrationVersion := 0
	query := fmt.Sprintf(
		`SELECT %s FROM %s LIMIT 1;`,
		ColumnVersion, TableMigrations,
	)
	err := db.QueryRow(ctx, query).Scan(&migrationVersion)
	if errors.Is(err, sql.ErrNoRows) {
		insertMigrationRow()
	} else if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("error getting migration version")
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

func updateMigrationVersion(db *pgxpool.Pool, migrationVersion int, file string) {
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

	loopIncrement := 1
	loopFileGet := func(i int) string {
		files := m.migrationFileMap[i]
		return files.up
	}
	// change previous if migration direction is down
	if desiredVersion < m.currentVersion {
		loopIncrement = -1
		loopFileGet = func(i int) string {
			// when we want to migrate down we should use the previos migration file
			files := m.migrationFileMap[i+1]
			return files.down
		}
	}

	// migrations files only start from 1
	loopStart := max(m.currentVersion, 1)
	if loopIncrement < 0 && len(m.migrationFileMap) == loopStart {
		// down migrations use the i+1 to get migration file that migrates to i
		loopStart -= 1
	}
	for i := loopStart; ; i = i + loopIncrement {
		file := loopFileGet(i)
		m.runMigrationFile(file)
		m.currentVersion = i
		m.updateCurrentMigrationVersion(file)
		if i == desiredVersion {
			break
		}
	}
}

// file is used just to make info gathering easier later
func (m *migrations) updateCurrentMigrationVersion(file string) {
	updateMigrationVersion(m.db, m.currentVersion, file)
	log.Info().
		Int("migrationVersion", m.currentVersion).
		Str("file", file).
		Msg("successful migration")
}

// file is used with [migrationFS.ReadFile]
func (m *migrations) runMigrationFile(file string) {
	fContent, err := migrationFS.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).
			Str("fileName", file).
			Str("from", "migrationFS").
			Msg("error reading file")
	}

	sqlContent := string(fContent)
	_, err = m.db.Exec(context.Background(), sqlContent)
	if err != nil {
		log.Fatal().Err(err).
			Str("fileName", file).
			Msg("failed to run migration file")
	}
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
func (m *migrations) validateMigrations() error {
	migrationLength := len(m.migrationFileMap)
	migrationsValidated := 0

	// TODO: would be nice to aggregate these errors together
	// low priority however so won't spend time on it
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
		migrationsValidated += 1
	}

	// if starting from 1 did not validate all migrations
	// then there are migrations that are below 0 version
	if migrationLength != migrationsValidated {
		return errors.New("migrations must start from 1" + mee)
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
	parsedVersion, err := strconv.ParseInt(matches[1], 10, 8)
	if err != nil {
		log.Fatal().Err(err).
			Str("versionFromRegex", matches[1]).
			Str("filename", baseName).
			Msg("could not get version from filename")
	}

	isUpSQL := matches[3] == "up"
	mVersion := int(parsedVersion)

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
