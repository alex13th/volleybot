package res

import (
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"
)

type StaticVolleyResourceLoader struct{}

func (rl StaticVolleyResourceLoader) GetResource() (or VolleyResources) {
	or.ReserveView = reserve.NewTelegramResourcesRu()
	or.Volley.Actions = volley.NewActionsResourcesRu()
	or.Volley.Activity = volley.NewAcivityResourcesRu()
	or.Volley.Cancel = volley.NewCancelResourcesRu()
	or.Volley.Courts = volley.NewCourtsResourcesRu()
	or.Volley.Description = volley.NewDescResourcesRu()
	or.Volley.Join = volley.NewJoinPlayersResourcesRu()
	or.Volley.Level = volley.NewLevelResourcesRu()
	or.Volley.List = volley.NewListResourcesRu()
	or.Volley.Main = volley.NewMainResourcesRu()
	or.Volley.MaxPlayer = volley.NewMaxPlayersResourcesRu()
	or.Volley.Price = volley.NewPriceResourcesRu()
	or.Volley.RemovePlayer = volley.RemovePlayerResourcesRu()
	or.Volley.Settings = volley.NewSettingsResourcesRu()
	or.Volley.Show = volley.NewShowResourcesRu()
	or.Volley.Sets = volley.NewSetsResourcesRu()
	or.Volley.BackBtn = "Назад"
	or.Volley.DescMessage = "Отлично. Отправьте мне в чат описание активности."
	or.Command.Command = "volley"
	or.Command.Description = "Пляжный волейбол"
	return
}

type VolleyResources struct {
	Command     telegram.BotCommand
	Location    location.Location
	ReserveView reserve.TelegramViewResources
	Volley      volley.MessageProcessorResources
}

type OrderResourceLoader interface {
	GetResource() VolleyResources
}

type StaticPersonResourceLoader struct{}

func (rl StaticPersonResourceLoader) GetResource() (r PersonResources) {
	r.ProfileCommand.Command = "profile"
	r.ProfileCommand.Description = "настройки профиля пользователя"
	return
}
