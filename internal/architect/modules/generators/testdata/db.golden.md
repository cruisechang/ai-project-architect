# DB_SCHEMA

## Phase 0

### Database
- Engine: postgresql

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

#### `artifacts`
- id (PK)
- project_id (FK -> projects.id)
- artifact_type
- file_path
- generated_at
- checksum

#### `generation_runs`
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
- Confirm entity design aligns with `SPEC Phase 0` and `API Phase 0`.

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
- Phase 1 can start only after schema validation and persistence checks pass, the schema remains aligned with `SPEC Phase 0` and `API Phase 0`, and unresolved data-model risks are explicitly documented.

### Phase Completion Report
- Phase: `Phase 0`
- Goal:
  - Define the data model required for the highest-priority generation workflow.
- Delivered Scope:
  - Record the schema delivered in Phase 0.
- Tests Executed:
  - Record schema validation and persistence checks run for Phase 0.
- Checks Performed:
  - Record entity review, relationship inspection, and constraint verification completed.
- Completion Result:
  - `PASS`, `PARTIAL`, or `BLOCKED`
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record schema constraints and data-model caveats the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after schema validation and persistence checks pass, the schema remains aligned with `SPEC Phase 0` and `API Phase 0`, and unresolved data-model risks are explicitly documented.
- Next Phase:
  - `Phase 1` or `None`
