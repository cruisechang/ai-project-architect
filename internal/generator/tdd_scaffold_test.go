package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"project-generator/internal/config"
)

func createTestProject(t *testing.T, opts config.CreateOptions) Result {
	t.Helper()
	dir := t.TempDir()
	opts.ParentPath = dir
	opts.Normalize()
	opts.EnsureDefaults()
	result, err := CreateProject(opts)
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	return result
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file %s to exist", path)
	}
}

func assertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("file %s missing %q", filepath.Base(path), substr)
	}
}

func TestGoBackendIncludesTestFile(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-go", ProjectType: "web-app", BackendType: "go",
		FrontendType: "none", Architecture: "backend-service",
		TechStack: "go", DockerCompose: "no",
	})

	testFile := filepath.Join(result.ProjectRoot, "backend", "cmd", "server", "main_test.go")
	assertFileExists(t, testFile)
	assertFileContains(t, testFile, "testing")
	assertFileContains(t, testFile, "httptest")
}

func TestPythonBackendIncludesTestFile(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-py", ProjectType: "web-app", BackendType: "python",
		FrontendType: "none", Architecture: "backend-service",
		TechStack: "python", DockerCompose: "no",
	})

	testFile := filepath.Join(result.ProjectRoot, "backend", "test_main.py")
	assertFileExists(t, testFile)
	assertFileContains(t, testFile, "def test_")

	reqFile := filepath.Join(result.ProjectRoot, "backend", "requirements.txt")
	assertFileContains(t, reqFile, "pytest")
}

func TestNodeBackendIncludesTestFile(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-node", ProjectType: "web-app", BackendType: "node",
		FrontendType: "none", Architecture: "backend-service",
		TechStack: "node", DockerCompose: "no",
	})

	testFile := filepath.Join(result.ProjectRoot, "backend", "src", "index.test.js")
	assertFileExists(t, testFile)
	assertFileContains(t, testFile, "test(")

	pkgFile := filepath.Join(result.ProjectRoot, "backend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestReactFrontendIncludesTestDeps(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-react", ProjectType: "web-app", BackendType: "none",
		FrontendType: "react", Architecture: "frontend-app",
		TechStack: "react", DockerCompose: "no",
	})

	pkgFile := filepath.Join(result.ProjectRoot, "frontend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestNextFrontendIncludesTestDeps(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-next", ProjectType: "web-app", BackendType: "none",
		FrontendType: "next", Architecture: "frontend-app",
		TechStack: "next", DockerCompose: "no",
	})

	pkgFile := filepath.Join(result.ProjectRoot, "frontend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestVueFrontendIncludesTestDeps(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-vue", ProjectType: "web-app", BackendType: "none",
		FrontendType: "vue", Architecture: "frontend-app",
		TechStack: "vue", DockerCompose: "no",
	})

	pkgFile := filepath.Join(result.ProjectRoot, "frontend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestNuxtFrontendIncludesTestDeps(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-nuxt", ProjectType: "web-app", BackendType: "none",
		FrontendType: "nuxt", Architecture: "frontend-app",
		TechStack: "nuxt", DockerCompose: "no",
	})

	pkgFile := filepath.Join(result.ProjectRoot, "frontend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestTypeScriptFrontendIncludesTestDeps(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-ts", ProjectType: "web-app", BackendType: "none",
		FrontendType: "pure-typescript", Architecture: "frontend-app",
		TechStack: "typescript", DockerCompose: "no",
	})

	pkgFile := filepath.Join(result.ProjectRoot, "frontend", "package.json")
	assertFileContains(t, pkgFile, "vitest")
	assertFileContains(t, pkgFile, "\"test\"")
}

func TestMakefileHasRealTestCommands(t *testing.T) {
	// The Makefile uses native toolchain commands for each backend.
	// All backends produce the same Makefile content.
	tests := []struct {
		name    string
		backend string
	}{
		{"go", "go"},
		{"python", "python"},
		{"node", "node"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createTestProject(t, config.CreateOptions{
				Name: "tdd-mk-" + tt.name, ProjectType: "web-app",
				BackendType: tt.backend, FrontendType: "none",
				Architecture: "backend-service", TechStack: tt.backend,
				DockerCompose: "no",
			})
			mkFile := filepath.Join(result.ProjectRoot, "Makefile")
			assertFileExists(t, mkFile)
			assertFileContains(t, mkFile, "test:")
		})
	}
}

// requiredReadmeFiles lists every README the policy mandates.
// README-en.md must never appear.
var requiredReadmeFiles = []string{
	"README.md",
	"README.zh-CN.md",
	"README.zh-TW.md",
	"README.de.md",
	"README.fr.md",
	"README.es.md",
	"README.ja.md",
	"README.ko.md",
}

