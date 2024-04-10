package service

import (
	"context"

	"github.com/superwhys/venkit/lg"
)

type workerFunc func(ctx context.Context) error
type worker struct {
	name string
	fn   workerFunc
}

func WithWorker(fn workerFunc) ServiceOption {
	return func(vs *VkService) {
		lg.Debugf("Add worker: %v", lg.FuncName(fn))
		vs.workers = append(vs.workers, &worker{
			name: lg.FuncName(fn),
			fn:   fn,
		})
	}
}
