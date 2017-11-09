/**
	log
	封装log相关操作
	方便以后统一管理日志打印
 */
package mylog

import (
	"log"
	"strings"
)

type Logger interface {
	Debug(a ...interface{})
	Warn(a ...interface{})
	Info(a ...interface{})
	Error(a ...interface{})
}

const (
	DEBUG = iota
	WARN
	INFO
	ERROR
)

type SimpleLogger struct {
	level int
}

func (simpleLogger SimpleLogger) Debug(a ...interface{}) {
	if DEBUG < simpleLogger.level {
		return
	}
	log.Println("DEBUG  ", a)
}

func (simpleLogger SimpleLogger) Warn(a ...interface{}) {
	if WARN < simpleLogger.level {
		return
	}
	log.Println("WARN  ", a)
}

func (simpleLogger SimpleLogger) Info(a ...interface{}) {
	if INFO < simpleLogger.level {
		return
	}
	log.Println("INFO  ", a)
}

func (simpleLogger SimpleLogger) Error(a ...interface{}) {
	if ERROR < simpleLogger.level {
		return
	}
	log.Println("ERROR  ", a)
}

var (
	logger Logger = SimpleLogger{level: DEBUG}
)

func InitLogger(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug", "0":
		logger = SimpleLogger{level: DEBUG}
	case "warn", "1":
		logger = SimpleLogger{level: WARN}
	case "info", "2":
		logger = SimpleLogger{level: INFO}
	case "error", "3":
		logger = SimpleLogger{level: ERROR}
	}
}

func Debug(a ...interface{}) {
	logger.Debug(a)
}

func Warn(a ...interface{}) {
	logger.Warn(a)
}

func Info(a ...interface{}) {
	logger.Info(a)
}

func Error(a ...interface{}) {
	logger.Error(a)
}
