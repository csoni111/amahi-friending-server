package main

type AmahiUser struct {
	BaseModel
	Email     string
	Systems   []System     `json:"-"`
	FriendsOf []FriendUser `json:"-"`
}
