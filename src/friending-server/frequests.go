package main

import (
	"time"
	"net/http"
	"errors"
	"fmt"
	"crypto/rand"
	"github.com/gorilla/mux"
	"strconv"
)

type RequestStatus int

const (
	Active RequestStatus = iota
	Expired
	Accepted
	Rejected
)

type FriendRequest struct {
	BaseModel
	Status          RequestStatus `gorm:"default:0" json:"status"`
	AmahiUserID     uint          `gorm:"not null" json:"-"`
	Pin             string        `gorm:"not null" json:"-"`
	InviteToken     string        `gorm:"unique;not null" json:"-"`
	SystemID        uint          `gorm:"not null" json:"-"`
	AmahiUser       AmahiUser     `json:"amahi_user"`
	LastRequestedAt time.Time     `json:"last_requested_at"`
}

type NewFR struct {
	Email string
	Pin   string
}

func tokenGenerator() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (nfr *NewFR) OK() error {
	if len(nfr.Email) == 0 {
		return ErrMissingField("email")
	}
	if len(nfr.Pin) == 0 {
		return ErrMissingField("pin")
	}
	return nil
}

func (fr *FriendRequest) sendEmailNotification() {
	fr.LastRequestedAt = time.Now()
	// TODO send email to friend specifying his pin
}

func getFRs(w http.ResponseWriter, r *http.Request) {
	db, err := getDb()
	defer db.Close()
	handle(err)

	system := checkApiKeyHeader(w, r)
	if system != nil {
		var reqs []FriendRequest
		db.Debug().Model(&system).Preload("amahi_user").Related(&reqs)
		respond(w, http.StatusOK, reqs)
	}
}

func newFR(w http.ResponseWriter, r *http.Request) {
	// validate ApiKey from headers
	system := checkApiKeyHeader(w, r)
	if system == nil {
		return
	}

	// validate post data
	var nfr NewFR
	if err := decode(r, &nfr); err != nil {
		respond(w, http.StatusBadRequest, err)
		return
	}

	// open db connection
	db, err := getDb()
	defer db.Close()
	handle(err)

	// validate friend's email
	user := new(AmahiUser)
	if db.Where("email = ?", nfr.Email).Take(user).RecordNotFound() {
		respond(w, http.StatusBadRequest, errors.New("no such user exists"))
		return
	}

	// set pin, email, invite token and system id
	fr := new(FriendRequest)
	fr.Pin = nfr.Pin
	fr.AmahiUserID = user.ID
	fr.InviteToken = tokenGenerator()
	fr.SystemID = system.ID
	fr.sendEmailNotification()

	// save new entry into database
	if rowsAffected := db.Create(fr).RowsAffected; rowsAffected > 0 {
		respond(w, http.StatusCreated, "created")
	} else {
		respond(w, http.StatusInternalServerError, db.Error)
	}
}

func removeFR(w http.ResponseWriter, r *http.Request) {
	// validate ApiKey from headers
	system := checkApiKeyHeader(w, r)
	if system == nil {
		return
	}

	// get request id from url vars
	reqId, err := strconv.Atoi(mux.Vars(r)["id"])

	// attempt to delete and send response
	db, err := getDb()
	defer db.Close()
	handle(err)
	if db.Where("id = ? AND system_id = ?", reqId, system.ID).Delete(&FriendRequest{}).RecordNotFound() {
		respond(w, http.StatusOK, "deleted")
	} else {
		respond(w, http.StatusNotFound, "not found")
	}
}

func resendFR(w http.ResponseWriter, r *http.Request) {
	// validate ApiKey from headers
	system := checkApiKeyHeader(w, r)
	if system == nil {
		return
	}

	// get request id from url vars
	reqId, err := strconv.Atoi(mux.Vars(r)["id"])

	// resend email notification
	db, err := getDb()
	defer db.Close()
	handle(err)
	var fr FriendRequest
	if db.Where("id = ? AND system_id = ?", reqId, system.ID).First(&fr).RecordNotFound() {
		fr.InviteToken = tokenGenerator()
		fr.sendEmailNotification()
		db.Save(&fr)
		respond(w, http.StatusOK, "request resent")
	} else {
		respond(w, http.StatusNotFound, "not found")
	}

}
