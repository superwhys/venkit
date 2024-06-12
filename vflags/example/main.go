package main

import (
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vflags"
)

type Config struct {
	Name    string
	Address []map[string]any
}

var (
	confFlags = vflags.Struct("testConfig", &Config{}, "test config")
)

func main() {
	vflags.Parse()

	config := Config{}

	lg.PanicError(confFlags(&config))

	lg.Info(lg.Jsonify(config))
}
