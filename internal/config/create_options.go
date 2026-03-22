package config

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ProjectTypes      = []string{"web-app", "ai-app", "devops-tool", "internal-tool", "platform-service"}
	AIFeatureTypes    = []string{"none", "prompt-workflow", "rag", "agent-system"}
	BackendTypes      = []string{"go", "python", "node", "none"}
	FrontendTypes     = []string{"react", "next", "nuxt", "vue", "pure-typescript", "none"}
	ArchitectureTypes = []string{"cli-tool", "backend-service", "frontend-app", "fullstack-web-app", "frontend-backend"}
	DocsTypes         = []string{"basic", "full"}
	YesNoTypes        = []string{"yes", "no"}
	AgentTypes        = []string{"codex", "claude-code"}
)

type CreateOptions struct {
	Name             string
	ParentPath       string
	ProjectType      string
	AIFeature        string
	AIAgent          string
	BackendType      string
	FrontendType     string
	DocsType         string
	GlobalSkillsPath string
	SelectedSkills   []string
	PromptTitle      string
	Description      string
	Architecture     string
	TechStack        string
	UnitTest         string
	IntegrationTest  string
	E2ETest          string
	PerformanceTest  string
	SecurityTest     string
	UATTest          string
	DockerCompose    string
	DryRun           bool
	Overwrite        bool
}

func ParseSkillsCSV(raw string) []string {
	items := strings.Split(raw, ",")
	result := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func (o *CreateOptions) Normalize() {
	o.Name = strings.TrimSpace(o.Name)
	o.ParentPath = strings.TrimSpace(o.ParentPath)
	o.ProjectType = strings.ToLower(strings.TrimSpace(o.ProjectType))
	o.AIFeature = strings.ToLower(strings.TrimSpace(o.AIFeature))
	switch o.AIFeature {
	case "prompt_workflow":
		o.AIFeature = "prompt-workflow"
	case "agent_system":
		o.AIFeature = "agent-system"
	}
	o.AIAgent = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(o.AIAgent)), " ", "-")
	o.BackendType = strings.ToLower(strings.TrimSpace(o.BackendType))
	o.FrontendType = strings.ToLower(strings.TrimSpace(o.FrontendType))
	switch o.FrontendType {
	case "typescript", "pure_ts", "pure-typescript", "pure typescript", "pure-ts":
		o.FrontendType = "pure-typescript"
	case "nuxtjs":
		o.FrontendType = "nuxt"
	}
	o.DocsType = strings.ToLower(strings.TrimSpace(o.DocsType))
	o.GlobalSkillsPath = strings.TrimSpace(o.GlobalSkillsPath)
	o.PromptTitle = strings.TrimSpace(o.PromptTitle)
	o.Description = strings.TrimSpace(o.Description)
	o.Architecture = strings.ToLower(strings.TrimSpace(o.Architecture))
	switch o.Architecture {
	case "microservices", "microservice", "service", "backend":
		o.Architecture = "backend-service"
	case "cli":
		o.Architecture = "cli-tool"
	case "frontend":
		o.Architecture = "frontend-app"
	case "fullstack":
		o.Architecture = "fullstack-web-app"
	case "frontend+backend", "frontend_backend", "frontend-backend", "fe-be":
		o.Architecture = "frontend-backend"
	}
	o.TechStack = strings.TrimSpace(o.TechStack)
	o.UnitTest = strings.ToLower(strings.TrimSpace(o.UnitTest))
	o.IntegrationTest = strings.ToLower(strings.TrimSpace(o.IntegrationTest))
	o.E2ETest = strings.ToLower(strings.TrimSpace(o.E2ETest))
	o.PerformanceTest = strings.ToLower(strings.TrimSpace(o.PerformanceTest))
	o.SecurityTest = strings.ToLower(strings.TrimSpace(o.SecurityTest))
	o.UATTest = strings.ToLower(strings.TrimSpace(o.UATTest))
	o.DockerCompose = strings.ToLower(strings.TrimSpace(o.DockerCompose))
	o.SelectedSkills = ParseSkillsCSV(strings.Join(o.SelectedSkills, ","))
}

func (o *CreateOptions) EnsureDefaults() {
	if o.ParentPath == "" {
		o.ParentPath = "."
	}
	if o.DocsType == "" {
		o.DocsType = "basic"
	}
	if o.AIFeature == "" {
		o.AIFeature = "none"
	}
	if o.AIAgent == "" {
		o.AIAgent = "codex"
	}
	if o.PromptTitle == "" {
		if o.Name != "" {
			o.PromptTitle = fmt.Sprintf("%s Prompt", o.Name)
		} else {
			o.PromptTitle = "Project Prompt"
		}
	}
}

