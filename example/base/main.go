package main

import (
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vflags"
)

var (
	testStr = vflags.String("testKye", "defaultVal", "test vflags string")
)

func main() {
	vflags.Parse()

	lg.Info(testStr())
}
