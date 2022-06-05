package main

import (
	"context"
	"net/http"
	"os"
	"volleybot/pkg/handlers"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StartHandler struct {
	Bot     *telegram.Bot
	Command telegram.BotCommand
}

func (h *StartHandler) StartCmd(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	if msg.Chat.Id <= 0 {
		return
	}
	text := "*Привет!*\n" +
		"Я достаточно молодой волейбольный бот, но кое-что я могу.\n\n" +
		"*Вот те команды, которые я уже понимаю:*\n" +
		"/list - можно просмотреть список уже заказанных площадок;\n" +
		"/order - забронировать площадки для себя и друзей;\n" +
		"/start - посмотреть это приветствие"
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

	oservice, _ := services.NewOrderService(
		services.WithPgPersonRepository(dbpool),
		services.WithPgLocationRepository(dbpool),
		services.WithPgReserveRepository(dbpool))
	tb, _ := telegram.NewBot(&telegram.Bot{Token: os.Getenv("TOKEN")})
	tb.Client = &http.Client{}

	orderHandler := handlers.NewOrderHandler(tb, oservice, handlers.DefaultResourceLoader{})
	if os.Getenv("LOCATION") != "" {
		orderHandler.Resources.Location.Name = os.Getenv("LOCATION")
	} else {
		orderHandler.Resources.Location.Name = "default"
	}

	lp, _ := tb.NewPoller()
	lp.UpdateHandlers[0].AppendMessageHandler(&orderHandler)
	lp.UpdateHandlers[0].AppendCallbackHandler(&orderHandler)

	sh := StartHandler{Bot: tb}
	startcmd := telegram.CommandHandler{
		Command: "start", Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return sh.StartCmd(m, nil)
		}}
	lp.UpdateHandlers[0].AppendMessageHandler(&startcmd)

	sh.Command.Command = "start"
	sh.Command.Description = "начать работу с ботом"
	cmds := []telegram.BotCommand{sh.Command}
	cmds = append(cmds, orderHandler.GetCommands()...)

	_, err = tb.SendRequest(&telegram.SetMyCommandsRequest{
		Commands: cmds, Scope: telegram.BotCommandScope{Type: "all_private_chats"}})
	if err == nil {
		lp.Run()
	}
}
