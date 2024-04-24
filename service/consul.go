package service

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/superwhys/venkit/discover"
	"github.com/superwhys/venkit/internal/shared"
	"github.com/superwhys/venkit/lg"
)

func (vs *VkService) registerIntoConsul(listener net.Listener) {
	if vs.serviceName == "" || !shared.GetIsUseConsul() {
		return
	}

	mountFn := func(ctx context.Context) error {
		addr := listener.Addr().String()
		if vs.tag == "" {
			vs.tag = "dev"
		}
		if err := discover.GetServiceFinder().RegisterServiceWithTag(vs.serviceName, addr, vs.tag); err != nil {
			lg.Errorf("register consul error: %v", err)
			return errors.Wrap(err, "Register-Consul")
		}
		lg.Infof("Register service: %v", vs.serviceName)
		if len(vs.tag) > 0 {
			lg.Infof("Registered with tags: %v", vs.tag)
		}

		<-ctx.Done()

		// programe down deregister
		discover.GetConsulServiceFinder().Close()
		return nil
	}

	vs.mounts = append(vs.mounts, mountFn)
}

func DiscoverServiceWithTag(service, tag string) string {
	return discover.GetServiceFinder().GetAddressWithTag(service, tag)
}

func DiscoverService(service string) string {
	return DiscoverServiceWithTag(service, "")
}
