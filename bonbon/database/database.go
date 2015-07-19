package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // provide sqlite3 driver
)

// GetDB start connection to database
func GetDB() (gorm.DB) {
	db, err := gorm.Open("sqlite3", "/tmp/bonbon.db")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Account{}, &Friend{}, &ChatRoom{})
	return db
}
