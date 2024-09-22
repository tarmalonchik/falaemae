package bootstrap

import (
	"context"
	"net/url"
	"strings"

	"github.com/tarmalonchik/falaemae/internal/app/core/config"
	"github.com/tarmalonchik/falaemae/internal/pkg/inmemory"
	"github.com/tarmalonchik/falaemae/internal/pkg/logger"
	"github.com/tarmalonchik/falaemae/internal/pkg/postgresql"
	"github.com/tarmalonchik/falaemae/internal/pkg/storage"
	"github.com/tarmalonchik/falaemae/internal/pkg/telegram"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
	"github.com/tarmalonchik/falaemae/internal/pkg/version"
)

type ClientsContainer struct {
	tgClient     *telegram.Client
	db           *postgresql.Postgres
	sqlx         *postgresql.SQLXClient
	storage      *storage.Model
	logger       *logger.Client
	documentsMap inmemory.InMemory[string, string]
}

func getClients(ctx context.Context, conf *config.Config) (*ClientsContainer, error) {
	clients := &ClientsContainer{
		logger:       logger.NewClient(conf.TelegramLogger, version.Service),
		documentsMap: inmemory.New[string, string](),
	}

	pgClient, err := postgresql.New(conf.Postgres)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "init postgres")
	}
	clients.db = pgClient
	clients.sqlx = postgresql.NewSQLXClient(clients.db.GetDB(), conf.Postgres)

	purl, err := url.Parse(conf.Postgres.GetPGDSN())
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "dns url parse")
	}

	if err = postgresql.RunMigrations(clients.db.GetDB(), conf.Postgres.GetPGMigrationsPath(), strings.Trim(purl.Path, "/")); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "migration")
	}

	clients.tgClient, err = telegram.NewClient(conf.Telegram, clients.documentsMap)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "init telegram")
	}

	clients.storage = storage.NewModel(clients.sqlx)

	return clients, nil
}
