package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateCurrency = errors.New("duplicate currency")
)

type Models struct {
	Currencies CurrencyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Currencies: CurrencyModel{DB: db},
	}
}
