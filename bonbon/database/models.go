package database

// Account database model for account
type Account struct {
	ID          int    `sql:"AUTO_INCREMENT"`
	Email       string
	Username    string
	AccessToken string
	Friends     []Friend
}

// Friend database model for friend
type Friend struct {
	FriendID int
}

type ChatRoom struct {
	ID int
}
