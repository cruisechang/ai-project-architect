# APA Integration Skill

Use this skill when integrating with an external service, third-party API, or internal service boundary.

## Inputs

- API spec or SDK documentation URL
- Auth mechanism (API key, OAuth, mTLS, etc.)
- Expected request/response shapes
- SLA / rate limit constraints

## Steps

1. **Read the spec** — Understand auth, pagination, error codes, and rate limits before writing code
2. **Define the interface** — Create an interface/protocol that your code depends on (not the SDK directly)
3. **Write a thin adapter** — Implement the interface against the external API; keep it minimal
4. **Handle errors explicitly** — Map external error codes to your domain errors; never let raw HTTP errors leak
5. **Add retries with backoff** — For transient failures (5xx, network timeouts); use exponential backoff
6. **Respect rate limits** — Honor `Retry-After` headers; add a client-side rate limiter if needed
7. **Write integration tests** — Use a test double or recorded fixture; verify the happy path and key error cases
8. **Set timeouts** — Every outbound call must have a deadline; never use the zero-value context
9. **Log request IDs** — Capture correlation IDs for debugging cross-service failures

## Rules

- Accept interfaces, return concrete types
- Integration code lives in a dedicated package (`internal/infra/`, `adapters/`, etc.)
- No business logic in adapters
- All secrets come from environment variables

## Output

- Working adapter with integration tests
- Error handling covers auth failure, rate limit, and timeout
- No API keys in source code
