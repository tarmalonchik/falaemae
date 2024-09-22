package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

var sendLogsTo = []int64{496869421}

const (
	tokenTemp  = "/bot%s"
	clientName = "tg-logger"
)

type Client struct {
	serverType string
	logger     *logrus.Logger
}

func NewClient(conf Config, serverType string) *Client {
	cl := &Client{
		serverType: serverType,
	}
	cl.logger = logrus.New()
	cl.logger.Hooks.Add(newHooks(conf, sendLogsTo, serverType))
	if !conf.IsCronJob {
		cl.Infof("logger started")
	}
	return cl
}

func (c *Client) Errorf(err error, format string, args ...any) {
	c.logger.WithError(err).Errorf(format, args...)
}

func (c *Client) Error(err error) {
	c.logger.WithError(err).Error()
}

func (c *Client) Infof(format string, args ...any) {
	c.logger.Infof(format, args...)
}

func (c *Client) Warnf(format string, args ...any) {
	c.logger.Warnf(format, args...)
}

func (c *Client) Log(ctx context.Context, message string) {
	c.logger.WithContext(ctx).WithError(fmt.Errorf("application runner")).Errorf(message)
}
