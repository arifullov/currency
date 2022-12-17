package services

import (
	"errors"

	"github.com/arifullov/currency/internal/clients"
	"github.com/arifullov/currency/internal/data"
	"github.com/sirupsen/logrus"
)

type FreecurrencyService struct {
	clients clients.Clients
	models  data.Models
	logger  *logrus.Logger
}

func (s FreecurrencyService) ActualizeCurrency(currency data.Currency) error {
	latestCurrencyData, err := s.clients.Freecurrency.GetLatestExchangeRates(
		currency.CurrencyFrom, []string{currency.CurrencyTo},
	)
	if err != nil {
		return err
	}

	well, ok := latestCurrencyData.Data[currency.CurrencyTo]
	if !ok {
		return errors.New("currency not found")
	}

	currency.Well = well
	err = s.models.Currencies.Update(&currency)
	return err
}

func (s FreecurrencyService) ActualizeCurrencies() {
	s.logger.Info("Actualize currencies start")
	currenciesForUpdate := map[string][]string{}

	currencies, err := s.models.Currencies.GetAll()
	if err != nil {
		s.logger.Error(err)
		return
	}

	for _, currency := range currencies {
		currenciesForUpdate[currency.CurrencyFrom] = append(currenciesForUpdate[currency.CurrencyFrom], currency.CurrencyTo)
	}

	for currency, updatingCurrencies := range currenciesForUpdate {
		latestCurrencyData, err := s.clients.Freecurrency.GetLatestExchangeRates(currency, updatingCurrencies)
		if err != nil {
			s.logger.Error(err)
			return
		}

		for currencyName, value := range latestCurrencyData.Data {
			currencyData := data.Currency{
				CurrencyFrom: currency,
				CurrencyTo:   currencyName,
				Well:         value,
			}
			err = s.models.Currencies.Update(&currencyData)
			if err != nil {
				s.logger.Error(currencyName)
			}
		}
	}

	s.logger.Info("Actualize currencies end")
}
