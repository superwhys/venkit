package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	
	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"github.com/superwhys/venkit/lg/v2/common"
)

type logable interface {
	Output(calldepth int, s string) error
}

type contextKey string

var (
	logContextKey contextKey = "logContext"
)

type LogContext struct {
	msg    []string
	keys   []string
	values []string
}

func ParseFromContext(ctx context.Context) *LogContext {
	if ctx == nil {
		return nil
	}
	
	val := ctx.Value(logContextKey)
	lc, ok := val.(*LogContext)
	if !ok {
		return nil
	}
	return lc
}

func cloneLogContext(c *LogContext) *LogContext {
	if c == nil {
		return nil
	}
	clone := &LogContext{
		msg:    common.SliceClone(c.msg),
		keys:   common.SliceClone(c.keys),
		values: common.SliceClone(c.values),
	}
	
	return clone
}

func (lc *LogContext) LogFmt() string {
	msg := strings.Join(lc.msg, " ")
	if len(lc.keys) != len(lc.values) {
		fmt.Println("Invalid numbers of keys and values")
		return msg
	}
	
	var buf bytes.Buffer
	
	encoder := logfmt.NewEncoder(&buf)
	
	for i := 0; i < len(lc.keys); i++ {
		encoder.EncodeKeyval(lc.keys[i], lc.values[i])
	}
	str := buf.String()
	if str == "" {
		return msg
	}
	
	return msg + " " + color.MagentaString(str)
}
