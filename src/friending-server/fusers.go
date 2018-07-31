package main

import "net/http"

type FriendUser struct {
	BaseModel
	ApiKey string
	Email  string
}

func getFUs(w http.ResponseWriter, r *http.Request) {

}

func removeFU(w http.ResponseWriter, r *http.Request) {

}
