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
	Time  string
	Order int
}

type SendFromServer struct {
	Cmd  string
	Who  int
	Time string
	Msg  string
}

func respondToSend(req SendRequest, now string, exist bool) SendResponse {
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

// structure for read
type ReadRequest struct {
	Cmd      string
	With_who int
	Time     string
}

type ReadResponse struct {
	OK       bool
	Cmd      string
	With_who int
	Time     string
}

type ReadFromOtherClient struct {
}

// structure for history
type HistoryMsg struct {
	Text string
	From int
	Time string
}

type HistoryRequest struct {
	Cmd      string
	With_who int
	When     string
	Number   int
	Order    int
}

type HistoryResponse struct {
	Cmd      string
	With_who int
	Number   int
	Order    int
	Msgs     []HistoryMsg
}

// structure for set nickname
type SetNickNameRequest struct {
	Cmd  string
	Who  int
	Nick string
}

type SetNickNameResponse struct {
	OK   bool
	Cmd  string
	Who  int
	Nick string
}

// structure for set setting
type Setting struct {
	Sign   string
	Avatar string
}

type UpdateSettingsRequest struct {
	Cmd     string
	Setting Setting
}

type UpdateSettingsResponse struct {
	OK      bool
	Cmd     string
	Setting Setting
}

// structure for change_sign
type SignCmd struct {
	Cmd  string
	Who  int
	Sign string
}

// structure for change_avatar
type AvatarCmd struct {
	Cmd    string
	Who    int
	Avatar string
}

// structure for bonbon

type BonbonRequest struct {
	Cmd string
}

type BonbonResponse struct {
	Cmd string
	OK  bool
}

type NewFriendCmd struct {
	Cmd  string
	Who  int
	Nick string
}

type simpleResponse struct {
	OK bool
}

// structure for connect
type ConnectRequest struct {
	Cmd  string
	Type string
}

type ConnectResponse struct {
	OK  bool
	Cmd string
}

type ConnectSuccess struct {
	Cmd    string
	Sign   string
	Avatar string
}

// structure for init
type Friend struct {
	ID       int
	Sign     string
	Avatar   string
	Nick     string
	Status   string
	LastRead string
}

type setting struct {
	Sign   string
	Avatar string
}

type InitCmd struct {
	Cmd     string
	OK      bool
	Setting setting
	Friends []Friend
}
