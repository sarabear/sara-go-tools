package main

import (
	"flag"
	"time"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	logPath := flag.String("log-dir", "logs", "日志目录路径")

	flag.Parse()

	// 初始化日志系统
	initLogger(*logPath)
	logInfo("程序启动，配置文件: %s, 日志目录: %s", *configPath, *logPath)

	config, err := loadConfig(*configPath)
	if err != nil {
		logError("加载配置失败: %v", err)
	}

	client, err := newDNSClient(config)
	if err != nil {
		logError("初始化DNS客户端失败: %v", err)
	}

	ipv6Prefix := getLocalIPv6Prefix()
	if ipv6Prefix == "" {
		logWarn("未获取到本机IPv6地址，程序退出。")
		return
	}
	logInfo("本机IPv6前缀: %s", ipv6Prefix)

	// 初始运行
	runOnce(client, config, ipv6Prefix)

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		logInfo("定时任务执行中...")
		ipv6Prefix = getLocalIPv6Prefix()
		logInfo("本机IPv6前缀: %s", ipv6Prefix)
		if ipv6Prefix == "" {
			logWarn("未获取到本机IPv6地址，定时任务跳过。")
			break
		}
		runOnce(client, config, ipv6Prefix)
	}
}
