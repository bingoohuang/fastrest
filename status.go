package fastrest

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bingoohuang/gg/pkg/v"
)

type Echo struct{ DummyService }

func (p *Echo) Process(c *Context) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功", Data: map[string]interface{}{
		"RemoteAddr": c.Ctx.RemoteAddr().String(),
		"RequestURI": string(c.Ctx.RequestURI()),
		"TimeStamp":  time.Now().Format(http.TimeFormat),
	}}, nil
}

type Status struct{ DummyService }

func (p *Status) Process(*Context) (interface{}, error) {
	return &Rsp{Status: 200, Message: "成功"}, nil
}

type Version struct{ DummyService }

func (p *Version) Process(ctx *Context) (interface{}, error) {
	// log.Printf("%s E! version request received from %s", logEnv, ctx.Ctx.RemoteAddr())
	switch LogTypeEnv {
	case LogOn:
		log.Printf("E! version request received from %s", ctx.Ctx.RemoteAddr())
	case LogAsync:
		log.Printf("[LOG_ASYNC] E! version request received from %s", ctx.Ctx.RemoteAddr())
	}

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
