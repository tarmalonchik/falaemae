package telegram

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	AdminCommand   = "Админка ⚙️"
	StartCommand   = "/start"
	PayCommand     = "Подписка 💸"
	SupportCommand = "Поддержка 🧑"
	InfoCommand    = "Как подключить ❓"
	ServersList    = "Поменять страну 🔄"
	ReferralLink   = "Пригласить друга 🙋🏻‍"
)

type AdminEvent string

const (
	EventOrderCreated            = AdminEvent("Пользователь разместил заказ")
	EventDemoFinishedAndNotPayed = AdminEvent("У пользователя кончилась демо-версия и далее не была оплачена")
	EventNewUser                 = AdminEvent("К нам пришел новый пользователь")
	EventUserBuyVpn              = AdminEvent("Пользователь купил впн")
)

type admin struct {
	ChatID        int64
	TelegramNick  string
	SendAdminLogs bool
}

type KeyBoardType struct {
	defaultButtons  *tgbotapi.ReplyKeyboardMarkup
	adminEnrich     []tgbotapi.KeyboardButton
	admins          []admin
	callBackButtons *tgbotapi.InlineKeyboardMarkup
}

func (s *KeyBoardType) GetButtons(chatID int64) interface{} {
	var resp tgbotapi.ReplyKeyboardMarkup

	if s.callBackButtons != nil {
		return s.callBackButtons
	}

	if s.defaultButtons != nil {
		resp = *s.defaultButtons
	}

	for i := range s.admins {
		if s.admins[i].ChatID == chatID {
			resp.Keyboard = append(resp.Keyboard, s.adminEnrich)
		}
	}
	return &resp
}

func NewKeyboard(req *NewKeyboardRequest) (*KeyBoardType, error) {
	callbackButtons := generateInlineTemplate(req)
	return &KeyBoardType{callBackButtons: callbackButtons}, nil
}

var KeyBoardItem = KeyBoardType{
	admins: []admin{
		{
			TelegramNick:  "tarmalonchik",
			ChatID:        496869421,
			SendAdminLogs: true,
		},
	},
	adminEnrich: []tgbotapi.KeyboardButton{
		{
			Text: AdminCommand,
		},
	},
	defaultButtons: &tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{
					Text: PayCommand,
				},
				{
					Text: InfoCommand,
				},
			},
			{
				{
					Text: ServersList,
				},
				{
					Text: SupportCommand,
				},
			},
			{
				{
					Text: ReferralLink,
				},
			},
		},
	},
}

func (s *KeyBoardType) getAdmins() []admin {
	return s.admins
}

func (s *KeyBoardType) IsAdmin(in interface{}) bool {
	switch val := in.(type) {
	case string:
		for i := range s.admins {
			if s.admins[i].TelegramNick == val {
				return true
			}
		}
	case int64:
		for i := range s.admins {
			if s.admins[i].ChatID == val {
				return true
			}
		}
	case int:
		for i := range s.admins {
			if s.admins[i].ChatID == int64(val) {
				return true
			}
		}
	}
	return false
}
