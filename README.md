# Hosts文件解析器

[English Version](README_EN.md)

这是一个用Go语言编写的hosts文件解析器，它使用词法分析和语法分析技术来解析和修改hosts文件。作为一个库项目，它提供了简单易用的API来操作和管理hosts文件。

## 功能特点

- 使用词法分析和语法分析技术解析hosts文件
- 支持添加、删除和查找hosts条目
- 保留注释和空行
- 支持IPv4和IPv6地址格式验证
- 支持IP地址对应多个域名的映射关系
- 支持按主机名操作hosts条目
- 支持加载和保存hosts文件
- 灵活的分级日志系统，支持自定义日志接口

## 安装

```bash
go get github.com/limitcool/hostsparser
```

## 库使用方法

### 导入包

```go
import (
    "github.com/limitcool/hostsparser/hosts"
    "github.com/limitcool/hostsparser/logger"
)
```

### 加载hosts文件

```go
// 加载系统hosts文件
hostsPath := hosts.GetSystemHostsPath()
hostsFile, err := hosts.LoadHostsFile(hostsPath)
if err != nil {
    // 处理错误
}

// 或者创建新的hosts文件
hostsFile := hosts.NewHostsFile()
```

### 修改hosts条目

```go
// 设置主机名对应的IP
err := hostsFile.SetHostIP("example.com", "127.0.0.1")

// 批量设置多个主机名映射到同一IP
err = hostsFile.SetMultipleHostIPs([]string{"example.com", "www.example.com"}, "127.0.0.1")

// 删除主机名映射
modified, err := hostsFile.RemoveHost("example.com")
```

### 查询hosts条目

```go
// 获取主机名对应的IP
ip, err := hostsFile.GetHostIP("example.com")

// 获取IP对应的所有主机名
domains, err := hostsFile.GetHostsByIP("127.0.0.1")

// 获取所有IP-域名对
pairs := hostsFile.GetAllIPDomainPairs()

// 按IP或域名筛选
filteredByIP := hostsFile.FilterIPDomainPairs("127.0.0.1", "")
filteredByDomain := hostsFile.FilterIPDomainPairs("", "example.com")
```

### 保存hosts文件

```go
// 保存到原文件
err := hostsFile.SaveHostsFile("")

// 保存到新文件
err = hostsFile.SaveHostsFile("./hosts.new")
```

### 日志系统

本库提供了灵活的分级日志系统，支持Debug、Info、Warn和Error四个日志级别。

#### 基本日志使用

```go
// 直接使用包级函数记录日志
logger.Debug("这是调试信息")
logger.Info("这是普通信息")
logger.Warn("这是警告信息")
logger.Error("这是错误信息")

// 使用格式化函数
logger.Debugf("调试: %s", "详细信息")
logger.Infof("信息: %s", "详细信息")
logger.Warnf("警告: %s", "详细信息")
logger.Errorf("错误: %s", "详细信息")
```

#### 禁用日志

```go
// 完全禁用所有日志输出
logger.DisableLogging()
```

#### 自定义日志接口

您可以实现自己的日志接口，例如集成zap、logrus等高性能日志库：

```go
// 1. 实现Logger接口
type MyCustomLogger struct {
    // 自定义字段
}

func (l *MyCustomLogger) Debug(args ...interface{}) {
    // 自定义实现
}

func (l *MyCustomLogger) Debugf(format string, args ...interface{}) {
    // 自定义实现
}

// 实现其他方法...

// 2. 设置为全局日志记录器
logger.SetLogger(&MyCustomLogger{})
```

#### Zap日志库适配示例

```go
import (
    "github.com/limitcool/hostsparser/logger"
    "go.uber.org/zap"
)

// 在主程序中使用
func main() {
    // 创建zap日志实例
    zapLogger, _ := zap.NewProduction()
    defer zapLogger.Sync()

    // 直接使用zap的Sugar接口适配Logger接口
    logger.SetLogger(zapLogger.Sugar())

    // 现在所有日志都会通过zap记录
}
```

## 项目结构

- `lexer/` - 词法分析器，将hosts文件内容分解为词法单元
- `parser/` - 语法分析器，解析词法单元并构建hosts条目结构
- `hosts/` - hosts文件操作的核心功能
- `logger/` - 分级日志系统，支持自定义日志接口

## 实现细节

### 词法分析和语法分析

本项目使用经典的编译器前端技术来解析hosts文件：

1. **词法分析（Lexical Analysis）**：将输入文本分解成有意义的标记（token），如IP地址、域名、注释等。

2. **语法分析（Syntax Analysis）**：分析标记序列，构建结构化的hosts条目。

这种方法使得解析过程更加健壮，能够处理各种格式的hosts文件，包括注释、空行和格式不正确的行。

### Hosts文件格式

标准hosts文件格式如下：

```
# 这是注释行
127.0.0.1 localhost
::1 localhost ipv6-localhost ipv6-loopback
192.168.1.1 example.com www.example.com # 这是行末注释
```

## 许可证

MIT

## 贡献

欢迎提交 Pull Requests 和 Issues。