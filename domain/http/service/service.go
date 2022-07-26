package service

import "github.com/sunist-c/CeobeBotFramework/domain/http"

type IService interface {
	RequestInfo() RequestInfo
	ServiceInfo() ServiceInfo
	Handler(ctx *Context)
}

type RequestInfo struct {
	UrlParams    map[string]bool
	BodyType     RequestBody
	BodyStruct   interface{}
	Headers      map[string]bool
	Dependencies map[string]bool
}

type ServiceInfo struct {
	Name        string
	Url         string
	Router      *IRouter
	Method      http.Method
	Middlewares []Middleware
}
