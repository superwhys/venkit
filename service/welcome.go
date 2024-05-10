package service

import (
	"embed"
	"fmt"
	"net"

	"github.com/lukesampson/figlet/figletlib"
	"github.com/superwhys/venkit/lg"
)

//go:embed fonts/standard.flf
var fontStandard embed.FS

func (vs *VkService) welcome(lis net.Listener) {
	if vs.serviceName != "" {
		vs.showServiceName()
	}

	if vs.tag != "" {
		lg.Infoc(vs.ctx, "Service tag: %v", vs.tag)
	}

	lg.Infoc(vs.ctx, "Listening addr: %v", lis.Addr().String())
	if vs.grpcUI {
		lg.Infoc(vs.ctx, fmt.Sprintf("GRPCUI start in: http://%s/debug/grpc/ui", vs.grpcSelfConn.Target()))
	}

	lg.Infoc(vs.ctx, "VenKit Server Version: %v", version)
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
