package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env string
}

type application struct{
	logger *log.Logger
	config config
}


func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Set port value")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development/staging/environment)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)


	app := &application{
		logger: logger,
		config: cfg,
	}

	server := &http.Server {
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on port %d", cfg.env, cfg.port)
	logger.Fatal(server.ListenAndServe())
}