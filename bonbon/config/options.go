package config

import (
	"fmt"
	"github.com/huandu/facebook"
)

// in this file, default values are present unless being set by LoadConfigFile()

// Hostname the hostname server works on
var Hostname = "0.0.0.0"

// Port the port server listens to
var Port = 8080

// Address the address of server
var Address = fmt.Sprintf("%s:%d", Hostname, Port)

// Mode set the mode of server. values can be either "release", "debug", "test"
var Mode = "debug"

// ElectiveNickNames candidate nicknames for naming anonymous friends
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
