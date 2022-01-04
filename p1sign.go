package fastrest

type P1Sign struct{ DummyService }

func (p *P1Sign) CreateReq() (interface{}, error) {
	return &P1SignReq{}, nil
}
func (p *P1Sign) Process(dtx *Context, serviceName string, r interface{}) (interface{}, error) {
	return p.process(dtx, serviceName, r.(*P1SignReq))
}
func (p *P1Sign) process(dtx *Context, serviceName string, req *P1SignReq) (interface{}, error) {
	return &P1SignRsp{Source: req.Source}, nil
}
