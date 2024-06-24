package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/superwhys/venkit/internal/shared"
	"github.com/superwhys/venkit/lg/common"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	callDepth   int
	enableDebug bool
	infoLog     *log.Logger
	debugLog    *log.Logger
	warnLog     *log.Logger
	errLog      *log.Logger
	fatalLog    *log.Logger

	infoFlag  int
	debugFlag int
	warnFlag  int
	errorFlag int
	fatalFlag int
}

type Option func(*Logger)

func WithCalldepth(callDepth int) Option {
	return func(l *Logger) {
		l.callDepth = callDepth
	}
}

func WithFileOption(filename string, maxSize, maxBackup, maxAge int, logCompress bool) Option {
	return func(l *Logger) {
		logger := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackup,
			MaxAge:     maxAge,
			Compress:   logCompress,
		}
		l.SetLoggerOutput(logger, logger)
	}
}

func WithInfoFlag(flag int) Option {
	return func(l *Logger) {
		l.infoFlag = flag
	}
}

func WithDebugFlag(flag int) Option {
	return func(l *Logger) {
		l.debugFlag = flag
	}
}

func WithWarnFlag(flag int) Option {
	return func(l *Logger) {
		l.warnFlag = flag
	}
}

func WithErrorFlag(flag int) Option {
	return func(l *Logger) {
		l.errorFlag = flag
	}
}

func New(options ...Option) *Logger {
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr

	l := &Logger{
		callDepth: 3,
	}
	for _, opt := range options {
		opt(l)
	}
	l.defaultFlag()

	l.infoLog = log.New(stdout, colorString(color.FgCyan, formatPrefix("[INFO]")), l.infoFlag)
	l.debugLog = log.New(stdout, colorString(color.FgWhite, formatPrefix("[DEBUG]")), l.debugFlag)
	l.errLog = log.New(stderr, colorString(color.FgRed, formatPrefix("[ERROR]")), l.errorFlag)
	l.warnLog = log.New(stdout, colorString(color.FgYellow, formatPrefix("[WARN]")), l.warnFlag)
	l.fatalLog = log.New(stderr, colorString(color.FgRed, formatPrefix("[FATAL]")), l.fatalFlag)

	return l
}

func colorString(col color.Attribute, str string) string {
	return color.New(col).Add(color.Bold).Sprintf(str)
}

func formatPrefix(prefix string) string {
	return fmt.Sprintf("%-8s", prefix)
}

func (l *Logger) defaultFlag() {
	if l.infoFlag == 0 {
		l.infoFlag = log.LstdFlags | log.Ldate | log.Ltime
	}

	if l.debugFlag == 0 {
		l.debugFlag = log.LstdFlags | log.Ldate | log.Ltime | log.Lshortfile
	}

	if l.errorFlag == 0 {
		l.errorFlag = log.LstdFlags | log.Ldate | log.Ltime | log.Lshortfile
	}

	if l.warnFlag == 0 {
		l.warnFlag = log.LstdFlags | log.Ldate | log.Ltime
	}

	if l.fatalFlag == 0 {
		l.fatalFlag = log.LstdFlags | log.Llongfile | log.Ldate | log.Ltime
	}
}

func (l *Logger) EnableDebug() {
	l.enableDebug = true
}

func (l *Logger) IsDebug() bool {
	return l.enableDebug
}

func (l *Logger) SetDefaultLoggerOutput(stdout, stderr io.Writer) {
	l.SetLoggerOutput(stdout, stderr)
}

func (l *Logger) EnableLogToFile(logConf *shared.LogConfig) {
	shared.PtrLogConfig = logConf
	jackLogger := &lumberjack.Logger{
		Filename:   logConf.FileName,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackup,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	l.Infof("set logger to file: %v", logConf.FileName)
	l.SetDefaultLoggerOutput(jackLogger, jackLogger)
}

func (l *Logger) SetLoggerOutput(stdout, stderr io.Writer) {
	l.infoLog.SetOutput(stdout)
	l.debugLog.SetOutput(stdout)
	l.errLog.SetOutput(stderr)
	l.warnLog.SetOutput(stdout)
	l.fatalLog.SetOutput(stderr)
}

func (l *Logger) doLog(log *log.Logger, msg string) {
	for _, line := range strings.Split(msg, "\n") {
		log.Output(l.callDepth, line)
	}
}

func (l *Logger) ClearContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, nil)
}

