package fastrest

type Encrypt struct{ DummyService }

func (p *Encrypt) CreateReq() (interface{}, error) {
	return &EncryptReq{}, nil
}

func (p *Encrypt) Process(dtx *Context) (interface{}, error) {
	return p.process(dtx, dtx.Req.(*EncryptReq))
}

func (p *Encrypt) process(dtx *Context, req *EncryptReq) (interface{}, error) {
	return &EncryptRsp{Data: req.PlainText}, nil
}
