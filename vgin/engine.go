package vgin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/snail"
	"github.com/superwhys/venkit/vflags"
)

var (
	isTest = vflags.Bool("isTest", false, "whether gin mode is test")
)

func init() {
	snail.RegisterObject("ginTest", func() error {
		if isTest() {
			gin.SetMode(gin.TestMode)
		}
		return nil
	})
}

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
	if lg.IsDebug() {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

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
			ctx:         context.Background(),
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

func (g *RouterGroup) Static(relativePath, root string) gin.IRoutes {
	return g.RouterGroup.Static(relativePath, root)
}

func (g *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return g.RouterGroup.StaticFS(relativePath, fs)
}

func (g *RouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
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
	handlerName := lg.StructName(handlers.Last())
	lg.Debugc(g.ctx, "Add router --> method=%s router=%s handler=%s", method, absolutePath, handlerName)
}

func (g *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.BasePath(), relativePath)
}

func (g *RouterGroup) RegisterRouter(method, path string, handlers HandlersChain) {
	absPath := g.calculateAbsolutePath(path)
	g.Handle(method, path, wrapHandler(g.ctx, handlers...)...)
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
