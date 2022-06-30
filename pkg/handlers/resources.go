package handlers

import (
	"fmt"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/telegram"

	"github.com/goodsign/monday"
)

type DateTimeResources struct {
	DayCount    int
	DateMessage string
	DateButton  string
	TimeMessage string
	TimeButton  string
}

type CourtResources struct {
	Message    string
	Button     string
	Min        int
	Max        int
	MaxPlayers int
}

type PlayerLevelResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type ActivityResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type SetResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type MaxPlayerResources struct {
	Message          string
	CountError       string
	GroupChatWarning string
	Button           string
	Min              int
	Max              int
}

type DescriptionResources struct {
	Message     string
	DoneMessage string
	Button      string
}

type JoinPlayerResources struct {
	Message          string
	Button           string
	LeaveButton      string
	MultiButton      string
	MultiButtonEmoji string
	MultiButtonText  string
}

type PriceResources struct {
	Message string
	Button  string
	Min     int
	Max     int
	Step    int
}

type CancelResources struct {
	Message string
	Button  string
	Confirm string
	Abort   string
}

type OrderResources struct {
	Location          location.Location
	ActionsBtn        string
	SettingsBtn       string
	BackBtn           string
	CopyBtn           string
	CopyMessage       string
	PublishBtn        string
	RefreshBtn        string
	ListCommand       telegram.BotCommand
	OrderCommand      telegram.BotCommand
	Locale            monday.Locale
	Description       DescriptionResources
	DateTime          DateTimeResources
	Court             CourtResources
	Activity          ActivityResources
	Level             PlayerLevelResources
	Set               SetResources
	MaxPlayer         MaxPlayerResources
	JoinPlayer        JoinPlayerResources
	Price             PriceResources
	Cancel            CancelResources
	RenewMessage      string
	ReservesMessage   string
	NoReservesMessage string
	NoReservesAnswer  string
	OkAnswer          string
}

type OrderResourceLoader interface {
	GetResource() OrderResources
}

type StaticOrderResourceLoader struct{}

func (rl StaticOrderResourceLoader) GetResource() (or OrderResources) {
	or.ListCommand.Command = "list"
	or.ListCommand.Description = "поиск подходящей площадки"
	or.OrderCommand.Command = "order"
	or.OrderCommand.Description = "заказать площадку(и)"
	or.ActionsBtn = "Действия"
	or.SettingsBtn = "Настройки"
	or.BackBtn = "Назад"
	or.CopyBtn = "🫂 Копировать"
	or.CopyMessage = "‼️ *КОПИЯ СДЕЛАНА* ‼️"
	or.PublishBtn = "Опубликовать"
	or.RefreshBtn = "Обновить"
	or.Description.Button = "Описание"
	or.Description.Message = "Отлично. Отправьте мне в чат описание активности."
	or.Description.DoneMessage = "Успешно! Описание обновлено."
	or.Locale = monday.LocaleRuRU
	or.DateTime.DateMessage = "❓Какая дата❓"
	or.DateTime.DateButton = "📆 Дата"
	or.DateTime.DayCount = 30
	or.DateTime.TimeMessage = "❓В какое время❓"
	or.DateTime.TimeButton = "⏰ Время"
	or.JoinPlayer.Message = "❓Сколько игроков записать❓"
	or.JoinPlayer.Button = "😀 Буду"
	or.JoinPlayer.MultiButtonEmoji = "🤩"
	or.JoinPlayer.MultiButtonText = "Буду не один"
	or.JoinPlayer.MultiButton = fmt.Sprintf("%s %s", or.JoinPlayer.MultiButtonEmoji, or.JoinPlayer.MultiButtonText)
	or.JoinPlayer.LeaveButton = "😞 Не смогу"
	or.Activity.Message = "❓Какой будет вид активности❓"
	or.Activity.Button = "Вид активности"
	or.Level.Message = "❓Какой минимальный уровень игроков❓"
	or.Level.Button = "💪 Уровень"
	or.Set.Message = "❓Количество часов❓"
	or.Set.Button = "⏱ Кол-во часов"
	or.Set.Max = 12
	or.Court.Message = "❓Сколько нужно кортов❓"
	or.Court.Button = "🏐 Площадки"
	or.Court.Max = 6
	or.Court.MaxPlayers = 6
	or.MaxPlayer.Message = "❓Максимальное количество игроков❓"
	or.MaxPlayer.CountError = "Ошибка количества игроков!"
	or.MaxPlayer.GroupChatWarning = fmt.Sprintf("⚠️*Внимание* - здесь функция *\"%s\"* ограничена числом игроков записи. "+
		"В чате с ботом можно добавить больше игроков в резерв!", or.JoinPlayer.MultiButtonText)
	or.MaxPlayer.Button = "👫 Мест"
	or.MaxPlayer.Min = 1
	or.MaxPlayer.Max = or.Court.Max * or.Court.MaxPlayers
	or.Price.Message = "❓Почем будет поиграть❓"
	or.Price.Button = "💰 Стоимость"
	or.Price.Min = 0
	or.Price.Max = 2000
	or.Price.Step = 100
	or.Cancel.Button = "💥Отменить"
	or.Cancel.Message = fmt.Sprintf("\n🧨*ВНИМАНИЕ!!!*🧨\nИгра будет отменена для всех участников. Если есть желание только выписаться, лучше воспользоваться кнопкой \"%s\"",
		or.JoinPlayer.LeaveButton)
	or.Cancel.Confirm = "🧨 Уверен"
	or.Cancel.Abort = "👌 Передумал"
	or.RenewMessage = "Запись обновлена и перемещена в конец чата"
	or.ReservesMessage = "❓Какую запись показать ❓"
	or.NoReservesMessage = "На дату %s нет доступных записей."
	or.NoReservesAnswer = "Резервы отсутствуют"
	or.OkAnswer = "Ok"

	return
}

type PersonResources struct {
	CommonHandler
	ProfileCommand telegram.BotCommand
	Level          PlayerLevelResources
}

type PersonResourceLoader interface {
	GetResource() PersonResources
}

type StaticPersonResourceLoader struct{}

func (rl StaticPersonResourceLoader) GetResource() (r PersonResources) {
	r.ProfileCommand.Command = "profile"
	r.ProfileCommand.Description = "настройки профиля пользователя"
	return
}
