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

func WithWorker(fn workerFunc) ServiceOption {
	return WithNameWorker(lg.FuncName(fn), fn)
}

func WithNameWorker(name string, fn workerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugc(vs.ctx, "Add worker: %v", name)
		vs.workers = append(vs.workers, &worker{
			name:       name,
			fn:         fn,
			isWithName: lg.FuncName(fn) != name,
		})
	}
}
