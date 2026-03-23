# API

## Phase 0

### API Style
OpenAPI-style REST design (can be upgraded to strict OpenAPI YAML later).

### Base Path
`/api/v1`

### Authentication
- Mode: jwt

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
- Confirm endpoint scope matches `PRD Phase 0` and `SPEC Phase 0`.
- Confirm request and response fields cover the implemented workflow.

### Error Handling
- 400: invalid request
- 404: context not found
- 409: path conflict
- 500: generation failure

### Completion Criteria
- All Phase 0 endpoints are documented and verified.

### Next-Phase Gate
- Phase 1 can start only after contract and integration tests pass, the endpoint set remains aligned with `PRD Phase 0` and `SPEC Phase 0`, and error behaviors are documented.

### Phase Completion Report
- Phase: `Phase 0`
- Goal:
  - Define the highest-priority API surface required for the current workflow.
- Delivered Scope:
  - Record the API surface delivered in Phase 0.
- Tests Executed:
  - Record contract, integration, and error-handling tests run for Phase 0.
- Checks Performed:
  - Record request/response review, endpoint inspection, and API validation completed.
- Completion Result:
  - `PASS`, `PARTIAL`, or `BLOCKED`
- Open Issues / Risks:
  - None
- Assumptions:
  - None
- Handoff Notes:
  - Record API constraints and behavior contracts the next phase must preserve.
- Next-Phase Gate:
  - Phase 1 can start only after contract and integration tests pass, the endpoint set remains aligned with `PRD Phase 0` and `SPEC Phase 0`, and error behaviors are documented.
- Next Phase:
  - `Phase 1` or `None`
