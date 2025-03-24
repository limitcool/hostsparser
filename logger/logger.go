package logger

import (
	"fmt"
	"log"
)

// Logger 接口定义了分级日志记录的方法
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// 默认日志记录器，使用标准库实现
type defaultLogger struct{}

func (l *defaultLogger) Debug(args ...interface{}) {
	log.Print("[DEBUG] ", fmt.Sprint(args...))
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *defaultLogger) Info(args ...interface{}) {
	log.Print("[INFO] ", fmt.Sprint(args...))
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *defaultLogger) Warn(args ...interface{}) {
	log.Print("[WARN] ", fmt.Sprint(args...))
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func (l *defaultLogger) Error(args ...interface{}) {
	log.Print("[ERROR] ", fmt.Sprint(args...))
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// 静默日志记录器，忽略所有日志
type silentLogger struct{}

func (l *silentLogger) Debug(args ...interface{})                 {}
func (l *silentLogger) Debugf(format string, args ...interface{}) {}
func (l *silentLogger) Info(args ...interface{})                  {}
func (l *silentLogger) Infof(format string, args ...interface{})  {}
func (l *silentLogger) Warn(args ...interface{})                  {}
func (l *silentLogger) Warnf(format string, args ...interface{})  {}
func (l *silentLogger) Error(args ...interface{})                 {}
func (l *silentLogger) Errorf(format string, args ...interface{}) {}

// 全局日志实例
var globalLogger Logger = &defaultLogger{}

// SetLogger 设置全局日志记录器
func SetLogger(logger Logger) {
	if logger != nil {
		globalLogger = logger
	}
}

// DisableLogging 禁用所有日志
func DisableLogging() {
	globalLogger = &silentLogger{}
}

// GetLogger 获取当前全局日志记录器
func GetLogger() Logger {
	return globalLogger
}

// 方便直接调用的包级函数
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}
