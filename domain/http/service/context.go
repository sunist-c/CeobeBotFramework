package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/CeobeBotFramework/domain/bus"
	"github.com/sunist-c/CeobeBotFramework/infrastructure/authenticator"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type RequestBody string

const (
	JsonBody RequestBody = "json"
	RawBody  RequestBody = "byte"
	YamlBody RequestBody = "yaml"
	XmlBody  RequestBody = "xml"
	FormBody RequestBody = "form"
)

type Context struct {
	ctx          *gin.Context
	ServiceName  string
	Errors       []error
	UrlParams    map[string]string
	Headers      map[string]string
	Body         any
	Dependencies map[string]any
}

func (c *Context) SaveUploadFiles(file *multipart.FileHeader, dst string) (err error) {
	return c.ctx.SaveUploadedFile(file, dst)
}

func (c *Context) Stream(step func(w io.Writer) bool) bool {
	return c.ctx.Stream(step)
}

func (c *Context) GetCookie(key string) (string, error) {
	return c.ctx.Cookie(key)
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *Context) GetContentType() string {
	return c.ctx.ContentType()
}

func (c *Context) SetAccepted(args ...string) {
	c.ctx.SetAccepted(args...)
}

func (c *Context) SetSameSite(site http.SameSite) {
	c.ctx.SetSameSite(site)
}

func (c *Context) JsonResponse(code int, body any) {
	c.ctx.JSON(code, body)
}

func (c *Context) JsonpResponse(code int, body any) {
	c.ctx.JSONP(code, body)
}

func (c *Context) ProtoBufResponse(code int, body any) {
	c.ctx.ProtoBuf(code, body)
}

func (c *Context) PureJsonResponse(code int, body any) {
	c.ctx.PureJSON(code, body)
}

func (c *Context) XmlResponse(code int, body any) {
	c.ctx.XML(code, body)
}

func (c *Context) StringResponse(code int, format string, args ...any) {
	c.ctx.String(code, format, args...)
}

func (c *Context) RawResponse(code int, contentType string, raw []byte) {
	c.ctx.Data(code, contentType, raw)
}

func (c *Context) FileAttachmentResponse(filename, filepath string) {
	c.ctx.FileAttachment(filepath, filename)
}

func (c *Context) FileResponse(filepath string) {
	c.ctx.File(filepath)
}

func (c *Context) FileSystemResponse(filepath string, system http.FileSystem) {
	c.ctx.FileFromFS(filepath, system)
}

func (c *Context) Redirect(code int, dst string) {
	c.ctx.Redirect(code, dst)
}

func NewContext(ctx *gin.Context, serviceName string, info RequestInfo, bus bus.IScheduler) *Context {
	context := &Context{
		ctx:          ctx,
		ServiceName:  serviceName,
		Errors:       []error{},
		UrlParams:    map[string]string{},
		Headers:      map[string]string{},
		Body:         nil,
		Dependencies: map[string]interface{}{},
	}

	// query params
	for key, necessary := range info.UrlParams {
		value, status := "", false
		if strings.HasPrefix(key, ":") {
			value, status = context.ctx.Params.Get(key)
		} else {
			value, status = context.ctx.GetQuery(key)
		}

		if necessary && !status {
			context.Errors = append(context.Errors, NecessaryFieldNotFoundError{
				field: key,
				where: "url-params",
			})
		} else {
			context.UrlParams[key] = value
		}
	}

	// load headers
	for key, necessary := range info.Headers {
		value := context.ctx.GetHeader(key)
		if necessary && value == "" {
			context.Errors = append(context.Errors, NecessaryFieldNotFoundError{
				field: key,
				where: "headers",
			})
		} else {
			context.Headers[key] = value
		}
	}

	// bind body
	switch info.BodyType {
	case JsonBody:
		context.Body = info.BodyStruct
		if err := context.ctx.ShouldBindJSON(&context.Body); err != nil {
			context.Errors = append(context.Errors, UnMarshalError{
				mapTo: context.Body,
				info:  err.Error(),
			})
		}
	case XmlBody:
		context.Body = info.BodyStruct
		if err := context.ctx.ShouldBindXML(&context.Body); err != nil {
			context.Errors = append(context.Errors, UnMarshalError{
				mapTo: context.Body,
				info:  err.Error(),
			})
		}
	case YamlBody:
		context.Body = info.BodyStruct
		if err := context.ctx.ShouldBindYAML(&context.Body); err != nil {
			context.Errors = append(context.Errors, UnMarshalError{
				mapTo: context.Body,
				info:  err.Error(),
			})
		}
	case FormBody:
		if err := context.ctx.Request.ParseForm(); err != nil {
			context.Errors = append(context.Errors, UnMarshalError{
				mapTo: url.Values{},
				info:  err.Error(),
			})
		} else {
			context.Body = context.ctx.Request.PostForm
		}
	case RawBody:
		if buff, err := ioutil.ReadAll(context.ctx.Request.Body); err != nil {
			context.Errors = append(context.Errors, UnMarshalError{
				mapTo: []byte{},
				info:  err.Error(),
			})
		} else {
			context.Body = buff
		}
	default:
		context.Errors = append(context.Errors, UnknownBodyTypeError{t: string(info.BodyType)})
	}

	// get dependencies
	for key, necessary := range info.Dependencies {
		result, err := bus.GetComponent(key, authenticator.NewUser(serviceName))
		if err != nil && necessary {
			context.Errors = append(context.Errors, err)
			context.Errors = append(context.Errors, NecessaryFieldNotFoundError{
				field: key,
				where: "dependencies",
			})
		} else {
			context.Dependencies[key] = result
		}
	}

	return context
}
