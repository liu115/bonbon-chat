package communicate

import (
	"errors"
)

func getUserByID(id int) (*user, error) {
	onlineLock.RLock()
	u := onlineUser[id]
	onlineLock.RUnlock()
	if u == nil {
		return nil, errors.New("getUserByID: user not online")
	} else {
		return u, nil
	}
}

func sendJsonByUser(user *user, json interface{}) error {
	user.lock.Lock()
	l := len(user.conns)
	if l == 0 {
		user.lock.Unlock()
		return errors.New("sendJsonByUser: user has no conns")
	}
	for i := 0; i < l; i++ {
		user.conns[i].WriteJSON(json)
	}
	user.lock.Unlock()
	return nil
}

func sendJsonByID(id int, json interface{}) error {
	u, err := getUserByID(id)
	if err != nil {
		return errors.New("sendJsonByID, " + err.Error())
	}
	err = sendJsonByUser(u, json)
	if err != nil {
		return errors.New("sendJsonByID, " + err.Error())
	}
	return nil
}
