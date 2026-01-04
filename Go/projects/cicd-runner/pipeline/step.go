package pipeline

// Step 定义单个执行步骤
type Step struct {
	Name      string            `yaml:"name"`       // 步骤名称
	Image     string            `yaml:"image"`      // 使用的镜像（可选，用于容器化执行）
	Commands  []string          `yaml:"commands"`   // 执行的命令列表
	Env       map[string]string `yaml:"env"`        // 步骤特定的环境变量
	When      string            `yaml:"when"`       // 执行条件（可选）
	Timeout   int               `yaml:"timeout"`    // 超时时间（秒）
	OnSuccess []string          `yaml:"on_success"` // 成功时执行的命令
	OnFailure []string          `yaml:"on_failure"` // 失败时执行的命令
}

// ShouldRun 判断步骤是否应该执行
func (s *Step) ShouldRun() bool {
	// 简化实现：when 字段为空或为 "always" 时执行
	if s.When == "" || s.When == "always" {
		return true
	}
	// 可以扩展支持更多条件，如 "on_success", "on_failure" 等
	return s.When == "on_success" || s.When == "on_failure"
}

// Validate 验证步骤配置
func (s *Step) Validate() error {
	if s.Name == "" {
		return ErrStepNameRequired
	}
	if len(s.Commands) == 0 {
		return ErrStepCommandsRequired
	}
	return nil
}
