package main

import (
	"context"
	"runtime"

	"github.com/tarmalonchik/falaemae/internal/app/core/bootstrap"
	"github.com/tarmalonchik/falaemae/internal/app/core/config"
	"github.com/tarmalonchik/falaemae/internal/pkg/version"

	"github.com/sirupsen/logrus"
	"github.com/vkidmode/server-core/pkg/core"
)

func init() {
	if version.Service == "" {
		version.Service = "core"
	}
}

func main() {
	ctx := context.Background()

	runtime.GOMAXPROCS(runtime.NumCPU())

	conf, err := config.GetConfig()
	if err != nil {
		logrus.Errorf("failed to load environment: %v", err)
		return
	}

	services, err := bootstrap.GetServices(ctx, conf)
	if err != nil {
		logrus.Errorf("failed to initiate service locator: %v", err)
		return
	}

	app := core.NewCore(services.GetCustomLogger(), conf.Core.GracefulTimeout, 50)
	app.AddRunner(services.GetTgWorker().Run, true)

	err = app.Launch(ctx)
	if err != nil {
		logrus.Errorf("application error: %v", err)
	}
}
