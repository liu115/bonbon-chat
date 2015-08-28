## description
In package "bonbon/database", _models.go_ provides object definitions in database, and _database.go_ provides helper functions to manipulate defined objects.

## functions
* func CreateAccountByToken(token string) (\*Account, error)
* func GetAccountByID(ID int) (\*Account, error)
* func GetFriendships(accountID int) ([]Friendship, error)
* func MakeFriendship(leftID int, rightID int) error
* func RemoveFriendship(leftID int, rightID int) error
* func SetSignature(id int, signature string) error
* func GetSignature(id int) (\*string, error)
* func SetNickNameOfFriendship(accountID int, friendID int, nickName string) error
* func GetFacebookFriends(id int) ([]Account, error)
* func GetFacebookFriendsOfFriends(id int, degree int) ([]Account, error)
* func AppendActivityLog(accountID int, action string, description string) error

## database table schema
Table schemas is identical the fields in structs defined in _models.go_. Take the struct _Friend_ for instance. The declaration goes as the following.
```
type Friendship struct {
	ID        int     `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AccountID int     `sql:"index"`
	NickName  string
	FriendID  int
}
```

This corresponds to a table _friendship_ in database, with columns _id_, _account\_id_, _nick\_name_ and _friend\_id_. The raw strings (`sql:"index"`, etc) describe the additional properties of the columns.

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

    // update and get the signature of an account
    err = database.SetSignature(id, "geek, rather than git")
    if err != nil {...}

    sign, err := database.GetSignature(id)
    if err != nil {...}

    // append an activity log of an account
    err = database.AppendActivityLog(id, "message", "Do you think Bonbon is censoring our chatting?")
    if err != nil {...}

    // list friends of an account by id
    friendAccounts, err := database.GetFacebookFriends(id)
    if err != nil {...}
    for _, friend := range friendAccounts {
        fmt.Println(friend.ID)
    }

    // list friends of friends up to Nth degree
    fofAccounts, err := database.GetFacebookFriendsOfFriends(id, n
    if err != nil {...}
    for _, account := range fofAccounts {
        fmt.Println(account.ID)
    }
}

```
