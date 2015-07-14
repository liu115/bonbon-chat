package database

import (
	"github.com/jinzhu/gorm"
)

// SQLConnection struct type for database connection hadle
type SQLConnection struct {
	db gorm.DB
}

// Connect start connection to database
func (conn *SQLConnection) Connect() {
	db, err := gorm.Open("sqlite3", "/tmp/bonbon.db")

	if err != nil {
		panic(err)
	}
	conn.db = db
}

// CreateAccount helper function for creating accounts
func (conn *SQLConnection) CreateAccount(email string, username string) {
	// TODO sanity check
	// account := Account{Email: email, Username: username}
}

// RemoveAccount helper function for removing accounts
func (conn *SQLConnection) RemoveAccount(email string, username string) {
	// TODO implementation
}
