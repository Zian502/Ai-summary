package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrStepNameRequired     = fmt.Errorf("step name is required")
	ErrStepCommandsRequired = fmt.Errorf("step commands are required")
	ErrPipelineNameRequired = fmt.Errorf("pipeline name is required")
)

// Pipeline 定义完整的 CI/CD 流水线
type Pipeline struct {
	Name        string            `yaml:"name"`        // Pipeline 名称
	Version     string            `yaml:"version"`     // Pipeline 版本
	Steps       []Step            `yaml:"steps"`       // 执行步骤列表
	Env         map[string]string `yaml:"env"`         // 全局环境变量
	Workspace   string            `yaml:"workspace"`   // 工作空间路径
	Concurrency int               `yaml:"concurrency"` // 并发执行数
}

// Load 从文件加载 Pipeline
func Load(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pipeline file: %w", err)
	}

	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse pipeline file: %w", err)
	}

	// 验证 Pipeline
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pipeline: %w", err)
	}

	return &p, nil
}

// Validate 验证 Pipeline 配置
func (p *Pipeline) Validate() error {
	if p.Name == "" {
		return ErrPipelineNameRequired
	}

	if len(p.Steps) == 0 {
		return fmt.Errorf("pipeline must have at least one step")
	}

	// 验证每个步骤
	for i, step := range p.Steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("step %d (%s): %w", i, step.Name, err)
		}
	}

	return nil
}

// GetStep 根据名称获取步骤
func (p *Pipeline) GetStep(name string) *Step {
	for i := range p.Steps {
		if p.Steps[i].Name == name {
			return &p.Steps[i]
		}
	}
	return nil
}

// GetEnv 获取环境变量，优先使用步骤特定的环境变量
func (p *Pipeline) GetEnv(step *Step, key string) string {
	// 先检查步骤特定的环境变量
	if step != nil && step.Env != nil {
		if val, ok := step.Env[key]; ok {
			return val
		}
	}
	// 再检查全局环境变量
	if p.Env != nil {
		if val, ok := p.Env[key]; ok {
			return val
		}
	}
	return ""
}
