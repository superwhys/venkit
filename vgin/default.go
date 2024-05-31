package vgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
)

type DefaultHandler struct{}

func (dh *DefaultHandler) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	return (&Ret{}).SuccessRet("default handler")
}

type ginHandlerFuncHandler struct {
	handlerFunc gin.HandlerFunc
	name        string
}

func (h *ginHandlerFuncHandler) Name() string {
	if h.name != "" {
		return h.name
	}
	return lg.FuncName(h.handlerFunc)
}

func (h *ginHandlerFuncHandler) InitHandler() Handler {
	return &ginHandlerFuncHandler{
		h.handlerFunc,
		h.name,
	}
}

func (h *ginHandlerFuncHandler) HandleFunc(_ context.Context, c *Context) HandleResponse {
	h.handlerFunc(c.Context)
	return nil
}

func WrapGinHandlerFunc(handlerFunc gin.HandlerFunc) Handler {
	return &ginHandlerFuncHandler{handlerFunc: handlerFunc}
}

func WrapGinHandlerFuncWithName(name string, handlerFunc gin.HandlerFunc) Handler {
	return &ginHandlerFuncHandler{handlerFunc: handlerFunc, name: name}
}
