# IMPLEMENTATION_PLAN

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
- Add `new`, `generate`, and `autopilot` commands.
- Ensure each pipeline stage reruns independently.

## Phase 4 - Validation & Hardening
- Add tests for planner/generator/scaffold behavior.
- Validate docs quality and deterministic output.
- Improve errors and UX messages.
