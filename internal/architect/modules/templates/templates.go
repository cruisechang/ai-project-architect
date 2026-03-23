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

## Phase 0

### Problem Statement
Users need a reliable way to execute the core idea: {{ .ProjectIdea }}.
Current workflows are fragmented, manual, or slow for repeated operations.

### Goals
- Deliver a production-ready foundation for {{ .ProjectName }}.
- Provide clear functional scope and engineering handoff.
- Enable incremental delivery with measurable outcomes.

### User Personas
- Product owner: defines priorities and validates business outcomes.
- End user: consumes the core product workflow.
- Operator/engineer: maintains system reliability and observability.

### Core Features
- Project planning and architecture baseline generation.
- Document generation: PRD, SPEC, ARCHITECTURE, API, DB schema, implementation plan.
- Repeatable scaffold generation aligned to architecture.

### Required Tests
- Validate the first end-to-end blueprint flow for the core product idea.
- Verify documentation outputs remain aligned across PRD, SPEC, API, and DB schema.

### Inspection Checklist
- Confirm Phase 0 scope is limited to the first shippable foundation.
- Confirm assumptions and non-goals are explicit.

### Non-Goals
- Building every feature in one release.
- Locking implementation to a single cloud/provider.

### Success Metrics
- Time to first blueprint generated.
- Documentation completeness score in review.
- Lead time reduction from idea to implementation start.

### Completion Criteria
- Phase 0 scope is fully documented, implemented, and validated.
- No blocking gap remains between product scope and technical handoff.

### Next-Phase Gate
- Phase 1 can start only after Phase 0 scope, tests, and PRD/SPEC/API alignment are all verified and open scope gaps are documented.

### Phase Completion Report
- Phase: ` + "`Phase 0`" + `
- Goal:
  - Deliver the highest-priority product foundation for the current project idea.
- Delivered Scope:
  - Record what Phase 0 delivered.
- Tests Executed:
  - Record the end-to-end, acceptance, and supporting tests run for Phase 0.
- Checks Performed:
  - Record document alignment checks, review steps, and manual verification completed.
- Completion Result:
  - ` + "`PASS`" + `, ` + "`PARTIAL`" + `, or ` + "`BLOCKED`" + `
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record product constraints and delivery context the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after Phase 0 scope, tests, and PRD/SPEC/API alignment are all verified and open scope gaps are documented.
- Next Phase:
  - ` + "`Phase 1`" + ` or ` + "`None`" + `
`,
	"SPEC": `# SPEC

## Technical Scope
- Project: {{ .ProjectName }}
- Project Type: {{ .ProjectType }}
- Frontend Framework: {{ .FrontendFramework }}
- Backend: {{ .BackendLanguage }} / {{ .BackendFramework }}
- Database: {{ .Database }}

## Phase 0

### System Modules
- Planner module: parse and normalize project requirements.
- Document generators: produce deterministic markdown artifacts.
- Scaffold generator: create repository baseline and folders.
- Template engine: render reusable templates into final docs.

### Feature Descriptions
- Parse project idea and infer stack deterministically.
- Generate engineering documents with fixed sections.
- Create implementation plan and scaffold for execution.

### Functional Requirements
- Must generate PRD.md, SPEC.md, ARCHITECTURE.md, API.md, DB_SCHEMA.md.
- Must generate IMPLEMENTATION_PLAN.md.
- Must support optional REPOSITORY_STRUCTURE.md output.
- Each generation step must be independently rerunnable.

### Required Tests
- Generator unit tests for deterministic outputs.
- Integration tests covering the full generation pipeline for Phase 0.

### Inspection Checklist
- Confirm the SPEC matches ` + "`PRD Phase 0`" + ` and ` + "`API Phase 0`" + `.
- Confirm module boundaries support rerunnable generation.

### Non-Functional Requirements
- Deterministic output for same input context.
- Modular implementation with isolated generators.
- CLI-first UX with scriptable commands.

### Edge Cases
- Missing project idea.
- Unsupported stack keywords.
- Existing target directory conflicts.

### Constraints
- Markdown output only.
- Context-driven generation with minimal hidden assumptions.

### Completion Criteria
- Phase 0 technical scope is implemented and verified.
- Output artifacts remain consistent across reruns.

