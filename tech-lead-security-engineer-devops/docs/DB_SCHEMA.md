# DB_SCHEMA

## Database
- Engine: postgresql

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
