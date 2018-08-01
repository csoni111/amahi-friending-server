package main

type AmahiUser struct {
	BaseModel
	Email       string       `gorm:"not null" json:"email"`
	Systems     []System     `json:"-"`
	FriendUsers []FriendUser `json:"-"`
}
