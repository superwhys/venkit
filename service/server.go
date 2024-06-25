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

	"github.com/gorilla/mux"
	gwRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"github.com/superwhys/venkit/lg"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type baseMount struct {
	fn func(ctx context.Context) error
}

type mountFn struct {
	baseMount
	daemon bool
}

type cronMountFn struct {
	baseMount
	running bool
	name    string
	cron    string
	sched   cron.Schedule
}

type VkService struct {
	ctx         context.Context
	serviceName string
	tags        []string

	cmux    cmux.CMux
	httpLst net.Listener
	grpcLst net.Listener

	httpCORS    bool
	httpMux     *mux.Router
	httpHandler http.Handler

	grpcUI                bool
	grpcServer            *grpc.Server
	grpcOptions           []grpc.ServerOption
	grpcUnaryInterceptors []grpc.UnaryServerInterceptor
	grpcServersFunc       []func(*grpc.Server)
	grpcSelfConn          *grpc.ClientConn

	// grpc gateway
	grpcGwServeMuxOption       []gwRuntime.ServeMuxOption
	grpcIncomingHeaderMapping  map[string]string
	grpcOutgoingHeaderMapping  map[string]string
	gatewayAPIPrefix           []string
	gatewayHandlers            []gatewayFunc
	gatewayMiddlewaresHandlers [][]gatewatMiddlewareHandler

	workers    []worker
	mounts     []mountFn
	cronMounts []cronMountFn
}

type ServiceOption func(*VkService)

func NewVkService(opts ...ServiceOption) *VkService {
	s := &VkService{
		ctx:     lg.With(context.Background(), "Framework", "Venkit"),
		httpMux: mux.NewRouter(),
	}
	s.httpHandler = s.httpMux

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (vs *VkService) notiKill() mountFn {
	return mountFn{
		baseMount: baseMount{
			fn: func(ctx context.Context) error {
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
					lg.Infoc(vs.ctx, "Graceful stopped server successfully")

					return errors.Errorf("Signal: %s", sg.String())
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		},
		daemon: true,
	}
}

func (vs *VkService) runFinalMount() error {
	grp, ctx := errgroup.WithContext(lg.ClearContext(vs.ctx))

	// run simple worker
	for _, mount := range vs.mounts {
		mf := mount
		grp.Go(func() (err error) {
			if mf.daemon {
				err = waitContext(ctx, func() error {
					return mf.fn(ctx)
				})
			} else {
				err = mf.fn(ctx)
			}

			if err != nil && mf.daemon {
				return err
			}
			return nil
		})
	}

	// run cron worker
	if len(vs.cronMounts) != 0 {
		grp.Go(func() error {
			c := cron.New()

			for _, cw := range vs.cronMounts {
				cw := cw
				runFn := func() {
					defer lg.Debugc(ctx, "Cron Worker: %v Next scheduler time: %v", cw.name, cw.sched.Next(time.Now()))
					err := waitContext(ctx, func() error {
						if cw.running {
							return errors.New("job still running")
						}
						cw.running = true
						defer func() { cw.running = false }()

						return cw.fn(ctx)
					})
					if err != nil {
						lg.Errorc(ctx, "Run cron worker: %v error: %v", cw.name, err)
						return
					}
				}
				c.AddFunc(cw.cron, runFn)
			}

			err := waitContext(ctx, func() error {
				c.Run()
				return nil
			})

			if err != nil {
				c.Stop()
				lg.Errorc(ctx, "Run cron error: %v", err)
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
		lg.Debugc(ctx, "Worker force close after 5 seconds")
		time.Sleep(time.Second * 5)
		stop <- errors.Wrap(ctx.Err(), "Force close")
	}()

	return <-stop
}

func (vs *VkService) mountCronWorker(worker *cronWorker) cronMountFn {
	sched, err := cron.ParseStandard(worker.cron)
	if err != nil {
		lg.Fatal("cron worker cron invalid. err: %v", err)
	}

	fn := func(ctx context.Context) error {
		if worker.isWithName {
			ctx = lg.With(ctx, "Worker", worker.name)
		}

		if err := worker.fn(ctx); err != nil {
			lg.Errorc(ctx, "worker: %v run error: %v", worker.name, err)
			return errors.Wrap(err, worker.name)
		}
		return nil
	}

	return cronMountFn{
		baseMount: baseMount{
			fn: fn,
		},
		name:  worker.name,
		cron:  worker.cron,
		sched: sched,
	}
}

func (vs *VkService) mountWorker(worker *simpleWorker) mountFn {
	fn := func(ctx context.Context) error {
		var c context.Context
		if worker.daemon {
			c = ctx
		} else {
			c = context.TODO()
		}

		if worker.isWithName {
			c = lg.With(c, "Worker", worker.name)
		}

		if err := worker.fn(c); err != nil {
			lg.Errorc(c, "worker: %v run error: %v", worker.name, err)
			return errors.Wrap(err, worker.name)
		}
		return nil
	}

	return mountFn{
		baseMount: baseMount{
			fn: fn,
		},
		daemon: worker.daemon,
	}
}

func (vs *VkService) wrapWorker() {
	for _, worker := range vs.workers {
		switch w := worker.(type) {
		case *simpleWorker:
			vs.mounts = append(vs.mounts, vs.mountWorker(w))
		case *cronWorker:
			vs.cronMounts = append(vs.cronMounts, vs.mountCronWorker(w))
		default:
			lg.Warnc(vs.ctx, "Unknown worker type. worker: %v", lg.StructName(w))
		}
	}
}

func (vs *VkService) setHTTPCORS() {
	if !vs.httpCORS {
		return
	}
	vs.httpHandler = cors.AllowAll().Handler(vs.httpHandler)
}

func (vs *VkService) beginCmux(listener net.Listener) {
	vs.cmux = cmux.New(listener)
	vs.grpcLst = vs.cmux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	vs.httpLst = vs.cmux.Match(cmux.HTTP1Fast(), cmux.HTTP2())
}

func (vs *VkService) listenCmux() mountFn {
	return mountFn{
		baseMount: baseMount{
			fn: func(ctx context.Context) error {
				return vs.cmux.Serve()
			},
		},
		daemon: true,
	}
}

func (vs *VkService) loadServiceName() {
	if vs.serviceName != "" {
		return
	}
}

func (vs *VkService) serve(listener net.Listener) error {
	vs.mounts = []mountFn{
		vs.notiKill(),
	}

	if len(vs.grpcServersFunc) != 0 {
		vs.beginCmux(listener)
		vs.beginGrpc()
		vs.mounts = append(vs.mounts, vs.listenHttpServer(vs.httpLst))
		vs.mounts = append(vs.mounts, vs.listenGrpcServer(vs.grpcLst))
		vs.mounts = append(vs.mounts, vs.listenCmux())
	} else {
		vs.mounts = append(vs.mounts, vs.listenHttpServer(listener))
	}

	// grpc self connection will be used in grpcUI
	if err := vs.prepareGrpcSelfConnect(listener); err != nil {
		return errors.Wrap(err, "prepare selfConn")
	}

	vs.loadServiceName()
	vs.registerIntoConsul(listener)
	vs.mountGRPCRestfulGateway()
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
