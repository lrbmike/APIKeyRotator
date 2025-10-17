package logger

import (
	"log"
	"os"
)

// Setup 配置日志记录器
func Setup() {
	// 设置日志输出格式，包含日期、时间、文件名和行号
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// 设置日志输出到标准输出
	log.SetOutput(os.Stdout)
	
	log.Println("Logging configured successfully.")
}

// Info 记录信息级别日志
func Info(v ...interface{}) {
	log.Println("[INFO]", v)
}

// Infof 记录格式化信息级别日志
func Infof(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Warning 记录警告级别日志
func Warning(v ...interface{}) {
	log.Println("[WARNING]", v)
}

// Warningf 记录格式化警告级别日志
func Warningf(format string, v ...interface{}) {
	log.Printf("[WARNING] "+format, v...)
}

// Error 记录错误级别日志
func Error(v ...interface{}) {
	log.Println("[ERROR]", v)
}

// Errorf 记录格式化错误级别日志
func Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// Fatal 记录致命错误并退出程序
func Fatal(v ...interface{}) {
	log.Fatal("[FATAL]", v)
}

// Fatalf 记录格式化致命错误并退出程序
func Fatalf(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}