package main

import (
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/service"
	"github.com/superwhys/venkit/v2/service/example/grpc/examplepb"
	exampleSrv "github.com/superwhys/venkit/v2/service/example/grpc/service"
	"github.com/superwhys/venkit/v2/vflags"
	"google.golang.org/grpc"
)

func main() {
	vflags.Parse()
	
	grpcSrv := exampleSrv.NewExampleService()
	
	cs := service.NewVkService(
		service.WithServiceName(vflags.GetServiceName()),
		service.WithHTTPCORS(),
		service.WithPprof(),
		// service.WithRestfulGateway("/", examplepb.RegisterExampleHelloServiceHandler),
		service.WithGrpcUI(),
		service.WithGrpcServer(func(srv *grpc.Server) {
			examplepb.RegisterExampleHelloServiceServer(srv, grpcSrv)
		}),
	)
	
	if err := cs.Run(0); err != nil {
		lg.PanicError(err)
	}
}
