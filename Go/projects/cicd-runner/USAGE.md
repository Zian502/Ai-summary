# 使用指南

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行示例（Mock 模式）

```bash
# 使用 Mock 模式运行，不会实际执行命令
go run main.go -mock -pipeline examples/pipeline.yaml
```

### 3. 运行示例（Local 模式）

```bash
# 使用 Local 模式运行，实际执行命令
go run main.go -pipeline examples/pipeline.yaml -config examples/config.yaml
```

### 4. 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定包的测试
go test ./config -v
go test ./executor -v
go test ./pipeline -v
go test ./runner -v
```

## 命令行参数

- `-config <path>`: 指定配置文件路径（可选）
- `-pipeline <path>`: 指定 Pipeline 配置文件路径（默认：examples/pipeline.yaml）
- `-mock`: 使用 Mock 模式（不实际执行命令）
- `-version`: 显示版本信息

## 配置说明

### 系统配置（config.yaml）

```yaml
runner:
  capacity: 10        # 并发执行容量
  timeout: 3600s      # 超时时间
  workspace: /tmp/cicd-workspace  # 工作空间目录

executor:
  type: local         # 执行器类型：local 或 mock
  env:
    GO_VERSION: "1.21"
    CGO_ENABLED: "0"

log:
  level: info         # 日志级别：debug, info, warn, error
  format: text        # 日志格式：json, text
```

### Pipeline 配置（pipeline.yaml）

```yaml
name: example-pipeline
version: "1.0"
workspace: /tmp/cicd-workspace

env:
  PROJECT_NAME: "cicd-runner"
  BUILD_VERSION: "0.1.0"

concurrency: 1  # 并发执行数，1 表示串行执行

steps:
  - name: build
    commands:
      - echo "Building..."
      - go build -o app ./...
    env:
      CGO_ENABLED: "0"
    timeout: 300
    on_success:
      - echo "Build succeeded"
    on_failure:
      - echo "Build failed"
  
  - name: test
    commands:
      - go test -v ./...
    timeout: 300
    when: always
```

## 环境变量配置

可以通过环境变量覆盖配置：

- `CICD_RUNNER_CAPACITY`: Runner 并发容量
- `CICD_RUNNER_TIMEOUT`: Runner 超时时间
- `CICD_RUNNER_WORKSPACE`: 工作空间目录
- `CICD_EXECUTOR_TYPE`: 执行器类型（local/mock）
- `CICD_LOG_LEVEL`: 日志级别

## 架构说明

### 核心组件

1. **Config**: 配置管理模块
   - 支持从 YAML 文件加载配置
   - 支持从环境变量加载配置
   - 支持配置验证和默认值

2. **Pipeline**: Pipeline 定义和解析
   - 支持多步骤定义
   - 支持步骤条件执行
   - 支持全局和步骤级环境变量

3. **Executor**: 执行器接口和实现
   - **LocalExecutor**: 本地执行器，实际执行命令
   - **MockExecutor**: Mock 执行器，用于测试和开发

4. **Runner**: Runner 核心
   - 协调 Pipeline 的执行
   - 管理并发执行
   - 处理超时和错误

### 执行流程

1. 加载配置（Config）
2. 加载 Pipeline 定义
3. 创建执行器（Executor）
4. 设置执行环境
5. 依次执行步骤（Steps）
6. 收集执行结果
7. 清理执行环境

## 扩展开发

### 添加新的执行器

1. 实现 `Executor` 接口
2. 在 `NewExecutor` 函数中注册

```go
func NewExecutor(executorType string) Executor {
    switch executorType {
    case "mock":
        return NewMockExecutor()
    case "local":
        return NewLocalExecutor()
    case "docker":  // 新增
        return NewDockerExecutor()
    default:
        return NewLocalExecutor()
    }
}
```

### 添加新的步骤条件

在 `Step.ShouldRun()` 方法中添加新的条件逻辑。

## 故障排查

### 问题：命令执行失败

- 检查命令是否正确
- 检查工作空间目录是否存在
- 检查环境变量是否正确设置

### 问题：超时错误

- 增加步骤的 `timeout` 值
- 检查命令是否卡住

### 问题：环境变量未生效

- 检查环境变量的优先级（步骤级 > 全局 > 配置级）
- 检查环境变量名称是否正确

