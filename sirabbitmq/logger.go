package sirabbitmq

import (
	"log"
	"os"
)

type LogLevel uint8

const (
	DebugLevel LogLevel = 1 + iota
	InfoLevel
	WarnLevel
	ErrorLevel
	IgnoreLevel
)

var (
	logLevel = DebugLevel

	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	debugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func SetLevel(lvl LogLevel) {
	logLevel = lvl
}

func Debug(msg string) {
	if logLevel <= DebugLevel {
		debugLogger.Println("[debug] " + msg)
	}
}
func Debugf(msg string, v ...any) {
	if logLevel <= DebugLevel {
		debugLogger.Printf("[debug] "+msg, v...)
	}
}

func Info(msg string) {
	if logLevel <= InfoLevel {
		infoLogger.Println("[info] " + msg)
	}
}
func Infof(msg string, v ...any) {
	if logLevel <= InfoLevel {
		infoLogger.Printf("[info] "+msg, v...)
	}
}

func Warn(msg string) {
	if logLevel <= WarnLevel {
		warnLogger.Println("[warn] " + msg)
	}
}
func Warnf(msg string, v ...any) {
	if logLevel <= WarnLevel {
		warnLogger.Printf("[warn] "+msg, v...)
	}
}

func Error(msg string) {
	if logLevel <= ErrorLevel {
		errorLogger.Println("[error] " + msg)
	}
}
func Errorf(msg string, v ...any) {
	if logLevel <= ErrorLevel {
		errorLogger.Printf("[error] "+msg, v...)
	}
}
