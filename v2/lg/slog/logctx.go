package slog

import (
	"context"
	"log/slog"
)

const (
	slContextKey = "slogContext"
)

type SlContext struct {
	childLogger *slog.Logger
}

func parseFromContext(ctx context.Context) *SlContext {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(slContextKey)
	lc, ok := val.(*SlContext)
	if !ok {
		return nil
	}
	return lc
}

func cloneSlogContext(c *SlContext) *SlContext {
	if c == nil {
		return nil
	}

	return &SlContext{
		childLogger: c.childLogger,
	}
}
