package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/bingoohuang/golog"

	"github.com/bingoohuang/fastrest"
	_ "github.com/bingoohuang/fastrest/validators/v10" // 引入请求结构体自动校验
)

// InitAssets is the initial assets.
//go:embed initassets
var InitAssets embed.FS

func main() {
	golog.Setup()

	// 注册路由
	router := fastrest.New(map[string]fastrest.Service{
		"GET /":        &fastrest.Version{},
		"GET /status":  &fastrest.Status{},
		"GET /echo":    &fastrest.Echo{},
		"POST /p1sign": &fastrest.P1Sign{},
		"GET /panic":   &fastrest.PanicService{},
	}, fastrest.WithPreProcessor(fastrest.PreProcessorFn(func(dtx *fastrest.Context) error {
		// 全局前置处理器
		return nil
	})), fastrest.WithPostProcessor(fastrest.PostProcessorFn(func(dtx *fastrest.Context) error {
		// 全局后置处理器
		return nil
	})), fastrest.WithPanicProcessor(fastrest.PanicProcessorFn(func(dtx *fastrest.Context, err interface{}) {
		dtx.Ctx.SetStatusCode(http.StatusInternalServerError)
		dtx.Ctx.SetBodyString(fmt.Sprintf("panic: %v", err))
	})),
	)

	args := fastrest.ParseArgs(&InitAssets)
	args.Run(router)
}
