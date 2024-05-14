# vgin
vgin is a toolkit that lightly encapsulates the gin framework.

In this toolkit, a `Handler` is defined, and all the interface logic can be write in HandleFunc
```Go
type Handler interface {
	HandleFunc(ctx context.Context, c *gin.Context) HandleResponse
}
```

If your api needs to accept some parameters, you can choose to parse the parameters into the structure yourself, 

or you can choose the `ParamsIn` provided here to automatically bind `Body`, `Form`, `Path`, `Query`, `Header` data into the `Handler`

When you use `ParamsIn` for automatic binding, you also need to upgrade the `Handler` to `IsolatedHandler` to prevent fields in the Handler from being used simultaneously in multiple threads.

## Example
```GO
type HelloHandler struct {
	Id          int `vpath:"user_id"`
	Name        string
	Age         int
	HeaderToken int `vheader:"Token"`
}

func (h *HelloHandler) InitHandler() IsolatedHandler {
	return &HelloHandler{}
}

func (h *HelloHandler) HandleFunc(ctx context.Context, c *gin.Context) vgin.HandleResponse {
	ret := &vgin.Ret{
		Code: 200,
		Data: h,
	}

	return ret
}

func main() {
	engine := vgin.New()

	engine.POST("/hello", &HelloHandler{})
	engine.POST("/hello/auto_bind/:user_id", vgin.ParamsIn(&HelloHandler{}))

	engine.Run(":8080")
}
```
