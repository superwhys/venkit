package service

import (
	"strings"

	"github.com/superwhys/venkit/lg"
)

func WithServiceName(name string) ServiceOption {
	return func(vs *VkService) {
		lg.Debug("With service name", name)

		segs := strings.SplitN(name, ":", 2)
		if len(segs) < 2 {
			vs.serviceName = name
		} else {
			vs.serviceName = segs[0]
			vs.tag = segs[1]
		}
	}
}

func WithTag(tag string) ServiceOption {
	return func(vs *VkService) {
		vs.tag = tag
	}
}

func WithHTTPCORS() ServiceOption {
	return func(vs *VkService) {
		lg.Debug("Enabled HTTP CORS")
		vs.httpCORS = true
	}
}
