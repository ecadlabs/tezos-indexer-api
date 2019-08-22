package service

import (
	"context"
	"net/http"

	"github.com/ecadlabs/tezos-indexer-api/errors"
	"github.com/ecadlabs/tezos-indexer-api/middleware"
	"github.com/ecadlabs/tezos-indexer-api/storage/pg"
	"github.com/ecadlabs/tezos-indexer-api/utils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	pool    *pgxpool.Pool
	storage *pg.PostgresStorage
	config  *Config
	logger  log.FieldLogger
}

func (s *Service) log() log.FieldLogger {
	if s.logger != nil {
		return s.logger
	}
	return log.StandardLogger()
}

func NewService(c *Config, logger log.FieldLogger) (*Service, error) {
	poolConfig, err := pgxpool.ParseConfig(c.URI)
	if err != nil {
		return nil, err
	}

	if poolConfig.MaxConns == 0 {
		poolConfig.MaxConns = int32(c.MaxConnections)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	storage := &pg.PostgresStorage{DB: pool}

	return &Service{
		pool:    pool,
		storage: storage,
		config:  c,
		logger:  logger,
	}, nil
}

func (s *Service) NewAPIHandler() http.Handler {
	h := &Handler{
		Storage: s.storage,
		Logger:  s.logger,
		Timeout: s.config.Timeout,
	}

	m := mux.NewRouter()
	if s.config.LogHTTP {
		m.Use((&middleware.Logging{}).Handler)
	}
	m.Use((&middleware.Recover{}).Handler)

	m.Methods("GET").Path("/balances/{pkh}").HandlerFunc(h.GetBalanceUpdate)

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.JSONError(w, errors.ErrResourceNotFound)
	})

	return m
}
