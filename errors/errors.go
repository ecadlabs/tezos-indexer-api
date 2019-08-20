package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Error interface {
	error
	Code() Code
}

type Code interface {
	String() string
	Status() int
}

type withCode struct {
	error
	code Code
}

func (w *withCode) Code() Code { return w.code }
func (w *withCode) Format(state fmt.State, verb rune) {
	w.error.(fmt.Formatter).Format(state, verb) // implemented by all "github.com/pkg/errors" types
}

func New(e string, code Code) Error {
	return &withCode{
		errors.New(e), // reuse "github.com/pkg/errors" stack tracing and formatting functionality
		code,
	}
}

type withErr withCode

func Wrap(e error, code Code) Error {
	return &withErr{
		errors.WithStack(e),
		code,
	}
}

func (w *withErr) Code() Code                        { return (*withCode)(w).Code() }
func (w *withErr) Format(state fmt.State, verb rune) { (*withCode)(w).Format(state, verb) }
func (w *withErr) Cause() error                      { return w.error }

type stdCode string

func (s stdCode) String() string { return string(s) }

func (s stdCode) Status() int {
	if status, ok := httpStatus[s]; ok {
		return status
	}
	return http.StatusInternalServerError
}

var (
	CodeUnknown          Code = stdCode("unknown")
	CodeResourceNotFound Code = stdCode("resource_not_found")
	CodeBadRequest       Code = stdCode("bad_request")
	CodeUnauthorized     Code = stdCode("unauthorized")
	CodeForbidden        Code = stdCode("forbidden")
	CodeEndpointNotFound Code = stdCode("endpoint_not_found")
	CodeLimitTooBig      Code = stdCode("limit_too_big")
)

var httpStatus = map[stdCode]int{
	CodeUnknown.(stdCode):          http.StatusInternalServerError,
	CodeResourceNotFound.(stdCode): http.StatusNotFound,
	CodeBadRequest.(stdCode):       http.StatusBadRequest,
	CodeForbidden.(stdCode):        http.StatusForbidden,
	CodeUnauthorized.(stdCode):     http.StatusUnauthorized,
	CodeEndpointNotFound.(stdCode): http.StatusNotFound,
	CodeLimitTooBig.(stdCode):      http.StatusBadRequest,
}

// Some predefined errors
var (
	ErrResourceNotFound = New("Resource not found", CodeResourceNotFound)
	ErrForbidden        = New("Forbidden", CodeForbidden)
	ErrEndpointNotFound = New("Endpoint not found", CodeEndpointNotFound)
)
