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
  apa iterate                  # output 'keep iterating until done' AI prompt
                              # use apa-loop + apa-implement to enforce the round-based delivery loop
                              # /apa-loop --max-iterations 30
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
  apa init --idea "online food ordering platform" --name food-platform --backend python --agent claude-code

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

  apa iterate      Output 'keep iterating until done' AI prompt (run at any stage)
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
  Inside repo                →  apa iterate → apa-loop + apa-implement guide agent → make test (repeat)

Output:
  .architect/context.json          inferred tech stack
  docs/                            PRD / SPEC / ARCHITECTURE / API / DB Schema / implementation plan
  backend/ frontend/ tests/        runnable starter files
  Makefile                         make test / make lint / make build
  CLAUDE.md or PROMPT.md           agent config (depends on --agent)
  agents/ skills/                  agent and repo-local skill templates`,

		"init.flag.idea":             "Product idea (triggers auto-inference of tech stack)",
		"init.flag.idea-file":        "Read idea from file (use - to read from stdin)",
		"init.flag.name":             "Project name (inferred from idea if omitted)",
		"init.flag.path":             "Parent directory path for the project",
		"init.flag.type":             "Project type: web-app | ai-app | devops-tool | internal-tool | platform-service",
		"init.flag.ai-feature":       "AI feature: none | prompt-workflow | rag | agent-system",
		"init.flag.agent":            "AI agent type: codex | claude-code",
		"init.flag.backend":          "Backend language: go | python | node | none",
		"init.flag.frontend":         "Frontend framework: react | next | nuxt | vue | pure-typescript | none",
		"init.flag.architecture":     "Architecture: cli-tool | backend-service | frontend-app | fullstack-web-app | frontend-backend",
		"init.flag.stack":            "Tech stack description",
		"init.flag.docs":             "Docs type: basic | full",
		"init.flag.unit-test":        "Add unit tests: yes | no",
		"init.flag.integration-test": "Add integration tests: yes | no",
		"init.flag.e2e-test":         "Add E2E tests: yes | no",
		"init.flag.docker-compose":   "Add Docker Compose: yes | no",
		"init.flag.skills":           "Skill names to copy, comma-separated",
		"init.flag.skills-path":      "Skills source directory (default: ./skills in current repo)",
		"init.flag.description":      "Project description",
		"init.flag.force":            "Backup existing directory and overwrite",
		"init.flag.dry-run":          "Preview files and directories to be created without writing",

		// iterate
		"iterate.short": "Output a 'keep iterating until done' AI prompt (can be run at any stage)",
		"iterate.long": `Reads current repo state and outputs a 'keep iterating until done' prompt for AI.

Paste the output into an AI (e.g. Claude Code) and pair it with ` + "`apa-loop`" + ` when you want an enforced delivery loop.
The AI will:
  1. Inventory current state (docs, tasks, tests, CI status)
  2. Loop: implement → test → fix → update docs
  3. Until all core requirements are done and tests pass

