package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) SetupRoutes(fiberApp *fiber.App) {
	api := fiberApp.Group("/")

	api.Get("/healthcheck", healthcheckHandler)

	//currency
	currency := api.Group("/api")
	currency.Get("/currency", app.getAllCurrenciesHandler)
	currency.Post("/currency", app.createCurrencyHandler)
	currency.Put("/currency", app.getCurrentCurrencyHandler)
}
