package fastrest

import (
	"log"

	"github.com/bingoohuang/gg/pkg/v"
)

type Arg struct {
	Config string `flag:"c" usage:"yaml Config filepath"`

	Addr      string `val:":14142"`
	Init      bool   `usage:"init example conf.yml/ctl and then exit"`
	Version   bool   `usage:"print version then exit"`
	ReusePort bool   `usage:"Reuse port"`
	MaxProcs  int    `usage:"GOMAXPROCS"`
}

func (c *Arg) VersionInfo() string { return v.Version() }

func (c *Arg) Run(router *Router) {
	if err := router.Serve(c.Addr, c.ReusePort); err != nil {
		log.Fatalf("error to serve: %s", err)
	}
}
