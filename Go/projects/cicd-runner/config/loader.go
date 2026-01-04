package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 应用默认值
	applyDefaults(&cfg)

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	cfg := DefaultConfig()

	// 从环境变量覆盖配置
	if val := os.Getenv("CICD_RUNNER_CAPACITY"); val != "" {
		fmt.Sscanf(val, "%d", &cfg.Runner.Capacity)
	}
	if val := os.Getenv("CICD_RUNNER_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			cfg.Runner.Timeout = timeout
		}
	}
	if val := os.Getenv("CICD_RUNNER_WORKSPACE"); val != "" {
		cfg.Runner.Workspace = val
	}
	if val := os.Getenv("CICD_EXECUTOR_TYPE"); val != "" {
		cfg.Executor.Type = val
	}
	if val := os.Getenv("CICD_LOG_LEVEL"); val != "" {
		cfg.Log.Level = val
	}

	return cfg
}

// applyDefaults 应用默认值
func applyDefaults(cfg *Config) {
	if cfg.Runner.Capacity == 0 {
		cfg.Runner.Capacity = 10
	}
	if cfg.Runner.Timeout == 0 {
		cfg.Runner.Timeout = 3600 * time.Second
	}
	if cfg.Runner.Workspace == "" {
		cfg.Runner.Workspace = "/tmp/cicd-workspace"
	}
	if cfg.Executor.Type == "" {
		cfg.Executor.Type = "local"
	}
	if cfg.Executor.Env == nil {
		cfg.Executor.Env = make(map[string]string)
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "text"
	}
}
