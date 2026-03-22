package output

import (
	"fmt"

	"project-generator/internal/generator"
)

func PrintDryRun(plan generator.Plan) {
	fmt.Println("=== DRY RUN ===")
	fmt.Printf("Project root: %s\n", plan.ProjectRoot)

	if len(plan.Questions) > 0 {
		fmt.Println("\nQuestions that would be asked:")
		for _, q := range plan.Questions {
			fmt.Printf("- %s\n", q)
		}
	}

	fmt.Println("\nDirectories to create:")
	for _, dir := range plan.Dirs {
		fmt.Printf("- %s\n", dir)
	}

	fmt.Println("\nFiles to create:")
	for _, file := range plan.Files {
		fmt.Printf("- %s\n", file)
	}

	fmt.Println("\nSkills copy plan:")
	for _, line := range plan.SkillsToCopy {
		fmt.Printf("- %s\n", line)
	}
}

func PrintCreateSummary(result generator.Result) {
	fmt.Println("=== CREATE SUMMARY ===")
	fmt.Printf("Project root: %s\n", result.ProjectRoot)
	if result.BackupPath != "" {
		fmt.Printf("Existing project was backed up to: %s\n", result.BackupPath)
	}

	fmt.Printf("Created directories: %d\n", len(result.CreatedDirs))
	for _, dir := range result.CreatedDirs {
		fmt.Printf("- %s\n", dir)
	}

	fmt.Printf("Created files: %d\n", len(result.CreatedFiles))
	for _, file := range result.CreatedFiles {
		fmt.Printf("- %s\n", file)
	}

	if len(result.SkillCopies) > 0 {
		fmt.Println("Skill copy results:")
		for _, skill := range result.SkillCopies {
			if skill.Success {
				fmt.Printf("- [OK] %s (%s -> %s)\n", skill.Name, skill.Source, skill.Dest)
			} else {
				fmt.Printf("- [FAIL] %s: %s\n", skill.Name, skill.Error)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("- %s\n", warning)
		}
	}
}
