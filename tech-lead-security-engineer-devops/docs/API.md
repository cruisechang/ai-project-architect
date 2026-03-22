# API

## API Style
OpenAPI-style REST design (can be upgraded to strict OpenAPI YAML later).

## Base Path
`/api/v1`

## Authentication
- Mode: jwt

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
