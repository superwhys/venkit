package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg/v2"
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
	return mountFn{
		baseMount: baseMount{
			fn: func(ctx context.Context) error {
				return vs.grpcServer.Serve(lis)
			},
		},
		daemon: true,
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
	
	fn := func(ctx context.Context) error {
		handler, err := standalone.HandlerViaReflection(ctx, vs.grpcSelfConn, vs.serviceName)
		if err != nil {
			return errors.Wrap(err, "start grpcUI")
		}
		
		vs.httpMux.Handle("/debug/grpc/ui/", http.StripPrefix("/debug/grpc/ui", handler))
		<-ctx.Done()
		return nil
	}
	
	vs.mounts = append(vs.mounts, mountFn{
		baseMount: baseMount{
			fn: fn,
		},
		daemon: true,
	})
}

func UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	td := lg.TimeFuncDuration()
	
	prefix := fmt.Sprintf("[%s]", strings.TrimPrefix(info.FullMethod, "/"))
	
	ctx = lg.With(ctx, prefix)
	
	ret, err := handler(ctx, req)
	duration := td()
	if err != nil {
		lg.Infoc(ctx, "Failed to handle method %s handle_time=%s handle_err=%s", info.FullMethod, duration, err)
	} else {
		lg.Infoc(ctx, "Succeed to handle method %s handle_time=%s", info.FullMethod, duration)
	}
	return ret, err
}

func StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	
	prefix := fmt.Sprintf("[%s]", strings.TrimPrefix(info.FullMethod, "/"))
	ctx = lg.With(ctx, prefix)
	td := lg.TimeFuncDuration()
	// Wrap a server stream implementation to modify the context to include the data.
	err := handler(srv, ss)
	duration := td()
	if err != nil {
		lg.Infoc(ctx, "Failed to handle stream method %s handle_time=%s handle_err=%s", info.FullMethod, duration, err)
	} else {
		lg.Infoc(ctx, "Succeed to handle stream method %s handle_time=%s", info.FullMethod, duration)
	}
	
	return err
}
