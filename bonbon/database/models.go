package database

import (
	"time"
)

// Account database model for account
type Account struct {
	ID              int        `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FacebookID      string
	FacebookName    string
	AccessToken     string
	Friends         []Friendship
	CreateAt        time.Time  `sql:"DEFAULT:current_timestamp"`
}

// Friendship database model for a friend relation to another
type Friendship struct {
	ID        int     `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AccountID int     `sql:"index"`
	NickName  string
	FriendID  int
}
