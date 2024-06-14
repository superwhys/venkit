package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/superwhys/venkit/lg/common"
)

type Logger struct {
	*slog.Logger
	lv *slog.LevelVar
}

func NewSlogLogger() *Logger {
	return &Logger{
		Logger: slog.Default(),
	}
}

func NewSlogWithHandler(handler slog.Handler, lv *slog.LevelVar) *Logger {
	return &Logger{
		Logger: slog.New(handler),
		lv:     lv,
	}
}

func relativeToGOROOT(path string) string {
	gopath := os.Getenv("GOPATH")
	path, _ = filepath.Rel(gopath, path)
	return path
}

func getSrouce() string {
	_, file, _, _ := runtime.Caller(4)
	return relativeToGOROOT(file)
}

func NewSlogTextLogger(w ...io.Writer) *Logger {
	lv := &slog.LevelVar{}
	lv.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{
		// AddSource: true,
		Level: lv,
	}

	var writer io.Writer
	if len(w) == 0 {
		writer = os.Stdout
	} else {
		writer = w[0]
	}

	handler := slog.NewTextHandler(writer, opts)
	return NewSlogWithHandler(handler, lv)
}

func (sl *Logger) EnableDebug() {
	if sl.lv == nil {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		sl.lv.Set(slog.LevelDebug)
	}
}

func (sl *Logger) IsDebug() bool {
	return sl.Logger.Enabled(context.TODO(), slog.LevelDebug)
}

func (sl *Logger) ClearContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, slContextKey, nil)
}

func (sl *Logger) PanicError(err error, msg ...any) {
	var s string
	if err != nil {
		if len(msg) > 0 {
			s = err.Error() + ":" + fmt.Sprint(msg...)
		} else {
			s = err.Error()
		}
		sl.Error(s)
		panic(err)
	}
}

func (sl *Logger) Infof(msg string, v ...any) {
	ctx := context.TODO()
	cl := sl.currentLogger(ctx).InfoContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Debugf(msg string, v ...any) {
	ctx := context.TODO()
	cl := sl.currentLogger(ctx).DebugContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Warnf(msg string, v ...any) {
	ctx := context.TODO()
	cl := sl.currentLogger(ctx).WarnContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Errorf(msg string, v ...any) {
	ctx := context.TODO()
	cl := sl.currentLogger(ctx).ErrorContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Fatalf(msg string, v ...any) {
	ctx := context.TODO()
	cl := sl.currentLogger(ctx).ErrorContext
	sl.logc(ctx, cl, msg, v...)

	os.Exit(1)
}

func (sl *Logger) fmtMsg(keys []string, values []string, attrs []slog.Attr, v []any) []any {
	if len(keys) != len(values) {
		sl.Errorf("Invalid numbers of keys and values")
		return nil
	}
	var groupMses []any

	if len(keys) == 0 {
		// no keys, it can use as same as log/slog
		groupMses = append(groupMses, v...)
	} else {
		for i, key := range keys {
			groupMses = append(groupMses, key, values[i])
		}

		for _, attr := range attrs {
			groupMses = append(groupMses, attr)
		}
	}

	return groupMses
}

func (sl *Logger) currentLogger(ctx context.Context) *slog.Logger {
	sc := parseFromContext(ctx)
	if sc == nil || sc.childLogger == nil {
		return sl.Logger
	}

	return sc.childLogger
}

// parseKVAndAttr  parse Infoc(ctx, "this is log, addr: %v, name=%v age=%v", addr, name, age, slog.String("city", city))
// to `time=2024-06-13T21:01:46.131+08:00 level=INFO msg="this is log, addr: ..." name=aaa age=18 city=city`
// keys=[name, age], values=[aaa, 18]
func (sl *Logger) parseKVAndAttr(msg string, v ...any) (m string, keys, values []string, remains []any, attrs []slog.Attr, err error) {
	for _, v := range v {
		if a, ok := v.(slog.Attr); ok {
			attrs = append(attrs, a)
		}
	}

	m, keys, values, remains, err = common.ParseFmtKeyValue(msg, v...)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	return
}

func (sl *Logger) With(ctx context.Context, msg string, v ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(msg) == 0 && len(v) == 0 {
		return ctx
	}

	sc := parseFromContext(ctx)

	if sc == nil {
		sc = &SlContext{}
	}
	newSc := cloneSlogContext(sc)

	cl := sl.currentLogger(ctx)

	var nl *slog.Logger

	if len(v) == 0 {
		// v is empty, it means that it need to create a new group
		nl = cl.WithGroup(msg)
	} else {
		// v is not empty, this means that it need to put v into a group whose key is msg and create a child logger
		m, keys, values, remains, attrs, err := sl.parseKVAndAttr(msg, v...)
		if err != nil {
			sl.Errorf("Error parsing message: %v", err)
			return ctx
		}

		as := sl.fmtMsg(keys, values, attrs, remains)

		nl = cl.With(slog.Attr{Key: m, Value: slog.GroupValue(argsToAttrSlice(as)...)})
	}

	newSc.childLogger = nl

	return context.WithValue(ctx, slContextKey, newSc)
}

func (sl *Logger) logc(ctx context.Context, cl contextLogger, msg string, v ...any) {
	// parse the msg and v, get the common msg and key-value pairs in msg or slog.Attr
	m, keys, values, remains, attrs, err := sl.parseKVAndAttr(msg, v...)
	if err != nil {
		sl.Errorf("KV invalid: %v", err)
		return
	}

	args := sl.fmtMsg(keys, values, attrs, remains)

	sl.contextLog(cl, ctx, m, args...)
}

func (sl *Logger) Infoc(ctx context.Context, msg string, v ...any) {
	cl := sl.currentLogger(ctx).InfoContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Debugc(ctx context.Context, msg string, v ...any) {
	if !sl.IsDebug() {
		return
	}

	cl := sl.currentLogger(ctx).DebugContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Warnc(ctx context.Context, msg string, v ...any) {
	cl := sl.currentLogger(ctx).WarnContext
	sl.logc(ctx, cl, msg, v...)
}

func (sl *Logger) Errorc(ctx context.Context, msg string, v ...any) {
	cl := sl.currentLogger(ctx).ErrorContext
	sl.logc(ctx, cl, msg, v...)
}

type contextLogger func(ctx context.Context, msg string, args ...any)

func (sl *Logger) contextLog(cl contextLogger, ctx context.Context, msg string, args ...any) {
	args = append(args, slog.String("source", getSrouce()))
	cl(ctx, msg, args...)
}
