package communicate

// 命名規則
// 若為請求：加上Request後綴
// 若為回應：加上Response後綴
// 若為伺服器主動行為：加上Cmd後綴
// TODO: 使結構符合命名規則

// 告知status
type StatusCmd struct {
	Cmd    string
	Who    int
	Status string
}

// structure for send
type SendRequest struct {
	Cmd   string
	Who   int
	Msg   string
	Order int
}

type SendResponse struct {
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

func respondToSend(req SendRequest, now int64, exist bool) SendResponse {
	res := SendResponse{
		OK:    exist,
		Who:   req.Who,
		Cmd:   req.Cmd,
		Msg:   req.Msg,
		Time:  now,
		Order: req.Order,
	}
	return res
}

// structure for set nickname
type setNickNameRequest struct {
	Cmd      string
	Who      int
	NickName string
}

// structure for set setting
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

// structure for bonbon

// type bonbonRequest struct {
// 	Cmd string
// }

type bonbonResponse struct {
	Cmd string
	OK  bool
}

type newFriendCmd struct {
	Cmd  string
	Who  int
	Nick string
}

type simpleResponse struct {
	OK bool
}

// structure for connect
type connectRequest struct {
	Cmd  string
	Type string
}

type connectResponse struct {
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

type InitCmd struct {
	Cmd     string
	OK      bool
	Setting setting
	Friends []friend
}
