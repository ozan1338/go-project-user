package app

import (
	"log"
	"net/http"

	logger "project/log"
	"project/pkg/postgresql"
	"project/router"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartApp() {
	r := mux.NewRouter()

	postgresql.DatabaseInit()
	router.RouterInit(r.PathPrefix("/api/v1").Subrouter())

	logger.Info("Start App")

	// Setup allowed Header, Method, and Origin for CORS on this below code ...
	var AllowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var AllowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"})
	var AllowedOrigins = handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(AllowedHeaders,AllowedMethods,AllowedOrigins)(r)))
}