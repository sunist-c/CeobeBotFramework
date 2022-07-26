package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/CeobeBotFramework/infrastructure/logging"
	"net/http"
)

type IEngine interface {
	Serve(address string)
	Group(url string) IRouter
	Close()
}

type Engine struct {
	exitChan chan struct{}
	e        *gin.Engine
}

func (e *Engine) Group(url string) IRouter {
	g := e.e.Group(url)
	return &Router{router: g}
}

func (e *Engine) serveFunction(address string) {
	go e.e.Run(address)
	logging.Info("started http service", "")
	for {
		select {
		case <-e.exitChan:
			logging.Info("closed http service", "")
			return
		}
	}
}

func (e *Engine) Serve(address string) {
	go e.serveFunction(address)
}

func (e *Engine) Close() {
	logging.Info("closing http service", "")
	e.exitChan <- struct{}{}
}

func DefaultNoRouteFunction(ctx *gin.Context) {
	logging.Info("not found", "%v - from: %v, url: %v", ctx.Request.Method, ctx.ClientIP(), ctx.Request.URL.String())
	ctx.JSON(http.StatusNotFound, struct {
		Message   string `json:"message"`
		Reference string `json:"reference"`
		RemoteIP  string `json:"remote_ip"`
		ClientIP  string `json:"client_ip"`
		UserAgent string `json:"user_agent"`
		Method    string `json:"method"`
	}{
		Message:   "page not found",
		Reference: ctx.Request.URL.String(),
		RemoteIP:  ctx.RemoteIP(),
		ClientIP:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		Method:    ctx.Request.Method,
	})
}

func DefaultNoMethodFunction(ctx *gin.Context) {
	logging.Info("not allowed", "%v - from: %v, url: %v", ctx.Request.Method, ctx.ClientIP(), ctx.Request.URL.String())
	ctx.JSON(http.StatusMethodNotAllowed, struct {
		Message   string `json:"message"`
		Reference string `json:"reference"`
		RemoteIP  string `json:"remote_ip"`
		ClientIP  string `json:"client_ip"`
		UserAgent string `json:"user_agent"`
		Method    string `json:"method"`
	}{
		Message:   "method not allowed",
		Reference: ctx.Request.URL.String(),
		RemoteIP:  ctx.RemoteIP(),
		ClientIP:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		Method:    ctx.Request.Method,
	})
}

func NewEngine() IEngine {
	engine := &Engine{
		exitChan: make(chan struct{}),
		e:        gin.New(),
	}
	engine.e.NoRoute(DefaultNoRouteFunction)
	engine.e.NoMethod(DefaultNoMethodFunction)
	return engine
}
