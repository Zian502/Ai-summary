# 简易版 Go CI/CD Runner

一个参考 Drone-runner-go 架构设计的简易 CI/CD 执行器，支持基础的 CI/CD 场景。

## 功能特性

- ✅ 支持 Pipeline 和 Step 定义
- ✅ 支持 Mock 模式（用于测试和开发）
- ✅ 支持 Test 执行
- ✅ 可自定义配置参数
- ✅ 简洁的架构设计

## 项目结构

```
projects/
├── main.go              # 主程序入口
├── config/              # 配置管理
│   ├── config.go       # 配置结构定义
│   └── loader.go       # 配置加载器
├── pipeline/            # Pipeline 定义
│   ├── pipeline.go     # Pipeline 结构
│   └── step.go         # Step 结构
├── executor/            # 执行器
│   ├── executor.go     # 执行器接口
│   ├── local.go        # 本地执行器
│   └── mock.go         # Mock 执行器
├── runner/              # Runner 核心
│   └── runner.go       # Runner 实现
└── examples/            # 示例配置
    ├── pipeline.yaml   # Pipeline 配置示例
    └── config.yaml     # 系统配置示例
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行示例

```bash
# 使用默认配置运行
go run main.go

# 指定配置文件
go run main.go -config examples/config.yaml

# 使用 Mock 模式
go run main.go -mock

# 运行测试
go test ./...
```

## 配置说明

### Pipeline 配置

Pipeline 配置文件定义了 CI/CD 的执行流程：

```yaml
name: example-pipeline
steps:
  - name: build
    image: golang:1.21
    commands:
      - go build -o app ./...
  
  - name: test
    image: golang:1.21
    commands:
      - go test ./...
  
  - name: lint
    image: golangci/golangci-lint:latest
    commands:
      - golangci-lint run
```

### 系统配置

系统配置文件定义了 Runner 的运行参数：

```yaml
runner:
  capacity: 10
  timeout: 3600
  workspace: /tmp/cicd-workspace

executor:
  type: local  # local 或 mock
  env:
    - GO_VERSION=1.21
    - CGO_ENABLED=0
```

## 使用示例

### 基本使用

```go
package main

import (
    "github.com/example/cicd-runner/config"
    "github.com/example/cicd-runner/runner"
)

func main() {
    cfg := config.Load("examples/config.yaml")
    r := runner.New(cfg)
    r.Run("examples/pipeline.yaml")
}
```

### Mock 模式

Mock 模式用于测试和开发，不会实际执行命令：

```go
cfg := config.Load("examples/config.yaml")
cfg.Executor.Type = "mock"
r := runner.New(cfg)
r.Run("examples/pipeline.yaml")
```

## 测试

运行所有测试：

```bash
go test ./... -v
```

运行特定包的测试：

```bash
go test ./executor -v
```

## 架构设计

本实现参考了 Drone-runner-go 的架构：

1. **Config**: 配置管理模块，支持从文件或环境变量加载配置
2. **Pipeline**: Pipeline 和 Step 的定义和解析
3. **Executor**: 执行器接口，支持多种执行器实现（Local、Mock）
4. **Runner**: Runner 核心，负责协调 Pipeline 的执行

## License

MIT

