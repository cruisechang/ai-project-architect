package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	templatefs "project-generator/templates"

	"project-generator/internal/config"
)

type TemplateData struct {
	ProjectName        string
	ProjectType        string
	AIAgent            string
	ProjectDescription string
	BackendType        string
	FrontendType       string
	DocsType           string
	Architecture       string
	TechStack          string
	UnitTest           string
	IntegrationTest    string
	E2ETest            string
	PerformanceTest    string
	SecurityTest       string
	UATTest            string
	DockerCompose      string
	PromptTitle        string
	IncludedSkills     []string
	IncludedSkillsCSV  string
	GeneratedAt        string
	Goals              []string
	Requirements       []string
	Notes              []string
}

func buildTemplateData(opts config.CreateOptions) TemplateData {
	goals, requirements, notes := defaultsByProjectType(opts.ProjectType)
	skillsCSV := "none"
	if len(opts.SelectedSkills) > 0 {
		skillsCSV = strings.Join(opts.SelectedSkills, ", ")
	}
	description := opts.Description
	if description == "" {
		description = "TBD"
	}
	architecture := opts.Architecture
	if architecture == "" {
		architecture = "TBD"
	}
	techStack := opts.TechStack
	if techStack == "" {
		techStack = "TBD"
	}

	return TemplateData{
		ProjectName:        opts.Name,
		ProjectType:        opts.ProjectType,
		AIAgent:            opts.AIAgent,
		ProjectDescription: description,
		BackendType:        opts.BackendType,
		FrontendType:       opts.FrontendType,
		DocsType:           opts.DocsType,
		Architecture:       architecture,
		TechStack:          techStack,
		UnitTest:           yesNoLabel(opts.UnitTest),
		IntegrationTest:    yesNoLabel(opts.IntegrationTest),
		E2ETest:            yesNoLabel(opts.E2ETest),
		PerformanceTest:    yesNoLabel(opts.PerformanceTest),
		SecurityTest:       yesNoLabel(opts.SecurityTest),
		UATTest:            yesNoLabel(opts.UATTest),
		DockerCompose:      yesNoLabel(opts.DockerCompose),
		PromptTitle:        opts.PromptTitle,
		IncludedSkills:     opts.SelectedSkills,
		IncludedSkillsCSV:  skillsCSV,
		GeneratedAt:        time.Now().Format(time.RFC3339),
		Goals:              goals,
		Requirements:       requirements,
		Notes:              notes,
	}
}

func renderTemplate(templateName string, data TemplateData) (string, error) {
	raw, err := templatefs.FS.ReadFile(templateName)
	if err != nil {
		return "", fmt.Errorf("cannot load template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return buf.String(), nil
}

func defaultsByProjectType(projectType string) ([]string, []string, []string) {
	switch projectType {
	case "ai-app":
		return []string{
				"Build a reliable AI app with strong prompt design and model quality.",
				"Keep agent workflow and knowledge retrieval explicit and measurable.",
			}, []string{
				"Define prompt design guidelines and prompt versioning strategy.",
				"Select model(s) with clear tradeoffs for quality, latency, and reliability.",
				"Design an evaluation framework with baseline and regression checks.",
				"Define agent workflow and tool boundaries for each task step.",
				"Define RAG strategy (indexing, retrieval quality, and grounding).",
				"Set cost control rules (token budget, usage tracking, and optimization).",
			}, []string{
				"Track evaluation scores and cost trends for each release.",
			}
	case "web-app":
		return []string{
				"Deliver a web product with coherent UI, API, and secure authentication.",
				"Keep critical user flows and frontend architecture explicit from day one.",
			}, []string{
				"Define UI scope and key screens for the first release.",
				"Define API contracts and integration boundaries.",
				"Define authentication and authorization approach.",
				"Document critical user flows and edge paths.",
				"Define frontend architecture and state/data flow strategy.",
			}, []string{
				"Validate product behavior end-to-end across UI/API/Auth.",
			}
	case "devops-tool":
		return []string{
				"Provide safe DevOps automation for delivery and operations.",
				"Make deployment and recovery steps repeatable.",
			}, []string{
				"Define CI/CD pipeline stages and quality gates.",
				"Define deployment process and environment promotion flow.",
				"Define rollback strategy with clear trigger criteria.",
				"Define infrastructure operation procedures and ownership.",
				"Automate repetitive operations with observable logs and alerts.",
			}, []string{
				"Treat failures as actionable diagnostics with fast rollback paths.",
			}
	case "internal-tool":
		return []string{
				"Ship maintainable internal tooling that improves team workflow efficiency.",
				"Provide visibility and control for daily operations.",
			}, []string{
				"Define target workflow and automation opportunities.",
				"Define data management model, ownership, and quality checks.",
				"Document onboarding flow for new users and operators.",
				"Define operation dashboard scope and key metrics.",
			}, []string{
				"Prefer straightforward architecture and fast operator feedback loops.",
			}
	case "platform-service":
		return []string{
				"Provide shared infrastructure capabilities as internal platform services.",
				"Enable developer teams through stable platform abstractions.",
			}, []string{
				"Define shared infrastructure capabilities and boundaries.",
				"Define internal platform interfaces and team responsibilities.",
				"Define developer platform experience and self-service workflows.",
				"Define service layer contracts, versioning, and compatibility policy.",
			}, []string{
				"Treat platform contracts as long-term interfaces with strong backward compatibility.",
			}
	default:
		return []string{"Define clear project outcomes."}, []string{"Document technical decisions."}, []string{"Keep scope focused."}
	}
}

func yesNoLabel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes":
		return "Yes"
	case "no":
		return "No"
	default:
		return "TBD"
	}
}
