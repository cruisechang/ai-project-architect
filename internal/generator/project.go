package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"project-generator/internal/config"
	"project-generator/internal/fsutil"
)

type Plan struct {
	ProjectRoot  string
	Questions    []string
	Dirs         []string
	Files        []string
	SkillsToCopy []string
}

type SkillCopyResult struct {
	Name    string
	Source  string
	Dest    string
	Success bool
	Error   string
}

type Result struct {
	ProjectRoot  string
	BackupPath   string
	CreatedDirs  []string
	CreatedFiles []string
	SkillCopies  []SkillCopyResult
	Warnings     []string
}

func (r Result) FailedSkillCopies() int {
	failed := 0
	for _, item := range r.SkillCopies {
		if !item.Success {
			failed++
		}
	}
	return failed
}

type plannedFile struct {
	RelPath      string
	TemplateName string
	RawContent   string
	Mode         os.FileMode
}

func BuildPlan(opts config.CreateOptions, questions []string) (Plan, error) {
	o := opts
	o.Normalize()
	o.EnsureDefaults()

	parent := o.ParentPath
	if parent == "" {
		parent = "."
	}
	parentAbs, err := config.ExpandPath(parent)
	if err != nil {
		return Plan{}, err
	}

	name := o.Name
	if name == "" {
		name = "<project-name>"
	}
	projectRoot := filepath.Join(parentAbs, name)

	dirs, files := buildArtifacts(o)
	plan := Plan{
		ProjectRoot: projectRoot,
		Questions:   questions,
		Dirs:        make([]string, 0, len(dirs)),
		Files:       make([]string, 0, len(files)),
	}

	for _, dir := range dirs {
		plan.Dirs = append(plan.Dirs, filepath.Join(projectRoot, dir))
	}
	for _, file := range files {
		plan.Files = append(plan.Files, filepath.Join(projectRoot, file.RelPath))
	}

	if o.GlobalSkillsPath != "" {
		for _, skill := range o.SelectedSkills {
			plan.SkillsToCopy = append(plan.SkillsToCopy,
				fmt.Sprintf("%s -> %s", filepath.Join(o.GlobalSkillsPath, skill), filepath.Join(projectRoot, "skills", skill)))
		}
	}
	if o.GlobalSkillsPath == "" && len(o.SelectedSkills) > 0 {
		for _, skill := range o.SelectedSkills {
			plan.SkillsToCopy = append(plan.SkillsToCopy, fmt.Sprintf("(missing skills source path) %s", skill))
		}
	}
	if len(o.SelectedSkills) == 0 {
		plan.SkillsToCopy = append(plan.SkillsToCopy, "(none)")
	}

	sort.Strings(plan.Dirs)
	sort.Strings(plan.Files)
	return plan, nil
}

