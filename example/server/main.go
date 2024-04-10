package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/service"
	"github.com/superwhys/venkit/vflags"
)

func main() {
	vflags.Parse()

	srv := service.NewVkService(
		service.WithServiceName("serviceName"),
		service.WithPprof(),
		service.WithHTTPCORS(),
		service.WithWorker(func(ctx context.Context) error {
			for {
				fmt.Println(10)
				time.Sleep(time.Second)
			}
		}),
		service.WithHttpHandler("/api/", func() http.Handler {
			router := gin.Default()

			router.GET("/hello", func(ctx *gin.Context) {
				ctx.JSON(200, "helloworld")
			})

			return router
		}()),
	)

	lg.PanicError(srv.Run(28080))
}
