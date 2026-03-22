package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"project-generator/internal/config"
	"project-generator/internal/i18n"
)

var ErrAborted = errors.New("operation aborted")

var frontendFrameworkOptions = []string{"react", "vue", "pure-typescript"}
var frontendLanguageOptions = []string{"typescript", "javascript"}
var backendLanguageOptions = []string{"go", "python", "node"}
var fullstackFrameworkOptions = []string{"next", "nuxt"}
var databaseOptions = []string{"postgres", "mysql", "sqlite", "none"}
var deploymentTargetOptions = []string{"docker compose", "kubernetes", "vm/server", "serverless"}
var authOptions = []string{"none", "jwt", "oauth"}

func projectTypeDescriptions() map[string]string {
	return map[string]string{
		"web-app":          i18n.T("wizard.project-type.web-app"),
		"ai-app":           i18n.T("wizard.project-type.ai-app"),
		"devops-tool":      i18n.T("wizard.project-type.devops-tool"),
		"internal-tool":    i18n.T("wizard.project-type.internal-tool"),
		"platform-service": i18n.T("wizard.project-type.platform-service"),
	}
}

func aiFeatureDescriptions() map[string]string {
	return map[string]string{
		"none":            i18n.T("wizard.ai-feature.none"),
		"prompt-workflow": i18n.T("wizard.ai-feature.prompt-workflow"),
		"rag":             i18n.T("wizard.ai-feature.rag"),
		"agent-system":    i18n.T("wizard.ai-feature.agent-system"),
	}
}

func aiAgentDescriptions() map[string]string {
	return map[string]string{
		"codex":       i18n.T("wizard.ai-agent.codex"),
		"claude-code": i18n.T("wizard.ai-agent.claude-code"),
	}
}

func architectureDescriptions() map[string]string {
	return map[string]string{
		"cli-tool":          i18n.T("wizard.architecture.cli-tool"),
		"backend-service":   i18n.T("wizard.architecture.backend-service"),
		"frontend-app":      i18n.T("wizard.architecture.frontend-app"),
		"fullstack-web-app": i18n.T("wizard.architecture.fullstack-web-app"),
		"frontend-backend":  i18n.T("wizard.architecture.frontend-backend"),
	}
}

func fullstackFrameworkDescriptions() map[string]string {
	return map[string]string{
		"next": i18n.T("wizard.fullstack.next"),
		"nuxt": i18n.T("wizard.fullstack.nuxt"),
	}
}

