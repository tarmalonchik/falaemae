package config

import (
	"github.com/tarmalonchik/falaemae/internal/app/core/svc/core"
	"github.com/tarmalonchik/falaemae/internal/pkg/config"
	"github.com/tarmalonchik/falaemae/internal/pkg/logger"
	"github.com/tarmalonchik/falaemae/internal/pkg/postgresql"
	"github.com/tarmalonchik/falaemae/internal/pkg/telegram"
)

// Config contains all environment variables
type Config struct {
	TelegramLogger logger.Config
	Telegram       telegram.Config
	Postgres       postgresql.Config
	Core           core.Config

	Debug bool `envconfig:"DEBUG" required:"true"`
}

func GetConfig() (conf *Config, err error) {
	conf = &Config{}
	err = config.Load(conf)
	return conf, err
}
