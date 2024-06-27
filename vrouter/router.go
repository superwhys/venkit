package vrouter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/superwhys/venkit/lg/v2"
)

const (
	DebugMode = iota
	ReleaseMode
)

var (
	vrouterMode = DebugMode
)

func SetMode(value int64) {
	switch value {
	case DebugMode:
		vrouterMode = DebugMode
	case ReleaseMode:
		vrouterMode = ReleaseMode
	default:
		panic("Vrouter mode unknown: " + strconv.FormatInt(value, 10) + " (available mode: debug release test)")
	}
}

type Vrouter struct {
	RouterGroup
	host        string
	scheme      string
	middlewares []Middleware
}

type RouterOption func(v *Vrouter)

func WithHost(host string) RouterOption {
	return func(v *Vrouter) {
		v.host = host
	}
}

func WithScheme(scheme string) RouterOption {
	return func(v *Vrouter) {
		v.scheme = scheme
	}
}

func WithNotFoundHandler(handler http.Handler) RouterOption {
	return func(v *Vrouter) {
		v.router.NotFoundHandler = handler
	}
}

func WithMethodNotAllowedHandler(handler http.Handler) RouterOption {
	return func(v *Vrouter) {
		v.router.MethodNotAllowedHandler = handler
	}
}

func (v *Vrouter) parseOptions(opts ...RouterOption) {
	for _, opt := range opts {
		opt(v)
	}
	if v.host != "" {
		v.router = v.router.Host(v.host).Subrouter()
	}
	
	if v.scheme != "" {
		v.router = v.router.Schemes(v.host).Subrouter()
	}
}

func New(opts ...RouterOption) *Vrouter {
	m := mux.NewRouter()
	
	v := &Vrouter{
		RouterGroup: newGroupWithRouter(m),
		middlewares: make([]Middleware, 0),
	}
	v.RouterGroup.root = true
	v.RouterGroup.vrouter = v
	v.parseOptions(opts...)
	
	return v
}

func NewVrouter(opts ...RouterOption) *Vrouter {
	v := New(opts...)
	v.UseMiddleware(NewLogMiddleware())
	
	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = WriteJSON(w, http.StatusNotFound, ErrorResponse(http.StatusNotFound, "page not found"))
	})
	v.router.NotFoundHandler = notFoundHandler
	v.router.MethodNotAllowedHandler = notFoundHandler
	return v
}
func (v *Vrouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.router.ServeHTTP(w, r)
}

func (v *Vrouter) ServeHandler() *mux.Router {
	return v.router
}

func (v *Vrouter) Run(addr string) error {
	srv := http.Server{
		Addr:    addr,
		Handler: v,
	}
	return srv.ListenAndServe()
}

func (v *Vrouter) initRouter(r iRoute) {
	f := v.makeHttpHandler(r)
	
	vr := r.router.Path(r.Path())
	if r.Method() != "" {
		vr = vr.Methods(r.Method())
	}
	
	if r.routeOption != nil {
		vr = r.routeOption(vr)
	}
	
	mr := vr.Handler(f)
	
	v.debugPrintRoute(r.Method(), mr, r.Handler())
}

func (v *Vrouter) UseMiddleware(m ...Middleware) {
	v.middlewares = append(v.middlewares, m...)
}

func (v *Vrouter) handleGlobalMiddleware(handler HandleFunc) HandleFunc {
	h := handler
	for _, m := range v.middlewares {
		h = m.WrapHandler(h)
	}
	
	return h
}

func (v *Vrouter) makeHttpHandler(wr iRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := lg.With(context.Background(), "handler", lg.FuncName(wr.Handler()))
		r = r.WithContext(ctx)
		
		// TODO: parse body data
		
		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}
		
		handlerFunc := v.handleGlobalMiddleware(wr.Handler())
		handlerFunc = wr.handleSpecifyMiddleware(handlerFunc)
		
		resp := handlerFunc(ctx, w, r, vars)
		if resp == nil {
			return
		}
		
		if resp.GetError() != nil {
			lg.Errorc(ctx, "handle error", "err", resp.GetError())
		}
		_ = WriteJSON(w, resp.GetCode(), resp.GetData())
	}
}

func (v *Vrouter) debugPrintRoute(method string, route *mux.Route, handler HandleFunc) {
	if vrouterMode != DebugMode {
		return
	}
	if method == "" {
		method = "ANY"
	}
	
	handlerName := lg.FuncName(handler)
	url, err := route.GetPathTemplate()
	if err != nil {
		lg.Error("get iRoute url error", "err", err, "handler", handlerName)
	}
	routerMsg := color.MagentaString(fmt.Sprintf("Method=%-6s Router=%-26s Handler=%s", method, url, handlerName))
	lg.Info(routerMsg)
}
