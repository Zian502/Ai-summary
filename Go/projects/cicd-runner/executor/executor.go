package executor

import (
	"context"
	"time"

	"github.com/projects/cicd-runner/pipeline"
)

// Result 执行结果
type Result struct {
	Success  bool           // 是否成功
	ExitCode int            // 退出码
	Output   string         // 输出内容
	Error    string         // 错误信息
	Duration time.Duration  // 执行耗时
	Step     *pipeline.Step // 执行的步骤
}

// Executor 执行器接口
type Executor interface {
	// Execute 执行单个步骤
	Execute(ctx context.Context, step *pipeline.Step, env map[string]string, workspace string) (*Result, error)

	// Setup 设置执行环境
	Setup(ctx context.Context, workspace string) error

	// Teardown 清理执行环境
	Teardown(ctx context.Context, workspace string) error

	// Type 返回执行器类型
	Type() string
}

// NewExecutor 根据类型创建执行器
func NewExecutor(executorType string) Executor {
	switch executorType {
	case "mock":
		return NewMockExecutor()
	case "local":
		return NewLocalExecutor()
	default:
		return NewLocalExecutor()
	}
}
