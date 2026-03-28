package config

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ProjectTypes      = []string{"frontend-only", "backend-only", "full-stack"}
	AIFeatureTypes    = []string{"none", "prompt-workflow", "rag", "agent-system"}
	BackendTypes      = []string{"go", "python", "node", "none"}
	FrontendTypes     = []string{"react", "next", "nuxt", "vue", "pure-typescript", "none"}
	ArchitectureTypes = []string{"cli", "server", "web-app-server", "mobile-app-server", "web-app", "mobile-app"}
	DocsTypes         = []string{"basic", "full"}
	YesNoTypes        = []string{"yes", "no"}
	AgentTypes        = []string{"codex", "claude-code", "universal"}
)

type CreateOptions struct {
	Name             string
	Idea             string
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
	APITest          string
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
	o.Idea = strings.TrimSpace(o.Idea)
	o.ParentPath = strings.TrimSpace(o.ParentPath)
	o.ProjectType = strings.ToLower(strings.TrimSpace(o.ProjectType))
	switch o.ProjectType {
	case "frontend-only", "frontend_only", "frontend only", "frontend":
		o.ProjectType = "frontend-only"
	case "backend-only", "backend_only", "backend only", "backend":
		o.ProjectType = "backend-only"
	case "full-stack", "full_stack", "full stack", "fullstack":
		o.ProjectType = "full-stack"
	// Backward-compatible aliases from previous versions.
	case "web-app", "ai-app", "devops-tool", "platform-service":
		o.ProjectType = "full-stack"
	case "internal-tool":
		o.ProjectType = "backend-only"
	}
	o.AIFeature = strings.ToLower(strings.TrimSpace(o.AIFeature))
	switch o.AIFeature {
	case "prompt_workflow":
		o.AIFeature = "prompt-workflow"
	case "agent_system":
		o.AIFeature = "agent-system"
	}
	o.AIAgent = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(o.AIAgent)), " ", "-")
	switch o.AIAgent {
	case "claude", "claudecode":
		o.AIAgent = "claude-code"
	case "both", "all":
		o.AIAgent = "universal"
	}
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
	case "cli", "cli-tool", "cli-only":
		o.Architecture = "cli"
	case "microservices", "microservice", "service", "backend", "backend-service", "server", "server-only":
		o.Architecture = "server"
	case "frontend", "frontend-app", "web-app", "web":
		o.Architecture = "web-app"
	case "mobile", "mobile-app", "mobile-app-only":
		o.Architecture = "mobile-app"
	case "fullstack", "fullstack-web-app", "frontend+backend", "frontend_backend", "frontend-backend", "fe-be":
		o.Architecture = "web-app-server"
	}
	if o.ProjectType == "" && o.Architecture != "" {
		o.ProjectType = InferProjectTypeFromArchitecture(o.Architecture)
	}
	if o.Architecture == "" && o.ProjectType != "" {
		o.Architecture = InferArchitectureFromProjectType(o.ProjectType)
	}
	o.TechStack = strings.TrimSpace(o.TechStack)
	o.UnitTest = strings.ToLower(strings.TrimSpace(o.UnitTest))
	o.APITest = strings.ToLower(strings.TrimSpace(o.APITest))
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
	if o.ProjectType == "" {
		o.ProjectType = "full-stack"
	}
	if o.Architecture == "" {
		o.Architecture = InferArchitectureFromProjectType(o.ProjectType)
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
	// backend / frontend requirements depend on architecture
	architecture := strings.TrimSpace(o.Architecture)
	if architecture == "" {
		architecture = InferArchitectureFromProjectType(o.ProjectType)
	}
	if strings.TrimSpace(architecture) == "" {
		missing = append(missing, "architecture")
		return missing
	}
	switch architecture {
	case "cli", "server":
		if strings.TrimSpace(o.BackendType) == "" {
			missing = append(missing, "backend")
		}
	case "web-app", "mobile-app":
		if strings.TrimSpace(o.FrontendType) == "" {
			missing = append(missing, "frontend")
		}
	case "web-app-server", "mobile-app-server":
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

func InferArchitectureFromProjectType(projectType string) string {
	switch strings.ToLower(strings.TrimSpace(projectType)) {
	case "frontend-only", "frontend_only", "frontend only", "frontend":
		return "web-app"
	case "backend-only", "backend_only", "backend only", "backend", "internal-tool":
		return "server"
	case "full-stack", "full_stack", "full stack", "fullstack", "web-app", "ai-app", "devops-tool", "platform-service":
		return "web-app-server"
	default:
		return ""
	}
}

func InferProjectTypeFromArchitecture(architecture string) string {
	switch strings.ToLower(strings.TrimSpace(architecture)) {
	case "web-app", "mobile-app", "web-app-only", "mobile-app-only", "frontend-app":
		return "frontend-only"
	case "cli", "server", "cli-only", "server-only", "cli-tool", "backend-service":
		return "backend-only"
	case "web-app-server", "mobile-app-server", "fullstack-web-app", "frontend-backend":
		return "full-stack"
	default:
		return ""
	}
}

func (o CreateOptions) AnyInputProvided() bool {
	if o.Name != "" || o.Idea != "" || o.ParentPath != "" || o.ProjectType != "" || o.AIFeature != "" || o.AIAgent != "" || o.BackendType != "" || o.FrontendType != "" {
		return true
	}
	if o.DocsType != "" || o.GlobalSkillsPath != "" || len(o.SelectedSkills) > 0 || o.PromptTitle != "" || o.Description != "" {
		return true
	}
	if o.Architecture != "" || o.TechStack != "" || o.UnitTest != "" || o.APITest != "" || o.IntegrationTest != "" || o.E2ETest != "" {
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
		return fmt.Errorf("invalid --type %q, valid values: %s", o.Architecture, strings.Join(ArchitectureTypes, ", "))
	}
	if o.DocsType != "" && !slices.Contains(DocsTypes, o.DocsType) {
		return fmt.Errorf("invalid --docs %q, valid values: %s", o.DocsType, strings.Join(DocsTypes, ", "))
	}
	if err := validateYesNoOption("--unit-test", o.UnitTest); err != nil {
		return err
	}
	if err := validateYesNoOption("--api-test", o.APITest); err != nil {
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
