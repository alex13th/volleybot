package res

import (
	"volleybot/pkg/bvbot"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/telegram"
)

type VolleyResourceLoader interface {
	GetResources() (or VolleyResources)
}

type StaticVolleyResourceLoader struct{}

func (rl StaticVolleyResourceLoader) GetResources() (res VolleyResources) {
	res.ReserveView = reserve.NewTelegramResourcesRu()
	res.Resources.Actions = bvbot.NewActionsResourcesRu()
	res.Resources.Activity = bvbot.NewAcivityResourcesRu()
	res.Resources.Cancel = bvbot.NewCancelResourcesRu()
	res.Resources.Courts = bvbot.NewCourtsResourcesRu()
	res.Resources.Description = bvbot.NewDescResourcesRu()
	res.Resources.Join = bvbot.NewJoinPlayersResourcesRu()
	res.Resources.Level = bvbot.NewLevelResourcesRu()
	res.Resources.List = bvbot.NewListResourcesRu()
	res.Resources.Main = bvbot.NewMainResourcesRu()
	res.Resources.MaxPlayer = bvbot.NewMaxPlayersResourcesRu()
	res.Resources.Price = bvbot.NewPriceResourcesRu()
	res.Resources.Profile = bvbot.NewProfileResourcesRu()
	res.Resources.RemovePlayer = bvbot.RemovePlayerResourcesRu()
	res.Resources.Settings = bvbot.NewSettingsResourcesRu()
	res.Resources.Show = bvbot.NewShowResourcesRu()
	res.Resources.Sets = bvbot.NewSetsResourcesRu()
	res.Resources.BackBtn = "Назад"
	res.Resources.DescMessage = "Отлично. Отправьте мне в чат описание активности."
	res.Command.Command = "volley"
	res.Command.Description = "Пляжный волейбол"
	return
}

type VolleyResources struct {
	Command     telegram.BotCommand
	Location    location.Location
	ReserveView reserve.TelegramViewResources
	Resources   bvbot.Resources
}
