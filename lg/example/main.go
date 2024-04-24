package main

import (
	"context"

	"github.com/superwhys/venkit/lg"
)

type Data struct {
	Name string
	Age  int
}

func main() {
	lg.EnableDebug()

	lg.Info("this is log")
	lg.Debug("this is log")
	lg.Warn("this is log")
	lg.Error("this is log")

	ctx := context.Background()
	lg.Infoc(ctx, "this is log")
	lg.Debugc(ctx, "this is log")
	lg.Warnc(ctx, "this is log")
	lg.Errorc(ctx, "this is log")

	ctx = lg.With(ctx, "[test]")
	lg.Infoc(ctx, "this is log")
	lg.Debugc(ctx, "this is log")
	lg.Warnc(ctx, "this is log")
	lg.Errorc(ctx, "this is log")

	lg.Infof("this is log: %v", 1)
	lg.Debugf("this is log: %v", 1)
	lg.Warnf("this is log: %v", 1)
	lg.Errorf("this is log: %v", 1)

	data := &Data{
		Name: "hoven",
		Age:  16,
	}

	lg.Warnc(ctx, lg.Jsonify(data))
	lg.Infoc(ctx, lg.Jsonify(data))
	lg.Debugc(ctx, lg.Jsonify(data))
	lg.Errorc(ctx, lg.Jsonify(data))

	lg.Fatal("test")
}
