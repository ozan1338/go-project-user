package producthandler

import (
	"net/http"
	"project/helpers"
)

type productHandler struct {
	helper helpers.HelpersInterface
}

func NewProductHandler(helper helpers.HelpersInterface) *productHandler {
	return &productHandler{
		helper: helper,
	}
}

func (h productHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	h.helper.WriteResponse(w,http.StatusNotImplemented, "not yet ready")
}