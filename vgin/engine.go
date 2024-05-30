package vgin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
)

var (
	anyMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

type RouterGroup struct {
	*gin.RouterGroup
	ctx context.Context
}

type Engine struct {
	*RouterGroup
	engine *gin.Engine
}

func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	if lg.IsDebug() {
		gin.SetMode(gin.DebugMode)
	}

	engine.MaxMultipartMemory = 100 << 20
	engine.Use(lg.LoggerMiddleware(), gin.Recovery())
	engine.Use(middlewares...)

	return engine
}

func New(middlewares ...gin.HandlerFunc) *Engine {
	engine := NewGinEngine()
	gin.SetMode(gin.ReleaseMode)

	return NewWithEngine(engine, middlewares...)
}

func NewWithEngine(engine *gin.Engine, middlewares ...gin.HandlerFunc) *Engine {
	engine.Use(middlewares...)
	return &Engine{
		RouterGroup: &RouterGroup{
			RouterGroup: &engine.RouterGroup,
			ctx:         lg.With(context.Background(), "[Vgin]"),
		},
		engine: engine,
	}
}

func (e *Engine) Run(addr ...string) error {
	return e.engine.Run(addr...)
}

func (e *Engine) GetGinEngine() *gin.Engine {
	return e.engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.engine.ServeHTTP(w, req)
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.engine.LoadHTMLGlob(pattern)
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	e.engine.LoadHTMLFiles(files...)
}

func (g *RouterGroup) Static(relativePath, root string) gin.IRoutes {
	return g.RouterGroup.Static(relativePath, root)
}

func (g *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return g.RouterGroup.StaticFS(relativePath, fs)
}

func (g *RouterGroup) Group(relativePath string, handlers ...Handler) *RouterGroup {
	return &RouterGroup{
		RouterGroup: g.RouterGroup.Group(relativePath, WrapHandler(g.ctx, handlers...)...),
		ctx:         g.ctx,
	}
}

func (g *RouterGroup) GroupOrigin(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		RouterGroup: g.RouterGroup.Group(relativePath, handlers...),
		ctx:         g.ctx,
	}
}

type HandlersChain []Handler

func (c HandlersChain) Last() Handler {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

func (g *RouterGroup) debugPrintRoute(method string, absolutePath string, handlers HandlersChain) {
	handler := handlers.Last()

	handlerName := guessHandlerName(handler)

	routerMsg := color.MagentaString(fmt.Sprintf("Method=%-6s Router=%-26s Handler=%s", method, absolutePath, handlerName))
	lg.Debugc(g.ctx, "Add router --> %v", routerMsg)
}

func (g *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.BasePath(), relativePath)
}

func (g *RouterGroup) RegisterRouter(method, path string, handlers HandlersChain) {
	absPath := g.calculateAbsolutePath(path)
	g.Handle(method, path, WrapHandler(g.ctx, handlers...)...)
	g.debugPrintRoute(method, absPath, handlers)
}

func (g *RouterGroup) GET(path string, handler ...Handler) {
	g.RegisterRouter(http.MethodGet, path, handler)
}

func (g *RouterGroup) POST(path string, handler ...Handler) {
	g.RegisterRouter(http.MethodPost, path, handler)
}

func (g *RouterGroup) PUT(path string, handler ...Handler) {
	g.RegisterRouter(http.MethodPut, path, handler)
}

func (g *RouterGroup) DELETE(path string, handler ...Handler) {
	g.RegisterRouter(http.MethodDelete, path, handler)
}

func (g *RouterGroup) Any(path string, handler ...Handler) {
	for _, method := range anyMethods {
		g.RegisterRouter(method, path, handler)
	}
}

func (g *RouterGroup) Specify(path string, methods []string, handler ...Handler) {
	for _, method := range methods {
		g.RegisterRouter(method, path, handler)
	}
}
