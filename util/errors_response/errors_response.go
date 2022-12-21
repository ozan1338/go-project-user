package errors_response

import (
	"errors"
	"fmt"
	"net/http"
)

//go:generate mockgen -destination=../../mocks/util/errors_response/mockErrorResponse.go -package=errors_response project/util/errors_reponse RespError
type RespError interface {
	GetMessage() string
	GetStatus() int
	GetError() string
}

type respError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Err     string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func (e respError) GetError() string {
	return fmt.Sprintf(e.Err)
}

func (e respError) GetStatus() int {
	return e.Status
}

func (e respError) GetMessage() string {
	return fmt.Sprintf(e.Message)
}

func NewRespError(message string, status int, err string) RespError{
	return respError{
		Message: message,
		Status: status,
		Err: err,
	}
}

func NewBadRequestError(message string) RespError {
	return respError{
		Message: message,
		Status:  http.StatusBadRequest,
		Err: "bad_request",
	}
}
