package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var (
	logger *log.Logger = log.StandardLogger()
)

func init() {
	LogLevel(DebugLevel)
	//logger.SetFormatter(&log.JSONFormatter{})
	Debug("init infrastructure/logging", "Set LogLevel to %v", DebugLevel)
}

func LogLevel(level Level) {
	logger.SetLevel(log.Level(level))
}

func generateFields(operation string) log.Fields {
	pc, file, line, _ := runtime.Caller(2)
	class := runtime.FuncForPC(pc).Name()
	where := fmt.Sprintf("%v-%v:%v", class, file, line)
	return log.Fields{
		"where":     where,
		"when":      time.Now().Format("2006-01-02:15-04-05.MST"),
		"operation": operation,
	}
}

func replaceFormat(format string) string {
	return fmt.Sprintf("%v", format)
}

func Info(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Infof(replaceFormat(format), args...)
}

func Error(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Errorf(replaceFormat(format), args...)
}

func Debug(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Debugf(replaceFormat(format), args...)
}

func Trace(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Tracef(replaceFormat(format), args...)
}

func Warning(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Warningf(replaceFormat(format), args...)
}

func Fatal(operation string, format string, args ...interface{}) {
	logger.WithFields(generateFields(operation)).Fatalf(replaceFormat(format), args...)
}
