package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg Config) (*Postgres, error) {
	if cfg.GetPGDSN() != "" {
		sql.Register("instrumented-postgres", pq.Driver{})
		dbSQL, err := sql.Open("instrumented-postgres", cfg.GetPGDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to connect database: %w", err)
		}

		dbSQL.SetMaxIdleConns(cfg.GetPGMaxIdleConns())
		dbSQL.SetMaxOpenConns(cfg.GetPGMaxOpenConns())
		dbSQL.SetConnMaxLifetime(cfg.GetPGTimeout())

		return &Postgres{db: dbSQL}, nil
	}

	return nil, nil
}

func (p *Postgres) GetDB() *sql.DB {
	return p.db
}

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}
