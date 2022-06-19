package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"volleybot/pkg/handlers"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StartHandler struct {
	Bot           *telegram.Bot
	Command       telegram.BotCommand
	orderHandler  *handlers.OrderBotHandler
	personHandler *handlers.PersonBotHandler
}

func (h *StartHandler) StartCmd(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	if msg.Chat.Id <= 0 {
		return
	}
	text := "*Привет!*\n" +
		"Я достаточно молодой волейбольный бот, но кое-что я могу.\n\n" +
		"*Вот те команды, которые я уже понимаю:*"
	cmds := []telegram.BotCommand{h.Command}
	cmds = append(cmds, h.orderHandler.GetCommands(msg.From)...)
	cmds = append(cmds, h.personHandler.GetCommands(msg.From)...)
	for _, cmd := range cmds {
		text += fmt.Sprintf("\n/%s - %s", cmd.Command, cmd.Description)
	}
	h.Bot.SendRequest(&telegram.SetMyCommandsRequest{
		Commands: cmds, Scope: telegram.BotCommandScopeChat{Type: "chat", ChatId: msg.From.Id}})

	mr := &telegram.MessageRequest{
		ChatId:    msg.Chat.Id,
		Text:      text,
		ParseMode: "Markdown"}
	return h.Bot.SendMessage(mr), nil
}

func main() {

	url := os.Getenv("PGURL")
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	pservice, _ := services.NewPersonService(services.WithPgPersonRepository(dbpool))
	oservice, _ := services.NewOrderService(pservice,
		services.WithPgLocationRepository(dbpool),
		services.WithPgReserveRepository(dbpool))
	tb, _ := telegram.NewBot(&telegram.Bot{Token: os.Getenv("TOKEN")})
	tb.Client = &http.Client{}

	orderHandler := handlers.NewOrderHandler(tb, oservice, handlers.StaticOrderResourceLoader{})
	orderHandler.StateRepository = telegram.NewMemoryStateRepository()
	personHandler := handlers.NewPersonHandler(tb, pservice, handlers.StaticPersonResourceLoader{})
	if os.Getenv("LOCATION") != "" {
		orderHandler.Resources.Location.Name = os.Getenv("LOCATION")
	} else {
		orderHandler.Resources.Location.Name = "default"
	}

	lp, _ := tb.NewPoller()
	lp.UpdateHandlers[0].AppendMessageHandlers(orderHandler.GetMessageHandler()...)
	lp.UpdateHandlers[0].AppendCallbackHandlers(orderHandler.GetCallbackHandlers()...)
	lp.UpdateHandlers[0].AppendMessageHandlers(personHandler.GetMessageHandler()...)
	lp.UpdateHandlers[0].AppendCallbackHandlers(personHandler.GetCallbackHandlers()...)

	sh := StartHandler{Bot: tb}
	startcmd := telegram.CommandHandler{
		Command: "start", Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return sh.StartCmd(m, nil)
		}}
	lp.UpdateHandlers[0].AppendMessageHandlers(&startcmd)

	sh.orderHandler = &orderHandler
	sh.personHandler = &personHandler
	sh.Command.Command = "start"
	sh.Command.Description = "начать работу с ботом"
	cmds := []telegram.BotCommand{sh.Command}

	_, err = tb.SendRequest(&telegram.SetMyCommandsRequest{
		Commands: cmds, Scope: telegram.BotCommandScope{Type: "all_private_chats"}})
	if err == nil {
		lp.Run()
	}
}
