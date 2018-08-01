package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	migrateDb()
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/frnd").Subrouter()
	apiRouter.HandleFunc("/requests", getFRs).Methods("GET")
	apiRouter.HandleFunc("/users", getFUs).Methods("GET")
	apiRouter.HandleFunc("/request", newFR).Methods("POST")
	apiRouter.HandleFunc("/request/{id:[0-9]+}", removeFR).Methods("DELETE")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", removeFU).Methods("DELETE")
	apiRouter.HandleFunc("/request/{id:[0-9]+}/resend", resendFR).Methods("PUT")

	server := new(http.Server)
	server.Handler = router
	err := server.ListenAndServe()
	handle(err)
}

func migrateDb() {
	db, err := getDb()
	defer db.Close()
	handle(err)
	db.LogMode(true)

	// create tables and insert sample data
	if !db.HasTable(&AmahiUser{}) {
		db.AutoMigrate(&AmahiUser{}, &System{}, &FriendUser{}, &FriendRequest{})

		// add foreign keys
		db.Model(&System{}).AddForeignKey("amahi_user_id", "amahi_users(id)",
			"RESTRICT", "RESTRICT")
		db.Model(&FriendRequest{}).AddForeignKey("amahi_user_id", "amahi_users(id)",
			"RESTRICT", "RESTRICT")
		db.Model(&FriendRequest{}).AddForeignKey("system_id", "systems(id)",
			"RESTRICT", "RESTRICT")
		db.Model(&FriendUser{}).AddForeignKey("amahi_user_id", "amahi_users(id)",
			"RESTRICT", "RESTRICT")
		db.Model(&FriendUser{}).AddForeignKey("system_id", "systems(id)",
			"RESTRICT", "RESTRICT")

		// insert sample data
		db.Create(&AmahiUser{Email: "abc@temp.com"})
		db.Create(&AmahiUser{Email: "bcd@temp.com"})
		db.Create(&AmahiUser{Email: "cde@temp.com"})
		db.Create(&System{AmahiUserID: 1, ApiKey: "abcdef"})
		db.Create(&FriendRequest{AmahiUserID: 1, Pin: "1234", InviteToken: tokenGenerator(),
			SystemID: 1, LastRequestedAt: time.Now(), Status: Accepted})
		db.Create(&FriendUser{AmahiUserID: 1, SystemID: 1})

	}

}
