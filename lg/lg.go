package lg

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/superwhys/venkit/internal/shared"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *Logger
)

func init() {
	logger = New()
}

func SetDefaultLoggerOutput(stdout, stderr io.Writer) {
	logger.SetLoggerOutput(stdout, stderr)
}

func IsDebug() bool {
	return logger.enableDebug
}

func EnableDebug() {
	logger.EnableDebug()
}

func EnableLogToFile(logConf *shared.LogConfig) {
	shared.PtrLogConfig = logConf
	logger := &lumberjack.Logger{
		Filename:   logConf.FileName,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackup,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	Infof("set logger to file: %v", logConf.FileName)
	SetDefaultLoggerOutput(logger, logger)
}

func Error(v ...any) {
	logger.Error(v...)
}

func PanicError(err error, msg ...any) {
	logger.PanicError(err, msg...)
}

func Warn(v ...any) {
	logger.Warn(v...)
}

func Info(v ...any) {
	logger.Info(v...)
}

func Debug(v ...any) {
	logger.Debug(v...)
}

func Fatal(v ...any) {
	logger.Fatal(v...)
}

func Fatalf(msg string, v ...any) {
	logger.Fatalf(msg, v...)
}

func Jsonify(v any) string {
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Error(err)
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
