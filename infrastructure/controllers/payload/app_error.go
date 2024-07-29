package payload

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"-"`
	Key        string `json:"-"`
}

func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    fmt.Sprintf("Log: %s\nMessage:%s", log, msg),
		Log:        log,
		Key:        key,
	}
}

func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    fmt.Sprintf("Log: %s\nMessage:%s", log, msg),
		Log:        log,
		Key:        key,
	}
}

func NewUnauthorized(root error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
	}
}

func NewPermissionDenied(root error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		RootErr:    root,
		Message:    msg,
		Key:        key,
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

func NewCustomError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, root.Error(), key)
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in the server", err.Error(), "ErrInternal")
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot list %s", strings.ToLower(entity)), fmt.Sprintf("ErrCannotList%s", entity))
}

func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)), fmt.Sprintf("ErrCannotDelete%s", entity))
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot update %s", strings.ToLower(entity)), fmt.Sprintf("ErrCannotUpdate%s", entity))
}

func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot get %s", strings.ToLower(entity)), fmt.Sprintf("ErrCannotGet%s", entity))
}

func ErrEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s deleted", strings.ToLower(entity)), fmt.Sprintf("Err%sDeleted", entity))
}

func ErrEntityExisted(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s already exists", strings.ToLower(entity)), fmt.Sprintf("Err%sAlreadyExists", entity))
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewFullErrorResponse(http.StatusNotFound, err, fmt.Sprintf("%s not found", strings.ToLower(entity)), err.Error(), fmt.Sprintf("Err%sNotFound"))
}

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot create %s", strings.ToLower(entity)), fmt.Sprintf("ErrCannotCreate%s", entity))
}

func ErrNoPermission(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("You have no permission"), fmt.Sprintf("ErrNoPermission"))
}

func ErrBindingData(err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("error binding data"), fmt.Sprintf("ErrBindingData"))
}

func ErrValidateFailed(err error) *AppError {
	return NewCustomError(err, err.Error(), fmt.Sprintf("ErrValidateFailed"))
}

func ErrParamRequired(err error) *AppError {
	return NewCustomError(err, err.Error(), fmt.Sprintf("ErrParamRequired"))
}

func ErrUploadFile(err error) *AppError {
	return NewCustomError(err, err.Error(), fmt.Sprintf("ErrUploadFile"))
}

func ErrDetectFileType(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrDetectFileType")
}

func ErrResetFilePointer(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrResetFilePointer")
}

func ErrHashPassword(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrHashPassword")
}

func ErrInvalidHashPassword(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrInvalidHashPassword")
}

func ErrExisted(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrExisted")
}

func ErrWrongPassword(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrIncorrectAccount")
}

func ErrGenerateToken(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrGenerateToken")
}

func ErrInvalidToken(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrInvalidToken")
}

func ErrPermissionDenied(err error) *AppError {
	return NewCustomError(err, err.Error(), "ErrPermissionDenied")
}