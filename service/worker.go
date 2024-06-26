package service

import (
	"context"
	
	"github.com/superwhys/venkit/lg/v2"
)

type worker interface {
	Fn(ctx context.Context) error
}

type base struct {
	name       string
	isWithName bool
	fn         workerFunc
}

type simpleWorker struct {
	*base
	daemon bool
}

func (s *simpleWorker) Fn(ctx context.Context) error {
	return s.fn(ctx)
}

type cronWorker struct {
	*base
	cron string
}

func (s *cronWorker) Fn(ctx context.Context) error {
	return s.fn(ctx)
}

type workerFunc func(ctx context.Context) error
type WorkerFunc workerFunc

func WithWorker(fn WorkerFunc) ServiceOption {
	return WithNameWorker(lg.FuncName(fn), fn)
}

func WithNameWorker(name string, fn WorkerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Add Daemon worker. WorkerName=%v", name)
		vs.workers = append(vs.workers, &simpleWorker{
			base: &base{
				name:       name,
				fn:         workerFunc(fn),
				isWithName: lg.FuncName(fn) != name,
			},
			daemon: true,
		})
	}
}

func WithTransientWorker(name string, fn WorkerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Add Transient worker. WorkerName=%v", name)
		vs.workers = append(vs.workers, &simpleWorker{
			base: &base{
				name:       name,
				fn:         workerFunc(fn),
				isWithName: lg.FuncName(fn) != name,
			},
			daemon: false,
		})
	}
}

func WithCronWorker(name, cron string, fn WorkerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Add Cron worker. WorkerName=%v Cron=%v", name, cron)
		vs.workers = append(vs.workers, &cronWorker{
			base: &base{
				name:       name,
				fn:         workerFunc(fn),
				isWithName: lg.FuncName(fn) != name,
			},
			cron: cron,
		})
	}
}
