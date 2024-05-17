package main

import (
	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/service"
	"github.com/superwhys/venkit/vflags"
)

func main() {
	vflags.Parse()

	router := gin.Default()

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})

	srv := service.NewVkService(
		service.WithServiceName("ServiceTest"),
		service.WithHttpHandler("", router),
	)

	lg.PanicError(srv.Run(0))
}
