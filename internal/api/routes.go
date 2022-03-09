package api

import (
	"github.com/gorilla/mux"
)

func ConfigureRoutes(router *mux.Router) {
	router.Path("/frames").Methods("GET").HandlerFunc(listFrames)
	router.Path("/frames").Methods("POST").HandlerFunc(postFrame)
}
