package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 10, cfg.Runner.Capacity)
	assert.Equal(t, 3600*time.Second, cfg.Runner.Timeout)
	assert.Equal(t, "/tmp/cicd-workspace", cfg.Runner.Workspace)
	assert.Equal(t, "local", cfg.Executor.Type)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Runner: RunnerConfig{
					Capacity:  10,
					Timeout:   3600 * time.Second,
					Workspace: "/tmp/test",
				},
				Executor: ExecutorConfig{
					Type: "local",
					Env:  make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid capacity",
			config: &Config{
				Runner: RunnerConfig{
					Capacity:  0,
					Timeout:   3600 * time.Second,
					Workspace: "/tmp/test",
				},
				Executor: ExecutorConfig{
					Type: "local",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid executor type",
			config: &Config{
				Runner: RunnerConfig{
					Capacity:  10,
					Timeout:   3600 * time.Second,
					Workspace: "/tmp/test",
				},
				Executor: ExecutorConfig{
					Type: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetEnv("TEST_KEY", "config_value")

	// 设置环境变量
	os.Setenv("TEST_KEY", "env_value")
	defer os.Unsetenv("TEST_KEY")

	// 环境变量优先
	assert.Equal(t, "env_value", cfg.GetEnv("TEST_KEY"))

	// 只有配置中的值
	os.Unsetenv("TEST_KEY")
	assert.Equal(t, "config_value", cfg.GetEnv("TEST_KEY"))

	// 都不存在
	assert.Equal(t, "", cfg.GetEnv("NON_EXISTENT"))
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("CICD_RUNNER_CAPACITY", "20")
	os.Setenv("CICD_RUNNER_TIMEOUT", "7200s")
	os.Setenv("CICD_RUNNER_WORKSPACE", "/tmp/custom")
	os.Setenv("CICD_EXECUTOR_TYPE", "mock")
	os.Setenv("CICD_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("CICD_RUNNER_CAPACITY")
		os.Unsetenv("CICD_RUNNER_TIMEOUT")
		os.Unsetenv("CICD_RUNNER_WORKSPACE")
		os.Unsetenv("CICD_EXECUTOR_TYPE")
		os.Unsetenv("CICD_LOG_LEVEL")
	}()

	cfg := LoadFromEnv()

	assert.Equal(t, 20, cfg.Runner.Capacity)
	assert.Equal(t, 7200*time.Second, cfg.Runner.Timeout)
	assert.Equal(t, "/tmp/custom", cfg.Runner.Workspace)
	assert.Equal(t, "mock", cfg.Executor.Type)
	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestLoad(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
runner:
  capacity: 5
  timeout: 1800s
  workspace: /tmp/test-workspace

executor:
  type: mock
  env:
    TEST_VAR: "test_value"

log:
  level: debug
  format: json
`

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, 5, cfg.Runner.Capacity)
	assert.Equal(t, 1800*time.Second, cfg.Runner.Timeout)
	assert.Equal(t, "/tmp/test-workspace", cfg.Runner.Workspace)
	assert.Equal(t, "mock", cfg.Executor.Type)
	assert.Equal(t, "test_value", cfg.Executor.Env["TEST_VAR"])
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
}
