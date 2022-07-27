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
	GetBtnData(interface{}) string
	GetData() string
	GetKeyboard() [][]InlineKeyboardButton
	GetState() string
	GetText() string
	Parse(string) error
	SetData(string)
}

func NewDateKeyboardHelper(msg string, state string) DateKeyboardHelper {
	return DateKeyboardHelper{Action: "get",
		Days: 6, Columns: 2, DateFormat: "Mon, 02.01", Locale: monday.LocaleRuRU,
		Msg: msg, State: state}
}

type DateKeyboardHelper struct {
	Msg        string
	State      string
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

func (kh DateKeyboardHelper) GetState() string {
	return kh.State
}

func (kh DateKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh *DateKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh DateKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s_%s", kh.State, kh.Data, "set", dt.Format("2006-02-01"))
}

func (kh DateKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	currDate := time.Now()
	for i := 1; i <= kh.Days; i++ {
		btnDate := currDate.AddDate(0, 0, i-1)
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

func NewTimeKeyboardHelper(msg string, state string) TimeKeyboardHelper {
	return TimeKeyboardHelper{Action: "get", StartHour: 7, EndHour: 21,
		Columns: 3, TimeFormat: "15:04", Locale: monday.LocaleRuRU,
		Msg: msg, State: state}
}

type TimeKeyboardHelper struct {
	Msg         string
	State       string
	Action      string
	Time        time.Time
	Data        string
	BackData    string
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
	Step        int
	TimeFormat  string
	Columns     int
	Locale      monday.Locale
}

func (kh TimeKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh TimeKeyboardHelper) GetState() string {
	return kh.State
}

func (kh TimeKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh *TimeKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh TimeKeyboardHelper) GetBtnData(val interface{}) string {
	dt := val.(time.Time)
	return fmt.Sprintf("%s_%s_%s_%s", kh.State, kh.Data, "set", dt.Format("15:04"))
}

func (kh *TimeKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for i := kh.StartHour; i <= kh.EndHour; i++ {
		btnTime := time.Date(0, 0, 0, i, 0, 0, 0, time.Local)
		count++
		btnText := monday.Format(btnTime, kh.TimeFormat, kh.Locale)
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText, CallbackData: kh.GetBtnData(btnTime)})
		if count%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
		}
		if kh.Step > 0 {
			for i := kh.Step; i+kh.Step <= 60; i += kh.Step {
				count++
				btnTime = btnTime.Add(time.Minute * time.Duration(kh.Step))
				btnText := monday.Format(btnTime, kh.TimeFormat, kh.Locale)
				kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText, CallbackData: kh.GetBtnData(btnTime)})
				if count%kh.Columns == 0 {
					kbd = append(kbd, kbdRow)
					kbdRow = []InlineKeyboardButton{}
				}
			}
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

func NewCountKeyboardHelper(msg string, state string, min int, max int) CountKeyboardHelper {
	return CountKeyboardHelper{Min: min, Max: max, Step: 1,
		Columns: 4, Msg: msg, State: state}
}

type CountKeyboardHelper struct {
	Msg      string
	State    string
	Action   string
	Count    int
	BackData string
	Data     string
	Min      int
	Max      int
	Step     int
	Columns  int
}

func (kh CountKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh CountKeyboardHelper) GetState() string {
	return kh.State
}

func (kh *CountKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh CountKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh CountKeyboardHelper) GetBtnData(val interface{}) string {
	count := val.(int)
	return fmt.Sprintf("%s_%s_%s_%s", kh.State, kh.Data, "set", strconv.Itoa(count))
}

func (kh CountKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for i := kh.Min; i <= kh.Max; i = i + kh.Step {
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
	Data   string
	Action string
}

type ActionsKeyboardHelper struct {
	State    string
	Msg      string
	Data     string
	BackData string
	Columns  int
	Actions  []ActionButton
}

func (kh ActionsKeyboardHelper) GetText() string {
	return kh.Msg
}

func (kh ActionsKeyboardHelper) GetState() string {
	return kh.State
}

func (kh ActionsKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh *ActionsKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh ActionsKeyboardHelper) GetBtnData(val interface{}) string {
	act := val.(ActionButton)
	data := kh.Data
	if act.Data != "" {
		data = act.Data
	}
	return fmt.Sprintf("%s_%s", act.Action, data)
}

func (kh ActionsKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	if kh.Columns == 0 {
		kh.Columns = 1
	}
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

func (tkh *ActionsKeyboardHelper) Parse(Data string) (err error) {
	splitedData := strings.Split(Data, "_")
	if len(splitedData) < 2 {
		err = HelperError{
			Msg:       "incorrect CallbackQuery data format",
			AnswerMsg: "Can't parse data"}
		return
	}
	tkh.Data = splitedData[1]
	tkh.State = splitedData[0]
	return
}

type EnumItem struct {
	Id   string
	Item string
}

func NewEnumKeyboardHelper(msg string, state string, enums []EnumItem) EnumKeyboardHelper {
	return EnumKeyboardHelper{Enums: enums,
		Columns: 2, Msg: msg, State: state}
}

type EnumKeyboardHelper struct {
	Msg      string
	State    string
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

func (kh EnumKeyboardHelper) GetState() string {
	return kh.State
}

func (kh *EnumKeyboardHelper) SetData(data string) {
	kh.Data = data
}

func (kh EnumKeyboardHelper) GetData() string {
	return kh.Data
}

func (kh EnumKeyboardHelper) GetBtnData(val interface{}) string {
	return fmt.Sprintf("%s_%s_%s_%s", kh.State, kh.Data, "set", val)
}

func (kh EnumKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for _, val := range kh.Enums {
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: val.Item, CallbackData: kh.GetBtnData(val.Id)})
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
