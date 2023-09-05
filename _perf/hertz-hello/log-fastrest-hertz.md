# fastrest vs hertz

## environment

`sysinfo -show host -format md`

| # | Hostname              |  Uptime | Uptime Human | Procs | OS    | Platform | Host ID                              | Platform Version | Kernel Version        | Kernel Arch | Os Release                             | Mem Available             | Num CPU | Cpu Mhz | Cpu Model                                |
|--:|-----------------------|--------:|--------------|------:|-------|----------|--------------------------------------|------------------|-----------------------|-------------|----------------------------------------|---------------------------|--------:|--------:|------------------------------------------|
| 1 | localhost.localdomain | 1010700 | 11 days      |   409 | linux | centos   | 00000000-0000-0000-0000-002590c24096 | 7.5.1804         | 3.10.0-862.el7.x86_64 | x86_64      | NAME="CentOS Linux" VERSION="7 (Core)" | 54.51GiB/62.73GiB, 00.87% |      32 |    3300 | Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz |

## log

```sh
# hertz-hello
2022/07/19 08:12:29.593278 engine.go:522: [Debug] HERTZ: Method=GET    absolutePath=/someJSON                 --> handlerName=main.main.func1 (num=2 handlers)
2022/07/19 08:12:29.593427 engine.go:522: [Debug] HERTZ: Method=GET    absolutePath=/moreJSON                 --> handlerName=main.main.func2 (num=2 handlers)
2022/07/19 08:12:29.593447 engine.go:522: [Debug] HERTZ: Method=GET    absolutePath=/pureJson                 --> handlerName=main.main.func3 (num=2 handlers)
2022/07/19 08:12:29.593507 engine.go:522: [Debug] HERTZ: Method=GET    absolutePath=/someData                 --> handlerName=main.main.func4 (num=2 handlers)
2022/07/19 08:12:29.594588 transport.go:91: [Info] HERTZ: HTTP server listening on address=127.0.0.1:8080
^C2022/07/19 08:13:50.931756 hertz.go:72: [Info] HERTZ: Begin graceful shutdown, wait at most num=5 seconds...
# fastrest
^C
```

```sh
# berf :8080/someJSON
Berf benchmarking http://127.0.0.1:8080/someJSON using 100 goroutine(s), 32 GoMaxProcs.

Summary:
  Elapsed                1m2.027s
  Count/RPS   18148325 292583.510
    200       18148325 292583.510
  ReadWrite  393.232 386.210 Mbps

Statistics     Min       Mean     StdDev     Max
  Latency     36µs       334µs     432µs   52.747ms
  RPS       279546.16  292537.75  5870.41  301120.4

Latency Percentile:
  P50     P75    P90    P95     P99     P99.9    P99.99
  264µs  407µs  588µs  748µs  1.283ms  4.061ms  16.795ms
```

```sh
# berf :14142/status
Berf benchmarking http://127.0.0.1:14142/status using 100 goroutine(s), 32 GoMaxProcs.

Summary:
  Elapsed                1m0.868s
  Count/RPS   13185791 216627.319
    200       13185791 216627.319
  ReadWrite  301.545 284.215 Mbps

Statistics     Min       Mean     StdDev      Max
  Latency     36µs       436µs     463µs   22.046ms
  RPS       195293.24  216532.81  3400.55  223771.49

Latency Percentile:
  P50     P75    P90     P95      P99     P99.9   P99.99
  290µs  558µs  949µs  1.294ms  2.147ms  4.412ms  7.332ms
```

## code

[hertz-examples](https://github.com/cloudwego/hertz-examples/blob/main/render/json/main.go)

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	// utils.H is a shortcut for map[string]interface{}
	h.GET("/someJSON", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "hey", "status": consts.StatusOK})
	})

	h.GET("/moreJSON", func(ctx context.Context, c *app.RequestContext) {
		// You also can use a struct
		var msg struct {
			Company  string `json:"company"`
			Location string
			Number   int
		}
		msg.Company = "company"
		msg.Location = "location"
		msg.Number = 123
		// Note that msg.Company becomes "company" in the JSON
		// Will output  :   {"company": "company", "Location": "location", "Number": 123}
		c.JSON(consts.StatusOK, msg)
	})

	h.GET("/pureJson", func(ctx context.Context, c *app.RequestContext) {
		c.PureJSON(consts.StatusOK, utils.H{
			"html": "<p> Hello World </p>",
		})
	})

	h.GET("/someData", func(ctx context.Context, c *app.RequestContext) {
		c.Data(consts.StatusOK, "text/plain; charset=utf-8", []byte("hello"))
	})

	h.Spin()
}
```
