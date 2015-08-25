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

func (wq *waitingQueue) match(id int) int {
	onlineUser[id].matchType = wq.Type
	switch wq.Type {
	case "stranger":
		wq.accept = func(s int) bool { return true }
		fmt.Printf("match with stranger, wq.queue length is %d\n", len(wq.queue))
	case "L1_FB_friend":
	case "L2_FB_friend":
	}
	disconnectByID(id)
	for i := 0; i < len(wq.queue); i++ {
		if id == wq.queue[i] {
			return -1
		} else if wq.accept(wq.queue[i]) {
			stranger := wq.queue[i]
			wq.queue = append(wq.queue[:i], wq.queue[i+1:]...)
			globalMatchLock.Lock()
			onlineUser[id].match = stranger
			onlineUser[stranger].match = id
			globalMatchLock.Unlock()
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

// TODO: 將三種類型分開平行處理
// MatchConsumer : 專門進行match的routine
func MatchConsumer() {
	waitingQueues := make(map[string]*waitingQueue)
	waitingQueues["stranger"] = new(waitingQueue)
	waitingQueues["stranger"].queue = make([]int, 0)
	waitingQueues["stranger"].Type = "stranger"
	for {
		var ans int
		req := <-matchRequestChannel
		switch req.Cmd {
		case "in":
			ans = waitingQueues[req.Type].match(req.ID)
		case "out":
			waitingQueues[req.Type].remove(req.ID)
		default:
		}
		matchDoneChannel <- ans
	}
}

// 實作隨機連結(connect) API
func handleConnect(msg []byte, id int, u *user) {
	fmt.Printf("start handle Connect\n")
	var req connectRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		fmt.Printf("unmarshal connect cmd, %s\n", err.Error())
		return
	}
	sendJsonToOnlineID(id, connectResponse{OK: true, Cmd: "connect"})
	matchRequestChannel <- matchRequest{Cmd: "in", ID: id, Type: req.Type}
	stranger := <-matchDoneChannel
	fmt.Printf("stranger is %d\n", stranger)
	// get signatures of both
	selfSignature, err := database.GetSignature(id)
	if err != nil {
		// TODO handle error
	}

	strangerSignature, err := database.GetSignature(stranger)
	if err != nil {
		// TODO handle error
	}

	if stranger != -1 {
		fmt.Printf("%d connect to %d\n", id, stranger)
		sendJsonToUnknownStatusID(
			stranger,
			connectSuccess{Cmd: "connected", Sign: *selfSignature},
			false,
		)
		sendJsonByUserWithLock(u, connectSuccess{Cmd: "connected", Sign: *strangerSignature})
	}
}

func disconnectByID(id int) {
	var stranger int
	globalMatchLock.Lock()
	if stranger = onlineUser[id].match; stranger != -1 {
		onlineUser[id].match = -1
		onlineUser[stranger].match = -1
	}
	globalMatchLock.Unlock()
	// 將io取出鎖外操作
	fmt.Printf("%d disconnect with %d\n", id, stranger)
	if stranger > 0 {
		sendJsonToUnknownStatusID(
			stranger,
			map[string]interface{}{"Cmd": "disconnected"},
			false,
		)
	}
}

// 實作斷線
func handleDisconnect(id int) {
	disconnectByID(id)
	sendJsonToOnlineID(id, map[string]interface{}{"OK": true, "Cmd": "disconnect"})
}
