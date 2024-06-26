package service

import (
	"context"
	"net"
	"net/http"
	"strings"
	
	"github.com/superwhys/venkit/lg/v2"
)

func (vs *VkService) listenHttpServer(lis net.Listener) mountFn {
	return mountFn{
		baseMount: baseMount{
			fn: func(ctx context.Context) error {
				return http.Serve(lis, vs.httpHandler)
			},
		},
		daemon: true,
	}
}

func WithHttpHandler(pattern string, handler http.Handler) ServiceOption {
	return func(vs *VkService) {
		if !strings.HasPrefix(pattern, "/") {
			pattern = "/" + pattern
		}
		
		defer lg.Infof("Registered http endpoint prefix. Prefix=%s", pattern)
		
		if strings.HasSuffix(pattern, "/") {
			vs.httpMux.Handle(pattern, http.StripPrefix(strings.TrimSuffix(pattern, "/"), handler))
			return
		}
		
		vs.httpMux.Handle(pattern, http.StripPrefix(pattern, handler))
	}
}
