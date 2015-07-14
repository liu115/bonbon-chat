package database

import (
	"time"
)

// Account database model for account
type Account struct {
	ID          int    `sql:"AUTO_INCREMENT"`
	Email       string
	Username    string
	AccessToken string
	Friends     []Friend
	IsObsolete  bool
	CreateAt    time.Time
	LastLoginAt time.Time
}

// Friend database model for friend
type Friend struct {
	FriendID int
}

type ChatRoom struct {
	ID int
}
