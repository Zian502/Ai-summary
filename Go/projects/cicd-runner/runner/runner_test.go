package runner

import (
	"os"
	"testing"

	"github.com/projects/cicd-runner/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := config.DefaultConfig()
	r := New(cfg)

	assert.NotNil(t, r)
	assert.Equal(t, cfg, r.config)
	assert.NotNil(t, r.executor)
}

func TestRun(t *testing.T) {
	// 创建临时 Pipeline 文件
	tmpFile, err := os.CreateTemp("", "test-pipeline-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	pipelineContent := `
name: test-pipeline
version: "1.0"

steps:
  - name: step1
    commands:
      - echo "Hello from step1"
    timeout: 10
  
  - name: step2
    commands:
      - echo "Hello from step2"
    timeout: 10
`

	_, err = tmpFile.WriteString(pipelineContent)
	require.NoError(t, err)
	tmpFile.Close()

	// 使用 Mock 执行器
	cfg := config.DefaultConfig()
	cfg.Executor.Type = "mock"
	cfg.Runner.Workspace = "/tmp/test-runner-workspace"

	r := New(cfg)
	err = r.Run(tmpFile.Name())
	assert.NoError(t, err)
}

func TestRunWithFailure(t *testing.T) {
	// 创建临时 Pipeline 文件，包含一个会失败的步骤
	tmpFile, err := os.CreateTemp("", "test-pipeline-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	pipelineContent := `
name: test-pipeline
version: "1.0"

steps:
  - name: success-step
    commands:
      - echo "This will succeed"
    timeout: 10
  
  - name: fail-step
    commands:
      - echo "This will fail"
    timeout: 10
`

	_, err = tmpFile.WriteString(pipelineContent)
	require.NoError(t, err)
	tmpFile.Close()

	// 使用 Mock 执行器，并设置失败结果
	cfg := config.DefaultConfig()
	cfg.Executor.Type = "mock"
	cfg.Runner.Workspace = "/tmp/test-runner-workspace"

	r := New(cfg)

	// Mock 执行器会自动将包含 "fail" 的步骤标记为失败
	err = r.Run(tmpFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pipeline failed")
}

func TestRunWithInvalidPipeline(t *testing.T) {
	cfg := config.DefaultConfig()
	r := New(cfg)

	// 尝试运行不存在的 Pipeline 文件
	err := r.Run("/nonexistent/pipeline.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load pipeline")
}
