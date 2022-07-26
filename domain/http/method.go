package http

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	HEAD    Method = "HEAD"
	OPTIONS Method = "OPTIONS"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
)
