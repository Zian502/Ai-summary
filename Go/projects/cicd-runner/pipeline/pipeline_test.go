package pipeline

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepValidate(t *testing.T) {
	tests := []struct {
		name    string
		step    Step
		wantErr bool
	}{
		{
			name: "valid step",
			step: Step{
				Name:     "test",
				Commands: []string{"echo hello"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			step: Step{
				Commands: []string{"echo hello"},
			},
			wantErr: true,
		},
		{
			name: "missing commands",
			step: Step{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "empty commands",
			step: Step{
				Name:     "test",
				Commands: []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.step.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStepShouldRun(t *testing.T) {
	tests := []struct {
		name string
		step Step
		want bool
	}{
		{
			name: "always run",
			step: Step{
				When: "always",
			},
			want: true,
		},
		{
			name: "empty when",
			step: Step{
				When: "",
			},
			want: true,
		},
		{
			name: "on_success",
			step: Step{
				When: "on_success",
			},
			want: true,
		},
		{
			name: "on_failure",
			step: Step{
				When: "on_failure",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.step.ShouldRun())
		})
	}
}

func TestPipelineValidate(t *testing.T) {
	tests := []struct {
		name     string
		pipeline Pipeline
		wantErr  bool
	}{
		{
			name: "valid pipeline",
			pipeline: Pipeline{
				Name: "test-pipeline",
				Steps: []Step{
					{
						Name:     "step1",
						Commands: []string{"echo hello"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			pipeline: Pipeline{
				Steps: []Step{
					{
						Name:     "step1",
						Commands: []string{"echo hello"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no steps",
			pipeline: Pipeline{
				Name: "test-pipeline",
			},
			wantErr: true,
		},
		{
			name: "invalid step",
			pipeline: Pipeline{
				Name: "test-pipeline",
				Steps: []Step{
					{
						// Missing name and commands
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pipeline.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPipelineGetStep(t *testing.T) {
	p := &Pipeline{
		Name: "test",
		Steps: []Step{
			{Name: "step1", Commands: []string{"cmd1"}},
			{Name: "step2", Commands: []string{"cmd2"}},
		},
	}

	step := p.GetStep("step1")
	assert.NotNil(t, step)
	assert.Equal(t, "step1", step.Name)

	step = p.GetStep("nonexistent")
	assert.Nil(t, step)
}

func TestPipelineGetEnv(t *testing.T) {
	p := &Pipeline{
		Env: map[string]string{
			"GLOBAL_VAR": "global_value",
		},
		Steps: []Step{
			{
				Name: "step1",
				Env: map[string]string{
					"STEP_VAR": "step_value",
				},
				Commands: []string{"echo"},
			},
		},
	}

	step := &p.Steps[0]

	// 步骤特定的环境变量优先
	assert.Equal(t, "step_value", p.GetEnv(step, "STEP_VAR"))

	// 全局环境变量
	assert.Equal(t, "global_value", p.GetEnv(step, "GLOBAL_VAR"))

	// 不存在的变量
	assert.Equal(t, "", p.GetEnv(step, "NON_EXISTENT"))

	// 没有步骤时
	assert.Equal(t, "global_value", p.GetEnv(nil, "GLOBAL_VAR"))
}

func TestLoad(t *testing.T) {
	// 创建临时 Pipeline 文件
	tmpFile, err := os.CreateTemp("", "test-pipeline-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	pipelineContent := `
name: test-pipeline
version: "1.0"
workspace: /tmp/test

env:
  TEST_VAR: "test_value"

steps:
  - name: step1
    commands:
      - echo "Hello"
      - echo "World"
    timeout: 60
`

	_, err = tmpFile.WriteString(pipelineContent)
	require.NoError(t, err)
	tmpFile.Close()

	p, err := Load(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, "test-pipeline", p.Name)
	assert.Equal(t, "1.0", p.Version)
	assert.Equal(t, "/tmp/test", p.Workspace)
	assert.Equal(t, "test_value", p.Env["TEST_VAR"])
	assert.Len(t, p.Steps, 1)
	assert.Equal(t, "step1", p.Steps[0].Name)
	assert.Len(t, p.Steps[0].Commands, 2)
}
