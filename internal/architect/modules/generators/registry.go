package generators

import (
	"fmt"
	"path/filepath"

	"project-generator/internal/architect/generators/api_generator"
	"project-generator/internal/architect/generators/architecture_generator"
	"project-generator/internal/architect/generators/db_generator"
	"project-generator/internal/architect/generators/implementation_plan_generator"
	"project-generator/internal/architect/generators/prd_generator"
	"project-generator/internal/architect/generators/repository_structure_generator"
	"project-generator/internal/architect/generators/spec_generator"
	"project-generator/internal/architect/model"
)

type Artifact struct {
	Key      string
	FilePath string
}

var orderedArtifacts = []Artifact{
	{Key: "prd", FilePath: filepath.Join("docs", "PRD.md")},
	{Key: "spec", FilePath: filepath.Join("docs", "SPEC.md")},
	{Key: "architecture", FilePath: filepath.Join("docs", "ARCHITECTURE.md")},
	{Key: "api", FilePath: filepath.Join("docs", "API.md")},
	{Key: "db", FilePath: filepath.Join("docs", "DB_SCHEMA.md")},
	{Key: "implementation-plan", FilePath: filepath.Join("docs", "IMPLEMENTATION_PLAN.md")},
	{Key: "repository-structure", FilePath: filepath.Join("docs", "REPOSITORY_STRUCTURE.md")},
}

func DefaultArtifacts() []Artifact {
	defaultKeys := map[string]struct{}{
		"prd":                 {},
		"spec":                {},
		"architecture":        {},
		"api":                 {},
		"db":                  {},
		"implementation-plan": {},
	}
	out := make([]Artifact, 0, len(defaultKeys))
	for _, artifact := range orderedArtifacts {
		if _, ok := defaultKeys[artifact.Key]; ok {
			out = append(out, artifact)
		}
	}
	return out
}

func ResolveTarget(target string) ([]Artifact, error) {
	if target == "project-structure" {
		target = "repository-structure"
	}

	switch target {
	case "docs", "all", "":
		return DefaultArtifacts(), nil
	case "architecture", "api", "scaffold", "prd", "spec", "db", "implementation-plan", "repository-structure":
		if target == "scaffold" {
			return nil, nil
		}
		for _, artifact := range orderedArtifacts {
			if artifact.Key == target {
				return []Artifact{artifact}, nil
			}
		}
		return nil, fmt.Errorf("unsupported artifact target: %s", target)
	default:
		return nil, fmt.Errorf("unsupported artifact target: %s", target)
	}
}

func Generate(ctx model.ProjectContext, artifact Artifact) (string, error) {
	switch artifact.Key {
	case "prd":
		return prd_generator.Generate(ctx)
	case "spec":
		return spec_generator.Generate(ctx)
	case "architecture":
		return architecture_generator.Generate(ctx)
	case "api":
		return api_generator.Generate(ctx)
	case "db":
		return db_generator.Generate(ctx)
	case "implementation-plan":
		return implementation_plan_generator.Generate(ctx)
	case "repository-structure":
		return repository_structure_generator.Generate(ctx)
	default:
		return "", fmt.Errorf("unsupported artifact key: %s", artifact.Key)
	}
}
