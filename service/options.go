package service

import (
	"strings"

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
