package apa

import (
	"project-generator/internal/i18n"

	"github.com/spf13/cobra"
)

func Execute() error {
	i18n.Init()

	root := &cobra.Command{
		Use:           "apa",
		Short:         i18n.T("root.short"),
		Long:          i18n.T("root.long"),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Register --lang as a persistent flag so cobra does not reject it.
	// The actual value is read from os.Args by i18n.Init() before commands are built.
	root.PersistentFlags().String("lang", "en", "output language (en, zh-TW, zh-CN, ja, ko, de, es, fr)")

	root.CompletionOptions.DisableDefaultCmd = true
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.AddCommand(newInitCmd())
	root.AddCommand(newIterateCmd())
	root.AddCommand(newListSkillsCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newVersionCmd())

	return root.Execute()
}
