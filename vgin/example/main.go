package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vgin"
)

type BeforeRequestMiddleware struct {
}

func (m *BeforeRequestMiddleware) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	lg.Infoc(ctx, "into before request middleware")
	c.Next()
	lg.Infoc(ctx, "before request middleware done")
	return nil
}

type AfterRequestMiddleware struct {
}

func (m *AfterRequestMiddleware) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	lg.Infoc(ctx, "into after request middleware")
	return nil
}

type HelloHandler struct {
	Id          int `vpath:"user_id"`
	Name        string
	Age         int
	HeaderToken int `vheader:"Token"`
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

	engine.POST("/hello/:user_id", &BeforeRequestMiddleware{}, vgin.ParamsIn(&HelloHandler{}), &AfterRequestMiddleware{})

	engine.Run(":8080")
}