func (l *Logger) PanicError(err error, msg ...any) {
	var s string
	if err != nil {
		if len(msg) > 0 {
			s = err.Error() + ":" + fmt.Sprint(msg...)
		} else {
			s = err.Error()
		}
		l.doLog(l.errLog, s)
		panic(err)
	}
}

func (l *Logger) Fatalf(msg string, v ...any) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.fatalLog)
	os.Exit(1)
}

func (l *Logger) Errorf(msg string, v ...any) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.errLog)
}

func (l *Logger) Warnf(msg string, v ...any) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.warnLog)
}

func (l *Logger) Infof(msg string, v ...any) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.infoLog)
}

func (l *Logger) Debugf(msg string, v ...any) {
	if !l.enableDebug {
		return
	}

	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.debugLog)
}

func (l *Logger) logc(ctx context.Context, lb logable) {
	lc := ParseFromContext(ctx)
	if lc == nil {
		return
	}

	msg := lc.LogFmt()
	for _, line := range strings.Split(msg, "\n") {
		lb.Output(l.callDepth, line)
	}
}

const badKey = "!BADKEY"

func getKVParis(kvs []any) (string, string, []any) {
	switch x := kvs[0].(type) {
	case string:
		if len(kvs) == 1 {
			return badKey, x, nil
		}
		return x, fmt.Sprintf("%v", kvs[1]), kvs[2:]
	default:
		return badKey, fmt.Sprintf("%v", x), kvs[1:]
	}
}

func (l *Logger) With(ctx context.Context, msg string, v ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(msg) == 0 && len(v) == 0 {
		return ctx
	}

	lc := ParseFromContext(ctx)
	if lc == nil {
		lc = &LogContext{}
	}

	newLc := cloneLogContext(lc)

	msg, keys, values, remains, err := common.ParseFmtKeyValue(msg, v...)
	if err != nil {
		l.Errorf("%v", err.Error())
		return ctx
	}

	remainsParser := func(remains []any) {
		var key, val string
		for len(remains) > 0 {
			key, val, remains = getKVParis(remains)
			keys = append(keys, key)
			values = append(values, val)
		}
	}

	var msgPrefix bool

	// l.With(ctx, "prefix", "logPrefix")
	// output: "[INFO] this is a log prefix=logPrefix"
	if len(remains) == len(v) {
		// This means that msg does not have any formatting symbols
		if len(v) == 1 {
			// This means it use msg for key and v[0] for value
			keys = append(keys, msg)
			values = append(values, fmt.Sprintf("%v", v[0]))
		} else if len(v)%2 == 0 {
			// This means that v is made up of key-value pairs
			msgPrefix = true
			remainsParser(remains)
		} else {
			keys = append(keys, msg)
			values = append(values, fmt.Sprintf("%v", v[0]))

			remainsParser(remains[1:])
		}
	} else {
		// This means that msg has some formatting symbols like `%v`
		// We just need to worry about the remaining key-value pairs in remains
		msgPrefix = true
		remainsParser(remains)
	}

	if msg != "" && msgPrefix {
		newLc.msg = append(newLc.msg, msg)
	}

	newLc.keys = append(newLc.keys, keys...)
	newLc.values = append(newLc.values, values...)
	return context.WithValue(ctx, logContextKey, newLc)
}

func (l *Logger) Infoc(ctx context.Context, msg string, v ...any) {
	if len(msg) > 0 || len(v) > 0 {
		ctx = l.With(ctx, msg, v...)
	}
	l.logc(ctx, l.infoLog)
}

func (l *Logger) Debugc(ctx context.Context, msg string, v ...any) {
	if !l.enableDebug {
		return
	}

	if len(msg) > 0 || len(v) > 0 {
		ctx = l.With(ctx, msg, v...)
	}
	l.logc(ctx, l.debugLog)
}

func (l *Logger) Errorc(ctx context.Context, msg string, v ...any) {
	if len(msg) > 0 || len(v) > 0 {
		ctx = l.With(ctx, msg, v...)
	}
	l.logc(ctx, l.errLog)
}

func (l *Logger) Warnc(ctx context.Context, msg string, v ...any) {
	if len(msg) > 0 || len(v) > 0 {
		ctx = l.With(ctx, msg, v...)
	}
	l.logc(ctx, l.warnLog)
}
