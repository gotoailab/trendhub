package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	globalLogger *Logger
	once         sync.Once
)

// Logger 全局日志记录器
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
	mu          sync.Mutex
}

// Init 初始化全局 logger
func Init(logFilePath string) error {
	var err error
	once.Do(func() {
		globalLogger = &Logger{}
		err = globalLogger.setup(logFilePath)
	})
	return err
}

// setup 配置日志输出
func (l *Logger) setup(logFilePath string) error {
	var writers []io.Writer
	
	// 总是输出到标准输出
	writers = append(writers, os.Stdout)
	
	// 如果指定了日志文件，同时输出到文件
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		l.logFile = file
		writers = append(writers, file)
	}
	
	multiWriter := io.MultiWriter(writers...)
	l.infoLogger = log.New(multiWriter, "[INFO] ", log.LstdFlags)
	l.errorLogger = log.New(multiWriter, "[ERROR] ", log.LstdFlags)
	
	return nil
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// GetGlobalLogger 获取全局 logger 实例
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// 如果没有初始化，使用默认配置（仅输出到 stdout）
		_ = Init("")
	}
	return globalLogger
}

// Info 打印 Info 级别日志
func Info(v ...interface{}) {
	GetGlobalLogger().infoLogger.Println(v...)
}

// Infof 格式化打印 Info 级别日志
func Infof(format string, v ...interface{}) {
	GetGlobalLogger().infoLogger.Printf(format, v...)
}

// Error 打印 Error 级别日志
func Error(v ...interface{}) {
	GetGlobalLogger().errorLogger.Println(v...)
}

// Errorf 格式化打印 Error 级别日志
func Errorf(format string, v ...interface{}) {
	GetGlobalLogger().errorLogger.Printf(format, v...)
}

// Fatal 打印 Fatal 级别日志并退出程序
func Fatal(v ...interface{}) {
	GetGlobalLogger().errorLogger.Fatalln(v...)
}

// Fatalf 格式化打印 Fatal 级别日志并退出程序
func Fatalf(format string, v ...interface{}) {
	GetGlobalLogger().errorLogger.Fatalf(format, v...)
}

// Close 关闭全局 logger
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

