package config

import (
	"fmt"
	"os"
	"time"
)

// Config 系统配置结构
type Config struct {
	Runner   RunnerConfig   `yaml:"runner"`
	Executor ExecutorConfig `yaml:"executor"`
	Log      LogConfig      `yaml:"log"`
}

// RunnerConfig Runner 配置
type RunnerConfig struct {
	Capacity  int           `yaml:"capacity"`  // 并发执行容量
	Timeout   time.Duration `yaml:"timeout"`   // 超时时间（秒）
	Workspace string        `yaml:"workspace"` // 工作空间目录
}

// ExecutorConfig 执行器配置
type ExecutorConfig struct {
	Type string            `yaml:"type"` // 执行器类型：local, mock
	Env  map[string]string `yaml:"env"`  // 环境变量
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `yaml:"level"`  // 日志级别：debug, info, warn, error
	Format string `yaml:"format"` // 日志格式：json, text
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Runner: RunnerConfig{
			Capacity:  10,
			Timeout:   3600 * time.Second,
			Workspace: "/tmp/cicd-workspace",
		},
		Executor: ExecutorConfig{
			Type: "local",
			Env:  make(map[string]string),
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Runner.Capacity <= 0 {
		return fmt.Errorf("runner capacity must be greater than 0")
	}
	if c.Runner.Timeout <= 0 {
		return fmt.Errorf("runner timeout must be greater than 0")
	}
	if c.Executor.Type != "local" && c.Executor.Type != "mock" {
		return fmt.Errorf("executor type must be 'local' or 'mock'")
	}
	return nil
}

// GetEnv 获取环境变量，优先使用系统环境变量
func (c *Config) GetEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return c.Executor.Env[key]
}

// SetEnv 设置环境变量
func (c *Config) SetEnv(key, value string) {
	c.Executor.Env[key] = value
}
