package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/lighten/internal/data"
	"github.com/lighten/internal/jsonlog"
)

const version = "1.0.0"

// Holds configuration values
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps float64
		burst int
		enabled bool
	}
}

// Holds the application logic and dependencies
type application struct {
	logger *jsonlog.Logger
	config config
	models data.Models
}

// openDB opens a connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// custom config of the DB connection
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

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Set port value")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development/staging/environment)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("LIGHTEN_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle connections time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
	logger.PrintInfo("database connection pool established", nil)
	defer db.Close()

	app := &application{
		logger: logger,
		config: cfg,
		models: data.NewModels(db),
	}

	err = app.serve()

	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
