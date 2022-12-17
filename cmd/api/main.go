package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/arifullov/currency/internal/clients"
	"github.com/arifullov/currency/internal/data"
	"github.com/arifullov/currency/internal/services"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	freecurrencyApikey string
}

type application struct {
	config   config
	services services.Services
	logger   *logrus.Logger
	models   data.Models
	wg       sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "API server port")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CURRENCY_DB_DSN"), "PosgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.freecurrencyApikey, "freecurrency-apikey", os.Getenv("FREECURRENCY_APIKEY"), "FreecurrencyAPI api key")

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	app.services = services.NewServices(clients.NewClients(cfg.freecurrencyApikey), app.models, app.logger)

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
