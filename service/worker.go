package service

import (
	"context"

	"github.com/superwhys/venkit/lg"
)

type workerFunc func(ctx context.Context) error
type worker struct {
	name       string
	isWithName bool
	fn         workerFunc
}

type WorkerFunc workerFunc

func WithWorker(fn WorkerFunc) ServiceOption {
	return WithNameWorker(lg.FuncName(fn), fn)
}

func WithNameWorker(name string, fn WorkerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Add worker. WorkerName=%v", name)
		vs.workers = append(vs.workers, &worker{
			name:       name,
			fn:         workerFunc(fn),
			isWithName: lg.FuncName(fn) != name,
		})
	}
}
