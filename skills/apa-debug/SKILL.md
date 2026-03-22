# APA Debug Skill

Use this skill when diagnosing and fixing a bug or unexpected behavior.

## Inputs

- Bug report or failing test output
- Reproduction steps
- Relevant stack trace or logs

## Steps

1. **Reproduce** — Confirm you can reproduce the failure locally before touching anything
2. **Isolate** — Narrow down to the smallest failing case (unit test or minimal repro)
3. **Read the error** — Parse the full stack trace; find the first frame in your own code
4. **Form a hypothesis** — State clearly what you believe is wrong before looking at code
5. **Verify hypothesis** — Add a targeted log or assertion to confirm/deny
6. **Fix** — Make the minimal change that addresses root cause, not symptoms
7. **Regression test** — Add a test that would have caught this bug
8. **Verify no side effects** — Run the full test suite (`make test`)
9. **Clean up** — Remove any debug logs or temporary assertions added during investigation

## Rules

- Never change multiple things at once — one variable per experiment
- If the hypothesis is wrong, discard it fully and re-read the evidence
- Fix the root cause, not the symptom

## Output

- Root cause documented in commit message
- Regression test added
- All tests passing
