package main

import (
	"net/http"
	"os"
	"volleybot/pkg/handlers"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"
)

func main() {

	url := os.Getenv("PGURL")
	oservice, _ := services.NewOrderService(
		services.WithPgPersonRepository(url),
		services.WithPgLocationRepository(url),
		services.WithPgReserveRepository(url))
	tb, _ := telegram.NewBot(&telegram.Bot{Token: os.Getenv("TOKEN")})
	tb.Client = &http.Client{}

	orderHandler := handlers.NewOrderHandler(tb, oservice)
	lp, _ := tb.NewPoller()
	lp.UpdateHandlers[0].AppendMessageHandler(&orderHandler)
	lp.UpdateHandlers[0].AppendCallbackHandler(&orderHandler)

	lp.Run()
}
