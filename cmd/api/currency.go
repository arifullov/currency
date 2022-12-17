package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/arifullov/currency/internal/data"
	"github.com/gofiber/fiber/v2"
)

func (app *application) createCurrencyHandler(c *fiber.Ctx) error {
	var input struct {
		CurrencyFrom string `json:"currencyFrom"`
		CurrencyTo   string `json:"currencyTo"`
	}

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	currency := data.Currency{
		CurrencyFrom: strings.ToUpper(input.CurrencyFrom),
		CurrencyTo:   strings.ToUpper(input.CurrencyTo),
	}

	err := app.models.Currencies.Insert(&currency)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCurrency):
			c.Status(http.StatusUnprocessableEntity)
			return c.JSON(map[string]string{
				"error": "currency already exist",
			})
		}
		return err
	}

	go func() {
		err = app.services.FreecurrencyService.ActualizeCurrency(currency)
		if err != nil {
			app.logger.Error(err)
		}
	}()

	c.Status(http.StatusCreated)
	return c.JSON(input)
}

func (app *application) getCurrentCurrencyHandler(c *fiber.Ctx) error {
	var input struct {
		CurrencyFrom string  `json:"currencyFrom"`
		CurrencyTo   string  `json:"currencyTo"`
		Value        float64 `json:"value"`
	}

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	currency := data.Currency{
		CurrencyFrom: input.CurrencyFrom,
		CurrencyTo:   input.CurrencyTo,
		Well:         input.Value,
	}

	err := app.models.Currencies.Update(&currency)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			c.Status(http.StatusNotFound)
			return c.JSON(map[string]string{
				"error": "currency not found",
			})
		}
		return err
	}

	return c.JSON(currency)
}

func (app *application) getAllCurrenciesHandler(c *fiber.Ctx) error {
	currencies, err := app.models.Currencies.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(currencies)
}
