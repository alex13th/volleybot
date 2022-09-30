package bvbot

import "time"

type Resources struct {
	Actions      ActionsResources
	Activity     AcivityResources
	Config       ConfigResources
	Courts       CourtsResources
	Cancel       CancelResources
	Description  DescResources
	Join         JoinResources
	Level        LevelResources
	List         ListResources
	Main         MainResources
	MaxPlayer    MaxPlayersResources
	Profile      ProfileResources
	RemovePlayer RemovePlayerResources
	Price        PriceResources
	Settings     SettingsResources
	Sets         SetsResources
	Show         ShowResources
	BackBtn      string
	DescMessage  string
}

type MainResources struct {
	ListCaption       string        `json:"list_caption"`
	ListDateBtn       string        `json:"List_date_btn"`
	NewReserveBtn     string        `json:"new_reserve_msg"`
	NoReservesMessage string        `json:"no_reserve_msg"`
	ParseMode         string        `json:"parse_mode"`
	PreviewDuration   time.Duration `json:"duration"`
	ProfileBtn        string        `json:"profile_btn"`
	SettingsBtn       string        `json:"settings_btn"`
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
	r.SettingsBtn = "🛠 Настройки"
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

type ActionsResources struct {
	BackBtn         string `json:"back_btn"`
	CancelBtn       string `json:"cancel_btn"`
	CopyBtn         string `json:"copy_btn"`
	CopyDoneMessage string `json:"copy_done_msg"`
	PublishBtn      string `json:"publish_btn"`
	RemovePlayerBtn string `json:"remove_player_btn"`
}

func NewActionsResourcesRu() (r ActionsResources) {
	r.BackBtn = "Назад"
	r.CancelBtn = "💥Отменить"
	r.CopyBtn = "🫂 Копировать"
	r.CopyDoneMessage = "Копия сделана! 👆"
	r.PublishBtn = "Опубликовать"
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
		MinPlayers:    "Игороков (min)",
		MinPlayersBtn: "Игороков (min)",
		MaxPlayers:    "Игороков (max)",
		MaxPlayersBtn: "Игороков (max)",
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