func (o CreateOptions) MissingRequired() []string {
	missing := make([]string, 0, 7)
	if strings.TrimSpace(o.Name) == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(o.ProjectType) == "" {
		missing = append(missing, "type")
	}
	if strings.TrimSpace(o.Architecture) == "" {
		missing = append(missing, "architecture")
	}

	// backend / frontend requirements depend on architecture
	switch o.Architecture {
	case "cli-tool", "backend-service":
		if strings.TrimSpace(o.BackendType) == "" {
			missing = append(missing, "backend")
		}
	case "frontend-app", "fullstack-web-app":
		if strings.TrimSpace(o.FrontendType) == "" {
			missing = append(missing, "frontend")
		}
	case "frontend-backend":
		if strings.TrimSpace(o.BackendType) == "" {
			missing = append(missing, "backend")
		}
		if strings.TrimSpace(o.FrontendType) == "" {
			missing = append(missing, "frontend")
		}
	default:
		// architecture unknown or empty: require both
		if strings.TrimSpace(o.BackendType) == "" {
			missing = append(missing, "backend")
		}
		if strings.TrimSpace(o.FrontendType) == "" {
			missing = append(missing, "frontend")
		}
	}

	if strings.TrimSpace(o.TechStack) == "" {
		missing = append(missing, "stack")
	}
	if strings.TrimSpace(o.DockerCompose) == "" {
		missing = append(missing, "docker-compose")
	}
	return missing
}

func (o CreateOptions) AnyInputProvided() bool {
	if o.Name != "" || o.ParentPath != "" || o.ProjectType != "" || o.AIFeature != "" || o.AIAgent != "" || o.BackendType != "" || o.FrontendType != "" {
		return true
	}
	if o.DocsType != "" || o.GlobalSkillsPath != "" || len(o.SelectedSkills) > 0 || o.PromptTitle != "" || o.Description != "" {
		return true
	}
	if o.Architecture != "" || o.TechStack != "" || o.UnitTest != "" || o.IntegrationTest != "" || o.E2ETest != "" {
		return true
	}
	if o.PerformanceTest != "" || o.SecurityTest != "" || o.UATTest != "" || o.DockerCompose != "" {
		return true
	}
	if o.Overwrite {
		return true
	}
	return false
}

func (o CreateOptions) ValidateKnownValues() error {
	if o.ProjectType != "" && !slices.Contains(ProjectTypes, o.ProjectType) {
		return fmt.Errorf("invalid --type %q, valid values: %s", o.ProjectType, strings.Join(ProjectTypes, ", "))
	}
	if o.AIFeature != "" && !slices.Contains(AIFeatureTypes, o.AIFeature) {
		return fmt.Errorf("invalid --ai-feature %q, valid values: %s", o.AIFeature, strings.Join(AIFeatureTypes, ", "))
	}
	if o.AIAgent != "" && !slices.Contains(AgentTypes, o.AIAgent) {
		return fmt.Errorf("invalid --agent %q, valid values: %s", o.AIAgent, strings.Join(AgentTypes, ", "))
	}
	if o.BackendType != "" && !slices.Contains(BackendTypes, o.BackendType) {
		return fmt.Errorf("invalid --backend %q, valid values: %s", o.BackendType, strings.Join(BackendTypes, ", "))
	}
	if o.FrontendType != "" && !slices.Contains(FrontendTypes, o.FrontendType) {
		return fmt.Errorf("invalid --frontend %q, valid values: %s", o.FrontendType, strings.Join(FrontendTypes, ", "))
	}
	if o.Architecture != "" && !slices.Contains(ArchitectureTypes, o.Architecture) {
		return fmt.Errorf("invalid --architecture %q, valid values: %s", o.Architecture, strings.Join(ArchitectureTypes, ", "))
	}
	if o.DocsType != "" && !slices.Contains(DocsTypes, o.DocsType) {
		return fmt.Errorf("invalid --docs %q, valid values: %s", o.DocsType, strings.Join(DocsTypes, ", "))
	}
	if err := validateYesNoOption("--unit-test", o.UnitTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--integration-test", o.IntegrationTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--e2e-test", o.E2ETest); err != nil {
		return err
	}
	if err := validateYesNoOption("--performance-test", o.PerformanceTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--security-test", o.SecurityTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--uat-test", o.UATTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--docker-compose", o.DockerCompose); err != nil {
		return err
	}
	return nil
}

func (o CreateOptions) ValidateForCreate() error {
	if len(o.MissingRequired()) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(o.MissingRequired(), ", "))
	}
	if err := o.ValidateKnownValues(); err != nil {
		return err
	}
	if strings.Contains(o.Name, "/") || strings.Contains(o.Name, "\\") {
		return fmt.Errorf("project name must not contain path separators: %q", o.Name)
	}
	if o.Name == "." || o.Name == ".." {
		return fmt.Errorf("invalid project name: %q", o.Name)
	}
	if strings.TrimSpace(o.ParentPath) == "" {
		return fmt.Errorf("parent path cannot be empty")
	}

	cleanBase := filepath.Base(o.Name)
	if cleanBase != o.Name {
		return fmt.Errorf("invalid project name: %q", o.Name)
	}

	return nil
}

func validateYesNoOption(flagName, value string) error {
	if value == "" {
		return nil
	}
	if slices.Contains(YesNoTypes, value) {
		return nil
	}
	return fmt.Errorf("invalid %s %q, valid values: %s", flagName, value, strings.Join(YesNoTypes, ", "))
}
