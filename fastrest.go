package fastrest

import (
	"embed"
	"encoding/json"
	"github.com/bingoohuang/easyjson"
	"github.com/bingoohuang/easyjson/bytebufferpool"
	"github.com/bingoohuang/easyjson/jwriter"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/sigx"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

type Context struct {
	Returners   []bytebufferpool.PoolReturner
	ServiceName string

	Req interface{}
	Rsp interface{}

	Ctx *fasthttp.RequestCtx
}

func (c *Context) Release() {
	for _, returner := range c.Returners {
		returner.ReturnPool()
	}
}

var Pool = &bytebufferpool.Pool{}

func (c *Context) ApplyPoolBuf(size int) []byte {
	buf := Pool.Get(bytebufferpool.WithMinSize(size))
	c.AppendPoolReturner(bytebufferpool.PoolReturnFn(func() { Pool.Put(buf) }))
	(*jwriter.Slice)(unsafe.Pointer(&buf.B)).Len = size // increase length manually
	return buf.B
}

func (c *Context) AppendPoolReturner(r bytebufferpool.PoolReturner) {
	if r != nil {
		c.Returners = append(c.Returners, r)
	}
}

func (c *Context) createRspData() (data []byte, err error) {
	var pt bytebufferpool.PoolReturner
	if v, ok := c.Rsp.(easyjson.Marshaler); ok {
		data, pt, err = easyjson.MarshalPool(Pool, v)
		c.AppendPoolReturner(pt)
		return data, err
	} else if c.Rsp != nil {
		return json.Marshal(c.Rsp)
	}
	return nil, nil
}

type Service interface {
	CreateReq() (interface{}, error)
	Process(dtx *Context) (interface{}, error)
}

type PostProcessor interface {
	PostProcess(dtx *Context) error
}

type DummyService struct{}

func (d *DummyService) CreateReq() (interface{}, error)      { return nil, nil }
func (d DummyService) Process(*Context) (interface{}, error) { return nil, nil }

type Router struct {
	routers       map[string]Service
	routerService map[string]string
}

func New(m map[string]Service) *Router {
	return &Router{
		routers:       m,
		routerService: makeRouterServices(m),
	}
}

func (r *Router) Serve(port string) error {
	log.Printf("Start to ListenAndServe %s", port)
	ln, err := reuseport.Listen("tcp4", port)
	if err != nil {
		return err
	}

	return fasthttp.Serve(ln, r.handle)
}

func (r *Router) handle(ctx *fasthttp.RequestCtx) {
	dtx := &Context{Ctx: ctx}
	defer dtx.Release()

	if err := r.handleService(dtx); err != nil {
		log.Printf("E! failed to handleService, error: %v", err)
		ctx.SetStatusCode(500)
	}
}

func (r *Router) handleService(dtx *Context) error {
	serviceName, s := r.findService(dtx)
	if s == nil {
		dtx.Ctx.NotFound()
		return nil
	}

	dtx.ServiceName = serviceName
	req, err := s.CreateReq()
	if err != nil {
		return err
	}
	dtx.Req = req

	if v, ok := req.(easyjson.Unmarshaler); ok {
		if pt, err := easyjson.UnmarshalPool(Pool, dtx.Ctx.Request.Body(), v); pt != nil {
			dtx.AppendPoolReturner(pt)
		} else if err != nil {
			return err
		}
	}

	dtx.Rsp, err = s.Process(dtx)
	if err != nil {
		return err
	}

	if p, ok := s.(PostProcessor); ok {
		if err := p.PostProcess(dtx); err != nil {
			return err
		}
	}

	data, err := dtx.createRspData()
	if err != nil {
		return err
	}

	dtx.Ctx.SetContentType("application/json; charset=utf-8")
	_, err = dtx.Ctx.Write(data)
	return err
}

func (r *Router) findService(dtx *Context) (string, Service) {
	path := string(dtx.Ctx.Path())
	method := string(dtx.Ctx.Method())
	key := method + " " + path
	if service, ok := r.routers[key]; ok {
		return r.routerService[key], service
	}

	if service, ok := r.routers[path]; ok {
		return r.routerService[path], service
	}

	return "", nil
}

func makeRouterServices(r map[string]Service) map[string]string {
	m := make(map[string]string)

	for k, v := range r {
		name := reflect.TypeOf(v).Elem().String()
		if p := strings.LastIndexByte(name, '.'); p > 0 {
			name = name[p+1:]
		}
		m[k] = name
	}

	return m
}

func ParseArgs(initFiles *embed.FS) Arg {
	c := Arg{}
	flagparse.Parse(&c,
		flagparse.AutoLoadYaml("c", "conf.yml"),
		flagparse.ProcessInit(initFiles))

	log.Printf("Arg parsed: %+v", c)

	// 注册性能采集信号，用法:
	// 第一步，通知开始采集：touch jj.cpu; kill -USR1 `pidof dsvs2`;
	// 第二部，压力测试开始（或者其他手工测试，等待程序运行一段时间，比如5分钟）
	// 第三步，通知结束采集，生成 cpu.profile 文件，命令与第一步相同
	// 第四步，下载 cpu.profile 文件，`go tool pprof -http :9402 cpu.profile` 开启浏览器查看
	sigx.RegisterSignalProfile(nil)

	maxProcs := int(4 * float64(runtime.GOMAXPROCS(0)))
	runtime.GOMAXPROCS(maxProcs)
	log.Printf("Changed runtime.GOMAXPROCS to %d", maxProcs)

	return c
}
