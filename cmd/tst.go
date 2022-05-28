package main

import (
	"net/http"
	"os"
	"volleybot/pkg/handlers"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"
)

type StartHandler struct {
	Bot *telegram.Bot
}

func (h *StartHandler) StartCmd(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	text := "*Привет!*\n" +
		"Я достаточно молодой волейбольный бот, но кое что я могу.\n\n" +
		"*Вот те команды, которые я уже понимаю:*\n" +
		"/list - можно просмотреть список уже заказанных площвдок;\n" +
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
	oservice, _ := services.NewOrderService(
		services.WithPgPersonRepository(url),
		services.WithPgLocationRepository(url),
		services.WithPgReserveRepository(url))
	tb, _ := telegram.NewBot(&telegram.Bot{Token: os.Getenv("TOKEN")})
	tb.Client = &http.Client{}

	orderHandler := handlers.NewOrderHandler(tb, oservice, handlers.DefaultResourceLoader{})
	lp, _ := tb.NewPoller()
	lp.UpdateHandlers[0].AppendMessageHandler(&orderHandler)
	lp.UpdateHandlers[0].AppendCallbackHandler(&orderHandler)

	sh := StartHandler{Bot: tb}
	startcmd := telegram.CommandHandler{
		Command: "start", Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return sh.StartCmd(m, nil)
		}}
	lp.UpdateHandlers[0].AppendMessageHandler(&startcmd)

	lp.Run()
}