func CreateProject(opts config.CreateOptions) (Result, error) {
	o := opts
	o.Normalize()
	o.EnsureDefaults()

	if err := o.ValidateForCreate(); err != nil {
		return Result{}, err
	}

	parentAbs, err := config.ExpandPath(o.ParentPath)
	if err != nil {
		return Result{}, err
	}

	skillsSourcePath := ""
	if o.GlobalSkillsPath != "" {
		skillsSourcePath, err = config.ExpandPath(o.GlobalSkillsPath)
		if err != nil {
			return Result{}, err
		}
		info, statErr := os.Stat(skillsSourcePath)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				return Result{}, fmt.Errorf("skills source path does not exist: %s", skillsSourcePath)
			}
			return Result{}, statErr
		}
		if !info.IsDir() {
			return Result{}, fmt.Errorf("skills source path is not a directory: %s", skillsSourcePath)
		}
	}

	if err := fsutil.EnsureDir(parentAbs, false); err != nil {
		return Result{}, err
	}

	result := Result{}
	result.ProjectRoot = filepath.Join(parentAbs, o.Name)

	exists, err := fsutil.Exists(result.ProjectRoot)
	if err != nil {
		return result, err
	}

	if exists {
		if !o.Overwrite {
			return result, fmt.Errorf("target project path already exists: %s", result.ProjectRoot)
		}
		backupPath := fmt.Sprintf("%s_backup_%s", result.ProjectRoot, time.Now().Format("20060102150405"))
		if err := os.Rename(result.ProjectRoot, backupPath); err != nil {
			return result, fmt.Errorf("cannot backup existing project directory: %w", err)
		}
		result.BackupPath = backupPath
	}

	if err := fsutil.EnsureDir(result.ProjectRoot, false); err != nil {
		return result, err
	}

	dirs, files := buildArtifacts(o)
	for _, dir := range dirs {
		absDir := filepath.Join(result.ProjectRoot, dir)
		if err := fsutil.EnsureDir(absDir, false); err != nil {
			return result, err
		}
		result.CreatedDirs = append(result.CreatedDirs, absDir)
	}

	data := buildTemplateData(o)
	for _, file := range files {
		absFile := filepath.Join(result.ProjectRoot, file.RelPath)
		content := file.RawContent
		if file.TemplateName != "" {
			rendered, renderErr := renderTemplate(file.TemplateName, data)
			if renderErr != nil {
				return result, fmt.Errorf("template render failed (%s): %w", file.TemplateName, renderErr)
			}
			content = rendered
		}

		if err := fsutil.WriteFile(absFile, []byte(content), file.Mode, false); err != nil {
			return result, err
		}
		result.CreatedFiles = append(result.CreatedFiles, absFile)
	}

	skillCopies := copySelectedSkills(result.ProjectRoot, skillsSourcePath, o.SelectedSkills)
	result.SkillCopies = append(result.SkillCopies, skillCopies...)
	for _, item := range skillCopies {
		if !item.Success {
			result.Warnings = append(result.Warnings, fmt.Sprintf("skill copy failed [%s]: %s", item.Name, item.Error))
		}
	}

	sort.Strings(result.CreatedDirs)
	sort.Strings(result.CreatedFiles)
	return result, nil
}

