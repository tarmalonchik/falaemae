package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	tgt "github.com/vkidmode/telegram_tree"

	"github.com/tarmalonchik/falaemae/internal/entities"
	storageSdk "github.com/tarmalonchik/falaemae/internal/pkg/storage"
	"github.com/tarmalonchik/falaemae/internal/pkg/tools"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

func (t *Service) generateDriverRoot() tgt.Node {
	return tgt.NewNode(
		tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Админка"))),
		tgt.WithProc(processorWrap(t.generateDriverNodes)),
	)
}

func (t *Service) generateDriverNodes(ctx context.Context, meta MetaData) ([]tgt.Node, error) {
	out := make([]tgt.Node, 0, 3)

	user, err := t.storage.GetUserByChatID(ctx, meta.ChatID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting user")
	}

	drives, err := t.storage.GetLast10Drives(ctx, user.ID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting drives")
	}

	out = append(out, []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("Цхинвал → Владиквказ"),
					tgt.WithMessage(fmt.Sprintf("• Направление: <%s>\n\nВыберите дату:", directionPayloadKey)),
					tgt.WithColumns(2),
				),
			),
			tgt.WithProc(processorWrap(t.driverSetDate)),
			tgt.WithPayload(tgt.NewPayload(directionPayloadKey, strconv.Itoa(int(entities.DirectionTypeToTskhinvali)))),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("Владиквказ → Цхинвал"),
					tgt.WithMessage(fmt.Sprintf("• Направление: <%s>\n\nВыберите дату:", directionPayloadKey)),
					tgt.WithColumns(2),
				),
			),
			tgt.WithProc(processorWrap(t.driverSetDate)),
			tgt.WithPayload(tgt.NewPayload(directionPayloadKey, strconv.Itoa(int(entities.DirectionTypeToVladikavkaz)))),
		),
	}...)

	if len(drives) > 0 {
		out = append(out, tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("Мои поездки"),
				),
			),
			tgt.WithProc(processorWrap(t.driverShowMyDrives)),
		))
	}
	return out, nil
}

func (t *Service) driverShowMyDrives(ctx context.Context, meta MetaData) ([]tgt.Node, error) {
	driver, err := t.storage.GetUserByChatID(ctx, meta.ChatID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting driver")
	}

	drives, err := t.storage.GetLast10Drives(ctx, driver.ID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting drives")
	}

	out := make([]tgt.Node, 0, len(drives))

	for i := range drives {
		date := tools.NewDataWithTime(drives[i].Time)

		tabTxt := date.PrettyPrintedDayMonth()
		tabTxt += " " + date.PrettyPrintedHHMM()
		if drives[i].Direction == entities.DirectionTypeToVladikavkaz {
			tabTxt += "Цх → Вл "
		} else {
			tabTxt += "Вл → Цх "
		}

		out = append(out, tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt(tabTxt),
				)),
			tgt.WithProc(processorWrap(t.driverSetHours)),
		))
	}
	return out, nil
}

func (t *Service) driverSetDate(_ context.Context, meta MetaData) ([]tgt.Node, error) {
	const days = 10

	out := make([]tgt.Node, 0, days)

	today := tools.NewDate()

	for i := 0; i < days; i++ {
		out = append(out, tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt(today.PrettyPrintedDayMonth()),
					tgt.WithMessage(fmt.Sprintf("• Направление: <%s>\n• Дата: <%s>\n\nВыберите час", directionPayloadKey, datePayloadKey)),
					tgt.WithColumns(4),
				)),
			tgt.WithPayload(tgt.NewPayload(datePayloadKey, today.PayloadPrinted())),
			tgt.WithProc(processorWrap(t.driverSetHours)),
		))
		today.Incr()
	}
	return out, nil
}

func (t *Service) driverSetHours(ctx context.Context, meta MetaData) ([]tgt.Node, error) {
	msg := tgt.WithMessage(
		fmt.Sprintf("• Направление: <%s>\n• Дата: <%s>\n• Время: <%s>\n\nВыберите минуты",
			directionPayloadKey,
			datePayloadKey,
			hoursPayloadKey,
		),
	)

	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("00:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "00")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("01:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "01")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("02:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "02")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("03:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "03")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("04:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "04")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("05:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "05")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("06:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "06")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("07:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "07")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("08:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "08")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("09:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "09")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("10:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "10")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("11:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "11")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("12:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "12")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("13:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "13")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("14:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "14")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("15:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "15")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("16:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "16")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("17:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "17")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("18:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "18")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("19:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "19")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("20:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "20")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("21:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "21")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("22:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "22")),
		),
		tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt("23:00"),
					msg,
				)),
			tgt.WithProc(processorWrap(t.driverSetMinutes)),
			tgt.WithPayload(tgt.NewPayload(hoursPayloadKey, "23")),
		),
	}, nil
}

