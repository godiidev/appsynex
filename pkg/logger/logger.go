// File: pkg/logger/logger.go
// Tạo tại: pkg/logger/logger.go
// Mục đích: Logging utilities for the application

package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type Logger struct {
	level  Level
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(INFO)
}

func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

func SetLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		defaultLogger.level = DEBUG
	case "INFO":
		defaultLogger.level = INFO
	case "WARN":
		defaultLogger.level = WARN
	case "ERROR":
		defaultLogger.level = ERROR
	case "FATAL":
		defaultLogger.level = FATAL
	default:
		defaultLogger.level = INFO
	}
}

func (l *Logger) log(level Level, v ...interface{}) {
	if level < l.level {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]
	message := fmt.Sprint(v...)
	
	l.logger.Printf("[%s] %s: %s", timestamp, levelName, message)
	
	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) logf(level Level, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]
	message := fmt.Sprintf(format, v...)
	
	l.logger.Printf("[%s] %s: %s", timestamp, levelName, message)
	
	if level == FATAL {
		os.Exit(1)
	}
}

// Package level functions
func Debug(v ...interface{}) {
	defaultLogger.log(DEBUG, v...)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.logf(DEBUG, format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.log(INFO, v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.logf(INFO, format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.log(WARN, v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.logf(WARN, format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.log(ERROR, v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.logf(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.log(FATAL, v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.logf(FATAL, format, v...)
}