func buildArtifacts(opts config.CreateOptions) ([]string, []plannedFile) {
	agent := strings.TrimSpace(strings.ToLower(opts.AIAgent))
	if agent == "" {
		agent = "codex"
	}
	aiFeatureEnabled := strings.TrimSpace(strings.ToLower(opts.AIFeature)) != "" && strings.TrimSpace(strings.ToLower(opts.AIFeature)) != "none"
	agentsTemplate := "agents_md.tmpl"
	if aiFeatureEnabled {
		agentsTemplate = "agents_ai_md.tmpl"
	}

	dirs := []string{"docs", "src", "tests", "tests/unit", "tests/integration", "skills", "agents", "agents/coder-agent", "agents/planner-agent", "agents/reviewer-agent", "agents/test-agent", "tasks", "tasks/queue"}
	files := []plannedFile{
		{RelPath: "AGENTS.md", TemplateName: agentsTemplate, Mode: 0o644},
		{RelPath: "agents/coder-agent/AGENT.md", TemplateName: "coder_agent_md.tmpl", Mode: 0o644},
		{RelPath: "agents/planner-agent/AGENT.md", TemplateName: "planner_agent_md.tmpl", Mode: 0o644},
		{RelPath: "agents/reviewer-agent/AGENT.md", TemplateName: "reviewer_agent_md.tmpl", Mode: 0o644},
		{RelPath: "agents/test-agent/AGENT.md", TemplateName: "test_agent_md.tmpl", Mode: 0o644},
		{RelPath: "tasks/queue/TASK_TEMPLATE.md", TemplateName: "task_template_md.tmpl", Mode: 0o644},
		{RelPath: "README.md", TemplateName: "readme_md.tmpl", Mode: 0o644},
		{RelPath: "README.zh-TW.md", TemplateName: "readme_zh_tw_md.tmpl", Mode: 0o644},
		{RelPath: "README.zh-CN.md", TemplateName: "readme_zh_cn_md.tmpl", Mode: 0o644},
		{RelPath: "README.de.md", TemplateName: "readme_de_md.tmpl", Mode: 0o644},
		{RelPath: "README.fr.md", TemplateName: "readme_fr_md.tmpl", Mode: 0o644},
		{RelPath: "README.es.md", TemplateName: "readme_es_md.tmpl", Mode: 0o644},
		{RelPath: "README.ja.md", TemplateName: "readme_ja_md.tmpl", Mode: 0o644},
		{RelPath: "README.ko.md", TemplateName: "readme_ko_md.tmpl", Mode: 0o644},
		{RelPath: "Makefile", TemplateName: "makefile.tmpl", Mode: 0o644},
	}

	files = append(files, docsFiles(opts.DocsType)...)
	dirs = append(dirs, backendDirectories(opts.BackendType)...)
	dirs = append(dirs, frontendDirectories(opts.FrontendType)...)
	files = append(files, backendFiles(opts.BackendType)...)
	files = append(files, frontendFiles(opts.FrontendType)...)

	switch agent {
	case "claude-code":
		dirs = append(dirs,
			".claude",
			".claude/agents",
			".claude/commands",
			"scripts",
			"skills/apa-catalog",
			"skills/apa-debug",
			"skills/apa-devops",
			"skills/apa-docs",
			"skills/apa-feature",
			"skills/apa-implement",
			"skills/apa-integration",
			"skills/apa-loop",
			"skills/apa-review",
			"skills/apa-tdd",
		)
		files = append(files,
			plannedFile{RelPath: "CLAUDE.md", TemplateName: "claude_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".claude/settings.json", TemplateName: "claude_settings_json.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".claude/memory.md", TemplateName: "agent_memory_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".claude/commands/apa-loop.md", TemplateName: "apa_loop_command_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".claude/commands/cancel-apa-loop.md", TemplateName: "cancel_apa_loop_command_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "scripts/apa-loop-hook.sh", TemplateName: "apa_loop_hook_sh.tmpl", Mode: 0o755},
			plannedFile{RelPath: "scripts/apa-loop-setup.sh", TemplateName: "apa_loop_setup_sh.tmpl", Mode: 0o755},
			plannedFile{RelPath: "scripts/apa-loop-cancel.sh", TemplateName: "apa_loop_cancel_sh.tmpl", Mode: 0o755},
			plannedFile{RelPath: "skills/apa-catalog/SKILL.md", TemplateName: "apa_catalog_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-debug/SKILL.md", TemplateName: "apa_debug_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-devops/SKILL.md", TemplateName: "apa_devops_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-docs/SKILL.md", TemplateName: "apa_docs_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-feature/SKILL.md", TemplateName: "feature_development_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-implement/SKILL.md", TemplateName: "apa_implement_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-integration/SKILL.md", TemplateName: "apa_integration_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-loop/SKILL.md", TemplateName: "apa_loop_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-review/SKILL.md", TemplateName: "apa_review_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-tdd/SKILL.md", TemplateName: "apa_tdd_skill_md.tmpl", Mode: 0o644},
		)
	default:
		promptTemplate := "prompt_md.tmpl"
		if aiFeatureEnabled {
			promptTemplate = "prompt_ai_md.tmpl"
		}
		dirs = append(dirs,
			"scripts",
			".codex",
			".codex/cache",
			"skills/apa-catalog",
			"skills/apa-debug",
			"skills/apa-devops",
			"skills/apa-docs",
			"skills/apa-feature",
			"skills/apa-implement",
			"skills/apa-integration",
			"skills/apa-loop",
			"skills/apa-review",
			"skills/apa-tdd",
		)
		files = append(files,
			plannedFile{RelPath: "PROMPT.md", TemplateName: promptTemplate, Mode: 0o644},
			plannedFile{RelPath: "PLANS.md", TemplateName: "plans_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".codex/config.json", TemplateName: "codex_config_json.tmpl", Mode: 0o644},
			plannedFile{RelPath: ".codex/memory.md", TemplateName: "agent_memory_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-catalog/SKILL.md", TemplateName: "apa_catalog_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-debug/SKILL.md", TemplateName: "apa_debug_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-devops/SKILL.md", TemplateName: "apa_devops_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-docs/SKILL.md", TemplateName: "apa_docs_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-feature/SKILL.md", TemplateName: "feature_development_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-implement/SKILL.md", TemplateName: "apa_implement_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-integration/SKILL.md", TemplateName: "apa_integration_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-loop/SKILL.md", TemplateName: "apa_loop_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-review/SKILL.md", TemplateName: "apa_review_skill_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "skills/apa-tdd/SKILL.md", TemplateName: "apa_tdd_skill_md.tmpl", Mode: 0o644},
		)
	}

	dirs = append(dirs, aiFeatureGeneratedDirs(opts.AIFeature)...)
	files = append(files, aiFeatureFiles(opts.AIFeature)...)

	dirs = uniqStrings(dirs)
	files = uniqFiles(files)
	return dirs, files
}

func aiFeatureGeneratedDirs(aiFeature string) []string {
	feature := strings.TrimSpace(strings.ToLower(aiFeature))
	if feature == "" || feature == "none" {
		return nil
	}

	dirs := []string{
		"ai",
		"ai/prompts",
		"ai/models",
		"ai/evaluation",
		"api",
		"data",
		"data/documents",
		"deployment",
		"deployment/docker",
		"deployment/compose",
		"scripts",
		"tests/ai_tests",
		"tests/api_tests",
	}

	if feature == "rag" {
		dirs = append(dirs,
			"ai/embeddings",
			"ai/retrieval",
			"ai/ingestion",
		)
	}
	if feature == "agent-system" {
		dirs = append(dirs,
			"agent",
			"agent/planner",
			"agent/tools",
			"agent/memory",
			"agent/workflows",
		)
	}
	return dirs
}

func aiFeatureFiles(aiFeature string) []plannedFile {
	feature := strings.TrimSpace(strings.ToLower(aiFeature))
	if feature == "" || feature == "none" {
		return nil
	}

	files := []plannedFile{
		{RelPath: "ai/prompts/system_prompt.md", RawContent: aiSystemPromptContent(), Mode: 0o644},
		{RelPath: "ai/prompts/user_prompt_templates.md", RawContent: aiUserPromptTemplatesContent(), Mode: 0o644},
		{RelPath: "ai/prompts/prompt_examples.md", RawContent: aiPromptExamplesContent(), Mode: 0o644},
		{RelPath: "ai/models/model_config.yaml", RawContent: aiModelConfigContent(), Mode: 0o644},
		{RelPath: "ai/models/model_selection.md", RawContent: aiModelSelectionContent(), Mode: 0o644},
		{RelPath: "ai/evaluation/eval_cases.yaml", RawContent: aiEvalCasesContent(), Mode: 0o644},
		{RelPath: "ai/evaluation/evaluation_runner.py", RawContent: aiEvaluationRunnerContent(), Mode: 0o644},
		{RelPath: "scripts/run_eval.sh", RawContent: aiRunEvalScriptContent(), Mode: 0o755},
		{RelPath: "scripts/build_index.sh", RawContent: aiBuildIndexScriptContent(), Mode: 0o755},
		{RelPath: "deployment/docker/Dockerfile", RawContent: aiDockerfileContent(), Mode: 0o644},
		{RelPath: "deployment/compose/docker-compose.yml", RawContent: aiDockerComposeContent(), Mode: 0o644},
	}

	if feature == "rag" {
		files = append(files,
			plannedFile{RelPath: "ai/embeddings/embedding_pipeline.py", RawContent: aiEmbeddingPipelineContent(), Mode: 0o644},
			plannedFile{RelPath: "ai/embeddings/embedding_config.yaml", RawContent: aiEmbeddingConfigContent(), Mode: 0o644},
			plannedFile{RelPath: "ai/retrieval/retriever.py", RawContent: aiRetrieverContent(), Mode: 0o644},
			plannedFile{RelPath: "ai/retrieval/search_pipeline.py", RawContent: aiSearchPipelineContent(), Mode: 0o644},
			plannedFile{RelPath: "ai/ingestion/document_loader.py", RawContent: aiDocumentLoaderContent(), Mode: 0o644},
			plannedFile{RelPath: "ai/ingestion/index_builder.py", RawContent: aiIndexBuilderContent(), Mode: 0o644},
		)
	}

	if feature == "agent-system" {
		files = append(files,
			plannedFile{RelPath: "agent/planner/task_planner.py", RawContent: aiTaskPlannerContent(), Mode: 0o644},
			plannedFile{RelPath: "agent/tools/tool_registry.py", RawContent: aiToolRegistryContent(), Mode: 0o644},
			plannedFile{RelPath: "agent/tools/tool_examples.py", RawContent: aiToolExamplesContent(), Mode: 0o644},
			plannedFile{RelPath: "agent/memory/agent_memory.py", RawContent: aiAgentMemoryContent(), Mode: 0o644},
			plannedFile{RelPath: "agent/workflows/agent_loop.py", RawContent: aiAgentLoopContent(), Mode: 0o644},
		)
	}
	return files
}

func aiSystemPromptContent() string {
	return `# System Prompt

This project follows the AI Project Template Structure.

Rules:

- prompts must be placed in ai/prompts
- model configuration must be in ai/models
- evaluation cases must be in ai/evaluation
- RAG pipelines must be in ai/retrieval
- document ingestion must be in ai/ingestion
- agent logic must be in agent/
`
}

func aiUserPromptTemplatesContent() string {
	return `# User Prompt Templates

## template_1
- Goal:
- Input:
- Constraints:
- Expected output:

## template_2
- Goal:
- Input:
- Constraints:
- Expected output:
`
}

func aiPromptExamplesContent() string {
	return `# Prompt Examples

## Example 1
User:
Assistant:

## Example 2
User:
Assistant:
`
}

func aiModelConfigContent() string {
	return `default_model: gpt-4.1-mini
fallback_model: gpt-4.1
temperature: 0.2
max_tokens: 1200
`
}

func aiModelSelectionContent() string {
	return `# Model Selection

- Primary model:
- Fallback model:
- Selection rationale:
- Known tradeoffs:
`
}

func aiEvalCasesContent() string {
	return `cases:
  - id: basic_case
    input: ""
    expected: ""
`
}

func aiEvaluationRunnerContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement evaluation runner")
`
}

func aiEmbeddingPipelineContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement embedding pipeline")
`
}

func aiEmbeddingConfigContent() string {
	return `provider: openai
model: text-embedding-3-small
batch_size: 64
`
}

func aiRetrieverContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement retriever")
`
}

func aiSearchPipelineContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement search pipeline")
`
}

func aiDocumentLoaderContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement document loader")
`
}

func aiIndexBuilderContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement index builder")
`
}

func aiTaskPlannerContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement task planner")
`
}

func aiToolRegistryContent() string {
	return `TOOLS = []
`
}

func aiToolExamplesContent() string {
	return `# TODO: add tool usage examples
`
}

func aiAgentMemoryContent() string {
	return `# TODO: implement agent memory management
`
}

func aiAgentLoopContent() string {
	return `#!/usr/bin/env python3
print("TODO: implement agent loop")
`
}

func aiRunEvalScriptContent() string {
	return `#!/usr/bin/env bash
set -euo pipefail
python3 ai/evaluation/evaluation_runner.py
`
}

func aiBuildIndexScriptContent() string {
	return `#!/usr/bin/env bash
set -euo pipefail
python3 ai/ingestion/index_builder.py
`
}

func aiDockerfileContent() string {
	return `FROM python:3.12-slim
WORKDIR /app
COPY . .
CMD ["python3", "ai/evaluation/evaluation_runner.py"]
`
}

func aiDockerComposeContent() string {
	return `services:
  app:
    build:
      context: ..
      dockerfile: deployment/docker/Dockerfile
`
}

func copySelectedSkills(projectRoot, skillsSourcePath string, skills []string) []SkillCopyResult {
	results := make([]SkillCopyResult, 0, len(skills))
	for _, skill := range skills {
		skillName := strings.TrimSpace(skill)
		if skillName == "" {
			continue
		}

		entry := SkillCopyResult{Name: skillName}
		if skillsSourcePath == "" {
			entry.Success = false
			entry.Error = "skills source path is empty"
			results = append(results, entry)
			continue
		}

		src := filepath.Join(skillsSourcePath, skillName)
		dst := filepath.Join(projectRoot, "skills", skillName)
		entry.Source = src
		entry.Dest = dst

		info, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				entry.Success = false
				entry.Error = "skill not found"
				results = append(results, entry)
				continue
			}
			entry.Success = false
			entry.Error = err.Error()
			results = append(results, entry)
			continue
		}
		if !info.IsDir() {
			entry.Success = false
			entry.Error = "skill path is not a directory"
			results = append(results, entry)
			continue
		}

		if err := fsutil.CopyDir(src, dst); err != nil {
			entry.Success = false
			entry.Error = err.Error()
			results = append(results, entry)
			continue
		}

		entry.Success = true
		results = append(results, entry)
	}
	return results
}

func uniqStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func uniqFiles(values []plannedFile) []plannedFile {
	seen := map[string]struct{}{}
	out := make([]plannedFile, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value.RelPath]; ok {
			continue
		}
		seen[value.RelPath] = struct{}{}
		out = append(out, value)
	}
	return out
}
