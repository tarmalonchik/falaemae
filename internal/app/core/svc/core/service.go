package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/uuid"
	tgt "github.com/vkidmode/telegram_tree"

	"github.com/tarmalonchik/falaemae/internal/pkg/logger"
	storageSdk "github.com/tarmalonchik/falaemae/internal/pkg/storage"
	"github.com/tarmalonchik/falaemae/internal/pkg/telegram"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

type storage interface {
	UpdateTelegramNick(ctx context.Context, chatID int64, nick *string) error
	GetUserByChatID(ctx context.Context, chatID int64) (user storageSdk.User, err error)
	HandlePGTransaction(pgTx storageSdk.PGTransactionFn) error
	InitUser(ctx context.Context, inputUser storageSdk.User) (user storageSdk.User, err error)
	GetUserPayload(ctx context.Context, userID uuid.UUID, orderID *uuid.UUID) (string, error)
}

type telegramClient interface {
	DeleteMessage(messageID int64, chatID int64)
	SendOrUpdateMessage(chatID int64, messageID int64, buttons *telegram.KeyBoardType, message string, disablePrev bool) error
	SendMessage(chatID sql.NullInt64, message string) error
	SendMessageForAdmins(event telegram.AdminEvent, msg string) error
}

type Service struct {
	ctx            context.Context
	conf           Config
	storage        storage
	logger         *logger.Client
	telegramClient telegramClient
	commandsMap    map[string]processorWithBranch
}

func NewService(
	ctx context.Context,
	conf Config,
	storage storage,
	logger *logger.Client,
	telegramClient telegramClient,
) (service *Service, err error) {
	service = &Service{
		ctx:            ctx,
		conf:           conf,
		storage:        storage,
		logger:         logger,
		telegramClient: telegramClient,
	}

	service.commandsMap = map[string]processorWithBranch{
		telegram.InfoCommand:    {service.genInstructionsBranch(), false},
		telegram.PayCommand:     {service.genPaymentBranch(), false},
		telegram.AdminCommand:   {service.isAdminCommand(service.genAdminBranch()), true},
		telegram.SupportCommand: {service.genSendSupportBranch(), true},
	}

	if err := tgt.ReplaceSymbols(symbolsToNum); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "replace symbols in tree")
	}

	if _, err := service.genNewTree(ctx, 0); err != nil { // for checking template
		return nil, err
	}
	return service, nil
}

func (t *Service) isAdminCommand(do commandsProcessorFunc) commandsProcessorFunc {
	return func(ctx context.Context, info MetaData) error {
		user, err := t.storage.GetUserByChatID(ctx, info.ChatID)
		if err != nil {
			return fmt.Errorf("tgupdates.adminCommand error getting user: %w", err)
		}
		if !user.IsAdmin() {
			return fmt.Errorf("tgupdates.adminCommand not admin")
		}
		return do(ctx, info)
	}
}

func (t *Service) ProcessUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if err := t.updateUserNick(ctx, update); err != nil {
		t.logger.Errorf(err, "updating user info")
	}

	if update.CallbackQuery != nil {
		err := t.commonProcessor(ctx, &MetaData{
			MessageID: int64(update.CallbackQuery.Message.MessageID),
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			Callback:  update.CallbackQuery.Data,
		})
		if err != nil {
			return trace.FuncNameWithErrorMsg(err, "doing callback")
		}
		return nil
	}
	if update.Message == nil {
		return nil
	}

	if val, ok := t.commandsMap[update.Message.Text]; ok {
		if val.worksForever {
			if err := val.processor(ctx, MetaData{ChatID: update.Message.Chat.ID}); err != nil {
				return trace.FuncNameWithErrorMsg(err, "processing")
			}
			return nil
		}
	}

	user, err := t.storage.GetUserByChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		if errors.Is(err, storageSdk.ErrUserNotFound) {
			return t.processNewUserUpdate(ctx, update)
		}
		return trace.FuncNameWithErrorMsg(err, "getting user")
	}

	if haveStartPrefix(update.Message.Text) {
		_ = t.telegramClient.SendMessage(user.ChatID, pleaseUseButtons)
		return nil
	}

	if processorBranch, ok := t.commandsMap[update.Message.Text]; ok {
		if err = processorBranch.processor(ctx, MetaData{ChatID: update.Message.Chat.ID}); err != nil {
			return trace.FuncNameWithErrorMsg(err, "commands processor")
		}
		return nil
	}

	//if err = t.processCustomMessages(ctx, update); err != nil {
	//	return trace.FuncNameWithErrorMsg(err, "custom processor")
	//}

	return nil
}

func (t *Service) processNewUserUpdate(ctx context.Context, update *tgbotapi.Update) (err error) {
	if update == nil {
		return nil
	}
	if haveStartPrefix(update.Message.Text) {
		if err = t.processNewUserStart(ctx, update); err != nil {
			return trace.FuncNameWithErrorMsg(err, "processing")
		}
		return nil
	}
	t.sendTapStart(update.Message.Chat.ID)
	return nil
}

