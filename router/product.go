package router

import (
	"net/http"
	"project/handlers/producthandler"
	"project/helpers"

	"github.com/gorilla/mux"
)

func productRouter(r *mux.Router) {

	helper := helpers.NewHelper()

	p := producthandler.NewProductHandler(helper)

	subRouter := r.PathPrefix("/product").Subrouter()
	subRouter.Use(authRequired)
	subRouter.HandleFunc("/upload-file", uploadFiles(p.UploadFile)).Methods(http.MethodPost)
}