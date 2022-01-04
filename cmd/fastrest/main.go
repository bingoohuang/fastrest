package main

import (
	"embed"
	"github.com/bingoohuang/fastrest"
)

// InitAssets is the initial assets.
//go:embed initassets
var InitAssets embed.FS

func main() {
	// 注册路由
	router := fastrest.New(map[string]fastrest.Service{
		"GET /status":  &fastrest.Status{},
		"POST /p1sign": &fastrest.P1Sign{},
	})

	args := fastrest.ParseArgs(&InitAssets)
	args.Run(router)
}
