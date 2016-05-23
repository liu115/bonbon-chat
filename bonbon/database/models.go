package database

import (
	"time"
)

// Account database model for account
type Account struct {
	ID              int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FacebookID      string
	FacebookName    string
	AccessToken     string
	Signature       string
	Avatar          string
	Friends         []Friendship
	FacebookFriends []byte
	CreateAt        time.Time `sql:"DEFAULT:current_timestamp"`
}

// Friendship database model for a friend relation to another
type Friendship struct {
	ID        int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AccountID int `sql:"index"`
	NickName  string
	LastRead  time.Time
	FriendID  int
}

// Message the log structure of activities
type Message struct {
	ID            int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FromAccountID int
	ToAccountID   int
	Type          int
	Context       string
	Time          time.Time `sql:"DEFAULT:current_timestamp"`
}

// ActivityLog the log structure of activities
type ActivityLog struct {
	ID          int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AccountID   int `sql:"index"`
	Action      string
	Description string
	Time        time.Time `sql:"DEFAULT:current_timestamp"`
}
