
# fastrest

fast restful framework for golang.

![img.png](_doc/architect.png)

2. Create your app directory, like `mkdir myapp; cd myapp; go mod init myapp;`
3. Create initial config.toml in a folder `initassets`, [example](cmd/fastrest/initassets/conf.yml)
---
---
```yaml
   ---
   addr: ":14142"
```
---
4. Create main code, [example](cmd/fastrest/main.go)
---
---
```go
    package main
    
    import (
    	"embed"
    	"github.com/bingoohuang/fastrest"
    	_ "github.com/bingoohuang/fastrest/validators/v10"  // 引入请求结构体自动校验
    )
    
    // InitAssets is the initial assets.
    //go:embed initassets
    var InitAssets embed.FS
    
    func main() {
    	// 注册路由
    	router := fastrest.New(map[string]fastrest.Service{
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
```
```output
Main function is generated automatically. Please remove func main()
```
---
5. Create Makefile, [example](Makefile)
6. Build `make`
7. Create initial conf.toml and ctl: `myapp -init`
8. Startup `./ctl start`, you can set env `export GOLOG_STDOUT=true` before startup to view the log in stdout for
   debugging.
9. Performance testing using [berf](https://github.com/bingoohuang/berf): `berf :14142/status -d15s -v`
10. Or single test `berf :14142/p1sign -v source=bingoo bizType=abc -n1`

---
```sh
➜  fastrest git:(main) ✗ berf :14142/p1sign source=bingoo bizType=abc -pRr -n1
### 127.0.0.1:63079->127.0.0.1:14142 time: 2022-01-05T14:19:36.312775+08:00 cost: 575.239µs
POST /p1sign HTTP/1.1
User-Agent: blow
Host: 127.0.0.1:14142
Content-Type: application/json; charset=utf-8
Content-Length: 36
Accept-Encoding: gzip, deflate

{"bizType":"abc","source":"bingoo"}

HTTP/1.1 200 OK
Server: fasthttp
Date: Wed, 05 Jan 2022 06:19:36 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 19

{"source":"bingoo"}
```
---
## easyjson marshalling and unmarshalling

1. Install [easyjson tool](https://github.com/bingoohuang/easyjson)
1. Tag the model, see the following example.
2. Generate easyjson codes: `easyjson yourmodel.go`

---
```go
//easyjson:json
type P1SignReq struct {
	Source  string `json:"source"`
	BizType string `json:"bizType"`
}

//easyjson:json
type P1SignRsp struct {
	Source string `json:"source"`
}
```
---

## 性能测试

1. 空接口 `/status` TPS 30 万.

|   # | Hostname             |   Uptime | Uptime Human | Procs | OS    | Platform | Host ID                              | Platform Version | Kernel Version               | Kernel Arch | Os Release                       | Mem Available             | Num CPU | Cpu Mhz | Cpu Model                                |
|----:|----------------------|---------:|--------------|------:|-------|----------|--------------------------------------|------------------|------------------------------|-------------|----------------------------------|---------------------------|--------:|--------:|------------------------------------------|
|   1 | fs04-192-168-126-184 | 14173428 | 5 months     |   373 | linux | centos   | ea4bc56f-c6da-4914-afc6-4d9e54267d41 | 8                | 4.18.0-240.22.1.el8_3.x86_64 | x86_64      | NAME="CentOS Stream" VERSION="8" | 57.25GiB/62.65GiB, 00.91% |      16 |    2300 | Intel(R) Xeon(R) Gold 5218 CPU @ 2.30GHz |

```sh
[footstone@fs04-192-168-126-184 ~]$ berf :14142/status -c500
Berf benchmarking http://127.0.0.1:14142/status using 500 goroutine(s), 16 GoMaxProcs.

Summary:
  Elapsed                 43.346s
  Count/RPS   13280606 306379.932
    200       13280606 306379.932
  ReadWrite  426.481 401.970 Mbps

Statistics     Min       Mean     StdDev     Max
  Latency     27µs      1.612ms   1.773ms  56.755ms
  RPS       285977.74  306261.61  8250.69  326172.9

Latency Percentile:
  P50        P75      P90     P95      P99     P99.9     P99.99
  1.324ms  2.059ms  2.749ms  3.85ms  9.626ms  18.188ms  29.256ms
[footstone@fs04-192-168-126-184 ~]$ berf :14142/status -c500 -d1m
Berf benchmarking http://127.0.0.1:14142/status for 1m0s using 500 goroutine(s), 16 GoMaxProcs.

Summary:
  Elapsed                    1m0s
  Count/RPS   18266731 304441.873
    200       18266731 304441.873
  ReadWrite  423.783 399.428 Mbps

Statistics     Min       Mean     StdDev      Max
  Latency     27µs      1.621ms   1.766ms  69.199ms
  RPS       278032.17  304437.17  9856.28  323089.77

Latency Percentile:
  P50        P75      P90      P95      P99     P99.9     P99.99
  1.376ms  2.065ms  2.713ms  3.702ms  9.401ms  19.398ms  33.875ms
```