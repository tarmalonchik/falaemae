package bootstrap

import (
	"context"

	"github.com/tarmalonchik/falaemae/internal/app/core/config"
	"github.com/tarmalonchik/falaemae/internal/app/core/svc/core"
	"github.com/tarmalonchik/falaemae/internal/app/core/workers/tgupdates"
	"github.com/tarmalonchik/falaemae/internal/pkg/logger"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

type ServiceContainer struct {
	conf     *config.Config
	clients  *ClientsContainer
	tgWorker *tgupdates.TelegramWorker
	coreSvc  *core.Service
}

func GetServices(ctx context.Context, conf *config.Config) (*ServiceContainer, error) {
	var (
		err error
		sv  = &ServiceContainer{conf: conf}
	)

	if sv.clients, err = getClients(ctx, conf); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting clients")
	}

	sv.coreSvc, err = core.NewService(
		ctx,
		conf.Core,
		sv.clients.storage,
		sv.clients.logger,
		sv.clients.tgClient,
	)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "creating core service")
	}

	sv.tgWorker = tgupdates.NewTelegramWorker(sv.clients.tgClient, sv.coreSvc, sv.clients.logger)
	return sv, nil
}

func (s *ServiceContainer) GetTgWorker() *tgupdates.TelegramWorker {
	return s.tgWorker
}

func (s *ServiceContainer) GetCustomLogger() *logger.Client {
	return s.clients.logger
}
