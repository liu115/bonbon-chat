package database

import (
	"time"
	"bonbon/config"
)

// Account database model for account
type Account struct {
	ID          int        `sql:"AUTO_INCREMENT"`
	FacebookID  string
	Name        string
	AccessToken string
	Friends     []Friend
	ChatRooms   []ChatRoom
	IsEnabled   bool       `sql:"DEFAULT:true"`
	CreateAt    time.Time  `sql:"DEFAULT:current_timestamp"`
	LastLoginAt time.Time
}

// NewAccount struct Account initializer
func NewAccount(token string) *Account {
	fbSession := config.GlobalApp.Session(token)

	err := fbSession.Validate()
	if err != nil {
		panic(err)
	}

	res, _ := fbSession.Get("/me", nil)

	var facebookID string
	var facebookName string

	res.DecodeField("id", &facebookID)
	res.DecodeField("name", &facebookName)

	acc := Account{AccessToken: token, Name: facebookName, FacebookID: facebookID}
	return &acc
}

// Friend database model for friend
type Friend struct {
	ID       int `sql:"AUTO_INCREMENT"`
	FriendID int
}

// ChatRoom database model for chat room
type ChatRoom struct {
	ID      int       `sql:"AUTO_INCREMENT"`
	Members []Account
}
