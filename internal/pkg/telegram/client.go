package telegram

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	telegram "github.com/Syfaro/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/vkidmode/telegram_tree"

	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

const (
	userIsDeactivated      = "user is deactivated"
	botWasBlockedByTheUser = "bot was blocked by the user"
	chatNotFound           = "chat not found"
)

type inMemory interface {
	Get(key string) (value string, ok bool)
	AddData(key string, value string)
}

type Client struct {
	conf         Config
	tg           *telegram.BotAPI
	httpClient   http.Client
	updates      telegram.UpdatesChannel
	documentsMap inMemory
}

func NewClient(conf Config, inMemory inMemory) (*Client, error) {
	cl := &Client{}
	cl.conf = conf

	tgAPI, err := telegram.NewBotAPI(conf.GetTgBotToken())
	if err != nil {
		return nil, fmt.Errorf("error init telegram client: %w", err)
	}
	cl.tg = tgAPI
	cl.httpClient = http.Client{
		Timeout: 30 * time.Second,
	}
	cl.documentsMap = inMemory

	return cl, nil
}

func (c *Client) Get() (telegram.UpdatesChannel, error) {
	var err error

	if c.updates != nil {
		return c.updates, nil
	}
	c.updates, err = c.tg.GetUpdatesChan(telegram.NewUpdate(0))
	if err != nil {
		return nil, err
	}
	return c.updates, nil
}

func (c *Client) Stop() {
	c.tg.StopReceivingUpdates()
}

func (c *Client) GetBotName() (string, error) {
	me, err := c.tg.GetMe()
	if err != nil {
		return "", fmt.Errorf("telegram.GetBotName error getting me: %w", err)
	}
	return me.String(), nil
}

func (c *Client) SendMessage(chatID sql.NullInt64, message string) error {
	return c.SendMessageWithSettings(chatID, &KeyBoardItem, message, true, ParseModeHtml)
}

func (c *Client) SendOrUpdateMessage(chatID int64, messageID int64, buttons *KeyBoardType, message string, disablePrev bool) error {
	if messageID != 0 {
		return c.UpdateMessage(chatID, messageID, buttons, message, disablePrev)
	}
	return c.SendMessageWithSettings(sql.NullInt64{Valid: true, Int64: chatID}, buttons, message, disablePrev, ParseModeHtml)
}

