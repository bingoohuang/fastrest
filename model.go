package fastrest

//go:generate easyjson -no_std_marshalers model.go

//easyjson:json
type Rsp struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status,omitempty"`
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

//easyjson:json
type EncryptReq struct {
	TransId   string `json:"transId"`
	AppId     string `json:"appId"`
	KeyId     string `json:"keyId"`
	Mode      string `json:"mode"`
	Padding   string `json:"padding"`
	PlainText string `json:"plainText"`
}

//easyjson:json
type EncryptRsp struct {
	Data string `json:"data"`
}
