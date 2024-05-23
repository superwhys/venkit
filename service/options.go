package service

import (
	"strings"

	gwRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/superwhys/venkit/lg"
	"google.golang.org/grpc"
)

func WithServiceName(name string) ServiceOption {
	return func(vs *VkService) {
		segs := strings.SplitN(name, ":", 2)
		if len(segs) < 2 {
			vs.serviceName = name
		} else {
			vs.serviceName = segs[0]
			vs.tags = append(vs.tags, segs[1])
		}
	}
}

func WithTag(tag string) ServiceOption {
	return func(vs *VkService) {
		vs.tags = append(vs.tags, tag)
	}
}

func WithHTTPCORS() ServiceOption {
	return func(vs *VkService) {
		vs.httpCORS = true
	}
}

func WithGrpcServer(grpcSrv func(srv *grpc.Server)) ServiceOption {
	return func(vs *VkService) {
		vs.grpcServersFunc = append(vs.grpcServersFunc, grpcSrv)
	}
}

func WithGrpcOptions(opt grpc.ServerOption) ServiceOption {
	return func(vs *VkService) {
		vs.grpcOptions = append(vs.grpcOptions, opt)
	}
}

func WithGrpcUI() ServiceOption {
	return func(vs *VkService) {
		vs.grpcUI = true
	}
}

func WithGrpcGwServeMuxOption(options ...gwRuntime.ServeMuxOption) ServiceOption {
	return func(vs *VkService) {
		vs.grpcGwServeMuxOption = append(vs.grpcGwServeMuxOption, options...)
	}
}

func WithIncomingHeaderMatcher(mapping map[string]string) ServiceOption {
	return func(vs *VkService) {
		vs.grpcIncomingHeaderMapping = mapping
	}
}

func WithOutgoingHeaderMatcher(mapping map[string]string) ServiceOption {
	return func(vs *VkService) {
		vs.grpcOutgoingHeaderMapping = mapping
	}
}

// WithRestfulGateway binds GRPC gateway handler for all registered grpc service.
func WithRestfulGateway(apiPrefix string, handler gatewayFunc, middleware ...gatewatMiddlewareHandler) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Enabled GRPC HTTP Gateway. ApiPrefix=%s", apiPrefix)
		prefix := strings.TrimSuffix(apiPrefix, "/")
		vs.gatewayAPIPrefix = append(vs.gatewayAPIPrefix, prefix)
		vs.gatewayHandlers = append(vs.gatewayHandlers, handler)
		if len(middleware) > 0 {
			ms := make([]gatewatMiddlewareHandler, 0, len(middleware))
			// append middleware with reverse order into vs.gatewayMiddlewaresHandlers
			for i := len(middleware) - 1; i >= 0; i-- {
				ms = append(ms, middleware[i])
			}
			vs.gatewayMiddlewaresHandlers = append(vs.gatewayMiddlewaresHandlers, ms)
		} else {
			vs.gatewayMiddlewaresHandlers = append(vs.gatewayMiddlewaresHandlers, nil)
		}
	}
}

func WithGRPCUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServiceOption {
	return func(vs *VkService) {
		vs.grpcUnaryInterceptors = append(vs.grpcUnaryInterceptors, interceptors...)
	}
}
