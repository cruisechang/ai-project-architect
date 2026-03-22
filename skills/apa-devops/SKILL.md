# APA DevOps Skill

Use this skill when handling deployment, CI/CD, containerization, or infrastructure tasks.

## Inputs

- Target environment (dev / staging / prod)
- Deployment artifact (binary, container image, Helm chart)
- Required env vars and secrets

## Steps

### Containerization
1. Write a minimal `Dockerfile` (multi-stage build, non-root user, no dev deps in final image)
2. Test image locally: `docker build -t app:local . && docker run --rm app:local`
3. Verify image size and attack surface are reasonable

### CI/CD
1. Ensure pipeline runs: lint → test → build → (deploy)
2. Fail fast: lint and test before any build steps
3. Never deploy from a failing build
4. Pin external action versions (avoid `@main` or `@latest`)

### Deployment
1. Validate all required env vars exist before starting
2. Perform a health check after deploy before routing traffic
3. Keep rollback plan ready: know the previous stable tag
4. Drain connections gracefully on shutdown (`SIGTERM` handler)

### Secrets
- Never hardcode secrets in source or Dockerfiles
- Use env vars injected at runtime or a secrets manager
- Rotate any secret that may have been exposed immediately

## Output

- Passing CI pipeline
- Service reachable and health check green
- No secrets in version control
