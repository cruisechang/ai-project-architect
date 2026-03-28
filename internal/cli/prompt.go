package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"project-generator/internal/architect/model"
	"project-generator/internal/architect/modules/planner"
	"project-generator/internal/config"
	"project-generator/internal/i18n"
)

var ErrAborted = errors.New("operation aborted")

var frontendFrameworkOptions = []string{"react", "vue", "pure-typescript"}
var frontendLanguageOptions = []string{"typescript", "javascript"}
var backendLanguageOptions = []string{"go", "python", "node"}
var databaseOptions = []string{"postgres", "mysql", "sqlite", "none"}
var deploymentTargetOptions = []string{"docker compose", "kubernetes", "vm/server", "serverless"}
var authOptions = []string{"none", "jwt", "oauth"}

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
		"universal":   i18n.T("wizard.ai-agent.universal"),
	}
}

func architectureDescriptions() map[string]string {
	return map[string]string{
		"cli":               i18n.T("wizard.architecture.cli"),
		"server":            i18n.T("wizard.architecture.server"),
		"web-app-server":    i18n.T("wizard.architecture.web-app-server"),
		"mobile-app-server": i18n.T("wizard.architecture.mobile-app-server"),
		"web-app":           i18n.T("wizard.architecture.web-app"),
		"mobile-app":        i18n.T("wizard.architecture.mobile-app"),
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

	idea, err := askMultilineText(reader, "1) Project idea")
	if err != nil {
		return opts, err
	}
	opts.Idea = strings.TrimSpace(idea)
	if opts.Idea != "" {
		inferred := planner.New().Infer(opts.Idea, opts.Name)
		applyContextToCreateOptions(inferred, &opts)
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	value, err := askText(reader, "2) Project name? (example: ops-platform)", opts.Name, true)
	if err != nil {
		return opts, err
	}
	opts.Name = value

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	parentPath, err := askText(reader, "3) Project directory path? (example: ./projects)", opts.ParentPath, true)
	if err != nil {
		return opts, err
	}
	opts.ParentPath = parentPath

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	defaultArchitecture := defaultFromList(opts.Architecture, config.ArchitectureTypes, config.ArchitectureTypes[2])
	architecture, err := askChoiceWithDescriptions(reader, "4) Type?", config.ArchitectureTypes, architectureDescriptions(), defaultArchitecture)
	if err != nil {
		return opts, err
	}
	opts.Architecture = architecture
	opts.ProjectType = config.InferProjectTypeFromArchitecture(opts.Architecture)

	fmt.Println()
	defaultAIFeature := defaultFromList(opts.AIFeature, config.AIFeatureTypes, config.AIFeatureTypes[0])
	aiFeature, err := askChoiceWithDescriptions(reader, "5) AI feature?", config.AIFeatureTypes, aiFeatureDescriptions(), defaultAIFeature)
	if err != nil {
		return opts, err
	}
	opts.AIFeature = aiFeature

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	defaultAIAgent := defaultFromList(opts.AIAgent, config.AgentTypes, "codex")
	aiAgent, err := askChoiceWithDescriptions(reader, "6) AI agent?", config.AgentTypes, aiAgentDescriptions(), defaultAIAgent)
	if err != nil {
		return opts, err
	}
	opts.AIAgent = aiAgent

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	frontendAppType := "web-app"
	backendRuntime := "server"
	switch opts.Architecture {
	case "mobile-app", "mobile-app-server":
		frontendAppType = "mobile-app"
	}
	switch opts.Architecture {
	case "cli":
		backendRuntime = "cli"
	case "server", "web-app-server", "mobile-app-server":
		backendRuntime = "server"
	}

	frontendSummary := "none"
	backendSummary := "none"
	database := "none"
	stackParts := []string{opts.Architecture}
	if opts.AIFeature != "none" {
		stackParts = append(stackParts, opts.AIFeature)
	}
	if frontendAppType != "" && opts.Architecture != "cli" && opts.Architecture != "server" {
		stackParts = append(stackParts, frontendAppType)
	}
	if opts.Architecture == "cli" || opts.Architecture == "server" {
		stackParts = append(stackParts, backendRuntime)
	}

	switch opts.Architecture {
	case "cli":
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
	case "server":
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
	case "web-app", "mobile-app":
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
	case "web-app-server", "mobile-app-server":
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

	deploymentTarget, err := askChoice(reader, "11) Deployment target?", deploymentTargetOptions, deploymentTargetOptions[0])
	if err != nil {
		return opts, err
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	ciCD, err := askChoiceWithDescriptions(reader, "12) Include CI/CD?", config.YesNoTypes, map[string]string{
		"yes": "GitHub Actions",
		"no":  "disabled",
	}, "yes")
	if err != nil {
		return opts, err
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	authType, err := askChoice(reader, "13) Authentication?", authOptions, authOptions[0])
	if err != nil {
		return opts, err
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	unitTest, err := askChoiceWithDescriptions(reader, "14) Unit tests?", config.YesNoTypes, map[string]string{
		"yes": "verify individual functions/modules",
		"no":  "skip unit-level checks",
	}, defaultYesNoValue(opts.UnitTest, true))
	if err != nil {
		return opts, err
	}
	integrationTest, err := askChoiceWithDescriptions(reader, "15) Integration tests?", config.YesNoTypes, map[string]string{
		"yes": "verify module/service interactions",
		"no":  "skip cross-module checks",
	}, defaultYesNoValue(opts.IntegrationTest, true))
	if err != nil {
		return opts, err
	}
	apiTest, err := askChoiceWithDescriptions(reader, "16) API tests?", config.YesNoTypes, map[string]string{
		"yes": "verify API contracts, status codes, and payloads",
		"no":  "skip API contract checks",
	}, defaultYesNoValue(opts.APITest, true))
	if err != nil {
		return opts, err
	}
	e2eTest, err := askChoiceWithDescriptions(reader, "17) E2E tests?", config.YesNoTypes, map[string]string{
		"yes": "verify end-to-end user flows in realistic scenarios",
		"no":  "skip full-flow scenario checks",
	}, defaultYesNoValue(opts.E2ETest, true))
	if err != nil {
		return opts, err
	}
	opts.UnitTest = unitTest
	opts.APITest = apiTest
	opts.IntegrationTest = integrationTest
	opts.E2ETest = e2eTest

	opts.DockerCompose = "no"
	if deploymentTarget == "docker compose" {
		opts.DockerCompose = "yes"
	}
	opts.TechStack = strings.Join(stackParts, " | ")
	archDescs := architectureDescriptions()

	fmt.Println()
	fmt.Println("Final configuration summary")
	fmt.Println()
	fmt.Printf("Project name: %s\n", opts.Name)
	fmt.Printf("Type: %s\n", archDescs[opts.Architecture])
	fmt.Printf("AI feature: %s\n", opts.AIFeature)
	fmt.Printf("AI agent: %s\n", opts.AIAgent)
	if opts.Architecture != "cli" && opts.Architecture != "server" {
		fmt.Printf("Frontend app type: %s\n", frontendAppType)
	}
	if opts.Architecture != "web-app" && opts.Architecture != "mobile-app" {
		fmt.Printf("Backend runtime: %s\n", backendRuntime)
	}
	fmt.Printf("Frontend: %s\n", frontendSummary)
	fmt.Printf("Backend: %s\n", backendSummary)
	fmt.Printf("Database: %s\n", database)
	fmt.Printf("Deployment: %s\n", deploymentTarget)
	fmt.Printf("CI/CD: %s\n", yesNoLabel(ciCD))
	fmt.Printf("Auth: %s\n", authType)
	fmt.Printf("Tests: unit=%s, api=%s, integration=%s, e2e=%s\n",
		yesNoLabel(opts.UnitTest), yesNoLabel(opts.APITest), yesNoLabel(opts.IntegrationTest), yesNoLabel(opts.E2ETest))
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
	questions := make([]string, 0, 16)
	if askAll || strings.TrimSpace(opts.Idea) == "" {
		questions = append(questions, "Project idea")
	}
	if askAll || opts.Name == "" {
		questions = append(questions, "Project name")
	}
	if askAll || strings.TrimSpace(opts.ParentPath) == "" {
		questions = append(questions, "Project directory path")
	}
	if askAll || opts.Architecture == "" {
		questions = append(questions, "Type")
	}
	if askAll || opts.AIFeature == "" {
		questions = append(questions, "AI feature")
	}
	if askAll || strings.TrimSpace(opts.AIAgent) == "" {
		questions = append(questions, "AI agent")
	}
	questions = append(questions, "Architecture-specific options")
	questions = append(questions, "Deployment target")
	questions = append(questions, "Include CI/CD")
	questions = append(questions, "Authentication")
	questions = append(questions, "Unit tests")
	questions = append(questions, "Integration tests")
	questions = append(questions, "API tests")
	questions = append(questions, "E2E tests")
	questions = append(questions, "Confirm before create")
	return questions
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

func askMultilineText(reader *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s (%s):\n", label, i18n.T("prompt.multiline.hint"))
	var lines []string
	emptyStreak := 0
	for {
		line, err := reader.ReadString('\n')
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			emptyStreak++
			if err != nil {
				break
			}
			if emptyStreak >= 2 {
				break
			}
			lines = append(lines, "")
		} else {
			emptyStreak = 0
			lines = append(lines, trimmed)
		}
		if err != nil {
			break
		}
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n"), nil
}

func applyContextToCreateOptions(ctx model.ProjectContext, opts *config.CreateOptions) {
	if opts.Name == "" {
		opts.Name = ctx.ProjectName
	}
	if opts.BackendType == "" {
		opts.BackendType = ctx.BackendLanguage
	}
	if opts.FrontendType == "" {
		opts.FrontendType = ctx.FrontendFramework
	}
	if opts.Architecture == "" {
		opts.Architecture = inferArchitectureFromContext(ctx)
	}
	if opts.TechStack == "" {
		opts.TechStack = buildTechStack(opts.Architecture, opts.BackendType, opts.FrontendType)
	}
	if opts.DockerCompose == "" {
		if strings.Contains(strings.ToLower(ctx.Deployment), "docker") {
			opts.DockerCompose = "yes"
		} else {
			opts.DockerCompose = "no"
		}
	}
	if opts.ProjectType == "" {
		opts.ProjectType = config.InferProjectTypeFromArchitecture(opts.Architecture)
	}
}

func inferArchitectureFromContext(ctx model.ProjectContext) string {
	frontend := strings.ToLower(ctx.FrontendFramework)
	switch ctx.ProjectType {
	case "backend":
		return "server"
	case "frontend":
		if frontend == "react-native" || frontend == "flutter" || frontend == "swiftui" || frontend == "kotlin" || frontend == "android" || frontend == "ios" {
			return "mobile-app"
		}
		return "web-app"
	case "fullstack":
		if frontend == "react-native" || frontend == "flutter" || frontend == "swiftui" || frontend == "kotlin" || frontend == "android" || frontend == "ios" {
			return "mobile-app-server"
		}
		return "web-app-server"
	}
	return ""
}

func buildTechStack(architecture, backend, frontend string) string {
	parts := []string{}
	if architecture != "" {
		parts = append(parts, architecture)
	}
	if backend != "" && backend != "none" {
		parts = append(parts, backend)
	}
	if frontend != "" && frontend != "none" {
		parts = append(parts, frontend)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " | ")
}
