package fastrest

import (
	"github.com/bingoohuang/gg/pkg/v"
	"log"
)

type Status struct{ DummyService }

func (p *Status) Process(*Context) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功"}, nil
}

type Version struct{ DummyService }

func (p *Version) Process(ctx *Context) (interface{}, error) {
	log.Printf("version request received from %s", ctx.Ctx.RemoteAddr())

	return &Rsp{Status: 200, Message: "成功", Data: map[string]interface{}{
		"gitCommit":  v.GitCommit,
		"buildTime":  v.BuildTime,
		"goVersion":  v.GoVersion,
		"appVersion": v.AppVersion,
	}}, nil
}
