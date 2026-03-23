package apa

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	architectruntime "project-generator/internal/architect/runtime"
	"project-generator/internal/config"
	"project-generator/internal/docphase"
	"project-generator/internal/i18n"
)

func newIterateCmd() *cobra.Command {
	var root string

	cmd := &cobra.Command{
		Use:   "iterate",
		Short: i18n.T("iterate.short"),
		Long:  i18n.T("iterate.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if root == "" {
				root = "."
			}
			rootAbs, err := config.ExpandPath(root)
			if err != nil {
				return err
			}
			printIteratePrompt(rootAbs)
			return nil
		},
	}

	cmd.Flags().StringVar(&root, "root", "", i18n.T("iterate.flag.root"))
	return cmd
}

const iterateSep = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

func printIteratePrompt(rootAbs string) {
	ctx, ctxErr := architectruntime.LoadContext(rootAbs)
	hasCtx := ctxErr == nil

	fmt.Println(i18n.T("iterate.prompt.intro"))
	fmt.Println()

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.project-info"))
	fmt.Println(iterateSep)
	fmt.Printf("%s %s\n", i18n.T("iterate.prompt.root-label"), rootAbs)
	if hasCtx {
		if ctx.ProjectName != "" {
			fmt.Printf("%s %s\n", i18n.T("iterate.prompt.name-label"), ctx.ProjectName)
		}
		if ctx.ProjectIdea != "" {
			fmt.Printf("%s %s\n", i18n.T("iterate.prompt.idea-label"), ctx.ProjectIdea)
		}
		if stack := iterateTechStack(ctx.BackendLanguage, ctx.FrontendFramework, ctx.Database); stack != "" {
			fmt.Printf("%s %s\n", i18n.T("iterate.prompt.stack-label"), stack)
		}
	} else {
		fmt.Println(i18n.T("iterate.prompt.no-context"))
	}
	fmt.Println()

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.docs-status"))
	fmt.Println(iterateSep)
	docFiles := []string{
		"docs/PRD.md",
		"docs/SPEC.md",
		"docs/ARCHITECTURE.md",
		"docs/API.md",
		"docs/DB_SCHEMA.md",
		"docs/IMPLEMENTATION_PLAN.md",
	}
	anyDoc := false
	var nonPhaseDocs []string
	for _, rel := range docFiles {
		status, err := docphase.Check(rootAbs, rel)
		if err != nil {
			fmt.Printf("  [%s] %s (%v)\n", i18n.T("iterate.prompt.missing"), rel, err)
			continue
		}
		if status.Exists {
			fmt.Printf("  [%s] %s\n", i18n.T("iterate.prompt.exists"), rel)
			anyDoc = true
			if !status.PhaseBased {
				nonPhaseDocs = append(nonPhaseDocs, rel)
			}
		} else {
			fmt.Printf("  [%s] %s\n", i18n.T("iterate.prompt.missing"), rel)
		}
	}
	if !anyDoc {
		fmt.Println(i18n.T("iterate.prompt.no-docs"))
	}
	fmt.Println()

	if len(nonPhaseDocs) > 0 {
		fmt.Println(iterateSep)
		fmt.Println(i18n.T("iterate.prompt.phase-warning"))
		fmt.Println(iterateSep)
		fmt.Println(i18n.T("iterate.prompt.phase-warning-items"))
		for _, rel := range nonPhaseDocs {
			fmt.Printf("  - %s\n", rel)
		}
		fmt.Println(i18n.T("iterate.prompt.phase-warning-action"))
		fmt.Println()
	}

	tasksDir := filepath.Join(rootAbs, "tasks", "queue")
	if info, err := os.Stat(tasksDir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(tasksDir)
		if names := iterateMarkdownFiles(entries); len(names) > 0 {
			fmt.Println(iterateSep)
			fmt.Println(i18n.T("iterate.prompt.tasks"))
			fmt.Println(iterateSep)
			for _, name := range names {
				fmt.Printf("  tasks/queue/%s\n", name)
			}
			fmt.Println()
		}
	}

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.workflow"))
	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.workflow-steps"))
	fmt.Println()

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.done"))
	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.done-items"))
	fmt.Println()

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.constraints"))
	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.constraints-items"))
	fmt.Println()

	fmt.Println(iterateSep)
	fmt.Println(i18n.T("iterate.prompt.start"))
	if len(nonPhaseDocs) > 0 {
		fmt.Printf(i18n.T("iterate.prompt.start-items.phase"), strings.Join(nonPhaseDocs, ", "))
		fmt.Println()
	} else {
		fmt.Println(i18n.T("iterate.prompt.start-items"))
	}
	fmt.Println(iterateSep)
}

func iterateTechStack(backend, frontend, database string) string {
	var parts []string
	for _, s := range []string{backend, frontend, database} {
		if s != "" && s != "none" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " | ")
}

func iterateMarkdownFiles(entries []os.DirEntry) []string {
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, e.Name())
		}
	}
	return names
}
