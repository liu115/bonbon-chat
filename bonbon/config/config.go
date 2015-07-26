package config

import (
	"github.com/huandu/facebook"
)

// ElectiveNickNames
var ElectiveNickNames = []string{"Yuri", "Mighty", "7122"}

// NumFriendsLimit maximum number of friends per account
var NumFriendsLimit = 50

// FBAppID App ID for Facebook graph api
var FBAppID = "915780538494020"

// FBAppSecret App Secret for Facebook graph api
var FBAppSecret = "d7b698973175d799a27e0d129f9e19ba"

// GlobalApp global app instance
var GlobalApp = facebook.New(FBAppID, FBAppSecret)

// DatabaseDriver driver string for gorm.Open()
var DatabaseDriver = "sqlite3"

// DatabaseArgs arguments for gorm.Open()
var DatabaseArgs = "/tmp/bonbon.db"
