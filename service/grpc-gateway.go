package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	
	gwRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/superwhys/venkit/v2/lg"
	"google.golang.org/grpc"
)

type gatewatMiddlewareHandler func(http.Handler) http.Handler
type gatewayFunc func(ctx context.Context, mux *gwRuntime.ServeMux, conn *grpc.ClientConn) error

func (vs *VkService) httpIncomingHeaderMatcher(headerName string) (mdName string, ok bool) {
	if len(vs.grpcIncomingHeaderMapping) == 0 {
		return "", false
	}
	
	key := strings.ToLower(headerName)
	mdName, exists := vs.grpcIncomingHeaderMapping[key]
	return mdName, exists
}

func (vs *VkService) httpOutgoingHeaderMatcher(headerName string) (mdName string, ok bool) {
	if len(vs.grpcOutgoingHeaderMapping) == 0 {
		return "", false
	}
	
	key := strings.ToLower(headerName)
	mdName, exists := vs.grpcOutgoingHeaderMapping[key]
	return mdName, exists
}

func (vs *VkService) mountGRPCRestfulGateway() {
	if len(vs.gatewayHandlers) <= 0 {
		return
	}
	fixGatewayVerb := func(h http.Handler, middlewares []gatewatMiddlewareHandler) http.Handler {
		for _, m := range middlewares {
			h = m(h)
		}
		
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lg.Info(fmt.Sprintf("receive request: %v", r.URL.Path))
			h.ServeHTTP(w, r)
		})
	}
	
	fn := func(ctx context.Context) error {
		opts := []gwRuntime.ServeMuxOption{
			gwRuntime.WithIncomingHeaderMatcher(vs.httpIncomingHeaderMatcher),
			gwRuntime.WithOutgoingHeaderMatcher(vs.httpOutgoingHeaderMatcher),
		}
		opts = append(opts, vs.grpcGwServeMuxOption...)
		gwmux := gwRuntime.NewServeMux(
			opts...,
		)
		
		for i := 0; i < len(vs.gatewayHandlers); i++ {
			if err := vs.gatewayHandlers[i](vs.ctx, gwmux, vs.grpcSelfConn); err != nil {
				lg.Error(fmt.Sprintf("Register %d gateway handler: %s", i, err.Error()))
				continue
			}
			vs.httpMux.Handle(vs.gatewayAPIPrefix[i]+"/", fixGatewayVerb(http.StripPrefix(vs.gatewayAPIPrefix[i], gwmux), vs.gatewayMiddlewaresHandlers[i]))
		}
		<-vs.ctx.Done()
		return nil
	}
	
	vs.mounts = append(vs.mounts, mountFn{
		baseMount: baseMount{
			fn: fn,
		},
		daemon: true,
	})
}
