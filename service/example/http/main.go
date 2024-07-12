package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/service"
	"github.com/superwhys/venkit/v2/vflags"
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
			t := time.NewTicker(time.Second * 2)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					lg.Errorf("ctx done")
					return ctx.Err()
				case <-t.C:
					lg.Infoc(ctx, "2 second run")
				}
			}
		}),
		service.WithCronWorker("cronWorker", "@every 5s", func(ctx context.Context) error {
			lg.Infoc(ctx, "cron run")
			return nil
		}),
	)

	lg.PanicError(srv.Run(0))
}
