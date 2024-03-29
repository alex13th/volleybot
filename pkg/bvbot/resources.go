package bvbot

import (
	"time"
	"volleybot/pkg/telegram"
)

type Resources struct {
	Actions       ActionsResources
	Activity      AcivityResources
	Config        ConfigResources
	Courts        CourtsResources
	Cancel        CancelResources
	Description   DescResources
	Join          JoinResources
	Level         LevelResources
	List          ListResources
	Main          MainResources
	MaxPlayer     MaxPlayersResources
	Profile       ProfileResources
	RemovePlayer  RemovePlayerResources
	Price         PriceResources
	Settings      SettingsResources
	Sets          SetsResources
	Show          ShowResources
	SendResources SendResources
	BackBtn       string
	DescMessage   string
}

type MainResources struct {
	ListCaption       string        `json:"list_caption"`
	ListDateBtn       string        `json:"List_date_btn"`
	NewReserveBtn     string        `json:"new_reserve_msg"`
	NoReservesMessage string        `json:"no_reserve_msg"`
	ParseMode         string        `json:"parse_mode"`
	PreviewDuration   time.Duration `json:"duration"`
	ProfileBtn        string        `json:"profile_btn"`
	ConfigBtn         string        `json:"config_btn"`
	Text              string        `json:"text"`
	TodayBtn          string        `json:"today_btn"`
}

func NewMainResourcesRu() (r MainResources) {
	r.ListCaption = "* Ближайшие активности *"
	r.ListDateBtn = "Найти по дате"
	r.NewReserveBtn = "✨ Забронировать"
	r.NoReservesMessage = "На ближайшее время активности не запланированы"
	r.Text = "Выберите действие"
	r.ParseMode = "Markdown"
	r.ProfileBtn = "😎 Профиль"
	r.ConfigBtn = "🛠 Настройки"
	r.TodayBtn = "Сегодня"
	return
}

type ListResources struct {
	ListCaption       string `json:"list_caption"`
	NoReservesMessage string `json:"no_reserve_message"`
	ParseMode         string `json:"parse_mode"`
	Text              string `json:"text"`
}

func NewListResourcesRu() (r ListResources) {
	r.ListCaption = "* Ближайшие активности *"
	r.NoReservesMessage = "На ближайшее время активности не запланированы"
	r.Text = "Выберите действие"
	r.ParseMode = "Markdown"
	return
}

type ShowResources struct {
	DateTime       telegram.DateTimeResources
	ActionsBtn     string
	DescriptionBtn string
	JoinBtn        string
	JoinLeaveBtn   string
	JoinMultiBtn   string
	JoinTimeBtn    string
	PayBtn         string
	RefreshBtn     string
	SetsBtn        string
	SettingsBtn    string
}

func NewShowResourcesRu() (r ShowResources) {
	r.DateTime = telegram.NewDateTimeResourcesRu()
	r.ActionsBtn = "Действия"
	r.DescriptionBtn = "Описание"
	r.JoinBtn = "😀 Буду"
	r.JoinLeaveBtn = "😞 Не смогу"
	r.JoinMultiBtn = "🤩 Буду не один"
	r.JoinTimeBtn = "🏃‍♂️ Опоздаю"
	r.PayBtn = "💰 Оплатить"
	r.RefreshBtn = "Обновить"
	r.SetsBtn = "⏱ Кол-во часов"
	r.SettingsBtn = "Настройки"
	return
}

type ActionsResources struct {
	BackBtn         string `json:"back_btn"`
	CancelBtn       string `json:"cancel_btn"`
	CopyBtn         string `json:"copy_btn"`
	CopyDoneMessage string `json:"copy_done_msg"`
	PaidBtn         string `json:"paid"`
	PublishBtn      string `json:"publish_btn"`
	SendBtn         string `json:"send_btn"`
	RemovePlayerBtn string `json:"remove_player_btn"`
}

func NewActionsResourcesRu() (r ActionsResources) {
	r.BackBtn = "Назад"
	r.CancelBtn = "💥Отменить"
	r.CopyBtn = "🫂 Копировать"
	r.CopyDoneMessage = "Копия сделана! 👆"
	r.PaidBtn = "💰 Оплаты"
	r.PublishBtn = "Опубликовать"
	r.SendBtn = "Отправить"
	r.RemovePlayerBtn = "Удалить игрока"
	return
}

type CancelResources struct {
	AbortBtn   string `json:"abort_btn"`
	BackBtn    string `json:"back_btn"`
	Text       string `json:"text"`
	ConfirmBtn string `json:"confirm_btn"`
}

func NewCancelResourcesRu() (r CancelResources) {
	r.BackBtn = "Передумал"
	r.ConfirmBtn = "🧨 Уверен"
	r.Text = "\n🧨*ВНИМАНИЕ!!!*🧨\nИгра будет отменена для всех участников. Если есть желание только выписаться, лучше воспользоваться кнопкой \"Не буду\""
	return
}

type SendResources struct {
	Message string `json:"message"`
	SendBtn string `json:"send_btn"`
}

