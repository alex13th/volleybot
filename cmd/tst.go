package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"volleybot/pkg/postgres"
	"volleybot/pkg/res"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StartHandler struct {
	Bot            *telegram.Bot
	Command        telegram.BotCommand
	ReserveService *services.VolleyBotService
}

func (h *StartHandler) StartCmd(msg *telegram.Message, chanr chan telegram.MessageResponse) error {
	if msg.Chat.Id <= 0 {
		return nil
	}
	text := "*Привет!*\n" +
		"Я достаточно молодой волейбольный бот, но кое-что я могу.\n\n" +
		"*Вот те команды, которые я уже понимаю:*"
	cmds := []telegram.BotCommand{h.Command}
	cmds = append(cmds, h.ReserveService.GetCommands(msg.From.Id)...)
	for _, cmd := range cmds {
		text += fmt.Sprintf("\n/%s - %s", cmd.Command, cmd.Description)
	}
	h.Bot.SendRequest(&telegram.SetMyCommandsRequest{
		Commands: cmds, Scope: telegram.BotCommandScopeChat{Type: "chat", ChatId: msg.From.Id}})

	mr := &telegram.MessageRequest{
		ChatId:    msg.Chat.Id,
		Text:      text,
		ParseMode: "Markdown"}
	h.Bot.SendMessage(mr)
	return nil
}

func main() {

	url := os.Getenv("PGURL")
	tb, _ := telegram.NewBot(&telegram.Bot{Token: os.Getenv("TOKEN")})
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	vres := res.StaticVolleyResourceLoader{}.GetResource()
	lrep, _ := postgres.NewLocationRepository(dbpool)
	lrep.UpdateDB()
	prep, _ := postgres.NewPersonPgRepository(dbpool)
	prep.UpdateDB()
	rrep, _ := postgres.NewVolleyPgRepository(dbpool, &prep, &lrep)
	rrep.UpdateDB()
	strep, _ := postgres.NewStateRepository(dbpool)
	strep.UpdateDB()
	vservice := services.VolleyBotService{Bot: tb, Resources: &vres, StateRepository: &strep,
		LocationRepository: &lrep, VolleyRepository: &rrep, PersonRepository: &prep}

	if os.Getenv("LOCATION") != "" {
		vres.Location.Name = os.Getenv("LOCATION")
	} else {
		vres.Location.Name = "default"
	}

	tb.Client = &http.Client{}
	lp, _ := tb.NewPoller()

	sh := StartHandler{Bot: tb}
	startcmd := telegram.CommandHandler{
		Command: "start", Handler: func(m *telegram.Message) error {
			return sh.StartCmd(m, nil)
		}}
	lp.UpdateHandlers[0].AppendMessageHandlers(&startcmd)
	lp.UpdateHandlers[0].AppendMessageHandlers(&vservice)
	lp.UpdateHandlers[0].AppendCallbackHandlers(&vservice)

	sh.ReserveService = &vservice
	sh.Command.Command = "start"
	sh.Command.Description = "начать работу с ботом"
	cmds := []telegram.BotCommand{sh.Command}

	_, err = tb.SendRequest(&telegram.SetMyCommandsRequest{
		Commands: cmds, Scope: telegram.BotCommandScope{Type: "all_private_chats"}})
	if err == nil {
		lp.Run()
	}
}
