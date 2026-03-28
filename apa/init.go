package apa

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"project-generator/internal/architect/model"
	architectgenerators "project-generator/internal/architect/modules/generators"
	"project-generator/internal/architect/modules/planner"
	architectruntime "project-generator/internal/architect/runtime"
	"project-generator/internal/cli"
	"project-generator/internal/config"
	"project-generator/internal/generator"
	"project-generator/internal/i18n"
	"project-generator/internal/output"
)

func newInitCmd() *cobra.Command {
	var opts config.CreateOptions
	var skillsRaw string
	var idea string
	var ideaFile string

	cmd := &cobra.Command{
		Use:   "init",
		Short: i18n.T("init.short"),
		Long:  i18n.T("init.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if idea, err = resolveIdea(idea, ideaFile); err != nil {
				return err
			}
			opts.Idea = idea

			interactive := isInteractive()

			// Infer tech stack from idea and pre-fill opts
			hasIdea := strings.TrimSpace(opts.Idea) != ""
			var inferredCtx model.ProjectContext
			if hasIdea {
				inferredCtx = planner.New().Infer(opts.Idea, opts.Name)
				mapContextToOpts(inferredCtx, &opts)
			}

			opts.SelectedSkills = config.ParseSkillsCSV(skillsRaw)
			if len(opts.SelectedSkills) > 0 && strings.TrimSpace(opts.GlobalSkillsPath) == "" {
				defaultSkillsPath, err := defaultLocalSkillsPath()
				if err != nil {
					return err
				}
				opts.GlobalSkillsPath = defaultSkillsPath
			}
			opts.Normalize()

			if err := opts.ValidateKnownValues(); err != nil {
				return err
			}

			// Run wizard for missing required fields, and also ensure --type/--agent
			// are part of the interactive flow when not explicitly provided.
			missing := opts.MissingRequired()
			runWizard := len(missing) > 0
			if interactive && (!cmd.Flags().Changed("type") || !cmd.Flags().Changed("agent")) {
				runWizard = true
			}
			if runWizard {
				if !interactive {
					return fmt.Errorf("missing required options: %v (provide via flags or run interactively)", missing)
				}
				filled, err := cli.CollectCreateOptions(opts, false)
				if err != nil {
					if errors.Is(err, cli.ErrAborted) {
						fmt.Println("Init cancelled.")
						return nil
					}
					return err
				}
				opts = filled
			}

			hasIdea = strings.TrimSpace(opts.Idea) != ""
			if hasIdea {
				inferredCtx = planner.New().Infer(opts.Idea, opts.Name)
			}

			opts.EnsureDefaults()
			opts.Normalize()

			if err := opts.ValidateForCreate(); err != nil {
				return err
			}

			if opts.DryRun {
				plan, err := generator.BuildPlan(opts, nil)
				if err != nil {
					return err
				}
				output.PrintDryRun(plan)
				return nil
			}

			// Confirm overwrite if target directory already exists
			if err := maybeConfirmOverwrite(&opts); err != nil {
				if errors.Is(err, cli.ErrAborted) {
					fmt.Println("Init cancelled.")
					return nil
				}
				return err
			}

			fmt.Println("[1/4] infer tech stack from idea")

			// [2/4] Generate code scaffold (creates dirs, code, and template docs)
			fmt.Println("[2/4] generate code scaffold")
			result, err := generator.CreateProject(opts)
			if err != nil {
				return err
			}

			// [3/4] Save context + overwrite template docs with idea-based content
			fmt.Println("[3/4] generate design docs")

			ctxToSave := inferredCtx
			if !hasIdea {
				ctxToSave = model.ProjectContext{
					ProjectName:       opts.Name,
					BackendLanguage:   opts.BackendType,
					FrontendFramework: opts.FrontendType,
				}
			}
			ctxToSave.ProjectName = opts.Name
			ctxToSave.EnsureDefaults()

			if saveErr := architectruntime.SaveContext(result.ProjectRoot, ctxToSave); saveErr != nil {
				fmt.Fprintf(os.Stderr, "warning: could not save context: %v\n", saveErr)
			}

			archDocCount := 0
			if hasIdea {
				archArtifacts, archErr := architectgenerators.ResolveTarget("architecture")
				if archErr != nil {
					return archErr
				}
				archFiles, archErr := writeArtifacts(result.ProjectRoot, ctxToSave, archArtifacts)
				if archErr != nil {
					return archErr
				}
				docArtifacts, docErr := architectgenerators.ResolveTarget("docs")
				if docErr != nil {
					return docErr
				}
				docFiles, docErr := writeArtifacts(result.ProjectRoot, ctxToSave, docArtifacts)
				if docErr != nil {
					return docErr
				}
				archDocCount = len(archFiles) + len(docFiles)
			}

			// [4/4] Done
			fmt.Println("[4/4] done")
			fmt.Println()
			fmt.Println("=== INIT SUMMARY ===")
			fmt.Printf("Project root:   %s\n", result.ProjectRoot)
			if result.BackupPath != "" {
				fmt.Printf("Backup:         %s\n", result.BackupPath)
			}
			fmt.Printf("Scaffold dirs:  %d\n", len(result.CreatedDirs))
			fmt.Printf("Scaffold files: %d\n", len(result.CreatedFiles))
			if hasIdea {
				fmt.Printf("Design docs:    %d files\n", archDocCount)
				fmt.Printf("Context:        %s/.architect/context.json\n", result.ProjectRoot)
			}
			for _, skill := range result.SkillCopies {
				if !skill.Success {
					fmt.Printf("- [FAIL] %s: %s\n", skill.Name, skill.Error)
				}
			}
			if failed := result.FailedSkillCopies(); failed > 0 {
				return fmt.Errorf("init completed with %d failed skill copy operations", failed)
			}
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Printf("  cd %s\n", result.ProjectRoot)
			fmt.Println("  make test     # start TDD: red → green → refactor")
			fmt.Println("  apa prompt   # generate AI prompt for continuous iteration")
			return nil
		},
	}

	cmd.Flags().StringVar(&idea, "idea", "", i18n.T("init.flag.idea"))
	cmd.Flags().StringVar(&ideaFile, "idea-file", "", i18n.T("init.flag.idea-file"))
	cmd.Flags().StringVar(&opts.Name, "name", "", i18n.T("init.flag.name"))
	cmd.Flags().StringVar(&opts.ParentPath, "path", "", i18n.T("init.flag.path"))
	cmd.Flags().StringVar(&opts.AIFeature, "ai-feature", "", i18n.T("init.flag.ai-feature"))
	cmd.Flags().StringVar(&opts.AIAgent, "agent", "", i18n.T("init.flag.agent"))
	cmd.Flags().StringVar(&opts.BackendType, "backend", "", i18n.T("init.flag.backend"))
	cmd.Flags().StringVar(&opts.FrontendType, "frontend", "", i18n.T("init.flag.frontend"))
	cmd.Flags().StringVar(&opts.Architecture, "type", "", i18n.T("init.flag.type"))
	cmd.Flags().StringVar(&opts.TechStack, "stack", "", i18n.T("init.flag.stack"))
	cmd.Flags().StringVar(&opts.DocsType, "docs", "", i18n.T("init.flag.docs"))
	cmd.Flags().StringVar(&opts.UnitTest, "unit-test", "", i18n.T("init.flag.unit-test"))
	cmd.Flags().StringVar(&opts.APITest, "api-test", "", i18n.T("init.flag.api-test"))
	cmd.Flags().StringVar(&opts.IntegrationTest, "integration-test", "", i18n.T("init.flag.integration-test"))
	cmd.Flags().StringVar(&opts.E2ETest, "e2e-test", "", i18n.T("init.flag.e2e-test"))
	cmd.Flags().StringVar(&opts.DockerCompose, "docker-compose", "", i18n.T("init.flag.docker-compose"))
	cmd.Flags().StringVar(&skillsRaw, "skills", "", i18n.T("init.flag.skills"))
	cmd.Flags().StringVar(&opts.GlobalSkillsPath, "skills-path", "", i18n.T("init.flag.skills-path"))
	cmd.Flags().StringVar(&opts.GlobalSkillsPath, "global-skills", "", "deprecated: use --skills-path")
	_ = cmd.Flags().MarkHidden("global-skills")
	cmd.Flags().StringVar(&opts.Description, "description", "", i18n.T("init.flag.description"))
	cmd.Flags().BoolVar(&opts.Overwrite, "force", false, i18n.T("init.flag.force"))
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, i18n.T("init.flag.dry-run"))

	return cmd
}

