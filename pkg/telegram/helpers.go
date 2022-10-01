package telegram

import (
	"fmt"
	"strconv"
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

type MessageRequestHelper interface {
	GetEditMR() EditMessageTextRequest
	GetMR() MessageRequest
}

type CallbackDataParser interface {
	GetAction() string
	GetPrefix() string
	GetState() State
	GetValue() string
	Parse(string) error
	SetState(state State)
}

type KeyboardHelper interface {
	GetKeyboard() [][]InlineKeyboardButton
	GetText() string
}

type BaseKeyboardHelper struct {
	State
	BackData string
	Text     string
}

func (kh BaseKeyboardHelper) GetText() string {
	return kh.Text
}

func (kh *BaseKeyboardHelper) SetBackData(data string) {
	kh.BackData = data
}

type DateTimeResources struct {
	DateBtn  string `json:"date_btn"`
	DateMsg  string `json:"date_msg"`
	DayCount int    `json:"daye_count"`
	TimeBtn  string `json:"time_btn"`
	TimeMsg  string `json:"time_msg"`
}

func NewDateTimeResourcesRu() DateTimeResources {
	return DateTimeResources{
		DateBtn: "📆 Дата", DateMsg: "❓Какая дата❓", DayCount: 30,
		TimeBtn: "⏰ Время", TimeMsg: "❓В какое время❓",
	}
}

func NewDateKeyboardHelper() DateKeyboardHelper {
	kh := DateKeyboardHelper{Days: 6, Columns: 2, DateFormat: "Mon, 02.01", Location: time.Local}
	return kh
}

func NewDateKeyboardHelperRu() (h DateKeyboardHelper) {
	h = NewDateKeyboardHelper()
	h.Locale = monday.LocaleRuRU
	return
}

type DateKeyboardHelper struct {
	BaseKeyboardHelper
	Date       time.Time
	Location   *time.Location
	Days       int
	DateFormat string
	Columns    int
	Locale     monday.Locale
}

func (kh *DateKeyboardHelper) Parse() (err error) {
	if kh.Action == "set" {
		if kh.Date, err = time.ParseInLocation("2006-01-02", kh.Value, kh.Location); err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse date error: %s", err.Error()),
				AnswerMsg: "Can't parse date"}
			return
		}
		kh.Date = time.Date(kh.Date.Year(), kh.Date.Month(), kh.Date.Day(), 0, 0, 0, 0, kh.Location)
	}
	return
}

func (kh DateKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	currDate := time.Now()
	for i := 1; i <= kh.Days; i++ {
		st := kh.State
		btnDate := currDate.AddDate(0, 0, i-1)
		btnText := monday.Format(btnDate, kh.DateFormat, kh.Locale)
		st.Action = "set"
		st.Value = btnDate.Format("2006-01-02")

		kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText,
			CallbackData: st.String()})
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

func NewTimeKeyboardHelper() TimeKeyboardHelper {
	kh := TimeKeyboardHelper{StartHour: 7, EndHour: 21, Columns: 3, TimeFormat: "15:04", Location: time.Local}
	return kh
}

func NewTimeKeyboardHelperRu() (h TimeKeyboardHelper) {
	h = NewTimeKeyboardHelper()
	h.Locale = monday.LocaleRuRU
	return
}

type TimeKeyboardHelper struct {
	BaseKeyboardHelper
	Location    *time.Location
	Time        time.Time
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
	Step        int
	TimeFormat  string
	Columns     int
	Locale      monday.Locale
}

