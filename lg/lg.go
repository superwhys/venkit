package lg

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	
	"github.com/superwhys/venkit/lg/v2/log"
	"github.com/superwhys/venkit/lg/v2/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger Logger
	Ctx    context.Context
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)
	logger = log.New(log.WithCalldepth(4))
	Ctx = logger.With(context.Background(), "service", "Venkit")
}

func SetSlog() {
	logger = slog.NewSlogLogger(slog.WithCallDepth(5))
}

func SetLogger(l Logger) {
	logger = l
}

func IsDebug() bool {
	return logger.IsDebug()
}

func EnableDebug() {
	logger.EnableDebug()
}

func FileLoggerWriter(filename string, maxSize, maxBackup, maxAge int, logCompress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   logCompress,
	}
}

func salvageMsg(v ...any) (msg string, remain []any) {
	first := v[0]
	
	if s, ok := first.(string); ok {
		msg = s
	}
	
	if len(v) > 1 {
		remain = v[1:]
	}
	
	return
}

func PanicError(err error, msg ...any) {
	logger.PanicError(err, msg...)
}

func Error(v ...any) {
	if len(v) == 0 {
		return
	}
	
	msg, remain := salvageMsg(v...)
	if len(remain) == 0 {
		logger.Errorf(msg)
	} else {
		logger.Errorf(msg, remain...)
	}
}

func Warn(v ...any) {
	if len(v) == 0 {
		return
	}
	
	msg, remain := salvageMsg(v...)
	if len(remain) == 0 {
		logger.Errorf(msg)
	} else {
		logger.Errorf(msg, remain...)
	}
}

func Info(v ...any) {
	if len(v) == 0 {
		return
	}
	
	msg, remain := salvageMsg(v...)
	if len(remain) == 0 {
		logger.Infof(msg)
	} else {
		logger.Infof(msg, remain...)
	}
}

func Debug(v ...any) {
	if len(v) == 0 {
		return
	}
	
	msg, remain := salvageMsg(v...)
	if len(remain) == 0 {
		logger.Debugf(msg)
	} else {
		logger.Debugf(msg, remain...)
	}
}

func Fatal(v ...any) {
	if len(v) == 0 {
		return
	}
	
	msg, remain := salvageMsg(v...)
	if len(remain) == 0 {
		logger.Fatalf(msg)
	} else {
		logger.Fatalf(msg, remain...)
	}
}

func Fatalf(msg string, v ...any) {
	logger.Fatalf(msg, v...)
}

func Jsonify(v any) string {
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Errorf("jsonify error: %v", err)
		panic(err)
	}
	return string(d)
}

func Errorf(msg string, v ...any) {
	logger.Errorf(msg, v...)
}

func Warnf(msg string, v ...any) {
	logger.Warnf(msg, v...)
}

func Infof(msg string, v ...any) {
	logger.Infof(msg, v...)
}

func Debugf(msg string, v ...any) {
	logger.Debugf(msg, v...)
}

func ClearContext(ctx context.Context) context.Context {
	return logger.ClearContext(ctx)
}

func With(ctx context.Context, msg string, v ...any) context.Context {
	return logger.With(ctx, msg, v...)
}

func Infoc(ctx context.Context, msg string, v ...any) {
	logger.Infoc(ctx, msg, v...)
}

func Debugc(ctx context.Context, msg string, v ...any) {
	logger.Debugc(ctx, msg, v...)
}

func Warnc(ctx context.Context, msg string, v ...any) {
	logger.Warnc(ctx, msg, v...)
}

func Errorc(ctx context.Context, msg string, v ...any) {
	logger.Errorc(ctx, msg, v...)
}

// TimeFuncDuration returns the duration consumed by function.
// It has specified usage like:
//
//	    f := TimeFuncDuration()
//		   DoSomething()
//		   duration := f()
func TimeFuncDuration() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}

func TimeDurationDefer(prefix ...string) func() {
	ps := "operation"
	if len(prefix) != 0 {
		ps = strings.Join(prefix, ", ")
	}
	start := time.Now()
	
	return func() {
		Infof("%v elapsed time: %v", ps, time.Since(start))
	}
}
