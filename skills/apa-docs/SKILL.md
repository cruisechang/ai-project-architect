# APA Docs Skill

Use this skill when creating, refreshing, or expanding project documentation from the current repository state.

This skill supports two modes:
- full-doc mode: generate or update the complete documentation set
- targeted-doc mode: generate or update only the specified document

## Inputs

- Requested mode: `all-docs` or `single-doc`
- Requested document scope:
  - full set: `all docs`
  - single file: `PRD`, `SPEC`, `ARCHITECTURE`, `API`, `DB_SCHEMA`, or `IMPLEMENTATION_PLAN`
- Current repository structure and code
- Existing files under `docs/`
- Any product or technical constraints from the user

## Steps

1. **Read existing docs first** — Inspect `docs/` and avoid overwriting useful decisions blindly
2. **Read the code and structure** — Use the actual repository layout, entrypoints, tests, and config as evidence
3. **Identify gaps** — Separate what is explicitly supported by code from what is still planned or assumed
4. **Pick the target docs**
   - In `all-docs` mode, update the full documentation set under `docs/`
   - In `single-doc` mode, update only the requested document
5. **Write concrete content** — Prefer specific modules, flows, interfaces, and constraints over vague architecture prose
6. **Mark uncertainty clearly** — If a detail is inferred rather than implemented, label it as an assumption
7. **Keep docs aligned**
   - In `all-docs` mode, normalize terminology and system boundaries across all docs
   - In `single-doc` mode, keep neighboring docs in mind but do not rewrite them unless the user asks
8. **Verify after writing** — Re-read the changed docs and confirm they match the current repo state

## Rules

- Do not invent implemented features that are not present in the repo
- Do not rewrite every doc if only one document needs updating
- Prefer concise, operational documentation over generic template text
- Use current file paths, module names, and commands from the repo

## Document Targets

- `PRD` -> `docs/PRD.md`
- `SPEC` -> `docs/SPEC.md`
- `ARCHITECTURE` -> `docs/ARCHITECTURE.md`
- `API` -> `docs/API.md`
- `DB_SCHEMA` -> `docs/DB_SCHEMA.md`
- `IMPLEMENTATION_PLAN` -> `docs/IMPLEMENTATION_PLAN.md`

## Usage Notes

- If the user says "update docs" without narrowing scope, default to `all-docs`
- If the user names one document, treat it as `single-doc`
- If the requested document does not exist yet, create it in `docs/`

## Output

- Updated document files under `docs/`
- Either:
  - the full set of docs, or
  - only the requested target doc
- Assumptions called out explicitly where code evidence is missing
- Terminology and structure consistent with the current repository
