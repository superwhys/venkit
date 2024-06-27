package vrouter

import (
	"net/http"
	
	"github.com/gorilla/mux"
)

type RouterGroup struct {
	vrouter     *Vrouter
	router      *mux.Router
	routes      []route
	middlewares []Middleware
	root        bool
}

func newGroupWithRouter(router *mux.Router) RouterGroup {
	return RouterGroup{
		router:      router,
		routes:      make([]route, 0),
		middlewares: make([]Middleware, 0),
	}
}

func (rg *RouterGroup) UseMiddleware(m ...Middleware) {
	rg.middlewares = append(rg.middlewares, m...)
}

func (rg *RouterGroup) HandlerRouter(routers ...Router) {
	wrapRoutes := func(routes []Route) {
		for _, r := range routes {
			rg.vrouter.initRouter(route{
				Route:      NewRoute(r.Method(), r.Path(), r.Handler()),
				router:     rg.router,
				middleware: rg.middlewares,
			})
		}
	}
	
	for _, router := range routers {
		wrapRoutes(router.Routes())
	}
}

func (rg *RouterGroup) HandleRoute(method, path string, handler HandleFunc) {
	r := route{
		Route:  NewRoute(method, path, handler),
		router: rg.router,
	}
	if !rg.root {
		r.middleware = rg.middlewares
	}
	
	rg.vrouter.initRouter(r)
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	router := rg.router.PathPrefix(prefix).Subrouter()
	g := newGroupWithRouter(router)
	g.vrouter = rg.vrouter
	return &g
}

func (rg *RouterGroup) GET(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodGet, path, handler)
}

func (rg *RouterGroup) POST(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodPost, path, handler)
}

func (rg *RouterGroup) PUT(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodPut, path, handler)
}

func (rg *RouterGroup) PATCH(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodPatch, path, handler)
}

func (rg *RouterGroup) DELETE(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodDelete, path, handler)
}

func (rg *RouterGroup) OPTIONS(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodOptions, path, handler)
}

func (rg *RouterGroup) HEAD(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodHead, path, handler)
}

func (rg *RouterGroup) TRACE(path string, handler HandleFunc) {
	rg.HandleRoute(http.MethodTrace, path, handler)
}

func (rg *RouterGroup) Any(path string, handler HandleFunc) {
	rg.HandleRoute("", path, handler)
}
