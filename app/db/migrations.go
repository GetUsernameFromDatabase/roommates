package db

// will panic on most errors in this package
// as processes here are considered critical for the function of this app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"roommates/logger"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var log = logger.MigrationLoggger

// consts used for sql
const (
	table = "migration_version"
	// version column
	cVersion = "version"
	// this just for more information should it make something easier in the future
	cFileName = "file_name"
	// to know when migration happened
	cMigratedAt = "migrated_at"
)

const (
	// migration error end, used for error messages
	mee = "\nNote: migrations must start from 1 and increment by 1"
)

// File Name Info regular expression, used to get information from the file name
//
// Regex group info:
//  0. full match, in this case the entire filename
//  1. migration version
//  2. migration name
//  3. migration type (up or down)
var fniRegex = regexp.MustCompile(`^(\d{5})_(.+)\.(up|down)\.sql$`)
var ctx = context.Background()

type migrationFiles struct {
	up   string
	down string
}

type migrations struct {
	currentVersion int
	maxVersion     int
	// key is the version of the migration 0-
	migrationFileMap map[int]*migrationFiles
	db               *pgxpool.Pool
}

func MigrateToLatest(db *pgxpool.Pool, migrationDir string) {
	m := NewMigrations(db, migrationDir)
	m.MigrateTo(m.maxVersion)
}

func NewMigrations(db *pgxpool.Pool, migrationDir string) *migrations {
	mVersion := getMigrationVersion(db)

	m := migrations{
		currentVersion:   mVersion,
		maxVersion:       mVersion, // changed by validateMigrations
		migrationFileMap: make(map[int]*migrationFiles),
		db:               db,
	}
	m.loadMigrationFiles(migrationDir)
	return &m
}

func SelectMigrationVersion(db *pgxpool.Pool) int {
	migrationVersion := 0

	// note: only one row is expected to be
	query := fmt.Sprintf(
		`SELECT %s FROM %s LIMIT 1;`,
		cVersion, table,
	)
	err := db.QueryRow(ctx, query).Scan(&migrationVersion)
	if errors.Is(err, sql.ErrNoRows) {
		return 0
	} else if err != nil {
		log.Error().Err(err).Caller().
			Str("query", query).
			Msg("error getting migration version")
		panic(err)
	}

	log.Info().Msgf("current migration version is %v", migrationVersion)
	return migrationVersion
}

// ensure that at least one migration row exists since we will only be updating it
func ensureMigrationRow(db *pgxpool.Pool) {
	insertMigrationRow := func() {
		query := fmt.Sprintf(
			`INSERT INTO %s (%s, %s) VALUES ($1, $2);`,
			table, cFileName, cVersion,
		)
		_, err := db.Exec(ctx, query, "NO_FILE", 0)
		if err != nil {
			log.Error().Err(err).Caller().
				Str("query", query).
				Msg("error initializing migration")
			panic(err)
		}
	}

	migrationVersion := 0
	query := fmt.Sprintf(
		`SELECT %s FROM %s LIMIT 1;`,
		cVersion, table,
	)
	err := db.QueryRow(ctx, query).Scan(&migrationVersion)
	if errors.Is(err, sql.ErrNoRows) {
		insertMigrationRow()
	} else if err != nil {
		log.Error().Err(err).Caller().
			Str("query", query).
			Msg("error getting migration version")
		panic(err)
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
		table, cVersion, cFileName, cMigratedAt,
	)
	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Error().Err(err).Caller().
			Str("query", query).
			Msg("error creating migration table")
		panic(err)
	}
	ensureMigrationRow(db)
}

func updateMigrationVersion(db *pgxpool.Pool, migrationVersion int, file string) {
	query := fmt.Sprintf(`UPDATE %s SET
	%s = $1,
	%s = $2,
	%s = $3;`,
		table, cVersion, cFileName, cMigratedAt,
	)
	_, err := db.Exec(ctx, query, migrationVersion, file, time.Now())
	if err != nil {
		log.Error().Err(err).Caller().
			Str("query", query).
			Msg("error inserting migration version")
		panic(err)
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
		err := fmt.Errorf("migration version must be between 0 and %v (inclusive)", m.maxVersion)
		panic(err)
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

func (m *migrations) updateCurrentMigrationVersion(file string) {
	updateMigrationVersion(m.db, m.currentVersion, file)
	log.Info().
		Int("migrationVersion", m.currentVersion).
		Str("file", file).
		Msg("successful migration")
}

func (m *migrations) runMigrationFile(filePath string) {
	fContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Caller().
			Str("filePath", filePath).
			Msg("error reading file")
		panic(err)
	}

	sqlContent := string(fContent)
	_, err = m.db.Exec(context.Background(), sqlContent)
	if err != nil {
		log.Error().Err(err).Caller().
			Str("filePath", filePath).
			Msg("failed to run migration file")
		panic(err)
	}
}

// Reads migration files from the specified directory
//
// Note: does not recursively read directories
func (m *migrations) loadMigrationFiles(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Error().Err(err).Caller().
			Str("filePath", dir).
			Msg("error reading contents of directory")
		panic(err)
	}
	// it is assumed that if there is no error with os.ReadDir
	//  then this won't cause an error as well
	dirPathAbs, _ := filepath.Abs(dir)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		m.modifyAddMigration(file.Name(), dirPathAbs)
	}

	err = m.validateMigrations()
	if err != nil {
		log.Error().Err(err).Caller().Msg("error during migration validation")
		panic(err)
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
func (m *migrations) modifyAddMigration(migrationFileName string, dirPathAbs string) {
	matches := fniRegex.FindStringSubmatch(migrationFileName)
	if matches == nil {
		log.Warn().
			Msgf("could not get migration info from \"%s\" file name. Ignoring it",
				migrationFileName,
			)
		return
	}
	parsedVersion, err := strconv.ParseInt(matches[1], 10, 8)
	if err != nil {
		log.Error().Err(err).Caller().
			Str("versionFromRegex", matches[1]).
			Str("migrationFileName", migrationFileName).
			Msg("error parsing version from filename")
		panic(err)
	}

	isUpSql := matches[3] == "up"
	fullPath := filepath.Join(dirPathAbs, migrationFileName)
	mVersion := int(parsedVersion)

	if _, ok := m.migrationFileMap[mVersion]; ok {
		if isUpSql {
			m.migrationFileMap[mVersion].up = fullPath
		} else {
			m.migrationFileMap[mVersion].down = fullPath
		}
	} else {
		var migration migrationFiles
		if isUpSql {
			migration = migrationFiles{
				up: fullPath,
			}
		} else {
			migration = migrationFiles{
				down: fullPath,
			}
		}

		m.migrationFileMap[mVersion] = &migration
	}
}
