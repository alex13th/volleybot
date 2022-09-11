package res

import (
	"volleybot/pkg/bvbot"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/telegram"
)

type StaticVolleyResourceLoader struct{}

func (rl StaticVolleyResourceLoader) GetResource() (or VolleyResources) {
	or.ReserveView = reserve.NewTelegramResourcesRu()
	or.Volley.Actions = bvbot.NewActionsResourcesRu()
	or.Volley.Activity = bvbot.NewAcivityResourcesRu()
	or.Volley.Cancel = bvbot.NewCancelResourcesRu()
	or.Volley.Courts = bvbot.NewCourtsResourcesRu()
	or.Volley.Description = bvbot.NewDescResourcesRu()
	or.Volley.Join = bvbot.NewJoinPlayersResourcesRu()
	or.Volley.Level = bvbot.NewLevelResourcesRu()
	or.Volley.List = bvbot.NewListResourcesRu()
	or.Volley.Main = bvbot.NewMainResourcesRu()
	or.Volley.MaxPlayer = bvbot.NewMaxPlayersResourcesRu()
	or.Volley.Price = bvbot.NewPriceResourcesRu()
	or.Volley.Profile = bvbot.ProfileResourcesRu()
	or.Volley.RemovePlayer = bvbot.RemovePlayerResourcesRu()
	or.Volley.Settings = bvbot.NewSettingsResourcesRu()
	or.Volley.Show = bvbot.NewShowResourcesRu()
	or.Volley.Sets = bvbot.NewSetsResourcesRu()
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
	Volley      bvbot.StateResources
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
