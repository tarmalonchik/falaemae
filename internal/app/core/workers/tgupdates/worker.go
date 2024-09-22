package tgupdates

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/tarmalonchik/falaemae/internal/app/core/svc/core"
	"github.com/tarmalonchik/falaemae/internal/pkg/logger"
	"github.com/tarmalonchik/falaemae/internal/pkg/telegram"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

type TelegramWorker struct {
	core     *core.Service
	tgClient *telegram.Client
	logger   *logger.Client
}

func NewTelegramWorker(
	tgClient *telegram.Client,
	core *core.Service,
	logger *logger.Client,
) *TelegramWorker {
	return &TelegramWorker{
		core:     core,
		tgClient: tgClient,
		logger:   logger,
	}
}

func (t *TelegramWorker) Run(ctx context.Context) error {
	updates, err := t.tgClient.Get()
	if err != nil {
		return fmt.Errorf("tgupdates.Run error getting updates channel: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			t.tgClient.Stop()
			logrus.Infof("%s stopped successfull", trace.FuncName())
			return nil
		case i := <-updates:
			if err = t.core.ProcessUpdate(ctx, &i); err != nil {
				t.logger.Errorf(trace.FuncNameWithError(err), "processing")
			}
		}
	}
}
