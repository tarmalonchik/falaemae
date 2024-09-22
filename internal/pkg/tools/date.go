package tools

import (
	"fmt"
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

func (d *Date) PayloadPrinted() (out string) {
	y, m, day := d.time.Date()
	return fmt.Sprintf("%d:%d:%d", y, m, day)
}
