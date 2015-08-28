package database

import (
	"fmt"
	"log"
	"errors"
	"bytes"
	"math/rand"
	"encoding/binary"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // provide sqlite3 driver
	"bonbon/config"
)

// InitDatabase the database package initialization function
func InitDatabase() error {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		return fmt.Errorf("cannot connect to database %v://%v", config.DatabaseDriver, config.DatabaseArgs)
	}

	db.AutoMigrate(&Account{}, &Friendship{})
	return nil
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

	// update Facebook friends
	err = UpdateFacebookFriends(account.ID)
	if err != nil {
		return nil, err
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

// UpdateFacebookFriends refresh the Facebook friends of an account
func UpdateFacebookFriends(id int) error {
	// get account
	db, err := GetDB()
	if err != nil {
		return err
	}

	var account Account
	query := db.Where("id = ?", id).First(&account)
	if query.Error != nil {
		return query.Error
	}

	// get friend ids from Facebook Graph API
	fbSession := config.GlobalApp.Session(account.AccessToken)
	res, err := fbSession.Get("/me/friends", nil)
	if err != nil {
		return err
	}

	var facebookFriendIDs []string

	paging, err := res.Paging(fbSession)
	if err != nil {
		return err
	}

	for {
		data := paging.Data()

		for _, item := range data {
			facebookFriendIDs = append(facebookFriendIDs, item["id"].(string))
		}

		noMore, err := paging.Next()
		if err != nil {
			return err
		} else if noMore {
			break
		}
	}

	// find Facebook friends in database
	var accountFriends []Account
	query = db.Where("facebook_id in (?)", facebookFriendIDs).Find(&accountFriends)
	if query.Error != nil {
		return err
	}

	// store friend ids in binary format
	// friend ids are stored in the form [number of friends] [friend id 1] [friend id 2] ... [friend id n]
	// each bracket pair [ ] here is a 64-bit little endian integer
	// e.g. a list of friend ids 1, 5, 6 is formatted as
	// 0x03 00 00 00 00 00 00 00  0x01 00 00 00 00 00 00 00  0x05 00 00 00 00 00 00 00  0x06 00 00 00 00 00 00 00
	bufFriendIDs := new(bytes.Buffer)
	binary.Write(bufFriendIDs, binary.LittleEndian, int64(len(accountFriends)))

	for _, friend := range accountFriends {
		binary.Write(bufFriendIDs, binary.LittleEndian, int64(friend.ID))
	}

	account.FacebookFriends = bufFriendIDs.Bytes()
	db.Save(&account)

	return nil
}

// GetFacebookFriends get a list of friends of an account
func GetFacebookFriends(id int) ([]*Account, error) {
	// get account
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var account Account
	query := db.Where("id = ?", id).First(&account)
	if query.Error != nil {
		return nil, query.Error
	}

	// parse friend ids from blob stored in database
	bufFriendIDs := bytes.NewBuffer(account.FacebookFriends)

	var numFacebookFriends int64
	err = binary.Read(bufFriendIDs, binary.LittleEndian, &numFacebookFriends)
	if err != nil {
		return nil, err
	}

	friendIDs := make([]int, int(numFacebookFriends))
	for i := 0; i < int(numFacebookFriends); i++ {
		var currFriendID int64
		err = binary.Read(bufFriendIDs, binary.LittleEndian, &currFriendID)
		if err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, int(currFriendID))
	}

	var friendAccounts []Account
	query = db.Where("id in (?)", friendIDs).Find(&friendAccounts)
	if query.Error != nil {
		return nil, query.Error
	}

	var friendPointerAccounts []*Account
	for _, account := range friendAccounts {
		friendPointerAccounts = append(friendPointerAccounts, &account)
	}

	return friendPointerAccounts, nil
}

// GetFacebookFriendsOfFriends get friends of friends up to N degrees of separation
func GetFacebookFriendsOfFriends(id int, degree int) ([]*Account, error) {
	if degree < 1 {
		return nil, fmt.Errorf("invalid degree: degree = %d", degree)
	}

	account, err := GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	// run BFS
	openFriends := make(map[int]*Account)
	closedFriends := make(map[int]*Account)

	openFriends[id] = account

	for i := 0; i <= degree; i++ {
		newOpenFriends := make(map[int]*Account)

		for friendID, friendAccount := range openFriends {
			if _, ok := closedFriends[friendID]; !ok {
				closedFriends[friendID] = friendAccount

				neighborFriends, err := GetFacebookFriends(friendID)
				if err != nil {
					return nil, err
				}

				for _, neighborAccount := range neighborFriends {
					newOpenFriends[neighborAccount.ID] = neighborAccount
				}
			}
		}

		openFriends = newOpenFriends
	}

	delete(closedFriends, id)

	var friendsOfFriends []*Account
	for _, friendAccount := range closedFriends {
		friendsOfFriends = append(friendsOfFriends, friendAccount)
	}

	return friendsOfFriends, nil
}

// AppendActivityLog push a new activity log to database
func AppendActivityLog(accountID int, action string, description string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	log := ActivityLog{
		AccountID: accountID,
		Action: action,
		Description: description,
	}

	db.Create(&log)
	return nil
}
