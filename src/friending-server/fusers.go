package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

type FriendUser struct {
	BaseModel
	SystemID    uint
	AmahiUserID uint
}

func getFUs(w http.ResponseWriter, r *http.Request) {
	db, err := getDb()
	defer db.Close()
	handle(err)

	system := checkApiKeyHeader(w, r)
	if system != nil {
		var users []AmahiUser
		/*rows, err := db.Raw("SELECT * FROM amahi_users WHERE name = ?", 3).Rows()
		defer rows.Close()
		handle(err)
		for rows.Next() {
			rows.Scan(&name, &age, &email)
		}*/
		db.Debug().Model(&FriendUser{SystemID: system.ID}).Related(&users)
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
	if db.Where("user_id = ? AND system_id = ?", userId, system.ID).Delete(&FriendRequest{}).RecordNotFound() {
		respond(w, http.StatusOK, "deleted")
	} else {
		respond(w, http.StatusNotFound, "not found")
	}
}
