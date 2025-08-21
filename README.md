# IPv6 DNS更新器

一个自动获取本机IPv6地址并更新阿里云DNS解析的工具。

## 功能特性

1. 自动获取本机IPv6地址前4段
2. 支持配置多个域名和对应的IPv6后4段
3. 自动组装完整IPv6地址
4. 使用阿里云DNS SDK V2更新域名解析
5. 定时任务每10分钟执行一次
6. 智能判断地址变化，只在有变化时更新
7. 完善的日志记录

## GO项目开发说明

1.增加依赖项，如下，执行后能自动添加到go.mod中
```bash
go get github.com/alibabacloud-go/darabonba-openapi/v2/client
```

2.打linux下的arm64包
```bash
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -ldflags="-s -w" -o sara-go-tools .
```