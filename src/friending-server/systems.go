package main

// System contains the info for a specific HDA.
type System struct {
	BaseModel
	AmahiUserID uint   `gorm:"not null"`
	ApiKey      string `gorm:"unique;not null"`
	Frs         []FriendRequest
	Fus         []FriendUser
}

// Get system for a given api key from db
func getSystem(apiKey string) (sys *System) {
	db, err := getDb()
	defer db.Close()
	handle(err)
	sys = new(System)
	if db.Where("api_key = ?", apiKey).First(sys).RecordNotFound() {
		return nil
	}
	return
}
