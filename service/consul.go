package service

import (
	"context"
	"fmt"
	"net"
	"strings"
	
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/v2/discover"
	"github.com/superwhys/venkit/v2/internal/shared"
	"github.com/superwhys/venkit/v2/lg"
)

func (vs *VkService) registerIntoConsul(listener net.Listener) {
	if vs.serviceName == "" || !shared.GetIsUseConsul() {
		return
	}
	
	fn := func(ctx context.Context) error {
		addr := listener.Addr().String()
		if len(vs.tags) == 0 {
			vs.tags = append(vs.tags, "dev")
		}
		
		if len(vs.grpcServersFunc) != 0 {
			vs.tags = append(vs.tags, GrpcTag)
		}
		
		if err := discover.GetServiceFinder().RegisterServiceWithTags(vs.serviceName, addr, vs.tags); err != nil {
			lg.Errorf("register consul error: %v", err)
			return errors.Wrap(err, "Register-Consul")
		}
		
		var logArgs []any
		logText := "Registered into consul success. Service=%v"
		logArgs = append(logArgs, vs.serviceName)
		if len(vs.tags) > 0 {
			logText = fmt.Sprintf("%v %v", logText, "Tag=%v")
			logArgs = append(logArgs, strings.Join(vs.tags, ","))
		}
		
		lg.Infoc(vs.ctx, logText, logArgs...)
		
		<-ctx.Done()
		
		// programe down deregister
		discover.GetServiceFinder().Close()
		
		return nil
	}
	
	vs.mounts = append(vs.mounts, mountFn{
		baseMount: baseMount{
			fn: fn,
		},
		daemon: true,
	})
}

func DiscoverServiceWithTag(service, tag string) string {
	return discover.GetServiceFinder().GetAddressWithTag(service, tag)
}

func DiscoverService(service string) string {
	return DiscoverServiceWithTag(service, "")
}