func (c *Client) SendMessageWithSettings(chatID sql.NullInt64, buttons *KeyBoardType, message string, disablePrev bool, mode ParseMode) error {
	if !chatID.Valid {
		return nil
	}

	tgMsg := telegram.NewMessage(chatID.Int64, message)

	if buttons != nil {
		tgMsg.ReplyMarkup = buttons.GetButtons(chatID.Int64)
	}

	switch mode {
	case ParseModeHtml:
		tgMsg.ParseMode = telegram.ModeHTML
	case ParseModeMarkdown:
		tgMsg.ParseMode = telegram.ModeMarkdown
	}

	tgMsg.DisableWebPagePreview = disablePrev

	_, err := c.tg.Send(tgMsg)
	if err != nil {
		if strings.Contains(err.Error(), userIsDeactivated) ||
			strings.Contains(err.Error(), botWasBlockedByTheUser) || strings.Contains(err.Error(), chatNotFound) {
			return nil
		}
		logrus.Error(trace.FuncNameWithErrorMsg(err, "unexpected sending message error"))
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}

func (c *Client) UpdateMessage(chatID, messageID int64, buttons *KeyBoardType, message string, disablePrev bool) error {
	updateMsg := telegram.NewEditMessageText(chatID, int(messageID), message)
	if buttons != nil {
		replyMarkup, ok := buttons.GetButtons(chatID).(*telegram.InlineKeyboardMarkup)
		if !ok {
			return fmt.Errorf("error while interface conversion")
		}
		updateMsg.ReplyMarkup = replyMarkup
	}

	updateMsg.ParseMode = telegram.ModeHTML
	updateMsg.DisableWebPagePreview = disablePrev

	_, err := c.tg.Send(updateMsg)
	if err != nil {
		if strings.Contains(err.Error(), userIsDeactivated) ||
			strings.Contains(err.Error(), botWasBlockedByTheUser) || strings.Contains(err.Error(), chatNotFound) {
			return nil
		}
		return fmt.Errorf("error while updating message: %w", err)
	}
	return nil
}

func (c *Client) DeleteMessage(messageID int64, chatID int64) {
	delMsg := telegram.DeleteMessageConfig{
		MessageID: int(messageID),
		ChatID:    chatID,
	}
	_, _ = c.tg.DeleteMessage(delMsg)
}

func enrichWithFooter(in *telegram.InlineKeyboardMarkup, callbackBack, callbackSkip string) {
	var (
		disabled      = " "
		callBackClose = telegram_tree.CallBackClose
		backText      string
		forwardText   string
	)

	if callbackBack == "" {
		callbackBack = telegram_tree.CallBackIgnore
		backText = disabled
	} else {
		backText = "⬅️"
	}

	if callbackSkip == "" {
		callbackSkip = telegram_tree.CallBackIgnore
		forwardText = disabled
	} else {
		forwardText = "➡️"
	}

	in.InlineKeyboard = append(in.InlineKeyboard, []telegram.InlineKeyboardButton{})
	in.InlineKeyboard[len(in.InlineKeyboard)-1] =
		append(in.InlineKeyboard[len(in.InlineKeyboard)-1], []telegram.InlineKeyboardButton{
			{
				Text:         backText,
				CallbackData: &callbackBack,
			},
			{
				Text:         "❌",
				CallbackData: &callBackClose,
			},
			{
				Text:         forwardText,
				CallbackData: &callbackSkip,
			},
		}...,
		)
}

func generateInlineTemplate(req *NewKeyboardRequest) *telegram.InlineKeyboardMarkup {
	if req == nil {
		return nil
	}

	response := &telegram.InlineKeyboardMarkup{}

	for i := range req.NextNodes {
		if i > 20 {
			break
		}
		if req.NextNodes[i] != nil {
			var (
				callBack                     string
				switchInlineQueryCurrentChat *string
			)
			if req.NextNodes[i].GetTelegram().GetSwitchInlineQueryCurrentChat() == nil {
				callBack = req.NextNodes[i].GetCallback()
			} else {
				switchInlineQueryCurrentChat = req.NextNodes[i].GetTelegram().GetSwitchInlineQueryCurrentChat()
			}

			response.InlineKeyboard = append(response.InlineKeyboard, []telegram.InlineKeyboardButton{})
			response.InlineKeyboard[len(response.InlineKeyboard)-1] =
				append(
					response.InlineKeyboard[len(response.InlineKeyboard)-1],
					[]telegram.InlineKeyboardButton{
						{
							Text:                         req.NextNodes[i].GetTelegram().GetTabTxt(),
							CallbackData:                 &callBack,
							SwitchInlineQueryCurrentChat: switchInlineQueryCurrentChat,
						},
					}...,
				)
		}
	}
	if !req.HideBar {
		enrichWithFooter(response, req.CallbackBack, req.CallbackSkip)
	}
	return response
}

func (c *Client) SendMessageForAdmins(event AdminEvent, msg string) error {
	resultMSG := ""
	if event != "" {
		resultMSG += "<b>" + string(event) + "</b>"
	}
	resultMSG += msg

	destList := make([]int64, 0)

	for _, val := range KeyBoardItem.getAdmins() {
		if val.SendAdminLogs {
			destList = append(destList, val.ChatID)
		}
	}

	for i := range destList {
		if err := c.SendMessage(sql.NullInt64{Valid: true, Int64: destList[i]}, resultMSG); err != nil {
			return trace.FuncNameWithErrorMsg(err, "sending message")
		}
	}
	return nil
}
