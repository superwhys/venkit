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
