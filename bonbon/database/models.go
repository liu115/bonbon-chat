package database

import (
	// "fmt"
	"time"
	"bonbon/config"
)

// Account database model for account
type Account struct {
	ID           int        `sql:"AUTO_INCREMENT"`
	FacebookID   string
	FacebookName string
	AccessToken  string
	Friends      []Friend
	ChatRooms    []ChatRoom
	IsEnabled    bool       `sql:"DEFAULT:true"`
	CreateAt     time.Time  `sql:"DEFAULT:current_timestamp"`
}

// GetAccount struct Account initializer
func GetAccount(token string) *Account {
	// get Facebook session
	fbSession := config.GlobalApp.Session(token)

	err := fbSession.Validate()
	if err != nil {
		return nil
	}

	// get name and id from Facebook
	res, _ := fbSession.Get("/me", nil)

	var facebookID string
	var facebookName string

	res.DecodeField("id", &facebookID)
	res.DecodeField("name", &facebookName)

	// update account database
	db := GetDB()
	user := Account{}
	query := db.Where("facebook_id = ?", facebookID).First(&user)

	// create account exists
	if query.Error != nil {
		user = Account{AccessToken: token,
			FacebookID: facebookID,
			FacebookName: facebookName,
			IsEnabled: true,
		}
		query := db.Create(&user)

		if query.Error != nil {
			return nil
		}
	} else {
		user.AccessToken = token
		user.FacebookName = facebookName
		query := db.Save(&user)

		if query.Error != nil {
			return nil
		}
	}

	return &user
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
