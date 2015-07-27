package database

import (
	"bonbon/config"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // provide sqlite3 driver
	"log"
	"math/rand"
)

func init() {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		log.Fatalf("cannot connect to database %v://%v", config.DatabaseDriver, config.DatabaseArgs)
		return
	}

	db.AutoMigrate(&Account{}, &Friendship{})
}

// GetDB start connection to database
func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// CreateAccountByToken struct Account initializer
func CreateAccountByToken(token string) (*Account, error) {
	// get Facebook session
	fbSession := config.GlobalApp.Session(token)

	err := fbSession.Validate()
	if err != nil {
		return nil, err
	}

	// get name and id from Facebook
	res, err := fbSession.Get("/me", nil)

	if err != nil {
		return nil, err
	}

	var facebookID string
	var facebookName string

	res.DecodeField("id", &facebookID)
	res.DecodeField("name", &facebookName)

	// update account database
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var account Account
	query := db.Where("facebook_id = ?", facebookID).First(&account)

	// create account if not exist
	if query.Error != nil {
		account = Account{AccessToken: token,
			FacebookID:   facebookID,
			FacebookName: facebookName,
		}
		query := db.Create(&account)

		if query.Error != nil {
			return nil, query.Error
		}
	} else {
		account.AccessToken = token
		account.FacebookName = facebookName
		query := db.Save(&account)

		if query.Error != nil {
			return nil, query.Error
		}
	}

	return &account, nil
}

// GetAccountByID get account object by id
func GetAccountByID(id int) (*Account, error) {
	// update account database
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var account Account
	query := db.Where("id = ?", id).First(&account)

	// create account if not exist
	if query.Error != nil {
		return nil, query.Error
	}

	return &account, nil
}

// GetFriendships get friends of an account
func GetFriendships(accountID int) ([]Friendship, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var friendShips []Friendship
	query := db.Where("account_id = ?", accountID).Find(&friendShips)

	if query.Error != nil {
		return nil, query.Error
	}

	return friendShips, nil
}

// MakeFriendship establish a friend relation on two accounts
func MakeFriendship(leftID int, rightID int) error {
	// obtain accounts from database
	db, err := GetDB()
	if err != nil {
		return err
	}

	var leftAccount Account
	var rightAccount Account

	query := db.Where("id = ?", leftID).First(&leftAccount)
	if query.Error != nil {
		return query.Error
	}

	query = db.Where("id = ?", rightID).First(&rightAccount)
	if query.Error != nil {
		return query.Error
	}

	// check if both accounts are identical
	if leftAccount.ID == rightAccount.ID {
		return errors.New("database: make friendship of two identical accounts")
	}

	// sanity check on friend relations
	var leftFriends []Friendship
	db.Model(&leftAccount).Related(&leftFriends, "AccountID")

	var rightFriends []Friendship
	db.Model(&rightAccount).Related(&rightFriends, "AccountID")

	if len(leftFriends) > config.NumFriendsLimit || len(rightFriends) > config.NumFriendsLimit {
		return errors.New("database: limit of number of friends is exceeded")
	}

	leftHasFriendship := false
	rightHasFriendship := false

	for _, friendShip := range leftFriends {
		if friendShip.FriendID == rightAccount.ID {
			leftHasFriendship = true
			break
		}
	}

	for _, friendShip := range rightFriends {
		if friendShip.FriendID == leftAccount.ID {
			rightHasFriendship = true
			break
		}
	}

	if leftHasFriendship != rightHasFriendship {
		log.Print("warning: malformed friend relation detected")
	} else if leftHasFriendship && rightHasFriendship {
		return errors.New("database: friendship had been established")
	}

	// append to friends of both
	numElectiveNickNames := len(config.ElectiveNickNames)

	if !leftHasFriendship {
		leftFriendship := Friendship{
			AccountID: leftAccount.ID,
			NickName:  config.ElectiveNickNames[rand.Intn(numElectiveNickNames)],
			FriendID:  rightAccount.ID,
		}
		db.Create(&leftFriendship)
	}

	if !rightHasFriendship {
		rightFriendship := Friendship{
			AccountID: rightAccount.ID,
			NickName:  config.ElectiveNickNames[rand.Intn(numElectiveNickNames)],
			FriendID:  leftAccount.ID,
		}
		db.Create(&rightFriendship)
	}

	return nil
}

// RemoveFriendship remove the friend relation on two accounts
func RemoveFriendship(leftID int, rightID int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	// get accounts from database
	var leftAccount Account
	var rightAccount Account

	query := db.Where("id = ?", leftID).First(&leftAccount)
	if query.Error != nil {
		return query.Error
	}

	query = db.Where("id = ?", rightID).First(&rightAccount)
	if query.Error != nil {
		return query.Error
	}

	// sanity check on friend relations
	leftHasFriendship := true
	rightHasFriendship := true

	var leftFriendship Friendship
	query = db.Where("account_id = ? and friend_id = ?", leftAccount.ID, rightAccount.ID).First(&leftFriendship)
	if query.Error != nil {
		leftHasFriendship = false
	}

	var rightFriendship Friendship
	query = db.Where("account_id = ? and friend_id = ?", rightAccount.ID, leftAccount.ID).First(&rightFriendship)
	if query.Error != nil {
		rightHasFriendship = false
	}

	if leftHasFriendship != rightHasFriendship {
		log.Print("warning: malformed friend relation detected")
	} else if !leftHasFriendship && !rightHasFriendship {
		return errors.New("database: friend relation does not exist")
	}

	// remove relation from friends of both
	if leftHasFriendship {
		db.Delete(&leftFriendship)
	}

	if rightHasFriendship {
		db.Delete(&rightFriendship)
	}

	return nil
}

// GetSignature get the signature string of an account
func GetSignature(id int) (*string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var account Account
	query := db.Where("id = ?", id).First(&account)

	if query.Error != nil {
		return nil, query.Error
	}

	return &account.Signature, nil
}

// SetSignature set the signature string of an account
func SetSignature(id int, signature string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var account Account
	query := db.Where("id = ?", id).First(&account)

	if query.Error != nil {
		return query.Error
	}

	account.Signature = signature
	query = db.Save(&account)

	if query.Error != nil {
		return query.Error
	}

	return nil
}

// SetNickNameOfFriendship set the nickname to a friend
func SetNickNameOfFriendship(accountID int, friendID int, nickName string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var friendship Friendship
	query := db.Where("account_id = ? and friend_id = ?", accountID, friendID).First(&friendship)
	if query.Error != nil {
		return query.Error
	}

	friendship.NickName = nickName
	db.Save(&friendship)
	return nil
}
