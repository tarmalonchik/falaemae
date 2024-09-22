package core

import (
	"context"
	"fmt"

	tgt "github.com/vkidmode/telegram_tree"

	"github.com/tarmalonchik/falaemae/internal/pkg/tools"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

func (t *Service) generateDriverRoot() tgt.Node {
	return tgt.NewNode(
		tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Админка"))),
		tgt.WithProc(processorWrap(t.generateDriverNodes)),
	)
}

func (t *Service) generateDriverNodes(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Запланировать поездку"))),
			tgt.WithProc(processorWrap(t.driverCreateDrive)),
		),
	}, nil
}

func (t *Service) driverCreateDrive(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(tgt.WithTabTxt("Указать дату"))),
			tgt.WithProc(processorWrap(t.driverSetDate)),
		),
	}, nil
}

func (t *Service) driverSetDate(_ context.Context, meta MetaData) ([]tgt.Node, error) {
	const days = 7

	out := make([]tgt.Node, 0, days)

	today := tools.NewDate()

	payload, err := tgt.ExtractPayload(meta.GetCallback())
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting payload")
	}

	fmt.Println(payload[datePayloadKey])

	for i := 0; i < days; i++ {
		out = append(out, tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt(fmt.Sprintf("Я еду %s", today.PrettyPrinted())),
					tgt.WithMessage("Выберите час"),
				)),
			tgt.WithPayload(tgt.NewPayload(datePayloadKey, today.PayloadPrinted())),
			tgt.WithProc(processorWrap(t.driverSetHours)),
		))
		today.Incr()
	}
	return out, nil
}

func (t *Service) driverSetHours(ctx context.Context, meta MetaData) ([]tgt.Node, error) {
	tree, err := t.genNewTree(ctx, meta.ChatID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting tree")
	}

	_, err = tree.GetNode(ctx, &meta)
	if err != nil {
		t.telegramClient.DeleteMessage(meta.MessageID, meta.ChatID)
		return nil, trace.FuncNameWithErrorMsg(err, "getting node")
	}

	//callbackSkip, err := node.GetCallbackSkip()
	//if err != nil {
	//	return nil, trace.FuncNameWithErrorMsg(err, "getting skip")
	//}
	//
	//callbackBack, err := node.GetCallbackBack()
	//if err != nil {
	//	return nil, trace.FuncNameWithErrorMsg(err, "getting back")
	//}
	//
	//keyboard, err := telegram.NewKeyboard(&telegram.NewKeyboardRequest{
	//	NextNodes:    node.GetChildren(),
	//	Meta:         &meta,
	//	HideBar:      node.GetTelegram().GetHideBar(),
	//	CallbackSkip: callbackSkip,
	//	CallbackBack: callbackBack,
	//})
	//if err != nil {
	//	return nil, trace.FuncNameWithErrorMsg(err, "creating keyboard")
	//}
	//_ = t.telegramClient.SendOrUpdateMessage(meta.ChatID, meta.MessageID, keyboard, "kaka сообщение", !node.GetTelegram().GetEnablePreview())

	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 12 часов ночью"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 1 час ночью"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 2 часа ночью"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 3 часа ночью"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 4 часа утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 5 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 6 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 7 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 8 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 10 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 11 часов утром"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		), tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 12 часов дня"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		), tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 1 час дня"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 2 часа дня"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 3 часа дня"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 4 часа дня"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 5 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 6 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 7 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 8 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 9 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 10 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("В 11 часов вечером"))),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
		),
	}, nil
}

func (t *Service) driverSetMinutes(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("0 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("10 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("20 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("30 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("40 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("50 минут"))),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
	}, nil
}

func (t *Service) driverSetPrice(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("0 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("100 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("200 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("300 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("400 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("500 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("600 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("700 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("800 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("900 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("1000 рублей"))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
		),
	}, nil
}

func (t *Service) driverSelectPlacesCount(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("1 место"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("2 места"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("3 места"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("4 места"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("5 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("6 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("7 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("8 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("9 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("10 мест"))),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
		),
	}, nil
}

func (t *Service) driverSelectCar(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		//tgt.NewNode(
		//	tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("1 место"))),
		//	tgt.WithProc(processorWrap(t.statisticsProcessor)),
		//),
	}, nil
}