func CollectCreateOptions(input config.CreateOptions, _ bool) (config.CreateOptions, error) {
	reader := bufio.NewReader(os.Stdin)
	opts := input
	if strings.TrimSpace(opts.ParentPath) == "" {
		opts.ParentPath = "."
	}
	if strings.TrimSpace(opts.AIAgent) == "" {
		opts.AIAgent = "codex"
	}

	fmt.Println("Project Setup Wizard")
	fmt.Println("====================")
	fmt.Println()

	value, err := askText(reader, "1) Project name? (example: ops-platform)", opts.Name, true)
	if err != nil {
		return opts, err
	}
	opts.Name = value

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	parentPath, err := askText(reader, "2) Project directory path? (example: ./projects)", opts.ParentPath, true)
	if err != nil {
		return opts, err
	}
	opts.ParentPath = parentPath

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	defaultProjectType := defaultFromList(opts.ProjectType, config.ProjectTypes, config.ProjectTypes[0])
	projectType, err := askChoiceWithDescriptions(reader, "3) Project type?", config.ProjectTypes, projectTypeDescriptions(), defaultProjectType)
	if err != nil {
		return opts, err
	}
	opts.ProjectType = projectType

	fmt.Println()
	defaultAIFeature := defaultFromList(opts.AIFeature, config.AIFeatureTypes, config.AIFeatureTypes[0])
	aiFeature, err := askChoiceWithDescriptions(reader, "4) AI feature?", config.AIFeatureTypes, aiFeatureDescriptions(), defaultAIFeature)
	if err != nil {
		return opts, err
	}
	opts.AIFeature = aiFeature

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	defaultAIAgent := defaultFromList(opts.AIAgent, config.AgentTypes, "codex")
	aiAgent, err := askChoiceWithDescriptions(reader, "5) AI agent?", config.AgentTypes, aiAgentDescriptions(), defaultAIAgent)
	if err != nil {
		return opts, err
	}
	opts.AIAgent = aiAgent

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	archDescs := architectureDescriptions()
	defaultArchitecture := defaultFromList(opts.Architecture, config.ArchitectureTypes, inferArchitectureFromExisting(opts))
	architecture, err := askChoiceWithDescriptions(reader, "6) Architecture?", config.ArchitectureTypes, archDescs, defaultArchitecture)
	if err != nil {
		return opts, err
	}
	opts.Architecture = architecture

	frontendSummary := "none"
	backendSummary := "none"
	database := "none"
	stackParts := []string{opts.Architecture}
	if opts.AIFeature != "none" {
		stackParts = append(stackParts, opts.AIFeature)
	}

	switch opts.Architecture {
	case "cli-tool":
		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()

		backendLangDefault := defaultFromList(opts.BackendType, backendLanguageOptions, backendLanguageOptions[0])
		backendLang, err := askChoice(reader, "7) CLI language?", backendLanguageOptions, backendLangDefault)
		if err != nil {
			return opts, err
		}
		opts.BackendType = backendLang
		opts.FrontendType = "none"
		backendSummary = fmt.Sprintf("%s (cli)", backendLang)
		stackParts = append(stackParts, backendLang)
	case "backend-service":
		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()

		backendLangDefault := defaultFromList(opts.BackendType, backendLanguageOptions, backendLanguageOptions[0])
		backendLang, err := askChoice(reader, "7) Backend language?", backendLanguageOptions, backendLangDefault)
		if err != nil {
			return opts, err
		}
		db, err := askChoice(reader, "8) Database?", databaseOptions, databaseOptions[0])
		if err != nil {
			return opts, err
		}
		opts.BackendType = backendLang
		opts.FrontendType = "none"
		backendSummary = backendLang
		database = db
		stackParts = append(stackParts, backendLang, db)
	case "frontend-app":
		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()

		frontendDefault := defaultFromList(opts.FrontendType, frontendFrameworkOptions, frontendFrameworkOptions[0])
		frontendFramework, err := askChoice(reader, "7) Frontend framework?", frontendFrameworkOptions, frontendDefault)
		if err != nil {
			return opts, err
		}
		frontendLang, err := askChoice(reader, "8) Frontend language?", frontendLanguageOptions, frontendLanguageOptions[0])
		if err != nil {
			return opts, err
		}
		opts.FrontendType = frontendFramework
		opts.BackendType = "none"
		frontendSummary = fmt.Sprintf("%s (%s)", frontendFramework, frontendLang)
		stackParts = append(stackParts, frontendFramework, frontendLang)
	case "fullstack-web-app":
		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()

		defaultFullstackFramework := "next"
		if opts.FrontendType == "nuxt" {
			defaultFullstackFramework = "nuxt"
		}
		fullstackFramework, err := askChoiceWithDescriptions(reader, "7) Framework?", fullstackFrameworkOptions, fullstackFrameworkDescriptions(), defaultFullstackFramework)
		if err != nil {
			return opts, err
		}
		frameworkLang, err := askChoice(reader, "8) Language?", frontendLanguageOptions, frontendLanguageOptions[0])
		if err != nil {
			return opts, err
		}
		opts.FrontendType = fullstackFramework
		opts.BackendType = "none"
		frontendSummary = fmt.Sprintf("%s (%s)", fullstackFramework, frameworkLang)
		backendSummary = fmt.Sprintf("%s runtime", fullstackFramework)
		stackParts = append(stackParts, fullstackFramework, frameworkLang)
	case "frontend-backend":
		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()

		frontendDefault := defaultFromList(opts.FrontendType, frontendFrameworkOptions, frontendFrameworkOptions[0])
		frontendFramework, err := askChoice(reader, "7) Frontend framework?", frontendFrameworkOptions, frontendDefault)
		if err != nil {
			return opts, err
		}
		frontendLang, err := askChoice(reader, "8) Frontend language?", frontendLanguageOptions, frontendLanguageOptions[0])
		if err != nil {
			return opts, err
		}
		backendLangDefault := defaultFromList(opts.BackendType, backendLanguageOptions, backendLanguageOptions[0])
		backendLang, err := askChoice(reader, "9) Backend language?", backendLanguageOptions, backendLangDefault)
		if err != nil {
			return opts, err
		}

		opts.FrontendType = frontendFramework
		opts.BackendType = backendLang
		frontendSummary = fmt.Sprintf("%s (%s)", frontendFramework, frontendLang)
		backendSummary = backendLang
		stackParts = append(stackParts, frontendFramework, frontendLang, backendLang)
	}

	if opts.FrontendType != "none" && frontendSummary == "none" {
		frontendSummary = opts.FrontendType
	}
	if opts.BackendType != "none" && backendSummary == "none" {
		backendSummary = opts.BackendType
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	deploymentTarget, err := askChoice(reader, "12) Deployment target?", deploymentTargetOptions, deploymentTargetOptions[0])
	if err != nil {
		return opts, err
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	ciCD, err := askChoiceWithDescriptions(reader, "13) Include CI/CD?", config.YesNoTypes, map[string]string{
		"yes": "GitHub Actions",
		"no":  "disabled",
	}, "yes")
	if err != nil {
		return opts, err
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	authType, err := askChoice(reader, "14) Authentication?", authOptions, authOptions[0])
	if err != nil {
		return opts, err
	}

	opts.DockerCompose = "no"
	if deploymentTarget == "docker compose" {
		opts.DockerCompose = "yes"
	}
	opts.TechStack = strings.Join(stackParts, " | ")

	fmt.Println()
	fmt.Println("Final configuration summary")
	fmt.Println()
	fmt.Printf("Project name: %s\n", opts.Name)
	fmt.Printf("Project type: %s\n", opts.ProjectType)
	fmt.Printf("AI feature: %s\n", opts.AIFeature)
	fmt.Printf("AI agent: %s\n", opts.AIAgent)
	fmt.Printf("Architecture: %s\n", archDescs[opts.Architecture])
	fmt.Printf("Frontend: %s\n", frontendSummary)
	fmt.Printf("Backend: %s\n", backendSummary)
	fmt.Printf("Database: %s\n", database)
	fmt.Printf("Deployment: %s\n", deploymentTarget)
	fmt.Printf("CI/CD: %s\n", yesNoLabel(ciCD))
	fmt.Printf("Auth: %s\n", authType)
	fmt.Println()

	yes, err := askYesNo(reader, "Confirm before create?", true)
	if err != nil {
		return opts, err
	}
	if !yes {
		return opts, ErrAborted
	}

	return opts, nil
}

func ConfirmOverwrite(projectRoot string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	question := fmt.Sprintf("Target already exists: %s. Backup existing directory and continue?", projectRoot)
	return askYesNo(reader, question, false)
}

func QuestionsFor(opts config.CreateOptions, askAll bool) []string {
	questions := make([]string, 0, 12)
	if askAll || opts.Name == "" {
		questions = append(questions, "Project name")
	}
	if askAll || strings.TrimSpace(opts.ParentPath) == "" {
		questions = append(questions, "Project directory path")
	}
	if askAll || opts.ProjectType == "" {
		questions = append(questions, "Project type")
	}
	if askAll || opts.AIFeature == "" {
		questions = append(questions, "AI feature")
	}
	if askAll || strings.TrimSpace(opts.AIAgent) == "" {
		questions = append(questions, "AI agent")
	}
	if askAll || opts.Architecture == "" {
		questions = append(questions, "Architecture")
	}
	questions = append(questions, "Architecture-specific options")
	questions = append(questions, "Deployment target")
	questions = append(questions, "Include CI/CD")
	questions = append(questions, "Authentication")
	questions = append(questions, "Confirm before create")
	return questions
}

func inferArchitectureFromExisting(opts config.CreateOptions) string {
	switch opts.Architecture {
	case "cli-tool", "backend-service", "frontend-app", "fullstack-web-app", "frontend-backend":
		return opts.Architecture
	}
	frontend := strings.TrimSpace(strings.ToLower(opts.FrontendType))
	backend := strings.TrimSpace(strings.ToLower(opts.BackendType))
	if frontend == "next" || frontend == "nuxt" {
		return "fullstack-web-app"
	}
	if frontend != "" && frontend != "none" && backend != "" && backend != "none" {
		return "frontend-backend"
	}
	if frontend != "" && frontend != "none" {
		return "frontend-app"
	}
	if backend != "" && backend != "none" {
		return "backend-service"
	}
	return config.ArchitectureTypes[0]
}

func defaultFromList(value string, options []string, fallback string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	for _, option := range options {
		if normalized == option {
			return option
		}
	}
	return fallback
}

func yesNoLabel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes":
		return "enabled"
	case "no":
		return "disabled"
	default:
		return value
	}
}

func boolToYesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func yesNoAsBool(value string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes", "y", "true":
		return true
	case "no", "n", "false":
		return false
	default:
		return fallback
	}
}

