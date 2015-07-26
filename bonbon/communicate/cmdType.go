package communicate

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
