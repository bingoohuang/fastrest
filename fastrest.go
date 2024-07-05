package fastrest

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/fastrest/fgrpc/service"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/bingoohuang/easyjson"
	"github.com/bingoohuang/easyjson/bytebufferpool"
	"github.com/bingoohuang/easyjson/jwriter"
	"github.com/bingoohuang/fastrest/fgrpc/server"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/iox"
	"github.com/bingoohuang/gg/pkg/sigx"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/soheilhy/cmux"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Context struct {
	Req interface{}
	Rsp interface{}

	Ctx         *fasthttp.RequestCtx
	ServiceName string

	Returners []bytebufferpool.PoolReturner
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

type Service interface {
	CreateReq() (interface{}, error)
	Process(dtx *Context) (interface{}, error)
}

type PanicProcessor interface {
	PanicProcess(dtx *Context, recovered interface{})
}
type PanicProcessorFn func(dtx *Context, recovered interface{})

func (f PanicProcessorFn) PanicProcess(dtx *Context, recovered interface{}) {
	f(dtx, recovered)
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

const (
	ResultSendFile = "__ResultSendFile"
)

type DummyService struct{}

func (d *DummyService) CreateReq() (interface{}, error)       { return nil, nil }
func (d *DummyService) Process(*Context) (interface{}, error) { return nil, nil }

type Router struct {
	routers       map[string]Service
	routerService map[string]string

	Config *RouterConfig
}

type RouterConfig struct {
	PanicProcessor  PanicProcessor
	AccessLogDir    string
	PreProcessors   []PreProcessor
	PostProcessors  []PostProcessor
	ErrorProcessors []ErrorProcessor

	NotFoundHandler    func(dtx *Context)
	MaxRequestBodySize int

	*fasthttp.Server
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

func WithNotFoundHandler(v func(ctx *Context)) RouterConfigFn {
	return func(r *RouterConfig) {
		r.NotFoundHandler = v
	}
}

func WithAccessLogDir(v string) RouterConfigFn {
	return func(r *RouterConfig) {
		r.AccessLogDir = v
	}
}

func WithPanicProcessor(v PanicProcessor) RouterConfigFn {
	return func(r *RouterConfig) {
		r.PanicProcessor = v
	}
}

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

// WithMaxRequestBodySize set Maximum request body size.
func WithMaxRequestBodySize(maxRequestBodySize int) RouterConfigFn {
	return func(r *RouterConfig) {
		r.MaxRequestBodySize = maxRequestBodySize
	}
}

// WithFastHTTPServer set customized fasthttp.Server, warn: Handler will be ignored
func WithFastHTTPServer(server *fasthttp.Server) RouterConfigFn {
	return func(r *RouterConfig) {
		r.Server = server
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

func IsEnvOff(name string) bool {
	s := strings.ToLower(os.Getenv(name))
	return s == "0" || s == "off" || s == "no"
}

func (r *Router) Serve(port string, reusePort bool) error {
	log.Printf("Start to ListenAndServe %s", port)

	listenFn := net.Listen
	if reusePort {
		listenFn = reuseport.Listen
	}

	ln, err := listenFn("tcp4", port)
	if err != nil {
		return err
	}

	return r.ServeListener(ln)
}

func (r *Router) ServeListener(ln net.Listener) error {
	m := cmux.New(ln)
	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	// Java gRPC Clients: Java gRPC client blocks until it receives a SETTINGS frame from the server.
	// If you are using the Java client to connect to a cmux'ed gRPC server please match with writers:
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	// All the rest is assumed to be HTTP
	httpL := m.Match(cmux.Any())

	grpcS := grpc.NewServer()
	service.RegisterServiceServer(grpcS, &server.GrpcServer{})

	if !IsEnvOff("GRPC_REFLECTION") {
		reflection.Register(grpcS)
	}
	handle := r.handle

	if r.Config.AccessLogDir != "" {
		accessLogWriter, err := GetDailyLogWriter(filepath.Join(r.Config.AccessLogDir, "access.log"))
		if err != nil {
			return err
		}
		defer iox.Close(accessLogWriter)

		accessLogger := log.New(accessLogWriter, "", 0)
		handle = Combined(r.handle, accessLogger.Printf)
	}

	httpS := r.Config.Server
	if httpS == nil {
		httpS = &fasthttp.Server{}
	}

	httpS.Handler = handle
	if r.Config.MaxRequestBodySize > 0 {
		httpS.MaxRequestBodySize = r.Config.MaxRequestBodySize
	}

	if httpS.MaxConnsPerIP == 0 {
		httpS.MaxConnsPerIP = 100
	}
	if httpS.IdleTimeout == 0 {
		httpS.IdleTimeout = 10 * time.Second
	}
	if httpS.ReadTimeout == 0 {
		httpS.ReadTimeout = 10 * time.Second
	}
	if httpS.WriteTimeout == 0 {
		httpS.WriteTimeout = 10 * time.Second
	}
	if httpS.ReadBufferSize == 0 {
		httpS.ReadBufferSize = 2048
	}

	// Use the muxed listeners for your servers.
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	return m.Serve()
}

func (r *Router) recover(dtx *Context) {
	recovered := recover()
	if recovered == nil {
		return
	}

	log.Printf("E! recovered: %+v", recovered)

	if p := r.Config.PanicProcessor; p != nil {
		p.PanicProcess(dtx, recovered)
		return
	}

	dtx.Ctx.SetStatusCode(http.StatusInternalServerError)
	dtx.Ctx.SetBodyString(fmt.Sprintf("%+v", recovered))
}

func (r *Router) handle(ctx *fasthttp.RequestCtx) {
	dtx := &Context{Ctx: ctx}
	defer dtx.Release()
	defer r.recover(dtx)

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
		if r.Config.NotFoundHandler != nil {
			r.Config.NotFoundHandler(dtx)
		} else {
			dtx.Ctx.NotFound()
		}
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

	if err := unmarshalJSON(dtx, req); err != nil {
		return err
	}

	dtx.Rsp, err = s.Process(dtx)
	if err != nil {
		return err
	}

	if dtx.Rsp == ResultSendFile {
		return nil
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

	data, err := marshalJSON(dtx, dtx.Rsp)
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
	if svc, ok := r.routers[key]; ok {
		return r.routerService[key], svc
	}

	if svc, ok := r.routers[path]; ok {
		return r.routerService[path], svc
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

	if c.MaxProcs > 0 {
		runtime.GOMAXPROCS(c.MaxProcs)
		log.Printf("Changed runtime.GOMAXPROCS to %d", c.MaxProcs)
	}

	return c
}

func unmarshalJSON(dtx *Context, req interface{}) error {
	if req == nil {
		return nil
	}

	if v, ok := req.(easyjson.Unmarshaler); ok {
		pt, err := easyjson.UnmarshalPool(Pool, dtx.Ctx.Request.Body(), v)
		if pt != nil {
			dtx.AppendPoolReturner(pt)
		}

		return err
	}

	if ss.Contains(string(dtx.Ctx.Request.Header.Peek("Content-Type")), "json") {
		return json.Unmarshal(dtx.Ctx.Request.Body(), req)
	}

	return nil
}

func marshalJSON(c *Context, rsp interface{}) (data []byte, err error) {
	if rsp == nil {
		return nil, nil
	}

	var pt bytebufferpool.PoolReturner
	if v, ok := rsp.(easyjson.Marshaler); ok {
		data, pt, err = easyjson.MarshalPool(Pool, v)
		c.AppendPoolReturner(pt)
		return data, err
	}

	return json.Marshal(rsp)
}
