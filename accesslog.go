package fastrest

import (
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/valyala/fasthttp"
)

// LogPrint logs a message using the given format and optional arguments.
// The usage of format and arguments is similar to that for fmt.Printf().
// LogPrint should be thread safe.
type LogPrint func(format string, a ...interface{})

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// https://github.com/AubSs/fasthttplogger/blob/master/fasthttplogger.go

// Tiny format:
// <method> <url> - <status> - <response-time us>
// GET / - 200 - 11.925 us
func Tiny(req fasthttp.RequestHandler, logFunc LogPrint) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		logFunc("%s %s - %v - %v",
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	}
}

// Short format:
// <remote-addr> | <HTTP/:http-version> | <method> <url> - <status> - <response-time us>
// 127.0.0.1:53324 | HTTP/1.1 | GET /hello - 200 - 44.8µs
func Short(req fasthttp.RequestHandler, logFunc LogPrint) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		logFunc("%v | %s | %s %s - %v - %v",
			GetClientIP(ctx),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	}
}

// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url>  <request body size> - <status> <response body size> <response-time us> | <user-agent>
// [2017-05-31 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello 12345 - 200  12345 48.279µs | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func Combined(req fasthttp.RequestHandler, logFunc LogPrint) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		logFunc("[%s] %s | %s | %s %s %d - %d %d %s | %s",
			end.Format("2006-01-02 15:04:05"),
			GetClientIP(ctx),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			len(ctx.Request.Body()),
			ctx.Response.Header.StatusCode(),
			len(ctx.Response.Body()),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	}
}

// GetClientIP returns the originating IP for a request.
func GetClientIP(ctx *fasthttp.RequestCtx) string {
	ip := string(ctx.Request.Header.Peek("X-Real-IP"))
	if ip == "" {
		ip = string(ctx.Request.Header.Peek("X-Forwarded-For"))
		if ip == "" {
			ip = ctx.RemoteAddr().String()
		}
	}
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// GetDailyLogWriter 日志文件切割，按天
func GetDailyLogWriter(filename string) (*rotatelogs.RotateLogs, error) {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	return rotatelogs.New(
		filename+"_%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
}
