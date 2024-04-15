package lg

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	infoLog  *log.Logger
	debugLog *log.Logger
	warnLog  *log.Logger
	errLog   *log.Logger
	fatalLog *log.Logger
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

func New(options ...Option) *Logger {
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr

	l := &Logger{
		infoLog:  log.New(stdout, color.GreenString("[INFO]"), log.LstdFlags|log.LUTC),
		debugLog: log.New(stdout, color.CyanString("[DEBUG]"), log.LstdFlags|log.Lshortfile|log.LUTC),
		errLog:   log.New(stderr, color.RedString("[ERROR]"), log.LstdFlags|log.Lshortfile|log.LUTC),
		warnLog:  log.New(stdout, color.YellowString("[WARN]"), log.LstdFlags|log.LUTC),
		fatalLog: log.New(stderr, color.RedString("[FATAL]"), log.LstdFlags|log.Llongfile|log.LUTC),
	}

	for _, opt := range options {
		opt(l)
	}

	return l
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
		log.Output(3, line)
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
	if debug && v[0] != nil {
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
	var s string
	if len(v) != 0 {
		s = strings.TrimSuffix(fmt.Sprintf(msg, v...), "\n")
	} else {
		s = msg
	}
	l.doLog(l.fatalLog, strings.TrimSuffix(s, "\n"))
	os.Exit(1)
}

func (l *Logger) Errorf(msg string, v ...interface{}) {
	var s string
	if len(v) != 0 {
		s = strings.TrimSuffix(fmt.Sprintf(msg, v...), "\n")
	} else {
		s = msg
	}
	l.doLog(l.errLog, strings.TrimSuffix(s, "\n"))

}

func (l *Logger) Warnf(msg string, v ...interface{}) {
	var s string
	if len(v) != 0 {
		s = strings.TrimSuffix(fmt.Sprintf(msg, v...), "\n")
	} else {
		s = msg
	}
	l.doLog(l.warnLog, s)
}

func (l *Logger) Infof(msg string, v ...interface{}) {
	var s string
	if len(v) != 0 {
		s = strings.TrimSuffix(fmt.Sprintf(msg, v...), "\n")
	} else {
		s = msg
	}
	l.doLog(l.infoLog, s)
}

func (l *Logger) Debugf(msg string, v ...interface{}) {
	if !debug {
		return
	}

	var s string
	if len(v) != 0 {
		s = strings.TrimSuffix(fmt.Sprintf(msg, v...), "\n")
	} else {
		s = msg
	}
	l.doLog(l.debugLog, s)
}
