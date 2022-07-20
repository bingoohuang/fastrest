package fastrest

//go:generate easyjson -no_std_marshalers model.go

//easyjson:json
type Rsp struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

//easyjson:json
type P1SignReq struct {
	Source  string `json:"source"`
	BizType string `json:"bizType"`
}

//easyjson:json
type P1SignRsp struct {
	Source string `json:"source"`
}
