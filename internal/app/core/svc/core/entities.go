package core

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/tarmalonchik/falaemae/internal/entities"
	"github.com/tarmalonchik/falaemae/internal/pkg/tools"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

const (
	pleaseUseButtons               = "Пожалуйста воспользуйтесь кнопками"
	wrongMessage                   = "Обратитесь пожалуйста в поддержку @vpnchik_support"
	customMessageFormat            = "%s\n\n"
	customMessageWithPayloadFormat = "(%s) %s\n\n"
	enterCarModelText              = "Введите марку автомобиля ⬇️"
	datePayloadKey                 = "!"
	directionPayloadKey            = "?"
	hoursPayloadKey                = "*"
	minutesPayloadKey              = "#"
	pricePayloadKey                = "@"
	carSlotsPayloadKey             = "&"
	carModelPayloadKey             = "+"
)

var carFleeSlotsRus = map[string]string{
	"1":  "1 место",
	"2":  "2 места",
	"3":  "3 места",
	"4":  "4 места",
	"5":  "5 мест",
	"6":  "6 мест",
	"7":  "7 мест",
	"8":  "8 мест",
	"9":  "9 мест",
	"10": "10 мест",
}

type Config struct {
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" default:"40s"`
}

type commandsProcessorFunc func(ctx context.Context, meta MetaData) error

type processorWithBranch struct {
	processor    commandsProcessorFunc
	worksForever bool
}

type tgCallbackDriveData struct {
	Time      int64                  `json:"t"`
	Direction entities.DirectionType `json:"d"`
	Price     int                    `json:"p"`
	CarSlots  int                    `json:"c"`
	MessageID int64                  `json:"m"`
	CarID     *int                   `json:"car,omitempty"`
}

func (t *tgCallbackDriveData) fillFromPayload(payload map[string]string, meta MetaData) error {
	priceString, ok := payload[pricePayloadKey]
	if !ok {
		return trace.FuncNameWithError(errors.New("no price payload"))
	}
	price, err := strconv.Atoi(priceString)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "price payload corrupted")
	}

	slotsString, ok := payload[carSlotsPayloadKey]
	if !ok {
		return trace.FuncNameWithErrorMsg(err, "no slots payload")
	}
	slots, err := strconv.Atoi(slotsString)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "slots payload corrupted")
	}

	date, ok := payload[datePayloadKey]
	if !ok {
		return trace.FuncNameWithErrorMsg(err, "no date payload")
	}

	minutesString, ok := payload[minutesPayloadKey]
	if !ok {
		return trace.FuncNameWithErrorMsg(err, "no minutes payload")
	}
	minutes, err := strconv.Atoi(minutesString)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "minutes payload corrupted")
	}

	hoursString, ok := payload[hoursPayloadKey]
	if !ok {
		return trace.FuncNameWithErrorMsg(err, "no hours payload")
	}
	hours, err := strconv.Atoi(hoursString)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "hours payload corrupted")
	}

	directionString, ok := payload[directionPayloadKey]
	if !ok {
		return trace.FuncNameWithErrorMsg(err, "no direction payload")
	}
	directionInt, err := strconv.Atoi(directionString)
	if err != nil {
		return trace.FuncNameWithErrorMsg(err, "direction payload corrupted")
	}
	direction := entities.DirectionType(directionInt)

	carIDString, ok := payload[carModelPayloadKey]
	if ok {
		carID, err := strconv.Atoi(carIDString)
		if err != nil {
			return trace.FuncNameWithErrorMsg(err, "car payload corrupted")
		}
		t.CarID = &carID
	}

	payloadDate := tools.NewDate()
	payloadDate.ParsePayloadPrinted(date)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t.Time = time.Date(payloadDate.Time().Year(), payloadDate.Time().Month(), payloadDate.Time().Day(), hours, minutes, 0, 0, moscowLocation).UTC().Unix()
	t.Direction = direction
	t.CarSlots = slots
	t.Price = price
	t.MessageID = meta.MessageID
	return nil
}
