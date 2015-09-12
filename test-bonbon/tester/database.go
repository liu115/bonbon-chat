package main

import (
	"bonbon/config"
	"bonbon/database"
	"fmt"
	"github.com/jinzhu/gorm"
)

func clearDB() error {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		return fmt.Errorf("cannot connect to database %v://%v", config.DatabaseDriver, config.DatabaseArgs)
	}
	db.DropTable(&database.Account{})
	db.DropTable(&database.Friendship{})
	database.InitDatabase()
	return nil
}

func createAccount(ID int, signature string) error {
	user := database.Account{ID: ID, Signature: signature}
	db, err := database.GetDB()
	if err != nil {
		return err
	}
	db.Create(&user)
	return nil
}