func (kh *TimeKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	if kh.StartHour < 0 {
		kh.StartHour = 24 + kh.StartHour
	}
	if kh.StartMinute < 0 {
		kh.StartMinute = 60 + kh.StartMinute
	}
	if kh.EndHour < 0 {
		kh.EndHour = 24 + kh.EndHour
	}
	if kh.EndMinute < 0 {
		kh.EndMinute = 60 + kh.EndMinute
	}
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for i := kh.StartHour; i <= kh.EndHour; i++ {
		btnTime := time.Date(0, 0, 0, i, 0, 0, 0, time.Local)
		count++
		btnText := monday.Format(btnTime, kh.TimeFormat, kh.Locale)
		st := kh.State
		st.Value = btnTime.Format("15:04")
		st.Action = "set"
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText,
			CallbackData: st.String()})
		if count%kh.Columns == 0 {
			kbd = append(kbd, kbdRow)
			kbdRow = []InlineKeyboardButton{}
		}
		if kh.Step > 0 {
			for i := kh.Step; i+kh.Step <= 60; i += kh.Step {
				count++
				btnTime = btnTime.Add(time.Minute * time.Duration(kh.Step))
				btnText := monday.Format(btnTime, kh.TimeFormat, kh.Locale)
				st.Value = btnTime.Format("15:04")
				st.Action = "set"
				kbdRow = append(kbdRow, InlineKeyboardButton{Text: btnText,
					CallbackData: st.String()})
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

func (kh *TimeKeyboardHelper) Parse() (err error) {
	if kh.Action == "set" {
		kh.Time, err = time.Parse("15:04", kh.Value)
		if err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse time error: %s", err.Error()),
				AnswerMsg: "Can't parse time"}
			return
		}
		kh.Time = time.Date(0, 0, 0, kh.Time.Hour(), kh.Time.Minute(), 0, 0, kh.Location)
	}
	return
}

func NewCountKeyboardHelper() CountKeyboardHelper {
	kh := CountKeyboardHelper{Min: 1, Max: 4, Step: 1, Columns: 4}
	return kh
}

type CountKeyboardHelper struct {
	BaseKeyboardHelper
	AlwaysZero bool
	Count      int
	Min        int
	Max        int
	Step       int
	Columns    int
}

func (kh CountKeyboardHelper) GetButton(v int) (btn InlineKeyboardButton) {
	st := kh.State
	st.Action = "set"
	st.Value = strconv.Itoa(v)
	return InlineKeyboardButton{Text: strconv.Itoa(v), CallbackData: st.String()}
}

func (kh CountKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	if kh.Min*kh.Max > 0 && kh.AlwaysZero {
		kbdRow = append(kbdRow, kh.GetButton(0))
		kbd = append(kbd, kbdRow)
		kbdRow = []InlineKeyboardButton{}
	}
	count := 0
	for i := kh.Min; i <= kh.Max; i = i + kh.Step {
		kbdRow = append(kbdRow, kh.GetButton(i))
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

func (kh *CountKeyboardHelper) Parse() (err error) {
	if kh.Action == "set" {
		kh.Count, err = strconv.Atoi(kh.Value)
		if err != nil {
			err = HelperError{
				Msg:       fmt.Sprintf("parse count error: %s", err.Error()),
				AnswerMsg: "Can't parse count"}
			return
		}
	}
	return
}

type ActionButton struct {
	Text   string
	Data   string
	Action string
}

func NewActionsKeyboardHelper() ActionsKeyboardHelper {
	kh := ActionsKeyboardHelper{Columns: 2}
	return kh
}

type ActionsKeyboardHelper struct {
	BaseKeyboardHelper
	Columns int
	Actions []ActionButton
}

func (kh ActionsKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	if kh.Columns == 0 {
		kh.Columns = 1
	}
	kbdRow := []InlineKeyboardButton{}
	for i, act := range kh.Actions {
		st := kh.State
		st.Action = act.Action
		if act.Data != "" {
			st.Data = act.Data
		}
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: act.Text, CallbackData: st.String()})
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

type EnumItem struct {
	Id   string
	Item string
}

func NewEnumKeyboardHelper(enums []EnumItem) EnumKeyboardHelper {
	return EnumKeyboardHelper{Enums: enums, Columns: 2}
}

type EnumKeyboardHelper struct {
	BaseKeyboardHelper
	Choice  string
	Columns int
	Enums   []EnumItem
}

func (kh EnumKeyboardHelper) GetKeyboard() (kbd [][]InlineKeyboardButton) {
	kbdRow := []InlineKeyboardButton{}
	count := 0
	for _, val := range kh.Enums {
		st := kh.State
		st.Action = "set"
		st.Value = val.Id
		kbdRow = append(kbdRow, InlineKeyboardButton{Text: val.Item, CallbackData: st.String()})
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
