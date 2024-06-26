package main

import (
	"fmt"
	
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/vflags"
)

type StructConf struct {
	Addr string `desc:"addr config"`
}

var (
	conf1 = vflags.String("conf1", "defaultConf1", "String vflags usage")
	conf2 = vflags.StringRequired("conf2", "requiredString usage")
	
	// this will show `structConf1` field in help
	conf3 = vflags.Struct("structConf1", &StructConf{Addr: "localhost:1"}, "struct usage with default value")
	// this will not show `structConf1` field in help
	conf4 = vflags.Struct("structConf2", (*StructConf)(nil), "struct usage with nil value")
)

func main() {
	vflags.Parse()
	
	fmt.Println(conf1())
	fmt.Println(conf2())
	
	structConf1 := &StructConf{}
	lg.PanicError(conf3(structConf1))
	fmt.Println(structConf1)
	
	structConf2 := &StructConf{}
	lg.PanicError(conf4(structConf2))
	fmt.Println(structConf2)
}
