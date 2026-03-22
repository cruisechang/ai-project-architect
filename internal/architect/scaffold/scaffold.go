package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

type Result struct {
	CreatedDirs  []string
	CreatedFiles []string
}

func Generate(root string) (Result, error) {
	dirs := []string{
		"apps",
		"services",
		"libs",
		"docs",
		"infrastructure",
		"modules/planner",
		"modules/generators",
		"modules/templates",
		"generators/prd_generator",
		"generators/spec_generator",
		"generators/architecture_generator",
		"generators/api_generator",
		"generators/db_generator",
		"generators/scaffold_generator",
		"templates",
		".architect",
		".codex",
		".codex/cache",
	}

	files := map[string]string{
		"apps/README.md":                              "# apps\n\nApplication entry points and deployable apps.\n",
		"services/README.md":                          "# services\n\nService domain modules and orchestration layers.\n",
		"libs/README.md":                              "# libs\n\nShared libraries and utility components.\n",
		"infrastructure/README.md":                    "# infrastructure\n\nDeployment and infrastructure definitions.\n",
		"modules/planner/README.md":                   "# planner module\n\nProject context planning pipeline.\n",
		"modules/generators/README.md":                "# generators module\n\nDocument generator orchestration and registry.\n",
		"modules/templates/README.md":                 "# templates module\n\nReusable markdown templates.\n",
		"generators/prd_generator/README.md":          "# prd_generator\n\nGenerates PRD.md\n",
		"generators/spec_generator/README.md":         "# spec_generator\n\nGenerates SPEC.md\n",
		"generators/architecture_generator/README.md": "# architecture_generator\n\nGenerates ARCHITECTURE.md\n",
		"generators/api_generator/README.md":          "# api_generator\n\nGenerates API.md\n",
		"generators/db_generator/README.md":           "# db_generator\n\nGenerates DB_SCHEMA.md\n",
		"generators/scaffold_generator/README.md":     "# scaffold_generator\n\nGenerates repository scaffold structure.\n",
		"templates/README.md":                         "# templates\n\nTemplate source files for generated documents.\n",
		".codex/config.json":                          "{\n  \"agent\": \"codex\",\n  \"approvalMode\": \"on-request\",\n  \"sandboxMode\": \"workspace-write\"\n}\n",
		".codex/memory.md":                            "# Agent Memory\n\nUse this file to keep project-specific notes and decisions.\n",
	}

	result := Result{}
	for _, d := range dirs {
		abs := filepath.Join(root, d)
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return result, fmt.Errorf("create directory %s: %w", abs, err)
		}
		result.CreatedDirs = append(result.CreatedDirs, abs)
	}

	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			return result, fmt.Errorf("write file %s: %w", abs, err)
		}
		result.CreatedFiles = append(result.CreatedFiles, abs)
	}

	return result, nil
}
