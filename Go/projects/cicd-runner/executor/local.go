package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/projects/cicd-runner/pipeline"
)

// LocalExecutor 本地执行器，在本地执行命令
type LocalExecutor struct{}

// NewLocalExecutor 创建本地执行器
func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

// Type 返回执行器类型
func (e *LocalExecutor) Type() string {
	return "local"
}

// Setup 设置执行环境
func (e *LocalExecutor) Setup(ctx context.Context, workspace string) error {
	return os.MkdirAll(workspace, 0755)
}

// Teardown 清理执行环境
func (e *LocalExecutor) Teardown(ctx context.Context, workspace string) error {
	// 可选：清理工作空间
	// return os.RemoveAll(workspace)
	return nil
}

// Execute 执行单个步骤
func (e *LocalExecutor) Execute(ctx context.Context, step *pipeline.Step, env map[string]string, workspace string) (*Result, error) {
	startTime := time.Now()

	// 创建带超时的上下文
	stepCtx := ctx
	if step.Timeout > 0 {
		var cancel context.CancelFunc
		stepCtx, cancel = context.WithTimeout(ctx, time.Duration(step.Timeout)*time.Second)
		defer cancel()
	}

	// 准备环境变量
	execEnv := os.Environ()
	for k, v := range env {
		execEnv = append(execEnv, fmt.Sprintf("%s=%s", k, v))
	}
	if step.Env != nil {
		for k, v := range step.Env {
			execEnv = append(execEnv, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// 执行所有命令
	var output strings.Builder
	var lastErr error
	var exitCode int

	for _, cmdStr := range step.Commands {
		// 解析命令 - 支持 shell 命令
		// 如果命令包含管道、重定向等 shell 特性，或包含环境变量，使用 sh -c 执行
		needsShell := strings.Contains(cmdStr, "|") || strings.Contains(cmdStr, ">") ||
			strings.Contains(cmdStr, "<") || strings.Contains(cmdStr, "&&") ||
			strings.Contains(cmdStr, "||") || strings.Contains(cmdStr, ";") ||
			strings.Contains(cmdStr, "$")

		if needsShell {
			// 使用 shell 执行复杂命令
			cmd := exec.CommandContext(stepCtx, "sh", "-c", cmdStr)
			cmd.Dir = workspace
			cmd.Env = execEnv
			cmd.Stdout = &output
			cmd.Stderr = &output

			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					exitCode = 1
				}
				lastErr = err
				output.WriteString(fmt.Sprintf("\nError: %v\n", err))

				// 如果步骤失败，执行 on_failure 命令
				if len(step.OnFailure) > 0 {
					e.executeHooks(stepCtx, step.OnFailure, execEnv, workspace, &output)
				}
				break
			}
			continue
		}

		// 简单命令：直接解析
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.CommandContext(stepCtx, parts[0], parts[1:]...)
		cmd.Dir = workspace
		cmd.Env = execEnv
		cmd.Stdout = &output
		cmd.Stderr = &output

		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = 1
			}
			lastErr = err
			output.WriteString(fmt.Sprintf("\nError: %v\n", err))

			// 如果步骤失败，执行 on_failure 命令
			if len(step.OnFailure) > 0 {
				e.executeHooks(stepCtx, step.OnFailure, execEnv, workspace, &output)
			}
			break
		}
	}

	// 如果所有命令成功，执行 on_success 命令
	if lastErr == nil && len(step.OnSuccess) > 0 {
		e.executeHooks(stepCtx, step.OnSuccess, execEnv, workspace, &output)
	}

	duration := time.Since(startTime)

	result := &Result{
		Success:  lastErr == nil,
		ExitCode: exitCode,
		Output:   output.String(),
		Error:    "",
		Duration: duration,
		Step:     step,
	}

	if lastErr != nil {
		result.Error = lastErr.Error()
	}

	return result, nil
}

// executeHooks 执行钩子命令
func (e *LocalExecutor) executeHooks(ctx context.Context, hooks []string, env []string, workspace string, output *strings.Builder) {
	for _, hook := range hooks {
		parts := strings.Fields(hook)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		cmd.Dir = workspace
		cmd.Env = env
		cmd.Stdout = output
		cmd.Stderr = output

		_ = cmd.Run() // 忽略钩子命令的错误
	}
}

// ensureWorkspace 确保工作空间存在
func (e *LocalExecutor) ensureWorkspace(workspace string) error {
	return os.MkdirAll(workspace, 0755)
}