func newMLProject(t *testing.T, name string) Result {
	t.Helper()
	return createTestProject(t, config.CreateOptions{
		Name: name, ProjectType: "web-app", BackendType: "go",
		FrontendType: "none", Architecture: "backend-service",
		TechStack: "go", DockerCompose: "no",
	})
}

// TestMultilingualREADMEsExist verifies every required README file is generated.
func TestMultilingualREADMEsExist(t *testing.T) {
	result := newMLProject(t, "ml-exist")
	for _, f := range requiredReadmeFiles {
		assertFileExists(t, filepath.Join(result.ProjectRoot, f))
	}
}

// TestREADMEEnNotGenerated enforces the policy: README-en.md must never be created.
func TestREADMEEnNotGenerated(t *testing.T) {
	result := newMLProject(t, "ml-no-en")
	forbidden := []string{"README-en.md", "README.en.md"}
	for _, f := range forbidden {
		path := filepath.Join(result.ProjectRoot, f)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("policy violation: %s must not be generated", f)
		}
	}
}

// TestREADMETableSwitcherPresent verifies every README uses the table-format
// language switcher (not an inline format) and links to README.md for English.
func TestREADMETableSwitcherPresent(t *testing.T) {
	result := newMLProject(t, "ml-table")
	for _, f := range requiredReadmeFiles {
		path := filepath.Join(result.ProjectRoot, f)
		// Table format always contains a pipe-delimited row pointing to README.md
		assertFileContains(t, path, "| [README.md](README.md)")
		// All files must link to every other language
		assertFileContains(t, path, "README.zh-CN.md")
		assertFileContains(t, path, "README.ko.md")
	}
}

// TestLocalizedREADMEsAreNotEnglishDuplicates checks that each non-English
// README contains language-specific marker strings that cannot appear in the
// English README, proving the prose is genuinely translated.
func TestLocalizedREADMEsAreNotEnglishDuplicates(t *testing.T) {
	result := newMLProject(t, "ml-lang")

	// Each entry: (file, unique-to-that-language marker, absent-in-english)
	checks := []struct {
		file   string
		marker string // must be present in the localized file
	}{
		// Traditional Chinese: "快速開始" = "Quick Start"
		{"README.zh-TW.md", "快速開始"},
		// Simplified Chinese: "快速开始" = "Quick Start"
		{"README.zh-CN.md", "快速开始"},
		// German: "Schnellstart" = "Quick Start"
		{"README.de.md", "Schnellstart"},
		// French: "Démarrage rapide" = "Quick Start"
		{"README.fr.md", "Démarrage rapide"},
		// Spanish: "Inicio rápido" = "Quick Start"
		{"README.es.md", "Inicio rápido"},
		// Japanese: "クイックスタート" = "Quick Start"
		{"README.ja.md", "クイックスタート"},
		// Korean: "빠른 시작" = "Quick Start"
		{"README.ko.md", "빠른 시작"},
	}

	englishPath := filepath.Join(result.ProjectRoot, "README.md")

	for _, c := range checks {
		path := filepath.Join(result.ProjectRoot, c.file)
		// 1. Localized file must contain the native-language marker.
		assertFileContains(t, path, c.marker)
		// 2. English README must NOT contain that native-language marker,
		//    confirming the marker is unique to the translated file.
		data, err := os.ReadFile(englishPath)
		if err != nil {
			t.Fatalf("cannot read README.md: %v", err)
		}
		if strings.Contains(string(data), c.marker) {
			t.Errorf("README.md must not contain %q (a %s-specific string)", c.marker, c.file)
		}
	}
}

// TestREADMEMDIsEnglish verifies the canonical README.md is English
// by checking for English-only section headings.
func TestREADMEMDIsEnglish(t *testing.T) {
	result := newMLProject(t, "ml-en-canon")
	path := filepath.Join(result.ProjectRoot, "README.md")
	assertFileContains(t, path, "Quick Start")
	assertFileContains(t, path, "Development Commands")
	assertFileContains(t, path, "Development Principles")
}

func TestTestDirectoryStructure(t *testing.T) {
	result := createTestProject(t, config.CreateOptions{
		Name: "tdd-dirs", ProjectType: "web-app", BackendType: "go",
		FrontendType: "react", Architecture: "fullstack-web-app",
		TechStack: "go+react", DockerCompose: "no",
	})

	dirs := []string{"tests", "tests/unit", "tests/integration"}
	for _, d := range dirs {
		path := filepath.Join(result.ProjectRoot, d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected dir %s to exist", d)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
	}
}
