package lg

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
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

	l := &Logger{}
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
		l.infoFlag = log.LstdFlags | log.LUTC
	}

	if l.debugFlag == 0 {
		l.debugFlag = log.LstdFlags | log.LUTC | log.Lshortfile
	}

	if l.errorFlag == 0 {
		l.errorFlag = log.LstdFlags | log.LUTC | log.Lshortfile
	}

	if l.warnFlag == 0 {
		l.warnFlag = log.LstdFlags | log.LUTC
	}

	if l.fatalFlag == 0 {
		l.fatalFlag = log.LstdFlags | log.Llongfile | log.LUTC
	}
}

func (l *Logger) EnableDebug() {
	l.enableDebug = true
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
		log.Output(4, line)
	}
}

func (l *Logger) Error(v ...any) {
	if v[0] != nil {
		l.doLog(l.errLog, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
	}
}

func (l *Logger) PanicError(err error, msg ...interface{}) {
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

func (l *Logger) Warn(v ...interface{}) {
	if v[0] != nil {
		l.doLog(l.warnLog, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if v[0] != nil {
		l.doLog(l.infoLog, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.enableDebug && v[0] != nil {
		l.doLog(l.debugLog, strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	var msg []string
	for _, i := range v {
		msg = append(msg, fmt.Sprintf("%v", i))
	}
	l.doLog(l.fatalLog, strings.Join(msg, " "))
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, v ...interface{}) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.fatalLog)
	os.Exit(1)
}

func (l *Logger) Errorf(msg string, v ...interface{}) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.errLog)
}

func (l *Logger) Warnf(msg string, v ...interface{}) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.warnLog)
}

func (l *Logger) Infof(msg string, v ...interface{}) {
	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.infoLog)
}

func (l *Logger) Debugf(msg string, v ...interface{}) {
	if !l.enableDebug {
		return
	}

	ctx := context.TODO()
	ctx = l.With(ctx, msg, v...)
	l.logc(ctx, l.debugLog)
}
