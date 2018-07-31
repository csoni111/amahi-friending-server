package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"gopkg.in/matryer/respond.v1"
	"log"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key; auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	opts := &respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request,
			status int, data interface{}) (int, interface{}) {
			if err, ok := data.(error); ok {
				return status, map[string]interface{}{"error": err.Error()}
			}
			return status, data
		},
		After: func(w http.ResponseWriter, r *http.Request,
			status int, data interface{}) {
			log.Println("->", status, data)
		},
	}

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/frnd/").Subrouter()
	apiRouter.Handle("requests", getFRs).Methods("GET")
	apiRouter.HandleFunc("users", getFUs).Methods("GET")
	apiRouter.HandleFunc("new", newFR).Methods("POST")
	apiRouter.HandleFunc("request/{id:[0-9]+}", removeFR).Methods("DELETE")
	apiRouter.HandleFunc("user/{id:[0-9]+}", removeFU).Methods("DELETE")
	apiRouter.HandleFunc("request/{id:[0-9]+}/resend", resendFR).Methods("PUT")

	server := new(http.Server)
	server.Handler = router
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
