package vrouter

import (
	"context"
	"net/http"
	
	"github.com/gorilla/mux"
)

const (
	statusCodeKey = "vrouter-status-key"
)

// route the final route group which will use to register into mux.Router
type route struct {
	Route
	router     *mux.Router
	middleware []Middleware
}

func (r *route) handleSpecifyMiddleware(handler HandleFunc) HandleFunc {
	next := handler
	for _, m := range r.middleware {
		next = m.WrapHandler(next)
	}
	
	return next
}

type Response interface {
	GetCode() int
	GetMessage() string
	GetJson() string
	GetData() any
	GetError() error
}

type HandleFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) Response

type Router interface {
	Routes() []Route
}

type Route interface {
	Handler() HandleFunc
	Path() string
	Method() string
}

type Middleware interface {
	WrapHandler(handler HandleFunc) HandleFunc
}

func (f HandleFunc) WrapHandler(handler HandleFunc) HandleFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) Response {
		if resp := f(ctx, w, r, vars); resp != nil && resp.GetError() != nil {
			return resp
		}
		return handler(ctx, w, r, vars)
	}
}

type defaultRoute struct {
	method  string
	path    string
	handler HandleFunc
}

func (r *defaultRoute) Handler() HandleFunc {
	return r.handler
}

func (r *defaultRoute) Path() string {
	return r.path
}

func (r *defaultRoute) Method() string {
	return r.method
}

func NewRoute(method, path string, handler HandleFunc) Route {
	return &defaultRoute{method, path, handler}
}

func NewGetRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodGet, path, handler)
}

func NewPostRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodPost, path, handler)
}

func NewPutRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodPut, path, handler)
}

func NewDeleteRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodDelete, path, handler)
}

func NewOptionsRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodOptions, path, handler)
}

func NewHeadRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodHead, path, handler)
}

func NewPatchRoute(path string, handler HandleFunc) Route {
	return NewRoute(http.MethodPatch, path, handler)
}

func NewAnyRoute(path string, handler HandleFunc) Route {
	return NewRoute("", path, handler)
}
