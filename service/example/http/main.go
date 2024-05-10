package main

import (
	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/service"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})

	srv := service.NewVkService(
		service.WithHttpHandler("", router),
	)

	lg.PanicError(srv.Run(0))
}
