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
	"github.com/superwhys/venkit/internal/shared"
	"github.com/superwhys/venkit/lg"
	"golang.org/x/sync/errgroup"
)

type mountFn func(ctx context.Context) error

type VkService struct {
	ctx         context.Context
	serviceName string
	tag         string

	httpCORS bool

	workers     []*worker
	httpMux     *http.ServeMux
	httpHandler http.Handler
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

func (vs *VkService) runFinalMount(mounts []mountFn) error {
	grp, ctx := errgroup.WithContext(vs.ctx)
	for _, mount := range mounts {
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
		if err := worker.fn(ctx); err != nil {
			lg.Errorf("worker: %v run error: %v", worker.name, err)
			return errors.Wrap(err, worker.name)
		}
		return nil
	}
}

func (vs *VkService) wrapWorker() []mountFn {
	mounts := make([]mountFn, 0, len(vs.workers))

	for _, worker := range vs.workers {
		mounts = append(mounts, vs.mountWorker(worker))
	}

	return mounts
}

func (vs *VkService) setHTTPCORS() error {
	vs.httpHandler = cors.AllowAll().Handler(vs.httpHandler)
	return nil
}

func (vs *VkService) welcome(lis net.Listener) {
	lg.Infof("Listening addr: %v", lis.Addr().String())
	lg.Infof("VenKit Server Version: %v", version)
}

func (vs *VkService) serve(listener net.Listener) error {

	mounts := []mountFn{
		vs.listenHttpServer(listener),
		vs.notiKill,
	}

	if vs.serviceName != "" && shared.GetIsUseConsul() {
		mounts = append(mounts, vs.registerIntoConsul(listener))
	}

	if vs.httpCORS {
		vs.setHTTPCORS()
	}
	mounts = append(mounts, vs.wrapWorker()...)

	vs.welcome(listener)
	return vs.runFinalMount(mounts)
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
