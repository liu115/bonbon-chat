package database

import (
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// SQLConnection struct type for database connection hadle
type SQLConnection struct {
	db gorm.DB
}

func (conn *DatabaseConnection) connect() {
	db, err := gorm.Open("sqlite3", "/tmp/bonbon.db")

	if err != nil {
		panic(err)
	}
	conn.db = db
}
