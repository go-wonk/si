package sirabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
)

type LogLevel int8

const (
	DebugLevel LogLevel = -1 + iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	IgnoreLevel
)

var (
	logLevel   = DebugLevel
	defaultLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	_logger    Logger
	// debugLogger *log.Logger
	// infoLogger  *log.Logger
	// warnLogger  *log.Logger
	// errorLogger *log.Logger
)

func init() {
	_logger = NewDefaultLogger(defaultLog)
	// debugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	// infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	// warnLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	// errorLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func SetLevel(lvl LogLevel) {
	logLevel = lvl
}

func SetLogger(logger Logger) {
	_logger = logger
}

func Debug(msg string) {
	_logger.Debug(context.TODO(), msg)
}
func Debugf(msg string, v ...any) {
	_logger.Debug(context.TODO(), fmt.Sprintf(msg, v...))
}

func Info(msg string) {
	_logger.Info(context.TODO(), msg)
}
func Infof(msg string, v ...any) {
	_logger.Info(context.TODO(), fmt.Sprintf(msg, v...))
}

func Warn(msg string) {
	_logger.Warn(context.TODO(), msg)
}
func Warnf(msg string, v ...any) {
	_logger.Warn(context.TODO(), fmt.Sprintf(msg, v...))
}

func Error(msg string) {
	_logger.Error(context.TODO(), msg)
}
func Errorf(msg string, v ...any) {
	_logger.Error(context.TODO(), fmt.Sprintf(msg, v...))
}

type Logger interface {
	Debug(ctx context.Context, msg string, data ...interface{})
	Info(ctx context.Context, msg string, data ...interface{})
	Warn(ctx context.Context, msg string, data ...interface{})
	Error(ctx context.Context, msg string, data ...interface{})
	Fatal(ctx context.Context, msg string, data ...interface{})
	Panic(ctx context.Context, msg string, data ...interface{})
}

type DefaultLogger struct {
	logger *log.Logger
}

func NewDefaultLogger(logger *log.Logger) *DefaultLogger {
	return &DefaultLogger{
		logger: logger,
	}
}
func (l *DefaultLogger) Debug(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > DebugLevel {
		return
	}
	l.logger.Printf("[debug] " + fmt.Sprint(append([]interface{}{msg}, data...)...))
}
func (l *DefaultLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > InfoLevel {
		return
	}
	l.logger.Printf("[info] " + fmt.Sprint(append([]interface{}{msg}, data...)...))

}
func (l *DefaultLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > WarnLevel {
		return
	}
	l.logger.Printf("[warn] " + fmt.Sprint(append([]interface{}{msg}, data...)...))

}
func (l *DefaultLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > ErrorLevel {
		return
	}
	l.logger.Printf("[error] " + fmt.Sprint(append([]interface{}{msg}, data...)...))

}
func (l *DefaultLogger) Fatal(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > FatalLevel {
		return
	}
	l.logger.Printf("[fatal] " + fmt.Sprint(append([]interface{}{msg}, data...)...))

}
func (l *DefaultLogger) Panic(ctx context.Context, msg string, data ...interface{}) {
	if logLevel > PanicLevel {
		return
	}
	l.logger.Printf("[panic] " + fmt.Sprint(append([]interface{}{msg}, data...)...))

}
