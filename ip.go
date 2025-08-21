package main

import (
	"net"
	"strings"
)

func getLocalIPv6Prefix() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logError("获取网络接口失败: %v", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() == nil {
			ip := ipnet.IP.String()
			if strings.HasPrefix(ip, "fe80") || strings.HasPrefix(ip, "::1") {
				continue
			}
			parts := strings.Split(ip, ":")
			if len(parts) >= 4 {
				return strings.Join(parts[:4], ":")
			}
		}
	}
	return ""
}
