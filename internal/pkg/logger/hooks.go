package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	serverType = "server_type"
	ipAddress  = "ip_address"
)

type hooks struct {
	httpClient  *http.Client
	conf        Config
	logChats    []int64
	extraFields map[string]string
}

func newHooks(conf Config, logChats []int64, serverTypeInput string) *hooks {
	return &hooks{
		conf: conf,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logChats: logChats,
		extraFields: map[string]string{
			serverType: serverTypeInput,
			ipAddress:  conf.SelfIPAddress,
		},
	}
}

func (s hooks) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (s hooks) Fire(in *logrus.Entry) error {
	s.addExtraFieldsToEntry(in)

	for i := range s.logChats {
		if err := s.sendSingleLog(in, s.logChats[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s hooks) sendSingleLog(in *logrus.Entry, chatID int64) error {
	ctx := context.Background()
	if in.Context != nil {
		ctx = in.Context
	}

	const (
		method        = "/sendMessage"
		requestMethod = http.MethodPost
	)

	type logMessage struct {
		ChatID    int64  `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}

	reqData := logMessage{
		Text:      resolveLoggerMessage(in),
		ChatID:    chatID,
		ParseMode: "HTML",
	}

	reqRaw, err := json.Marshal(reqData)
	if err != nil {
		logrus.Errorf("telegram_logger.SendLogMessage: %v", err)
		return err
	}

	u, err := url.Parse(s.conf.TgBotBaseURL)
	if err != nil {
		logrus.Errorf("[%s][client][%s][%s] parse url '%s' error %v",
			requestMethod, clientName, method, s.conf.TgBotBaseURL, err)
		return err
	}
	u.Path += fmt.Sprintf(tokenTemp, s.conf.TgBotToken) + method

	req, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), bytes.NewBuffer(reqRaw))
	if err != nil {
		logrus.Errorf("[%s][client][%s][%s] create request error: %v",
			requestMethod, clientName, method, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	respHTTP, err := s.httpClient.Do(req)
	if err != nil {
		logrus.Errorf("[%s][client][%s][%s] requesting error: %v",
			requestMethod, clientName, method, err)
		return err
	}
	defer func() { _ = respHTTP.Close }()

	if respHTTP.StatusCode >= 400 {
		var resp errResp
		resp.ReadFromReader(respHTTP.Body)
		if resp.isNotFound() {
			return nil
		}
		logrus.Errorf("[%s][client][%s][%s] bad status code: %d",
			requestMethod, clientName, method, respHTTP.StatusCode)
	}
	return err
}

func resolveLoggerMessage(in *logrus.Entry) string {
	resp := ""
	resp += fmt.Sprintf("<b>level=%s</b>\n", in.Level.String())
	for key, val := range in.Data {
		if key == "" || val == "" {
			continue
		}
		resp += fmt.Sprintf("<u>%s</u>=%s\n", key, val)
	}
	resp += fmt.Sprintf("<u>message</u>=%s\n", in.Message)
	return resp
}

func (s hooks) addExtraFieldsToEntry(in *logrus.Entry) {
	if in.Data == nil {
		in.Data = make(map[string]interface{})
	}
	for key, val := range s.extraFields {
		in.Data[key] = val
	}
}
