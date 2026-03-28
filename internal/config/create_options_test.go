package config

import (
	"strings"
	"testing"
)

func TestParseSkillsCSV(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{"empty", "", nil},
		{"single", "foo", []string{"foo"}},
		{"multiple", "foo,bar,baz", []string{"foo", "bar", "baz"}},
		{"with spaces", " foo , bar , baz ", []string{"foo", "bar", "baz"}},
		{"dedup", "foo,bar,foo,baz,bar", []string{"foo", "bar", "baz"}},
		{"empty items", "foo,,bar,,", []string{"foo", "bar"}},
		{"whitespace only items", "foo, , bar, ,", []string{"foo", "bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSkillsCSV(tt.raw)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ParseSkillsCSV(%q) = %v, want %v", tt.raw, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ParseSkillsCSV(%q)[%d] = %q, want %q", tt.raw, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		opts  CreateOptions
		check func(t *testing.T, o CreateOptions)
	}{
		{
			name: "trims whitespace",
			opts: CreateOptions{Name: "  my-project  ", ParentPath: " /tmp "},
			check: func(t *testing.T, o CreateOptions) {
				if o.Name != "my-project" {
					t.Errorf("Name = %q", o.Name)
				}
				if o.ParentPath != "/tmp" {
					t.Errorf("ParentPath = %q", o.ParentPath)
				}
			},
		},
		{
			name: "lowercases type fields",
			opts: CreateOptions{ProjectType: " Web-App ", BackendType: "GO", FrontendType: "REACT"},
			check: func(t *testing.T, o CreateOptions) {
				if o.ProjectType != "full-stack" {
					t.Errorf("ProjectType = %q", o.ProjectType)
				}
				if o.BackendType != "go" {
					t.Errorf("BackendType = %q", o.BackendType)
				}
				if o.FrontendType != "react" {
					t.Errorf("FrontendType = %q", o.FrontendType)
				}
			},
		},
		{
			name: "normalizes ai feature underscores",
			opts: CreateOptions{AIFeature: "prompt_workflow"},
			check: func(t *testing.T, o CreateOptions) {
				if o.AIFeature != "prompt-workflow" {
					t.Errorf("AIFeature = %q", o.AIFeature)
				}
			},
		},
		{
			name: "normalizes agent_system",
			opts: CreateOptions{AIFeature: "agent_system"},
			check: func(t *testing.T, o CreateOptions) {
				if o.AIFeature != "agent-system" {
					t.Errorf("AIFeature = %q", o.AIFeature)
				}
			},
		},
		{
			name: "normalizes frontend aliases",
			opts: CreateOptions{FrontendType: "typescript"},
			check: func(t *testing.T, o CreateOptions) {
				if o.FrontendType != "pure-typescript" {
					t.Errorf("FrontendType = %q", o.FrontendType)
				}
			},
		},
		{
			name: "normalizes nuxtjs",
			opts: CreateOptions{FrontendType: "nuxtjs"},
			check: func(t *testing.T, o CreateOptions) {
				if o.FrontendType != "nuxt" {
					t.Errorf("FrontendType = %q", o.FrontendType)
				}
			},
		},
		{
			name: "normalizes architecture aliases",
			opts: CreateOptions{Architecture: "microservices"},
			check: func(t *testing.T, o CreateOptions) {
				if o.Architecture != "server" {
					t.Errorf("Architecture = %q", o.Architecture)
				}
			},
		},
		{
			name: "normalizes architecture cli",
			opts: CreateOptions{Architecture: "cli"},
			check: func(t *testing.T, o CreateOptions) {
				if o.Architecture != "cli" {
					t.Errorf("Architecture = %q", o.Architecture)
				}
			},
		},
		{
			name: "normalizes architecture fe-be",
			opts: CreateOptions{Architecture: "fe-be"},
			check: func(t *testing.T, o CreateOptions) {
				if o.Architecture != "web-app-server" {
					t.Errorf("Architecture = %q", o.Architecture)
				}
			},
		},
		{
			name: "normalizes agent spaces to dashes",
			opts: CreateOptions{AIAgent: " Claude Code "},
			check: func(t *testing.T, o CreateOptions) {
				if o.AIAgent != "claude-code" {
					t.Errorf("AIAgent = %q", o.AIAgent)
				}
			},
		},
		{
			name: "normalizes universal aliases",
			opts: CreateOptions{AIAgent: " both "},
			check: func(t *testing.T, o CreateOptions) {
				if o.AIAgent != "universal" {
					t.Errorf("AIAgent = %q", o.AIAgent)
				}
			},
		},
		{
			name: "deduplicates skills",
			opts: CreateOptions{SelectedSkills: []string{"a,b,a,c"}},
			check: func(t *testing.T, o CreateOptions) {
				if len(o.SelectedSkills) != 3 {
					t.Errorf("SelectedSkills = %v", o.SelectedSkills)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.Normalize()
			tt.check(t, tt.opts)
		})
	}
}

func TestEnsureDefaults(t *testing.T) {
	o := CreateOptions{}
	o.EnsureDefaults()

	if o.ParentPath != "." {
		t.Errorf("ParentPath = %q, want %q", o.ParentPath, ".")
	}
	if o.DocsType != "basic" {
		t.Errorf("DocsType = %q, want %q", o.DocsType, "basic")
	}
	if o.AIFeature != "none" {
		t.Errorf("AIFeature = %q, want %q", o.AIFeature, "none")
	}
	if o.ProjectType != "full-stack" {
		t.Errorf("ProjectType = %q, want %q", o.ProjectType, "full-stack")
	}
	if o.Architecture != "web-app-server" {
		t.Errorf("Architecture = %q, want %q", o.Architecture, "web-app-server")
	}
	if o.AIAgent != "codex" {
		t.Errorf("AIAgent = %q, want %q", o.AIAgent, "codex")
	}
	if o.PromptTitle != "Project Prompt" {
		t.Errorf("PromptTitle = %q, want %q", o.PromptTitle, "Project Prompt")
	}
}

func TestEnsureDefaultsWithName(t *testing.T) {
	o := CreateOptions{Name: "my-app"}
	o.EnsureDefaults()

	if o.PromptTitle != "my-app Prompt" {
		t.Errorf("PromptTitle = %q, want %q", o.PromptTitle, "my-app Prompt")
	}
}

func TestMissingRequired(t *testing.T) {
	tests := []struct {
		name    string
		opts    CreateOptions
		wantLen int
	}{
		{
			name:    "all missing",
			opts:    CreateOptions{},
			wantLen: 2,
		},
		{
			name: "none missing",
			opts: CreateOptions{
				Name: "x", ProjectType: "full-stack", BackendType: "go",
				FrontendType: "react", Architecture: "web-app-server",
				TechStack: "go+react", DockerCompose: "yes",
			},
			wantLen: 0,
		},
		{
			name:    "whitespace only counts as missing",
			opts:    CreateOptions{Name: "  ", ProjectType: "web-app"},
			wantLen: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.MissingRequired()
			if len(got) != tt.wantLen {
				t.Errorf("MissingRequired() = %v (len %d), want len %d", got, len(got), tt.wantLen)
			}
		})
	}
}

func TestValidateKnownValues(t *testing.T) {
	tests := []struct {
		name    string
		opts    CreateOptions
		wantErr string
	}{
		{"valid empty", CreateOptions{}, ""},
		{"valid project type", CreateOptions{ProjectType: "full-stack"}, ""},
		{"invalid project type", CreateOptions{ProjectType: "invalid"}, "invalid --type"},
		{"invalid ai feature", CreateOptions{AIFeature: "bad"}, "invalid --ai-feature"},
		{"invalid agent", CreateOptions{AIAgent: "bad"}, "invalid --agent"},
		{"invalid backend", CreateOptions{BackendType: "rust"}, "invalid --backend"},
		{"invalid frontend", CreateOptions{FrontendType: "svelte"}, "invalid --frontend"},
		{"invalid architecture", CreateOptions{Architecture: "bad"}, "invalid --type"},
		{"invalid docs", CreateOptions{DocsType: "bad"}, "invalid --docs"},
		{"invalid unit test", CreateOptions{UnitTest: "maybe"}, "invalid --unit-test"},
		{"invalid api test", CreateOptions{APITest: "maybe"}, "invalid --api-test"},
		{"valid yes/no", CreateOptions{UnitTest: "yes", APITest: "yes", E2ETest: "no"}, ""},
		{"invalid docker compose", CreateOptions{DockerCompose: "bad"}, "invalid --docker-compose"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.ValidateKnownValues()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestValidateForCreate(t *testing.T) {
	valid := CreateOptions{
		Name: "my-project", ProjectType: "full-stack", BackendType: "go",
		FrontendType: "react", Architecture: "web-app-server",
		TechStack: "go+react", ParentPath: "/tmp", DockerCompose: "yes",
	}

	if err := valid.ValidateForCreate(); err != nil {
		t.Errorf("valid opts got error: %v", err)
	}

	tests := []struct {
		name    string
		modify  func(o *CreateOptions)
		wantErr string
	}{
		{"name with slash", func(o *CreateOptions) { o.Name = "a/b" }, "path separators"},
		{"name is dot", func(o *CreateOptions) { o.Name = "." }, "invalid project name"},
		{"name is dotdot", func(o *CreateOptions) { o.Name = ".." }, "invalid project name"},
		{"empty parent", func(o *CreateOptions) { o.ParentPath = "" }, "parent path cannot be empty"},
		{"missing fields", func(o *CreateOptions) { o.Name = "" }, "missing required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := valid
			tt.modify(&o)
			err := o.ValidateForCreate()
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestAnyInputProvided(t *testing.T) {
	if (CreateOptions{}).AnyInputProvided() {
		t.Error("empty opts should return false")
	}
	if !(CreateOptions{Name: "x"}).AnyInputProvided() {
		t.Error("Name set should return true")
	}
	if !(CreateOptions{Overwrite: true}).AnyInputProvided() {
		t.Error("Overwrite set should return true")
	}
	if !(CreateOptions{UATTest: "yes"}).AnyInputProvided() {
		t.Error("UATTest set should return true")
	}
	if !(CreateOptions{APITest: "yes"}).AnyInputProvided() {
		t.Error("APITest set should return true")
	}
}