func NewSendResourcesRu() (r SendResources) {
	r.Message = "Выберите чат для отправки объявления"
	r.SendBtn = "Куда отправим"
	return
}

type RemovePlayerResources struct {
	BackBtn         string
	Message         string
	RemovePlayerBtn string
}

func RemovePlayerResourcesRu() (r RemovePlayerResources) {
	r.BackBtn = "Назад"
	r.RemovePlayerBtn = "Удалить игрока"
	return
}

type ProfileResources struct {
	CancelNotifyBtn string
	LevelBtn        string
	NotifiesBtn     string
	NotifyBtn       string
	ParseMode       string
	SexBtn          string
	Text            string
}

func NewProfileResourcesRu() (r ProfileResources) {
	r.CancelNotifyBtn = "При отмене"
	r.LevelBtn = "Уровень"
	r.NotifiesBtn = "Оповещения"
	r.NotifyBtn = "При изменениях"
	r.ParseMode = "Markdown"
	r.SexBtn = "Пол"
	r.Text = ""
	return
}

type SettingsResources struct {
	ActivityBtn string
	BackBtn     string
	CourtBtn    string
	LevelBtn    string
	MaxBtn      string
	NetTypeBtn  string
	PriceBtn    string
}

func NewSettingsResourcesRu() (r SettingsResources) {
	r.ActivityBtn = "Вид активности"
	r.BackBtn = "Назад"
	r.CourtBtn = "🏐 Площадки"
	r.LevelBtn = "💪 Уровень"
	r.MaxBtn = "👫 Мест"
	r.NetTypeBtn = "📏 Вид сетки"
	r.PriceBtn = "💰 Стоимость"
	return
}

type MaxPlayersResources struct {
	BackBtn          string `json:"back_btn"`
	Columns          int    `json:"columns"`
	GroupChatWarning string `json:"group_chat_warning"`
	Message          string `json:"message"`
}

func NewMaxPlayersResourcesRu() (r MaxPlayersResources) {
	r.BackBtn = "Назад"
	r.Columns = 4
	r.GroupChatWarning = "⚠️*Внимание* - здесь функция добавления игроков ограничена числом игроков записи. " +
		"В чате с ботом можно добавить больше игроков в резерв!"
	return
}

type CourtsResources struct {
	Columns int    `json:"columns"`
	Message string `json:"message"`
}

func NewCourtsResourcesRu() CourtsResources {
	return CourtsResources{Columns: 4, Message: "❓Сколько нужно кортов❓"}
}

type PriceResources struct {
	Columns int
	Message string
}

func NewPriceResourcesRu() PriceResources {
	return PriceResources{Columns: 4, Message: "❓Почем будет поиграть❓"}
}

type AcivityResources struct {
	Columns int    `json:"columns"`
	Message string `json:"message"`
}

func NewAcivityResourcesRu() AcivityResources {
	return AcivityResources{Columns: 1, Message: "❓Какой будет вид активности❓"}
}

type LevelResources struct {
	Columns int    `json:"columns"`
	Message string `json:"message"`
}

func NewLevelResourcesRu() LevelResources {
	return LevelResources{Columns: 3, Message: "❓Какой минимальный уровень игроков❓"}
}

type ConfigResources struct {
	Courts    ConfigCourtsResources `json:"courts"`
	Price     ConfigPriceResources  `json:"price"`
	ParseMode string
}

func NewConfigResourcesRu() (cfg ConfigResources) {
	cfg.ParseMode = "markdown"
	cfg.Courts = NewConfigCourtsResourcesRu()
	cfg.Price = NewConfigPriceResourcesRu()
	return
}

type ConfigCourtsResources struct {
	CourtBtn      string `json:"courts_btn"`
	Max           string `json:"max"`
	MaxBtn        string `json:"max_btn"`
	MaxPlayers    string `json:"max_players"`
	MaxPlayersBtn string `json:"max_players_btn"`
	MinPlayers    string `json:"min_players"`
	MinPlayersBtn string `json:"min_players_btn"`
}

func NewConfigCourtsResourcesRu() ConfigCourtsResources {
	return ConfigCourtsResources{
		CourtBtn:      "Настройки площадок",
		Max:           "Площадок",
		MaxBtn:        "Площадки",
		MinPlayers:    "Игроков (min)",
		MinPlayersBtn: "Игроков (min)",
		MaxPlayers:    "Игроков (max)",
		MaxPlayersBtn: "Игроков (max)",
	}
}

type ConfigPriceResources struct {
	PriceBtn string `json:"price_btn"`
	Min      string `json:"min"`
	MinBtn   string `json:"min_btn"`
	Max      string `json:"max"`
	MaxBtn   string `json:"max_btn"`
	Step     string `json:"step"`
	StepBtn  string `json:"step_btn"`
}

func NewConfigPriceResourcesRu() ConfigPriceResources {
	return ConfigPriceResources{
		PriceBtn: "Настройки цены",
		Min:      "Цена (min)",
		MinBtn:   "Цена (min)",
		Max:      "Цена (max)",
		MaxBtn:   "Цена (max)",
		Step:     "Шаг",
		StepBtn:  "Шаг",
	}
}
