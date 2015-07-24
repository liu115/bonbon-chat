package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // provide sqlite3 driver
	"bonbon/config"
)

var isDatabaseInitialized = false

// GetDB start connection to database
func GetDB() (gorm.DB) {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		panic(err)
	}

	if !isDatabaseInitialized {
		db.AutoMigrate(&Account{}, &Friend{}, &ChatRoom{})
		isDatabaseInitialized = true
	}
	return db
}
