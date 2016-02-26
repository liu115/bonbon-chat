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

func checkFriendship(ID1 int, ID2 int) bool {
	check_point := 0
	friendships, err := database.GetFriendships(ID1)
	if err != nil {
		fmt.Printf("in checkFriendship %s", err.Error())
	}
	for _, friendship := range friendships {
		if friendship.FriendID == ID2 {
			check_point += 1
			break
		}
	}

	friendships, err = database.GetFriendships(ID2)
	if err != nil {
		fmt.Printf("in checkFriendship %s", err.Error())
	}
	for _, friendship := range friendships {
		if friendship.FriendID == ID1 {
			check_point += 1
			break
		}
	}

	return (check_point == 2)
}
