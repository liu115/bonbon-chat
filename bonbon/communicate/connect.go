package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
)

// TODO: 建造waitingStrangerQueue map，分別對應stranger, L1, L2
// 每個都是queue，只是匹配的限制不同，stranger就是每個都匹配
// removeFromStramgerQueue時，移除在此類型等待的id
// 必須在user裡紀錄試圖match的類型

type matchRequest struct {
	ID int
	// Cmd: in, out
	Cmd string
	// relation: stranger, L1_FB_friend, L2_FB_friend
	Type string
}

var matchRequestChannel = make(chan matchRequest)
var matchDoneChannel = make(chan int)

type waitingQueue struct {
	queue  []int
	Type   string
	accept func(int) bool
}

func inAcconts(accounts []*database.Account) func(int) bool {
	return func(s int) bool {
		for _, friend := range accounts {
			if friend.ID == s {
				return true
			}
		}
		return false
	}
}

func strangerAccept(id int) func(int) bool {
	friendShips, err := database.GetFriendships(id)
	if err != nil {
		fmt.Printf("in stranger Accept, %s", err.Error())
	}
	return func(s int) bool {
		for _, friend := range friendShips {
			if friend.FriendID == s {
				return false
			}
		}
		return true
	}
}

func L1_FB_friendAccept(id int) func(int) bool {
	L1_FB_friends, err := database.GetFacebookFriends(id)
	if err != nil {
		fmt.Printf("in L1_FB_friend  Accept, %s", err.Error())
	}
	friendShips, err := database.GetFriendships(id)
	if err != nil {
		fmt.Printf("in L1_FB_friend Accept, %s", err.Error())
	}
	return func(s int) bool {
		for _, friend := range friendShips {
			if friend.FriendID == s {
				return false
			}
		}
		for _, friend := range L1_FB_friends {
			if friend.ID == s {
				return true
			}
		}
		return false
	}
}

func L2_FB_friendAccept(id int) func(int) bool {
	L2_FB_friends, err := database.GetFacebookFriendsOfFriends(id, 2)
	if err != nil {
		fmt.Printf("in L2_FB_friend  Accept, %s", err.Error())
	}
	friendShips, err := database.GetFriendships(id)
	if err != nil {
		fmt.Printf("in L2_FB_friend Accept, %s", err.Error())
	}
	return func(s int) bool {
		for _, friend := range friendShips {
			if friend.FriendID == s {
				return false
			}
		}
		for _, friend := range L2_FB_friends {
			if friend.ID == s {
				return true
			}
		}
		return false
	}
}

func (wq *waitingQueue) match(id int) int {
	onlineUser[id].matchType = wq.Type
	switch wq.Type {
	case "stranger":
		wq.accept = strangerAccept(id)
	case "L1_FB_friend":
		wq.accept = L1_FB_friendAccept(id)
	case "L2_FB_friend":
		wq.accept = L2_FB_friendAccept(id)
	}
	disconnectByID(id)
	for i := 0; i < len(wq.queue); i++ {
		if id == wq.queue[i] {
			fmt.Printf("queue is %s\n", wq.queue)
			return -1
		} else if wq.accept(wq.queue[i]) {
			stranger := wq.queue[i]
			wq.queue = append(wq.queue[:i], wq.queue[i+1:]...)
			onlineUser[id].match = stranger
			onlineUser[id].matchType = ""
			onlineUser[stranger].match = id
			onlineUser[stranger].matchType = ""
			fmt.Printf("queue is %s\n", wq.queue)
			return stranger
		}
	}
	wq.queue = append(wq.queue, id)
	fmt.Printf("queue is %s\n", wq.queue)
	return -1
}

func (wq *waitingQueue) remove(id int) {
	for i := 0; i < len(wq.queue); i++ {
		if wq.queue[i] == id {
			wq.queue = append(wq.queue[:i], wq.queue[i+1:]...)
			break
		}
	}
}

func createWaitingQueue(Type string) *waitingQueue {
	w := new(waitingQueue)
	w.queue = make([]int, 0)
	w.Type = Type
	return w
}

// TODO: 將三種類型分開平行處理
// MatchConsumer : 專門進行match的routine
func MatchConsumer() {
	waitingQueues := make(map[string]*waitingQueue)
	waitingQueues["stranger"] = createWaitingQueue("stranger")
	waitingQueues["L1_FB_friend"] = createWaitingQueue("L1_FB_friend")
	waitingQueues["L2_FB_friend"] = createWaitingQueue("L2_FB_friend")
	for {
		var ans int
		req := <-matchRequestChannel
		if req.Type != "" {
			switch req.Cmd {
			case "in":
				fmt.Printf("queue in\n")
				if t := onlineUser[req.ID].matchType; t != "" {
					waitingQueues[t].remove(req.ID)
				}
				ans = waitingQueues[req.Type].match(req.ID)
			case "out":
				fmt.Printf("queue out\n")
				waitingQueues[req.Type].remove(req.ID)
				onlineUser[req.ID].matchType = ""
			default:
			}
		}
		matchDoneChannel <- ans
	}
}

// 實作隨機連結(connect) API
func handleConnect(msg []byte, id int, u *user) {
	fmt.Printf("start handle Connect\n")
	var req ConnectRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		fmt.Printf("unmarshal connect cmd, %s\n", err.Error())
		return
	}
	sendJsonToOnlineID(id, ConnectResponse{OK: true, Cmd: "connect"})
	matchRequestChannel <- matchRequest{Cmd: "in", ID: id, Type: req.Type}
	stranger := <-matchDoneChannel
	fmt.Printf("stranger is %d\n", stranger)
	if stranger != -1 {
		// get signatures of both
		selfSignature, err := database.GetSignature(id)
		if err != nil {
			// TODO handle error
		}

		strangerSignature, err := database.GetSignature(stranger)
		if err != nil {
			// TODO handle error
		}

		fmt.Printf("%d connect to %d\n", id, stranger)
		sendJsonToUnknownStatusID(
			stranger,
			ConnectSuccess{Cmd: "connected", Sign: *selfSignature},
		)
		sendJsonByUser(u, ConnectSuccess{Cmd: "connected", Sign: *strangerSignature})
	}
}

func disconnectByID(id int) bool {
	var stranger int
	disconnect := false
	if stranger = onlineUser[id].match; stranger != -1 {
		disconnect = true
		onlineUser[id].match = -1
		onlineUser[stranger].match = -1
	}
	fmt.Printf("%d disconnect with %d\n", id, stranger)
	if stranger > 0 {
		sendJsonToUnknownStatusID(
			stranger,
			map[string]interface{}{"Cmd": "disconnected"},
		)
	}
	return disconnect
}

// 實作斷線
func handleDisconnect(id int, u *user) {
	dis := disconnectByID(id)
	if u.matchType != "" {
		matchRequestChannel <- matchRequest{Cmd: "out", ID: id, Type: u.matchType}
		<-matchDoneChannel
	}
	sendJsonToOnlineID(id, map[string]interface{}{"OK": dis, "Cmd": "disconnect"})
}
