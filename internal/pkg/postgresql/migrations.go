package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations processes migrations
func RunMigrations(db *sql.DB, sourceURL, dbName string) error {
	if sourceURL == "" {
		return errors.New("no cfg.GetPGMigrationsPath provided")
	}

	px, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName:          dbName,
		MultiStatementMaxSize: postgres.DefaultMultiStatementMaxSize,
	})
	if err != nil {
		return fmt.Errorf("[migration] migration file %s error db instance %w", sourceURL, err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, dbName, px)
	if err != nil {
		return fmt.Errorf("[migration] migration file %s error migrate instance %w", sourceURL, err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("[migration] migration file %s error migrate %w", sourceURL, err)
	}

	return nil
}
