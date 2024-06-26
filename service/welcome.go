package service

import (
	"embed"
	"fmt"
	"net"
	"strings"
	
	"github.com/lukesampson/figlet/figletlib"
	"github.com/superwhys/venkit/lg/v2"
)

//go:embed fonts/standard.flf
var fontStandard embed.FS

func (vs *VkService) welcome(lis net.Listener) {
	lg.Infoc(vs.ctx, "Listening... Addr=%v", lis.Addr().String())
	if vs.grpcUI {
		lg.Infoc(vs.ctx, "GRPCUI enabled. URL=%s", fmt.Sprintf("http://%s/debug/grpc/ui", vs.grpcSelfConn.Target()))
	}
	
	if vs.serviceName != "" {
		vs.showServiceName()
	}
	
	if len(vs.tags) != 0 {
		lg.Infoc(vs.ctx, "Service Tag=%v", strings.Join(vs.tags, ","))
	}
	
	lg.Infoc(vs.ctx, "VenKit Service Started. Version=%v", version)
}

func (vs *VkService) showServiceName() {
	standardFont, err := fontStandard.ReadFile("fonts/standard.flf")
	if err != nil {
		return
	}
	f, err := figletlib.ReadFontFromBytes(standardFont)
	if err != nil {
		lg.Debugc(vs.ctx, "Can not show service name because of: %v", err)
		return
	}
	
	figletlib.PrintMsg(vs.serviceName, f, 80, f.Settings(), "left")
}
