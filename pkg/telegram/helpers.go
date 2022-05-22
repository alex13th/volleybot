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
	GetData() string
	SetData(string)
}

func NewDateKeyboardHelper(msg string, prefix string) DateKeyboardHelper {
	return DateKeyboardHelper{Action: "get",
		Days: 6, Columns: 2, DateFormat: "Mon, 02.01", Locale: monday.LocaleRuRU,
		Msg: msg, Prefix: prefix}
}

type DateKeyboardHelper struct {
	Msg        string
	Prefix     string
	Action     string
	Date       time.Time
	Data       string
	BackData   string
	Days       int
	DateFormat string
	Columns    int
	Locale     monday.Locale
}

func (kh DateKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh DateKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh *DateKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh DateKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s_%s", kh.Prefix, kh.Data, "set", dt.Format("2006-02-01"))
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

func (kh *DateKeyboardHelper) Parse(Data string) (err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse date"}
		return
	} else if len(splitedData) > 2 {
		kh.Action = splitedData[2]
	}
	if kh.Action == "set" {
		kh.Date, err = time.Parse("2006-02-01", splitedData[3])
		if err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse date error: %s", err.Error()),
				AnswerMsg: "Can't parse date"}
			return
		}
	}
	kh.Data = splitedData[1]
	return
}

func NewTimeKeyboardHelper(msg string, prefix string) TimeKeyboardHelper {
	return TimeKeyboardHelper{Action: "get", StartHour: 7, EndHour: 21,
		Columns: 3, TimeFormat: "15:04", Locale: monday.LocaleRuRU,
		Msg: msg, Prefix: prefix}
}

type TimeKeyboardHelper struct {
	Msg        string
	Prefix     string
	Action     string
	Time       time.Time
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

func (kh TimeKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh *TimeKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh TimeKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s_%s", kh.Prefix, kh.Data, "set", dt.Format("15:04"))
}

func (kh *TimeKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
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

func (kh *TimeKeyboardHelper) Parse(Data string) (err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse date"}
		return
	} else if len(splitedData) > 2 {
		kh.Action = splitedData[2]
	}

	if kh.Action == "set" {
		kh.Time, err = time.Parse("15:04", splitedData[3])
		if err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse time error: %s", err.Error()),
				AnswerMsg: "Can't parse time"}
			return
		}
	}
	kh.Data = splitedData[1]
	return
}

func NewCountKeyboardHelper(msg string, prefix string, min int, max int) CountKeyboardHelper {
	return CountKeyboardHelper{Min: min, Max: max, step: 1,
		Columns: 4, Msg: msg, Prefix: prefix}
}

type CountKeyboardHelper struct {
	Msg      string
	Prefix   string
	Action   string
	Count    int
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

func (kh *CountKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh CountKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh CountKeyboardHelper) GetBtnData(val interface{}) string {
	count := val.(int)
	return fmt.Sprintf("%s_%s_%s_%s", kh.Prefix, kh.Data, "set", strconv.Itoa(count))
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

func (kh *CountKeyboardHelper) Parse(Data string) (err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse count"}
		return
	} else if len(splitedData) > 2 {
		kh.Action = splitedData[2]
	}

	if kh.Action == "set" {
		kh.Count, err = strconv.Atoi(splitedData[3])
		if err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse count error: %s", err.Error()),
				AnswerMsg: "Can't parse count"}
			return
		}
	}
	kh.Data = splitedData[1]
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

func (kh ActionsKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh *ActionsKeyboardHelper) SetData(data string) {
	kh.Data = data
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

type EnumItem struct {
	Id   string
	Item fmt.Stringer
}

func NewEnumKeyboardHelper(msg string, prefix string, enums []EnumItem) EnumKeyboardHelper {
	return EnumKeyboardHelper{Enums: enums,
		Columns: 2, Msg: msg, Prefix: prefix}
}

type EnumKeyboardHelper struct {
	Msg      string
	Prefix   string
	Action   string
	Choice   string
	BackData string
	Data     string
	Enums    []EnumItem
	Columns  int
}

func (kh EnumKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh *EnumKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh EnumKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh EnumKeyboardHelper) GetBtnData(val interface{}) string {
	return fmt.Sprintf("%s_%s_%s_%s", kh.Prefix, kh.Data, "set", val)
}

func (kh EnumKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for _, val := range kh.Enums {
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: val.Item.String(), CallbackData: kh.GetBtnData(val.Id)})
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

func (kh *EnumKeyboardHelper) Parse(Data string) (err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse date"}
		return
	} else if len(splitedData) > 2 {
		kh.Action = splitedData[2]
	}

	if kh.Action == "set" {
		kh.Choice = splitedData[3]
	}
	kh.Data = splitedData[1]
	return
}