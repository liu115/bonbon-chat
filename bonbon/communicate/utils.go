package communicate

import (
	"errors"
)

func sendJsonByUser(user *user, json interface{}) error {
	user.lock.Lock()
	l := len(user.conns)
	errMsg := ""
	for i := 0; i < l; i++ {
		err := user.conns[i].WriteJSON(json)
		if err != nil {
			errMsg += (err.Error() + " ")
		}
	}
	user.lock.Unlock()
	if errMsg != "" {
		return errors.New("sendJsonByUser: " + errMsg)
	}
	return nil
}

func sendJsonByID(id int, json interface{}) error {
	u := onlineUser[id]
	err := sendJsonByUser(u, json)
	if err != nil {
		return err
	}
	return nil
}
