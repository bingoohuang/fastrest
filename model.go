package fastrest

//go:generate easyjson model.go

//easyjson:json
type Rsp struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

//easyjson:json
type P1SignReq struct {
	Source  string
	BizType string
}

//easyjson:json
type P1SignRsp struct {
	Source string
}
