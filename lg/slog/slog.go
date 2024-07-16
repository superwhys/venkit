package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/superwhys/venkit/lg/v2/common"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
	callDepth  int
	withSource bool
	lv         *slog.LevelVar
}

type Opt func(*Logger)

func WithCallDepth(callDepth int) Opt {
	return func(l *Logger) {
		l.callDepth = callDepth
	}
}

func WithSource() Opt {
	return func(l *Logger) {
		l.withSource = true
	}
}

func NewSlogLogger(opts ...Opt) *Logger {
	l := &Logger{
		Logger:    slog.Default(),
		callDepth: 4,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func NewSlogWithHandler(handler slog.Handler, lv *slog.LevelVar, opts ...Opt) *Logger {
	l := &Logger{
		Logger:    slog.New(handler),
		lv:        lv,
		callDepth: 4,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l

}

func relativeToGOROOT(path string) string {
	gopath := os.Getenv("GOPATH")
	path, _ = filepath.Rel(gopath, path)
	return path
}

func (sl *Logger) getSrouce() string {
	_, file, _, _ := runtime.Caller(sl.callDepth)
	return relativeToGOROOT(file)
}

func NewSlogTextLogger(w io.Writer, opts ...Opt) *Logger {
	lv := &slog.LevelVar{}
	lv.Set(slog.LevelInfo)
	slogOpts := &slog.HandlerOptions{
		Level: lv,
	}

	if w == nil {
		w = os.Stdout
	}
	handler := slog.NewTextHandler(w, slogOpts)
	return NewSlogWithHandler(handler, lv, opts...)
}

func (sl *Logger) EnableLogToFile(logConf *common.LogConfig) {
	jackLogger := &lumberjack.Logger{
		Filename:   logConf.FileName,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackup,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	debugMode := false
	if sl.IsDebug() {
		debugMode = true
	}

	lv := &slog.LevelVar{}
	slogOpts := &slog.HandlerOptions{Level: lv}
	sl.lv = lv
	sl.Logger = slog.New(slog.NewJSONHandler(jackLogger, slogOpts))
	if debugMode {
		sl.EnableDebug()
	}

	sl.Debugf("set logger to file with json")
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
		groupMses = append(groupMses, v...)
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

	m, keys, values, remains, attrs, err := sl.parseKVAndAttr(msg, v...)
	if err != nil {
		sl.Errorf("Error parsing message: %v", err)
		return ctx
	}

	if len(v) == 0 || len(attrs)+len(remains)+len(keys) == 0 {
		nl = cl.WithGroup(m)
	} else {
		var as []any
		if len(remains)%2 == 0 {
			as = append(as, sl.fmtMsg(keys, values, attrs, remains)...)
			nl = cl.With(slog.Attr{Key: m, Value: slog.GroupValue(argsToAttrSlice(as)...)})
		} else {
			as = []any{m, remains[0]}
			as = append(as, sl.fmtMsg(keys, values, attrs, remains[1:])...)
			nl = cl.With(as...)
		}
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
	if sl.withSource {
		args = append(args, slog.String("source", sl.getSrouce()))
	}
	cl(ctx, msg, args...)
}
