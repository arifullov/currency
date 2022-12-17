package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Currency struct {
	CurrencyFrom string    `json:"currency_from"`
	CurrencyTo   string    `json:"currency_to"`
	Well         float64   `json:"well"`
	UpdatedAt    time.Time `json:"-"`
}

type CurrencyModel struct {
	DB *sql.DB
}

func (m CurrencyModel) Get(currencyFrom, currencyTo string) (*Currency, error) {
	query := "SELECT currency_from, currency_to, well FROM currencies WHERE currency_from = $1 AND currency_to = $2"

	var currency Currency

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, currencyFrom, currencyTo).Scan(
		&currency.CurrencyFrom,
		&currency.CurrencyTo,
		&currency.Well,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &currency, nil
}

func (m CurrencyModel) Insert(currency *Currency) error {
	query := "INSERT INTO currencies (currency_from, currency_to, well) VALUES ($1, $2, $3)"

	args := []any{currency.CurrencyFrom, currency.CurrencyTo, currency.Well}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "currencies_pkey"`:
			return ErrDuplicateCurrency
		default:
			return err
		}
	}

	return err
}

func (m CurrencyModel) Update(currency *Currency) error {
	query := `
	UPDATE currencies
	SET well = $1, updated_at = $2
	WHERE currency_from = $3 AND currency_to = $4
	RETURNING well 
	`

	args := []any{
		currency.Well,
		time.Now(),
		currency.CurrencyFrom,
		currency.CurrencyTo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&currency.Well)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m CurrencyModel) GetAll() ([]*Currency, error) {
	query := `SELECT currency_from, currency_to, well FROM currencies`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	currencies := []*Currency{}

	for rows.Next() {
		var currency Currency

		err := rows.Scan(
			&currency.CurrencyFrom,
			&currency.CurrencyTo,
			&currency.Well,
		)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, &currency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
