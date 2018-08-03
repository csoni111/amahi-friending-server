package main

import (
	"time"
	"net/http"
	"errors"
	"fmt"
	"crypto/rand"
	"github.com/gorilla/mux"
	"strconv"
	"encoding/json"
)

type RequestStatus int

const (
	Active RequestStatus = iota
	Expired
	Accepted
	Rejected
)

func (e *RequestStatus) MarshalJSON() ([]byte, error) {
	value, ok := map[RequestStatus]string{
		Active: "Active", Expired: "Expired", Accepted: "Accepted", Rejected: "Rejected"}[*e]
	if !ok {
		return nil, errors.New("invalid status value")
	}
	return json.Marshal(value)
}

type FriendRequest struct {
	BaseModel
	Status          RequestStatus `gorm:"default:0" json:"status"`
	AmahiUserID     uint          `gorm:"not null;unique_index:idx_sys_user" json:"-"`
	Pin             string        `gorm:"not null" json:"-"`
	InviteToken     string        `gorm:"unique;not null" json:"-"`
	SystemID        uint          `gorm:"not null;unique_index:idx_sys_user" json:"-"`
	AmahiUser       *AmahiUser    `json:"amahi_user"`
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
		db.Model(&system).Preload("AmahiUser").Related(&reqs)
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

	// check if a friend request for given user already exists
	fr := new(FriendRequest)
	if db.Where("amahi_user_id = ? AND system_id = ?", user.ID, system.ID).First(&fr).RecordNotFound() {
		// set pin, user, invite token and system id
		fr.Pin = nfr.Pin
		fr.AmahiUser = user
		fr.InviteToken = tokenGenerator()
		fr.SystemID = system.ID
		fr.sendEmailNotification()

		// save new entry into database
		if rowsAffected := db.Create(fr).RowsAffected; rowsAffected > 0 {
			respond(w, http.StatusCreated, fr)
		} else {
			respond(w, http.StatusInternalServerError, db.Error)
		}
	} else {
		respond(w, http.StatusBadRequest, "request already exists for user")
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
	if rowsAffected := db.Where("id = ? AND system_id = ?", reqId, system.ID).Delete(&FriendRequest{}).
		RowsAffected; rowsAffected > 0 {
		respond(w, http.StatusOK, "deleted successfully")
	} else {
		respond(w, http.StatusNotFound, errors.New("not found"))
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
		respond(w, http.StatusNotFound, "not found")
	} else {
		fr.InviteToken = tokenGenerator()
		fr.sendEmailNotification()
		db.Save(&fr)
		respond(w, http.StatusOK, "request resent")
	}

}
