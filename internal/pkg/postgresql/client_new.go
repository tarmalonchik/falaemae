package postgresql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SQLXClient is a shortcut structure to a new sqlx DB
type SQLXClient struct {
	*sqlx.DB
}

// NewSQLXClient creates new database connection to postgres
func NewSQLXClient(db *sql.DB, _ Config) *SQLXClient {
	return &SQLXClient{DB: sqlx.NewDb(db, "postgres")}
}
