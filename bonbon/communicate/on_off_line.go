package communicate

import (
	"bonbon/database"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

func initOnline(id int, conn *websocket.Conn) (*user, error) {
	_, err := database.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("id %d login\n", id)

	friendships, err := database.GetFriendships(id)
	if err != nil {
		return nil, err
	}

	onlineLock.Lock()
	if onlineUser[id] == nil {
		onlineUser[id] = &user{
			match:  -1,
			conns:  []*websocket.Conn{conn},
			lock:   new(sync.Mutex),
			bonbon: false,
		}
		for i := 0; i < len(friendships); i++ {
			fmt.Printf("id %d try notify %d he is onlne\n", id, friendships[i].FriendID)
			sendJsonByID(friendships[i].FriendID, StatusCmd{Cmd: "Status", Who: id, Status: "on"})
		}
	} else {
		onlineUser[id].lock.Lock()
		onlineUser[id].conns = append(onlineUser[id].conns, conn)
		onlineUser[id].lock.Unlock()
	}
	onlineLock.Unlock()
	// TODO: 在初始化訊息送到前不開放給別人傳送訊息
	err = sendInitMsg(id)
	if err != nil {
		return nil, err
	}
	return onlineUser[id], nil //此時必定還存在
}

func getInitInfo(id int) (*initCmd, error) {
	account, err := database.GetAccountByID(id)
	if err != nil {
		return &initCmd{Cmd: "init", OK: false}, err
	}
	friendships, err := database.GetFriendships(id)
	if err != nil {
		return &initCmd{Cmd: "init", OK: false}, err
	}
	var friends []friend
	onlineLock.RLock() // 可等資料庫操作結束後再鎖，增進效能
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
			new_firiend := friend{
				ID:     friendships[i].FriendID,
				Sign:   friend_account.Signature,
				Nick:   friendships[i].NickName,
				Status: status,
			}
			friends = append(friends, new_firiend)
		} else {
			return &initCmd{Cmd: "init", OK: false}, err
		}
	}
	onlineLock.RUnlock()
	my_setting := setting{Sign: account.Signature}
	return &initCmd{Cmd: "init", OK: true, Setting: my_setting, Friends: friends}, nil
}

// 實作init訊息
func sendInitMsg(id int) error {
	msg, err := getInitInfo(id)
	if err != nil {
		return err
	}
	err = sendJsonByID(id, msg)
	if err != nil {
		return err
	}
	return nil
}

func clearOffline(id int, conn *websocket.Conn) {
	// 刪除本連線
	onlineLock.Lock()
	u := onlineUser[id]
	u.lock.Lock()
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
		removeFromStrangerQueue(id)
		// 若還在連線
		disconnectByID(id)
		// 傳送離線訊息
		friendships, err := database.GetFriendships(id)
		if err == nil {
			for i := 0; i < len(friendships); i++ {
				fmt.Printf("%d try to notify %d he is offline\n", i, friendships[i].FriendID)
				sendJsonByIDNoLock(friendships[i].FriendID, StatusCmd{Cmd: "Status", Who: id, Status: "off"})
			}
		}
		delete(onlineUser, id)
	}
	u.lock.Unlock()
	onlineLock.Unlock()
	fmt.Printf("id %d 下線了\n", id)
}
