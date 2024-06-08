package main

import (
	"context"

	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/v2/vgin"
)

type User struct {
	Username string `vquery:"user_name"`
}

type jsonHandler struct {
	JsonDataStr     string  `vjson:"name" form:"name"`
	JsonDataInt     int     `vjson:"age" form:"age"`
	JsonDataFloat64 float64 `vjson:"money" form:"money"`
}

func JsonHandler(ctx context.Context, c *vgin.Context, data *jsonHandler) vgin.HandleResponse {
	return vgin.SuccessRet(data)
}

func main() {
	lg.EnableDebug()
	engine := vgin.New()

	engine.GET("hello", func(ctx context.Context, c *vgin.Context, user *User) vgin.HandleResponse {
		return vgin.SuccessRet(user.Username)
	})

	engine.POST("json", JsonHandler)

	engine.Run(":8080")
}
