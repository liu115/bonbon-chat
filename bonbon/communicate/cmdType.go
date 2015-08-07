package communicate

import (
	"bonbon/database"
)

// structure for send
type SendCmd struct {
	Cmd   string
	Who   int
	Msg   string
	Order int
}

type setNickNameRequest struct {
	Cmd      string
	Who      int
	NickName string
}

type updateSettingsRequest struct {
	Cmd  string
	Settings struct {
		Signature string
	}
}

type simpleResponse struct {
	OK bool
}

type SendCmdResponse struct {
	OK    bool
	Who   int
	Cmd   string
	Time  int64
	Order int
}

type SendFromServer struct {
	Cmd  string
	Who  int
	Time int64
	Msg  string
}

func respondToSend(req SendCmd, now int64, exist bool) SendCmdResponse {
	res := SendCmdResponse{
		OK:    exist,
		Who:   req.Who,
		Cmd:   req.Cmd,
		Time:  now,
		Order: req.Order,
	}
	return res
}

// structure for connect
type connectCmd struct {
	Cmd  string
	Type string
}

type connectCmdResponse struct {
	OK  bool
	Cmd string
}

type connectSuccess struct {
	Cmd  string
	Sign string
}

// structure for init
type friend struct {
	ID   int
	Sign string
	Nick string
}

type setting struct {
	Sign string
}

type initCmd struct {
	Cmd     string
	OK      bool
	Setting setting
	Friends []friend
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
	for i := 0; i < len(friendships); i++ {
		// 這邊的檢查可能可以容錯高一點
		friend_account, err := database.GetAccountByID(friendships[i].FriendID)
		if err == nil {
			new_firiend := friend{
				ID:   friendships[i].FriendID,
				Sign: friend_account.Signature,
				Nick: friendships[i].NickName,
			}
			friends = append(friends, new_firiend)
		} else {
			return &initCmd{Cmd: "init", OK: false}, err
		}
	}
	my_setting := setting{Sign: account.Signature}
	return &initCmd{Cmd: "init", OK: true, Setting: my_setting, Friends: friends}, nil
}
