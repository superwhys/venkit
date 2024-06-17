package main

import (
	"context"

	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/lg/log"
	"github.com/superwhys/venkit/lg/slog"
)

type Data struct {
	Name string
	Age  int
}

func init() {
}

func main() {
	lg.Error("test error")

	logLogger := log.New()
	logLogger.EnableDebug()

	ctx := context.Background()
	lg.Infoc(ctx, "this is log: %v, name: %v, age=%v", "protocol", "super", 16, "protocol", 27)
	lg.Infoc(ctx, "this is log: %v", "info")
	lg.Infoc(ctx, "this is log: %v", "info", "badValue")
	lg.Infoc(ctx, "this is log: %v", "info", "protocol", 27)
	logLogger.Infoc(ctx, "this is log")
	logLogger.Debugc(ctx, "this is log")
	logLogger.Warnc(ctx, "this is log")
	logLogger.Errorc(ctx, "this is log")

	ctx = logLogger.With(ctx, "[test] prefix=%s", "logLogger")
	logLogger.Infoc(ctx, "this is log")
	logLogger.Debugc(ctx, "this is log")
	logLogger.Warnc(ctx, "this is log")
	logLogger.Errorc(ctx, "this is log", "name", "super")

	logLogger.Infof("this is log: %v", 1)
	logLogger.Debugf("this is log: %v", 1)
	logLogger.Warnf("this is log: %v", 1)
	logLogger.Errorf("this is log: %v", 1)

	data := &Data{
		Name: "hoven",
		Age:  16,
	}

	logLogger.Warnc(ctx, lg.Jsonify(data))
	logLogger.Infoc(ctx, lg.Jsonify(data))
	logLogger.Debugc(ctx, lg.Jsonify(data))
	logLogger.Errorc(ctx, lg.Jsonify(data))

	slogLogger := slog.NewSlogLogger()
	slogLogger.EnableDebug()
	slogLogger.Infof("this is slog: %s", "info", "city", "shenzhen")
	slogLogger.Warnf("this is slog: %v", "warn")
	slogLogger.Errorf("this is slog: %v", "error")
	slogLogger.Debugf("this is slog: %v", "debug")

	ctx = slogLogger.With(context.Background(), "[test] prefix=%s city=%s", "slogLogger", "shenzhen", "ani", "dog")

	slogLogger.Infoc(ctx, "this is slog context: %v, name=%s age=%d", "info", "super", 18)
	slogLogger.Debugc(ctx, "this is slog context: %v, name=%s age=%d", "debug", "super", 18)
	slogLogger.Errorc(ctx, "this is slog context: %v, name=%s age=%d", "error", "super", 18)
	slogLogger.Warnc(ctx, "this is slog context: %v, name=%s age=%d", "warn", "super", 18)

	ctx = slogLogger.With(ctx, "Group")

	slogLogger.Infoc(ctx, "this is group msg, day: %v", 1, "province", "guangdond")

	//logLogger.Fatal("test")
}
