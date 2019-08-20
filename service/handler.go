package service

import (
	"context"
	"net/http"
	"time"

	"github.com/ecadlabs/tezos-indexer-api/errors"
	"github.com/ecadlabs/tezos-indexer-api/storage/pg"
	"github.com/ecadlabs/tezos-indexer-api/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Storage *pg.PostgresStorage
	Logger  log.FieldLogger
	Timeout time.Duration
}

func (h *Handler) log() log.FieldLogger {
	if h.Logger != nil {
		return h.Logger
	}
	return log.StandardLogger()
}

func (h *Handler) context(r *http.Request) (context.Context, context.CancelFunc) {
	if h.Timeout != 0 {
		return context.WithTimeout(r.Context(), h.Timeout)
	}
	return r.Context(), func() {}
}

var schemaDecoder = schema.NewDecoder()

func (h *Handler) GetBalanceUpdate(w http.ResponseWriter, r *http.Request) {
	type getBalanceUpdateRequest struct {
		Start   time.Time `schema:"start"`
		End     time.Time `schema:"end"`
		Limit   int       `schema:"limit"`
		Compact bool      `schema:"compact"`
	}

	r.ParseForm()
	pkh := mux.Vars(r)["pkh"]

	req := getBalanceUpdateRequest{
		Compact: true,
	}

	if err := schemaDecoder.Decode(&req, r.Form); err != nil {
		utils.JSONError(w, errors.Wrap(err, errors.CodeBadRequest))
		return
	}

	ctx, cancel := h.context(r)
	defer cancel()

	ret, err := h.Storage.GetBalanceUpdate(ctx, pkh, req.Start, req.End, req.Limit)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	if req.Compact {
		type compactBalanceUpdate struct {
			BlockLevel     []int64     `json:"level"`
			BlockTimestamp []time.Time `json:"timestamp"`
			Diff           []int64     `json:"diff"`
			Value          []int64     `json:"value"`
		}

		compacted := compactBalanceUpdate{
			BlockLevel:     make([]int64, len(ret)),
			BlockTimestamp: make([]time.Time, len(ret)),
			Diff:           make([]int64, len(ret)),
			Value:          make([]int64, len(ret)),
		}

		for i, u := range ret {
			compacted.BlockLevel[i] = u.BlockLevel
			compacted.BlockTimestamp[i] = u.BlockTimestamp
			compacted.Diff[i] = u.Diff
			compacted.Value[i] = u.Value
		}

		utils.JSONResponse(w, http.StatusOK, &compacted)
		return
	}

	utils.JSONResponse(w, http.StatusOK, ret)
}
