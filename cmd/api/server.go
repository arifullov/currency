package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (app *application) serve() error {
	srv := fiber.New(fiber.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		AppName:      "Currency",
	})

	app.SetupRoutes(srv)

	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		for range ticker.C {
			app.services.FreecurrencyService.ActualizeCurrencies()
		}
	}()

	return srv.Listen(fmt.Sprintf(":%d", app.config.port))
}
