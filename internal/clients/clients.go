package clients

type Clients struct {
	Freecurrency FreecurrencyClient
}

func NewClients(freecurrencyApikey string) Clients {
	return Clients{
		Freecurrency: FreecurrencyClient{
			url:    "https://api.freecurrencyapi.com",
			apiKey: freecurrencyApikey,
		},
	}
}
