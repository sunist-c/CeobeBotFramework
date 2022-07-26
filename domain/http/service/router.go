package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/CeobeBotFramework/domain/bus"
	"github.com/sunist-c/CeobeBotFramework/infrastructure/logging"
	"reflect"
	"runtime"
)

type IRouter interface {
	BaseRoute() string
	Use(middleware Middleware, middlewareInfo RequestInfo, bus bus.IScheduler, serviceName string)
	Group(url string) IRouter
	BindGET(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindPOST(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindDELETE(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindPUT(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindPATCH(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindHEAD(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindOPTIONS(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
	BindANY(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string)
}

type Router struct {
	router *gin.RouterGroup
}

func (r *Router) BaseRoute() string {
	return r.router.BasePath()
}

func (r *Router) Use(middleware Middleware, middlewareInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.Use(func(context *gin.Context) {
		middleware(NewContext(context, serviceName, middlewareInfo, bus))
	})
	logging.Info("bind middleware succeed", "%v -> %v", r.BaseRoute(), runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name())
}

func (r *Router) Group(url string) IRouter {
	g := r.router.Group(url)
	return &Router{router: g}
}

func (r *Router) BindGET(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.GET(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "GET %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindPOST(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.POST(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "POST %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindDELETE(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.DELETE(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "DELETE %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindPUT(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.PUT(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "PUT %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindPATCH(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.PATCH(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "PATCH %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindHEAD(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.HEAD(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "HEAD %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindOPTIONS(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.router.OPTIONS(url, func(context *gin.Context) {
		handler(NewContext(context, serviceName, handlerInfo, bus))
	})
	logging.Info("bind handler succeed", "OPTIONS %v -> %v", r.BaseRoute()+url, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
}

func (r *Router) BindANY(url string, handler func(ctx *Context), handlerInfo RequestInfo, bus bus.IScheduler, serviceName string) {
	r.BindGET(url, handler, handlerInfo, bus, serviceName)
	r.BindPOST(url, handler, handlerInfo, bus, serviceName)
	r.BindDELETE(url, handler, handlerInfo, bus, serviceName)
	r.BindPUT(url, handler, handlerInfo, bus, serviceName)
	r.BindPATCH(url, handler, handlerInfo, bus, serviceName)
	r.BindHEAD(url, handler, handlerInfo, bus, serviceName)
	r.BindOPTIONS(url, handler, handlerInfo, bus, serviceName)
}
