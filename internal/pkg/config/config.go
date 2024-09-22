package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

func Load(conf interface{}) error {
	configFileNames := make([]string, 0, 1)

	path := "./configs/env"
	if _, err := os.Stat(path); err == nil {
		configFileNames = append(configFileNames, path)
	}

	if err := loadConfig(conf, "", configFileNames...); err != nil {
		return trace.FuncNameWithErrorMsg(err, "load environment")
	}
	return nil
}

func loadConfig(c interface{}, prefix string, filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.WithField("filenames", filenames).Info("config file not found, using defaults")
	}

	err = envconfig.Process(prefix, c)
	if err != nil {
		return fmt.Errorf("error env config loading: %w", err)
	}

	return nil
}
