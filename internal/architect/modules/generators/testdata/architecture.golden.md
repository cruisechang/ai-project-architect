# ARCHITECTURE

## Phase 0

### High-Level Architecture
The system follows a modular pipeline:

idea -> planner -> prd -> spec -> architecture -> api -> database -> scaffold

### Service Boundaries
- Planner: context extraction and stack inference.
- Generators: document-specific rendering units.
- Scaffold module: repository and folder bootstrap.

### Components
- CLI layer (`cmd/`): command routing and UX.
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
- Confirm architecture boundaries align with `PRD Phase 0` and `SPEC Phase 0`.
- Confirm generator responsibilities are explicit.

### Deployment Architecture
- Default target: docker-compose
- Can be adapted to container, Kubernetes, VM, or serverless.

### Scalability Considerations
- Add new generator without changing existing command contracts.
- Keep each artifact generation independently rerunnable.
- Keep context model versioned for compatibility.

### Completion Criteria
- Phase 0 architecture supports the committed blueprint workflow.

### Next-Phase Gate
- Phase 1 can start only after architecture verification passes, the module boundaries remain aligned with `PRD Phase 0` and `SPEC Phase 0`, and deployment assumptions are documented.

### Phase Completion Report
- Phase: `Phase 0`
- Goal:
  - Define the architecture required to support the highest-priority blueprint workflow.
- Delivered Scope:
  - Record the architecture delivered in Phase 0.
- Tests Executed:
  - Record pipeline, architecture, and integration checks run for Phase 0.
- Checks Performed:
  - Record service-boundary review, dependency inspection, and architecture validation completed.
- Completion Result:
  - `PASS`, `PARTIAL`, or `BLOCKED`
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record architecture invariants and deployment assumptions the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after architecture verification passes, the module boundaries remain aligned with `PRD Phase 0` and `SPEC Phase 0`, and deployment assumptions are documented.
- Next Phase:
  - `Phase 1` or `None`
