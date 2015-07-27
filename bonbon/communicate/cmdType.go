package communicate

// structure for send
type SendCmd struct {
	Cmd   string
	Who   int
	Msg   string
	Order int
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
	Msg  []string
}

type setting struct {
	Sign string
}

type initMsg struct {
	Cmd     string
	Setting setting
	Friends []friend
}

func getMyInitInfo() {}
