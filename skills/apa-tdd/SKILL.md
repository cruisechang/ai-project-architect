# APA TDD Skill

Use this skill when adding any non-trivial functionality. Tests come first — always.

## The Cycle

```
RED → GREEN → REFACTOR
```

1. **RED** — Write a failing test that describes the desired behavior
2. **GREEN** — Write the minimal code to make the test pass (no more, no less)
3. **REFACTOR** — Clean up without changing behavior; run tests again

## Steps

1. **Identify the behavior** — State in one sentence what the function/method must do
2. **Write the test first** — Name it `Test<What>_<When>_<Expected>` or equivalent
3. **Run it** — Confirm it fails (`RED`); if it passes, the test is wrong
4. **Implement minimally** — No speculative code, no "while I'm here" additions
5. **Run it again** — Confirm it passes (`GREEN`)
6. **Refactor** — Remove duplication, improve naming, simplify logic
7. **Run again** — All tests still green
8. **Repeat** for the next behavior

## Coverage Targets

- Unit tests: all pure logic functions
- Integration tests: service/DB/API boundaries
- E2E tests: critical user journeys only
- Minimum: 80% line coverage; focus on branch coverage for business logic

## Rules

- Never write implementation before the test
- One assertion per logical concept (can be multiple lines)
- Tests must be deterministic — no random, no clock dependencies without injection
- Mocks only at system boundaries (DB, HTTP, filesystem); not between internal modules

## Output

- Tests pass: `make test`
- Coverage at or above 80%: `make coverage`
- No production code without a corresponding test
