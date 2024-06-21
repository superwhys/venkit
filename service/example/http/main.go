package main

import (
	"context"
	"errors"
	"time"

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
		service.WithTransientWorker("transientWorker", func(ctx context.Context) error {
			time.Sleep(time.Second * 15)
			return errors.New("return")
		}),
		service.WithCronWorker("cronWorker", "@every 5s", func(ctx context.Context) error {
			lg.Infoc(ctx, "10s run")
			time.Sleep(time.Second * 20)
			return nil
		}),
	)

	lg.PanicError(srv.Run(0))
}
