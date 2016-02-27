package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
	"time"
)

func handleHistory(msg []byte, id int, user *user) {
	var request historyRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		fmt.Printf("unmarshal history cmd, %s\n", err.Error())
		return
	}
	msgs, err := database.GetMessagesBeforeTime(id, request.With_who, time.Unix(0, request.When), request.Number)
	if err != nil {
		fmt.Printf("getMessagesBeforeTime in history cmd, %s\n", err.Error())
		return
	}

	msgs_modify := make([]historyMsg, len(msgs))
	for i, msg := range msgs {
		msgs_modify[i] = historyMsg{Text: msg.Context, From: msg.FromAccountID, Time: msg.Time.UnixNano()}
	}

	sendJsonToOnlineID(id,
		historyResponse{
			Cmd:      "history",
			Order:    request.Order,
			With_who: request.With_who,
			Number:   len(msgs),
			Msgs:     msgs_modify,
		},
		false)
}