func (t *Service) driverSetMinutes(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	msg := tgt.WithMessage(
		fmt.Sprintf("• Направление: <%s>\n• Дата: <%s>\n• Время: <%s>:<%s>\n\nВыберите стоимость проезда",
			directionPayloadKey,
			datePayloadKey,
			hoursPayloadKey,
			minutesPayloadKey,
		),
	)

	return []tgt.Node{
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("0 минут"),
				tgt.WithColumns(2),
				msg,
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "00")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("10 минут"),
				tgt.WithColumns(2),
				msg,
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "10")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("20 минут"),
				tgt.WithColumns(2),
				msg,
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "20")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("30 минут"),
				tgt.WithColumns(2),
				msg,
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "30")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("40 минут"),
				tgt.WithColumns(2),
				msg,
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "40")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(
				tgt.WithTabTxt("50 минут"),
				msg,
				tgt.WithColumns(2),
			)),
			tgt.WithPayload(tgt.NewPayload(minutesPayloadKey, "50")),
			tgt.WithProc(processorWrap(t.driverSetPrice)),
		),
	}, nil
}

func (t *Service) driverSetPrice(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	msg := tgt.WithMessage(
		fmt.Sprintf("• Направление: <%s>\n• Дата: <%s>\n• Время: <%s>:<%s>\n• Стоимость: <%s>\n\nВыберите количество доступных мест",
			directionPayloadKey,
			datePayloadKey,
			hoursPayloadKey,
			minutesPayloadKey,
			pricePayloadKey,
		),
	)

	return []tgt.Node{
		tgt.NewNode(

			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("0 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "0")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("100 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "1")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("200 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "2")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("300 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "3")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("400 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "4")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("500 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "5")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("600 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "6")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("700 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "7")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("800 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "8")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("900 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "9")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("1000 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "10")),
		),
		tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("1100 рублей"), msg, tgt.WithColumns(2))),
			tgt.WithProc(processorWrap(t.driverSelectPlacesCount)),
			tgt.WithPayload(tgt.NewPayload(pricePayloadKey, "11")),
		),
	}, nil
}

func (t *Service) driverSelectPlacesCount(_ context.Context, _ MetaData) (out []tgt.Node, err error) {
	msg := tgt.WithMessage(
		fmt.Sprintf("• Направление: <%s>\n• Дата: <%s>\n• Время: <%s>:<%s>\n• Стоимость: <%s>\n• Свободно мест: <%s>\n\nВыберите марку автомобиля",
			directionPayloadKey,
			datePayloadKey,
			hoursPayloadKey,
			minutesPayloadKey,
			pricePayloadKey,
			carSlotsPayloadKey,
		),
	)

	out = make([]tgt.Node, 0, len(carFleeSlotsRus))
	for i := 1; i <= 11; i++ {
		out = append(out, tgt.NewNode(
			tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt(carFleeSlotsRus[strconv.Itoa(i)]), msg)),
			tgt.WithProc(processorWrap(t.driverSelectCar)),
			tgt.WithPayload(tgt.NewPayload(carSlotsPayloadKey, strconv.Itoa(i))),
		))
	}
	return out, nil
}

func (t *Service) driverSelectCar(ctx context.Context, meta MetaData) ([]tgt.Node, error) {
	payload, err := tgt.ExtractPayload(meta.GetCallback())
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting payload")
	}

	data := &tgCallbackDriveData{}
	if err = data.fillFromPayload(payload, meta); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "parsing from payload")
	}

	val, err := json.Marshal(data)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "payload json marshal")
	}

	valString := tools.ToBase64(val)

	driver, err := t.storage.GetUserByChatID(ctx, meta.ChatID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting driver")
	}

	carsID, err := t.storage.GetLastCars(ctx, driver.ID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting drives")
	}

	cars, err := t.storage.GetCarsByID(ctx, carsID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting cars by id")
	}

	out := make([]tgt.Node, 0, len(cars)+1)

	for i := range cars {
		out = append(out, tgt.NewNode(
			tgt.WithTg(
				tgt.NewTelegram(
					tgt.WithTabTxt(cars[i].Name), tgt.DeleteMsg()),
			),
			tgt.WithProc(processorWrap(t.driverCachedCarProcessor)),
			tgt.WithPayload(tgt.NewPayload(carModelPayloadKey, strconv.Itoa(int(cars[i].ID)))),
		))
	}

	out = append(out, tgt.NewNode(
		tgt.WithTg(
			tgt.NewTelegram(
				tgt.WithTabTxt(fmt.Sprintf("Ввести данные о марке машины")),
				tgt.WithSwitchInline(fmt.Sprintf(customMessageWithPayloadFormat, valString, enterCarModelText)),
			),
		),
	))
	return out, nil
}

func (t *Service) driverCachedCarProcessor(ctx context.Context, meta MetaData) (out []tgt.Node, err error) {
	payload, err := tgt.ExtractPayload(meta.GetCallback())
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting payload")
	}

	data := &tgCallbackDriveData{}
	if err = data.fillFromPayload(payload, meta); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "parsing from payload")
	}

	driver, err := t.storage.GetUserByChatID(ctx, meta.ChatID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting data")
	}

	if data.CarID == nil {
		return nil, trace.FuncNameWithError(errors.New("car is null"))
	}

	car, err := t.storage.GetCarByID(ctx, *data.CarID)
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting car by id")
	}

	_, err = t.storage.CreateDrive(ctx, storageSdk.Drive{
		DriverID:  driver.ID,
		Time:      time.Unix(data.Time, 0),
		Direction: data.Direction,
		Car:       car.ID,
		Price:     int64(data.Price * 100),
		Slots:     data.CarSlots,
	})
	if err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "creating drive")
	}

	t.telegramClient.DeleteMessage(meta.MessageID, meta.ChatID)

	return nil, nil
}
