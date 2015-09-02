package communicate

import (
	"errors"
)

func getUserByID(id int) *user {
	onlineLock.RLock()
	u := onlineUser[id]
	onlineLock.RUnlock()
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

func sendJsonByUserWithLock(user *user, json interface{}) error {
	user.lock.Lock()
	sendJsonByUser(user, json)
	user.lock.Unlock()
	return nil
}

func sendJsonToUnknownStatusID(id int, json interface{}, lock bool) error {
	if !lock {
		onlineLock.RLock()
	}
	u := onlineUser[id]
	if !lock {
		onlineLock.RUnlock()
	}
	if u == nil {
		return errors.New("sendJsonByID: ID is offline")
	}
	err := sendJsonByUserWithLock(u, json)
	if err != nil {
		return err
	}
	return nil
}

func sendJsonToOnlineID(id int, json interface{}) error {
	err := sendJsonByUserWithLock(onlineUser[id], json)
	if err != nil {
		return err
	}
	return nil
}
