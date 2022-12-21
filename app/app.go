package app

import (
	"log"
	"net/http"

	logger "project/log"
	"project/pkg/postgresql"
	"project/router"

	"github.com/gorilla/mux"
)

func StartApp() {
	r := mux.NewRouter()

	postgresql.DatabaseInit()
	router.RouterInit(r.PathPrefix("/api/v1").Subrouter())

	logger.Info("Start App")

	log.Fatal(http.ListenAndServe(":8080", r))
}