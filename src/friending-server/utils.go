package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
	"encoding/json"
	"log"
	"time"
	"errors"
)

type BaseModel struct {
	ID        uint      `gorm:"primary_key; auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"bad request"`
}

// ok represents types capable of validating themselves.
type ok interface {
	OK() error
}

// common error types
// missing field error in api calls
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return string(e) + " is required"
}

// to decode the post body and return errors if any
func decode(r *http.Request, v ok) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return v.OK()
}

// responds to a api request with given data, string or error
// in json format along with given status code
func respond(w http.ResponseWriter, status int, data interface{}) {
	msg := data
	if err, ok := data.(error); ok {
		msg = HTTPError{Code: status, Message: err.Error()}
	} else if str, ok := data.(string); ok {
		msg = map[string]interface{}{"code": status, "message": str}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
	log.Println("->", status, data)
}

// opens a new connection to db and returns the instance
// caller is responsible for closing the db connection after use
func getDb() (db *gorm.DB, err error) {
	return gorm.Open("mysql", DbCreds)
}

// checks for a valid api key and returns corresponding HDA or sends 401
func checkApiKeyHeader(w http.ResponseWriter, r *http.Request) (sys *System) {
	apiKey := r.Header.Get(ApiKeyHeader)
	// TODO check API Key validity
	if len(apiKey) > 0 {
		sys = getSystem(apiKey)
		if sys != nil {
			return
		}
	}
	respond(w, http.StatusUnauthorized, errors.New("authorization failed"))
	return
}
