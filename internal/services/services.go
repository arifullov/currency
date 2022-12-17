package services

import (
	"github.com/arifullov/currency/internal/clients"
	"github.com/arifullov/currency/internal/data"
	"github.com/sirupsen/logrus"
)

type Services struct {
	FreecurrencyService FreecurrencyService
}

func NewServices(clients clients.Clients, models data.Models, logger *logrus.Logger) Services {
	return Services{
		FreecurrencyService: FreecurrencyService{
			clients: clients,
			models:  models,
			logger:  logger,
		},
	}
}
