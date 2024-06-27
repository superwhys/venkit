package vrouter

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"
	
	"github.com/gorilla/mux"
)

type RouterGroup struct {
	// vrouter is the root instance
	vrouter *Vrouter
	// router is which used by the current group
	router      *mux.Router
	routes      []iRoute
	middlewares []Middleware
	root        bool
}

func newGroupWithRouter(router *mux.Router) RouterGroup {
	return RouterGroup{
		router:      router,
		routes:      make([]iRoute, 0),
		middlewares: make([]Middleware, 0),
	}
}

func (rg *RouterGroup) UseMiddleware(m ...Middleware) {
	rg.middlewares = append(rg.middlewares, m...)
}

func (rg *RouterGroup) HandlerRouter(routers ...Router) {
	wrapRoutes := func(routes []Route) {
		for _, r := range routes {
			var opt RouteOption
			switch tr := r.(type) {
			case OptRoute:
				opt = tr.Option
			default:
			}
			
			rg.vrouter.initRouter(iRoute{
				Route:       r,
				router:      rg.router,
				middleware:  rg.middlewares,
				routeOption: opt,
			})
		}
	}
	
	for _, router := range routers {
		wrapRoutes(router.Routes())
	}
}

func (rg *RouterGroup) HandleRoute(method, path string, handler HandleFunc, opt RouteOption) {
	r := iRoute{
		Route:       NewRoute(method, path, handler),
		router:      rg.router,
		routeOption: opt,
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

func (rg *RouterGroup) staticHandler(prefix string, fs http.FileSystem) HandleFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) Response {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			http.FileServer(fs).ServeHTTP(w, r2)
		} else {
			return ErrorResponse(http.StatusNotFound, "page not found")
		}
		return nil
	}
}

func (rg *RouterGroup) Static(path, root string, opt RouteOption) {
	rg.StaticFS(path, http.Dir(root), opt)
}

func (rg *RouterGroup) StaticFS(relativePath string, fs http.FileSystem, opt RouteOption) {
	urlPattern := path.Join(relativePath, "{filepath}")
	handler := rg.staticHandler(relativePath, fs)
	rg.GET(urlPattern, handler, opt)
}

func (rg *RouterGroup) GET(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodGet, path, handler, opt)
}

func (rg *RouterGroup) POST(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodPost, path, handler, opt)
}

func (rg *RouterGroup) PUT(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodPut, path, handler, opt)
}

func (rg *RouterGroup) PATCH(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodPatch, path, handler, opt)
}

func (rg *RouterGroup) DELETE(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodDelete, path, handler, opt)
}

func (rg *RouterGroup) OPTIONS(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodOptions, path, handler, opt)
}

func (rg *RouterGroup) HEAD(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodHead, path, handler, opt)
}

func (rg *RouterGroup) TRACE(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute(http.MethodTrace, path, handler, opt)
}

func (rg *RouterGroup) Any(path string, handler HandleFunc, opt RouteOption) {
	rg.HandleRoute("", path, handler, opt)
}
