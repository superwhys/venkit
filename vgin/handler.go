package vgin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

func WrapHandler(ctx context.Context, handlers ...Handler) []gin.HandlerFunc {
	handlerFuncs := make([]gin.HandlerFunc, 0, len(handlers))
	for _, handler := range handlers {
		handlerFuncs = append(handlerFuncs, wrapHandler(ctx, handler))
	}

	return handlerFuncs
}

func wrapHandler(ctx context.Context, handler Handler) gin.HandlerFunc {
	ctx = lg.With(ctx, "[%v]", guessHandlerName(handler))

	var (
		handlerGetter                         = checkIsolatedHandler(handler)
		isWebsocket                           = checkIsWebsocket(handler)
		websocketUpgrader *websocket.Upgrader = nil
	)

	return func(c *gin.Context) {
		cc := &Context{c}
		var ret HandleResponse
		nh := handlerGetter()
		if !isWebsocket {
			ret = nh.HandleFunc(ctx, cc)
		} else {
			wh := nh.(WebSocketHandler)
			websocketUpgrader = initWebsocketUpgrader()
			cc.Set(webSocketKey, websocketUpgrader)
			ret = wh.HandleFunc(ctx, cc)
			if checkRet(ctx, cc, ret) {
				return
			}

			conn, _ := cc.Get(webSocketConnKey)
			ret = wh.HandleWebSocket(ctx, cc, conn.(*websocket.Conn))
		}

		if checkRet(ctx, cc, ret) {
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

func checkRet(ctx context.Context, c *Context, ret HandleResponse) (hasErr bool) {
	if ret != nil && ret.GetError() != nil {
		lg.Errorc(ctx, "handle err: %v", ret.GetError())
		AbortWithError(c.Context, ret.GetCode(), ret.GetMessage())
		return true
	}
	return false
}

func checkIsWebsocket(handler Handler) bool {
	check := func(h Handler) bool {
		switch h.(type) {
		case WebSocketHandler:
			return true
		default:
			return false
		}
	}

	switch h := handler.(type) {
	case WrapInHandler:
		return check(h.OriginHandler())
	default:
		return check(h)
	}
}

func checkIsolatedHandler(handler Handler) (handlerGetter func() Handler) {
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
	return
}

func initWebsocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
