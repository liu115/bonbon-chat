package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
)

// -1代表目前無人
var waitingStranger = -1
var StrangerLock = new(sync.Mutex)

func removeFromStrangerQueue(id int) {
	StrangerLock.Lock()
	if id == waitingStranger {
		waitingStranger = -1
	}
	StrangerLock.Unlock()
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
	var stranger = -1
	switch req.Type {
	case "stranger":
		StrangerLock.Lock()
		disconnectByID(id)
		// 此時必為無配對，因為只能在strangerLock內建立match，而剛剛消除了match
		if waitingStranger == -1 || waitingStranger == id {
			waitingStranger = id
		} else {
			stranger = waitingStranger // 若在stranger為waitingStranger，則在strangerLock內不可能消失
			waitingStranger = -1
			globalMatchLock.Lock()
			u.match = stranger
			onlineUser[stranger].match = id
			globalMatchLock.Unlock()
		}
		StrangerLock.Unlock()

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

	case "L1_FB_friend":
	case "L2_FB_friend":
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
