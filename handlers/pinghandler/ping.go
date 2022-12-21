package pinghandler

import (
	"net/http"
	"project/helpers"
)

type pingHandler struct{
	Helper helpers.HelpersInterface
}

func PingHandler(Helpers helpers.HelpersInterface) *pingHandler {
	return &pingHandler{Helper: Helpers}
}

func (h pingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	h.Helper.WriteResponse(w,http.StatusOK, "OK")
}