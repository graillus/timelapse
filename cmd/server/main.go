package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/api", func(resp http.ResponseWriter, req *http.Request) {
		_, err := resp.Write([]byte("Hello World"))
		if err != nil {
			resp.WriteHeader(200)
			return
		}
		resp.WriteHeader(500)
	})

	_ = http.ListenAndServe("0.0.0.0:8990", apiRouter)
	log.Fatal("server error")
}
