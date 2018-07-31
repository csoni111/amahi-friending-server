package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/frnd/").Subrouter()
	apiRouter.HandleFunc("requests", getFRs).Methods("GET")
	apiRouter.HandleFunc("users", getFUs).Methods("GET")
	apiRouter.HandleFunc("new", newFR).Methods("POST")
	apiRouter.HandleFunc("request/{id:[0-9]+}", removeFR).Methods("DELETE")
	apiRouter.HandleFunc("user/{id:[0-9]+}", removeFU).Methods("DELETE")
	apiRouter.HandleFunc("request/{id:[0-9]+}/resend", resendFR).Methods("PUT")

	server := new(http.Server)
	server.Handler = router
	err := server.ListenAndServe()
	handle(err)
}
