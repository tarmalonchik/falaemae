package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var engMonthToRus = map[string]string{
	"January":   "Января",
	"February":  "Февраля",
	"March":     "Марта",
	"April":     "Апреля",
	"May":       "Мая",
	"June":      "Июня",
	"July":      "Июля",
	"August":    "Августа",
	"September": "Сентября",
	"October":   "Октября",
	"November":  "Ноября",
	"December":  "Декабря",
}

type Date struct {
	time time.Time
}

func NewDate() *Date {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	return &Date{
		time: time.Now().In(moscowLocation),
	}
}

func NewDataWithTime(inTime time.Time) *Date {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	return &Date{
		time: inTime.In(moscowLocation),
	}
}

func (d *Date) Time() time.Time {
	return d.time
}

func (d *Date) Incr() {
	d.time = d.time.Add(24 * time.Hour)
}

func (d *Date) PrettyPrinted() (out string) {
	y, m, day := d.time.Date()
	out += fmt.Sprintf("%d ", day)
	out += fmt.Sprintf("%s ", engMonthToRus[m.String()])
	out += fmt.Sprintf("%d г. ", y)
	return out
}

func (d *Date) PrettyPrintedDayMonth() (out string) {
	_, m, day := d.time.Date()
	out += fmt.Sprintf("%d ", day)
	out += fmt.Sprintf("%s ", engMonthToRus[m.String()])
	return out
}

func (d *Date) PrettyPrintedHHMM() (out string) {
	h := d.time.Hour()
	m := d.time.Minute()
	out += fmt.Sprintf("%s:%s ", hourMinuteToString(h), hourMinuteToString(m))
	return out
}

func (d *Date) PayloadPrinted() (out string) {
	y, m, day := d.time.Date()
	return fmt.Sprintf("%d:%d:%d", y, m, day)
}

func (d *Date) ParsePayloadPrinted(input string) (out string) {
	items := strings.Split(input, ":")
	if len(items) != 3 {
		return input
	}

	year, err := strconv.Atoi(items[0])
	if err != nil {
		return input
	}

	monthInt, err := strconv.Atoi(items[1])
	if err != nil {
		return input
	}

	month := time.Month(monthInt)

	day, err := strconv.Atoi(items[2])
	if err != nil {
		return input
	}

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	d.time = time.Date(year, month, day, 0, 0, 0, 0, moscowLocation)

	return d.PrettyPrinted()
}

func hourMinuteToString(in int) string {
	if in < 10 {
		return "0" + strconv.Itoa(in)
	}
	return strconv.Itoa(in)
}
