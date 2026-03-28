# Codex Review Skill (Claude Compatibility Alias)

Use this skill when the user asks Claude to run `apa-codex-review` for code review.

This skill is intentionally aligned with `apa-review` and should produce the same review quality bar.

## Inputs

- Diff or list of changed files
- Feature description / ticket context
- Test results

## Checklist

### Correctness
- [ ] Does the code do what the description says?
- [ ] Are edge cases handled (empty input, nulls, off-by-one, concurrent access)?
- [ ] Are errors handled explicitly — no silent swallowing?
- [ ] Is there a regression test for any bug fix?

### Design
- [ ] Is the change the simplest solution that solves the problem?
- [ ] Does it introduce unnecessary abstractions or over-engineering?
- [ ] Does it violate existing module boundaries?
- [ ] Are new dependencies justified?

### Readability
- [ ] Are names clear and intent-revealing?
- [ ] Are functions small (< 20 lines) and single-purpose?
- [ ] Is nesting depth <= 3 levels?
- [ ] Are comments explaining *why*, not *what*?

### Tests
- [ ] Do tests cover the happy path?
- [ ] Do tests cover key failure modes?
- [ ] Are tests isolated (no shared mutable state between tests)?
- [ ] Do tests run deterministically?

### Security
- [ ] No hardcoded secrets or credentials?
- [ ] User input validated before use?
- [ ] No SQL / command injection surface?
- [ ] Sensitive data not logged?

### Compatibility
- [ ] No breaking changes to public APIs without a major version bump?
- [ ] Backwards-compatible with existing data formats?

## Severity Labels

- **BLOCKER** — Must fix before merge (correctness, security)
- **MAJOR** — Should fix (design, test gaps)
- **MINOR** — Nice to fix (naming, style)
- **NIT** — Optional (micro-style, personal preference)

## Output

- Inline comments with severity labels
- Summary of blockers (if any) at the top
- Approval when all BLOCKERs and MAJORs are resolved
