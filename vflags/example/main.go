package main

import (
	"context"
	"time"

	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/service"
	"github.com/superwhys/venkit/v2/vflags"
)

type Config struct {
	Name    string
	Address []map[string]any
}

func (c *Config) Reload() {
	lg.Infof("config reload")
}

var (
	confFlags = vflags.Struct("testConfig", &Config{}, "test config")
)

func main() {
	vflags.Parse(
		vflags.EnableConsul(),
	)

	config := Config{}
	lg.PanicError(confFlags(&config))

	srv := service.NewVkService(
		service.WithServiceName(vflags.GetServiceName()),
		service.WithWorker(func(ctx context.Context) error {
			for {
				lg.Info(lg.Jsonify(config))

				time.Sleep(5 * time.Second)
			}
		}),
	)

	lg.PanicError(srv.Run(0))
}
