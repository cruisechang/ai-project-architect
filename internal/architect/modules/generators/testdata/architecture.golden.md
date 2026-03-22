# ARCHITECTURE

## High-Level Architecture
The system follows a modular pipeline:

idea -> planner -> prd -> spec -> architecture -> api -> database -> scaffold

## Service Boundaries
- Planner: context extraction and stack inference.
- Generators: document-specific rendering units.
- Scaffold module: repository and folder bootstrap.

## Components
- CLI layer (`cmd/`): command routing and UX.
- Planning layer: project context model construction.
- Generator layer: PRD/SPEC/ARCHITECTURE/API/DB/PLAN outputs.
- Template layer: reusable deterministic markdown templates.

## Data Flow
1. User provides project idea.
2. Planner infers normalized context.
3. Generators produce docs from context.
4. Scaffold writer creates repository baseline.

## Deployment Architecture
- Default target: docker-compose
- Can be adapted to container, Kubernetes, VM, or serverless.

## Scalability Considerations
- Add new generator without changing existing command contracts.
- Keep each artifact generation independently rerunnable.
- Keep context model versioned for compatibility.
