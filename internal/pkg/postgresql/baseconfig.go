package postgresql

import (
	"fmt"
	"time"
)

type Config struct {
	PgProto          string        `envconfig:"POSTGRES_PROTO" default:"tcp"`
	PgURL            string        `envconfig:"POSTGRES_URL" default:""`
	PgAddress        string        `envconfig:"POSTGRES_ADDRESS" default:"localhost"`
	PgPort           string        `envconfig:"POSTGRES_PORT" default:""`
	PgDb             string        `envconfig:"POSTGRES_DB" default:""`
	PgUser           string        `envconfig:"POSTGRES_USER" default:""`
	PgPassword       string        `envconfig:"POSTGRES_PASS" default:""`
	PgSSLMode        string        `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
	PgMigrationsPath string        `envconfig:"POSTGRES_MIGRATIONSPATH" default:"file:migrations"`
	PgTimeout        time.Duration `envconfig:"POSTGRES_TIMEOUT" default:"5s"`
	PGMaxIdleConns   int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"10"`
	PGMaxOpenConns   int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"15"`
}

func (c *Config) GetPGMigrationsPath() string { return c.PgMigrationsPath }
func (c *Config) GetPGMaxIdleConns() int      { return c.PGMaxIdleConns }
func (c *Config) GetPGMaxOpenConns() int      { return c.PGMaxOpenConns }
func (c *Config) GetPGTimeout() time.Duration { return c.PgTimeout }

func (c *Config) GetPGAddress() string {
	return c.PgAddress + ":" + c.PgPort
}

func (c *Config) GetPGDSN() string {
	pgURL := c.PgURL
	if pgURL == "" {
		pgURL = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			c.PgUser,
			c.PgPassword,
			c.GetPGAddress(),
			c.PgDb,
			c.PgSSLMode,
		)
	}
	return pgURL
}
