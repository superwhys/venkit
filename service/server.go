package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"github.com/superwhys/venkit/lg"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type mountFn func(ctx context.Context) error

type VkService struct {
	ctx         context.Context
	serviceName string
	tag         string

	cmux    cmux.CMux
	httpLst net.Listener
	grpcLst net.Listener

	httpCORS    bool
	httpMux     *http.ServeMux
	httpHandler http.Handler

	grpcServer      *grpc.Server
	grpcOptions     []grpc.ServerOption
	grpcServersFunc []func(*grpc.Server)
	grpcUI          bool
	grpcSelfConn    *grpc.ClientConn

	workers []*worker
	mounts  []mountFn
}

type ServiceOption func(*VkService)

func NewVkService(opts ...ServiceOption) *VkService {
	s := &VkService{
		ctx:     context.Background(),
		httpMux: http.NewServeMux(),
	}
	s.httpHandler = s.httpMux

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (vs *VkService) notiKill(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	select {
	case sg := <-ch:
		lg.Info("Graceful stopped server successfully")

		return errors.Errorf("Signal: %s", sg.String())
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (vs *VkService) runFinalMount() error {
	grp, ctx := errgroup.WithContext(vs.ctx)
	for _, mount := range vs.mounts {
		mf := mount
		grp.Go(func() error {
			err := waitContext(ctx, func() error {
				return mf(ctx)
			})
			if err != nil {
				return err
			}
			return nil
		})
	}
	return grp.Wait()
}

// waitContext used to detects whether ctx is disabled by other workers
func waitContext(ctx context.Context, fn func() error) error {
	stop := make(chan error)
	go func() {
		stop <- fn()
	}()

	go func() {
		<-ctx.Done()
		lg.Debug("Worker force close after 5 seconds")
		time.Sleep(time.Second * 5)
		stop <- errors.Wrap(ctx.Err(), "Force close")
	}()

	return <-stop
}

func (vs *VkService) mountWorker(worker *worker) mountFn {
	return func(ctx context.Context) error {
		if worker.isWithName {
			ctx = lg.With(ctx, "[%v]", worker.name)
		}

		if err := worker.fn(ctx); err != nil {
			lg.Errorf("worker: %v run error: %v", worker.name, err)
			return errors.Wrap(err, worker.name)
		}
		return nil
	}
}

func (vs *VkService) wrapWorker() {
	for _, worker := range vs.workers {
		vs.mounts = append(vs.mounts, vs.mountWorker(worker))
	}
}

func (vs *VkService) setHTTPCORS() {
	if !vs.httpCORS {
		return
	}
	vs.httpHandler = cors.AllowAll().Handler(vs.httpHandler)
}

func (vs *VkService) welcome(lis net.Listener) {
	lg.Infoc(vs.ctx, "Listening addr: %v", lis.Addr().String())
	if vs.grpcUI {
		lg.Infoc(vs.ctx, fmt.Sprintf("GRPCUI start in: http://%s/debug/grpc/ui", vs.grpcSelfConn.Target()))
	}

	lg.Infof("VenKit Server Version: %v", version)
}

func (vs *VkService) beginCmux(listener net.Listener) {
	vs.cmux = cmux.New(listener)
	vs.grpcLst = vs.cmux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	vs.httpLst = vs.cmux.Match(cmux.HTTP1Fast(), cmux.HTTP2())
}

func (vs *VkService) listenCmux(ctx context.Context) error {
	return vs.cmux.Serve()
}

func (vs *VkService) serve(listener net.Listener) error {
	vs.beginCmux(listener)
	vs.beginGrpc()

	vs.mounts = []mountFn{
		// initialize various connections
		vs.listenHttpServer(vs.httpLst),
		vs.listenGrpcServer(vs.grpcLst),
		vs.listenCmux,
		vs.notiKill,
	}

	// grpc self connection will be used in grpcUI
	if err := vs.prepareGrpcSelfConnect(listener); err != nil {
		return errors.Wrap(err, "prepare selfConn")
	}

	vs.registerIntoConsul(listener)
	vs.enableGrpcUI()
	vs.setHTTPCORS()
	vs.wrapWorker()
	vs.welcome(listener)
	return vs.runFinalMount()
}

func (vs *VkService) Run(port int) error {
	var addr string
	if port > 0 {
		addr = fmt.Sprintf(":%d", port)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return vs.serve(lis)
}
