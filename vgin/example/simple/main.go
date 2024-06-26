package main

import (
	"context"
	
	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/vgin/v2"
)

type User struct {
	Username string `vquery:"user_name"`
}

type jsonHandler struct {
	JsonDataStr     string  `json:"name" form:"name"`
	JsonDataInt     int     `json:"age" form:"age"`
	JsonDataFloat64 float64 `json:"money" form:"money"`
}

func JsonHandler(ctx context.Context, c *vgin.Context, data *jsonHandler) vgin.HandleResponse {
	return vgin.SuccessRet(data)
}

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lg.Info("Middleware")
	}
}

func main() {
	lg.EnableDebug()
	engine := vgin.New()
	
	engine.GET("hello", Middleware(), func(ctx context.Context, c *vgin.Context, user *User) vgin.HandleResponse {
		return vgin.SuccessRet(user.Username)
	})
	
	engine.POST("json", JsonHandler)
	
	engine.Run(":8080")
}
