package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vgin"
)

type BeforeRequestMiddleware struct {
}

func (m *BeforeRequestMiddleware) InitHandler() vgin.Handler {
	return &BeforeRequestMiddleware{}
}

func (m *BeforeRequestMiddleware) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	lg.Infoc(ctx, "into before request middleware")
	c.Next()
	lg.Infoc(ctx, "before request middleware done")
	return nil
}

type AfterRequestMiddleware struct {
}

func (m *AfterRequestMiddleware) InitHandler() vgin.Handler {
	return &AfterRequestMiddleware{}
}

func (m *AfterRequestMiddleware) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	lg.Infoc(ctx, "into after request middleware")
	return nil
}

type HelloHandler struct {
	Id          int `vpath:"user_id"`
	Name        string
	Age         int
	Money       float64
	HeaderToken int `vheader:"Token"`
}

func (h *HelloHandler) InitHandler() vgin.IsolatedHandler {
	return &HelloHandler{}
}

func (h *HelloHandler) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	lg.Info(lg.Jsonify(h))

	ret := &vgin.Ret{
		Code: 200,
		Data: h,
	}

	return ret
}

func main() {
	lg.EnableDebug()
	engine := vgin.New()

	engine.POST("/hello/:user_id", vgin.ParamsIn(&HelloHandler{}))

	engine.Run(":8080")
}
