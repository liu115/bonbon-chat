package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func handleRead(msg []byte, id int, user *user) {
	var request ReadRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		fmt.Printf("unmarshal read cmd, %s\n", err.Error())
		return
	}
	_time := time.Now()

	err = database.UpdateReadTime(id, request.With_who, _time)
	if err != nil {
		fmt.Printf("UpdateReadTime in handleRead, %s\n", err.Error())
		return
	}

	now_str := strconv.FormatInt(_time.UnixNano(), 10)
	sendJsonToOnlineID(id,
		ReadResponse{
			OK:       true,
			Cmd:      "read",
			With_who: request.With_who,
			Time:     now_str,
		})
}