Usage:
  apa iterate                        # run in current repo directory
  apa iterate --root ~/projects/foo  # specify project directory
  apa iterate | pbcopy               # copy to clipboard (macOS)
  apa iterate > prompt.md            # save to file`,

		"iterate.flag.root": "Project root path (default: current directory)",

		// iterate prompt output
		"iterate.prompt.intro":                "You are the primary implementation AI for this repo. Enter 'keep iterating until done' mode. Execute directly — do not just give suggestions.",
		"iterate.prompt.project-info":         "Project Info",
		"iterate.prompt.root-label":           "Project root:",
		"iterate.prompt.name-label":           "Name:",
		"iterate.prompt.idea-label":           "Idea:",
		"iterate.prompt.stack-label":          "Tech stack:",
		"iterate.prompt.no-context":           "(apa init has not been run — .architect/context.json not found)",
		"iterate.prompt.docs-status":          "Design Doc Status",
		"iterate.prompt.exists":               "FOUND",
		"iterate.prompt.missing":              "MISSING",
		"iterate.prompt.no-docs":              "(Run apa init to create a full project with design docs)",
		"iterate.prompt.phase-warning":        "Phase Rewrite Required",
		"iterate.prompt.phase-warning-items":  "The following existing docs are not written in priority-aligned `Phase 0`, `Phase 1`, ... format:",
		"iterate.prompt.phase-warning-action": "Call the `apa-docs` skill first and rewrite those docs into aligned phases before continuing implementation.",
		"iterate.prompt.tasks":                "Task Queue",
		"iterate.prompt.workflow":             "Workflow (loop until DONE)",
		"iterate.prompt.workflow-steps": `  1. Start with an inventory: read docs, tasks, tests, CI status. List incomplete items and risks.
  2. Each round: tackle 1–3 highest-priority tasks (measure by verifiable outcomes).
  3. After implementation, run required checks (at minimum: tests; run lint if available).
  4. Fix failures until they pass — never stop at half-done.
  5. Update docs and task status (what you did, what's left, what's next).
  6. Move immediately to the next round — do not wait for my confirmation unless you hit a blocking decision.`,
		"iterate.prompt.done": "Definition of Done (ALL must be satisfied)",
		"iterate.prompt.done-items": `  - All core requirements implemented and consistent with docs.
  - All tests passing (including new or fixed tests).
  - No blocking errors, no known P0/P1 issues.
  - Docs updated to handoff-ready state (README / how-to / design decisions / TODOs).
  - Final delivery summary and remaining risks (if any) listed.`,
		"iterate.prompt.constraints": "Constraints",
		"iterate.prompt.constraints-items": `  - Do not delete or revert content I have not asked to remove.
  - Each round: small commits of working results first, then expand next round.
  - For major trade-offs, present "options + recommendation + impact" first; otherwise just do it.`,
		"iterate.prompt.start":             "Begin now. First output:",
		"iterate.prompt.start-items":       "  A. Current inventory\n  B. The 1–3 tasks for round one\nThen enter the implementation loop until DONE.",
		"iterate.prompt.start-items.phase": "  A. Use the `apa-docs` skill to rewrite these docs into aligned phases: %s\n  B. Current inventory after the doc rewrite\n  C. The 1–3 tasks for round one\nThen enter the implementation loop until DONE.",

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
		"wizard.project-type.web-app":          "User-facing web app: UI, API, auth, user flow",
		"wizard.project-type.ai-app":           "AI app: prompt, model workflow, tool use, evaluation",
		"wizard.project-type.devops-tool":      "DevOps tool: CI/CD, deployment, rollback, automation, security ops",
		"wizard.project-type.internal-tool":    "Internal tool: process efficiency, ops support, data management, onboarding",
		"wizard.project-type.platform-service": "Platform service: shared capabilities, infra services, observability, security governance",

		"wizard.ai-feature.none":            "No AI functionality",
		"wizard.ai-feature.prompt-workflow": "Simple LLM app using prompt templates",
		"wizard.ai-feature.rag":             "Retrieval Augmented Generation system",
		"wizard.ai-feature.agent-system":    "Autonomous AI agent with tool usage and planning",

		"wizard.ai-agent.codex":       "OpenAI Codex style project structure",
		"wizard.ai-agent.claude-code": "Claude Code style project structure",

		"wizard.architecture.cli-tool":          "CLI tool",
		"wizard.architecture.backend-service":   "Backend service",
		"wizard.architecture.frontend-app":      "Frontend app",
		"wizard.architecture.fullstack-web-app": "Fullstack web app",
		"wizard.architecture.frontend-backend":  "Frontend + backend",

		"wizard.fullstack.next": "React fullstack framework",
		"wizard.fullstack.nuxt": "Vue fullstack framework",
	})
}

const iterSep = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
