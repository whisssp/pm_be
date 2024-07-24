package payload

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	ErrKey     string      `json:"-"`
	RootErr    error       `json:"-"`
	Data       interface{} `json:"data"`
}

func NewErrResponse(code int, message, key string, RootErr error) *AppError {
	return &AppError{
		StatusCode: code,
		Message:    fmt.Sprintf("%s: %s.", key, RootErr.Error()),
		ErrKey:     key,
		RootErr:    RootErr,
		Data:       nil,
	}
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func NewCustomError(code int, root error, msg, key string) *AppError {
	return NewErrResponse(code, root.Error(), key, errors.New(msg))
}

func ErrInvalidRequest(err error) *AppError {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), "ErrInvalidRequest")
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(http.StatusNotFound, err, err.Error(), fmt.Sprintf("ErrNotFound"))
}

func ErrBindingData(err error) *AppError {
	return NewCustomError(http.StatusInternalServerError, err, err.Error(), fmt.Sprintf("ErrBindingData"))
}

func ErrValidateFailed(err error) *AppError {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), fmt.Sprintf("ErrValidateFailed"))
}

func ErrDB(err error) *AppError {
	return NewErrResponse(http.StatusInternalServerError, err.Error(), "ERR_DB", err)
}

func ErrParamRequired(err error) *AppError {
	return NewCustomError(http.StatusBadRequest, err, fmt.Sprintf("missing path parameter. %s", err.Error()), fmt.Sprintf("ErrParamRequired"))
}