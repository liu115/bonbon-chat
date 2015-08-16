package communicate

// 命名規則
// 若為請求：加上Request後綴
// 若為回應：加上Response後綴
// 若為伺服器主動行為：加上Cmd後綴
// TODO: 使結構符合命名規則

type StatusCmd struct {
	Cmd    string
	Who    int
	Status string
}

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

type Setting struct {
	Sign string
}

type updateSettingsRequest struct {
	Cmd     string
	Setting Setting
}

type updateSettingsResponse struct {
	OK      bool
	Cmd     string
	Setting Setting
}

// type bonbonRequest struct {
// 	Cmd string
// }

type bonbonResponse struct {
	Cmd string
	OK  bool
}

type newFriendFromServer struct {
	Cmd  string
	Who  int
	Nick string
}

type simpleResponse struct {
	OK bool
}

type SendCmdResponse struct {
	OK    bool
	Who   int
	Cmd   string
	Msg   string
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
		Msg:   req.Msg,
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
	ID     int
	Sign   string
	Nick   string
	Status string
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
