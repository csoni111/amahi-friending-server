package main

import (
	"time"
	"net/http"
	"errors"
	"fmt"
	"crypto/rand"
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
	Status          RequestStatus `gorm:"default:0"`
	Email           string
	Pin             string
	InviteToken     string `gorm:"unique"`
	SystemID        uint
	LastRequestedAt time.Time
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

	db, err := getDb()
	defer db.Close()
	handle(err)
	var user AmahiUser
	if db.Where("email = ?", nfr.Email).Take(&user).RecordNotFound() {
		return errors.New("no such user exists")
	}
	return nil
}

func getFRs(w http.ResponseWriter, r *http.Request) {
	db, err := getDb()
	defer db.Close()
	handle(err)

	system := checkApiKeyHeader(w, r)
	if system != nil {
		var reqs []FriendRequest
		db.Model(&system).Related(&reqs)
		respond(w, http.StatusOK, reqs)
	}
}

func newFR(w http.ResponseWriter, r *http.Request) {
	// validate ApiKey from headers
	system := checkApiKeyHeader(w, r)

	// validate post data
	var nfr NewFR
	if err := decode(r, &nfr); err != nil {
		respond(w, http.StatusBadRequest, err)
	}

	// set pin, email, invite token and system id
	var fr FriendRequest
	fr.Pin = nfr.Pin
	fr.Email = nfr.Email
	fr.InviteToken = tokenGenerator()
	fr.SystemID = system.ID

	// save new entry into database
	db, err := getDb()
	defer db.Close()
	handle(err)
	db.Create(fr)
}

func removeFR(w http.ResponseWriter, r *http.Request) {

}

func resendFR(w http.ResponseWriter, r *http.Request) {

}
