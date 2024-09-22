package core

import (
	"context"
	"time"
)

const (
	pleaseUseButtons = "Пожалуйста воспользуйтесь кнопками"
)

type Config struct {
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" default:"40s"`
}

type commandsProcessorFunc func(ctx context.Context, meta MetaData) error

type processorWithBranch struct {
	processor    commandsProcessorFunc
	worksForever bool
}
