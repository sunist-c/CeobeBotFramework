package service

import (
	"github.com/sunist-c/CeobeBotFramework/infrastructure/logging"
	"time"
)

type Middleware func(ctx *Context)

func DefaultServiceTraceMiddleware(ctx *Context) {
	start := time.Now()
	ctx.ctx.Next()
	t := time.Now().Sub(start)
	ip := ctx.ctx.RemoteIP()
	method := ctx.ctx.Request.Method
	logging.Info("dial request", "%v - from: %v, use: %v", method, ip, t.String())
}
