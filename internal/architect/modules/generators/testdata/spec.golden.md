# SPEC

## Technical Scope
- Project: golden-project
- Project Type: fullstack
- Frontend Framework: react
- Backend: go / gin
- Database: postgresql

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
- Confirm the SPEC matches `PRD Phase 0` and `API Phase 0`.
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
- Phase 1 can start only after generator and integration tests pass, the SPEC remains aligned with `PRD Phase 0` and `API Phase 0`, and deferred risks are documented.

### Phase Completion Report
- Phase: `Phase 0`
- Goal:
  - Define the technical system baseline for the highest-priority delivery phase.
- Delivered Scope:
  - Record the technical scope completed in Phase 0.
- Tests Executed:
  - Record generator, integration, and verification tests run for Phase 0.
- Checks Performed:
  - Record architecture traceability checks, review outcomes, and technical inspections.
- Completion Result:
  - `PASS`, `PARTIAL`, or `BLOCKED`
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record technical constraints and decisions the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after generator and integration tests pass, the SPEC remains aligned with `PRD Phase 0` and `API Phase 0`, and deferred risks are documented.
- Next Phase:
  - `Phase 1` or `None`
