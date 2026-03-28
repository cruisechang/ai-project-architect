package apa

import (
	"bufio"
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

func newPromptCmd() *cobra.Command {
	var root string
	var docsOnly bool

	cmd := &cobra.Command{
		Use:   "prompt",
		Short: i18n.T("prompt.short"),
		Long:  i18n.T("prompt.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if root == "" {
				root = "."
			}
			resolvedDocsOnly, err := resolvePromptDocsOnly(cmd.Flags().Changed("docs-only"), docsOnly)
			if err != nil {
				return err
			}
			rootAbs, err := config.ExpandPath(root)
			if err != nil {
				return err
			}
			printPromptOutput(rootAbs, resolvedDocsOnly)
			return nil
		},
	}

	cmd.Flags().StringVar(&root, "root", "", i18n.T("prompt.flag.root"))
	cmd.Flags().BoolVar(&docsOnly, "docs-only", false, i18n.T("prompt.flag.docs-only"))
	return cmd
}

const promptSep = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

func printPromptOutput(rootAbs string, docsOnly bool) {
	ctx, ctxErr := architectruntime.LoadContext(rootAbs)
	hasCtx := ctxErr == nil

	introKey := "prompt.output.intro"
	workflowKey := "prompt.output.workflow-steps"
	doneKey := "prompt.output.done-items"
	constraintsKey := "prompt.output.constraints-items"
	startKey := "prompt.output.start-items"
	startPhaseKey := "prompt.output.start-items.phase"
	if docsOnly {
		introKey = "prompt.output.docs-only.intro"
		workflowKey = "prompt.output.docs-only.workflow-steps"
		doneKey = "prompt.output.docs-only.done-items"
		constraintsKey = "prompt.output.docs-only.constraints-items"
		startKey = "prompt.output.docs-only.start-items"
		startPhaseKey = "prompt.output.docs-only.start-items.phase"
	}

	fmt.Println(i18n.T(introKey))
	fmt.Println()

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.project-info"))
	fmt.Println(promptSep)
	fmt.Printf("%s %s\n", i18n.T("prompt.output.root-label"), rootAbs)
	if hasCtx {
		if ctx.ProjectName != "" {
			fmt.Printf("%s %s\n", i18n.T("prompt.output.name-label"), ctx.ProjectName)
		}
		if ctx.ProjectIdea != "" {
			fmt.Printf("%s %s\n", i18n.T("prompt.output.idea-label"), ctx.ProjectIdea)
		}
		if stack := promptTechStack(ctx.BackendLanguage, ctx.FrontendFramework, ctx.Database); stack != "" {
			fmt.Printf("%s %s\n", i18n.T("prompt.output.stack-label"), stack)
		}
	} else {
		fmt.Println(i18n.T("prompt.output.no-context"))
	}
	fmt.Println()

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.docs-status"))
	fmt.Println(promptSep)
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
			fmt.Printf("  [%s] %s (%v)\n", i18n.T("prompt.output.missing"), rel, err)
			continue
		}
		if status.Exists {
			fmt.Printf("  [%s] %s\n", i18n.T("prompt.output.exists"), rel)
			anyDoc = true
			if !status.PhaseBased {
				nonPhaseDocs = append(nonPhaseDocs, rel)
			}
		} else {
			fmt.Printf("  [%s] %s\n", i18n.T("prompt.output.missing"), rel)
		}
	}
	if !anyDoc {
		fmt.Println(i18n.T("prompt.output.no-docs"))
	}
	fmt.Println()

	if len(nonPhaseDocs) > 0 {
		fmt.Println(promptSep)
		fmt.Println(i18n.T("prompt.output.phase-warning"))
		fmt.Println(promptSep)
		fmt.Println(i18n.T("prompt.output.phase-warning-items"))
		for _, rel := range nonPhaseDocs {
			fmt.Printf("  - %s\n", rel)
		}
		fmt.Println(i18n.T("prompt.output.phase-warning-action"))
		fmt.Println()
	}

	tasksDir := filepath.Join(rootAbs, "tasks", "queue")
	if info, err := os.Stat(tasksDir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(tasksDir)
		if names := promptMarkdownFiles(entries); len(names) > 0 {
			fmt.Println(promptSep)
			fmt.Println(i18n.T("prompt.output.tasks"))
			fmt.Println(promptSep)
			for _, name := range names {
				fmt.Printf("  tasks/queue/%s\n", name)
			}
			fmt.Println()
		}
	}

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.workflow"))
	fmt.Println(promptSep)
	fmt.Println(i18n.T(workflowKey))
	fmt.Println()

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.done"))
	fmt.Println(promptSep)
	fmt.Println(i18n.T(doneKey))
	fmt.Println()

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.constraints"))
	fmt.Println(promptSep)
	fmt.Println(i18n.T(constraintsKey))
	fmt.Println()

	fmt.Println(promptSep)
	fmt.Println(i18n.T("prompt.output.start"))
	if len(nonPhaseDocs) > 0 {
		fmt.Printf(i18n.T(startPhaseKey), strings.Join(nonPhaseDocs, ", "))
		fmt.Println()
	} else {
		fmt.Println(i18n.T(startKey))
	}
	fmt.Println(promptSep)
}

func promptTechStack(backend, frontend, database string) string {
	var parts []string
	for _, s := range []string{backend, frontend, database} {
		if s != "" && s != "none" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " | ")
}

func promptMarkdownFiles(entries []os.DirEntry) []string {
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, e.Name())
		}
	}
	return names
}

func resolvePromptDocsOnly(explicitDocsOnly bool, docsOnly bool) (bool, error) {
	if explicitDocsOnly || !isInteractive() {
		return docsOnly, nil
	}

	reader := bufio.NewReader(os.Stdin)
	value, err := promptText(reader, i18n.T("prompt.mode.label"), i18n.T("prompt.mode.default"))
	if err != nil {
		return false, err
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "implementation", "implement", "impl", "loop":
		return false, nil
	case "docs", "doc", "docs-only", "doc-review", "review":
		return true, nil
	default:
		return false, fmt.Errorf(i18n.T("prompt.mode.invalid"), value)
	}
}
