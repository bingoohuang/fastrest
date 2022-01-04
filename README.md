# fastrest

fast restful framework for golang.

2. Create your app directory, like `mkdir myapp; cd myapp; go mod init myapp;`
3. Create  initial config.toml in a folder `initassets`, [example](cmd/fastrest/initassets/conf.yml)
   ```yaml
   ---
   addr: ":14142"
    ```
4. Create main code, [example](cmd/fastrest/main.go)
   ```go
   package main

   import (      
       "embed"
       "github.com/bingoohuang/fastrest"
   )

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
   ```
5. Create Makefile, [example](Makefile)
6. `make`
7. Create initial conf.toml and ctl: `myapp -init`
8. Startup `./ctl start`, you can set env `export GOLOG_STDOUT=true` before startup to view the log in stdout for debugging.
9. Performance testing using [berf](https://github.com/bingoohuang/berf): `berf :14142/status -d15s -v`

## easyjson marshalling and unmarshalling

1. Install [easyjson tool](https://github.com/bingoohuang/easyjson)
1. Tag the model, see the following example.
2. `easyjson yourmodel.go`

```go
//easyjson:json
type Service1Req struct {
	Source  string
	BizType string
}

//easyjson:json
type Service1Rsp struct {
	Source string
}
```