func (t *Service) sendTapStart(chatID int64) {
	msg := "Пожалуйста нажмите на команду /start"
	_ = t.telegramClient.SendMessage(sql.NullInt64{Valid: true, Int64: chatID}, msg)
}

func haveStartPrefix(text string) bool {
	if len(text) >= 6 {
		if text[:6] == telegram.StartCommand {
			return true
		}
	}
	return false
}

func (t *Service) updateUserNick(ctx context.Context, update *tgbotapi.Update) error {
	if update == nil {
		return nil
	}
	if update.CallbackQuery != nil {
		return t.storage.UpdateTelegramNick(ctx, update.CallbackQuery.Message.Chat.ID, &update.CallbackQuery.Message.Chat.UserName)
	}
	if update.Message != nil {
		return t.storage.UpdateTelegramNick(ctx, update.Message.Chat.ID, &update.Message.Chat.UserName)
	}
	return nil
}

func (t *Service) commonProcessor(ctx context.Context, meta *MetaData) (err error) {
	var (
		tree     *tgt.NodesHandler
		keyboard *telegram.KeyBoardType
	)

	if isCloseOrIgnoreCallback(meta.Callback) {
		if isCloseCallback(meta.Callback) {
			t.telegramClient.DeleteMessage(meta.MessageID, meta.ChatID)
		}
		return nil
	}

	if tree, err = t.genNewTree(ctx, meta.ChatID); err != nil {
		return trace.FuncNameWithErrorMsg(err, "getting tree")
	}

	node, err := tree.GetNode(ctx, meta)
	if err != nil {
		t.telegramClient.DeleteMessage(meta.MessageID, meta.ChatID)
		return trace.FuncNameWithErrorMsg(err, "getting node")
	}

	callbackSkip, err := node.GetCallbackSkip()
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "getting skip")
	}

	callbackBack, err := node.GetCallbackBack()
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "getting back")
	}

	if node.GetTelegram().DeleteMessage() {
		t.telegramClient.DeleteMessage(meta.MessageID, meta.ChatID)
	}
	if node.GetTelegram().GetResendMsg() {
		meta.MessageID = 0
	}

	if len(node.GetChildren()) == 0 {
		return nil
	}

	message := meta.Message
	if message == "" {
		message = node.GetTelegram().GetMessage()
	}

	keyboard, err = telegram.NewKeyboard(&telegram.NewKeyboardRequest{
		NextNodes:    node.GetChildren(),
		Meta:         meta,
		HideBar:      node.GetTelegram().GetHideBar(),
		CallbackSkip: callbackSkip,
		CallbackBack: callbackBack,
	})
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "creating keyboard")
	}

	_ = t.telegramClient.SendOrUpdateMessage(meta.ChatID, meta.MessageID, keyboard, message, !node.GetTelegram().GetEnablePreview())
	return nil
}

func (t *Service) genNewTree(ctx context.Context, chatID int64) (treeHandler *tgt.NodesHandler, err error) {
	treeHandler, err = tgt.NewNodesHandler(t.generateRootNodes(ctx, chatID), "Выбери:")
	if err != nil {
		return nil, fmt.Errorf("master.genNewTree creating handler: %w", err)
	}
	if treeHandler == nil {
		return nil, fmt.Errorf("master.genNewTree handler is nil")
	}
	return treeHandler, nil
}

func (t *Service) processNewUserStart(ctx context.Context, update *tgbotapi.Update) (err error) {
	newUser, err := t.storage.InitUser(ctx, storageSdk.User{
		TelegramNick: sql.NullString{Valid: update.Message.Chat.UserName != "", String: update.Message.Chat.UserName},
		ChatID:       sql.NullInt64{Valid: true, Int64: update.Message.Chat.ID},
	})
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "init user")
	}

	//startPaymentMessage := "Вас приветствует бот Falaemae😎\nЯ помогу Вам! Выберите, пожалуйста:"

	//err = t.ProcessCommandAsCallback(ctx, &MetaData{
	//	ChatID:   update.Message.Chat.ID,
	//	Callback: entities.PaymentRoot,
	//	Message:  startPaymentMessage,
	//})
	//if err != nil {
	//	return trace.FuncNameWithErrorMsg(err, "process payment")
	//}
	_ = t.sendAdminMessagesAboutNewUser(ctx, newUser)
	return nil
}

func (t *Service) sendAdminMessagesAboutNewUser(ctx context.Context, user storageSdk.User) error {
	payload, err := t.storage.GetUserPayload(ctx, user.ID, nil)
	if err != nil {
		return fmt.Errorf("tgupdates.sendAdminMessagesAboutNewUser error getting payload: %w", err)
	}

	if payload != "" {
		_ = t.telegramClient.SendMessageForAdmins(telegram.EventNewUser, payload)
	}
	return nil
}
