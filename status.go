package fastrest

import "time"

type Status struct{ DummyService }

func (p *Status) Process(*Context, string, interface{}) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功", Data: time.Now().UnixNano() / 1e6}, nil
}