func defaultYesNoValue(value string, fallback bool) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "yes" || normalized == "no" {
		return normalized
	}
	return boolToYesNo(fallback)
}

func askText(reader *bufio.Reader, label, defaultValue string, required bool) (string, error) {
	for {
		if defaultValue != "" {
			fmt.Printf("%s [%s]: ", label, defaultValue)
		} else {
			fmt.Printf("%s: ", label)
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		value := strings.TrimSpace(line)
		if value == "" {
			value = defaultValue
		}
		if required && strings.TrimSpace(value) == "" {
			fmt.Println("This field is required.")
			continue
		}
		return value, nil
	}
}

func askChoice(reader *bufio.Reader, label string, options []string, defaultValue string) (string, error) {
	return askChoiceWithDescriptions(reader, label, options, nil, defaultValue)
}

func askChoiceWithDescriptions(reader *bufio.Reader, label string, options []string, descriptions map[string]string, defaultValue string) (string, error) {
	for {
		fmt.Printf("%s\n", label)
		for i, option := range options {
			display := option
			if desc := strings.TrimSpace(descriptions[option]); desc != "" {
				display = fmt.Sprintf("%s - %s", option, desc)
			}
			fmt.Printf("  %d) %s\n", i+1, display)
		}

		prompt := "Select"
		if defaultValue != "" {
			prompt = fmt.Sprintf("Select [%s]", defaultValue)
		}
		fmt.Printf("%s: ", prompt)

		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(strings.ToLower(line))
		if value == "" {
			if defaultValue != "" {
				return defaultValue, nil
			}
			return options[0], nil
		}

		if index, convErr := strconv.Atoi(value); convErr == nil {
			if index >= 1 && index <= len(options) {
				return options[index-1], nil
			}
		}

		for _, option := range options {
			if value == option {
				return option, nil
			}
		}

		fmt.Println("Invalid choice. Use index or value from the list.")
	}
}

func askYesNo(reader *bufio.Reader, label string, defaultYes bool) (bool, error) {
	defaultLabel := "y/N"
	if defaultYes {
		defaultLabel = "Y/n"
	}

	for {
		fmt.Printf("%s (%s): ", label, defaultLabel)
		line, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		value := strings.TrimSpace(strings.ToLower(line))
		if value == "" {
			return defaultYes, nil
		}
		switch value {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Please enter y/yes or n/no.")
		}
	}
}
