# APA Doc Review Skill

Use this skill when the user wants to iteratively refine project documentation with the agent, round by round, and explicitly does not want implementation to start yet.

This skill is for interactive document review and revision. It is not the same as `apa-docs`:
- `apa-docs` generates or rewrites documentation
- `apa-doc-review` runs a feedback loop with the user until the docs are approved

## Goal

Improve and align project documents through repeated user feedback until the user explicitly approves the docs for implementation handoff.

## Scope

- `README*.md`
- `docs/PRD.md`
- `docs/SPEC.md`
- `docs/ARCHITECTURE.md`
- `docs/API.md`
- `docs/DB_SCHEMA.md`
- `docs/IMPLEMENTATION_PLAN.md`
- `docs/IMPLEMENTATION_STATUS.md` when needed for planning clarity

## Core Rules

- Do not implement code
- Do not modify application source files unless the user explicitly changes scope away from doc review
- Do not start `apa-loop`
- Do not use `apa-implement`
- After each document revision round, stop and wait for user feedback
- Do not transition into implementation until the user explicitly says `docs approved`

## Process

1. Read the current docs before proposing changes.
2. Identify gaps, ambiguity, conflicts, missing assumptions, and doc-to-doc inconsistency.
3. Propose the smallest useful doc revision for the current round.
4. Update only the relevant docs for that round.
5. Report:
   - what changed
   - what is still unclear
   - recommended next doc improvements
6. Stop and wait for feedback.

## Review Priorities

Review documents in this order unless the user directs otherwise:

1. Scope and priorities in `PRD`
2. Technical behavior and constraints in `SPEC`
3. Public contract shape in `API`
4. Data ownership and schema in `DB_SCHEMA`
5. System boundaries in `ARCHITECTURE`
6. Delivery sequencing in `IMPLEMENTATION_PLAN`
7. Operator-facing clarity in `README*.md`

## Alignment Rules

- Keep documentation phase-based where applicable (`Phase 0`, `Phase 1`, ...)
- Keep PRD, SPEC, API, DB schema, architecture, and implementation plan aligned
- If a statement is inferred rather than confirmed, label it as an assumption
- If a decision is deferred, mark it clearly instead of hiding the gap
- Prefer concrete wording over generic planning language

## Output Per Round

- Current doc issues
- Doc changes made this round
- Open questions or unresolved tradeoffs
- Recommended next review focus
- Explicit stop for user feedback

## Exit Condition

Only when the user explicitly says `docs approved` may the workflow transition into implementation. At that point, ask whether to switch to `apa-loop` + `apa-implement`, or wait for the user's next instruction.
