package fastrest

import (
	"embed"
	"encoding/json"
	"github.com/bingoohuang/gg/pkg/ss"
	"log"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/bingoohuang/easyjson"
	"github.com/bingoohuang/easyjson/bytebufferpool"
	"github.com/bingoohuang/easyjson/jwriter"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/sigx"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
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

type PreProcessor interface {
	PreProcess(dtx *Context) error
}
type PreProcessorFn func(dtx *Context) error

func (f PreProcessorFn) PreProcess(dtx *Context) error {
	return f(dtx)
}

type PostProcessor interface {
	PostProcess(dtx *Context) error
}

type PostProcessorFn func(dtx *Context) error

func (f PostProcessorFn) PostProcess(dtx *Context) error {
	return f(dtx)
}

type ErrorProcessor interface {
	ProcessError(dtx *Context, err error) error
}

type ErrorProcessorFn func(dtx *Context, err error) error

func (f ErrorProcessorFn) ProcessError(dtx *Context, err error) error {
	return f(dtx, err)
}

type DummyService struct{}

func (d *DummyService) CreateReq() (interface{}, error)      { return nil, nil }
func (d DummyService) Process(*Context) (interface{}, error) { return nil, nil }

type Router struct {
	routers       map[string]Service
	routerService map[string]string

	Config *RouterConfig
}

type RouterConfig struct {
	PreProcessors   []PreProcessor
	PostProcessors  []PostProcessor
	ErrorProcessors []ErrorProcessor
}

var (
	DefaultPreProcessors   []PreProcessor
	DefaultPostProcessors  []PostProcessor
	DefaultErrorProcessors []ErrorProcessor
)

func RegisterPreProcessor(processors ...PreProcessor) {
	DefaultPreProcessors = append(DefaultPreProcessors, processors...)
}

func RegisterPostProcessor(processors ...PostProcessor) {
	DefaultPostProcessors = append(DefaultPostProcessors, processors...)
}

func RegisterErrorProcessors(processors ...ErrorProcessor) {
	DefaultErrorProcessors = append(DefaultErrorProcessors, processors...)
}

type RouterConfigFn func(*RouterConfig)

func WithPreProcessor(v PreProcessor) RouterConfigFn {
	return func(r *RouterConfig) {
		r.PreProcessors = append(r.PreProcessors, v)
	}
}

func WithPostProcessor(v PostProcessor) RouterConfigFn {
	return func(r *RouterConfig) {
		r.PostProcessors = append(r.PostProcessors, v)
	}
}

func WithErrorProcessor(v ErrorProcessor) RouterConfigFn {
	return func(r *RouterConfig) {
		r.ErrorProcessors = append(r.ErrorProcessors, v)
	}
}

func New(m map[string]Service, fns ...RouterConfigFn) *Router {
	config := &RouterConfig{
		PreProcessors:   append([]PreProcessor{}, DefaultPreProcessors...),
		PostProcessors:  append([]PostProcessor{}, DefaultPostProcessors...),
		ErrorProcessors: append([]ErrorProcessor{}, DefaultErrorProcessors...),
	}
	for _, fn := range fns {
		fn(config)
	}
	return &Router{
		Config:        config,
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

	err := r.handleService(dtx)
	if err == nil {
		return
	}

	for _, p := range r.Config.ErrorProcessors {
		if err = p.ProcessError(dtx, err); err == nil {
			break
		}
	}

	if err != nil {
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

	for _, p := range r.Config.PreProcessors {
		if err := p.PreProcess(dtx); err != nil {
			return err
		}
	}

	if p, ok := s.(PreProcessor); ok {
		if err := p.PreProcess(dtx); err != nil {
			return err
		}
	}

	if v, ok := req.(easyjson.Unmarshaler); ok {
		if pt, err := easyjson.UnmarshalPool(Pool, dtx.Ctx.Request.Body(), v); pt != nil {
			dtx.AppendPoolReturner(pt)
		} else if err != nil {
			return err
		}
	} else if req != nil && ss.Contains(string(dtx.Ctx.Request.Header.Peek("Content-Type")), "json") {
		if err := json.Unmarshal(dtx.Ctx.Request.Body(), req); err != nil {
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

	for i := len(r.Config.PostProcessors) - 1; i >= 0; i-- {
		if err := r.Config.PostProcessors[i].PostProcess(dtx); err != nil {
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
	sigx.RegisterSignalProfile()

	maxProcs := int(4 * float64(runtime.GOMAXPROCS(0)))
	runtime.GOMAXPROCS(maxProcs)
	log.Printf("Changed runtime.GOMAXPROCS to %d", maxProcs)

	return c
}
