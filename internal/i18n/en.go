package i18n

func init() {
	register("en", map[string]string{
		// root
		"root.short": "CLI from product idea to running project",
		"root.long": `apa (AI Project Architect)

` + iterSep + `
Recommended workflow
` + iterSep + `

  [First time] (run outside the project directory)
  apa init --idea "online food ordering platform" --name food-platform --path ~/projects

  [Inside the repo, let skills guide the agent]
  cd ~/projects/food-platform
  apa list-skills              # list available skills
  apa prompt                   # output 'keep iterating until done' AI prompt
                              # tell the agent to use apa-loop + apa-implement
  apa prompt --docs-only       # output doc-review prompt before implementation starts
                              # tell the agent to use apa-doc-review only
                              # optional wrapper in hook/slash environments:
                              # /apa-loop --max-iterations 30 --reviewer agent-self
                              # /cancel-apa-loop
  make test                    # run tests (repo-native Makefile)

` + iterSep + `
apa init — first bootstrap (core command)
` + iterSep + `

  # Interactive mode (prompts for idea / name / path / tech options)
  apa init

  # Non-interactive mode (--idea auto-infers tech stack)
  apa init --idea "online food ordering platform" --name food-platform --path ~/projects

  # Advanced: override inferred values
  apa init --idea "online food ordering platform" --name food-platform --backend python --agent universal

  # Preview files to be created (no actual writes)
  apa init --dry-run

  Output:
    .architect/context.json   inferred tech stack
    docs/                     PRD, SPEC, ARCHITECTURE, API, DB Schema, implementation plan
    backend/ frontend/ tests/ runnable starter files
    Makefile                  make test / make lint / make build
    CLAUDE.md or PROMPT.md    agent config
    agents/ skills/           agent and skill templates

` + iterSep + `
Other commands
` + iterSep + `

  apa prompt      Output 'keep iterating until done' AI prompt (run at any stage)
  apa prompt --docs-only  Output doc-only review prompt before implementation
  apa list-skills  List available skills
  apa doctor       Check environment
  apa version      Show version

Run "apa <command> --help" for details on each command.`,

		// init
		"init.short": "Bootstrap a new project (context + docs + runnable scaffold in one step)",
		"init.long": `Create a complete new project from an idea.

Steps:
  [1/4] Infer tech stack from idea (any field can be overridden via flags)
  [2/4] Generate runnable code scaffold (Makefile / test setup / agent config)
  [3/4] Generate design docs (PRD / SPEC / ARCHITECTURE / API / DB Schema / implementation plan)
  [4/4] Done

Usage:
  First time (outside repo)  →  apa init
  Inside repo                →  apa prompt → tell the agent to use apa-loop + apa-implement → make test (repeat)

Output:
  .architect/context.json          inferred tech stack
  docs/                            PRD / SPEC / ARCHITECTURE / API / DB Schema / implementation plan
  backend/ frontend/ tests/        runnable starter files
  Makefile                         make test / make lint / make build
  AGENTS.md + PROMPT.md / CLAUDE.md agent config (depends on --agent)
  agents/ skills/                  agent and repo-local skill templates`,

		"init.flag.idea":             "Product idea (triggers auto-inference of tech stack)",
		"init.flag.idea-file":        "Read idea from file (use - to read from stdin)",
		"init.flag.name":             "Project name (inferred from idea if omitted)",
		"init.flag.path":             "Parent directory path for the project",
		"init.flag.type":             "Type: cli | server | web-app-server | mobile-app-server | web-app | mobile-app",
		"init.flag.ai-feature":       "AI feature: none | prompt-workflow | rag | agent-system",
		"init.flag.agent":            "AI agent type: codex | claude-code | universal",
		"init.flag.backend":          "Backend language: go | python | node | none",
		"init.flag.frontend":         "Frontend framework: react | next | nuxt | vue | pure-typescript | none",
		"init.flag.stack":            "Tech stack description",
		"init.flag.docs":             "Docs type: basic | full",
		"init.flag.unit-test":        "Add unit tests: yes | no",
		"init.flag.api-test":         "Add API tests: yes | no",
		"init.flag.integration-test": "Add integration tests: yes | no",
		"init.flag.e2e-test":         "Add E2E tests: yes | no",
		"init.flag.docker-compose":   "Add Docker Compose: yes | no",
		"init.flag.skills":           "Skill names to copy, comma-separated",
		"init.flag.skills-path":      "Skills source directory (default: ./skills in current repo)",
		"init.flag.description":      "Project description",
		"init.flag.force":            "Backup existing directory and overwrite",
		"init.flag.dry-run":          "Preview files and directories to be created without writing",

		// prompt
		"prompt.short": "Output a 'keep iterating until done' AI prompt (can be run at any stage)",
		"prompt.long": `Reads current repo state and outputs a 'keep iterating until done' prompt for AI.

Paste the output into an AI and pair it with ` + "`apa-loop`" + ` when you want an enforced delivery loop.
Use the same primary path in both Codex and Claude Code repos: explicitly ask the agent to use the ` + "`apa-loop`" + ` skill together with ` + "`apa-implement`" + `.
If the environment also exposes generated slash commands or hooks, treat them as optional wrappers around the same workflow.
The AI will:
  1. Inventory current state (docs, tasks, tests, CI status)
  2. Loop with ` + "`apa-loop`" + `: pick 1-3 tasks → RED → GREEN → REFACTOR → validate → update docs/status
  3. Until the completion gate is met and the agent can output ` + "`<promise>COMPLETE</promise>`" + `

Usage:
  apa prompt                        # run in current repo directory
  apa prompt --docs-only            # generate a doc-review prompt before implementation
  apa prompt --reviewer agent-self  # generate implementation prompt with persisted reviewer
  apa prompt --root ~/projects/foo  # specify project directory
  apa prompt | pbcopy               # copy to clipboard (macOS)
  apa prompt > prompt.md            # save to file`,

		"prompt.flag.root":        "Project root path (default: current directory)",
		"prompt.flag.docs-only":   "Generate a doc-review prompt that uses `apa-doc-review` and blocks implementation",
		"prompt.flag.reviewer":    "Reviewer to persist in the generated implementation prompt: agent-self | apa-codex-review | apa-claude-review",
		"prompt.mode.label":       "Prompt mode? (implementation/docs-only)",
		"prompt.mode.default":     "implementation",
		"prompt.mode.invalid":     "invalid prompt mode %q (use implementation or docs-only)",
		"prompt.reviewer.label":   "Reviewer? (agent-self/apa-codex-review/apa-claude-review)",
		"prompt.reviewer.default": "agent-self",
		"prompt.reviewer.invalid": "invalid reviewer %q (use agent-self, apa-codex-review, or apa-claude-review)",

		// prompt output
		"prompt.output.intro":                "You are the primary implementation AI for this repo. Enter 'keep iterating until done' mode. Use the `apa-loop` and `apa-implement` skills. Execute directly — do not just give suggestions.",
		"prompt.output.project-info":         "Project Info",
		"prompt.output.root-label":           "Project root:",
		"prompt.output.reviewer-label":       "Reviewer:",
		"prompt.output.name-label":           "Name:",
		"prompt.output.idea-label":           "Idea:",
		"prompt.output.stack-label":          "Tech stack:",
		"prompt.output.no-context":           "(apa init has not been run — .architect/context.json not found)",
		"prompt.output.docs-status":          "Design Doc Status",
		"prompt.output.exists":               "FOUND",
		"prompt.output.missing":              "MISSING",
		"prompt.output.no-docs":              "(Run apa init to create a full project with design docs)",
		"prompt.output.phase-warning":        "Phase Rewrite Required",
		"prompt.output.phase-warning-items":  "The following existing docs are not written in priority-aligned `Phase 0`, `Phase 1`, ... format:",
		"prompt.output.phase-warning-action": "Call the `apa-docs` skill first and rewrite those docs into aligned phases before continuing implementation.",
		"prompt.output.tasks":                "Task Queue",
		"prompt.output.workflow":             "Workflow (loop until DONE)",
		"prompt.output.workflow-steps": `  1. Start with an inventory: read docs, tasks, tests, CI status. List incomplete items and risks.
  2. Use the ` + "`apa-loop`" + ` and ` + "`apa-implement`" + ` skills as the primary workflow.
  3. Each round: tackle 1–3 highest-priority tasks measured by verifiable outcomes.
  4. Follow RED → GREEN → REFACTOR. Start with a failing test or equivalent executable verification before implementation.
  5. After implementation, run required checks (at minimum: tests; run lint if available).
  6. Fix failures until they pass — never stop at half-done.
  7. Use reviewer ` + "`%s`" + ` for every round unless I explicitly switch it.
  8. Update docs and repo-local status (` + "`docs/IMPLEMENTATION_STATUS.md`" + ` or ` + "`TASKS.md`" + `) with completed work, in-progress work, failing checks, blockers, assumptions, and next 1–3 tasks.
  9. Move immediately to the next round — do not wait for my confirmation unless you hit a blocking decision.`,
		"prompt.output.done": "Definition of Done (ALL must be satisfied)",
		"prompt.output.done-items": `  - All documented P0/core requirements implemented and consistent with docs.
  - All documented core flows runnable.
  - All tests / build / lint / required checks passing.
  - No blocking errors and no known high-severity unresolved issues.
  - Docs and repo-local status updated to handoff-ready state.
  - Only when the completion gate is fully met, output ` + "`<promise>COMPLETE</promise>`" + `.`,
		"prompt.output.constraints": "Constraints",
		"prompt.output.constraints-items": `  - Do not delete or revert content I have not asked to remove.
  - Do not stop after a single task or partial feature.
  - Default to direct execution when docs are sufficient; only ask when blocked by a genuinely unsafe or irreversible decision.
  - For major trade-offs, present "options + recommendation + impact" first; otherwise just do it.`,
		"prompt.output.start":             "Begin now. First output:",
		"prompt.output.start-items":       "  A. Current inventory\n  B. The 1–3 tasks for round one\n  C. Confirm that reviewer `%s` is active for this loop\nThen enter the apa-loop implementation cycle and continue until the completion gate is fully met.",
		"prompt.output.start-items.phase": "  A. Use the `apa-docs` skill to rewrite these docs into aligned phases: %s\n  B. Current inventory after the doc rewrite\n  C. The 1–3 tasks for round one\n  D. Confirm that reviewer `%s` is active for this loop\nThen enter the apa-loop implementation cycle and continue until the completion gate is fully met.",
		"prompt.output.docs-only.intro":   "You are the documentation review AI for this repo. Use the `apa-doc-review` skill only. Do not implement code. Revise the docs directly and stop after each revision for feedback.",
		"prompt.output.docs-only.workflow-steps": `  1. Start with a doc inventory: read README, PRD, SPEC, API, DB schema, architecture, and implementation plan.
  2. Use the ` + "`apa-doc-review`" + ` skill as the primary workflow.
  3. In each round, identify the smallest useful doc revision that improves clarity, alignment, or scope control.
  4. Update only the relevant docs for that round — do not change implementation code.
  5. After each revision, summarize what changed, what is still unclear, and the recommended next doc focus.
  6. Stop and wait for my feedback after every round.
  7. Do not start ` + "`apa-loop`" + ` or ` + "`apa-implement`" + ` unless I explicitly say ` + "`docs approved`" + `.`,
		"prompt.output.docs-only.done-items": `  - PRD, SPEC, API, DB schema, architecture, README, and implementation plan are aligned where applicable.
  - Open assumptions, unresolved trade-offs, and deferred decisions are explicit instead of hidden.
  - The docs are clear enough for implementation handoff.
  - Only when I explicitly say ` + "`docs approved`" + ` may you ask whether to switch into ` + "`apa-loop`" + ` + ` + "`apa-implement`" + `.`,
		"prompt.output.docs-only.constraints-items": `  - Do not write or change implementation code.
  - Do not start ` + "`apa-loop`" + `.
  - Do not use ` + "`apa-implement`" + `.
  - Keep revisions small and reviewable.
  - If something is uncertain, label it as an assumption or open question instead of inventing certainty.`,
		"prompt.output.docs-only.start-items":       "  A. Current doc issues and inconsistencies\n  B. The smallest useful doc changes for this round\nThen revise the docs and stop for feedback.",
		"prompt.output.docs-only.start-items.phase": "  A. Use the `apa-docs` skill to rewrite these docs into aligned phases first: %s\n  B. Current doc issues and inconsistencies after the rewrite\n  C. The smallest useful doc changes for this round\nThen revise the docs and stop for feedback.",

		// list-skills
		"list-skills.short": "List available skills in the specified directory",
		"list-skills.long": `Lists all skill subdirectories under the skills directory.

A skill is a feature module directory that can be copied into a project via "apa init --skills".
Defaults to ./skills in the current repo when --path is not specified.`,

		"list-skills.flag.path": "Skills directory path (default: ./skills)",

		// doctor
		"doctor.short": "Check environment (Go version, write permission, skills path)",
		"doctor.long": `Check environment conditions required by apa.

Checks:
  Go executable    Confirm go command exists in PATH
  Go version       Show installed Go version
  Filesystem write Confirm write permission in current directory
  Skills path      Verify local skills directory exists (optional)

Output format:
  [PASS] check name: detail
  [FAIL] check name: detail`,

		"doctor.flag.skills-path": "Skills directory path to verify (optional, default: ./skills)",
		"doctor.flag.check-write": "Whether to check write permission in current directory",

		// version
		"version.short": "Show version info (version / commit / build date)",

		// architect_helpers
		"prompt.multiline.hint": "press Enter twice or Ctrl+D to finish",

		// wizard / interactive descriptions
		"wizard.project-type.frontend-only": "Frontend only: UI-focused app with no backend service in this repo",
		"wizard.project-type.backend-only":  "Backend only: API/service/CLI focused project",
		"wizard.project-type.full-stack":    "Full stack: frontend + backend in one project",

		"wizard.ai-feature.none":            "No AI functionality",
		"wizard.ai-feature.prompt-workflow": "Simple LLM app using prompt templates",
		"wizard.ai-feature.rag":             "Retrieval Augmented Generation system",
		"wizard.ai-feature.agent-system":    "Autonomous AI agent with tool usage and planning",

		"wizard.ai-agent.codex":       "OpenAI Codex style project structure",
		"wizard.ai-agent.claude-code": "Claude Code style project structure",
		"wizard.ai-agent.universal":   "Generate both Codex and Claude Code wrappers on top of shared project files",

		"wizard.architecture.cli":               "CLI",
		"wizard.architecture.server":            "Server",
		"wizard.architecture.web-app-server":    "Web app + server",
		"wizard.architecture.mobile-app-server": "Mobile app + server",
		"wizard.architecture.web-app":           "Web app",
		"wizard.architecture.mobile-app":        "Mobile app",

		"wizard.fullstack.next": "React fullstack framework",
		"wizard.fullstack.nuxt": "Vue fullstack framework",
	})
}

const iterSep = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
