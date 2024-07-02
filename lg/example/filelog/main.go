package main

import (
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/lg/v2/common"
)

func main() {
	fileLog := &common.LogConfig{}
	fileLog.SetDefault()

	lg.SetSlog()
	lg.EnableLogToFile(fileLog)

	for i := range 20 {
		lg.Infof("this is a log: %v", i)
	}
}
