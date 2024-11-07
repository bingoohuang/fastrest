# fastrest

fast restful framework for golang.

![img.png](_doc/architect.png|

1. Create your app directory, like `mkdir myapp; cd myapp; go mod init myapp;`
2. Create initial config.toml in a folder `initassets`, [example](cmd/fastrest/initassets/conf.yml|

   ```yaml
   ---
   addr: ":14142"
   ```

3. Create main code, [example](cmd/fastrest/main.go|

   ```go
   package main
   
   import (
       "embed"
       "github.com/bingoohuang/fastrest"
       _ "github.com/bingoohuang/fastrest/validators/v10" // 引入请求结构体自动校验
   |
   
   // InitAssets is the initial assets.
   //go:embed initassets
   var InitAssets embed.FS
   
   func main(| {
       // 注册路由
       router := fastrest.New(map[string]fastrest.Service{
           "GET /status":  &fastrest.Status{},
           "POST /p1sign": &fastrest.P1Sign{},
       }, fastrest.WithPreProcessor(fastrest.PreProcessorFn(func(dtx *fastrest.Context| error {
           // 全局前置处理器
           return nil
       }||, fastrest.WithPostProcessor(fastrest.PostProcessorFn(func(dtx *fastrest.Context| error {
           // 全局后置处理器
           return nil
       }|||
   
       args := fastrest.ParseArgs(&InitAssets|
       args.Run(router|
   }
   ```

4. Create Makefile, [example](Makefile|
5. Build `make`
6. Create initial conf.toml and ctl: `myapp -init`
7. Startup `./ctl start`, you can set env `export GOLOG_STDOUT=true` before startup to view the log in stdout for
   debugging.
8. Performance testing using [berf](https://github.com/bingoohuang/berf|: `berf :14142/status -d15s -v`
9. Or single test `berf :14142/p1sign -v source=bingoo bizType=abc -n1`

```sh
➜  fastrest git:(main| ✗ berf :14142/p1sign source=bingoo bizType=abc -pRr -n1
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

## easyjson marshalling and unmarshalling

1. Install [easyjson tool](https://github.com/bingoohuang/easyjson|
2. Tag the model, see the following example.
3. Generate easyjson codes: `easyjson yourmodel.go`

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

## 性能测试

1. 空接口 `/status` TPS 30 万.

| # | Hostname             |   Uptime | Uptime Human | Procs | OS    | Platform | Host ID                              | Platform Version | Kernel Version               | Kernel Arch | Os Release                       | Mem Available             | Num CPU | Cpu Mhz | Cpu Model |
|--:|----------------------|---------:|--------------|------:|-------|----------|--------------------------------------|------------------|------------------------------|-------------|----------------------------------|---------------------------|--------:|--------:|-----------|
| 1 | fs04-192-168-126-184 | 14173428 | 5 months     |   373 | linux | centos   | ea4bc56f-c6da-4914-afc6-4d9e54267d41 | 8                | 4.18.0-240.22.1.el8_3.x86_64 | x86_64      | NAME="CentOS Stream" VERSION="8" | 57.25GiB/62.65GiB, 00.91% |      16 |    2300 | Intel(R   | Xeon(R| Gold 5218 CPU @ 2.30GHz |

```sh
[footstone@fs04-192-168-126-184 ~]$ berf :14142/status -c500
Berf benchmarking http://127.0.0.1:14142/status using 500 goroutine(s|, 16 GoMaxProcs.

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
Berf benchmarking http://127.0.0.1:14142/status for 1m0s using 500 goroutine(s|, 16 GoMaxProcs.

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

## AccessLog

1. 设置环境变量 ACCESS_LOG_DIR，启用 access log, e.g. `ACCESS_LOG_DIR=accesslog fastrest`

## 技术特点

fastrest 基于 Go 语言优秀的GPM并发模型之上，叠加多项优化技术：

1. 替换默认的 net/http 为 fasthttp，基准测试性能提升约10倍
2. 替换默认的 encoding/base64 为基于 Turbo-Base64 实现的 cristalhq/base64，是基准测试性能提升约3倍
3. 替换默认的 encoding/json 为基于代码生成免除反射调用的 easyjson，，基准测试性能提升约5倍
4. 调整默认的 最大并发数 GOMAXPROCS 8 为4倍值约30，能达到最高性能

## 支持环境变量

- 设置环境变量 STRUCT_ENV_VERBOSE=1 后，程序运行时会输出所有环境变量，便于调试。
- 设置环境变量 STRUCT_ENV_ON=0 后，可以关闭所有环境变量的支持。

| 变量名称                                             | 类型            | 默认值         | 示例            | 描述                                |
|--------------------------------------------------|---------------|-------------|---------------|-----------------------------------|
| FAST_HTTP_NAME                                   | string        | 空           | FastServer    | 响应头中的服务器名称.                       |
| FAST_HTTP_CONCURRENCY                            | int           | 256 * 1024  | 100           | 服务器最大并发数.                         |
| FAST_HTTP_READ_BUFFER_SIZE                       | int           | 2048        | 4096          | 读缓冲区大小                            |
| FAST_HTTP_WRITE_BUFFER_SIZE                      | int           | 4096        | 写缓冲区大小        |                                   |
| FAST_HTTP_READ_TIMEOUT                           | time.Duration | 0（无超时）      | 10s           | 请求读超时                             |
| FAST_HTTP_WRITE_TIMEOUT                          | time.Duration | 10s         | 5s            | 响应写超时                             |
| FAST_HTTP_IDLE_TIMEOUT                           | time.Duration | 10s         | 30s           | 连接最大空闲时间                          |
| FAST_HTTP_MAX_CONNS_PER_IP                       | int           | 100         | 1000          | 每个IP的最大连接数                        |
| FAST_HTTP_MAX_REQUESTS_PER_CONN                  | int           | 0           | 10000         | 每个连接的最大请求数（达到后关闭连接）               |
| FAST_HTTP_MAX_IDLE_WORKER_DURATION               | time.Duration | 10s         | 30s           | 协程最大空闲时间（达到后退出协程）                 |
| FAST_HTTP_TCP_KEEPALIVE_PERIOD                   | time.Duration | 操作系统默认值     | 60s           | TCP Keepalive 间隔时间                |
| FAST_HTTP_MAX_REQUEST_BODY_SIZE                  | int           | 4194304(4M) | 10485760(10M) | 最大请求体大小                           |
| FAST_HTTP_SLEEP_WHEN_CONCURRENCY_LIMITS_EXCEEDED | time.Duration | 0s          | 3s            | 当并发数超过最大值时，等待一段时间再继续处理请求          |
| FAST_HTTP_DISABLE_KEEPALIVE                      | bool          | false       | true          | 关闭长连接                             |
| FAST_HTTP_TCP_KEEPALIVE                          | bool          | false       | true          | 启用 TCP Keepalive                  |
| FAST_HTTP_REDUCE_MEMORY_USAGE                    | bool          | false       | true          | 减少内存占用                            |
| FAST_HTTP_GET_ONLY                               | bool          | false       | true          | 只处理 GET 请求(用于抗DOS攻击）              |
| FAST_HTTP_DISABLE_PRE_PARSE_MULTIPART_FORM       | bool          | false       | true          | 禁用 multipart/form-data 解析，以减少内存占用 |
| FAST_HTTP_LOG_ALL_ERRORS                         | bool          | false       | true          | 所有错误都记录到日志                        |
| FAST_HTTP_SECURE_ERROR_LOG_MESSAGE               | bool          | false       | true          | 错误日志信息脱敏                          |
| FAST_HTTP_DISABLE_HEADER_NAMES_NORMALIZING       | bool          | false       | true          | 禁用 header 名称标准化                   |
| FAST_HTTP_NO_DEFAULT_SERVER_HEADER               | bool          | false       | true          | 禁用默认的服务器 header                   |
| FAST_HTTP_NO_DEFAULT_DATE                        | bool          | false       | true          | 禁用默认的 Date header                 |
| FAST_HTTP_NO_DEFAULT_CONTENT_TYPE                | bool          | false       | true          | 禁用默认的 Content-Type header         |
| FAST_HTTP_KEEP_HIJACKED_CONNS                    | bool          | false       | true          | 保持 Hijacked 的连接                   |
| FAST_HTTP_CLOSE_ON_SHUTDOWN                      | bool          | false       | true          | 关闭服务时关闭连接                         |
| FAST_HTTP_STREAM_REQUEST_BODY                    | bool          | false       | true          | 流式处理请求体                           |
