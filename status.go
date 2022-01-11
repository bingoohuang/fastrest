package fastrest

type Status struct{ DummyService }

func (p *Status) Process(*Context) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功"}, nil
}
