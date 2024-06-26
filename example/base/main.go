package main

import (
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/vflags"
)

var (
	testStr = vflags.String("testKye", "defaultVal", "test vflags string")
)

func main() {
	vflags.Parse()
	
	lg.Info(testStr())
}
