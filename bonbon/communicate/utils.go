package communicate

import (
	"errors"
)

func getUserByID(id int) *user {
	u := onlineUser[id]
	return u
}

func sendJsonByUser(user *user, json interface{}) error {
	l := len(user.conns)
	errMsg := ""
	for i := 0; i < l; i++ {
		err := user.conns[i].WriteJSON(json)
		if err != nil {
			errMsg += (err.Error() + " ")
		}
	}
	if errMsg != "" {
		return errors.New("sendJsonByUser: " + errMsg)
	}
	return nil
}

func sendJsonToUnknownStatusID(id int, json interface{}) error {
	u := onlineUser[id]
	if u == nil {
		return errors.New("sendJsonByID: ID is offline")
	}
	err := sendJsonByUser(u, json)
	if err != nil {
		return err
	}
	return nil
}

func sendJsonToOnlineID(id int, json interface{}) error {
	u := onlineUser[id]
	err := sendJsonByUser(u, json)
	if err != nil {
		return err
	}
	return nil
}
