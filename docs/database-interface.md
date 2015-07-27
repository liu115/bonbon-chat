## description
In package "bonbon/database", _models.go_ provides object definitions in database, and _database.go_ provides helper functions to manipulate the objects.

## functions
* func CreateAccountByToken(token string) (\*Account, error)
* func GetAccountByID(ID int) (\*Account, error)
* func GetFriendships(accountID int) ([]Friendship, error)
* func MakeFriendship(leftID int, rightID int) error
* func RemoveFriendship(leftID int, rightID int) error

## example usage
```
import (
	...
	"bonbon/database"
)

func foo() {
	// get facebook token from some place
	token := ...

	// create an account regarding to the token
	// if such account exists, it returns the existing one
	account, err := database.CreateAccountByToken(token)

	// get account by id
	account, err = database.GetAccountByID(id)

	// get fields
	fmt.Print( account.ID, account.FacebookID, account.FacebookName, account.AccessToken )

	// get the list of friends of an account
	// it return the slice of struct Friendship defined in models.go
	friendShips, err := database.GetFriendships(account.ID)

	for _, friend := range friendShips {
		fmt.Print(, friend.FriendID, friend.NickName )
	}

	// establish friendship
	err = database.MakeFriendship(leftID, rightID)
	if err != nil {...}

	// remove friendship
	err = database.RemoveFriendship(leftID, rightID)
	if err != nil {...}
}

```
