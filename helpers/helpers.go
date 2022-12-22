package helpers

import (
	"encoding/json"
	"io"
	"net/http"
	resError "project/util/errors_response"
)

type HelpersInterface interface {
	WriteResponse(w http.ResponseWriter,code int, data any )
	ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) resError.RespError
}

type HelpersStruct struct{}

func (h HelpersStruct) WriteResponse(w http.ResponseWriter,code int, data any ) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func (h HelpersStruct) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) resError.RespError {
	maxBytes := 1024 * 1024 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// attempt to decode the data
	err := dec.Decode(data)
	if err != nil {
		return resError.NewBadRequestError(err.Error())
	}

	// make sure only one JSON value in payload
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return resError.NewBadRequestError("body must only contain a single JSON value")
	}

	return nil
}

func NewHelper() *HelpersStruct {
	return &HelpersStruct{}
}