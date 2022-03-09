package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graillus/timelapse/internal/api"
	"github.com/graillus/timelapse/internal/log"
)

type config struct {
	listenAddr  string
	storagePath string
}

func main() {
	cfg := loadConfig()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(resp http.ResponseWriter, _ *http.Request) {
		resp.WriteHeader(http.StatusOK)
	})

	apiRouter := router.PathPrefix("/api").Subrouter()
	api.ConfigureRoutes(apiRouter)
	api.StoragePath = cfg.storagePath

	go func() {
		log.Infof("Server listening to %s", cfg.listenAddr)
		err := http.ListenAndServe(cfg.listenAddr, router)
		log.Fatalf("Fatal server error: %v", err)
	}()

	select {}
}

func loadConfig() config {
	return config{
		listenAddr:  "0.0.0.0:8990",
		storagePath: "/tmp",
	}
}
