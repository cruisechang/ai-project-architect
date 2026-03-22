# APA Feature Skill

Use this skill when implementing a new feature end-to-end.

## Inputs

- Feature description or ticket
- Acceptance criteria
- Relevant existing code paths

## Steps

1. **Clarify requirements** — Restate the feature in your own words and confirm scope
2. **Define acceptance criteria** — List concrete, testable conditions for "done"
3. **Identify edge cases** — What inputs or states could break this?
4. **Impact analysis** — Which existing modules / APIs does this touch?
5. **Design solution** — Sketch the data flow and interfaces before coding
6. **Write tests first** — Unit tests for core logic, integration tests for boundaries (TDD)
7. **Implement** — Write the minimal code to make tests pass
8. **Refactor** — Clean up duplication, enforce naming conventions, remove dead code
9. **Update documentation** — Inline comments where logic is non-obvious, update relevant docs/
10. **Verify release safety** — No breaking API changes, all tests green, no debug prints

## Output

- Passing tests (`make test` or equivalent)
- Updated docs if public API changed
- Clean diff with no unrelated changes
