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
	SystemID    uint      `gorm:"not null" json:"-"`
	AmahiUserID uint      `gorm:"not null" json:"-"`
	AmahiUser   AmahiUser `json:"amahi_user"`
}

func getFUs(w http.ResponseWriter, r *http.Request) {
	db, err := getDb()
	defer db.Close()
	handle(err)

	system := checkApiKeyHeader(w, r)
	if system != nil {
		var users []FriendUser
		if err = db.Debug().Joins("JOIN amahi_users as au on au.id = friend_users.amahi_user_id").
			Where("friend_users.system_id = ?", system.ID).Preload("AmahiUser").
			Find(&users).Error; err != nil {
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
	if rowsAffected := db.Debug().Where("amahi_user_id = ? AND system_id = ?", userId, system.ID).
		Delete(&FriendUser{}).RowsAffected;
		rowsAffected > 0 {
		respond(w, http.StatusOK, "deleted")
	} else {
		respond(w, http.StatusNotFound, errors.New("not found"))
	}
}
