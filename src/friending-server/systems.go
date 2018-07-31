package main

// System contains the info for a specific HDA.
type System struct {
	BaseModel
	AmahiUserID uint
	ApiKey      string `gorm:"unique"`
	FRs         []FriendRequest
	FUs         []FriendUser
	// other fields can be added
}

// Get system for a given api key from db
func getSystem(apiKey string) (sys *System) {
	db, err := getDb()
	defer db.Close()
	handle(err)
	db.Where("api_key = ?", apiKey).First(sys)
	return
}