### Next-Phase Gate
- Phase 1 can start only after generator and integration tests pass, the SPEC remains aligned with ` + "`PRD Phase 0`" + ` and ` + "`API Phase 0`" + `, and deferred risks are documented.

### Phase Completion Report
- Phase: ` + "`Phase 0`" + `
- Goal:
  - Define the technical system baseline for the highest-priority delivery phase.
- Delivered Scope:
  - Record the technical scope completed in Phase 0.
- Tests Executed:
  - Record generator, integration, and verification tests run for Phase 0.
- Checks Performed:
  - Record architecture traceability checks, review outcomes, and technical inspections.
- Completion Result:
  - ` + "`PASS`" + `, ` + "`PARTIAL`" + `, or ` + "`BLOCKED`" + `
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record technical constraints and decisions the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after generator and integration tests pass, the SPEC remains aligned with ` + "`PRD Phase 0`" + ` and ` + "`API Phase 0`" + `, and deferred risks are documented.
- Next Phase:
  - ` + "`Phase 1`" + ` or ` + "`None`" + `
`,
	"ARCHITECTURE": `# ARCHITECTURE

## Phase 0

### High-Level Architecture
The system follows a modular pipeline:

idea -> planner -> prd -> spec -> architecture -> api -> database -> scaffold

### Service Boundaries
- Planner: context extraction and stack inference.
- Generators: document-specific rendering units.
- Scaffold module: repository and folder bootstrap.

### Components
- CLI layer (` + "`cmd/`" + `): command routing and UX.
- Planning layer: project context model construction.
- Generator layer: PRD/SPEC/ARCHITECTURE/API/DB/PLAN outputs.
- Template layer: reusable deterministic markdown templates.

### Data Flow
1. User provides project idea.
2. Planner infers normalized context.
3. Generators produce docs from context.
4. Scaffold writer creates repository baseline.

### Required Tests
- Verify the pipeline works end to end for Phase 0 blueprint generation.
- Verify architecture boundaries hold under reruns and partial target generation.

### Inspection Checklist
- Confirm architecture boundaries align with ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `.
- Confirm generator responsibilities are explicit.

### Deployment Architecture
- Default target: {{ .Deployment }}
- Can be adapted to container, Kubernetes, VM, or serverless.

### Scalability Considerations
- Add new generator without changing existing command contracts.
- Keep each artifact generation independently rerunnable.
- Keep context model versioned for compatibility.

### Completion Criteria
- Phase 0 architecture supports the committed blueprint workflow.

### Next-Phase Gate
- Phase 1 can start only after architecture verification passes, the module boundaries remain aligned with ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `, and deployment assumptions are documented.

### Phase Completion Report
- Phase: ` + "`Phase 0`" + `
- Goal:
  - Define the architecture required to support the highest-priority blueprint workflow.
- Delivered Scope:
  - Record the architecture delivered in Phase 0.
- Tests Executed:
  - Record pipeline, architecture, and integration checks run for Phase 0.
- Checks Performed:
  - Record service-boundary review, dependency inspection, and architecture validation completed.
- Completion Result:
  - ` + "`PASS`" + `, ` + "`PARTIAL`" + `, or ` + "`BLOCKED`" + `
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record architecture invariants and deployment assumptions the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after architecture verification passes, the module boundaries remain aligned with ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `, and deployment assumptions are documented.
- Next Phase:
  - ` + "`Phase 1`" + ` or ` + "`None`" + `
`,
	"API": `# API

## Phase 0

### API Style
OpenAPI-style REST design (can be upgraded to strict OpenAPI YAML later).

### Base Path
` + "`/api/v1`" + `

### Authentication
- Mode: {{ .Authentication }}

### Endpoints
#### GET /health
- Purpose: service health check
- Response: { "status": "ok" }

#### POST /projects/plan
- Purpose: infer project context from idea
- Request:
  - idea: string
  - project_name: string
- Response:
  - project_context object

#### POST /projects/docs
- Purpose: generate architecture documents
- Request:
  - context: object
  - targets: string[]
- Response:
  - generated_files: string[]

#### POST /projects/scaffold
- Purpose: generate project scaffold
- Request:
  - context: object
  - root_path: string
- Response:
  - created_directories: string[]
  - created_files: string[]

### Required Tests
- API contract tests for each Phase 0 endpoint.
- Integration tests covering planner, docs, and scaffold flows.

### Inspection Checklist
- Confirm endpoint scope matches ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `.
- Confirm request and response fields cover the implemented workflow.

