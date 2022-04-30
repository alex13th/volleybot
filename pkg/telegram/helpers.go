package telegram

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
)

type HelperError struct {
	Msg       string
	AnswerMsg string
}

func (e HelperError) Error() string {
	return e.Msg
}

type KeyboardHelper interface {
	GetText() string
	GetKeyboard() [][]InlineKeyboardButton
	GetBtnData(interface{}) string
	SetData(string) KeyboardHelper
}

func NewDateKeyboardHelper(msg string, prefix string) DateKeyboardHelper {
	return DateKeyboardHelper{
		Days: 6, Columns: 2, DateFormat: "Mon, 02.01", Locale: monday.LocaleRuRU,
		Msg: msg, Prefix: prefix}
}

type DateKeyboardHelper struct {
	Msg        string
	Prefix     string
	Data       string
	BackData   string
	Days       int
	DateFormat string
	Columns    int
	Locale     monday.Locale
}

func (kh DateKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh DateKeyboardHelper) SetData(data string) KeyboardHelper {
	kh.Data = data
	return kh
}

func (kh DateKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s", kh.Prefix, dt.Format("2006-02-01"), kh.Data)
}

func (kh DateKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	currDate := time.Now()
	for i := 1; i <= kh.Days; i++ {
		btnDate := currDate.AddDate(0, 0, i)
		btnText := monday.Format(btnDate, kh.DateFormat, kh.Locale)
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText, CallbackData: kh.GetBtnData(btnDate)})
		if i%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
		}
	}
	if len(kbdRow) > 0 {
		kbd = append(kbd, kbdRow)
	}
	if kh.BackData != "" {
		kbdRow := []InlineKeyboardButton{}
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: "Назад", CallbackData: kh.BackData})
		kbd = append(kbd, kbdRow)
	}
	return
}

func (dkh DateKeyboardHelper) Parse(Data string) (d time.Time, data string, err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 3 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse date"}
		return
	}
	d, perr := time.Parse("2006-02-01", splitedData[1])
	if perr != nil {
		err = HelperError{
			Msg:       fmt.Sprintf("parse date error: %s", perr.Error()),
			AnswerMsg: "Can't parse date"}
		return
	}
	data = splitedData[2]
	return
}

func NewTimeKeyboardHelper(msg string, prefix string) TimeKeyboardHelper {
	return TimeKeyboardHelper{StartHour: 7, EndHour: 21,
		Columns: 3, TimeFormat: "15:04", Locale: monday.LocaleRuRU,
		Msg: msg, Prefix: prefix}
}

type TimeKeyboardHelper struct {
	Msg        string
	Prefix     string
	Data       string
	BackData   string
	StartHour  int
	EndHour    int
	TimeFormat string
	Columns    int
	Locale     monday.Locale
}

func (kh TimeKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh TimeKeyboardHelper) SetData(data string) KeyboardHelper {
	kh.Data = data
	return kh
}

func (kh TimeKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s", kh.Prefix, dt.Format("15:04"), kh.Data)
}

func (kh TimeKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	for i := kh.StartHour; i <= kh.EndHour; i++ {
		btnTime := time.Date(0, 0, 0, i, 0, 0, 0, time.Local)
		btnText := monday.Format(btnTime, kh.TimeFormat, kh.Locale)
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText, CallbackData: kh.GetBtnData(btnTime)})
		if i%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
		}
	}
	if len(kbdRow) > 0 {
		kbd = append(kbd, kbdRow)
	}
	if kh.BackData != "" {
		kbdRow := []InlineKeyboardButton{}
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: "Назад", CallbackData: kh.BackData})
		kbd = append(kbd, kbdRow)
	}
	return
}

func (tkh TimeKeyboardHelper) Parse(Data string) (t time.Time, data string, err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 3 {
		err = HelperError{
			Msg:       "incorrect CallbackQuery data time format",
			AnswerMsg: "Can't parse time"}
		return
	}
	t, perr := time.Parse("15:04", splitedData[1])
	if perr != nil {
		err = HelperError{
			Msg:       fmt.Sprintf("parse time error: %s", perr.Error()),
			AnswerMsg: "Can't parse date"}
		return
	}
	data = splitedData[2]
	return
}

func NewCountKeyboardHelper(msg string, prefix string, min int, max int) CountKeyboardHelper {
	return CountKeyboardHelper{Min: min, Max: max, step: 1,
		Columns: 4, Msg: msg, Prefix: prefix}
}

type CountKeyboardHelper struct {
	Msg      string
	Prefix   string
	BackData string
	Data     string
	Min      int
	Max      int
	step     int
	Columns  int
}

func (kh CountKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh CountKeyboardHelper) SetData(data string) KeyboardHelper {
	kh.Data = data
	return kh
}

func (kh CountKeyboardHelper) GetBtnData(val interface{}) string {
	count := val.(int)
	return fmt.Sprintf("%s_%s_%s", kh.Prefix, strconv.Itoa(count), kh.Data)
}

func (kh CountKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for i := kh.Min; i <= kh.Max; i = i + kh.step {
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: strconv.Itoa(i), CallbackData: kh.GetBtnData(i)})
		count++
		if count%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
			count = 0
		}
	}
	if len(kbdRow) > 0 {
		kbd = append(kbd, kbdRow)
	}
	if kh.BackData != "" {
		kbdRow := []InlineKeyboardButton{}
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: "Назад", CallbackData: kh.BackData})
		kbd = append(kbd, kbdRow)
	}
	return
}

func (tkh CountKeyboardHelper) Parse(Data string) (count int, data string, err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 3 {
		err = HelperError{
			Msg:       "incorrect CallbackQuery data integer format",
			AnswerMsg: "Can't parse count"}
		return
	}
	count, perr := strconv.Atoi(splitedData[1])
	if perr != nil {
		err = HelperError{
			Msg:       fmt.Sprintf("parse integer error: %s", perr.Error()),
			AnswerMsg: "Can't parse count"}
		return
	}
	data = splitedData[2]
	return
}

type ActionButton struct {
	Text   string
	Prefix string
}

type ActionsKeyboardHelper struct {
	Msg      string
	Data     string
	BackData string
	Columns  int
	Actions  []ActionButton
}

func (kh ActionsKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh ActionsKeyboardHelper) SetData(data string) KeyboardHelper {
	kh.Data = data
	return kh
}

func (kh ActionsKeyboardHelper) GetBtnData(val interface{}) string {
	act := val.(ActionButton)
	return fmt.Sprintf("%s_%s", act.Prefix, kh.Data)
}

func (kh ActionsKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	for i, act := range kh.Actions {
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: act.Text, CallbackData: kh.GetBtnData(act)})
		if (i+1)%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
		}
	}
	if len(kbdRow) > 0 {
		kbd = append(kbd, kbdRow)
	}
	if kh.BackData != "" {
		kbdRow := []InlineKeyboardButton{}
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: "Назад", CallbackData: kh.BackData})
		kbd = append(kbd, kbdRow)
	}
	return
}

func (tkh ActionsKeyboardHelper) Parse(Data string) (data string, err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect CallbackQuery data format",
			AnswerMsg: "Can't parse data"}
		return
	}
	data = splitedData[1]
	return
}
