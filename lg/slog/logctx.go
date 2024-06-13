package slog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/superwhys/venkit/lg/common"
)

const (
	slContextKey = "slogContext"
)

type SlContext struct {
	childLogger *slog.Logger
	msgs        []string
	keys        []string
	values      []string
	attrs       []slog.Attr
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
		msgs:        common.SliceClone(c.msgs),
		keys:        common.SliceClone(c.keys),
		values:      common.SliceClone(c.values),
	}
}

func (sc *SlContext) LogFmt() (string, []any) {
	msg := strings.Join(sc.msgs, " ")

	if len(sc.keys) != len(sc.values) {
		fmt.Println("Invalid numbers of keys and values")
		return msg, nil
	}
	var groupMses []any
	for i, key := range sc.keys {
		groupMses = append(groupMses, key, sc.values[i])
	}

	for _, attr := range sc.attrs {
		groupMses = append(groupMses, attr)
	}

	return msg, groupMses
}
