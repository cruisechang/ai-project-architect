package apa

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"project-generator/internal/config"
	"project-generator/internal/i18n"
)

func newListSkillsCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "list-skills",
		Short: i18n.T("list-skills.short"),
		Long:  i18n.T("list-skills.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				defaultPath, err := defaultLocalSkillsPath()
				if err != nil {
					return err
				}
				if isInteractive() {
					reader := bufio.NewReader(os.Stdin)
					val, err := promptText(reader, "Skills path", defaultPath)
					if err != nil {
						return err
					}
					if strings.TrimSpace(val) != "" {
						path = val
					} else {
						path = defaultPath
					}
				} else {
					path = defaultPath
				}
			}

			absPath, err := config.ExpandPath(path)
			if err != nil {
				return err
			}

			info, err := os.Stat(absPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("skills path does not exist: %s", absPath)
				}
				return err
			}
			if !info.IsDir() {
				return fmt.Errorf("skills path is not a directory: %s", absPath)
			}

			entries, err := os.ReadDir(absPath)
			if err != nil {
				return err
			}

			skills := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					skills = append(skills, entry.Name())
				}
			}
			sort.Strings(skills)

			fmt.Printf("Skills path: %s\n", absPath)
			if len(skills) == 0 {
				fmt.Println("No skills found.")
				return nil
			}

			for _, skill := range skills {
				fmt.Printf("- %s\n", skill)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&path, "path", "", i18n.T("list-skills.flag.path"))
	return cmd
}