// mapContextToOpts fills in empty CreateOptions fields from an inferred ProjectContext.
// Explicit flags (non-empty values) always take priority over inferred values.
// This includes sensible defaults for TechStack, DockerCompose, and ProjectType so
// that non-interactive invocations with --idea do not require the full wizard.
func mapContextToOpts(ctx model.ProjectContext, opts *config.CreateOptions) {
	if opts.BackendType == "" {
		opts.BackendType = ctx.BackendLanguage
	}
	if opts.FrontendType == "" {
		opts.FrontendType = ctx.FrontendFramework
	}
	if opts.Architecture == "" {
		opts.Architecture = inferArchitectureFromContext(ctx)
	}
	if opts.TechStack == "" {
		opts.TechStack = buildTechStack(opts.Architecture, opts.BackendType, opts.FrontendType)
	}
	if opts.DockerCompose == "" {
		if strings.Contains(strings.ToLower(ctx.Deployment), "docker") {
			opts.DockerCompose = "yes"
		} else {
			opts.DockerCompose = "no"
		}
	}
	if opts.ProjectType == "" {
		opts.ProjectType = config.InferProjectTypeFromArchitecture(opts.Architecture)
	}
}

func inferArchitectureFromContext(ctx model.ProjectContext) string {
	frontend := strings.ToLower(ctx.FrontendFramework)
	switch ctx.ProjectType {
	case "backend":
		return "server"
	case "frontend":
		if frontend == "react-native" || frontend == "flutter" || frontend == "swiftui" || frontend == "kotlin" || frontend == "android" || frontend == "ios" {
			return "mobile-app"
		}
		return "web-app"
	case "fullstack":
		if frontend == "react-native" || frontend == "flutter" || frontend == "swiftui" || frontend == "kotlin" || frontend == "android" || frontend == "ios" {
			return "mobile-app-server"
		}
		return "web-app-server"
	}
	return ""
}

func maybeConfirmOverwrite(opts *config.CreateOptions) error {
	parent := opts.ParentPath
	if parent == "" {
		parent = "."
	}
	parentAbs, err := config.ExpandPath(parent)
	if err != nil {
		return err
	}
	projectRoot := filepath.Join(parentAbs, opts.Name)
	info, statErr := os.Stat(projectRoot)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return nil
		}
		return statErr
	}
	if !info.IsDir() {
		return fmt.Errorf("target path exists and is not a directory: %s", projectRoot)
	}
	if opts.Overwrite {
		return nil
	}
	yes, err := cli.ConfirmOverwrite(projectRoot)
	if err != nil {
		return err
	}
	if !yes {
		return cli.ErrAborted
	}
	opts.Overwrite = true
	return nil
}

func buildTechStack(architecture, backend, frontend string) string {
	parts := []string{}
	if architecture != "" {
		parts = append(parts, architecture)
	}
	if backend != "" && backend != "none" {
		parts = append(parts, backend)
	}
	if frontend != "" && frontend != "none" {
		parts = append(parts, frontend)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " | ")
}
