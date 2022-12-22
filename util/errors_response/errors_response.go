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

type RespErrorStruct struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Err     string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func (e RespErrorStruct) GetError() string {
	return fmt.Sprintf(e.Err)
}

func (e RespErrorStruct) GetStatus() int {
	return e.Status
}

func (e RespErrorStruct) GetMessage() string {
	return fmt.Sprintf(e.Message)
}

func NewRespError(message string, status int, err string) RespError{
	return RespErrorStruct{
		Message: message,
		Status: status,
		Err: err,
	}
}

func NewBadRequestError(message string) RespError {
	return RespErrorStruct{
		Message: message,
		Status:  http.StatusBadRequest,
		Err: "bad_request",
	}
}