### Error Handling
- 400: invalid request
- 404: context not found
- 409: path conflict
- 500: generation failure

### Completion Criteria
- All Phase 0 endpoints are documented and verified.

### Next-Phase Gate
- Phase 1 can start only after contract and integration tests pass, the endpoint set remains aligned with ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `, and error behaviors are documented.

### Phase Completion Report
- Phase: ` + "`Phase 0`" + `
- Goal:
  - Define the highest-priority API surface required for the current workflow.
- Delivered Scope:
  - Record the API surface delivered in Phase 0.
- Tests Executed:
  - Record contract, integration, and error-handling tests run for Phase 0.
- Checks Performed:
  - Record request/response review, endpoint inspection, and API validation completed.
- Completion Result:
  - ` + "`PASS`" + `, ` + "`PARTIAL`" + `, or ` + "`BLOCKED`" + `
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record API constraints and behavior contracts the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after contract and integration tests pass, the endpoint set remains aligned with ` + "`PRD Phase 0`" + ` and ` + "`SPEC Phase 0`" + `, and error behaviors are documented.
- Next Phase:
  - ` + "`Phase 1`" + ` or ` + "`None`" + `
`,
	"DB_SCHEMA": `# DB_SCHEMA

## Phase 0

### Database
- Engine: {{ .Database }}

### Tables
#### projects
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

#### ` + "`artifacts`" + `
- id (PK)
- project_id (FK -> projects.id)
- artifact_type
- file_path
- generated_at
- checksum

#### ` + "`generation_runs`" + `
- id (PK)
- project_id (FK -> projects.id)
- run_mode (new|generate|autopilot)
- status
- started_at
- finished_at

### Required Tests
- Schema validation for Phase 0 entities and relations.
- Persistence tests for project, artifact, and generation run records.

### Inspection Checklist
- Confirm the schema supports only the committed Phase 0 workflow.
- Confirm entity design aligns with ` + "`SPEC Phase 0`" + ` and ` + "`API Phase 0`" + `.

### Relationships
- projects 1:N artifacts
- projects 1:N generation_runs

### Indexes
- idx_projects_name
- idx_artifacts_project_type
- idx_runs_project_status

### Constraints
- projects.name unique
- artifacts.artifact_type + file_path unique per project

### Completion Criteria
- Phase 0 schema supports the required generation workflow and passes validation.

### Next-Phase Gate
- Phase 1 can start only after schema validation and persistence checks pass, the schema remains aligned with ` + "`SPEC Phase 0`" + ` and ` + "`API Phase 0`" + `, and unresolved data-model risks are explicitly documented.

### Phase Completion Report
- Phase: ` + "`Phase 0`" + `
- Goal:
  - Define the data model required for the highest-priority generation workflow.
- Delivered Scope:
  - Record the schema delivered in Phase 0.
- Tests Executed:
  - Record schema validation and persistence checks run for Phase 0.
- Checks Performed:
  - Record entity review, relationship inspection, and constraint verification completed.
- Completion Result:
  - ` + "`PASS`" + `, ` + "`PARTIAL`" + `, or ` + "`BLOCKED`" + `
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record schema constraints and data-model caveats the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after schema validation and persistence checks pass, the schema remains aligned with ` + "`SPEC Phase 0`" + ` and ` + "`API Phase 0`" + `, and unresolved data-model risks are explicitly documented.
- Next Phase:
  - ` + "`Phase 1`" + ` or ` + "`None`" + `
`,
	"IMPLEMENTATION_PLAN": `# IMPLEMENTATION_PLAN

## Phase 0 - Planning Foundation
- Define project context schema and planner inference rules.
- Build deterministic context generation from idea input.

## Phase 1 - Document Generators
- Implement PRD/SPEC/ARCHITECTURE/API/DB generators.
- Validate output format and section completeness.

## Phase 2 - Scaffold Generator
- Create default repository structure.
- Add modules/generators/templates baseline folders.

## Phase 3 - CLI Workflow
- Add ` + "`new`" + `, ` + "`generate`" + `, and ` + "`autopilot`" + ` commands.
- Ensure each pipeline stage reruns independently.

## Phase 4 - Validation & Hardening
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
