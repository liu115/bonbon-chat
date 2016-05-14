package communicate

import (
	"bonbon/database"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
)

func initOnline(id int, conn *websocket.Conn) (*user, error) {
	_, err := database.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("id %d login\n", id)

	friendships, err := database.GetFriendships(id)
	if err != nil {
		fmt.Printf("getFriendships: %s", err.Error())
		return nil, err
	}

	if onlineUser[id] == nil {
		onlineUser[id] = &user{
			match:     -1,
			matchType: "",
			conns:     []*websocket.Conn{conn},
			bonbon:    false,
		}
		err = sendInitMsg(id)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(friendships); i++ {
			fmt.Printf("id %d try notify %d he is online\n", id, friendships[i].FriendID)
			sendJsonToUnknownStatusID(
				friendships[i].FriendID,
				StatusCmd{Cmd: "status", Who: id, Status: "on"},
			)
		}
	} else {
		onlineUser[id].conns = append(onlineUser[id].conns, conn)
		err = sendInitMsg(id)
		if err != nil {
			return nil, err
		}
	}
	return onlineUser[id], nil //此時必定還存在
}

func getInitInfo(id int) (*InitCmd, error) {
	account, err := database.GetAccountByID(id)
	if err != nil {
		return &InitCmd{Cmd: "init", OK: false}, err
	}
	friendships, err := database.GetFriendships(id)
	if err != nil {
		return &InitCmd{Cmd: "init", OK: false}, err
	}
	var friends = []Friend{}
	for i := 0; i < len(friendships); i++ {
		// 這邊的檢查可能可以容錯高一點
		friend_account, err := database.GetAccountByID(friendships[i].FriendID)
		var status string
		if onlineUser[friendships[i].FriendID] == nil {
			status = "off"
		} else {
			status = "on"
		}
		if err == nil {
			new_firiend := Friend{
				ID:       friendships[i].FriendID,
				Sign:     friend_account.Signature,
				Nick:     friendships[i].NickName,
				Status:   status,
				LastRead: strconv.FormatInt(friendships[i].LastRead.UnixNano(), 10),
			}
			friends = append(friends, new_firiend)
		} else {
			return &InitCmd{Cmd: "init", OK: false}, err
		}
	}
	my_setting := setting{Sign: account.Signature}
	return &InitCmd{Cmd: "init", OK: true, Setting: my_setting, Friends: friends}, nil
}

// 實作init訊息
func sendInitMsg(id int) error {
	msg, err := getInitInfo(id)
	if err != nil {
		return err
	}
	err = sendJsonToOnlineID(id, msg)
	if err != nil {
		return err
	}
	return nil
}

func clearOffline(id int, conn *websocket.Conn) {
	// 刪除本連線
	u := onlineUser[id]
	if u == nil {
		fmt.Printf("id %d 下線了\n", id)
		return
	}
	conns := u.conns
	var which int
	for i := 0; i < len(conns); i++ {
		if conn == conns[i] {
			which = i
			break
		}
	}
	u.conns = append(conns[:which], conns[which+1:]...)
	if len(u.conns) == 0 {
		// 若還在等待陌生人
		matchRequestChannel <- matchRequest{Cmd: "out", ID: id, Type: u.matchType}
		<-matchDoneChannel
		// 若還在連線
		disconnectByID(id)
		// 傳送離線訊息
		friendships, err := database.GetFriendships(id)
		if err == nil {
			for i := 0; i < len(friendships); i++ {
				fmt.Printf("%d try to notify %d he is offline\n", i, friendships[i].FriendID)
				sendJsonToUnknownStatusID(
					friendships[i].FriendID,
					StatusCmd{Cmd: "status", Who: id, Status: "off"},
				)
			}
		} else {
			fmt.Printf("getFriendships: %s", err.Error())
		}
		delete(onlineUser, id)
	}
	fmt.Printf("id %d 下線了\n", id)
}
