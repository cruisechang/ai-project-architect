package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"project-generator/internal/architect/model"
)

var catalog = map[string]string{
	"PRD": `# PRD

## Product Overview
- Project: {{ .ProjectName }}
- Idea: {{ .ProjectIdea }}
- Target Type: {{ .ProjectType }}

## Problem Statement
Users need a reliable way to execute the core idea: {{ .ProjectIdea }}.
Current workflows are fragmented, manual, or slow for repeated operations.

## Goals
- Deliver a production-ready foundation for {{ .ProjectName }}.
- Provide clear functional scope and engineering handoff.
- Enable incremental delivery with measurable outcomes.

## User Personas
- Product owner: defines priorities and validates business outcomes.
- End user: consumes the core product workflow.
- Operator/engineer: maintains system reliability and observability.

## Core Features
- Project planning and architecture baseline generation.
- Document generation: PRD, SPEC, ARCHITECTURE, API, DB schema, implementation plan.
- Repeatable scaffold generation aligned to architecture.

## Non-Goals
- Building every feature in one release.
- Locking implementation to a single cloud/provider.

## Success Metrics
- Time to first blueprint generated.
- Documentation completeness score in review.
- Lead time reduction from idea to implementation start.
`,
	"SPEC": `# SPEC

## Technical Scope
- Project: {{ .ProjectName }}
- Project Type: {{ .ProjectType }}
- Frontend Framework: {{ .FrontendFramework }}
- Backend: {{ .BackendLanguage }} / {{ .BackendFramework }}
- Database: {{ .Database }}

## System Modules
- Planner module: parse and normalize project requirements.
- Document generators: produce deterministic markdown artifacts.
- Scaffold generator: create repository baseline and folders.
- Template engine: render reusable templates into final docs.

## Feature Descriptions
- Parse project idea and infer stack deterministically.
- Generate engineering documents with fixed sections.
- Create implementation plan and scaffold for execution.

## Functional Requirements
- Must generate PRD.md, SPEC.md, ARCHITECTURE.md, API.md, DB_SCHEMA.md.
- Must generate IMPLEMENTATION_PLAN.md.
- Must support optional REPOSITORY_STRUCTURE.md output.
- Each generation step must be independently rerunnable.

## Non-Functional Requirements
- Deterministic output for same input context.
- Modular implementation with isolated generators.
- CLI-first UX with scriptable commands.

## Edge Cases
- Missing project idea.
- Unsupported stack keywords.
- Existing target directory conflicts.

## Constraints
- Markdown output only.
- Context-driven generation with minimal hidden assumptions.
`,
	"ARCHITECTURE": `# ARCHITECTURE

## High-Level Architecture
The system follows a modular pipeline:

idea -> planner -> prd -> spec -> architecture -> api -> database -> scaffold

## Service Boundaries
- Planner: context extraction and stack inference.
- Generators: document-specific rendering units.
- Scaffold module: repository and folder bootstrap.

## Components
- CLI layer (` + "`cmd/`" + `): command routing and UX.
- Planning layer: project context model construction.
- Generator layer: PRD/SPEC/ARCHITECTURE/API/DB/PLAN outputs.
- Template layer: reusable deterministic markdown templates.

## Data Flow
1. User provides project idea.
2. Planner infers normalized context.
3. Generators produce docs from context.
4. Scaffold writer creates repository baseline.

## Deployment Architecture
- Default target: {{ .Deployment }}
- Can be adapted to container, Kubernetes, VM, or serverless.

## Scalability Considerations
- Add new generator without changing existing command contracts.
- Keep each artifact generation independently rerunnable.
- Keep context model versioned for compatibility.
`,
	"API": `# API

## API Style
OpenAPI-style REST design (can be upgraded to strict OpenAPI YAML later).

## Base Path
` + "`/api/v1`" + `

## Authentication
- Mode: {{ .Authentication }}

## Endpoints
### GET /health
- Purpose: service health check
- Response: { "status": "ok" }

### POST /projects/plan
- Purpose: infer project context from idea
- Request:
  - idea: string
  - project_name: string
- Response:
  - project_context object

### POST /projects/docs
- Purpose: generate architecture documents
- Request:
  - context: object
  - targets: string[]
- Response:
  - generated_files: string[]

### POST /projects/scaffold
- Purpose: generate project scaffold
- Request:
  - context: object
  - root_path: string
- Response:
  - created_directories: string[]
  - created_files: string[]

## Error Handling
- 400: invalid request
- 404: context not found
- 409: path conflict
- 500: generation failure
`,
	"DB_SCHEMA": `# DB_SCHEMA

## Database
- Engine: {{ .Database }}

## Tables
### projects
- id (PK)
- name
- idea
- project_type
- frontend_framework
- backend_language
- backend_framework
- database_engine
- authentication
- deployment
- created_at

### artifacts
- id (PK)
- project_id (FK -> projects.id)
- artifact_type
- file_path
- generated_at
- checksum

### generation_runs
- id (PK)
- project_id (FK -> projects.id)
- run_mode (new|generate|autopilot)
- status
- started_at
- finished_at

## Relationships
- projects 1:N artifacts
- projects 1:N generation_runs

## Indexes
- idx_projects_name
- idx_artifacts_project_type
- idx_runs_project_status

## Constraints
- projects.name unique
- artifacts.artifact_type + file_path unique per project
`,
	"IMPLEMENTATION_PLAN": `# IMPLEMENTATION_PLAN

## Phase 1 - Planning Foundation
- Define project context schema and planner inference rules.
- Build deterministic context generation from idea input.

## Phase 2 - Document Generators
- Implement PRD/SPEC/ARCHITECTURE/API/DB generators.
- Validate output format and section completeness.

## Phase 3 - Scaffold Generator
- Create default repository structure.
- Add modules/generators/templates baseline folders.

## Phase 4 - CLI Workflow
- Add ` + "`new`" + `, ` + "`generate`" + `, and ` + "`autopilot`" + ` commands.
- Ensure each pipeline stage reruns independently.

## Phase 5 - Validation & Hardening
- Add tests for planner/generator/scaffold behavior.
- Validate docs quality and deterministic output.
- Improve errors and UX messages.
`,
	"REPOSITORY_STRUCTURE": `# REPOSITORY_STRUCTURE

## Recommended Repository Layout

` + "```text" + `
{{ .ProjectName }}/
├─ apps/
├─ services/
├─ libs/
├─ docs/
│  ├─ PRD.md
│  ├─ SPEC.md
│  ├─ ARCHITECTURE.md
│  ├─ API.md
│  ├─ DB_SCHEMA.md
│  ├─ IMPLEMENTATION_PLAN.md
│  └─ REPOSITORY_STRUCTURE.md
├─ infrastructure/
├─ modules/
│  ├─ planner/
│  ├─ generators/
│  └─ templates/
├─ generators/
│  ├─ prd_generator/
│  ├─ spec_generator/
│  ├─ architecture_generator/
│  ├─ api_generator/
│  ├─ db_generator/
│  └─ scaffold_generator/
└─ templates/
` + "```" + `
`,
}

func Render(name string, ctx model.ProjectContext) (string, error) {
	raw, ok := catalog[name]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", name)
	}

	t, err := template.New(name).Parse(raw)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}
