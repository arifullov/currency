package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FreecurrencyClient struct {
	url    string
	apiKey string
}

type ExchangeRates struct {
	Data map[string]float64 `json:"data"`
}

func (c FreecurrencyClient) GetLatestExchangeRates(baseCurrency string, currencies []string) (*ExchangeRates, error) {
	requestURL := fmt.Sprintf(
		"%s/v1/latest?apikey=%s&base_currency=%s&currencies=%s",
		c.url, c.apiKey, baseCurrency, strings.Join(currencies, ","),
	)
	r, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	exchangeRates := ExchangeRates{}
	return &exchangeRates, json.NewDecoder(r.Body).Decode(&exchangeRates)
}
