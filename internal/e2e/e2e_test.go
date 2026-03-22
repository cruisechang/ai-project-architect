package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var testBinaryPath string

func TestMain(m *testing.M) {
	root := repoRoot()
	binaryPath := filepath.Join(os.TempDir(), fmt.Sprintf("project-generator-e2e-%d", time.Now().UnixNano()))
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = root
	buildCmd.Env = append(os.Environ(), "GOCACHE=/tmp/go-cache", "GOPATH=/tmp/go")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build project-generator test binary: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	defer os.Remove(binaryPath)

	testBinaryPath = binaryPath
	code := m.Run()
	os.Exit(code)
}

func TestDoctorCommand(t *testing.T) {
	out, err := runProjgen(t, nil, "doctor", "--check-write=true")
	if err != nil {
		t.Fatalf("doctor command failed: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "[PASS] Go executable")
	assertContains(t, out, "[PASS] Go version")
	assertContains(t, out, "[PASS] Skills path")
}

func TestVersionCommandDefaultBuildInfo(t *testing.T) {
	out, err := runProjgen(t, nil, "version")
	if err != nil {
		t.Fatalf("version command failed: %v\noutput:\n%s", err, out)
	}
	assertContains(t, out, "version=dev")
	assertContains(t, out, "commit=none")
	assertContains(t, out, "date=unknown")
}

func TestVersionCommandWithInjectedBuildInfo(t *testing.T) {
	root := repoRoot()
	tmp := t.TempDir()
	customBin := filepath.Join(tmp, "project-generator-custom")
	if runtime.GOOS == "windows" {
		customBin += ".exe"
	}

	ldflags := "-X project-generator/internal/buildinfo.Version=v9.9.9 -X project-generator/internal/buildinfo.Commit=abc1234 -X project-generator/internal/buildinfo.BuildDate=2026-03-08T00:00:00Z"
	buildCmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", customBin, ".")
	buildCmd.Dir = root
	buildCmd.Env = append(os.Environ(), "GOCACHE=/tmp/go-cache", "GOPATH=/tmp/go")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("custom build failed: %v\noutput:\n%s", err, string(out))
	}

	versionCmd := exec.Command(customBin, "version")
	versionCmd.Dir = root
	out, err := versionCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("running custom binary failed: %v\noutput:\n%s", err, string(out))
	}
	text := string(out)
	assertContains(t, text, "version=v9.9.9")
	assertContains(t, text, "commit=abc1234")
	assertContains(t, text, "date=2026-03-08T00:00:00Z")
}

