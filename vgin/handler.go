package vgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
)

type Context struct {
	*gin.Context
}

// Handler is the base Handler, you can use this when your handler not include parameters
// be sure to use IsolatedHandler when your handler include parameters
type Handler interface {
	HandleFunc(ctx context.Context, c *Context) HandleResponse
}

type NameHandler interface {
	Handler
	Name() string
}

// IsolatedHandler is used to prevent multiple concurrent use of variables with the same structure.
// InitHandler is called to create a new instance each time a request is processed
// be sure to use this when your handler include parameters
type IsolatedHandler interface {
	Handler
	InitHandler() IsolatedHandler
}

type DefaultHandler struct{}

func (dh *DefaultHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return (&Ret{}).SuccessRet("default handler")
}

func wrapDefaultHandler(ctx context.Context, handlers ...Handler) []gin.HandlerFunc {
	handlerFuncs := make([]gin.HandlerFunc, 0, len(handlers))
	for _, handler := range handlers {
		handlerFuncs = append(handlerFuncs, wrapHandler(ctx, handler))
	}

	return handlerFuncs
}

func wrapHandler(ctx context.Context, handler Handler) gin.HandlerFunc {
	ctx = lg.With(ctx, "[%v]", guessHandlerName(handler))

	var handlerGetter func() Handler
	switch h := handler.(type) {
	case IsolatedHandler:
		handlerGetter = func() Handler {
			return h.InitHandler()
		}
	default:
		handlerGetter = func() Handler {
			return handler
		}
	}

	return func(c *gin.Context) {
		ret := handlerGetter().HandleFunc(ctx, &Context{c})
		if ret != nil && ret.GetError() != nil {
			lg.Errorc(ctx, "handle err: %v", ret.GetError())
			AbortWithError(c, ret.GetCode(), ret.GetMessage())
			return
		}

		if c.IsAborted() {
			return
		}

		if ret != nil {
			ReturnWithStatus(c, ret.GetCode(), ret.GetData())
		}
	}
}

type ginHandlerFuncHandler struct {
	handlerFunc gin.HandlerFunc
}

func (h *ginHandlerFuncHandler) Name() string {
	return lg.FuncName(h.handlerFunc)
}

func (h *ginHandlerFuncHandler) InitHandler() Handler {
	return &ginHandlerFuncHandler{h.handlerFunc}
}

func (h *ginHandlerFuncHandler) HandleFunc(_ context.Context, c *Context) HandleResponse {
	h.handlerFunc(c.Context)
	return nil
}

func WrapGinHandlerFunc(handlerFunc gin.HandlerFunc) Handler {
	return &ginHandlerFuncHandler{handlerFunc}
}
