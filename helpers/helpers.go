package helpers

import (
	"encoding/json"
	"net/http"
)

type HelpersInterface interface {
	WriteResponse(w http.ResponseWriter,code int, data any )
}

type HelpersStruct struct{}

func (h HelpersStruct) WriteResponse(w http.ResponseWriter,code int, data any ) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func NewHelper() *HelpersStruct {
	return &HelpersStruct{}
}