package log

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	LevelInfo = iota
	LevelWarn
	LevelError
	LevelDebug

	PrefixError = "\033[31m[ERROR]\033[0m \u001B[34m"
	PrefixWarn  = "\033[33m[WARN]\033[0m \u001B[34m"
	PrefixInfo  = "\033[32m[INFO]\033[0m \u001B[34m"
	PrefixDebug = "\033[36m[DEBUG]\033[0m \u001B[34m"
)

var (
	prefixs      = []string{PrefixInfo, PrefixWarn, PrefixError, PrefixDebug}
	globalLogger = NewLogger(LevelDebug, os.Stdout)
)

type Logger struct {
	loggers []*log.Logger
}

func NewLogger(level int, out io.Writer) *Logger {
	if level < 0 {
		panic(errors.New("invalid log level"))
	}
	if level > LevelDebug {
		level = LevelDebug
	}
	l := new(Logger)
	l.loggers = make([]*log.Logger, LevelDebug+1)
	i := 0
	for ; i <= level; i++ {
		if i == LevelInfo {
			l.loggers[i] = log.New(out, prefixs[i], log.LstdFlags)
		} else {
			l.loggers[i] = log.New(out, prefixs[i], log.LstdFlags)
		}
	}
	for ; i <= LevelDebug; i++ {
		l.loggers[i] = log.New(ioutil.Discard, "", 0)
	}
	return l
}

func caller() string {
	_, file, line, _ := runtime.Caller(2)
	index := strings.LastIndex(file, "/")
	return fmt.Sprintf("%s:%d: ", file[index+1:], line)
}

func (l *Logger) output(level int, caller string, format string, args ...interface{}) {
	l.loggers[level].Printf(caller+"\033[0m"+format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.output(LevelInfo, caller(), format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.output(LevelWarn, caller(), format, args...)
}

func (l *Logger) Error(err error) {
	l.output(LevelError, caller(), err.Error())
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.output(LevelError, caller(), format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.output(LevelDebug, caller(), format, args...)
}

func (l *Logger) SetOutput(out io.Writer) {
	for _, logger := range l.loggers {
		logger.SetOutput(out)
	}
}

func Info(format string, args ...interface{}) {
	globalLogger.output(LevelInfo, caller(), format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.output(LevelWarn, caller(), format, args...)
}

func Error(err error) {
	globalLogger.output(LevelError, caller(), err.Error())
}

func Errorf(format string, args ...interface{}) {
	globalLogger.output(LevelError, caller(), format, args...)
}

func Debug(format string, args ...interface{}) {
	globalLogger.output(LevelDebug, caller(), format, args...)
}

func SetOutput(out io.Writer) {
	globalLogger.SetOutput(out)
}

func (l *Logger) SetLevel(level int) {
	if level < 0 {
		panic(errors.New("invalid log level"))
	}
	if level > LevelDebug {
		level = LevelDebug
	}
	for i := 0; i <= level; i++ {
		l.loggers[i].SetOutput(os.Stdout)
	}
	for i := level + 1; i <= LevelDebug; i++ {
		l.loggers[i].SetOutput(ioutil.Discard)
	}
}

func SetLevel(level int) {
	globalLogger.SetLevel(level)
}
