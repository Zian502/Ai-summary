package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/projects/cicd-runner/config"
	"github.com/projects/cicd-runner/runner"
)

var (
	configPath   = flag.String("config", "", "配置文件路径（可选）")
	pipelinePath = flag.String("pipeline", "examples/pipeline.yaml", "Pipeline 配置文件路径")
	mockMode     = flag.Bool("mock", false, "使用 Mock 模式（不实际执行命令）")
	version      = flag.Bool("version", false, "显示版本信息")
)

const (
	Version = "0.1.0"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("CI/CD Runner v%s\n", Version)
		return
	}

	// 加载配置
	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 使用默认配置或从环境变量加载
		cfg = config.LoadFromEnv()
	}

	// 如果指定了 mock 模式，覆盖配置
	if *mockMode {
		cfg.Executor.Type = "mock"
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		os.Exit(1)
	}

	// 检查 Pipeline 文件是否存在
	if _, err := os.Stat(*pipelinePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Pipeline file not found: %s\n", *pipelinePath)
		fmt.Fprintf(os.Stderr, "Usage: %s -pipeline <path> [-config <path>] [-mock]\n", os.Args[0])
		os.Exit(1)
	}

	// 创建并运行 Runner
	r := runner.New(cfg)
	if err := r.Run(*pipelinePath); err != nil {
		fmt.Fprintf(os.Stderr, "Pipeline execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Pipeline completed successfully!")
}
