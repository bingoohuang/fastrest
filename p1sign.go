package fastrest

type P1Sign struct{ DummyService }

func (p *P1Sign) CreateReq() (interface{}, error) {
	return &P1SignReq{}, nil
}

func (p *P1Sign) Process(dtx *Context) (interface{}, error) {
	return p.process(dtx, dtx.Req.(*P1SignReq))
}

func (p *P1Sign) process(dtx *Context, req *P1SignReq) (interface{}, error) {
	return &P1SignRsp{Source: req.Source}, nil
}
