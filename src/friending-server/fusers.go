package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"errors"
)

type FriendUser struct {
	BaseModel
	SystemID    uint `gorm:"not null"`
	AmahiUserID uint `gorm:"not null"`
}

func getFUs(w http.ResponseWriter, r *http.Request) {
	db, err := getDb()
	defer db.Close()
	handle(err)

	system := checkApiKeyHeader(w, r)
	if system != nil {
		var users []AmahiUser
		if err = db.Joins("JOIN friend_users as fu on fu.amahi_user_id = amahi_users.id").
			Where("fu.system_id = ?", system.ID).Find(&users).Error; err != nil {
			log.Fatal(err)
			respond(w, http.StatusInternalServerError, err)
			return
		}
		respond(w, http.StatusOK, users)
	}
}

func removeFU(w http.ResponseWriter, r *http.Request) {
	// validate ApiKey from headers
	system := checkApiKeyHeader(w, r)
	if system == nil {
		return
	}

	// get request id from url vars
	userId, err := strconv.Atoi(mux.Vars(r)["id"])

	// attempt to delete and send response
	db, err := getDb()
	defer db.Close()
	handle(err)
	if rowsAffected := db.Debug().Where("amahi_user_id = ? AND system_id = ?", userId, system.ID).Delete(&FriendRequest{}).RowsAffected;
	 rowsAffected > 0 {
		respond(w, http.StatusOK, "deleted")
	} else {
		respond(w, http.StatusNotFound, errors.New("not found"))
	}
}
