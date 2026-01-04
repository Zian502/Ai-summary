package executor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/projects/cicd-runner/pipeline"
)

// MockExecutor Mock 执行器，用于测试和开发，不实际执行命令
type MockExecutor struct {
	results map[string]*Result // 预定义的执行结果
}

// NewMockExecutor 创建 Mock 执行器
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		results: make(map[string]*Result),
	}
}

// SetResult 设置步骤的预定义结果（用于测试）
func (e *MockExecutor) SetResult(stepName string, result *Result) {
	e.results[stepName] = result
}

// Type 返回执行器类型
func (e *MockExecutor) Type() string {
	return "mock"
}

// Setup 设置执行环境（Mock 模式下不做任何操作）
func (e *MockExecutor) Setup(ctx context.Context, workspace string) error {
	return nil
}

// Teardown 清理执行环境（Mock 模式下不做任何操作）
func (e *MockExecutor) Teardown(ctx context.Context, workspace string) error {
	return nil
}

// Execute 执行单个步骤（Mock 模式）
func (e *MockExecutor) Execute(ctx context.Context, step *pipeline.Step, env map[string]string, workspace string) (*Result, error) {
	startTime := time.Now()

	// 检查是否有预定义的结果
	if result, ok := e.results[step.Name]; ok {
		result.Step = step
		result.Duration = time.Since(startTime)
		return result, nil
	}

	// 模拟执行时间
	time.Sleep(100 * time.Millisecond)

	// 生成模拟输出
	var output strings.Builder
	output.WriteString(fmt.Sprintf("[MOCK] Executing step: %s\n", step.Name))
	output.WriteString(fmt.Sprintf("[MOCK] Commands: %v\n", step.Commands))
	output.WriteString(fmt.Sprintf("[MOCK] Workspace: %s\n", workspace))

	if len(env) > 0 {
		output.WriteString("[MOCK] Environment variables:\n")
		for k, v := range env {
			output.WriteString(fmt.Sprintf("  %s=%s\n", k, v))
		}
	}

	// 模拟执行命令
	for _, cmd := range step.Commands {
		output.WriteString(fmt.Sprintf("[MOCK] Running: %s\n", cmd))
		output.WriteString("[MOCK] ✓ Command completed successfully\n")
	}

	duration := time.Since(startTime)

	// 默认返回成功结果
	result := &Result{
		Success:  true,
		ExitCode: 0,
		Output:   output.String(),
		Error:    "",
		Duration: duration,
		Step:     step,
	}

	// 如果步骤名称包含 "fail" 或 "error"，模拟失败
	stepNameLower := strings.ToLower(step.Name)
	if strings.Contains(stepNameLower, "fail") ||
		strings.Contains(stepNameLower, "error") {
		result.Success = false
		result.ExitCode = 1
		result.Error = "Mock execution failed (simulated)"
	}

	return result, nil
}
