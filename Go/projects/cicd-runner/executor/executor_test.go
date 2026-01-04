package executor

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/projects/cicd-runner/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name     string
		execType string
		wantType string
	}{
		{"local executor", "local", "local"},
		{"mock executor", "mock", "mock"},
		{"default executor", "unknown", "local"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutor(tt.execType)
			assert.Equal(t, tt.wantType, exec.Type())
		})
	}
}

func TestMockExecutor(t *testing.T) {
	exec := NewMockExecutor()
	assert.Equal(t, "mock", exec.Type())

	ctx := context.Background()
	workspace := "/tmp/test-workspace"

	// Setup
	err := exec.Setup(ctx, workspace)
	assert.NoError(t, err)

	// Teardown
	err = exec.Teardown(ctx, workspace)
	assert.NoError(t, err)

	// Execute
	step := &pipeline.Step{
		Name:     "test-step",
		Commands: []string{"echo hello", "echo world"},
	}

	result, err := exec.Execute(ctx, step, nil, workspace)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.NotEmpty(t, result.Output)
	assert.Equal(t, step, result.Step)
}

func TestMockExecutorWithPredefinedResult(t *testing.T) {
	exec := NewMockExecutor()

	step := &pipeline.Step{
		Name:     "test-step",
		Commands: []string{"echo hello"},
	}

	// 设置预定义的结果
	expectedResult := &Result{
		Success:  false,
		ExitCode: 1,
		Output:   "predefined output",
		Error:    "predefined error",
	}
	exec.SetResult("test-step", expectedResult)

	ctx := context.Background()
	result, err := exec.Execute(ctx, step, nil, "/tmp/test")
	require.NoError(t, err)

	assert.False(t, result.Success)
	assert.Equal(t, 1, result.ExitCode)
	assert.Equal(t, "predefined output", result.Output)
	assert.Equal(t, "predefined error", result.Error)
}

func TestMockExecutorFailureSimulation(t *testing.T) {
	exec := NewMockExecutor()

	// 步骤名称包含 "fail" 应该模拟失败
	step := &pipeline.Step{
		Name:     "test-fail-step",
		Commands: []string{"echo hello"},
	}

	ctx := context.Background()
	result, err := exec.Execute(ctx, step, nil, "/tmp/test")
	require.NoError(t, err)

	assert.False(t, result.Success)
	assert.Equal(t, 1, result.ExitCode)
	assert.Contains(t, result.Error, "Mock execution failed")
}

func TestLocalExecutor(t *testing.T) {
	exec := NewLocalExecutor()
	assert.Equal(t, "local", exec.Type())

	ctx := context.Background()
	workspace, err := os.MkdirTemp("", "test-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspace)

	// Setup
	err = exec.Setup(ctx, workspace)
	assert.NoError(t, err)

	// Execute a simple command
	step := &pipeline.Step{
		Name:     "test-step",
		Commands: []string{"echo hello world"},
		Timeout:  10,
	}

	result, err := exec.Execute(ctx, step, nil, workspace)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "hello world")
	assert.Greater(t, result.Duration, time.Duration(0))

	// Teardown
	err = exec.Teardown(ctx, workspace)
	assert.NoError(t, err)
}

func TestLocalExecutorWithFailure(t *testing.T) {
	exec := NewLocalExecutor()

	ctx := context.Background()
	workspace, err := os.MkdirTemp("", "test-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspace)

	// Execute a failing command
	step := &pipeline.Step{
		Name:     "failing-step",
		Commands: []string{"false"}, // false 命令总是返回非零退出码
		Timeout:  10,
	}

	result, err := exec.Execute(ctx, step, nil, workspace)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEqual(t, 0, result.ExitCode)
}

func TestLocalExecutorWithTimeout(t *testing.T) {
	exec := NewLocalExecutor()

	ctx := context.Background()
	workspace, err := os.MkdirTemp("", "test-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspace)

	// Execute a command that would timeout
	step := &pipeline.Step{
		Name:     "timeout-step",
		Commands: []string{"sleep", "10"}, // 睡眠 10 秒
		Timeout:  1,                       // 但超时设置为 1 秒
	}

	result, err := exec.Execute(ctx, step, nil, workspace)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEqual(t, 0, result.ExitCode)
}

func TestLocalExecutorWithEnv(t *testing.T) {
	exec := NewLocalExecutor()

	ctx := context.Background()
	workspace, err := os.MkdirTemp("", "test-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspace)

	env := map[string]string{
		"TEST_VAR": "test_value",
	}

	step := &pipeline.Step{
		Name:     "env-step",
		Commands: []string{"echo $TEST_VAR"},
		Timeout:  10,
	}

	result, err := exec.Execute(ctx, step, env, workspace)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "test_value")
}

func TestLocalExecutorWithMultipleCommands(t *testing.T) {
	exec := NewLocalExecutor()

	ctx := context.Background()
	workspace, err := os.MkdirTemp("", "test-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspace)

	step := &pipeline.Step{
		Name: "multi-command-step",
		Commands: []string{
			"echo 'first command'",
			"echo 'second command'",
			"echo 'third command'",
		},
		Timeout: 10,
	}

	result, err := exec.Execute(ctx, step, nil, workspace)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "first command")
	assert.Contains(t, result.Output, "second command")
	assert.Contains(t, result.Output, "third command")
}