func TestListSkillsCommand(t *testing.T) {
	skillsRoot := t.TempDir()
	mustMkdirAll(t, filepath.Join(skillsRoot, "planner"))
	mustMkdirAll(t, filepath.Join(skillsRoot, "implementer"))
	mustWriteFile(t, filepath.Join(skillsRoot, "README.txt"), "file")

	out, err := runProjgen(t, nil, "list-skills", "--path", skillsRoot)
	if err != nil {
		t.Fatalf("list-skills command failed: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "- planner")
	assertContains(t, out, "- implementer")
	assertNotContains(t, out, "README.txt")
}

// TestInitCommandCodexStack verifies that apa init creates a full project with
// code scaffold, design docs, and context when --idea is provided.
func TestInitCommandCodexStack(t *testing.T) {
	root := t.TempDir()
	globalSkills := filepath.Join(root, "global-skills")
	mustMkdirAll(t, filepath.Join(globalSkills, "planner"))
	mustWriteFile(t, filepath.Join(globalSkills, "planner", "SKILL.md"), "# Planner")

	projectParent := filepath.Join(root, "projects")
	out, err := runProjgen(t, nil,
		"init",
		"--idea", "Build an AI RAG system with Go backend and Next.js frontend",
		"--name", "demo-project",
		"--path", projectParent,
		"--type", "ai-app",
		"--ai-feature", "rag",
		"--agent", "codex",
		"--docs", "full",
		"--skills-path", globalSkills,
		"--skills", "planner",
		"--description", "demo desc",
		"--unit-test", "yes",
		"--integration-test", "yes",
		"--e2e-test", "yes",
		"--docker-compose", "yes",
	)
	if err != nil {
		t.Fatalf("init command failed: %v\noutput:\n%s", err, out)
	}

	projectRoot := filepath.Join(projectParent, "demo-project")

	// Code scaffold (from generator.CreateProject)
	assertPathExists(t, filepath.Join(projectRoot, "AGENTS.md"))
	assertPathExists(t, filepath.Join(projectRoot, "PROMPT.md"))
	assertPathExists(t, filepath.Join(projectRoot, "PLANS.md"))
	assertPathExists(t, filepath.Join(projectRoot, "Makefile"))
	assertPathExists(t, filepath.Join(projectRoot, "src"))
	assertPathExists(t, filepath.Join(projectRoot, "tests"))
	assertPathExists(t, filepath.Join(projectRoot, "tasks", "queue", "TASK_TEMPLATE.md"))
	assertPathExists(t, filepath.Join(projectRoot, "scripts"))
	assertPathExists(t, filepath.Join(projectRoot, "agents", "planner-agent", "AGENT.md"))
	assertPathExists(t, filepath.Join(projectRoot, "agents", "coder-agent", "AGENT.md"))
	assertPathExists(t, filepath.Join(projectRoot, "agents", "reviewer-agent", "AGENT.md"))
	assertPathExists(t, filepath.Join(projectRoot, "agents", "test-agent", "AGENT.md"))
	assertPathExists(t, filepath.Join(projectRoot, ".codex", "config.json"))
	assertPathExists(t, filepath.Join(projectRoot, ".codex", "memory.md"))
	assertPathExists(t, filepath.Join(projectRoot, ".codex", "cache"))
	assertPathExists(t, filepath.Join(projectRoot, "skills", "apa-feature", "SKILL.md"))
	assertPathExists(t, filepath.Join(projectRoot, "skills", "planner", "SKILL.md"))

	// AI/RAG scaffold
	assertPathExists(t, filepath.Join(projectRoot, "ai", "prompts", "system_prompt.md"))
	assertPathExists(t, filepath.Join(projectRoot, "ai", "models", "model_config.yaml"))
	assertPathExists(t, filepath.Join(projectRoot, "ai", "evaluation", "eval_cases.yaml"))
	assertPathExists(t, filepath.Join(projectRoot, "ai", "embeddings"))
	assertPathExists(t, filepath.Join(projectRoot, "ai", "retrieval"))
	assertPathExists(t, filepath.Join(projectRoot, "ai", "ingestion"))
	assertPathExists(t, filepath.Join(projectRoot, "deployment", "docker", "Dockerfile"))
	assertPathExists(t, filepath.Join(projectRoot, "deployment", "compose", "docker-compose.yml"))

	// Design docs (from architect generators, overwriting template docs)
	assertPathExists(t, filepath.Join(projectRoot, "docs", "PRD.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "SPEC.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "ARCHITECTURE.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "API.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "DB_SCHEMA.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "IMPLEMENTATION_PLAN.md"))

	// Context for apa iterate
	assertPathExists(t, filepath.Join(projectRoot, ".architect", "context.json"))

	assertFileContains(t, filepath.Join(projectRoot, "AGENTS.md"), "AI Agent Engineering Guidelines")
	assertFileContains(t, filepath.Join(projectRoot, "PROMPT.md"), "AI Prompt Engineering Guidelines")
	assertContains(t, out, "=== INIT SUMMARY ===")
}

// TestInitCommandClaudeAgent verifies claude-code agent structure and no codex artifacts.
func TestInitCommandClaudeAgent(t *testing.T) {
	projectParent := t.TempDir()
	out, err := runProjgen(t, nil,
		"init",
		"--idea", "Internal tool for managing deployments with Python backend",
		"--name", "demo-claude",
		"--path", projectParent,
		"--type", "internal-tool",
		"--ai-feature", "prompt-workflow",
		"--agent", "claude-code",
		"--backend", "python",
		"--frontend", "none",
		"--docs", "basic",
		"--architecture", "backend-service",
		"--unit-test", "yes",
		"--integration-test", "yes",
		"--e2e-test", "no",
		"--docker-compose", "no",
	)
	if err != nil {
		t.Fatalf("init command failed: %v\noutput:\n%s", err, out)
	}

	projectRoot := filepath.Join(projectParent, "demo-claude")
	assertPathExists(t, filepath.Join(projectRoot, "CLAUDE.md"))
	assertPathExists(t, filepath.Join(projectRoot, "AGENTS.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "PRD.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "SPEC.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "ARCHITECTURE.md"))
	assertPathExists(t, filepath.Join(projectRoot, ".claude", "settings.json"))
	assertPathExists(t, filepath.Join(projectRoot, ".claude", "memory.md"))
	assertPathExists(t, filepath.Join(projectRoot, ".architect", "context.json"))
	assertPathNotExists(t, filepath.Join(projectRoot, ".codex"))
	assertPathNotExists(t, filepath.Join(projectRoot, "PROMPT.md"))
	assertPathNotExists(t, filepath.Join(projectRoot, "PLANS.md"))
	assertContains(t, out, "=== INIT SUMMARY ===")
}

// TestInitCommandSkillFailureIsGraceful verifies that a missing skill causes a
// non-zero exit with a descriptive message while the project is still created.
func TestInitCommandSkillFailureIsGraceful(t *testing.T) {
	root := t.TempDir()
	globalSkills := filepath.Join(root, "global-skills")
	mustMkdirAll(t, filepath.Join(globalSkills, "planner"))
	mustWriteFile(t, filepath.Join(globalSkills, "planner", "SKILL.md"), "# Planner")

	projectParent := filepath.Join(root, "projects")
	out, err := runProjgen(t, nil,
		"init",
		"--idea", "Simple web app",
		"--name", "demo-fail-skill",
		"--path", projectParent,
		"--type", "web-app",
		"--agent", "codex",
		"--backend", "go",
		"--frontend", "react",
		"--architecture", "frontend-backend",
		"--skills-path", globalSkills,
		"--skills", "planner,missing-skill",
		"--unit-test", "yes",
		"--integration-test", "yes",
		"--e2e-test", "no",
		"--docker-compose", "no",
	)
	if err == nil {
		t.Fatalf("expected init command to return non-zero when skill copy fails\noutput:\n%s", out)
	}

	projectRoot := filepath.Join(projectParent, "demo-fail-skill")
	assertPathExists(t, filepath.Join(projectRoot, "skills", "planner", "SKILL.md"))
	assertContains(t, out, "[FAIL] missing-skill: skill not found")
	assertContains(t, out, "init completed with 1 failed skill copy operations")
}

// TestInitCommandGeneratesBlueprint verifies the full init flow produces both
// code scaffold and idea-derived design docs in one command.
func TestInitCommandGeneratesBlueprint(t *testing.T) {
	parent := t.TempDir()
	out, err := runProjgen(t, nil,
		"init",
		"--name", "architect-demo",
		"--path", parent,
		"--idea", "Build a DevOps monitoring dashboard with React frontend, Go backend, and PostgreSQL.",
	)
	if err != nil {
		t.Fatalf("init command failed: %v\noutput:\n%s", err, out)
	}

	projectRoot := filepath.Join(parent, "architect-demo")

	// Idea-derived design docs
	assertPathExists(t, filepath.Join(projectRoot, ".architect", "context.json"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "PRD.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "SPEC.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "ARCHITECTURE.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "API.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "DB_SCHEMA.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "IMPLEMENTATION_PLAN.md"))

	// Code scaffold
	assertPathExists(t, filepath.Join(projectRoot, "Makefile"))
	assertPathExists(t, filepath.Join(projectRoot, "agents", "planner-agent", "AGENT.md"))

	assertContains(t, out, "INIT SUMMARY")
}

// TestInitCommandIdeaFile verifies that --idea-file works as the idea source.
func TestInitCommandIdeaFile(t *testing.T) {
	parent := t.TempDir()
	ideaFile := filepath.Join(t.TempDir(), "idea.txt")
	mustWriteFile(t, ideaFile, `你是一位資深雲端與 DevSecOps 架構師，請幫我規劃一個「DevOps 管理平台」的第一階段：資安掃描模組。

目標：
- 建立可落地的 MVP 設計，之後可擴充到監控、操作自動化、資安管理、VM 管理。
- 目前先聚焦「資安掃描」端到端流程。

技術偏好：Go backend, React frontend, PostgreSQL, Kubernetes 部署。`)

	out, err := runProjgen(t, nil,
		"init",
		"--name", "devsec-platform",
		"--path", parent,
		"--idea-file", ideaFile,
	)
	if err != nil {
		t.Fatalf("init --idea-file command failed: %v\noutput:\n%s", err, out)
	}

	projectRoot := filepath.Join(parent, "devsec-platform")
	assertPathExists(t, filepath.Join(projectRoot, ".architect", "context.json"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "PRD.md"))
	assertPathExists(t, filepath.Join(projectRoot, "docs", "ARCHITECTURE.md"))
	assertPathExists(t, filepath.Join(projectRoot, "Makefile"))
	assertContains(t, out, "INIT SUMMARY")
}

// TestIterateCommandNoContext verifies apa iterate works without a context.json.
func TestIterateCommandNoContext(t *testing.T) {
	dir := t.TempDir()
	out, err := runProjgen(t, nil, "iterate", "--root", dir)
	if err != nil {
		t.Fatalf("iterate command failed: %v\noutput:\n%s", err, out)
	}
	assertContains(t, out, "keep iterating until done")
	assertContains(t, out, "Workflow")
	assertContains(t, out, "Definition of Done")
	assertContains(t, out, "Constraints")
	assertContains(t, out, "context.json")
}

// TestIterateCommandWithContext verifies apa iterate embeds project info when context exists.
func TestIterateCommandWithContext(t *testing.T) {
	parent := t.TempDir()
	projectName := "iterate-ctx-demo"

	// Bootstrap a project so context.json exists
	_, err := runProjgen(t, nil,
		"init",
		"--name", projectName,
		"--path", parent,
		"--idea", "A monitoring dashboard with Go backend and React frontend.",
	)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	projectRoot := filepath.Join(parent, projectName)
	out, iterErr := runProjgen(t, nil, "iterate", "--root", projectRoot)
	if iterErr != nil {
		t.Fatalf("iterate command failed: %v\noutput:\n%s", iterErr, out)
	}
	assertContains(t, out, "keep iterating until done")
	assertContains(t, out, "Workflow")
	assertContains(t, out, "Definition of Done")
	assertContains(t, out, projectName)
}

func runProjgen(t *testing.T, env []string, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command(testBinaryPath, args...)
	cmd.Dir = repoRoot()
	cmd.Env = append(os.Environ(), env...)
	// Pipe stdin from an empty reader so the subprocess is never a TTY.
	// This ensures interactive prompts are skipped in all test environments.
	cmd.Stdin = bytes.NewReader(nil)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func repoRoot() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir failed for %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed for %s: %v", path, err)
	}
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path exists: %s (%v)", path, err)
	}
}

func assertPathNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected path NOT exists: %s", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("failed to stat %s: %v", path, err)
	}
}

func assertFileContains(t *testing.T, path, needle string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	if !strings.Contains(string(content), needle) {
		t.Fatalf("expected file %s to contain %q\ncontent:\n%s", path, needle, string(content))
	}
}

func assertContains(t *testing.T, output, needle string) {
	t.Helper()
	if !strings.Contains(output, needle) {
		t.Fatalf("expected output to contain %q\noutput:\n%s", needle, output)
	}
}

func assertNotContains(t *testing.T, output, needle string) {
	t.Helper()
	if strings.Contains(output, needle) {
		t.Fatalf("expected output NOT to contain %q\noutput:\n%s", needle, output)
	}
}
