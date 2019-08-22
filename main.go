package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ecadlabs/tezos-indexer-api/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		config     service.Config
		configFile string
	)

	flag.StringVar(&configFile, "c", "", "Config file.")
	flag.BoolVar(&config.LogHTTP, "log-http", false, "Log HTTP requests.")
	flag.StringVar(&config.HTTPAddress, "address", ":8000", "HTTP address to listen on.")
	flag.DurationVar(&config.Timeout, "timeout", 0, "PostgreSQL request timeout.")
	flag.IntVar(&config.MaxConnections, "max-connections", 4, "Maximum number of PostgreSQL connections.")
	flag.StringVar(&config.URI, "db", "postgres://indexer:indexer@localhost:5432/mainnet", "PostgreSQL server URI.")

	flag.Parse()

	if configFile != "" {
		if err := config.Load(configFile); err != nil {
			log.Fatal(err)
		}
		// Override from command line
		flag.Parse()
	}

	svc, err := service.NewService(&config, log.StandardLogger())
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Addr:    config.HTTPAddress,
		Handler: svc.NewAPIHandler(),
	}

	log.Printf("HTTP server listening on %s", config.HTTPAddress)

	errChan := make(chan error)
	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}

		case s := <-signalChan:
			log.Printf("Captured %v. Exiting...\n", s)
			return
		}
	}
}
