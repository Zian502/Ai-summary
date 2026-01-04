package runner

import (
	"context"
	"fmt"
	"sync"

	"github.com/projects/cicd-runner/config"
	"github.com/projects/cicd-runner/executor"
	"github.com/projects/cicd-runner/pipeline"
)

// Runner CI/CD Runner 核心
type Runner struct {
	config   *config.Config
	executor executor.Executor
}

// New 创建新的 Runner
func New(cfg *config.Config) *Runner {
	exec := executor.NewExecutor(cfg.Executor.Type)
	return &Runner{
		config:   cfg,
		executor: exec,
	}
}

// Run 运行 Pipeline
func (r *Runner) Run(pipelinePath string) error {
	// 加载 Pipeline
	p, err := pipeline.Load(pipelinePath)
	if err != nil {
		return fmt.Errorf("failed to load pipeline: %w", err)
	}

	// 创建工作空间
	workspace := r.config.Runner.Workspace
	if p.Workspace != "" {
		workspace = p.Workspace
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Runner.Timeout)
	defer cancel()

	// 设置执行环境
	if err := r.executor.Setup(ctx, workspace); err != nil {
		return fmt.Errorf("failed to setup executor: %w", err)
	}
	defer r.executor.Teardown(ctx, workspace)

	// 准备环境变量
	env := r.prepareEnv(p)

	// 执行步骤
	results, err := r.executeSteps(ctx, p, env, workspace)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	// 打印结果
	r.printResults(results)

	// 检查是否有失败的步骤
	for _, result := range results {
		if !result.Success {
			return fmt.Errorf("pipeline failed at step: %s", result.Step.Name)
		}
	}

	return nil
}

// executeSteps 执行所有步骤
func (r *Runner) executeSteps(ctx context.Context, p *pipeline.Pipeline, env map[string]string, workspace string) ([]*executor.Result, error) {
	concurrency := p.Concurrency
	if concurrency <= 0 {
		concurrency = 1 // 默认串行执行
	}
	if concurrency > r.config.Runner.Capacity {
		concurrency = r.config.Runner.Capacity
	}

	var results []*executor.Result
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用信号量控制并发
	sem := make(chan struct{}, concurrency)

	for i := range p.Steps {
		step := &p.Steps[i]

		// 检查步骤是否应该执行
		if !step.ShouldRun() {
			continue
		}

		wg.Add(1)
		go func(s *pipeline.Step) {
			defer wg.Done()

			// 获取信号量
			sem <- struct{}{}
			defer func() { <-sem }()

			// 准备步骤特定的环境变量
			stepEnv := r.prepareStepEnv(p, s, env)

			// 执行步骤
			result, err := r.executor.Execute(ctx, s, stepEnv, workspace)
			if err != nil {
				result = &executor.Result{
					Success:  false,
					ExitCode: 1,
					Error:    err.Error(),
					Step:     s,
				}
			}

			// 保存结果
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(step)
	}

	wg.Wait()
	return results, nil
}

// prepareEnv 准备环境变量
func (r *Runner) prepareEnv(p *pipeline.Pipeline) map[string]string {
	env := make(map[string]string)

	// 从配置中获取环境变量
	for k, v := range r.config.Executor.Env {
		env[k] = v
	}

	// 从 Pipeline 中获取全局环境变量
	for k, v := range p.Env {
		env[k] = v
	}

	return env
}

// prepareStepEnv 准备步骤特定的环境变量
func (r *Runner) prepareStepEnv(p *pipeline.Pipeline, step *pipeline.Step, baseEnv map[string]string) map[string]string {
	env := make(map[string]string)

	// 复制基础环境变量
	for k, v := range baseEnv {
		env[k] = v
	}

	// 添加步骤特定的环境变量
	if step.Env != nil {
		for k, v := range step.Env {
			env[k] = v
		}
	}

	return env
}

// printResults 打印执行结果
func (r *Runner) printResults(results []*executor.Result) {
	fmt.Println("\n=== Pipeline Execution Results ===")
	for i, result := range results {
		status := "✓ SUCCESS"
		if !result.Success {
			status = "✗ FAILED"
		}

		fmt.Printf("\n[%d] Step: %s\n", i+1, result.Step.Name)
		fmt.Printf("  Status: %s\n", status)
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Exit Code: %d\n", result.ExitCode)

		if result.Output != "" {
			fmt.Printf("  Output:\n%s\n", indent(result.Output, "    "))
		}

		if result.Error != "" {
			fmt.Printf("  Error: %s\n", result.Error)
		}
	}
	fmt.Println("\n===================================")
}

// indent 缩进文本
func indent(text, prefix string) string {
	lines := []rune(text)
	var result []rune
	var currentLine []rune

	for _, char := range lines {
		if char == '\n' {
			if len(currentLine) > 0 {
				result = append(result, []rune(prefix)...)
				result = append(result, currentLine...)
				currentLine = []rune{}
			}
			result = append(result, '\n')
		} else {
			currentLine = append(currentLine, char)
		}
	}

	if len(currentLine) > 0 {
		result = append(result, []rune(prefix)...)
		result = append(result, currentLine...)
	}

	return string(result)
}
