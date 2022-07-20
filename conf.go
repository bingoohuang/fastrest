package fastrest

import (
	"log"

	"github.com/bingoohuang/gg/pkg/v"
)

type Arg struct {
	Config    string `flag:"c" usage:"yaml Config filepath"`
	Init      bool   `usage:"init example conf.yml/ctl and then exit"`
	Version   bool   `usage:"print version then exit"`
	ReusePort bool   `usage:"Reuse port"`
	Sonic     bool   `usage:"Using sonic JSON serializing & deserializing library, accelerated by JIT (just-in-time compiling) and SIMD (single-instruction-multiple-data)."`

	Addr string `val:":14142"`
}

func (c *Arg) VersionInfo() string { return v.Version() }

func (c *Arg) Run(router *Router) {
	if c.Sonic && !router.Config.UsingSonic {
		router.Config.UsingSonic = true
	}

	if err := router.Serve(c.Addr, c.ReusePort); err != nil {
		log.Fatalf("error to serve: %s", err)
	}
}
