package res

import (
	"volleybot/pkg/telegram"
)

type PersonResources struct {
	ProfileCommand telegram.BotCommand
	Level          PlayerLevelResources
}

type PersonResourceLoader interface {
	GetResource() PersonResources
}

type PlayerLevelResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}
