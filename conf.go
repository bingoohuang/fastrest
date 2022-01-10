package fastrest

import (
	"github.com/bingoohuang/gg/pkg/v"
	"log"
)

type Arg struct {
	Config  string `flag:"c" usage:"yaml Config filepath"`
	Init    bool   `usage:"init example conf.yml/ctl and then exit"`
	Version bool   `usage:"print version then exit"`

	Addr string `val:":14142"`
}

func (c *Arg) VersionInfo() string { return v.Version() }

func (c *Arg) Run(router *Router) {
	if err := router.Serve(c.Addr); err != nil {
		log.Fatalf("error to serve: %s", err)
	}
}
