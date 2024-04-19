# vgin
vgin is a toolkit that lightly encapsulates the gin framework.

In this toolkit, a `Handler` is defined, and all parameter binding is done automatically in the handler
```Go
type Handler interface {
	HandleFunc(ctx context.Context, c *gin.Context) HandleResponse
}
```

## Example
```GO
type HelloHandler struct {
	Id          int `vpath:"user_id"`
	Name        string
	Age         int
	HeaderToken int `vheader:"Token"`
}

func (h *HelloHandler) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	ret := &vgin.Ret{
		Code: 200,
		Data: h,
	}

	return ret
}

func main() {
	lg.EnableDebug()
	engine := vgin.New()

	engine.POST("/hello/:user_id", &HelloHandler{})
	engine.POST("/hello/auto_bind/:user_id", vgin.ParamsIn(&HelloHandler{}))

	engine.Run(":8080")
}
```

You can use `vgin.ParamsIn` handelr to let `vgin` auto bind the `Body`, `Form`, `Path`, `Query` data into the handler

