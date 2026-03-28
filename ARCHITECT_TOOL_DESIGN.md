# AI PROJECT ARCHITECT - Tool Design

## 1) System Architecture

Pipeline (rerunnable by stage):

idea -> planner -> prd -> spec -> architecture -> api -> database -> scaffold

Layers:
- CLI Layer: command routing (`init`, `prompt`, `list-skills`, `doctor`, `version`)
- Planner Layer: infer deterministic project context from idea
- Document Generator Layer: artifact-specific generators
- Template Layer: markdown templates for each artifact
- Scaffold Layer: repo skeleton generation

## 2) Repository Structure

```text
cmd/
internal/architect/
  model/
  modules/
    planner/
    generators/
    templates/
  generators/
    prd_generator/
    spec_generator/
    architecture_generator/
    api_generator/
    db_generator/
    implementation_plan_generator/
    repository_structure_generator/
  runtime/
  scaffold/
```

Generated project scaffold:

```text
/apps
/services
/libs
/docs
/infrastructure
/modules/planner
/modules/generators
/modules/templates
/generators/prd_generator
/generators/spec_generator
/generators/architecture_generator
/generators/api_generator
/generators/db_generator
/generators/scaffold_generator
/templates
```

## 3) Generator Module Design

Modules:
- Planner (`internal/architect/modules/planner`)
  - Parse idea and infer stack
  - Output normalized `ProjectContext`
- Generator Registry (`internal/architect/modules/generators`)
  - Resolve target artifacts (`docs`, `architecture`, `api`, ...)
  - Dispatch to specific generators
- Artifact Generators (`internal/architect/generators/*`)
  - Single responsibility per artifact
  - `Generate(context) -> markdown`
- Template Engine (`internal/architect/modules/templates`)
  - Deterministic template catalog
  - Shared render function

## 4) CLI Design

Commands:
- `apa init`
  - collect idea + project location
  - infer context
  - create scaffold, docs, and context file
- `apa prompt`
  - inspect the current repo state
  - output a structured prompt for iterative AI implementation
- `apa list-skills`
  - list repo-local skills under `./skills`
- `apa doctor`
  - validate environment assumptions such as Go availability and skills path
- `apa version`
  - show build metadata

## 5) Templates (PRD / SPEC / ARCHITECTURE / API)

Implemented template outputs:
- `PRD`: problem statement, goals, personas, core features, non-goals, success metrics
- `SPEC`: modules, functional/non-functional requirements, edge cases, constraints
- `ARCHITECTURE`: high-level architecture, boundaries, components, data flow, deployment, scalability
- `API`: OpenAPI-style endpoint design, request/response schema, auth, errors

Additional templates:
- `DB_SCHEMA`
- `IMPLEMENTATION_PLAN`
- `REPOSITORY_STRUCTURE`
