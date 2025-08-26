package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *log.Logger
)

func initLogger(logDir string) {
	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 配置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "ddns.log"),
		MaxSize:    1,    // 每个日志文件最大1MB
		MaxAge:     7,    // 保留7天
		MaxBackups: 5,    // 最多保留5个备份文件
		LocalTime:  true, // 使用本地时间
		Compress:   true, // 压缩旧日志
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)

	// 创建自定义日志格式
	logger = log.New(multiWriter, "", 0)
}

// 自定义日志函数
func logInfo(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("[INFO] "+timestamp+" "+format, v...)
}

func logError(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("[ERROR] "+timestamp+" "+format, v...)
}

func logWarn(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("[WARN] "+timestamp+" "+format, v...)
}
