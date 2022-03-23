package main

import (
	"embed"

	"github.com/bingoohuang/fastrest"
	_ "github.com/bingoohuang/fastrest/validators/v10" // 引入请求结构体自动校验
)

// InitAssets is the initial assets.
//go:embed initassets
var InitAssets embed.FS

func main() {
	// 注册路由
	router := fastrest.New(map[string]fastrest.Service{
		"GET /":        &fastrest.Version{},
		"GET /status":  &fastrest.Status{},
		"POST /p1sign": &fastrest.P1Sign{},
	}, fastrest.WithPreProcessor(fastrest.PreProcessorFn(func(dtx *fastrest.Context) error {
		// 全局前置处理器
		return nil
	})), fastrest.WithPostProcessor(fastrest.PostProcessorFn(func(dtx *fastrest.Context) error {
		// 全局后置处理器
		return nil
	})))

	args := fastrest.ParseArgs(&InitAssets)
	args.Run(router)
}
