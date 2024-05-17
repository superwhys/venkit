package service

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	GrpcTag = "grpc"
)

func (vs *VkService) beginGrpc() {
	if len(vs.grpcServersFunc) == 0 {
		return
	}

	vs.grpcServer = grpc.NewServer(vs.grpcOptions...)
	for _, fn := range vs.grpcServersFunc {
		fn(vs.grpcServer)
	}
	reflection.Register(vs.grpcServer)
}

func (vs *VkService) listenGrpcServer(lis net.Listener) mountFn {
	return func(ctx context.Context) error {
		return vs.grpcServer.Serve(lis)
	}
}

func (vs *VkService) prepareGrpcSelfConnect(listener net.Listener) error {
	if !vs.grpcUI {
		return nil
	}

	_, port, _ := net.SplitHostPort(listener.Addr().String())
	target := fmt.Sprintf("127.0.0.1:%s", port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16 * 1024 * 1024)),
	}
	conn, err := grpc.DialContext(vs.ctx, target, opts...)
	if err != nil {
		return errors.Wrap(err, "self connect to grpc")
	}
	vs.grpcSelfConn = conn
	return nil
}

func (vs *VkService) enableGrpcUI() {
	if !vs.grpcUI {
		return
	}

	mountFn := func(ctx context.Context) error {
		handler, err := standalone.HandlerViaReflection(ctx, vs.grpcSelfConn, vs.serviceName)
		if err != nil {
			return errors.Wrap(err, "start grpcUI")
		}

		vs.httpMux.Handle("/debug/grpc/ui/", http.StripPrefix("/debug/grpc/ui", handler))
		<-ctx.Done()
		return nil
	}

	vs.mounts = append(vs.mounts, mountFn)
}
