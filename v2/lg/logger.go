package lg

import "context"

type Logger interface {
	EnableDebug()
	IsDebug() bool
	With(ctx context.Context, msg string, v ...any) context.Context
	ClearContext(ctx context.Context) context.Context

	PanicError(error, ...any)

	Infof(string, ...any)
	Debugf(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)

	Infoc(context.Context, string, ...any)
	Debugc(context.Context, string, ...any)
	Warnc(context.Context, string, ...any)
	Errorc(context.Context, string, ...any)
}
