package main

type AmahiUser struct {
	BaseModel
	Email       string       `gorm:"not null"`
	Systems     []System     `json:"-"`
	FriendUsers []FriendUser `json:"-"`
}
