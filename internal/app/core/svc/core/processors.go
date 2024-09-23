package core

import (
	"context"
	"encoding/json"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	storageSdk "github.com/tarmalonchik/falaemae/internal/pkg/storage"
	"github.com/tarmalonchik/falaemae/internal/pkg/tools"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

func (t *Service) processCreateDrive(ctx context.Context, update *tgbotapi.Update, base64Data, carInfo string) error {
	var driveData tgCallbackDriveData

	data, err := tools.ParseBase64(base64Data)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "corrupted base64 data")
	}

	if err = json.Unmarshal(data, &driveData); err != nil {
		return trace.FuncNameWithErrorMsg(err, "unmarshal base64 info")
	}

	driver, err := t.storage.GetUserByChatID(ctx, update.Message.Chat.ID)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "getting data")
	}

	id, err := t.storage.CreateOrResolveCarID(ctx, carInfo)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "resolving car info")
	}

	_, err = t.storage.CreateDrive(ctx, storageSdk.Drive{
		DriverID:  driver.ID,
		Time:      time.Unix(driveData.Time, 0),
		Direction: driveData.Direction,
		Car:       id,
		Price:     int64(driveData.Price * 100),
		Slots:     driveData.CarSlots,
	})
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "creating drive")
	}

	t.telegramClient.DeleteMessage(driveData.MessageID, update.Message.Chat.ID)
	t.telegramClient.DeleteMessage(int64(update.Message.MessageID), update.Message.Chat.ID)

	// created drive post

	return nil
}
