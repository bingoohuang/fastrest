package fastrest

import (
	"github.com/bingoohuang/gg/pkg/v"
	"log"
	"os"
	"strings"
)

type Status struct{ DummyService }

func (p *Status) Process(*Context) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功"}, nil
}

type Version struct{ DummyService }

var logEnv = os.Getenv("LOG")

func (p *Version) Process(ctx *Context) (interface{}, error) {
	log.Printf("%s E! version request received from %s", logEnv, ctx.Ctx.RemoteAddr())

	return &Rsp{Status: 200, Message: "成功", Data: map[string]interface{}{
		"gitCommit":  v.GitCommit,
		"buildTime":  v.BuildTime,
		"goVersion":  v.GoVersion,
		"appVersion": v.AppVersion,
	}}, nil
}

type LogType int

const (
	LogOff LogType = iota
	LogOn
	LogAsync
)

var LogTypeEnv = func() LogType {
	switch v := os.Getenv("LOG_TYPE"); strings.ToLower(v) {
	case "0", "off", "no":
		return LogOff
	case "async":
		return LogAsync
	default:
		return LogOn
	}
}()
