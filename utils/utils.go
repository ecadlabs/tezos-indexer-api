package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ecadlabs/tezos-indexer-api/errors"
	errorsv2 "github.com/pkg/errors"
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
	Code  string `json:"code,omitempty"`
	Cause string `json:"cause,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func JSONError(w http.ResponseWriter, err error) {
	var code errors.Code

	if e, ok := err.(errors.Error); ok {
		code = e.Code()
	} else {
		code = errors.CodeUnknown
	}

	res := errorResponse{
		Error: err.Error(),
		Code:  code.String(),
	}

	if cause := errorsv2.Cause(err); cause != err {
		res.Cause = cause.Error()
	}

	JSONResponse(w, code.Status(), &res)
}

type Paginated struct {
	Value      interface{} `json:"value"`
	TotalCount *int        `json:"total_count,omitempty"`
	Next       string      `json:"next,omitempty"`
}
